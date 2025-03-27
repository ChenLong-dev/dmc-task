package grpc

import (
	protoc "dmc-task/cmd/dmctask/grpc/task"
)

type TaskServer struct {
	protoc.UnimplementedTaskServer
}
