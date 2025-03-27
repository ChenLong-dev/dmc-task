package fixedtimesingletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.BizCode), logx.Field("biz_id", req.BizId))

	// 1、格式校验
	if req.Type != int64(core.FixedTimeSingleTask) {
		err = errors.New("type is not fixed time task")
		logx.WithContext(ctx).Error(err)
		return
	}

	execTime := utils.GetTime(req.ExecTime)
	now := utils.GetUTCTime().Add(time.Second * 60)
	internal := execTime.Sub(now)
	if internal < 0 {
		err = fmt.Errorf("exec time must be later 60s than current time, current time is %s, exec time is %s, internal:%ds",
			now.Format(time.DateTime), execTime.Format(time.DateTime), internal)
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、添加任务（入库）
	taskId, err = cron.AddDataToCronTasks(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
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
	ctx = logx.ContextWithFields(ctx, logx.Field("id", req.Id))

	// 1、格式校验
	if req.Id == "" {
		err = fmt.Errorf("task id is empty")
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、删除任务
	err = cron.DelDataFromCronTasks(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
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
	ctx = logx.ContextWithFields(ctx, logx.Field("filter", req.Filter))

	// 1、查询数据库
	total, results, err := cron.QueryDataFromCronTasks(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, common.FixedTimeSingleTaskData{
			BaseData: common.BaseData{
				Id:     v.Id,
				Status: v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
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
	resp.Page = req.Page
	resp.Page.Total = total
	return
}
