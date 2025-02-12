package croncycletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/model/croncycletasks"
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
	// 1、格式校验
	if req.Type != int64(core.CronCycleTask) {
		err = errors.New("type is not cron cycle task")
		logx.Error(err)
		return
	}
	if req.Timeout <= 0 {
		req.Timeout = command.DefaultTimeout
	}

	// 2、检查是否已添加过定时任务
	err = cron.CheckCronCycle(ctx, req)
	if err != nil {
		logx.Error(err)
		return
	}
	// 3、添加任务（入库）
	taskId, err = cron.AddDataToCronCycleTask(ctx, 0, req)
	if err != nil {
		logx.Error(err)
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
	// 1、格式检查
	if req.Id == "" {
		err = errors.New("id is empty")
		logx.Error(err)
		return
	}

	// 2、删除定时任务
	entryId, err = cron.DelDataFromCronCycleTask(ctx, req.Id)
	if err != nil {
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
	// 1、格式检查
	if req.Id == "" {
		err = errors.New("id is nil")
		logx.Error(err)
		return
	}

	// 2、修改定时任务
	entryId, err = cron.ModDataFromCronCycleTask(ctx, req)
	if err != nil {
		return
	}

	return
}

func QueryCronCycle(ctx context.Context, req *common.QueryCronCycleTaskReq) (resp *common.QueryTaskConfigResp) {
	var err error
	resp = &common.QueryTaskConfigResp{}
	defer func() {
		if err != nil {
			resp.Code = core.CronCycleError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.CronCycleError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = core.Success.Msg
		}
	}()

	// 1、查询定时任务
	var results []*croncycletasks.TCronCycleTasks
	results, err = cron.QueryDataFromCronCycleTask(ctx, req)
	if err != nil {
		logx.Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, common.CronCycleTaskData{
			BaseData: common.BaseData{
				Id:     v.Id,
				Status: v.Status,
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
	return
}
