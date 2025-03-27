package crontasks

import (
	"context"
	"database/sql"
	"dmc-task/utils"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ TCronTasksModel = (*customTCronTasksModel)(nil)

type (
	// TCronTasksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTCronTasksModel.
	TCronTasksModel interface {
		tCronTasksModel
		withSession(session sqlx.Session) TCronTasksModel

		// TransactCtx 添加事务的实现
		TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		// ExecSql 执行sql语句
		ExecSql(ctx context.Context, sql string) (sql.Result, error)
		// GetCronTaskByStatus 通过状态获取待执行的定时任务
		GetCronTaskByStatus(ctx context.Context, status int64) (*TCronTasks, error)
		// GetCronTasksByStatus 通过状态获取待执行的定时任务（状态、执行时间范围）
		GetCronTasksByStatus(ctx context.Context, status int64, execTime time.Duration) ([]*TCronTasks, error)
		// GetCronTasksByStatus2 通过状态获取定时任务（状态、时间范围（小时）、条数限制）
		GetCronTasksByStatus2(ctx context.Context, status int64, timeHorizon, limit int64) ([]*TCronTasks, error)
		// GetCronTasksRowsWithPlaceHolder 获取CronTasksRowsWithPlaceHolder
		GetCronTasksRowsWithPlaceHolder() string
		// GetTableName 获取表名
		GetTableName() string
	}

	customTCronTasksModel struct {
		*defaultTCronTasksModel
	}
)

// NewTCronTasksModel returns a model for the database table.
func NewTCronTasksModel(conn sqlx.SqlConn) TCronTasksModel {
	return &customTCronTasksModel{
		defaultTCronTasksModel: newTCronTasksModel(conn),
	}
}

func (m *customTCronTasksModel) withSession(session sqlx.Session) TCronTasksModel {
	return NewTCronTasksModel(sqlx.NewSqlConnFromSession(session))
}

// TransactCtx 添加事务的实现
func (m *customTCronTasksModel) TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// ExecSql 执行sql语句
func (m *customTCronTasksModel) ExecSql(ctx context.Context, sql string) (sql.Result, error) {
	return m.conn.ExecCtx(ctx, sql)
}

// GetCronTaskByStatus 通过状态获取待执行的定时任务
func (m *customTCronTasksModel) GetCronTaskByStatus(ctx context.Context, status int64) (*TCronTasks, error) {
	oneSecondsBefore := utils.GetUTCTime2(-1 * time.Second)
	oneSecondsAfter := utils.GetUTCTime2(1 * time.Second)
	query := fmt.Sprintf("select %s from %s where `status` = ? and exec_time between ? and ? order by `exec_time`", tCronTasksRows, m.table)
	var resp *TCronTasks
	err := m.conn.QueryRowCtx(ctx, &resp, query, status, oneSecondsBefore, oneSecondsAfter)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetCronTasksByStatus 通过状态获取待执行的定时任务（状态、执行时间范围）
func (m *customTCronTasksModel) GetCronTasksByStatus(ctx context.Context, status int64, execTime time.Duration) ([]*TCronTasks, error) {
	oneSecondsBefore := utils.GetUTCTime2(0 * time.Second)
	oneSecondsAfter := utils.GetUTCTime2(execTime)
	query := fmt.Sprintf("select %s from %s where `status` = ? and exec_time between ? and ? order by `exec_time`", tCronTasksRows, m.table)
	var resp []*TCronTasks
	err := m.conn.QueryRowsCtx(ctx, &resp, query, status, oneSecondsBefore, oneSecondsAfter)
	switch err {
	case nil:
		if len(resp) == 0 {
			return nil, ErrNotFound
		}
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetCronTasksByStatus2 通过状态获取定时任务（状态、时间范围（小时）、条数限制）
func (m *customTCronTasksModel) GetCronTasksByStatus2(ctx context.Context, status int64, timeHorizon, limit int64) ([]*TCronTasks, error) {
	query := fmt.Sprintf("select %s from %s where `status` >= ? and update_time >= NOW() - INTERVAL %d HOUR order by `update_time` desc limit %d",
		tCronTasksRows, m.table, timeHorizon, limit)
	var resp []*TCronTasks
	err := m.conn.QueryRowsCtx(ctx, &resp, query, status)
	switch err {
	case nil:
		if len(resp) == 0 {
			return nil, ErrNotFound
		}
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetCronTasksRowsWithPlaceHolder 获取CronTasksRowsWithPlaceHolder
func (m *customTCronTasksModel) GetCronTasksRowsWithPlaceHolder() string {
	return tCronTasksRowsWithPlaceHolder
}

// GetTableName 获取表名
func (m *customTCronTasksModel) GetTableName() string {
	return m.tableName()
}
