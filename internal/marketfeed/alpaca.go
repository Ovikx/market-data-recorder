package marketfeed

import (
	"net/http"
	"os"

	"github.com/Ovikx/market-data-recorder/internal/utils"
	"github.com/gorilla/websocket"
)

func ConnectToAlpacaMarketFeed(wsUrl string, symbols []string) ([]*websocket.Conn, error) {
	header := http.Header{}
	header.Add("APCA-API-KEY-ID", os.Getenv("ALPACA_API_KEY"))
	header.Add("APCA-API-SECRET-KEY", os.Getenv("ALPACA_SECRET"))

	// Subscribe to trades
	conn, _, err := utils.SubscribeToWebsocket(wsUrl, header, struct {
		Action string   `json:"action"`
		Trades []string `json:"trades"`
		Quotes []string `json:"quotes"`
		Bars   []string `json:"bars"`
	}{
		Action: "subscribe",
		Trades: symbols,
		Quotes: symbols,
		Bars:   []string{},
	})
	if err != nil {
		return nil, err
	}

	return []*websocket.Conn{conn}, nil
}
