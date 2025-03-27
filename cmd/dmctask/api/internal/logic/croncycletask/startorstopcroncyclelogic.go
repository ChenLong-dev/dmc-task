package croncycletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/croncycletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartorstopCronCycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartorstopCronCycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartorstopCronCycleLogic {
	return &StartorstopCronCycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartorstopCronCycleLogic) StartorstopCronCycle(req *types.StartOrStopCronCycleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.StartOrStopCronCycleTaskReq{}
	r.Id = req.Id
	r.IsStart = req.IsStart
	res := croncycletask.StartOrStopCronCycle(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
