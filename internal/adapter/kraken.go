package adapter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/marketfeedmodels"
	"github.com/Ovikx/market-data-recorder/internal/strategyprimitives"
)

type blshKrakenAdapter struct{}

func NewKrakenAdapter() *blshKrakenAdapter {
	return &blshKrakenAdapter{}
}

func (a *blshKrakenAdapter) Reroute(data []byte, ticks chan Tick, orders chan Order, trades chan Trade) error {
	var channelNameStruct marketfeedmodels.KrakenChannelName
	err := json.Unmarshal(data, &channelNameStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	switch channelNameStruct.Channel {
	case "book":
		var msg marketfeedmodels.KrakenBookMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal book message: %v", err)
		}

		// Ignore snapshots
		if msg.Type == "snapshot" {
			return nil
		}

		// Kraken expects a singular object to be in the `Data` array
		if len(msg.Data) != 1 {
			return fmt.Errorf("unexpected number of items in book message data: expected 1, got %d", len(msg.Data))
		}

		data := msg.Data[0]

		time, err := time.Parse(time.RFC3339Nano, data.Timestamp)
		if err != nil {
			return fmt.Errorf("unable to parse time %s: %v", data.Timestamp, err)
		}

		for _, ask := range data.Asks {
			orders <- strategyprimitives.NewOrder(data.Symbol, 0, ask.Price, ask.Qty, time)
		}
		for _, bid := range data.Bids {
			orders <- strategyprimitives.NewOrder(data.Symbol, 1, bid.Price, bid.Qty, time)
		}
	case "trade":
		var msg marketfeedmodels.KrakenTradeMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal trade message: %v", err)
		}

		// Ignore snapshots
		if msg.Type == "snapshot" {
			return nil
		}

		for _, trade := range msg.Data {
			var side uint8 = 0
			if trade.Side == "buy" {
				side = 1
			} else if trade.Side != "sell" {
				return fmt.Errorf("expected side 'buy' or 'sell', instead received: %s", trade.Side)
			}

			time, err := time.Parse(time.RFC3339Nano, trade.Timestamp)
			if err != nil {
				return fmt.Errorf("unable to parse time %s: %v", trade.Timestamp, err)
			}
			trades <- strategyprimitives.NewTrade(trade.Symbol, side, trade.Price, trade.Qty, time)
		}
	case "ticker":
		var msg marketfeedmodels.KrakenTickerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal ticker message: %v", err)
		}

		// Ignore snapshots
		if msg.Type == "snapshot" {
			return nil
		}

		// Kraken expects a singular object to be in the `Data` array
		if len(msg.Data) != 1 {
			return fmt.Errorf("unexpected number of items in book message data: expected 1, got %d", len(msg.Data))
		}

		data := msg.Data[0]

		ticks <- strategyprimitives.NewTick(data.Symbol, data.Last, time.Now())
	}

	return nil
}
