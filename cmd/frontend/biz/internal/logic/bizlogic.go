package logic

import (
	"context"
	"dmc-task/cmd/frontend/internal"

	"dmc-task/cmd/frontend/biz/internal/svc"
	"dmc-task/cmd/frontend/biz/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BizLogic {
	return &BizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BizLogic) Biz(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &internal.BizRequest{}
	r.BizType = req.BizType
	r.TaskId = req.TaskId
	r.Start = req.Start
	r.End = req.End
	res := internal.Biz(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	resp.Data = res.Data
	return
}
