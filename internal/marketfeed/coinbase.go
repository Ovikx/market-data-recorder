package marketfeed

import (
	"github.com/Ovikx/market-data-recorder/internal/utils"
	"github.com/gorilla/websocket"
)

func ConnectToCoinbaseMarketFeed(wsUrl string, jwtGenerator func(uri string) (string, error), productIds []string, ticks, orders, trades bool) ([]*websocket.Conn, []func() (*websocket.Conn, error), error) {
	type subscriptionMsg struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channel    string   `json:"channel"`
		JWT        string   `json:"jwt"`
	}

	conns := make([]*websocket.Conn, 0)
	reconnectFuncs := make([]func() (*websocket.Conn, error), 0)

	chs := make([]string, 0)
	if ticks {
		chs = append(chs, "ticker")
	}
	if orders {
		chs = append(chs, "level2")
	}
	if trades {
		chs = append(chs, "market_trades")
	}
	for _, channel := range chs {
		for _, id := range productIds {
			connectFunc := func() (*websocket.Conn, error) {
				// Generate JWT for initial handshake
				jwt, err := jwtGenerator(wsUrl)
				if err != nil {
					return nil, err
				}

				// Subscribe to the ticker channel
				conn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, subscriptionMsg{
					Type:       "subscribe",
					ProductIDs: []string{id},
					Channel:    channel,
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
				return conn, err
			}

			conn, err := connectFunc()
			if err != nil {
				return nil, nil, err
			}

			conns = append(conns, conn)
			reconnectFuncs = append(reconnectFuncs, connectFunc)
		}
	}

	return conns, reconnectFuncs, nil

}
