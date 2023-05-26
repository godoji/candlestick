package candlestick

import (
	"strings"
)

type uniqueAssetIdentifier struct {
	Broker   string `json:"broker"`
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
}

type AssetIdentifier = *uniqueAssetIdentifier

func (u AssetIdentifier) ToString() string {
	return u.Broker + ":" + u.Exchange + ":" + u.Symbol
}

func ParseSymbol(s string) (AssetIdentifier, bool) {
	// Split symbol into segments
	xs := strings.Split(s, ":")
	if len(xs) != 3 {
		return nil, false
	}
	// Return struct
	return &uniqueAssetIdentifier{
		Broker:   xs[0],
		Exchange: xs[1],
		Symbol:   xs[2],
	}, true
}

func NewAssetIdentifier(broker string, exchange string, symbol string) AssetIdentifier {
	return &uniqueAssetIdentifier{Broker: broker, Exchange: exchange, Symbol: symbol}
}
