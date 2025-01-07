package blshadapters

import "github.com/Ovikx/market-data-recorder/internal/strategy"

type blsh interface {
	Trades() chan strategy.Trade
	Orders() chan strategy.Order
	Ticks() chan strategy.Tick
}
