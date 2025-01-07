package marketfeedmodels

type KrakenChannelName struct {
	Channel string `json:"channel"`
}

// book channel

type krakenBookAsk struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
}

type krakenBookBid struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
}

type krakenBookData struct {
	Asks      []krakenBookAsk `json:"asks"`
	Bids      []krakenBookBid `json:"bids"`
	Checksum  int             `json:"checksum"`
	Symbol    string          `json:"symbol"`
	Timestamp string          `json:"timestamp"`
}

type KrakenBookMessage struct {
	Channel string           `json:"channel"`
	Type    string           `json:"type"`
	Data    []krakenBookData `json:"data"`
}

// trade channel

type krakenTradeTrade struct {
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Qty       float64 `json:"qty"`
	Price     float64 `json:"price"`
	OrdType   string  `json:"ord_type"`
	TradeID   int     `json:"trade_id"`
	Timestamp string  `json:"timestamp"`
}

type KrakenTradeMessage struct {
	Channel string             `json:"channel"`
	Type    string             `json:"type"`
	Data    []krakenTradeTrade `json:"data"`
}

// ticker channel

type krakenTickerData struct {
	Ask       float64 `json:"ask"`
	AskQty    float64 `json:"ask_qty"`
	Bid       float64 `json:"bid"`
	BidQty    float64 `json:"bid_qty"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	High      float64 `json:"high"`
	Last      float64 `json:"last"`
	Low       float64 `json:"low"`
	Symbol    string  `json:"symbol"`
	Volume    float64 `json:"volume"`
	Vwap      float64 `json:"vwap"`
}

type KrakenTickerMessage struct {
	Channel string             `json:"channel"`
	Type    string             `json:"type"`
	Data    []krakenTickerData `json:"data"`
}
