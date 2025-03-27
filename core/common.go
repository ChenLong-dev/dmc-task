package core

import (
	"dmc-task/utils"
	"time"
)

// TaskType 定义任务类型
type TaskType int

const (
	Default             TaskType = iota
	RealTimeSingleTask           // 1 - 实时单任务
	FixedTimeSingleTask          // 2 - 固定时间单任务
	CronCycleTask                // 3 - 定时循环任务
)

var TaskTypeMap = map[TaskType]string{
	RealTimeSingleTask:  "实时单任务",
	FixedTimeSingleTask: "固定时间单任务",
	CronCycleTask:       "定时循环任务",
}

// TaskStatus 定义任务执行状态类型
type TaskStatus int

const (
	Stoped   TaskStatus = -4 // -4 - 已暂停
	Added    TaskStatus = -3 // -3 - 已添加
	Modified TaskStatus = -2 // -2 - 已修改
	Deleted  TaskStatus = -1 // -1 - 已删除
)

const (
	Init     TaskStatus = iota // 0 - 初始化
	Pending                    // 1 - 待执行
	Running                    // 2 - 执行中
	Failed                     // 3 - 失败
	Finished                   // 4 - 已完成
)

var TaskStatusMap = map[TaskStatus]string{
	Init:     "初始化",
	Pending:  "待执行",
	Running:  "执行中",
	Failed:   "失败",
	Finished: "已完成",
}

const (
	DefaultLimit       = 500
	DefaultTimeHorizon = 1
)

const (
	FixCycleSpec  = "@every 55s"     // 每55秒扫描一次数据库（为了解决扫描覆盖的问题）
	FixedScanTime = 59 * time.Second // 扫描最近59秒内将要执行的任务
)

type Result struct {
	Code   int         `json:"code"`   // 错误码
	Bid    string      `json:"bid"`    // 业务ID
	Msg    string      `json:"msg"`    // 错误信息
	Status string      `json:"status"` // 任务状态
	Data   interface{} `json:"data"`   // 返回数据
}

func GetResult(code int, bid string, msg string, status TaskStatus, data interface{}) string {
	res := Result{
		Code:   code,
		Bid:    bid,
		Msg:    msg,
		Status: TaskStatusMap[status],
		Data:   data,
	}
	resByte, err := utils.MarshalByJson(res)
	if err != nil {
		return ""
	}
	return string(resByte)
}
