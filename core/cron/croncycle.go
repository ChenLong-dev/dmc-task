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
