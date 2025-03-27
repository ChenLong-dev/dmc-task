package timewheel

import "math"

// ----------------------- time ------------------------------
const (
	OneThousand int64 = 1000
	Sixty       int64 = 60
	TwentyFour  int64 = 24
	Thirty      int64 = 30
	Twelve      int64 = 12
)

// Sec2msInt64 sec -> ms
func Sec2msInt64(sec int64) int64 {
	return sec * OneThousand
}

// Min2msInt64 min -> ms
func Min2msInt64(min int64) int64 {
	return Sec2msInt64(min * Sixty)
}

// Hour2msInt64 hour -> ms
func Hour2msInt64(hour int64) int64 {
	return Min2msInt64(hour * Sixty)
}

// Day2msInt64 day -> ms
func Day2msInt64(day int64) int64 {
	return Hour2msInt64(day * TwentyFour)
}

// Month2msInt64 month -> ms
func Month2msInt64(month int64) int64 {
	return Day2msInt64(month * Thirty)
}

// Year2MsInt64 year -> ms
func Year2MsInt64(year int64) int64 {
	return Month2msInt64(year * Twelve)
}

// ------------ float64 -----------------

// Sec2msFloat64 sec -> ms
func Sec2msFloat64(sec float64) int64 {
	return int64(math.Floor(sec * float64(OneThousand)))
}

// Min2msFloat64 min -> ms
func Min2msFloat64(min float64) int64 {
	return Sec2msFloat64(min * float64(Sixty))
}

// Hour2msFloat64 hour -> ms
func Hour2msFloat64(hour float64) int64 {
	return Min2msFloat64(hour * float64(Sixty))
}

// Day2msFloat64 day -> ms
func Day2msFloat64(day float64) int64 {
	return Hour2msFloat64(day * float64(TwentyFour))
}

// Month2msFloat64 month -> ms
func Month2msFloat64(month float64) int64 {
	return Day2msFloat64(month * float64(Thirty))
}

// Year2MsFloat64 year -> ms
func Year2MsFloat64(year float64) int64 {
	return Month2msFloat64(year * float64(Twelve))
}

// ----------------------- timewheel ----------

// getTimewheelScale get timewheel scale by interval
func getTimewheelScale(src TW, interval int64) (dst TW) {
	for i := src - 1; i >= TW_SmallMs; i-- {
		if (interval > 0 && twRange[i].min >= interval) ||
			(twRange[i].min < interval && twRange[i].max >= interval) {
			return i
		}
	}
	return TW_NO
}
