package strategyprimitives

import (
	"fmt"
	"time"
)

type tick struct {
	symbol    string
	price     float64
	timestamp time.Time
}

func NewTick(symbol string, price float64, timestamp time.Time) *tick {
	return &tick{symbol: symbol, price: price, timestamp: timestamp}
}

func (t *tick) Symbol() string {
	return t.symbol
}

func (t *tick) Price() float64 {
	return t.price
}

func (t *tick) Timestamp() time.Time {
	return t.timestamp
}

func (o *tick) String() string {
	return fmt.Sprintf("%s %f (%v)", o.symbol, o.price, o.timestamp)
}
