package model

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/common"
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
	if !server.SvrCtx.Config.App.IsDistributed { // 判断是否开启分布式锁
		server.SvrCtx.IsMasterSource = "master"
		return true
	}
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
	if !server.SvrCtx.Config.App.IsDistributed { // 判断是否开启分布式锁
		server.SvrCtx.IsMasterSource = ""
		return
	}
	if server.SvrCtx.IsMasterSource == "" {
		return
	}
	lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Unlock(server.SvrCtx.IsMasterSource)
}

func Renew() error {
	if !server.SvrCtx.Config.App.IsDistributed { // 判断是否开启分布式锁
		return nil
	}
	if server.SvrCtx.IsMasterSource == "" {
		err := fmt.Errorf("not found master source")
		return err
	}
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Renew(server.SvrCtx.IsMasterSource)
}

func Reset() error {
	if !server.SvrCtx.Config.App.IsDistributed { // 判断是否开启分布式锁
		return nil
	}
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Reset(server.SvrCtx.IsMasterSource)
}

// ***************************************************  分页获取

type PaginationRequest struct {
	Ctx      context.Context
	Conn     sqlx.SqlConn
	Table    string
	Where    string
	OrderBy  string
	Args     []interface{}
	Page     int
	PageSize int
}

type PaginationResult[T any] struct {
	Count int  `json:"count"`
	Data  []*T `json:"data"`
}

func Paginate[T any](req *PaginationRequest) (result *PaginationResult[T], err error) {
	if req == nil {
		err = fmt.Errorf("req is nil")
		logx.WithContext(req.Ctx).Error(err)
		return
	}

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", req.Table)
	if req.Where != "" {
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, req.Where)
	}

	var count int
	if err = req.Conn.QueryRowCtx(req.Ctx, &count, countQuery, req.Args...); err != nil {
		logx.WithContext(req.Ctx).Error(err)
		return
	}

	// 查询数据
	dataQuery := fmt.Sprintf("SELECT * FROM %s", req.Table)
	if req.Where != "" {
		dataQuery = fmt.Sprintf("%s WHERE %s", dataQuery, req.Where)
	}
	if req.OrderBy != "" {
		dataQuery = fmt.Sprintf("%s ORDER BY %s", dataQuery, req.OrderBy)
	}
	if req.Page > 0 && req.PageSize > 0 {
		dataQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", dataQuery, req.PageSize, offset)
	}
	logx.WithContext(req.Ctx).Debug(dataQuery)

	var data []*T
	if err = req.Conn.QueryRowsCtx(req.Ctx, &data, dataQuery, req.Args...); err != nil {
		logx.WithContext(req.Ctx).Error(err)
		return
	}
	return &PaginationResult[T]{
		Count: count,
		Data:  data,
	}, nil
}

func Query[T any](ctx context.Context, table string, filter common.FilterBase, page common.PageBase) (result *PaginationResult[T], err error) {
	// 1、查询条件
	query := PaginationRequest{}
	query.Ctx = ctx
	query.Conn = *server.SvrCtx.MysqlConn
	query.Table = table

	if filter.TimeType != "" {
		query.OrderBy = fmt.Sprintf("%s DESC", filter.TimeType)
	}

	query.Page = int(page.Page)
	query.PageSize = int(page.PageSize)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	where := ""
	if filter.Id != "" {
		where = fmt.Sprintf("id = '%s'", filter.Id)
	}

	if filter.BizCode != "" {
		if where != "" {
			where = fmt.Sprintf("%s AND biz_code = '%s'", where, filter.BizCode)
		} else {
			where = fmt.Sprintf("biz_code = '%s'", filter.BizCode)
		}
	}

	if filter.BizId != "" {
		if where != "" {
			where = fmt.Sprintf("%s AND biz_id = '%s'", where, filter.BizId)
		} else {
			where = fmt.Sprintf("biz_id = '%s'", filter.BizId)
		}
	}

	if filter.CronTaskId != "" {
		if where != "" {
			where = fmt.Sprintf("%s AND cron_task_id = '%s'", where, filter.CronTaskId)
		} else {
			where = fmt.Sprintf("cron_task_id = '%s'", filter.CronTaskId)
		}
	}

	if filter.Status > int64(core.Stoped) && filter.Status < int64(core.Finished) {
		if where != "" {
			where = fmt.Sprintf("%s AND status = %d", where, filter.Status)
		} else {
			where = fmt.Sprintf("status = %d", filter.Status)
		}
	}

	if filter.TimeType != "" && filter.Start != "" && filter.End != "" {
		if where != "" {
			where = fmt.Sprintf("%s AND %s BETWEEN '%s' AND '%s'", where, filter.TimeType, filter.Start, filter.End)
		} else {
			where = fmt.Sprintf("%s BETWEEN '%s' AND '%s'", filter.TimeType, filter.Start, filter.End)
		}
	}
	query.Where = where

	// 2、查询分页
	res, err := Paginate[T](&query)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	return res, nil
}
