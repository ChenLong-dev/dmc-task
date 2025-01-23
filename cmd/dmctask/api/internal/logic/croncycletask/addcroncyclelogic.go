package croncycletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/croncycletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCronCycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCronCycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCronCycleLogic {
	return &AddCronCycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type CronCycleTask struct {
	Type     string `json:"type" validate:"required"`
	BizCode  string `json:"biz_code" validate:"required"`
	Cron     string `json:"cron" validate:"required"`
	ExecPath string `json:"exec_path" validate:"required"`
	Param    string `json:"param" validate:"required"`
	ExtInfo  string `json:"ext_info,optional"`
}

func (l *AddCronCycleLogic) AddCronCycle(req *types.AddCronCycleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.AddCronCycleTaskReq{}
	r.Type = req.Type
	r.Cron = req.Cron
	r.BizCode = req.BizCode
	r.ExecPath = req.ExecPath
	r.Param = req.Param
	r.Timeout = req.Timeout
	r.ExtInfo = req.ExtInfo
	res := croncycletask.AddCronCycle(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
