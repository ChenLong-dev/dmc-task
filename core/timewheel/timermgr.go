package timewheel

import (
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type edTimerMgr struct {
	tws  map[TW]*TimeWheel
	ext  chan *TaskElement
	stop chan bool
}

func (this *edTimerMgr) addTask(ty TW, interval int64, runs int64, params interface{}, fn func(interface{})) {
	tRange, ok := twRange[ty]
	if !ok {
		logx.Error("can not find tw range data, ty: ", ty)
		return
	}
	logx.Debugf("ty:%d, interval:%d, runs:%d, tRange:%+v", ty, interval, runs, tRange)
	if tRange.min > interval || tRange.max < interval {
		logx.Error("tw interval over range, ty: ", ty, interval)
		return
	}
	this.tws[ty].add(interval, runs, params, fn)
}

func (this *edTimerMgr) wheelLoop(stopWheel chan bool, wg *sync.WaitGroup) {
	logx.Debug("add wheelLoop is success!")
	defer wg.Done()
	for {
		select {
		case ele := <-this.ext:
			this.tws[ele.ty].pushBack(ele.task.Value.(*WheelTask))
		case <-stopWheel:
			logx.Infof("wheelLoop is exit")
			return
		}
	}
}
