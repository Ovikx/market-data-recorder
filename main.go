package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/adapter"
	"github.com/Ovikx/market-data-recorder/internal/dbwriter"
	"github.com/Ovikx/market-data-recorder/internal/jwtgen"
	"github.com/Ovikx/market-data-recorder/internal/marketfeed"
	"github.com/Ovikx/market-data-recorder/internal/profileloader"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type strategyAdapter interface {
	Reroute(data []byte, ticks chan adapter.Tick, orders chan adapter.Order) error
}

// Calculates the number of seconds to wait on the n-th reconnect retry
func backoff(n int) int {
	return 1 + n*n
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
	liveBool := flag.Bool("l", false, "whether received market data should be logged")
	tickBool := flag.Bool("t", false, "whether to record ticks")
	orderBool := flag.Bool("o", false, "whether to record orders")
	flag.Parse()

	if *profilePathStr == "" {
		log.Fatal("no profile provided")
	}

	if !*liveBool {
		log.Printf("WARNING: NOT RECORDING DATA")
	}

	// Connect to the DB and start listening
	var ticks chan adapter.Tick
	var orders chan adapter.Order
	if *tickBool {
		ticks = make(chan adapter.Tick)
	}
	if *orderBool {
		orders = make(chan adapter.Order)
	}
	dbwriter, err := dbwriter.New(os.Getenv("POSTGRES_URL"), *liveBool, ticks, orders)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	go dbwriter.Record("ticks", "orders")
	defer dbwriter.Close()

	// Load the profile
	profile, err := profileloader.FromFile(*profilePathStr, "schemas/profile-schema.json")
	if err != nil {
		log.Fatalf("failed to load profile %v: %v", *profilePathStr, err)
	}

	// Adapter (for sending the right WS messages to the right Go channels)
	var strategyAdapter strategyAdapter

	// Connect to market feed (websocket streams)
	var marketFeedConns []*websocket.Conn
	var reconnectFuncs []func() (*websocket.Conn, error)

	switch profile.Provider {
	case "coinbase":
		marketFeedConns, reconnectFuncs, err = marketfeed.ConnectToCoinbaseMarketFeed(profile.WSUrl, jwtgen.CoinbaseJWT, profile.Symbols)
		strategyAdapter = adapter.NewCoinbaseAdapter()
	case "alpaca":
		marketFeedConns, err = marketfeed.ConnectToAlpacaMarketFeed(profile.WSUrl, profile.Symbols)
		strategyAdapter = adapter.NewCoinbaseAdapter()
	case "kraken":
		marketFeedConns, err = marketfeed.ConnectToKrakenMarketFeed(profile.WSUrl, profile.Symbols)
		strategyAdapter = adapter.NewKrakenAdapter()
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
	numRetries := make([]int, len(marketFeedConns))

	numLive.Store(int32(len(marketFeedConns)))
	for i := 0; i < len(marketFeedConns); i++ {
		go func() {
			// Close done channel once there are no more live connections
			defer func() {
				numLive.Add(-1)
				if numLive.Load() == 0 {
					close(done)
				}
			}()
			for {
				_, message, err := marketFeedConns[i].ReadMessage()
				if err != nil {
					log.Printf("read error on conn %d: %v", i, err)

					// Normal close
					if strings.Contains(err.Error(), "close 1000") {
						return
					}

					marketFeedConns[i].WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					err = fmt.Errorf("dummy error")
					var newConn *websocket.Conn
					for err != nil {
						// Fatalf after too many attempts
						if numRetries[i] > 3 {
							log.Fatalf("failed to reconnect on conn %d too many times", i)
						}

						// Wait
						<-time.After(time.Duration(backoff(numRetries[i])) * time.Second)

						// Attempt to reconnect
						log.Printf("attempting to reconnect on conn %d (retry %d)", i, numRetries[i])
						newConn, err = reconnectFuncs[i]()
						if err != nil {
							log.Printf("failed to reconnect on conn %d: %v", i, err)
							numRetries[i]++
						}
					}

					marketFeedConns[i] = newConn
					numRetries[i] = 0
					log.Printf("reconnected on conn %d", i)
				} else {
					err = strategyAdapter.Reroute(message, dbwriter.Ticks(), dbwriter.Orders())
					if err != nil {
						log.Println("failed to reroute market data:", err)
						return
					}
				}

			}
		}()
	}

	for {
		select {
		case e := <-dbwriter.Errors():
			log.Printf("db error: %v", e)
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
