package candlestick

type AssetInfo struct {
	Symbol             string           `json:"symbol"`
	Identifier         AssetIdentifier  `json:"identifier"`
	Pair               string           `json:"pair"`
	BaseAsset          string           `json:"baseAsset"`
	BaseAssetPrecision int              `json:"baseAssetPrecision"`
	QuoteAsset         string           `json:"quoteAsset"`
	QuotePrecision     int              `json:"quotePrecision"`
	Constraints        TradeConstraints `json:"constraints"`
	OnBoardDate        int64            `json:"onBoardDate"`
	Splits             []AssetSplit     `json:"splits"`
}

type AssetSplit struct {
	Time  int64   `json:"time"`
	Ratio float64 `json:"split"`
}

type TradeConstraints struct {
	MaxPrice     float64 `json:"maxPrice"`
	MinPrice     float64 `json:"minPrice"`
	TickSize     float64 `json:"tickSize"`
	MaxQuantity  float64 `json:"maxQuantity"`
	MinQuantity  float64 `json:"minQuantity"`
	StepSize     float64 `json:"stepSize"`
	MaxNumOrders int     `json:"maxNumOrders"`
	MinNotional  float64 `json:"minNotional"`
}

type ExchangeInfo struct {
	Name       string                `json:"name"`
	ExchangeId string                `json:"exchangeId"`
	BrokerId   string                `json:"brokerId"`
	LastUpdate int64                 `json:"lastUpdate"`
	Symbols    map[string]*AssetInfo `json:"symbols"`
	Resolution []int64               `json:"resolution"`
}

type BrokerInfo struct {
	Name string `json:"name"`
}

type ExchangeList struct {
	Exchanges  []*ExchangeInfo        `json:"exchanges"`
	BrokerInfo map[string]*BrokerInfo `json:"brokerInfo"`
}

func (e *ExchangeInfo) Symbol(symbol string) (*AssetInfo, bool) {
	v, ok := e.Symbols[symbol]
	return v, ok
}
