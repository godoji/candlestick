package candlestick

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func randomCandleSet() *CandleSet {
	data := &CandleSet{
		Candles: make([]Candle, CandleSetSize),
		Meta: DataSetMeta{
			UID:        "test_uid_name",
			Block:      534859,
			Complete:   true,
			LastUpdate: 1685903959,
			Symbol:     "AAPLUSD",
			Interval:   60,
		},
	}
	for i := range data.Candles {
		data.Candles[i] = Candle{
			Open:           rand.Float64(),
			High:           rand.Float64(),
			Low:            rand.Float64(),
			Close:          rand.Float64(),
			Volume:         rand.Float64(),
			TakerVolume:    rand.Float64(),
			NumberOfTrades: rand.Int63(),
			Time:           rand.Int63(),
			Missing:        rand.Intn(1) == 1,
		}
	}
	return data
}

func TestCandlesBinary(t *testing.T) {
	data := randomCandleSet()
	bin, err := EncodeCandleSet(data)
	if err != nil {
		log.Fatalln(err)
	}
	decoded, err := DecodeCandleSet(bin)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range data.Candles {
		if data.Candles[i].Time != decoded.Candles[i].Time {
			fmt.Printf("candle %d did not match:\n", i)
			fmt.Println(data.Candles[i])
			fmt.Println(decoded.Candles[i])
			t.FailNow()
		}
		if data.Meta.UID != decoded.Meta.UID {
			fmt.Printf("meta did not match:\n")
			fmt.Println(data.Meta)
			fmt.Println(decoded.Meta)
			t.FailNow()
		}
	}
}

func BenchmarkCandlesEncodeBinary(b *testing.B) {
	candles := randomCandleSet()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeCandleSet(candles)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkCandlesDecodeBinary(b *testing.B) {
	candles := randomCandleSet()
	bin, err := EncodeCandleSet(candles)
	if err != nil {
		log.Fatalln(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = DecodeCandleSet(bin)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkCandlesEncodeGob(b *testing.B) {
	candles := randomCandleSet()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := encoder.Encode(candles)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkCandlesDecodeGob(b *testing.B) {
	candles := randomCandleSet()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(candles)
	if err != nil {
		log.Fatalln(err)
	}
	result := &CandleSet{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := bytes.NewBuffer(buf.Bytes())
		err = gob.NewDecoder(b).Decode(result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TestUnixToBlock(t *testing.T) {

	if UnixToBlock(0, Interval1d) != 0 {
		t.FailNow()
	}
	if UnixToBlock(1, Interval1d) != 0 {
		t.FailNow()
	}
	if UnixToBlock(-431000000, Interval1d) != -1 {
		t.FailNow()
	}
	if UnixToBlock(-432000000, Interval1d) != -1 {
		t.FailNow()
	}
	if UnixToBlock(-433000000, Interval1d) != -2 {
		t.FailNow()
	}

}

func TestBlockToUnix(t *testing.T) {
	if BlockToUnix(-1, Interval1d) != -432000000 {
		t.FailNow()
	}
}
