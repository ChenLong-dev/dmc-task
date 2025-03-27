package grpc

import (
	"context"
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// ================================================
// 固定时间单任务属性

// AddFixedTimeSingleTask 添加固定时间单任务
func (t *TaskServer) AddFixedTimeSingleTask(ctx context.Context, req *protoc.AddFixedTimeSingleTaskReq) (resp *protoc.Response, err error) {
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
	if req.Task.Type != int64(core.FixedTimeSingleTask) {
		gerr = errors.New("type is not fixed time task")
		logx.WithContext(ctx).Error(gerr)
		return
	}

	execTime := utils.GetTime(req.Task.ExecTime)
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		gerr = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、添加任务（入库）
	taskId, gerr = cron.AddDataToCronTasks(ctx, &common.AddFixedTimeSingleTaskReq{
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
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 3、返回响应
	return
}

// DelFixedTimeSingleTask 删除固定时间单任务
func (t *TaskServer) DelFixedTimeSingleTask(ctx context.Context, req *protoc.DelFixedTimeSingleTaskReq) (resp *protoc.Response, err error) {
	var gerr error
	resp = &protoc.Response{}
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.FixCronError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.FixCronError.Msg, gerr.Error()),
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
		gerr = fmt.Errorf("task id is empty")
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、删除任务
	gerr = cron.DelDataFromCronTasks(ctx, &common.DelFixedTimeSingleTaskReq{
		Id: req.Id,
	})
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	return
}

// QueryFixedTimeSingleTask 查询固定时间单任务
func (t *TaskServer) QueryFixedTimeSingleTask(ctx context.Context, req *protoc.QueryFixedTimeSingleTaskReq) (resp *protoc.QueryFixedTimeSingleTaskResp, err error) {
	var gerr error
	resp = &protoc.QueryFixedTimeSingleTaskResp{}
	defer func() {
		if gerr != nil {
			resp.Base = &protoc.Base{
				Code: int32(core.FixCronError.Code),
				Msg:  fmt.Sprintf("%s: %s", core.FixCronError.Msg, gerr.Error()),
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

	// 1、查询数据库
	r := &common.QueryFixedTimeSingleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	r.Filter.BizId = req.Filter.BizId
	//r.Filter.CronTaskId = req.Filter.CronTaskId // 该类型任务没有定时任务ID
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	total, res, gerr := cron.QueryDataFromCronTasks(ctx, r)
	if gerr != nil {
		logx.WithContext(ctx).Error(gerr)
		return
	}
	// 2、组装
	for _, v := range res {
		resp.Data = append(resp.Data, &protoc.FixedTimeSingleTaskData{
			Base: &protoc.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
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
	resp.Page = &protoc.PageBase{}
	resp.Page.Total = total
	resp.Page.Page = r.Page.Page
	resp.Page.PageSize = r.Page.PageSize
	return
}
