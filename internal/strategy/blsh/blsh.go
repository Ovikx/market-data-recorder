package blsh

import (
	"log"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/datastructure/timewindow"
	"github.com/Ovikx/market-data-recorder/internal/strategy"
)

type blsh struct {
	trades   chan strategy.Trade
	orders   chan strategy.Order
	ticks    chan strategy.Tick
	executor strategy.Executor
}

func New() *blsh {
	return &blsh{trades: make(chan strategy.Trade), orders: make(chan strategy.Order), ticks: make(chan strategy.Tick)}
}

func (s *blsh) Trades() chan strategy.Trade {
	return s.trades
}

func (s *blsh) Orders() chan strategy.Order {
	return s.orders
}

func (s *blsh) Ticks() chan strategy.Tick {
	return s.ticks
}

func (s *blsh) Start() error {
	ticks := timewindow.New[strategy.Tick](5 * time.Second)
	var sumPrice float64
	for {
		select {
		case <-s.trades:
			// log.Println("TRADES:", msg)
		case <-s.orders:
			// if msg.Quantity() > 0 {
			// 	removed, ok := orders.Insert(msg, msg.Timestamp())
			// 	if ok {
			// 		sumPrice += msg.Price()
			// 	}
			// 	for _, r := range removed {
			// 		sumPrice -= r.Price()
			// 	}
			// }

			// log.Printf("%v, ORDERS: %v\n", sumPrice/float64(orders.Length()), msg)
		case msg := <-s.ticks:
			removed, ok := ticks.Insert(msg, msg.Timestamp())
			if ok {
				sumPrice += msg.Price()
			}
			for _, r := range removed {
				sumPrice -= r.Price()
			}
			log.Printf("TICKER: %v\n", sumPrice/float64(ticks.Length()))
		}
	}
}
