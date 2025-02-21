package grpc

import (
	"context"
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/model/croncycletasks"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TaskServer struct {
	protoc.UnimplementedTaskServer
}

// ================================================
// 循环定时任务属性

// AddCronCycleTask 添加循环定时任务
func (t *TaskServer) AddCronCycleTask(ctx context.Context, req *protoc.AddCronCycleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.Task.BizCode))

	// 1、格式校验
	if req.Task.Type != int64(core.CronCycleTask) {
		err = errors.New("type is not cron cycle task")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.Task.Timeout <= 0 {
		req.Task.Timeout = command.DefaultTimeout
	}

	// 2、检查是否已添加过定时任务
	task := common.AddCronCycleTaskReq{
		CronCycleTask: common.CronCycleTask{
			Type:     req.Task.Type,
			BizCode:  req.Task.BizCode,
			Cron:     req.Task.Cron,
			ExecPath: req.Task.ExecPath,
			Param:    req.Task.Param,
			Timeout:  int64(req.Task.Timeout),
			ExtInfo:  req.Task.ExtInfo,
		},
	}
	err = cron.CheckCronCycle(ctx, &task)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 3、添加任务（入库）
	taskId, err = cron.AddDataToCronCycleTask(ctx, 0, &task)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 4、返回响应
	return
}

// DelCronCycleTask 删除定时循环任务
func (t *TaskServer) DelCronCycleTask(ctx context.Context, req *protoc.DelCronCycleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、格式检查
	if req.Id == "" {
		err = errors.New("id is empty")
		logx.WithContext(ctx).Error(err)
		return
	}

	// 2、删除定时任务
	entryId, err = cron.DelDataFromCronCycleTask(ctx, req.Id)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	return
}

