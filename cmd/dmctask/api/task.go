package internal

import (
	"dmc-task/core"
	"github.com/zeromicro/go-zero/core/logx"

	"dmc-task/cmd/dmctask/api/internal/config"
	"dmc-task/cmd/dmctask/api/internal/handler"
	"dmc-task/cmd/dmctask/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func Run(conf *core.Config) {
	var c config.Config
	c.RestConf.Host = conf.Server.Host
	c.RestConf.Port = conf.Server.Port
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}
