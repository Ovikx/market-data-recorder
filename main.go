package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/jwtgen"
	"github.com/Ovikx/market-data-recorder/internal/marketfeed"
	"github.com/Ovikx/market-data-recorder/internal/profileloader"
	"github.com/Ovikx/market-data-recorder/internal/strategy/blsh"
	"github.com/Ovikx/market-data-recorder/internal/strategyadapters/blshadapters"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type strategyAdapter interface {
	Reroute(data []byte) error
}

type strategy interface {
	Start() error
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Load the .env file variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// Parse cmd line args
	profilePathStr := flag.String("p", "", "path of the profile JSON file to use")
	_ = flag.Bool("record", false, "Whether to record received market data")
	_ = flag.Bool("live", false, "Whether to trade using real money") // TODO: pass the returned val to whatever makes trades
	flag.Parse()

	if *profilePathStr == "" {
		log.Fatal("no profile provided")
	}

	// Load the profile
	profile, err := profileloader.FromFile(*profilePathStr, "schemas/profile-schema.json")
	if err != nil {
		log.Fatalf("failed to load profile %v: %v", *profilePathStr, err)
	}

	// Start the strategy
	// TODO: make dynamic strategy selection (from cmd line maybe idk)
	s := blsh.New()
	go s.Start()

	// Adapter (for sending the right WS messages to the right Go channels)
	var strategyAdapter strategyAdapter

	// Connect to market feed (websocket streams)
	var marketFeedConns []*websocket.Conn

	switch profile.Provider {
	case "coinbase":
		marketFeedConns, err = marketfeed.ConnectToCoinbaseMarketFeed(profile.WSUrl, jwtgen.CoinbaseJWT, profile.Symbols)
		strategyAdapter = blshadapters.NewCoinbaseAdapter(s)
	case "alpaca":
		marketFeedConns, err = marketfeed.ConnectToAlpacaMarketFeed(profile.WSUrl, profile.Symbols)
		strategyAdapter = blshadapters.NewCoinbaseAdapter(s)
	case "kraken":
		marketFeedConns, err = marketfeed.ConnectToKrakenMarketFeed(profile.WSUrl, profile.Symbols)
		strategyAdapter = blshadapters.NewKrakenAdapter(s)
	}

	if err != nil {
		log.Fatalf("error connecting to market feed: %v", err)
	}
	defer func() {
		for _, conn := range marketFeedConns {
			conn.Close()
		}
	}()

	// Continuously read incoming messages from all channels
	done := make(chan struct{})
	numLive := atomic.Int32{}

	numLive.Store(int32(len(marketFeedConns)))
	for i, conn := range marketFeedConns {
		go func() {
			// Close done channel once there are no more live connections
			defer func() {
				numLive.Add(-1)
				if numLive.Load() == 0 {
					close(done)
				}
			}()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("read error on conn %d: %v", i, err)
					return
				}

				err = strategyAdapter.Reroute(message)
				if err != nil {
					log.Println("failed to reroute market data:", err)
					return
				}
			}
		}()
	}

	for {
		select {
		case <-done:
			fmt.Println("Done")
			return
		case <-interrupt: // CTRL-C
			log.Println("Interrupt")
			for i, conn := range marketFeedConns {
				if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
					log.Printf("write error on close conn %d: %v", i, err)
					select {
					case <-done:
					case <-time.After(time.Second):
					}
					return
				}
			}
		}

	}

}
