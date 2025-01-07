package marketfeedmodels

type CoinbaseChannelName struct {
	Channel string `json:"channel"`
}

// market_trades channel

type coinbaseTradesTrade struct {
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"`
	Time      string `json:"time"`
}

type coinbaseTradesEvent struct {
	Type   string                `json:"string"`
	Trades []coinbaseTradesTrade `json:"trades"`
}

// Represents the schema for Coinbase's `market_trades` websocket channel. See https://docs.cdp.coinbase.com/advanced-trade/docs/ws-channels#market-trades-channel for more info.
type CoinbaseTradesMessage struct {
	Channel     string `json:"channel"`
	ClientID    string `json:"client_id"`
	Timestamp   string `json:"timestamp"`
	SequenceNum int    `json:"sequence_num"`
	Events      []coinbaseTradesEvent
}

// l2_data channel

type coinbaseL2Update struct {
	Side        string `json:"side"`
	EventTime   string `json:"event_time"`
	PriceLevel  string `json:"price_level"`
	NewQuantity string `json:"new_quantity"`
}

type coinbaseL2Event struct {
	Type      string             `json:"type"`
	ProductID string             `json:"product_id"`
	Updates   []coinbaseL2Update `json:"updates"`
}

// Represents the schema for Coinbase's `l2_data` websocket channel.
type CoinbaseL2Message struct {
	Channel     string            `json:"channel"`
	ClientID    string            `json:"client_id"`
	Timestamp   string            `json:"timestamp"`
	SequenceNum int               `json:"sequence_num"`
	Events      []coinbaseL2Event `json:"events"`
}

// ticker channel

type coinbaseTickerEvent struct {
	Type    string                 `json:"type"`
	Tickers []coinbaseTickerTicker `json:"tickers"`
}

type coinbaseTickerTicker struct {
	Type               string `json:"type"`
	ProductID          string `json:"product_id"`
	Price              string `json:"price"`
	Volume24H          string `json:"volume_24_h"`
	Low24H             string `json:"low_24_h"`
	High24H            string `json:"high_24_h"`
	Low52W             string `json:"low_52_w"`
	High52W            string `json:"high_52_w"`
	PricePercentChg24H string `json:"price_percent_chg_24_h"`
	BestBid            string `json:"best_bid"`
	BestBidQuantity    string `json:"best_bid_quantity"`
	BestAsk            string `json:"best_ask"`
	BestAskQuantity    string `json:"best_ask_quantity"`
}
type CoinbaseTickerMessage struct {
	Channel     string                `json:"channel"`
	ClientID    string                `json:"client_id"`
	Timestamp   string                `json:"timestamp"`
	SequenceNum int                   `json:"sequence_num"`
	Events      []coinbaseTickerEvent `json:"events"`
}
