package timewheel

type WheelTask struct {
	id       string
	create   int64
	start    int64
	interval int64
	runs     int64
	params   interface{}
	fn       func(interface{})
	scales   int64
	mark     int64
}

func (wt *WheelTask) reset() {
	wt.mark = 0
}
