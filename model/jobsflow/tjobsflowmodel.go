package jobsflow

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TJobsFlowModel = (*customTJobsFlowModel)(nil)

type (
	// TJobsFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTJobsFlowModel.
	TJobsFlowModel interface {
		tJobsFlowModel
		withSession(session sqlx.Session) TJobsFlowModel

		// TransactCtx 添加事务的实现
		TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		// ExecSql 执行sql语句
		ExecSql(ctx context.Context, sql string) (sql.Result, error)
		// GetJobFlowByStatus 通过状态获取Job流水记录表的记录
		GetJobFlowByStatus(ctx context.Context, status int64) (*TJobsFlow, error)
		// GetJobsFlowByStatus 通过状态获取Job流水记录表的记录（状态）
		GetJobsFlowByStatus(ctx context.Context, status int64) ([]*TJobsFlow, error)
		// GetJobsFlowByStatus2 通过状态获取Job流水记录表的记录（状态、时间范围（小时）、条数限制）
		GetJobsFlowByStatus2(ctx context.Context, status int64, timeHorizon, limit int64) ([]*TJobsFlow, error)
		// GetJobsFlowByCronTaskId 根据定时任务id查询job任务
		GetJobsFlowByCronTaskId(ctx context.Context, cronTaskId string) (*TJobsFlow, error)
		// GetJobsFlowRowsExpectAutoSet 获取tJobsFlowRowsExpectAutoSet
		GetJobsFlowRowsExpectAutoSet() string
		// GetJobsFlowRowsWithPlaceHolder 获取JobsFlowRowsWithPlaceHolder
		GetJobsFlowRowsWithPlaceHolder() string
		// GetTableName 获取表名
		GetTableName() string
	}

	customTJobsFlowModel struct {
		*defaultTJobsFlowModel
	}
)

// NewTJobsFlowModel returns a model for the database table.
func NewTJobsFlowModel(conn sqlx.SqlConn) TJobsFlowModel {
	return &customTJobsFlowModel{
		defaultTJobsFlowModel: newTJobsFlowModel(conn),
	}
}

func (m *customTJobsFlowModel) withSession(session sqlx.Session) TJobsFlowModel {
	return NewTJobsFlowModel(sqlx.NewSqlConnFromSession(session))
}

// TransactCtx 添加事务的实现
func (m *customTJobsFlowModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// ExecSql 执行sql语句
func (m *customTJobsFlowModel) ExecSql(ctx context.Context, sql string) (sql.Result, error) {
	return m.conn.ExecCtx(ctx, sql)
}

// GetJobFlowByStatus 通过状态获取Job流水记录表的记录
func (m *customTJobsFlowModel) GetJobFlowByStatus(ctx context.Context, status int64) (*TJobsFlow, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? order by `create_time` limit 1", tJobsFlowRows, m.table)
	var resp TJobsFlow
	err := m.conn.QueryRowCtx(ctx, &resp, query, status)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetJobsFlowByStatus 通过状态获取Job流水记录表的记录（状态）
func (m *customTJobsFlowModel) GetJobsFlowByStatus(ctx context.Context, status int64) ([]*TJobsFlow, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? order by `create_time` limit 50", tJobsFlowRows, m.table)
	var resp []*TJobsFlow
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

// GetJobsFlowByStatus2 通过状态获取Job流水记录表的记录（状态、时间范围（小时）、条数限制）
func (m *customTJobsFlowModel) GetJobsFlowByStatus2(ctx context.Context, status int64, timeHorizon, limit int64) ([]*TJobsFlow, error) {
	query := fmt.Sprintf("select %s from %s where `status` >= ? and update_time >= NOW() - INTERVAL %d HOUR order by `update_time` desc limit %d",
		tJobsFlowRows, m.table, timeHorizon, limit)
	var resp []*TJobsFlow
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

// GetJobsFlowByCronTaskId 根据定时任务id查询job任务
func (m *customTJobsFlowModel) GetJobsFlowByCronTaskId(ctx context.Context, cronTaskId string) (*TJobsFlow, error) {
	query := fmt.Sprintf("select %s from %s where `cron_task_id` = ? limit 1", tJobsFlowRows, m.table)
	var resp TJobsFlow
	err := m.conn.QueryRowCtx(ctx, &resp, query, cronTaskId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetJobsFlowRowsExpectAutoSet 获取tJobsFlowRowsExpectAutoSet
func (m *customTJobsFlowModel) GetJobsFlowRowsExpectAutoSet() string {
	return tJobsFlowRowsExpectAutoSet
}

// GetJobsFlowRowsWithPlaceHolder 获取JobsFlowRowsWithPlaceHolder
func (m *customTJobsFlowModel) GetJobsFlowRowsWithPlaceHolder() string {
	return tJobsFlowRowsWithPlaceHolder
}

// GetTableName 获取表名
func (m *customTJobsFlowModel) GetTableName() string {
	return m.tableName()
}
