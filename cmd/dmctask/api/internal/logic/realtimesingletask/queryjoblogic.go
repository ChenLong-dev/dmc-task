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
	r.Id = req.Id
	r.Status = req.Status
	r.TimeHorizon = req.TimeHorizon
	r.Limit = req.Limit
	res := realtimesingletask.QueryJob(l.ctx, r)
	resp = &types.QueryRealTimeSingleTaskResp{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	for _, v := range res.Data {
		resp.Data = append(resp.Data, types.RealTimeSingleTaskData{
			BaseData: types.BaseData{
				Id:     v.Id,
				Status: v.Status,
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
	return
}
