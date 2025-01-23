package gracefulstop

import (
	"dmc-task/core/cron"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/server"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown() {
	cron.Stop()
	logx.Info("cron stop success!")
	if server.SvrCtx.IsMasterSource != "" {
		timewheel.Stop()
		logx.Info("timewheel stop success!")
		model.Unlock()
		logx.Info("unlock success!")
	}
}

func GracefulShutdown() {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		Shutdown()
	}()
}
