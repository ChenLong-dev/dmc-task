package realtimesingletask

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/command"
	"dmc-task/core/common"
	"dmc-task/core/cron"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
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
	ctx = logx.ContextWithFields(ctx, logx.Field("biz_code", req.BizCode), logx.Field("biz_id", req.BizId))

	// 1、格式校验
	if req.Type != int64(core.RealTimeSingleTask) {
		err = errors.New("type is not real time task")
		logx.WithContext(ctx).Error(err)
		return
	}
	if req.Timeout <= 0 {
		req.Timeout = command.DefaultTimeout
	}

	// 2、入库+执行+更新
	taskId, err = cron.AddRealTimeTask(ctx, req.RealTimeSingleTask)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	// 3、返回响应
	return
}

func QueryJob(ctx context.Context, req *common.QueryRealTimeSingleTaskReq) (resp *common.QueryRealTimeSingleTaskResp) {
	logx.WithContext(ctx).Infof("[QueryJob] req:%+v", req)
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
	ctx = logx.ContextWithFields(ctx, logx.Field("filter", req.Filter))

	// 1、查询数据库
	total, results, err := cron.QueryDataFromJobsFlow(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	// 2、组装
	for _, v := range results {
		resp.Data = append(resp.Data, common.RealTimeSingleTaskData{
			BaseData: common.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: utils.GetTimeStr(v.UpdateTime),
				CreateTime: utils.GetTimeStr(v.CreateTime),
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
	resp.Page = req.Page
	resp.Page.Total = total
	return
}
