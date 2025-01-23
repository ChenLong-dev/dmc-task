package method

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type CronTaskData struct {
	Id            string    `json:"id"`             // 定时任务ID
	Type          int64     `json:"type"`           // 任务类型
	Bid           string    `json:"bid"`            // 业务ID
	CallbackParam string    `json:"callback_param"` // 定时任务回调参数
	ExecParams    string    `json:"exec_params"`    // 定时任务参数
	ExecTime      time.Time `json:"exec_time"`      // 定时任务执行的时间
}

type JobData struct {
	Id            string `json:"id"`             // job的任务ID
	Type          int64  `json:"type"`           // 任务类型
	CronTaskId    string `json:"cron_task_id"`   // 定时任务ID
	Bid           string `json:"bid"`            // 业务ID
	CallbackParam string `json:"callback_param"` // 回调参数
	ExecParams    string `json:"exec_params"`    // job的参数
}

type CallbackParam struct {
	CallbackUrl string `json:"callback_url"` // 回调地址
	Method      string `json:"method"`       // 请求方法
	Retry       int    `json:"retry"`        // 重试次数
	Timeout     int    `json:"timeout"`      // 超时时间
}

type DoInterface interface {
	DoJob(ctx context.Context, execParams interface{}) error
	DoCallback(ctx context.Context, callbackParam, result string) error
}

type CustomMethod struct {
}

func NewCustomMethod() *CustomMethod {
	return &CustomMethod{}
}

func (c *CustomMethod) DoJob(ctx context.Context, execParams interface{}) error {
	// TODO: 实现自定义的job逻辑
	time.Sleep(5 * time.Second) // 模拟耗时操作
	logx.Infof("==== ==== execParams: %v", execParams)
	return nil
}

func (c *CustomMethod) DoCallback(ctx context.Context, callbackParam, result string) error {
	// TODO: 实现自定义的回调逻辑
	time.Sleep(3 * time.Second) // 模拟耗时操作
	logx.Infof("<==== ====> callback result: %s", result)
	return nil
}
