package croncycletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/croncycletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCronCycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelCronCycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCronCycleLogic {
	return &DelCronCycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelCronCycleLogic) DelCronCycle(req *types.DelCronCycleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.DelCronCycleTaskReq{}
	r.Id = req.Id
	res := croncycletask.DelCronCycle(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
