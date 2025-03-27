package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

//////////////////////////////////////////////////////////////

func NewRedis(addr, password string, database int) (*redis.Client, error) {
	logx.Debugf("redis connecting! addr:%s, password:%s, database:%d", addr, password, database)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := rdb.Ping(cxt)
	if cmd.Err() != nil {
		logx.Error(cmd.Err())
		return nil, cmd.Err()
	}
	return rdb, nil
}
