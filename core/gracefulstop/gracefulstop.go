package gracefulstop

import (
	internal "dmc-task/cmd/dmctask/api"
	"dmc-task/cmd/dmctask/grpc"
	"dmc-task/core"
	"dmc-task/core/cron"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/server"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"syscall"
)

// Shutdown 优雅关闭
func Shutdown() {
	// 优雅关闭定时任务
	cron.Stop()
	logx.Info("cron stop success!")

	// 根据master标识判断，优雅关闭时间轮以及释放分布式锁
	if server.SvrCtx.IsMasterSource != "" {
		timewheel.Stop()
		logx.Info("timewheel is gracefully stopped!")
		model.Unlock()
		logx.Info("unlock success!")
	}

	// 根据api服务配置判断，优雅关闭Api服务
	if core.Cfg.ApiServer.Enabled && internal.ApiServer != nil {
		internal.ApiServer.Stop()
		logx.Infof("api server is gracefully stopped! api:%+v", core.Cfg.ApiServer)
	}

	// 根据grpc服务配置判断，优雅关闭Grpc服务
	if core.Cfg.GrpcServer.Enabled && grpc.GgpcServer != nil {
		grpc.GgpcServer.Stop()
		logx.Infof("grpc server is gracefully stopped! grpc:%+v", core.Cfg.GrpcServer)
	}
}

func GracefulShutdown(s chan<- bool) {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		Shutdown()
		s <- true
	}()
}
