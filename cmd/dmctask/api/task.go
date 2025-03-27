package internal

import (
	"dmc-task/core"
	"dmc-task/core/middleware"

	"github.com/zeromicro/go-zero/core/logx"

	"dmc-task/cmd/dmctask/api/internal/config"
	"dmc-task/cmd/dmctask/api/internal/handler"
	"dmc-task/cmd/dmctask/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var ApiServer *rest.Server

func Run(conf *core.Config) {
	logx.Debug("Api server starting ...")
	if !conf.ApiServer.Enabled {
		logx.Infof("Api Server is not enabled, exit. conf: %+v", conf.ApiServer)
		return
	}
	var c config.Config
	c.RestConf.Host = conf.ApiServer.Host
	c.RestConf.Port = conf.ApiServer.Port
	ApiServer = rest.MustNewServer(c.RestConf)
	ApiServer.Use(middleware.AuthWithMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(ApiServer, ctx)

	logx.Infof("Starting api server at %s:%d...", c.Host, c.Port)
	ApiServer.Start()
}
