package grpc

import (
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"dmc-task/core"
	"fmt"
	"net"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GgpcServer *grpc.Server

func Run(conf *core.Config) {
	logx.Debug("Grpc server starting ...")
	if !conf.GrpcServer.Enabled {
		logx.Infof("Grpc Server is not enabled, exit. conf: %+v", conf.GrpcServer)
		return
	}

	addr := fmt.Sprintf("%s:%d", conf.GrpcServer.Host, conf.GrpcServer.Port)

	//创建listen监听端口
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logx.Error(err)
		return
	}

	//创建 gRPC Server 对象
	GgpcServer = grpc.NewServer()

	// 在gRPC服务上启用反射服务
	//启用反射服务后，客户端可以使用 gRPC 反射 API 查询服务器支持的服务列表、服务下的方法列表等信息。
	//这对于开发和测试阶段非常有用，因为它允许客户端在没有预先定义 .proto 文件的情况下与服务器通信。
	reflection.Register(GgpcServer)
	//处理注册到grpc服务中
	protoc.RegisterTaskServer(GgpcServer, &TaskServer{})
	logx.Infof("Starting grpc server at %s...", addr)
	// 运行 grpc server
	if err = GgpcServer.Serve(listener); err != nil {
		logx.Error(err)
		return
	}
}
