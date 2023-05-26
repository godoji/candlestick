package candlestick

import "math"

type PositionSide string

const Long = PositionSide("LONG")
const Short = PositionSide("SHORT")

type OrderSide string

const Buy = OrderSide("BUY")
const Sell = OrderSide("SELL")

type OrderKind string

const OrderStopLoss = OrderKind("SL")
const OrderTakeProfit = OrderKind("TP")
const OrderLimit = OrderKind("LIMIT")
const OrderMarket = OrderKind("MARKET")

type Position struct {
	Entry      float64 `json:"entry"`
	Amount     float64 `json:"amount"`
	Symbol     string  `json:"symbol"`
	LastUpdate int64   `json:"lastUpdate"`
	LastPrice  float64 `json:"lastPrice"`
}

func (p *Position) Abs() float64 {
	return math.Abs(p.Amount)
}

type Order struct {
	Id           string       `json:"id"`
	Symbol       string       `json:"symbol"`
	Price        float64      `json:"price"`
	Side         OrderSide    `json:"side"`
	Kind         OrderKind    `json:"kind"`
	PositionSide PositionSide `json:"positionSide"`
	Close        bool         `json:"close"`
	Time         int64        `json:"time"`
	Amount       float64      `json:"amount"`
	Filled       float64      `json:"filled"`
}

func (o *Order) PositionAmount() float64 {
	if o.Side == Sell {
		return -o.Amount
	} else {
		return o.Amount
	}
}

type HedgedPosition struct {
	Long  Position `json:"long"`
	Short Position `json:"short"`
}

func (h *HedgedPosition) Position(side PositionSide) *Position {
	if side == Long {
		return &h.Long
	} else {
		return &h.Short
	}
}
