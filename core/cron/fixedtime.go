package cron

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/model/crontasks"
	"dmc-task/model/jobsflow"
	"dmc-task/server"
	"dmc-task/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"
)

// 从DB中添加固定时间任务
func addFixedTimeSingleTasksFromDB() {
	ctx := context.Background()
	mc := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	results, err := mc.GetCronTasksByStatus(ctx, int64(core.Init), core.FixedScanTime)
	if err != nil {
		if errors.Is(err, crontasks.ErrNotFound) {
			logx.Debugf("[addFixTimeSingleTasks] task not found")
			return
		}
		logx.Error(err)
		return
	}

	n := len(results)
	for i, v := range results {
		// 1、校验任务是否已经pending（已经pending的job不再pending）
		_, err = mj.GetJobsFlowByCronTaskId(ctx, v.Id)
		if err != nil && errors.Is(err, jobsflow.ErrNotFound) {
			// 2、将固定定时任务pending到job中
			err = pendingFixedTimeSingleTaskFromDB(v)
			if err != nil {
				logx.Error(err)
				continue
			}
			logx.Debugf("[add FixTimeSingleTasks from db] [%d-%d] %+v", n, i+1, v)
		}
	}
}

func pendingFixedTimeSingleTaskFromDB(cronTask *crontasks.TCronTasks) error {
	ctx := context.Background()
	// 1、在流水任务中增加执行任务流水
	jobId := uuid.New().String()
	job := jobsflow.TJobsFlow{
		Id:           jobId,
		Type:         cronTask.Type,
		CronTaskId:   cronTask.Id,
		BizCode:      cronTask.BizCode,
		BizId:        cronTask.BizId,
		ExecPath:     cronTask.ExecPath,
		Param:        cronTask.Param,
		Timeout:      cronTask.Timeout,
		StartTime:    sql.NullTime{},
		FinishTime:   sql.NullTime{},
		ExecInterval: 0,
		Status:       int64(core.Pending),
		ResultMsg:    core.GetResult(core.Success.Code, cronTask.BizId, core.TaskStatusMap[core.Pending], core.Pending, nil),
		ExtInfo:      "{}",
	}
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", job.BizCode), logx.Field("id", job.Id), logx.Field("id", jobId))

	// 3、准备任务前：更新任务状态（t_cron_tasks） 和 插入任务（t_jobs_flow）
	mc := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	err := mc.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新 t_cron_tasks 表状态为 Pending
		cronTask.StartTime = sql.NullTime{Time: utils.GetUTCTime(), Valid: true}
		cronTask.Status = int64(core.Pending)
		cronTask.ResultMsg = core.GetResult(core.Success.Code, cronTask.BizId, core.TaskStatusMap[core.Pending], core.Pending, nil)
		queryCron := fmt.Sprintf("update %s set %s where `id` = ?", mc.GetTableName(),
			mc.GetCronTasksRowsWithPlaceHolder())
		_, err := session.ExecCtx(ctx, queryCron, cronTask.Type, cronTask.BizCode, cronTask.BizId, cronTask.ExecPath,
			cronTask.Param, cronTask.Timeout, cronTask.StartTime, cronTask.FinishTime, cronTask.ExecTime,
			cronTask.ExecInterval, cronTask.Status, cronTask.ResultMsg, cronTask.ExtInfo, cronTask.Id)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		// 插入 t_jobs_flow 表
		queryJob := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			mj.GetTableName(), mj.GetJobsFlowRowsExpectAutoSet())
		_, err = session.ExecCtx(ctx, queryJob, job.Id, job.Type, job.CronTaskId, job.BizCode, job.BizId, job.ExecPath,
			job.Param, job.Timeout, job.StartTime, job.FinishTime, job.ExecInterval, job.Status, job.ResultMsg,
			job.ExtInfo)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return err
	}

	// 2、添加到时间轮中
	interval := float64(cronTask.ExecTime.Sub(utils.GetUTCTime()).Seconds())
	timewheel.AddTimer(timewheel.TW_Sec, timewheel.Sec2msFloat64(interval), 1, job,
		func(p interface{}) {
			_ = execFixedTimeTask(p.(jobsflow.TJobsFlow))
		})
	return nil
}

