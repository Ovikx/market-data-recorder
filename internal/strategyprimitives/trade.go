package strategyprimitives

import (
	"fmt"
	"time"
)

type trade struct {
	symbol    string
	side      uint8
	price     float64
	size      float64
	timestamp time.Time
}

func NewTrade(symbol string, side uint8, price, quantity float64, timestamp time.Time) *trade {
	return &trade{
		symbol:    symbol,
		side:      side,
		price:     price,
		size:      quantity,
		timestamp: timestamp,
	}
}

func (t *trade) Symbol() string {
	return t.symbol
}

func (t *trade) Side() uint8 {
	return t.side
}

func (t *trade) Price() float64 {
	return t.price
}

func (t *trade) Size() float64 {
	return t.size
}

func (t *trade) Timestamp() time.Time {
	return t.timestamp
}

func (o *trade) String() string {
	side := "S"
	if o.side == 1 {
		side = "B"
	}
	return fmt.Sprintf("%s %s @ %f (qty. %f) (%v)", o.symbol, side, o.price, o.size, o.timestamp)
}
