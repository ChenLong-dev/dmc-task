package grpc

import (
	"context"
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

// ================================================
// 循环定时任务属性

// AddCronCycleTask 添加循环定时任务
func (t *TaskServer) AddCronCycleTask(ctx context.Context, req *protoc.AddCronCycleTaskReq) (resp *protoc.Response, err error) {
	var gerr error
	resp = &protoc.Response{}
	var taskId string
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, gerr.Error()),
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
		gerr = errors.New("type is not cron cycle task")
		logx.WithContext(ctx).Error(gerr)
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
	gerr = cron.CheckCronCycle(ctx, &task)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 3、添加任务（入库）
	taskId, gerr = cron.AddDataToCronCycleTask(ctx, 0, &task)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 4、返回响应
	return
}

// DelCronCycleTask 删除定时循环任务
func (t *TaskServer) DelCronCycleTask(ctx context.Context, req *protoc.DelCronCycleTaskReq) (resp *protoc.Response, err error) {
	var gerr error
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, gerr.Error()),
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
		gerr = errors.New("id is empty")
		logx.WithContext(ctx).Error(gerr)
		return
	}

	// 2、删除定时任务
	entryId, gerr = cron.DelDataFromCronCycleTask(ctx, req.Id)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}

	return
}

// ModCronCycleTask 修改定时循环任务
func (t *TaskServer) ModCronCycleTask(ctx context.Context, req *protoc.ModCronCycleTaskReq) (resp *protoc.Response, err error) {
	var gerr error
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, gerr.Error()),
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
		gerr = errors.New("id is nil")
		logx.WithContext(ctx).Error(gerr)
		return
	}

	// 2、修改定时任务
	entryId, gerr = cron.ModDataFromCronCycleTask(ctx, &common.ModCronCycleTaskReq{
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
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}

	return
}

func (t *TaskServer) StartOrStopCronCycleTask(ctx context.Context, req *protoc.StartOrStopCronCycleTaskReq) (resp *protoc.Response, err error) {
	var gerr error
	resp = &protoc.Response{}
	var entryId int64
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, gerr.Error()),
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
		gerr = errors.New("id is nil")
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、启动或停止定时任务
	entryId, gerr = cron.StartOrStopDataFromCronCycleTask(ctx, &common.StartOrStopCronCycleTaskReq{
		Id:      req.Id,
		IsStart: req.IsStart,
	})
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	return
}

// QueryCronCycleTask 查询定时循环任务
func (t *TaskServer) QueryCronCycleTask(ctx context.Context, req *protoc.QueryCronCycleTaskReq) (resp *protoc.QueryCronCycleTaskResp, err error) {
	logx.WithContext(ctx).Info("QueryCronCycleTask: ", req)
	var gerr error
	resp = &protoc.QueryCronCycleTaskResp{}
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.CronCycleError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.CronCycleError.Msg, gerr.Error()),
			}
		} else {
			resp.Base = &protoc.Base{
				Code: int32(core.Success.Code),
				Msg:  core.Success.Msg,
			}
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("filter", req.Filter))
	if req.Filter == nil || req.Page == nil {
		gerr = errors.New("filter or page is nil")
		logx.WithContext(ctx).Error(gerr)
		return
	}

	// 1、查询定时任务
	r := &common.QueryCronCycleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	//r.Filter.BizId = req.Filter.BizId // 该类型任务没有定时任务ID
	//r.Filter.CronTaskId = req.Filter.CronTaskId // 该类型任务没有定时任务ID
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	total, res, gerr := cron.QueryDataFromCronCycleTask(ctx, r)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、组装
	for _, v := range res {
		resp.Data = append(resp.Data, &protoc.CronCycleTaskData{
			Base: &protoc.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
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
	resp.Page = &protoc.PageBase{}
	resp.Page.Total = total
	resp.Page.Page = r.Page.Page
	resp.Page.PageSize = r.Page.PageSize
	return
}
