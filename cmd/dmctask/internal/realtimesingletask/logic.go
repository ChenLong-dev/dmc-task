package realtimesingletask

import (
	"context"
	"database/sql"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/model/jobsflow"
	"dmc-task/server"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultLimit       = 500
	DefaultTimeHorizon = 1
)

func AddJob(ctx context.Context, req *common.AddRealTimeSingleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Code = core.JobError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId)
		}
	}()
	// 1、格式校验
	if req.Type != int64(core.RealTimeSingleTask) {
		err = errors.New("type is not real time task")
		logx.Error(err)
		return
	}
	if req.Timeout <= 0 {
		req.Timeout = command.DefaultTimeout
	}

	// 2、入库+执行+更新
	taskId, err = cron.AddRealTimeTask(ctx, req.RealTimeSingleTask)
	if err != nil {
		logx.Error(err)
		return
	}

	// 3、返回响应
	return
}

func QueryJob(ctx context.Context, req *common.QueryRealTimeSingleTaskReq) (resp *common.QueryRealTimeSingleTaskResp) {
	var err error
	resp = &common.QueryRealTimeSingleTaskResp{}
	defer func() {
		if err != nil {
			resp.Code = core.JobError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.JobError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = core.Success.Msg
		}
	}()

	// 1、查询数据库
	results, err := queryDataFromDB(ctx, req)
	if err != nil {
		logx.Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, common.RealTimeSingleTaskData{
			BaseData: common.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			RealTimeSingleTask: common.RealTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				ExtInfo:  v.ExtInfo,
			},
			StartTime:  utils.GetTimeStr(v.StartTime.Time),
			FinishTime: utils.GetTimeStr(v.FinishTime.Time),
			Interval:   v.ExecInterval,
			ResultMsg:  v.ResultMsg,
		})
	}
	return
}

// ----------------- 私有函数 -----------------

func addDataToDB(ctx context.Context, req *common.AddRealTimeSingleTaskReq) (taskId string, err error) {
	m := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	if req.ExtInfo == "" {
		req.ExtInfo = "{}"
	}
	taskId = uuid.New().String()
	_, err = m.Insert(ctx, &jobsflow.TJobsFlow{
		Id:         taskId,
		Type:       req.Type,
		BizCode:    req.BizCode,
		BizId:      req.BizId,
		ExecPath:   req.ExecPath,
		Param:      req.Param,
		Timeout:    req.Timeout,
		StartTime:  sql.NullTime{},
		FinishTime: sql.NullTime{},
		ResultMsg:  "{}",
		ExtInfo:    req.ExtInfo,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	logx.Infof("add realtime single task success, task is %s", taskId)
	return
}

func queryDataFromDB(ctx context.Context, req *common.QueryRealTimeSingleTaskReq) (results []*jobsflow.TJobsFlow, err error) {
	m := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	if req.Id != "" {
		var result = &jobsflow.TJobsFlow{}
		result, err = m.FindOne(ctx, req.Id)
		if err != nil {
			logx.Error(err)
			return
		}
		results = append(results, result)
		return
	}
	if req.Status < 0 || req.Status > int64(core.Finished) {
		err = errors.New("status is error")
		logx.Error(err)
		return
	}
	if req.TimeHorizon == 0 {
		req.TimeHorizon = DefaultTimeHorizon
	}
	if req.Limit == 0 {
		req.Limit = DefaultLimit
	}
	logx.Debugf("req:%+v", req)
	results, err = m.GetJobsFlowByStatus2(ctx, req.Status, req.TimeHorizon, req.Limit)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}
