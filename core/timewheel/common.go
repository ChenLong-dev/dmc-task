package timewheel

import (
	"container/list"
	"time"
)

// https://github.com/Peakchen/akTimeWheel/tree/main

type TW int

const (
	TW_NO      TW = 0
	TW_SmallMs TW = 1 //5ms  	(5-100]ms
	TW_BigMs   TW = 2 //50ms 	(100-1000]ms
	TW_Sec     TW = 3 //1s   	(1-60]s
	TW_Min     TW = 4 //60s  	(1-60]min
	TW_Day     TW = 5 //30min 	(1-24]h
	TW_Month   TW = 6 //12h    (1-30]d
	TW_Year    TW = 7 //15d    (1-12]month

	TW_Max TW = TW_Year
)

var TWMap = map[TW]string{
	TW_NO:      "NO",
	TW_SmallMs: "SmallMs",
	TW_BigMs:   "BigMs",
	TW_Sec:     "Sec",
	TW_Min:     "Min",
	TW_Day:     "Day",
	TW_Month:   "Month",
	TW_Year:    "Year",
}

type TWRange struct {
	min int64
	max int64
}

var (
	twTickers = map[TW]time.Duration{
		TW_SmallMs: 5 * time.Millisecond,
		TW_BigMs:   50 * time.Millisecond,
		TW_Sec:     1 * time.Second,
		TW_Min:     60 * time.Second,
		TW_Day:     30 * time.Minute,
		TW_Month:   12 * time.Hour,
		TW_Year:    15 * 24 * time.Hour,
	}

	twScale = map[TW]int64{
		TW_SmallMs: 5,
		TW_BigMs:   50,
		TW_Sec:     Sec2msInt64(1),
		TW_Min:     Sec2msInt64(60),
		TW_Day:     Min2msInt64(30),
		TW_Month:   Hour2msInt64(12),
		TW_Year:    Day2msInt64(15),
	}

	twRange = map[TW]*TWRange{
		TW_SmallMs: &TWRange{min: 5, max: 100},
		TW_BigMs:   &TWRange{min: 100, max: 1000},
		TW_Sec:     &TWRange{min: Sec2msInt64(1), max: Sec2msInt64(60)},
		TW_Min:     &TWRange{min: Min2msInt64(1), max: Min2msInt64(60)},
		TW_Day:     &TWRange{min: Hour2msInt64(1), max: Hour2msInt64(24)},
		TW_Month:   &TWRange{min: Day2msInt64(1), max: Day2msInt64(30)},
		TW_Year:    &TWRange{min: Month2msInt64(1), max: Month2msInt64(12)},
	}
)

const (
	TaskElementMax = 1000
)

type TaskElement struct {
	ty   TW
	task *list.Element
}
