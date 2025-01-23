package timewheel

import (
	"testing"
	"time"
)

func init() {
	Start()
}

// TW_SmallMs TW = 1 //5ms  	(5-100]ms
func Test_TW_SmallMs(t *testing.T) {
	var exit = make(chan int, 1)
	var count int64
	var interval int64 = 5
	var runs int64 = 1000
	var now = time.Now()
	var params = 123
	AddTimer(TW_SmallMs, interval, runs, params, func(p interface{}) {
		count++
		t.Log("TW_SmallMs hello...", count, p)
		if count >= runs && runs > 0 {
			exit <- 1
		}
	})
	select {
	case <-exit:
		t.Log("TW_SmallMs exit...", time.Since(now))
	}
}

// TW_BigMs   TW = 2 //50ms 	(100-1000]ms
func Test_TW_BigMs(t *testing.T) {
	var exit = make(chan int, 1)
	var count int64
	var interval int64 = 101
	var runs int64 = 2
	var now = time.Now()
	var params = "345"
	AddTimer(TW_BigMs, interval, runs, params, func(p interface{}) {
		count++
		t.Log("TW_BigMs hello...", count, p)
		if count >= runs && runs > 0 {
			exit <- 1
		}
	})
	select {
	case <-exit:
		t.Log("TW_BigMs exit...", time.Since(now))
	}
}

// TW_Sec     TW = 3 //1s   	(1-60]s
func Test_TW_Sec(t *testing.T) {
	var exit = make(chan int, 1)
	var count int64
	var interval int64 = Sec2msFloat64(2.101)
	var runs int64 = 2
	var now = time.Now()
	type pp struct {
		person string
		age    int
		addr   []string
	}
	var params = pp{
		person: "xxxx",
		age:    29,
		addr:   []string{"6", "7", "8"},
	}
	AddTimer(TW_Sec, interval, runs, params, func(p interface{}) {
		count++
		t.Log("TW_Sec hello...", count, p)
		if count >= runs && runs > 0 {
			exit <- 1
		}
	})
	select {
	case <-exit:
		t.Log("TW_Sec exit...", time.Since(now))
	}
}

// TW_Min     TW = 4 //60s  	(1-60]min
func Test_TW_Min(t *testing.T) {
	var exit = make(chan int, 1)
	var count int64
	var interval int64 = Min2msFloat64(1.101)
	var runs int64 = 2
	var now = time.Now()
	var params interface{}
	AddTimer(TW_Min, interval, runs, params, func(p interface{}) {
		count++
		t.Log("TW_Min hello...", count)
		if count >= runs && runs > 0 {
			exit <- 1
		}
	})
	select {
	case <-exit:
		t.Log("TW_Min exit...", time.Since(now))
	}
}

func TestTimeWheel(t *testing.T) {
	var exit = make(chan int, 1)
	var count int64
	var interval int64 = Sec2msInt64(2)
	var runs int64 = 10
	var params = "12+23"
	AddTimer(TW_Sec, interval, runs, params, func(p interface{}) {
		count++
		t.Log("hello...", count, p)
		if count >= runs {
			exit <- 1
		}
	})
	select {
	case <-exit:
		t.Log("exit...")
	}
}

func TestStopTimeWheel(t *testing.T) {
	var count int64
	var runs int64 = 2
	AddTimer(TW_Sec, Sec2msFloat64(1.101), runs, nil, func(p interface{}) {
		count++
		t.Log("TW_Sec-1 hello...", count, p)
		time.Sleep(time.Second * 1)
	})
	//time.Sleep(10 * time.Second)
	Stop()
	time.Sleep(3 * time.Second)
	t.Log("退出程序")
}
