package candlestick

const (
	Interval1m  = 60
	Interval3m  = 3 * Interval1m
	Interval5m  = 5 * Interval1m
	Interval15m = 15 * Interval1m
	Interval30m = 2 * Interval15m
	Interval45m = 3 * Interval15m
	Interval1h  = 4 * Interval15m
	Interval2h  = 2 * Interval1h
	Interval4h  = 2 * Interval2h
	Interval6h  = 3 * Interval2h
	Interval8h  = 2 * Interval4h
	Interval12h = 2 * Interval6h
	Interval1d  = 2 * Interval12h
	Interval3d  = 3 * Interval1d
)

var IntervalList = []int64{
	Interval1m,
	Interval3m,
	Interval5m,
	Interval15m,
	Interval30m,
	Interval45m,
	Interval1h,
	Interval2h,
	Interval4h,
	Interval6h,
	Interval8h,
	Interval12h,
	Interval1d,
	Interval3d,
}

var IntervalMap = map[int64]int64{
	Interval3m:  Interval1m,
	Interval5m:  Interval1m,
	Interval15m: Interval5m,
	Interval30m: Interval15m,
	Interval45m: Interval15m,
	Interval1h:  Interval15m,
	Interval2h:  Interval1h,
	Interval4h:  Interval2h,
	Interval6h:  Interval2h,
	Interval8h:  Interval4h,
	Interval12h: Interval6h,
	Interval1d:  Interval12h,
	Interval3d:  Interval1d,
}
