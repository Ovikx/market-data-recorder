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

	conns := make([]*websocket.Conn, 0)
	for _, id := range productIds {
		// Subscribe to the ticker channel
		conn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
			Type:       "subscribe",
			ProductIDs: []string{id},
			Channel:    "ticker",
			JWT:        jwt,
		})
		if err != nil {
			return nil, err
		}
		err = conn.WriteJSON(subscriptionMsg{
			Type:       "subscribe",
			ProductIDs: []string{},
			Channel:    "heartbeats",
			JWT:        jwt,
		})
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
	}

	return conns, err

}