// 执行定时循环任务
func execFixedTimeTask(taskParam jobsflow.TJobsFlow) error {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("biz_code", taskParam.BizCode), logx.Field("biz_id", taskParam.BizId), logx.Field("id", taskParam.Id))
	var err error
	var data []string
	taskParam.StartTime = sql.NullTime{Time: utils.GetUTCTime(), Valid: true}
	defer func() {
		taskParam.FinishTime = sql.NullTime{Time: utils.GetUTCTime(), Valid: true}
		taskParam.ExecInterval = int64(taskParam.FinishTime.Time.Sub(taskParam.StartTime.Time).Seconds())
		if err != nil {
			status := core.Failed
			taskParam.Status = int64(status)
			taskParam.ResultMsg = core.GetResult(core.FixCronError.Code, taskParam.BizId, err.Error(), status, nil)
		} else {
			status := core.Finished
			taskParam.Status = int64(status)
			taskParam.ResultMsg = core.GetResult(core.Success.Code, taskParam.BizId, core.TaskStatusMap[status], status, data)
		}
		err = updateRecord(ctx, taskParam)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
	}()
	// 1、开始任务前：更新任务状态（t_cron_tasks和t_jobs_flow） Running
	status := core.Running
	taskParam.Status = int64(status)
	taskParam.ResultMsg = core.GetResult(core.Success.Code, taskParam.BizId, core.TaskStatusMap[status], status, nil)
	err = updateRecord(ctx, taskParam)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return err
	}

	// 2、调用任务接口
	data, err = command.ExecCommand(ctx, taskParam.Timeout, taskParam.ExecPath, strings.Split(taskParam.Param, " "))
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return err
	}
	// 3、结束任务后：更新任务状态（t_cron_tasks和t_jobs_flow） Failed Finished
	return nil
}

func updateRecord(ctx context.Context, taskParam jobsflow.TJobsFlow) error {
	// 1、开始任务前：更新任务状态（t_cron_tasks和t_jobs_flow）
	mc := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	err := mc.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		rc, err := mc.FindOne(ctx, taskParam.CronTaskId)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		if taskParam.Status == int64(core.Failed) || taskParam.Status == int64(core.Finished) {
			rc.FinishTime = taskParam.FinishTime
		} else {
			rc.StartTime = taskParam.StartTime
		}
		rc.ExecInterval = taskParam.ExecInterval
		rc.Status = taskParam.Status
		rc.ResultMsg = taskParam.ResultMsg
		err = mc.Update(ctx, rc)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		err = mj.Update(ctx, &taskParam)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

///////////////////////////////////////////////////////////////////////
// 固定定时任务相关

// AddDataToCronTasks 向数据库中添加定时任务
func AddDataToCronTasks(ctx context.Context, req *common.AddFixedTimeSingleTaskReq) (taskId string, err error) {
	m := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	if req.ExtInfo == "" {
		req.ExtInfo = "{}"
	}
	taskId = uuid.New().String()
	_, err = m.Insert(context.Background(), &crontasks.TCronTasks{
		Id:         taskId,
		Type:       req.Type,
		BizCode:    req.BizCode,
		BizId:      req.BizId,
		ExecPath:   req.ExecPath,
		Param:      req.Param,
		Timeout:    req.Timeout,
		StartTime:  sql.NullTime{},
		FinishTime: sql.NullTime{},
		ExecTime:   utils.GetTime(req.ExecTime),
		ResultMsg:  "{}",
		ExtInfo:    req.ExtInfo,
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	logx.WithContext(ctx).Infof("[add fixedtime task to db], task is %s", taskId)
	return taskId, nil
}

// DelDataFromCronTasks 从定时任务中删除数据
func DelDataFromCronTasks(ctx context.Context, req *common.DelFixedTimeSingleTaskReq) (err error) {
	m := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	// 1、根据id查找任务是否存在
	result, err := m.FindOne(ctx, req.Id)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、校验时间（执行时间要大于1m前）
	execTime := result.ExecTime
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		err = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		logx.WithContext(ctx).Error(err)
		return
	}
	// 3、删除任务
	err = m.Delete(ctx, req.Id)
	if err != nil {
		logx.Error(err)
		return
	}
	logx.WithContext(ctx).Infof("delete cron task success, task is %+v", result)
	return
}

// QueryDataFromCronTasks 从定时任务中查询数据
func QueryDataFromCronTasks(ctx context.Context, req *common.QueryFixedTimeSingleTaskReq) (total int64, results []*crontasks.TCronTasks, err error) {
	res, err := model.Query[crontasks.TCronTasks](
		ctx,
		crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn).GetTableName(),
		req.Filter,
		req.Page)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return int64(res.Count), res.Data, nil
}
