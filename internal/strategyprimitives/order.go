package strategyprimitives

import (
	"fmt"
	"time"
)

type order struct {
	symbol    string
	side      uint8
	price     float64
	size      float64
	timestamp time.Time
}

func NewOrder(symbol string, side uint8, price float64, quantity float64, timestamp time.Time) *order {
	return &order{
		symbol:    symbol,
		side:      side,
		price:     price,
		size:      quantity,
		timestamp: timestamp,
	}
}

func (o *order) Symbol() string {
	return o.symbol
}

func (o *order) Side() uint8 {
	return o.side
}

func (o *order) Price() float64 {
	return o.price
}

func (o *order) Size() float64 {
	return o.size
}

func (o *order) Timestamp() time.Time {
	return o.timestamp
}

func (o *order) String() string {
	side := "A"
	if o.side == 1 {
		side = "B"
	}
	return fmt.Sprintf("%s %s @ %f (qty. %f) (%v)", o.symbol, side, o.price, o.size, o.timestamp)
}
