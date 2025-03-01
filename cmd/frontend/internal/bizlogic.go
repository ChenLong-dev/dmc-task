package internal

import (
	"context"
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"dmc-task/core"
	"dmc-task/core/common"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

var GrpcClient *grpc.ClientConn

type BizRequest struct {
	BizType int    `json:"biz_type" validate:"required"`
	TaskId  string `json:"task_id,optional"`
	Start   int64  `json:"start,optional"`
	End     int64  `json:"end,optional"`
}

type BizResponse struct {
	common.Base
	Data interface{} `json:"data"`
}

func Biz(ctx context.Context, req *BizRequest) (resp *BizResponse) {
	var err error
	resp = &BizResponse{}
	defer func() {
		if err != nil {
			resp.Code = core.InnerError.Code
			resp.Msg = fmt.Sprintf("%s: %s", core.InnerError.Msg, err.Error())
		} else {
			resp.Code = core.Success.Code
			resp.Msg = core.Success.Msg
		}
	}()

	switch core.TaskType(req.BizType) {
	case core.RealTimeSingleTask:
		// todo: real-time single task
		resp.Data, err = QueryRealTimeSingleTaskList(ctx, req.TaskId, req.Start, req.End)
	case core.FixedTimeSingleTask:
		// todo: fixed-time single task
		resp.Data, err = QueryFixedTimeSingleTaskList(ctx, req.TaskId, req.Start, req.End)
	case core.CronCycleTask:
		// todo: cron cycle task
		resp.Data, err = QueryCronCycleTaskList(ctx, req.TaskId, req.Start, req.End)
	default:
		err = fmt.Errorf("invalid biz_type: %d", req.BizType)
		logx.Error(err)
		return
	}

	return
}

func QueryRealTimeSingleTaskList(ctx context.Context, id string, start, end int64) (data interface{}, err error) {
	client := protoc.NewTaskClient(GrpcClient)
	resp, err := client.QueryRealTimeSingleTask(ctx, &protoc.QueryRealTimeSingleTaskReq{
		Id:          id,
		TimeHorizon: 300,
		Limit:       100,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	return resp.Data, nil
}

func QueryFixedTimeSingleTaskList(ctx context.Context, id string, start, end int64) (data interface{}, err error) {
	client := protoc.NewTaskClient(GrpcClient)
	resp, err := client.QueryFixedTimeSingleTask(ctx, &protoc.QueryFixedTimeSingleTaskReq{
		Id:          id,
		TimeHorizon: 300,
		Limit:       100,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	return resp.Data, nil
}

func QueryCronCycleTaskList(ctx context.Context, id string, start, end int64) (data interface{}, err error) {
	client := protoc.NewTaskClient(GrpcClient)
	resp, err := client.QueryCronCycleTask(ctx, &protoc.QueryCronCycleTaskReq{
		Id: id,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	return resp.Data, nil
}
