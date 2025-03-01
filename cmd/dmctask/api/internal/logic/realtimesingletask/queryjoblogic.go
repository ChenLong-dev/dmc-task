package realtimesingletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/realtimesingletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryJobLogic {
	return &QueryJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryJobLogic) QueryJob(req *types.QueryRealTimeSingleTaskReq) (resp *types.QueryRealTimeSingleTaskResp, err error) {
	// todo: add your logic here and delete this line
	r := &common.QueryRealTimeSingleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	r.Filter.BizId = req.Filter.BizId
	r.Filter.CronTaskId = req.Filter.CronTaskId
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	res := realtimesingletask.QueryJob(l.ctx, r)
	resp = &types.QueryRealTimeSingleTaskResp{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	for _, v := range res.Data {
		resp.Data = append(resp.Data, types.RealTimeSingleTaskData{
			BaseData: types.BaseData{
				Id:     v.Id,
				Status: v.Status,
				UpdateTime: v.UpdateTime,
				CreateTime: v.CreateTime,
			},
			RealTimeSingleTask: types.RealTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				ExtInfo:  v.ExtInfo,
			},
			StartTime:  v.StartTime,
			FinishTime: v.FinishTime,
			Interval:   v.Interval,
			ResultMsg:  v.ResultMsg,
		})
	}
	resp.Page.Total = res.Page.Total
	resp.Page.Page = res.Page.Page
	resp.Page.PageSize = res.Page.PageSize
	return
}
