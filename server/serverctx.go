package server

import (
	"dmc-task/core"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var SvrCtx *ServiceContext

type ServiceContext struct {
	Config         *core.Config
	MysqlConn      *sqlx.SqlConn
	RedisClient    *redis.Client
	IsMasterSource string
}

func NewServiceContext(c *core.Config) *ServiceContext {
	SvrCtx = &ServiceContext{
		Config:         c,
		MysqlConn:      nil,
		RedisClient:    nil,
		IsMasterSource: "",
	}

	return SvrCtx
}
