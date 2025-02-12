package cron

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/model/croncycletasks"
	"dmc-task/model/crontasks"
	"dmc-task/model/jobsflow"
	"dmc-task/server"
	"dmc-task/utils"
	"errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

// 添加定时循环任务
func addCronCycleInitTasks() error {
	logx.Debug("添加定时循环初始任务 1 ...")
	ctx := context.Background()
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	results, _ := m.GetCronTasks(ctx)
	n := len(results)
	for i, v := range results {
		if v.Status == int64(core.Deleted) {
			// 任务已删除，跳过
			_ = m.Delete(ctx, v.Id)
			continue
		}
		id, err := addDynamicTask(common.CronCycleTask{
			Type:     v.Type,
			BizCode:  v.BizCode,
			Cron:     v.Cron,
			ExecPath: v.ExecPath,
			Param:    v.Param,
			Timeout:  v.Timeout,
			ExtInfo:  v.ExtInfo,
		})
		if err != nil {
			logx.Error(err)
			return err
		}
		v.EntryId = id
		v.Status = int64(core.Running)
		err = m.Update(ctx, v)
		if err != nil {
			logx.Error(err)
			return err
		}
		time.Sleep(time.Millisecond * 500)
		logx.Debugf("[add CronCycleTasks] [%d-%d] %+v", n, i+1, v)
	}
	return nil
}

func addCronCycleTasks() {
	logx.Debug("循环添加定时循环任务 2 ....")
	ctx := context.Background()
	mc := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	results, err := mc.GetCronTasks(ctx)
	if err != nil {
		if errors.Is(err, crontasks.ErrNotFound) {
			logx.Debugf("[addCronCycleTasks] task not found")
			return
		}
		logx.Error(err)
		return
	}
	n := len(results)
	for i, v := range results {
		if v.Status >= int64(core.Init) {
			// 跳过
			continue
		}
		if v.Status == int64(core.Deleted) {
			// 任务删除任务，跳过
			_ = mc.Delete(ctx, v.Id)
			logx.Infof("[CronCycleTasks] [%d-%d] deleted id:%s", n, i+1, v.Id)
			continue
		}
		if v.Status == int64(core.Modified) {
			// 任务修改，移除任务，重新添加
			RemoveTask(v.EntryId)
			logx.Infof("[CronCycleTasks] [%d-%d] modified id:%s, entryID:%d", n, i+1, v.Id, v.EntryId)
		}
		if v.Status == int64(core.Added) {
			// 任务正在运行，跳过
			logx.Debugf("[CronCycleTasks] [%d-%d] added id:%s, entryID:%d", n, i+1, v.Id, v.EntryId)
		}
		entryID, err := addDynamicTask(common.CronCycleTask{
			Type:     v.Type,
			BizCode:  v.BizCode,
			Cron:     v.Cron,
			ExecPath: v.ExecPath,
			Param:    v.Param,
			Timeout:  v.Timeout,
			ExtInfo:  v.ExtInfo,
		})
		if err != nil {
			logx.Error(err)
			continue
		}
		v.EntryId = entryID
		v.Status = int64(core.Running)
		_ = mc.Update(ctx, v)
		logx.Infof("[CronCycleTasks] [%d-%d] add id:%s, entryID:%d", n, i+1, v.Id, entryID)
		continue
	}

	return
}

// 执行定时循环任务
func execCronCycle(taskParam common.CronCycleTask) {
	ctx := context.Background()
	// 1、查询定时循环任务（获取cron_task_id）
	mc := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	result, err := mc.GetCronTaskByBizCodeAndType(ctx, taskParam.Type, taskParam.BizCode)
	if err != nil {
		logx.Error(err)
		return
	}

	// 2、在流水任务中增加执行任务流水
	jobId := uuid.New().String()
	job := &jobsflow.TJobsFlow{
		Id:           jobId,
		Type:         taskParam.Type,
		CronTaskId:   result.Id,
		BizCode:      taskParam.BizCode,
		ExecPath:     taskParam.ExecPath,
		Param:        taskParam.Param,
		Timeout:      taskParam.Timeout,
		StartTime:    sql.NullTime{Time: utils.GetUTCTime(), Valid: true},
		FinishTime:   sql.NullTime{Time: utils.GetUTCTime(), Valid: false},
		ExecInterval: 0,
		Status:       int64(core.Running),
		ResultMsg:    core.GetResult(core.Success.Code, "", core.TaskStatusMap[core.Running], core.Running, nil),
		ExtInfo:      "{}",
	}
	// 插入t_jobs_flow表
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	_, err = mj.Insert(ctx, job)
	if err != nil {
		logx.Error(err)
		return
	}

	// 3、调用任务接口
	data, err := command.ExecCommand(ctx, taskParam.Timeout, taskParam.ExecPath, strings.Split(taskParam.Param, " "))
	var status int64
	var msg string
	if err != nil {
		logx.Error(err)
		status = int64(core.Failed)
		msg = err.Error()
	} else {
		status = int64(core.Finished)
		msg = core.TaskStatusMap[core.Finished]
	}

	// 4、更新流水任务中的任务状态
	job.Status = status
	job.FinishTime = sql.NullTime{Time: utils.GetUTCTime(), Valid: true}
	job.ExecInterval = int64(time.Now().Sub(job.StartTime.Time).Seconds())
	job.ResultMsg = core.GetResult(core.Success.Code, "", msg, core.TaskStatus(status), data)
	err = mj.Update(ctx, job)
	if err != nil {
		logx.Error(err)
		return
	}
}

///////////////////////////////////////////////////////////////////////
// 循环定时任务相关

// CheckCronCycle 检查定时任务是否存在
func CheckCronCycle(ctx context.Context, req *common.AddCronCycleTaskReq) (err error) {
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

// AddDataToCronCycleTask 向数据库中添加定时任务
func AddDataToCronCycleTask(ctx context.Context, id int64, req *common.AddCronCycleTaskReq) (taskId string, err error) {
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

// DelDataFromCronCycleTask 从数据库中删除定时任务
func DelDataFromCronCycleTask(ctx context.Context, id string) (entryId int64, err error) {
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
	RemoveTask(result.EntryId)
	result.Status = int64(core.Deleted)
	// 3、更新数据库中的任务信息
	err = m.Update(ctx, result)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}

// ModDataFromCronCycleTask 修改定时任务
func ModDataFromCronCycleTask(ctx context.Context, req *common.ModCronCycleTaskReq) (entryId int64, err error) {
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

// QueryDataFromCronCycleTask 从定时任务中查询数据
func QueryDataFromCronCycleTask(ctx context.Context, req *common.QueryCronCycleTaskReq) (results []*croncycletasks.TCronCycleTasks, err error) {
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
