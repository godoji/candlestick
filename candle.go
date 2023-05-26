package candlestick

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
)

const CandleSetSize = int64(5000)
const candleSetByteSize = 65

type Candle struct {
	Open           float64 `json:"o"`
	High           float64 `json:"h"`
	Low            float64 `json:"l"`
	Close          float64 `json:"c"`
	Volume         float64 `json:"v"`
	TakerVolume    float64 `json:"tv"`
	NumberOfTrades int64   `json:"not"`
	Time           int64   `json:"t"`
	Missing        bool    `json:"m"`
}

type DataSetMeta struct {
	UID        string `json:"uid"`
	Block      int64  `json:"block"`
	Complete   bool   `json:"complete"`
	LastUpdate int64  `json:"lastUpdate"`
	Symbol     string `json:"symbol"`
	Interval   int64  `json:"interval"`
}

type CandleSet struct {
	Candles []Candle    `json:"candles"`
	Meta    DataSetMeta `json:"meta"`
}

func (b *CandleSet) UID() string {
	return b.Meta.UID
}

func (b *CandleSet) BlockNumber() int64 {
	return b.Meta.Block
}

func (b *CandleSet) IsComplete() bool {
	return b.Meta.Complete
}

func (b *CandleSet) LastUpdate() int64 {
	return b.Meta.LastUpdate
}

func (b *CandleSet) Symbol() string {
	return b.Meta.Symbol
}

func (b *CandleSet) Interval() int64 {
	return b.Meta.Interval
}

func (b *CandleSet) TimeStampAtIndex(i int64) int64 {
	return b.UnixFirst() + i*b.Meta.Interval
}

func (b *CandleSet) AtTime(timeStamp int64) *Candle {
	return b.AtIndex(b.Index(timeStamp))
}

func (b *CandleSet) AtIndex(index int64) *Candle {
	return &b.Candles[index]
}

func (b *CandleSet) Index(timeStamp int64) int64 {
	return (timeStamp - b.UnixFirst()) / b.Interval()
}

func (b *CandleSet) UnixFirst() int64 {
	return b.Meta.Block * b.Meta.Interval * CandleSetSize
}

func (b *CandleSet) UnixLast() int64 {
	return (b.Meta.Block+1)*b.Meta.Interval*CandleSetSize - b.Meta.Interval
}

func EncodeCandleSet(b *CandleSet) ([]byte, error) {

	// encode meta data
	var metaBuf bytes.Buffer
	err := gob.NewEncoder(&metaBuf).Encode(b.Meta)
	if err != nil {
		return nil, err
	}
	metaBytes := metaBuf.Bytes()

	// create main buffer
	cSize := candleSetByteSize
	candleDataSize := 8 + cSize*len(b.Candles)
	buf := make([]byte, candleDataSize+len(metaBytes))

	// add number of candles
	binary.BigEndian.PutUint64(buf[0:], uint64(len(b.Candles)))

	// add candle data
	for i, c := range b.Candles {
		binary.BigEndian.PutUint64(buf[8+i*cSize+0:], math.Float64bits(c.Open))
		binary.BigEndian.PutUint64(buf[8+i*cSize+8:], math.Float64bits(c.High))
		binary.BigEndian.PutUint64(buf[8+i*cSize+16:], math.Float64bits(c.Low))
		binary.BigEndian.PutUint64(buf[8+i*cSize+24:], math.Float64bits(c.Close))
		binary.BigEndian.PutUint64(buf[8+i*cSize+32:], math.Float64bits(c.Volume))
		binary.BigEndian.PutUint64(buf[8+i*cSize+40:], math.Float64bits(c.TakerVolume))
		binary.BigEndian.PutUint64(buf[8+i*cSize+48:], uint64(c.NumberOfTrades))
		binary.BigEndian.PutUint64(buf[8+i*cSize+56:], uint64(c.Time))
		isMissing := uint8(0)
		if c.Missing {
			isMissing = 1
		}
		buf[8+i*cSize+64] = isMissing
	}

	// copy meta bytes
	copy(buf[8+len(b.Candles)*cSize:], metaBytes)

	return buf, nil
}

func DecodeCandleSet(data []byte) (*CandleSet, error) {

	cSize := candleSetByteSize
	numberOfCandles := int(binary.BigEndian.Uint64(data[0:]))

	candles := make([]Candle, numberOfCandles)
	for i := 0; i < numberOfCandles; i++ {
		candles[i] = Candle{
			Open:           math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+0:])),
			High:           math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+8:])),
			Low:            math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+16:])),
			Close:          math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+24:])),
			Volume:         math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+32:])),
			TakerVolume:    math.Float64frombits(binary.BigEndian.Uint64(data[8+i*cSize+40:])),
			NumberOfTrades: int64(binary.BigEndian.Uint64(data[8+i*cSize+48:])),
			Time:           int64(binary.BigEndian.Uint64(data[8+i*cSize+56:])),
			Missing:        data[8+i*cSize+64] == 1,
		}
	}

	metaBytes := bytes.NewReader(data[8+numberOfCandles*cSize:])
	var meta DataSetMeta
	err := gob.NewDecoder(metaBytes).Decode(&meta)
	if err != nil {
		return nil, err
	}

	cs := &CandleSet{
		Candles: candles,
		Meta:    meta,
	}

	return cs, nil
}

func BlockToUnix(block int64, interval int64) int64 {
	return block * CandleSetSize * interval
}

func UnixToBlock(unixTime int64, interval int64) int64 {
	b := unixTime / (CandleSetSize * interval)
	if unixTime < 0 && unixTime%(CandleSetSize*interval) != 0 {
		b--
	}
	return b
}
