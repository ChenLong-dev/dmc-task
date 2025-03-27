// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package types

type AddCronCycleTaskReq struct {
	CronCycleTask
}

type AddFixedTimeSingleTaskReq struct {
	FixedTimeSingleTask
}

type AddRealTimeSingleTaskReq struct {
	RealTimeSingleTask
}

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type BaseData struct {
	Id         string `json:"id"`
	Status     int64  `json:"status"`
	UpdateTime string `json:"update_time"`
	CreateTime string `json:"create_time"`
}

type CronCycleTask struct {
	Type     int64  `json:"type" validate:"required,min=1,max=3"`
	BizCode  string `json:"biz_code" validate:"required"`
	Cron     string `json:"cron" validate:"required"`
	ExecPath string `json:"exec_path" validate:"required"`
	Param    string `json:"param" validate:"required"`
	Timeout  int64  `json:"timeout" validate:"required,min=5"`
	ExtInfo  string `json:"ext_info,optional"`
}

type CronCycleTaskData struct {
	BaseData
	CronCycleTask
}

type DelCronCycleTaskReq struct {
	Id string `json:"id" validate:"required"`
}

type DelFixedTimeSingleTaskReq struct {
	Id string `json:"id" validate:"required"`
}

type FilterBase struct {
	Id         string `json:"id,optional"`
	BizCode    string `json:"biz_code,optional"`
	BizId      string `json:"biz_id,optional"`
	CronTaskId string `json:"cron_task_id,optional"`
	Status     int64  `json:"status,default=4" validate:"oneof=-4 -3 -2 -1 0 1 2 3 4"`
	TimeType   string `json:"time_type,default=create_time" validate:"oneof=create_time update_time start_time finish_time"`
	Start      string `json:"start,optional" validate:"omitempty,checkDate"`
	End        string `json:"end,optional" validate:"omitempty,checkDate"`
}

type FixedTimeSingleTask struct {
	Type     int64  `json:"type" validate:"required,min=1,max=3"`
	BizCode  string `json:"biz_code" validate:"required"`
	BizId    string `json:"biz_id,optional"`
	ExecPath string `json:"exec_path" validate:"required"`
	ExecTime int64  `json:"exec_time" validate:"required"`
	Param    string `json:"param" validate:"required"`
	Timeout  int64  `json:"timeout" validate:"required,min=5"`
	ExtInfo  string `json:"ext_info,optional"`
}

type FixedTimeSingleTaskData struct {
	BaseData
	FixedTimeSingleTask
	StartTime  string `json:"start_time"`
	FinishTime string `json:"finish_time"`
	Interval   int64  `json:"interval"`
	ResultMsg  string `json:"result_msg"`
}

type ModCronCycleTaskReq struct {
	Id string `json:"id" validate:"required"`
	CronCycleTask
}

type PageBase struct {
	Total    int64 `json:"total,optional"`
	Page     int64 `json:"page,optional"`
	PageSize int64 `json:"page_size,optional"`
}

type PostDemo struct {
	Name           string `json:"name" validate:"required"`                           // 姓名
	Age            int64  `json:"age" validate:"required,gte=1,lte=130"`              // 年龄
	Mobile         string `json:"mobile,optional" validate:"omitempty,checkMobile"`   // 手机号码
	Email          string `json:"email,optional" validate:"omitempty,checkEmail"`     // 邮箱地址
	Date           string `json:"date" validate:"omitempty,checkDate,checkAfterDate"` // 时间
	Password       string `json:"password" validate:"required"`                       // 密码
	ConfimPassword string `json:"confimPassword" validate:"eqfield=Password"`         // 确认密码
}

type PostDemoReq struct {
	PostDemo
}

type QueryCronCycleTaskReq struct {
	Filter FilterBase `json:"filter" validate:"required"`
	Page   PageBase   `json:"page" validate:"required"`
}

type QueryCronCycleTaskResp struct {
	Base
	Data []CronCycleTaskData `json:"data"`
	Page PageBase            `json:"page"`
}

type QueryFixedTimeSingleTaskReq struct {
	Filter FilterBase `json:"filter" validate:"required"`
	Page   PageBase   `json:"page" validate:"required"`
}

type QueryFixedTimeSingleTaskResp struct {
	Base
	Data []FixedTimeSingleTaskData `json:"data"`
	Page PageBase                  `json:"page"`
}

type QueryRealTimeSingleTaskReq struct {
	Filter FilterBase `json:"filter" validate:"required"`
	Page   PageBase   `json:"page" validate:"required"`
}

type QueryRealTimeSingleTaskResp struct {
	Base
	Data []RealTimeSingleTaskData `json:"data"`
	Page PageBase                 `json:"page"`
}

type RealTimeSingleTask struct {
	Type     int64  `json:"type" validate:"required,min=1,max=3"`
	BizCode  string `json:"biz_code" validate:"required"`
	BizId    string `json:"biz_id,optional"`
	ExecPath string `json:"exec_path" validate:"required"`
	Param    string `json:"param" validate:"required"`
	Timeout  int64  `json:"timeout" validate:"required,min=5"`
	ExtInfo  string `json:"ext_info,optional"`
}

type RealTimeSingleTaskData struct {
	BaseData
	RealTimeSingleTask
	StartTime  string `json:"start_time"`
	FinishTime string `json:"finish_time"`
	Interval   int64  `json:"interval"`
	ResultMsg  string `json:"result_msg"`
}

type Response struct {
	Base
}

type StartOrStopCronCycleTaskReq struct {
	Id      string `json:"id" validate:"required"`
	IsStart bool   `json:"is_start"`
}
