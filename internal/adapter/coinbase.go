package adapter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/marketfeedmodels"
	"github.com/Ovikx/market-data-recorder/internal/strategyprimitives"
)

type blshCoinbaseAdapter struct{}

func NewCoinbaseAdapter() *blshCoinbaseAdapter {
	return &blshCoinbaseAdapter{}
}

func (a *blshCoinbaseAdapter) Reroute(data []byte, ticks chan Tick) error {
	var channelNameStruct marketfeedmodels.CoinbaseChannelName
	err := json.Unmarshal(data, &channelNameStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	switch channelNameStruct.Channel {
	// case "l2_data":
	// 	var msg marketfeedmodels.CoinbaseL2Message
	// 	if err := json.Unmarshal(data, &msg); err != nil {
	// 		return fmt.Errorf("failed to unmarshal l2_data message: %v", err)
	// 	}

	// 	// Send all the updates through the chanenl
	// 	for _, e := range msg.Events {
	// 		if e.Type == "update" {
	// 			for _, u := range e.Updates {
	// 				var side uint8 = 0
	// 				if u.Side == "bid" {
	// 					side = 1
	// 				} else if u.Side != "offer" {
	// 					return fmt.Errorf("expected side 'bid' or 'offer', instead received: %s", u.Side)
	// 				}

	// 				priceLevel, err := strconv.ParseFloat(u.PriceLevel, 64)
	// 				if err != nil {
	// 					return fmt.Errorf("unable to parse '%s' into a float64: %v", u.PriceLevel, err)
	// 				}

	// 				newQuantity, err := strconv.ParseFloat(u.NewQuantity, 64)
	// 				if err != nil {
	// 					return fmt.Errorf("unable to parse '%s' into a float64: %v", u.NewQuantity, err)
	// 				}

	// 				time, err := time.Parse(time.RFC3339Nano, u.EventTime)
	// 				if err != nil {
	// 					return fmt.Errorf("unable to parse time %s: %v", u.EventTime, err)
	// 				}
	// 				a.Orders() <- strategyprimitives.NewOrder(e.ProductID, side, priceLevel, newQuantity, time)
	// 			}
	// 		}
	// 	}
	// case "market_trades":
	// 	var tradesMsg marketfeedmodels.CoinbaseTradesMessage
	// 	if err := json.Unmarshal(data, &tradesMsg); err != nil {
	// 		return fmt.Errorf("failed to unmarshal market_trades message: %v", err)
	// 	}

	// 	// Send all the updates through the chanenl
	// 	for _, e := range tradesMsg.Events {
	// 		for _, t := range e.Trades {
	// 			var side uint8 = 0
	// 			if t.Side == "BUY" {
	// 				side = 1
	// 			} else if t.Side != "SELL" {
	// 				return fmt.Errorf("expected side 'BUY' or 'SELL', instead received: %s", t.Side)
	// 			}

	// 			price, err := strconv.ParseFloat(t.Price, 64)
	// 			if err != nil {
	// 				return fmt.Errorf("unable to parse '%s' into a float64", t.Price)
	// 			}

	// 			size, err := strconv.ParseFloat(t.Size, 64)
	// 			if err != nil {
	// 				return fmt.Errorf("unable to parse '%s' into a float64", t.Size)
	// 			}

	// 			time, err := time.Parse(time.RFC3339Nano, t.Time)
	// 			if err != nil {
	// 				return fmt.Errorf("unable to parse time '%s': %v", t.Time, err)
	// 			}
	// 			a.Trades() <- strategyprimitives.NewTrade(t.ProductID, side, price, size, time)
	// 		}
	// 	}
	case "ticker":
		var msg marketfeedmodels.CoinbaseTickerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal ticker message: %v", err)
		}

		time, err := time.Parse(time.RFC3339Nano, msg.Timestamp)
		if err != nil {
			return fmt.Errorf("unable to parse time '%s': %v", msg.Timestamp, err)
		}

		for _, e := range msg.Events {
			if e.Type == "update" {
				for _, t := range e.Tickers {
					price, err := strconv.ParseFloat(t.Price, 64)
					if err != nil {
						return fmt.Errorf("unable to parse '%s' into a float64", t.Price)
					}
					ticks <- strategyprimitives.NewTick(t.ProductID, price, time)
				}
			}
		}
	}

	return nil
}
