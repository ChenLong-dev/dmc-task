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
// 实时单任务属性

// AddRealTimeSingleTask 添加实时单任务
func (t *TaskServer) AddRealTimeSingleTask(ctx context.Context, req *protoc.AddRealTimeSingleTaskReq) (resp *protoc.Response, err error) {
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
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.Task.BizCode), logx.Field("biz_id", req.Task.BizId))

	// 1、格式校验
	if req.Task.Type != int64(core.RealTimeSingleTask) {
		gerr = errors.New("type is not real time task")
		logx.WithContext(ctx).Error(gerr)
		return
	}
	if req.Task.Timeout <= 0 {
		req.Task.Timeout = command.DefaultTimeout
	}

	// 2、入库+执行+更新
	taskId, gerr = cron.AddRealTimeTask(ctx, common.RealTimeSingleTask{
		Type:     req.Task.Type,
		BizCode:  req.Task.BizCode,
		BizId:    req.Task.BizId,
		ExecPath: req.Task.ExecPath,
		Param:    req.Task.Param,
		Timeout:  int64(req.Task.Timeout),
		ExtInfo:  req.Task.ExtInfo,
	})
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}

	// 3、返回响应
	return
}

// QueryRealTimeSingleTask 查询实时单任务
func (t *TaskServer) QueryRealTimeSingleTask(ctx context.Context, req *protoc.QueryRealTimeSingleTaskReq) (resp *protoc.QueryRealTimeSingleTaskResp, err error) {
	var gerr error
	resp = &protoc.QueryRealTimeSingleTaskResp{}
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.JobError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.JobError.Msg, gerr.Error()),
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

	//// 1、查询数据库
	r := &common.QueryRealTimeSingleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	r.Filter.BizId = req.Filter.BizId
	r.Filter.CronTaskId = req.Filter.CronTaskId
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	total, res, gerr := cron.QueryDataFromJobsFlow(ctx, r)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、组装
	for _, v := range res {
		resp.Data = append(resp.Data, &protoc.RealTimeSingleTaskData{
			Base: &protoc.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
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
	resp.Page = &protoc.PageBase{}
	resp.Page.Total = total
	resp.Page.Page = r.Page.Page
	resp.Page.PageSize = r.Page.PageSize
	return
}
