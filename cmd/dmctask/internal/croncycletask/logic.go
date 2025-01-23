package croncycletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/model/croncycletasks"
	"dmc-task/server"
	"fmt"
	"github.com/google/uuid"
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
	err = checkCronCycle(ctx, req)
	if err != nil {
		logx.Error(err)
		return
	}
	// 3、添加任务（入库）
	taskId, err = addDataToDB(ctx, 0, req)
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
	entryId, err = removeCronCycleTask(ctx, req.Id)
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
	entryId, err = modCronCycleTask(ctx, req)
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
	results, err = queryCronCycleTask(ctx, req)
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

// ----------------- 私有函数 -----------------
func checkCronCycle(ctx context.Context, req *common.AddCronCycleTaskReq) (err error) {
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	result, err := m.GetCronTaskByBizCodeAndType(ctx, req.Type, req.BizCode)
	if err != nil {
		if errors.Is(err, croncycletasks.ErrNotFound) {
			return nil
		}
		logx.Error(err)
		return
	}
	if result != nil {
		err = errors.New("cron cycle task already exists")
		logx.Error(err)
		return
	}
	return
}

func addDataToDB(ctx context.Context, id int64, req *common.AddCronCycleTaskReq) (taskId string, err error) {
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	if req.ExtInfo == "" {
		req.ExtInfo = "{}"
	}
	taskId = uuid.New().String()
	_, err = m.Insert(ctx, &croncycletasks.TCronCycleTasks{
		Id:       taskId,
		EntryId:  id,
		Type:     req.Type,
		BizCode:  req.BizCode,
		Cron:     req.Cron,
		ExecPath: req.ExecPath,
		Param:    req.Param,
		Status:   int64(core.Added),
		Timeout:  req.Timeout,
		ExtInfo:  req.ExtInfo,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	return
}

func removeCronCycleTask(ctx context.Context, id string) (entryId int64, err error) {
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	// 1、通过id从数据库中获取任务信息
	result, err := m.GetCronTaskById(ctx, id)
	if err != nil {
		if errors.Is(err, croncycletasks.ErrNotFound) {
			return entryId, errors.New("task not found")
		}
		logx.Error(err)
		return
	}
	entryId = result.EntryId
	cron.RemoveTask(result.EntryId)
	result.Status = int64(core.Deleted)
	// 3、更新数据库中的任务信息
	err = m.Update(ctx, result)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}

func modCronCycleTask(ctx context.Context, req *common.ModCronCycleTaskReq) (entryId int64, err error) {
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	// 1、通过id从数据库中获取任务信息
	result, err := m.GetCronTaskById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, croncycletasks.ErrNotFound) {
			return 0, errors.New("task not found")
		}
		logx.Error(err)
		return
	}
	entryId = result.EntryId
	// 2、检查参数
	if req.Type <= 0 || req.Type != int64(core.CronCycleTask) {
		err = errors.New("type is not cron cycle task")
		logx.Error(err)
		return
	}
	if req.BizCode != result.BizCode || req.Type != result.Type {
		err = errors.New("biz code or type is not match")
		logx.Error(err)
		return
	}
	if req.Cron != "" {
		result.Cron = req.Cron
	}
	if req.ExecPath != "" {
		result.ExecPath = req.ExecPath
	}
	if req.Param != "" {
		result.Param = req.Param
	}
	if req.Timeout != 0 {
		if req.Timeout < 0 {
			result.Timeout = command.DefaultTimeout
		} else {
			result.Timeout = req.Timeout
		}
	}
	if req.ExtInfo != "" {
		result.ExtInfo = req.ExtInfo
	}

	// 3、修改数据库中的任务信息
	err = m.Update(ctx, &croncycletasks.TCronCycleTasks{
		Id:       result.Id,
		EntryId:  entryId,
		Type:     result.Type,
		BizCode:  result.BizCode,
		Cron:     result.Cron,
		ExecPath: result.ExecPath,
		Param:    result.Param,
		Timeout:  result.Timeout,
		Status:   int64(core.Modified),
		ExtInfo:  result.ExtInfo,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	// 4、返回修改的任务id
	return
}

func queryCronCycleTask(ctx context.Context, req *common.QueryCronCycleTaskReq) (results []*croncycletasks.TCronCycleTasks, err error) {
	logx.Debugf("queryCronCycleTask req: %v", req.Id)
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	if req.Id == "" {
		results, err = m.GetCronTasks(ctx)
		if err != nil {
			if errors.Is(err, croncycletasks.ErrNotFound) {
				return nil, errors.New("tasks not found")
			}
			logx.Error(err)
			return
		}
	} else {
		var result *croncycletasks.TCronCycleTasks
		result, err = m.GetCronTaskById(ctx, req.Id)
		if err != nil {
			if errors.Is(err, croncycletasks.ErrNotFound) {
				return nil, errors.New("task not found")
			}
			logx.Error(err)
			return
		}
		results = append(results, result)
	}
	return
}
