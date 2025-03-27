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

func (l *QueryCronCycleLogic) QueryCronCycle(req *types.QueryCronCycleTaskReq) (resp *types.QueryCronCycleTaskResp, err error) {
	// todo: add your logic here and delete this line
	r := &common.QueryCronCycleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	//r.Filter.BizId = req.Filter.BizId // 该类型任务没有定时任务ID
	//r.Filter.CronTaskId = req.Filter.CronTaskId // 该类型任务没有定时任务ID
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	res := croncycletask.QueryCronCycle(l.ctx, r)
	resp = &types.QueryCronCycleTaskResp{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	for _, v := range res.Data {
		resp.Data = append(resp.Data, types.CronCycleTaskData{
			BaseData: types.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: v.UpdateTime,
				CreateTime: v.CreateTime,
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
	resp.Page.Total = res.Page.Total
	resp.Page.Page = res.Page.Page
	resp.Page.PageSize = res.Page.PageSize
	return
}
