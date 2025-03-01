package fixedtimesingletask

import (
	"context"
	"dmc-task/cmd/dmctask/internal/fixedtimesingletask"
	"dmc-task/core/common"

	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryCronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCronLogic {
	return &QueryCronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCronLogic) QueryCron(req *types.QueryFixedTimeSingleTaskReq) (resp *types.QueryFixedTimeSingleTaskResp, err error) {
	// todo: add your logic here and delete this line
	r := &common.QueryFixedTimeSingleTaskReq{}
	r.Filter.Id = req.Filter.Id
	r.Filter.BizCode = req.Filter.BizCode
	r.Filter.BizId = req.Filter.BizId
	//r.Filter.CronTaskId = req.Filter.CronTaskId // 该类型任务没有定时任务ID
	r.Filter.Status = req.Filter.Status
	r.Filter.TimeType = req.Filter.TimeType
	r.Filter.Start = req.Filter.Start
	r.Filter.End = req.Filter.End
	r.Page.Page = req.Page.Page
	r.Page.PageSize = req.Page.PageSize
	res := fixedtimesingletask.QueryCron(l.ctx, r)
	resp = &types.QueryFixedTimeSingleTaskResp{}
	resp.Code = res.Code
	resp.Msg = res.Msg
	for _, v := range res.Data {
		resp.Data = append(resp.Data, types.FixedTimeSingleTaskData{
			BaseData: types.BaseData{
				Id:         v.Id,
				Status:     v.Status,
				UpdateTime: v.UpdateTime,
				CreateTime: v.CreateTime,
			},
			FixedTimeSingleTask: types.FixedTimeSingleTask{
				Type:     v.Type,
				BizCode:  v.BizCode,
				BizId:    v.BizId,
				ExecPath: v.ExecPath,
				ExecTime: v.ExecTime,
				Param:    v.Param,
				ExtInfo:  v.ExtInfo,
			},
			StartTime:  v.StartTime,
			FinishTime: v.FinishTime,
			Interval:   v.Interval,
			ResultMsg:  v.ResultMsg,
		})
	}
	resp.Page.Total = res.Page.Total
	resp.Page.Page = res.Page.Page
	resp.Page.PageSize = res.Page.PageSize
	return
}
