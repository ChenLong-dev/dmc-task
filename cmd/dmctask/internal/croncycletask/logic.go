package croncycletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/utils"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func AddCronCycle(ctx context.Context, req *common.AddCronCycleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId)
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.BizCode))

	// 1、格式校验
	if req.Type != int64(core.CronCycleTask) {
		err = errors.New("type is not cron cycle task")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.Timeout <= 0 {
		req.Timeout = command.DefaultTimeout
	}

	// 2、检查是否已添加过定时任务
	err = cron.CheckCronCycle(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 3、添加任务（入库）
	taskId, err = cron.AddDataToCronCycleTask(ctx, 0, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 4、返回响应
	return
}

func DelCronCycle(ctx context.Context, req *common.DelCronCycleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId)
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

func ModCronCycle(ctx context.Context, req *common.ModCronCycleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId)
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
	entryId, err = cron.ModDataFromCronCycleTask(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	return
}

func StartOrStopCronCycle(ctx context.Context, req *common.StartOrStopCronCycleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var entryId int64
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, entry id is %d", core.Success.Msg, entryId)
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
	entryId, err = cron.StartOrStopDataFromCronCycleTask(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return
}

func QueryCronCycle(ctx context.Context, req *common.QueryCronCycleTaskReq) (resp *common.QueryCronCycleTaskResp) {
	var err error
	resp = &common.QueryCronCycleTaskResp{}
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = core.Success.Msg
		}
	}()
	ctx = logx.ContextWithFields(ctx, logx.Field("filter", req.Filter))

	// 1、格式检查
	if req.Filter.TimeType != "" && req.Filter.TimeType != "create_time" && req.Filter.TimeType != "update_time" {
		err = errors.New("time_type is not create_time or update_time")
		logx.WithContext(ctx).Error(err)
		return
	}

	// 1、查询定时任务
	total, results, err := cron.QueryDataFromCronCycleTask(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, common.CronCycleTaskData{
			BaseData: common.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
			},
			CronCycleTask: common.CronCycleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				Cron:     v.Cron,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				Timeout:  v.Timeout,
				ExtInfo:  v.ExtInfo,
			},
		})
	}
	resp.Page = req.Page
	resp.Page.Total = total
	return
}
