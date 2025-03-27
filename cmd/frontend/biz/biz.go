package biz

import (
	"dmc-task/core"
	"github.com/zeromicro/go-zero/core/logx"

	"dmc-task/cmd/frontend/biz/internal/config"
	"dmc-task/cmd/frontend/biz/internal/handler"
	"dmc-task/cmd/frontend/biz/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var ApiServer *rest.Server

func Run(conf *core.FrontendConfig) {
	logx.Debug("Api server starting ...")
	var c config.Config
	c.RestConf.Host = conf.ApiServer.Host
	c.RestConf.Port = conf.ApiServer.Port
	ApiServer = rest.MustNewServer(c.RestConf)
	defer ApiServer.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(ApiServer, ctx)

	logx.Infof("Starting api server at %s:%d...", c.Host, c.Port)
	ApiServer.Start()
}
