package mockcoinbaseexecutor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type mockcoinbaseexecutor struct {
	baseUrl      string
	jwtGenerator func(uri string) (string, error)
}

func New(baseUrl string, jwtGenerator func(uri string) (string, error)) *mockcoinbaseexecutor {
	return &mockcoinbaseexecutor{baseUrl: baseUrl, jwtGenerator: jwtGenerator}
}

// TODO: figure out how to get order updates. look into merging the ws channels for those with the array of ws channels for market data
func (e *mockcoinbaseexecutor) MarketBuy(symbol string, quoteSize float64) error {
	type marketMarketIoc struct {
		QuoteSize string `json:"quote_size"`
	}

	type orderConfiguration struct {
		MarketMarketIoc marketMarketIoc `json:"market_market_ioc"`
	}

	type orderBody struct {
		ProductID          string             `json:"product_id"`
		Side               string             `json:"side"`
		OrderConfiguration orderConfiguration `json:"order_configuration"`
	}

	jwt, err := e.jwtGenerator(fmt.Sprintf("POST %s/orders/preview", strings.ReplaceAll(e.baseUrl, "https://", "")))
	if err != nil {
		return fmt.Errorf("failed to generate jwt: %v", err)
	}

	body := orderBody{
		ProductID: symbol,
		Side:      "BUY",
		OrderConfiguration: orderConfiguration{MarketMarketIoc: marketMarketIoc{
			QuoteSize: fmt.Sprintf("%f", quoteSize),
		}},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.baseUrl+"/orders/preview", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%v %v", resp, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected code 200, received %v, %v", resp.StatusCode, resp.Status)
	}

	return nil
}

func (e *mockcoinbaseexecutor) MarketSell(symbol string, baseSize float64) error {
	type marketMarketIoc struct {
		BaseSize string `json:"base_size"`
	}

	type orderConfiguration struct {
		MarketMarketIoc marketMarketIoc `json:"market_market_ioc"`
	}

	type orderBody struct {
		ProductID          string             `json:"product_id"`
		Side               string             `json:"side"`
		OrderConfiguration orderConfiguration `json:"order_configuration"`
	}

	jwt, err := e.jwtGenerator(fmt.Sprintf("POST %s/orders/preview", strings.ReplaceAll(e.baseUrl, "https://", "")))
	if err != nil {
		return fmt.Errorf("failed to generate jwt: %v", err)
	}

	body := orderBody{
		ProductID: symbol,
		Side:      "SELL",
		OrderConfiguration: orderConfiguration{MarketMarketIoc: marketMarketIoc{
			BaseSize: fmt.Sprintf("%f", baseSize),
		}},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.baseUrl+"/orders/preview", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%v %v", resp, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected code 200, received %v, %v", resp.StatusCode, resp.Status)
	}

	return nil
}

func (e *mockcoinbaseexecutor) LimitBuy(symbol string, size float64, limitPrice float64) error {
	return nil
}
func (e *mockcoinbaseexecutor) StopLimitSell(symbol string, size float64, limitPrice float64) error {
	return nil
}
