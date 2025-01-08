package adapter

import "time"

type Trade interface {
	Symbol() string
	// 0 = Sell ; 1 = Buy
	Side() uint8
	Price() float64
	Size() float64
	Timestamp() time.Time
}

type Order interface {
	Symbol() string
	// 0 = Sell ; 1 = Buy
	Side() uint8
	Price() float64
	Size() float64
	Timestamp() time.Time
}

type Tick interface {
	Symbol() string
	Price() float64
	Timestamp() time.Time
}
