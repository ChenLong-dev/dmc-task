package realtimesingletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/realtimesingletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddJobLogic {
	return &AddJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddJobLogic) AddJob(req *types.AddRealTimeSingleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.AddRealTimeSingleTaskReq{}
	r.Type = req.Type
	r.BizCode = req.BizCode
	r.BizId = req.BizId
	r.ExecPath = req.ExecPath
	r.Param = req.Param
	r.Timeout = req.Timeout
	r.ExtInfo = req.ExtInfo
	res := realtimesingletask.AddJob(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
