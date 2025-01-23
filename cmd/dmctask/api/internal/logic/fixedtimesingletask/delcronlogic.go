package fixedtimesingletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/fixedtimesingletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCronLogic {
	return &DelCronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelCronLogic) DelCron(req *types.DelFixedTimeSingleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.DelFixedTimeSingleTaskReq{}
	r.Id = req.Id
	res := fixedtimesingletask.DelCron(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
