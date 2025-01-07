package strategy

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

// Executes trades
type Executor interface {
	MarketBuy(symbol string, quoteSize float64) error
	MarketSell(symbol string, baseSize float64) error
	LimitBuy(symbol string, size float64, limitPrice float64) error
	StopLimitSell(symbol string, size float64, limitPrice float64) error
}
