package cron

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/model"
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
	results, err := m.GetCronTasks(ctx)
	if err != nil {
		logx.Error(err)
		return err
	}
	n := len(results)
	for i, v := range results {
		if v.Status == int64(core.Stoped) {
			// 任务已暂停，跳过
			continue
		}
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
		ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", v.BizCode), logx.Field("id", v.Id))
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		v.EntryId = id
		v.Status = int64(core.Running)
		err = m.Update(ctx, v)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		time.Sleep(time.Millisecond * 500)
		logx.WithContext(ctx).Debugf("[add CronCycleTasks] [%d-%d] %+v", n, i+1, v)
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
		if v.Status >= int64(core.Init) || v.Status == int64(core.Stoped) {
			// 正式任务和暂停任务，跳过
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
			RemoveTask(ctx, v.EntryId)
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
	ctx := logx.ContextWithFields(context.Background(), logx.Field("biz_code", taskParam.BizCode))
	// 1、查询定时循环任务（获取cron_task_id）
	mc := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	result, err := mc.GetCronTaskByBizCodeAndType(ctx, taskParam.Type, taskParam.BizCode)
	if err != nil {
		logx.WithContext(ctx).Error(err)
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
		logx.WithContext(ctx).Error(err)
		return
	}
	ctx = logx.ContextWithFields(ctx, logx.Field("id", jobId))

	// 3、调用任务接口
	data, err := command.ExecCommand(ctx, taskParam.Timeout, taskParam.ExecPath, strings.Split(taskParam.Param, " "))
	var status int64
	var msg string
	if err != nil {
		logx.WithContext(ctx).Error(err)
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
		logx.WithContext(ctx).Error(err)
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
		logx.WithContext(ctx).Error(err)
		return
	}
	ctx = logx.ContextWithFields(ctx, logx.Field("task_id", result.Id))
	if result != nil {
		err = errors.New("cron cycle task already exists")
		logx.WithContext(ctx).Error(err)
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
	ctx = logx.ContextWithFields(ctx, logx.Field("task_id", taskId))
	_, err = m.Insert(context.Background(), &croncycletasks.TCronCycleTasks{
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
		logx.WithContext(ctx).Error(err)
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
		logx.WithContext(ctx).Error(err)
		return
	}
	entryId = result.EntryId
	RemoveTask(ctx, result.EntryId)
	result.Status = int64(core.Deleted)
	// 3、更新数据库中的任务信息
	err = m.Update(ctx, result)
	if err != nil {
		logx.WithContext(ctx).Error(err)
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
		logx.WithContext(ctx).Error(err)
		return
	}
	entryId = result.EntryId
	// 2、检查参数
	if req.Type <= 0 || req.Type != int64(core.CronCycleTask) {
		err = errors.New("type is not cron cycle task")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.BizCode != result.BizCode || req.Type != result.Type {
		err = errors.New("biz code or type is not match")
		logx.WithContext(ctx).Error(err)
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
	result.Status = int64(core.Modified)
	err = m.Update(ctx, result)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 4、返回修改的任务id
	return
}

// StartOrStopDataFromCronCycleTask 启停定时任务
func StartOrStopDataFromCronCycleTask(ctx context.Context, req *common.StartOrStopCronCycleTaskReq) (entryId int64, err error) {
	m := croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn)
	// 1、通过id从数据库中获取任务信息
	result, err := m.GetCronTaskById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, croncycletasks.ErrNotFound) {
			return 0, errors.New("task not found")
		}
		logx.WithContext(ctx).Error(err)
		return
	}
	entryId = result.EntryId
	// 2、检查参数：
	if result.Status == int64(core.Stoped) { // 2.1 如果数据库记录是已暂停的定时任务
		// 请求是暂停
		if req.IsStart == false {
			logx.WithContext(ctx).Info("task is already stopped!")
			return
		}
		// 请求开始，更新数据库状态为添加
		result.Status = int64(core.Added)
		err = m.Update(ctx, result)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
		logx.WithContext(ctx).Info("add again cron cycle task!")
	} else if result.Status >= int64(core.Init) { // 2.2 如果数据库记录的是正常状态的定时任务
		// 请求开始
		if req.IsStart == true {
			logx.WithContext(ctx).Info("task is already started!")
			return
		}
		// 请求暂停，移除任务(不删数据库)，并更新数据库状态为暂停
		RemoveTask(ctx, result.EntryId)
		result.Status = int64(core.Stoped)
		err = m.Update(ctx, result)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
		logx.WithContext(ctx).Infof("stopped cron cycle task! entryId:%d", result.EntryId)
	} else {
		err = errors.New("task status is not valid")
		logx.WithContext(ctx).Error(err)
		return
	}

	return
}

// QueryDataFromCronCycleTask 从定时任务中查询数据
func QueryDataFromCronCycleTask(ctx context.Context, req *common.QueryCronCycleTaskReq) (total int64, results []*croncycletasks.TCronCycleTasks, err error) {
	res, err := model.Query[croncycletasks.TCronCycleTasks](
		ctx,
		croncycletasks.NewTCronCycleTasksModel(*server.SvrCtx.MysqlConn).GetTableName(),
		req.Filter,
		req.Page)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return int64(res.Count), res.Data, nil
}
