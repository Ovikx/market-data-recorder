package marketfeed

import (
	"github.com/Ovikx/market-data-recorder/internal/utils"
	"github.com/gorilla/websocket"
)

func ConnectToCoinbaseMarketFeed(wsUrl string, jwtGenerator func(uri string) (string, error), productIds []string) ([]*websocket.Conn, error) {
	// Generate JWT for initial handshake
	jwt, err := jwtGenerator(wsUrl)
	if err != nil {
		return nil, err
	}

	type subscriptionMsg struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channel    string   `json:"channel"`
		JWT        string   `json:"jwt"`
	}

	// Subscribe to the order book channel
	orderBookConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: productIds,
		Channel:    "level2",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}
	err = orderBookConn.WriteJSON(subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: []string{},
		Channel:    "heartbeats",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to the trades channel
	tradesConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: productIds,
		Channel:    "market_trades",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}
	err = tradesConn.WriteJSON(subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: []string{},
		Channel:    "heartbeats",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to the ticker channel
	tickerConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: productIds,
		Channel:    "ticker",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}
	err = tickerConn.WriteJSON(subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: []string{},
		Channel:    "heartbeats",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}

	return []*websocket.Conn{tradesConn, orderBookConn, tickerConn}, err

}
