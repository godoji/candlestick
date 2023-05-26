package candlestick

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
)

type SeriesType string
type AxisType string

const seriesFieldByteSize = 20

const (
	BarChart   = SeriesType("BAR_CHART")
	LineChart  = SeriesType("LINE_CHART")
	PathChart  = SeriesType("PATH_CHART")
	PriceAxis  = AxisType("PRICE_AXIS")
	CustomAxis = AxisType("CUSTOM_AXIS")
)

type Indicator struct {
	Series map[string]*IndicatorSeries `json:"series"`
	Meta   IndicatorMeta               `json:"meta"`
}

const indicatorValueByteSize = 9

type IndicatorValue struct {
	Value   float64 `json:"v"`
	Missing bool    `json:"m"`
}

const seriesByteSize = int(CandleSetSize)*indicatorValueByteSize + 3*seriesFieldByteSize // includes name

type IndicatorSeries struct {
	Values []IndicatorValue `json:"values"`
	Kind   SeriesType       `json:"kind"`
	Axis   AxisType         `json:"axis"`
}

func (b *Indicator) UID() string {
	return b.Meta.UID
}

func (b *Indicator) BlockNumber() int64 {
	return b.Meta.Block
}

func (b *Indicator) IsComplete() bool {
	return b.Meta.Complete
}

func (b *Indicator) LastUpdate() int64 {
	return b.Meta.LastUpdate
}

func (b *Indicator) Symbol() string {
	return b.Meta.Symbol
}

func (b *Indicator) Interval() int64 {
	return b.Meta.Interval
}

func (b *Indicator) TimeStampAtIndex(i int64) int64 {
	return b.UnixFirst() + i*b.Meta.Interval
}

func (b *Indicator) AtTime(series string, timeStamp int64) *IndicatorValue {
	return b.AtIndex(series, b.Index(timeStamp))
}

func (b *Indicator) AtIndex(series string, index int64) *IndicatorValue {
	return &b.Series[series].Values[index]
}

func (b *Indicator) Index(timeStamp int64) int64 {
	return (timeStamp - b.UnixFirst()) / b.Interval()
}

func (b *Indicator) UnixFirst() int64 {
	return b.Meta.Block * b.Meta.Interval * CandleSetSize
}

func (b *Indicator) UnixLast() int64 {
	return (b.Meta.Block+1)*b.Meta.Interval*CandleSetSize - b.Meta.Interval
}

type IndicatorMeta struct {
	UID          string `json:"uid"`
	Block        int64  `json:"block"`
	Complete     bool   `json:"complete"`
	LastUpdate   int64  `json:"lastUpdate"`
	Symbol       string `json:"symbol"`
	Interval     int64  `json:"interval"`
	BaseInterval int64  `json:"baseInterval"`
	Name         string `json:"name"`
	Parameters   []int  `json:"parameters"`
}

func EncodeIndicatorSet(ind *Indicator) ([]byte, error) {

	// encode meta data
	var metaBuf bytes.Buffer
	err := gob.NewEncoder(&metaBuf).Encode(ind.Meta)
	if err != nil {
		return nil, err
	}
	metaBytes := metaBuf.Bytes()

	// create main buffer
	buf := make([]byte, 8+seriesByteSize*len(ind.Series)+len(metaBytes))

	// add number of series
	binary.BigEndian.PutUint64(buf[0:], uint64(len(ind.Series)))

	// add each series
	index := 0
	for key, series := range ind.Series {
		copy(buf[8+index*seriesByteSize+0*seriesFieldByteSize:], key)
		copy(buf[8+index*seriesByteSize+1*seriesFieldByteSize:], series.Kind)
		copy(buf[8+index*seriesByteSize+2*seriesFieldByteSize:], series.Axis)
		offset := 8 + index*seriesByteSize + 3*seriesFieldByteSize
		for i, v := range series.Values {
			binary.BigEndian.PutUint64(buf[offset+i*indicatorValueByteSize:], math.Float64bits(v.Value))
			isMissing := uint8(0)
			if v.Missing {
				isMissing = 1
			}
			buf[offset+i*indicatorValueByteSize+8] = isMissing
		}
		index++
	}

	// copy meta bytes
	copy(buf[8+seriesByteSize*len(ind.Series):], metaBytes)

	return buf, nil
}

func DecodeIndicatorSet(data []byte) (*Indicator, error) {

	numberOfSeries := int(binary.BigEndian.Uint64(data[0:]))

	seriesMap := make(map[string]*IndicatorSeries)

	for index := 0; index < numberOfSeries; index++ {

		series := &IndicatorSeries{
			Values: make([]IndicatorValue, CandleSetSize),
			Kind:   "",
			Axis:   "",
		}

		// decode text fields
		o := 8 + index*seriesByteSize
		bytesKey := data[o+0*seriesFieldByteSize : o+1*seriesFieldByteSize]
		key := string(bytesKey[:bytes.IndexByte(bytesKey, 0)])
		bytesKind := data[o+1*seriesFieldByteSize : o+2*seriesFieldByteSize]
		series.Kind = SeriesType(bytesKind[:bytes.IndexByte(bytesKind, 0)])
		bytesAxis := data[o+2*seriesFieldByteSize : o+3*seriesFieldByteSize]
		series.Axis = AxisType(bytesAxis[:bytes.IndexByte(bytesAxis, 0)])

		// decode values
		o = 8 + index*seriesByteSize + 3*seriesFieldByteSize
		for i := 0; i < int(CandleSetSize); i++ {
			series.Values[i] = IndicatorValue{
				Value:   math.Float64frombits(binary.BigEndian.Uint64(data[o+i*indicatorValueByteSize:])),
				Missing: data[o+i*indicatorValueByteSize+8] == 1,
			}
		}

		seriesMap[key] = series
	}

	metaBytes := bytes.NewReader(data[8+seriesByteSize*numberOfSeries:])
	var meta IndicatorMeta
	err := gob.NewDecoder(metaBytes).Decode(&meta)
	if err != nil {
		return nil, err
	}

	ind := &Indicator{
		Series: seriesMap,
		Meta:   meta,
	}

	return ind, nil
}
