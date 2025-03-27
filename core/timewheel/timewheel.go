package timewheel

import (
	"container/list"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type TimeWheel struct {
	ty    TW
	tasks *list.List
	stop  chan bool
}

var (
	_edTimerMgr *edTimerMgr
)

func Start() {
	_edTimerMgr = &edTimerMgr{
		tws:  make(map[TW]*TimeWheel, TW_Max),
		ext:  make(chan *TaskElement, TaskElementMax),
		stop: make(chan bool),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go _edTimerMgr.wheelLoop(_edTimerMgr.stop, &wg)

	for wi := TW_SmallMs; wi <= TW_Max; wi++ {
		wg.Add(1)
		_edTimerMgr.tws[wi] = &TimeWheel{
			ty:    wi,
			tasks: list.New(),
		}
		_edTimerMgr.tws[wi].stop = make(chan bool)
		go _edTimerMgr.tws[wi].taskLoop(_edTimerMgr.ext, _edTimerMgr.stop, &wg)
	}
}

func Stop() {
	if _edTimerMgr == nil && _edTimerMgr.tws == nil {
		logx.Error("stop timewhell is failed! handler is nil!")
		return
	}
	for i := TW_Max; i >= TW_SmallMs; i-- {
		// 1、删除任务
		tasks := _edTimerMgr.tws[i].tasks
		num := tasks.Len()
		for e := tasks.Front(); e != nil; e = e.Next() {
			tasks.Remove(e)
		}

		// 2、停止协程
		logx.Debugf("stop ::::::: %s[%d] remove: %d", TWMap[i], i, num)
		_edTimerMgr.tws[i].stop <- true
	}
}

func AddTimer(ty TW, interval int64, runs int64, params interface{}, fn func(interface{})) {
	_edTimerMgr.addTask(ty, interval, runs, params, fn)
}

func (this *TimeWheel) add(interval int64, runs int64, params interface{}, fn func(interface{})) {
	this.tasks.PushBack(&WheelTask{
		create:   time.Now().Unix(),
		interval: interval,
		runs:     runs,
		params:   params,
		fn:       fn,
	})
}

func (this *TimeWheel) pushBack(e interface{}) {
	this.tasks.PushBack(e)
}

func (this *TimeWheel) pop(e interface{}) {
	this.tasks.Remove(e.(*list.Element))
}

func (this *TimeWheel) remove(id string) {
	for e := this.tasks.Front(); e != nil; e = e.Next() {
		task := e.Value.(*WheelTask)
		if task.id == id {
			this.tasks.Remove(e)
		}
	}
}

func (this *TimeWheel) taskLoop(ext chan *TaskElement, stop chan bool, wg *sync.WaitGroup) {
	logx.Infof("add taskLoop %s[%d] is success!", TWMap[this.ty], this.ty)
	ticker := time.NewTicker(twTickers[this.ty])
	defer func() {
		ticker.Stop()
		wg.Done()
	}()

	for {
		select {
		case <-ticker.C:
			if this.tasks.Len() == 0 {
				continue
			}

			for e := this.tasks.Front(); e != nil; e = e.Next() {
				task := e.Value.(*WheelTask)
				if task.mark == task.runs {
					this.pop(e)
					continue
				}
				if task.mark == 0 && task.start == 0 {
					task.start = time.Now().Unix()
				}
				task.scales += twScale[this.ty]
				subscale := task.interval - task.scales
				if subscale <= 0 {
					go func(t *WheelTask) {
						t.scales = 0
						t.mark++
						t.fn(task.params)
					}(task)
				} else if subscale < twScale[this.ty] && subscale > 0 {
					dst := getTimewheelScale(this.ty, subscale)
					if dst != TW_NO {
						logx.Debugf("[%d] get next ty: %s, src: %s, subscale: %v.", task.interval, TWMap[dst], TWMap[this.ty], subscale)
						ext <- &TaskElement{
							ty: dst,
							task: &list.Element{
								Value: &WheelTask{
									create:   time.Now().Unix(),
									interval: subscale,
									runs:     1,
									params:   task.params,
									fn:       task.fn,
								},
							},
						}
						task.scales = 0
						task.mark++
					}
				}
			}

		case <-this.stop:
			logx.Infof("taskLoop is exit, %s[%d]", TWMap[this.ty], this.ty)
			if this.ty == TW_SmallMs {
				stop <- true
			}
			return
		}
	}
}
