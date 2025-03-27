package fixedtimesingletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/fixedtimesingletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCronLogic {
	return &AddCronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type FixedTimeSingleTask struct {
	Type       int64  `json:"type"`
	BizCode    string `json:"biz_code"`
	BizId      string `json:"biz_id"`
	ExecPath   string `json:"exec_path"`
	ExecTime   int64  `json:"exec_time"`
	Param      string `json:"param"`
	StartTime  int64  `json:"start_time"`
	FinishTime int64  `json:"finish_time"`
	Interval   int64  `json:"interval"`
	ResultMsg  string `json:"result_msg"`
	ExtInfo    string `json:"ext_info"`
}

func (l *AddCronLogic) AddCron(req *types.AddFixedTimeSingleTaskReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	r := &common.AddFixedTimeSingleTaskReq{}
	r.Type = req.Type
	r.BizCode = req.BizCode
	r.BizId = req.BizId
	r.ExecPath = req.ExecPath
	r.ExecTime = req.ExecTime
	r.Param = req.Param
	r.Timeout = req.Timeout
	r.ExtInfo = req.ExtInfo
	res := fixedtimesingletask.AddCron(l.ctx, r)
	resp = &types.Response{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	return
}