// ModCronCycleTask 修改定时循环任务
func (t *TaskServer) ModCronCycleTask(ctx context.Context, req *protoc.ModCronCycleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、格式检查
	if req.Id == "" {
		err = errors.New("id is nil")
		logx.WithContext(ctx).Error(err)
		return
	}

	// 2、修改定时任务
	entryId, err = cron.ModDataFromCronCycleTask(ctx, &common.ModCronCycleTaskReq{
		Id: req.Id,
		CronCycleTask: common.CronCycleTask{
			Type:     req.Task.Type,
			BizCode:  req.Task.BizCode,
			Cron:     req.Task.Cron,
			ExecPath: req.Task.ExecPath,
			Param:    req.Task.Param,
			Timeout:  int64(req.Task.Timeout),
			ExtInfo:  req.Task.ExtInfo,
		},
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	return
}

func (t *TaskServer) StartOrStopCronCycleTask(ctx context.Context, req *protoc.StartOrStopCronCycleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、格式检查
	if req.Id == "" {
		err = errors.New("id is nil")
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、启动或停止定时任务
	entryId, err = cron.StartOrStopDataFromCronCycleTask(ctx, &common.StartOrStopCronCycleTaskReq{
		Id:      req.Id,
		IsStart: req.IsStart,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return
}

// QueryCronCycleTask 查询定时循环任务
func (t *TaskServer) QueryCronCycleTask(ctx context.Context, req *protoc.QueryCronCycleTaskReq) (resp *protoc.QueryCronCycleTaskResp, err error) {
	resp = &protoc.QueryCronCycleTaskResp{}
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  core.Success.Msg,
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、查询定时任务
	var results []*croncycletasks.TCronCycleTasks
	results, err = cron.QueryDataFromCronCycleTask(ctx, &common.QueryCronCycleTaskReq{
		Id: req.Id,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, &protoc.CronCycleTaskData{
			Base: &protoc.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			Task: &protoc.CronCycleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				Cron:     v.Cron,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				Timeout:  int32(v.Timeout),
				ExtInfo:  v.ExtInfo,
			},
		})
	}
	return
}

// ================================================
// 固定时间单任务属性

// AddFixedTimeSingleTask 添加固定时间单任务
func (t *TaskServer) AddFixedTimeSingleTask(ctx context.Context, req *protoc.AddFixedTimeSingleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.Task.BizCode), logx.Field("biz_id", req.Task.BizId))

	// 1、格式校验
	if req.Task.Type != int64(core.FixedTimeSingleTask) {
		err = errors.New("type is not fixed time task")
		logx.WithContext(ctx).Error(err)
		return
	}

	execTime := utils.GetTime(req.Task.ExecTime)
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		err = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、添加任务（入库）
	taskId, err = cron.AddDataToCronTasks(ctx, &common.AddFixedTimeSingleTaskReq{
		FixedTimeSingleTask: common.FixedTimeSingleTask{
			Type:     req.Task.Type,
			BizCode:  req.Task.BizCode,
			BizId:    req.Task.BizId,
			ExecPath: req.Task.ExecPath,
			ExecTime: req.Task.ExecTime,
			Param:    req.Task.Param,
			Timeout:  int64(req.Task.Timeout),
			ExtInfo:  req.Task.ExtInfo,
		},
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 3、返回响应
	return
}

// DelFixedTimeSingleTask 删除固定时间单任务
func (t *TaskServer) DelFixedTimeSingleTask(ctx context.Context, req *protoc.DelFixedTimeSingleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.FixCronError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.FixCronError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  core.Success.Msg,
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、格式校验
	if req.Id == "" {
		err = fmt.Errorf("task id is empty")
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、删除任务
	err = cron.DelDataFromCronTasks(ctx, &common.DelFixedTimeSingleTaskReq{
		Id: req.Id,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return
}

// QueryFixedTimeSingleTask 查询固定时间单任务
func (t *TaskServer) QueryFixedTimeSingleTask(ctx context.Context, req *protoc.QueryFixedTimeSingleTaskReq) (resp *protoc.QueryFixedTimeSingleTaskResp, err error) {
	resp = &protoc.QueryFixedTimeSingleTaskResp{}
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.FixCronError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.FixCronError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  core.Success.Msg,
			}
		}
	}()
	// 1、查询数据库
	results, err := cron.QueryDataFromCronTasks(ctx, &common.QueryFixedTimeSingleTaskReq{
		Id:          req.Id,
		Status:      req.Status,
		TimeHorizon: req.TimeHorizon,
		Limit:       req.Limit,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, &protoc.FixedTimeSingleTaskData{
			Base: &protoc.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			Task: &protoc.FixedTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				ExecTime: utils.GetTimestamp(v.ExecTime),
				Param:    v.Param,
				Timeout:  int32(v.Timeout),
				ExtInfo:  v.ExtInfo,
			},
			StartTime:  utils.GetTimeStr(v.StartTime.Time),
			FinishTime: utils.GetTimeStr(v.FinishTime.Time),
			Interval:   v.ExecInterval,
			ResultMsg:  v.ResultMsg,
		})
	}
	return
}

// ================================================
// 实时单任务属性

// AddRealTimeSingleTask 添加实时单任务
func (t *TaskServer) AddRealTimeSingleTask(ctx context.Context, req *protoc.AddRealTimeSingleTaskReq) (resp *protoc.Response, err error) {
	resp = &protoc.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId),
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.Task.BizCode), logx.Field("biz_id", req.Task.BizId))

	// 1、格式校验
	if req.Task.Type != int64(core.RealTimeSingleTask) {
		err = errors.New("type is not real time task")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.Task.Timeout <= 0 {
		req.Task.Timeout = command.DefaultTimeout
	}

	// 2、入库+执行+更新
	taskId, err = cron.AddRealTimeTask(ctx, common.RealTimeSingleTask{
		Type:     req.Task.Type,
		BizCode:  req.Task.BizCode,
		BizId:    req.Task.BizId,
		ExecPath: req.Task.ExecPath,
		Param:    req.Task.Param,
		Timeout:  int64(req.Task.Timeout),
		ExtInfo:  req.Task.ExtInfo,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	// 3、返回响应
	return
}

// QueryRealTimeSingleTask 查询实时单任务
func (t *TaskServer) QueryRealTimeSingleTask(ctx context.Context, req *protoc.QueryRealTimeSingleTaskReq) (resp *protoc.QueryRealTimeSingleTaskResp, err error) {
	resp = &protoc.QueryRealTimeSingleTaskResp{}
	defer func() {
		if err != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  core.Success.Msg,
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、查询数据库
	results, err := cron.QueryDataFromJobsFlow(ctx, &common.QueryRealTimeSingleTaskReq{
		Id:          req.Id,
		Status:      req.Status,
		TimeHorizon: req.TimeHorizon,
		Limit:       req.Limit,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, &protoc.RealTimeSingleTaskData{
			Base: &protoc.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			Task: &protoc.RealTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				Timeout:  int32(v.Timeout),
				ExtInfo:  v.ExtInfo,
			},
			StartTime:  utils.GetTimeStr(v.StartTime.Time),
			FinishTime: utils.GetTimeStr(v.FinishTime.Time),
			Interval:   v.ExecInterval,
			ResultMsg:  v.ResultMsg,
		})
	}
	return
}
