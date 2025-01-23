package croncycletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/croncycletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModCronCycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModCronCycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModCronCycleLogic {
	return &ModCronCycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModCronCycleLogic) ModCronCycle(req *types.ModCronCycleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.ModCronCycleTaskReq{}
	r.Id = req.Id
	r.Type = req.Type
	r.Cron = req.Cron
	r.BizCode = req.BizCode
	r.ExecPath = req.ExecPath
	r.Param = req.Param
	r.ExtInfo = req.ExtInfo
	res := croncycletask.ModCronCycle(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
