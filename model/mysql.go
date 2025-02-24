package model

import (
	"context"
	"dmc-task/model/lock"
	"dmc-task/server"
	"dmc-task/utils"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func InitMysql() {
	logx.Debug("init mysql... ", server.SvrCtx.Config.MySQL)
	config := server.SvrCtx.Config.MySQL
	conn, err := NewMysql(config.Username, config.Password, config.Host, config.Database, config.Port)
	if err != nil {
		panic(err)
	}
	server.SvrCtx.MysqlConn = conn

	sqlx.DisableStmtLog() // 禁用sql语句日志
	logx.Infof("[mysql] init success! config:%+v", config)
}

func NewMysql(username, password, host, database string, port int) (*sqlx.SqlConn, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		username, password, host, port, database)
	logx.Debugf("mysql connecting! dataSource:%s", dataSource)
	mysql := sqlx.NewMysql(dataSource)
	db, err := mysql.RawDB()
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = db.PingContext(cxt)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return &mysql, nil
}

func Lock() bool {
	var source string
	if server.SvrCtx.Config.ApiServer.Enabled {
		source = fmt.Sprintf("%s-%s:%d", utils.GetLocalIP(), server.SvrCtx.Config.ApiServer.Host, server.SvrCtx.Config.ApiServer.Port)
	} else {
		source = fmt.Sprintf("%s-%s:%d", utils.GetLocalIP(), server.SvrCtx.Config.GrpcServer.Host, server.SvrCtx.Config.GrpcServer.Port)
	}

	if ok := lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Lock(source); !ok {
		return false
	}
	server.SvrCtx.IsMasterSource = source
	return true
}

func Unlock() {
	if server.SvrCtx.IsMasterSource == "" {
		return
	}
	lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Unlock(server.SvrCtx.IsMasterSource)
}

func Renew() error {
	if server.SvrCtx.IsMasterSource == "" {
		err := fmt.Errorf("not found master source")
		return err
	}
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Renew(server.SvrCtx.IsMasterSource)
}

func Reset() error {
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Reset(server.SvrCtx.IsMasterSource)
}
