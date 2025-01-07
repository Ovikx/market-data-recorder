package marketfeed

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Ovikx/market-data-recorder/internal/utils"
	"github.com/gorilla/websocket"
)

func getKrakenSignature(urlPath string, data interface{}, secret string) (string, error) {
	var encodedData string

	switch v := data.(type) {
	case string:
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			return "", err
		}
		encodedData = jsonData["nonce"].(string) + v
	case map[string]interface{}:
		dataMap := url.Values{}
		for key, value := range v {
			dataMap.Set(key, fmt.Sprintf("%v", value))
		}
		encodedData = v["nonce"].(string) + dataMap.Encode()
	default:
		return "", fmt.Errorf("invalid data type")
	}
	sha := sha256.New()
	sha.Write([]byte(encodedData))
	shasum := sha.Sum(nil)

	message := append([]byte(urlPath), shasum...)
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha512.New, decoded)
	mac.Write(message)
	macsum := mac.Sum(nil)
	sigDigest := base64.StdEncoding.EncodeToString(macsum)
	return sigDigest, nil
}

func ConnectToKrakenMarketFeed(wsUrl string, symbols []string) ([]*websocket.Conn, error) {
	type Params struct {
		Channel string   `json:"channel"`
		Symbol  []string `json:"symbol"`
	}

	// Subscribe to book
	bookConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, struct {
		Method string `json:"method"`
		Params `json:"params"`
	}{
		Method: "subscribe",
		Params: Params{
			Channel: "book",
			Symbol:  symbols,
		},
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to trades
	tradesConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, struct {
		Method string `json:"method"`
		Params `json:"params"`
	}{
		Method: "subscribe",
		Params: Params{
			Channel: "trade",
			Symbol:  symbols,
		},
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to ticker
	tickerConn, _, err := utils.SubscribeToWebsocket(wsUrl, nil, struct {
		Method string `json:"method"`
		Params `json:"params"`
	}{
		Method: "subscribe",
		Params: Params{
			Channel: "ticker",
			Symbol:  symbols,
		},
	})
	if err != nil {
		return nil, err
	}

	return []*websocket.Conn{tradesConn, bookConn, tickerConn}, err
}
