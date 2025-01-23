package fixedtimesingletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/model/crontasks"
	"dmc-task/server"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const (
	DefaultLimit       = 500
	DefaultTimeHorizon = 1
)

func AddCron(ctx context.Context, req *common.AddFixedTimeSingleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	var taskId string
	defer func() {
		if err != nil {
			resp.Code = core.FixCronError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.FixCronError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = fmt.Sprintf("%s, task id is %s", core.Success.Msg, taskId)
		}
	}()
	// 1、格式校验
	if req.Type != int64(core.FixedTimeSingleTask) {
		err = errors.New("type is not fixed time task")
		logx.Error(err)
		return
	}

	execTime := utils.GetTime(req.ExecTime)
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		err = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		return
	}
	// 2、添加任务（入库）
	taskId, err = cron.AddFixedTimeTask(ctx, &req.FixedTimeSingleTask)
	if err != nil {
		logx.Error(err)
		return
	}
	// 3、返回响应
	return
}

func DelCron(ctx context.Context, req *common.DelFixedTimeSingleTaskReq) (resp *common.Response) {
	var err error
	resp = &common.Response{}
	defer func() {
		if err != nil {
			resp.Code = core.FixCronError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.FixCronError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = core.Success.Msg
		}
	}()
	// 1、格式校验
	if req.Id == "" {
		err = fmt.Errorf("task id is empty")
		return
	}
	// 2、删除任务
	err = delDataFromDB(ctx, req)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}

func QueryCron(ctx context.Context, req *common.QueryFixedTimeSingleTaskReq) (resp *common.QueryFixedTimeSingleTaskResp) {
	var err error
	resp = &common.QueryFixedTimeSingleTaskResp{}
	defer func() {
		if err != nil {
			resp.Code = core.FixCronError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.FixCronError.Msg, err.Error())
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
		resp.Data = append(resp.Data, common.FixedTimeSingleTaskData{
			BaseData: common.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			FixedTimeSingleTask: common.FixedTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				ExecTime: utils.GetTimestamp(v.ExecTime),
				Param:    v.Param,
				Timeout:  v.Timeout,
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

func delDataFromDB(ctx context.Context, req *common.DelFixedTimeSingleTaskReq) (err error) {
	m := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	// 1、根据id查找任务是否存在
	result, err := m.FindOne(ctx, req.Id)
	if err != nil {
		logx.Error(err)
		return
	}
	// 2、校验时间（执行时间要大于1m前）
	execTime := result.ExecTime
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		err = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		return
	}
	// 3、删除任务
	err = m.Delete(ctx, req.Id)
	if err != nil {
		logx.Error(err)
		return
	}
	logx.Infof("delete cron task success, task is %+v", result)
	return
}

func queryDataFromDB(ctx context.Context, req *common.QueryFixedTimeSingleTaskReq) (results []*crontasks.TCronTasks, err error) {
	m := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	if req.Id != "" {
		var result = &crontasks.TCronTasks{}
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
	results, err = m.GetCronTasksByStatus2(ctx, req.Status, req.TimeHorizon, req.Limit)
	if err != nil {
		logx.Error(err)
		return
	}
	return
}
