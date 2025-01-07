package updatefeed

import (
	"github.com/Ovikx/market-data-recorder/internal/utils"
	"github.com/gorilla/websocket"
)

func ConnectToCoinbaseUpdateFeed(wsUrl string, jwtGenerator func(uri string) (string, error), productIds []string) ([]*websocket.Conn, error) {
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
	userConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: productIds,
		Channel:    "user",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}

	err = userConn.WriteJSON(subscriptionMsg{
		Type:       "subscribe",
		ProductIDs: []string{},
		Channel:    "heartbeats",
		JWT:        jwt,
	})
	if err != nil {
		return nil, err
	}

	return []*websocket.Conn{userConn}, nil
}
