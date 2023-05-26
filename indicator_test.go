package candlestick

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func randomIndicatorSet() *Indicator {
	data := &Indicator{
		Series: map[string]*IndicatorSeries{},
		Meta: IndicatorMeta{
			UID:          "test_indicat_uid_name",
			Block:        534859,
			Complete:     true,
			LastUpdate:   1685903959,
			Symbol:       "AAPLUSD",
			Interval:     3600,
			BaseInterval: 60,
			Name:         "nice",
			Parameters:   []int{200, 103},
		},
	}
	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("series%d", i)
		values := make([]IndicatorValue, CandleSetSize)
		for i2 := range values {
			values[i2] = IndicatorValue{
				Value:   rand.Float64(),
				Missing: false,
			}
		}
		data.Series[key] = &IndicatorSeries{
			Values: values,
			Kind:   LineChart,
			Axis:   CustomAxis,
		}
	}
	return data
}

func TestIndicatorBinary(t *testing.T) {
	data := randomIndicatorSet()
	bin, err := EncodeIndicatorSet(data)
	if err != nil {
		log.Fatalln(err)
	}
	decoded, err := DecodeIndicatorSet(bin)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range data.Series {
		if len(data.Series[i].Values) != len(decoded.Series[i].Values) {
			fmt.Printf("series lengths did not match: %d and %d\n", len(data.Series[i].Values), len(decoded.Series[i].Values))
			t.FailNow()
		}
		for j := range data.Series[i].Values {
			if data.Series[i].Values[j].Value != decoded.Series[i].Values[j].Value {
				fmt.Printf("indicator value %d did not match:\n", j)
				fmt.Println(data.Series[i].Values[j])
				fmt.Println(decoded.Series[i].Values[j])
				t.FailNow()
			}
		}
		if data.Series[i].Kind != decoded.Series[i].Kind {
			fmt.Printf("series meta did not match:\n")
			t.FailNow()
		}
	}
	if data.Meta.UID != decoded.Meta.UID {
		fmt.Printf("meta did not match:\n")
		fmt.Println(data.Meta)
		fmt.Println(decoded.Meta)
		t.FailNow()
	}
}

func BenchmarkIndicatorEncodeBinary(b *testing.B) {
	ind := randomIndicatorSet()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeIndicatorSet(ind)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkIndicatorDecodeBinary(b *testing.B) {
	ind := randomIndicatorSet()
	bin, err := EncodeIndicatorSet(ind)
	if err != nil {
		log.Fatalln(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = DecodeIndicatorSet(bin)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkIndicatorEncodeGob(b *testing.B) {
	ind := randomIndicatorSet()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := encoder.Encode(ind)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkIndicatorDecodeGob(b *testing.B) {
	ind := randomIndicatorSet()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(ind)
	if err != nil {
		log.Fatalln(err)
	}
	result := &Indicator{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := bytes.NewBuffer(buf.Bytes())
		err = gob.NewDecoder(b).Decode(result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
