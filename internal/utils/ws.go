package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Opens a websocket connection with the server at `wsUrl`. `message` is sent to the server with optional headers `requestHeader`.
// Returns the opened websocket connection, the response, and any possible error.
func SubscribeToWebsocket(wsUrl string, requestHeader http.Header, message interface{}) (*websocket.Conn, *http.Response, error) {
	c, resp, err := websocket.DefaultDialer.Dial(wsUrl, requestHeader)
	if err != nil {
		return nil, nil, err
	}
	err = c.WriteJSON(message)
	if err != nil {
		return nil, nil, err
	}
	return c, resp, nil
}
