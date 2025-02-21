package cron

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/model/jobsflow"
	"dmc-task/server"
	"dmc-task/utils"
	"errors"
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
		Id:         taskId,
		Type:       taskParam.Type,
		BizCode:    taskParam.BizCode,
		BizId:      taskParam.BizId,
		ExecPath:   taskParam.ExecPath,
		Param:      taskParam.Param,
		Timeout:    taskParam.Timeout,
		StartTime:  sql.NullTime{Time: utils.GetUTCTime(), Valid: true},
		FinishTime: sql.NullTime{},
		ResultMsg:  "{}",
		ExtInfo:    taskParam.ExtInfo,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
	defer cancel()
	go func() {
		_ = execRealtime(job)
	}()

	// 插入t_jobs_flow表
	mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	_, err := mj.Insert(ctx, &job)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return "", nil
	}

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
	return nil
}

///////////////////////////////////////////////////////////////////////
// 实时单任务相关

// QueryDataFromJobsFlow 从JobsFlow表中查询数据
func QueryDataFromJobsFlow(ctx context.Context, req *common.QueryRealTimeSingleTaskReq) (results []*jobsflow.TJobsFlow, err error) {
	m := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	if req.Id != "" {
		var result = &jobsflow.TJobsFlow{}
		result, err = m.FindOne(ctx, req.Id)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return
		}
		results = append(results, result)
		return
	}
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	if req.Status < 0 || req.Status > int64(core.Finished) {
		err = errors.New("status is error")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.TimeHorizon == 0 {
		req.TimeHorizon = core.DefaultTimeHorizon
	}
	if req.Limit == 0 {
		req.Limit = core.DefaultLimit
	}

	results, err = m.GetJobsFlowByStatus2(ctx, req.Status, req.TimeHorizon, req.Limit)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return
}
