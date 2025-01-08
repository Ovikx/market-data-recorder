package main_test

import (
	"context"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/adapter"
	"github.com/Ovikx/market-data-recorder/internal/dbwriter"
	"github.com/Ovikx/market-data-recorder/internal/jwtgen"
	"github.com/Ovikx/market-data-recorder/internal/marketfeed"
	"github.com/joho/godotenv"
)

type rerouter interface {
	Reroute(data []byte, ticks chan adapter.Tick) error
}

func TestCoinbaseTickWrites(t *testing.T) {
	// Load the .env file variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to the DB and start listening
	dbwriter, err := dbwriter.New(os.Getenv("POSTGRES_URL"), true)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	go dbwriter.Record("test_ticks")
	defer dbwriter.Close()

	// Create test table
	dbwriter.Pool().Exec(context.Background(), "DROP TABLE test_ticks")
	_, err = dbwriter.Pool().Exec(context.Background(), "CREATE TABLE test_ticks (LIKE ticks INCLUDING ALL)")
	if err != nil {
		t.Errorf("failed to create dummy db")
	}
	defer func() {
		dbwriter.Pool().Exec(context.Background(), "DROP TABLE test_ticks")
	}()

	marketFeedConns, _, err := marketfeed.ConnectToCoinbaseMarketFeed("wss://advanced-trade-ws.coinbase.com", jwtgen.CoinbaseJWT, []string{"BTC-USD"})
	strategyAdapter := adapter.NewCoinbaseAdapter()

	done := make(chan struct{})
	numLive := atomic.Int32{}
	numAddedExpected := atomic.Int32{}

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

				err = strategyAdapter.Reroute(message, dbwriter.Ticks())
				if err != nil {
					t.Errorf("failed to reroute market data: %v", err)
				}
				if strings.Contains(string(message), "ticker") && strings.Contains(string(message), "update") {
					numAddedExpected.Add(1)
				}
			}
		}()
	}

	go func() {
		<-time.After(3 * time.Second)
		close(done)
	}()

	for {
		select {
		case e := <-dbwriter.Errors():
			t.Errorf("db error: %v", e)
		case <-done:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rows, err := dbwriter.Pool().Query(ctx, "SELECT id FROM test_ticks")
			if err != nil {
				t.Error(err)
			}
			defer rows.Close()
			numAddedActual := 0
			for rows.Next() {
				numAddedActual++
			}
			if numAddedActual != int(numAddedExpected.Load()) {
				t.Errorf("expected %d new rows but only found %d", numAddedExpected.Load(), numAddedActual)
			}
			return
		}
	}
}

func TestKrakenTickWrites(t *testing.T) {
	// Load the .env file variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to the DB and start listening
	dbwriter, err := dbwriter.New(os.Getenv("POSTGRES_URL"), true)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	go dbwriter.Record("test_ticks")
	defer dbwriter.Close()

	// Create test table
	dbwriter.Pool().Exec(context.Background(), "DROP TABLE test_ticks")
	_, err = dbwriter.Pool().Exec(context.Background(), "CREATE TABLE test_ticks (LIKE ticks INCLUDING ALL)")
	if err != nil {
		t.Errorf("failed to create dummy db")
	}
	defer func() {
		dbwriter.Pool().Exec(context.Background(), "DROP TABLE test_ticks")
	}()

	marketFeedConns, err := marketfeed.ConnectToKrakenMarketFeed("wss://ws.kraken.com/v2", []string{"BTC-USD"})
	strategyAdapter := adapter.NewKrakenAdapter()

	done := make(chan struct{})
	numLive := atomic.Int32{}
	numAddedExpected := atomic.Int32{}

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

				err = strategyAdapter.Reroute(message, dbwriter.Ticks())
				if err != nil {
					t.Errorf("failed to reroute market data: %v", err)
				}
				if strings.Contains(string(message), "ticker") && strings.Contains(string(message), "update") {
					numAddedExpected.Add(1)
				}
			}
		}()
	}

	go func() {
		<-time.After(3 * time.Second)
		close(done)
	}()

	for {
		select {
		case e := <-dbwriter.Errors():
			t.Errorf("db error: %v", e)
		case <-done:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rows, err := dbwriter.Pool().Query(ctx, "SELECT id FROM test_ticks")
			if err != nil {
				t.Error(err)
			}
			defer rows.Close()
			numAddedActual := 0
			for rows.Next() {
				log.Println(rows.Values())
				numAddedActual++
			}
			if numAddedActual != int(numAddedExpected.Load()) {
				t.Errorf("expected %d new rows but only found %d", numAddedExpected.Load(), numAddedActual)
			}
			return
		}
	}
}
