package croncycletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/croncycletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryCronCycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryCronCycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCronCycleLogic {
	return &QueryCronCycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCronCycleLogic) QueryCronCycle(req *types.QueryCronCycleTaskReq) (resp *types.QueryTaskConfigResp, err error) {
	// todo: add your logic here and delete this line
	r := &common.QueryCronCycleTaskReq{}
	r.Id = req.Id
	res := croncycletask.QueryCronCycle(l.ctx, r)
	resp = &types.QueryTaskConfigResp{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	for _, v := range res.Data {
		resp.Data = append(resp.Data, types.CronCycleTaskData{
			BaseData: types.BaseData{
				Id:     v.Id,
				Status: v.Status,
			},
			CronCycleTask: types.CronCycleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				Cron:     v.Cron,
				ExecPath: v.ExecPath,
				Param:    v.Param,
				ExtInfo:  v.ExtInfo,
			},
		})
	}
	return
}
