package mockcoinbaseexecutor_test

import (
	"testing"

	"github.com/Ovikx/market-data-recorder/internal/executor/mockcoinbaseexecutor"
	"github.com/Ovikx/market-data-recorder/internal/jwtgen"
)

func TestMarketBuy(t *testing.T) {
	t.Parallel()

	e := mockcoinbaseexecutor.New("https://api.coinbase.com/api/v3/brokerage", jwtgen.CoinbaseJWT)
	if err := e.MarketBuy("USDT-USD", 0.01); err != nil {
		t.Errorf("failed to market buy: %v", err)
	}
}

func TestMarketSell(t *testing.T) {
	t.Parallel()

	e := mockcoinbaseexecutor.New("https://api.coinbase.com/api/v3/brokerage", jwtgen.CoinbaseJWT)
	if err := e.MarketSell("USDT-USD", 0.01); err != nil {
		t.Errorf("failed to market sell: %v", err)
	}
}
