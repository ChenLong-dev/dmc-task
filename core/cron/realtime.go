package cron

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/model"
	"dmc-task/model/jobsflow"
	"dmc-task/server"
	"dmc-task/utils"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

func AddRealTimeTask(ctx context.Context, taskParam common.RealTimeSingleTask) (string, error) {
	return addRealTimeTask(ctx, taskParam)
}

func addRealTimeTask(ctx context.Context, taskParam common.RealTimeSingleTask) (string, error) {
	// 1、在流水任务中增加执行任务流水
	if taskParam.ExtInfo == "" {
		taskParam.ExtInfo = "{}"
	}
	taskId := uuid.New().String()
	job := jobsflow.TJobsFlow{
		Id:           taskId,
		Type:         taskParam.Type,
		BizCode:      taskParam.BizCode,
		BizId:        taskParam.BizId,
		ExecPath:     taskParam.ExecPath,
		Param:        taskParam.Param,
		Timeout:      taskParam.Timeout,
		StartTime:    sql.NullTime{Time: utils.GetUTCTime(), Valid: true},
		FinishTime:   sql.NullTime{},
		ExecInterval: 0,
		Status:       int64(core.Running),
		ResultMsg:    "{}",
		ExtInfo:      taskParam.ExtInfo,
	}
	ctx = logx.ContextWithFields(ctx, logx.Field("id", taskId))

	// 插入t_jobs_flow表
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	_, err := mj.Insert(context.Background(), &job)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return "", nil
	}
	logx.WithContext(ctx).Debugf("insert job flow is success! %d", job.Status)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
	defer cancel()
	go func() {
		_ = execRealtime(job)
	}()

	return taskId, nil
}

func execRealtime(job jobsflow.TJobsFlow) error {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("id", job.Id),
		logx.Field("biz_code", job.BizCode), logx.Field("biz_id", job.BizId))

	// 1、执行命令
	var data []string
	var status int64
	var msg string
	data, err := command.ExecCommand(ctx, job.Timeout, job.ExecPath, strings.Split(job.Param, " "))
	if err != nil {
		logx.WithContext(ctx).Error(err)
		status = int64(core.Failed)
		msg = err.Error()
	} else {
		status = int64(core.Finished)
		msg = core.TaskStatusMap[core.Finished]
	}

	// 2、更新任务状态
	job.Status = status
	job.FinishTime = sql.NullTime{Time: utils.GetUTCTime(), Valid: true}
	job.ExecInterval = int64(job.FinishTime.Time.Sub(job.StartTime.Time).Seconds())
	job.ResultMsg = core.GetResult(core.Success.Code, "", msg, core.TaskStatus(status), data)
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	err = mj.Update(ctx, &job)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return nil
	}
	logx.WithContext(ctx).Debugf("update job flow is success! %d", job.Status)
	return nil
}

///////////////////////////////////////////////////////////////////////
// 实时单任务相关

// QueryDataFromJobsFlow 从JobsFlow表中查询数据
func QueryDataFromJobsFlow(ctx context.Context, req *common.QueryRealTimeSingleTaskReq) (total int64, results []*jobsflow.TJobsFlow, err error) {
	res, err := model.Query[jobsflow.TJobsFlow](
		ctx,
		jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn).GetTableName(),
		req.Filter,
		req.Page)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return int64(res.Count), res.Data, nil
}
