package croncycletasks

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TCronCycleTasksModel = (*customTCronCycleTasksModel)(nil)

type (
	// TCronCycleTasksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTCronCycleTasksModel.
	TCronCycleTasksModel interface {
		tCronCycleTasksModel
		withSession(session sqlx.Session) TCronCycleTasksModel

		// TransactCtx 添加事务的实现
		TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		// ExecSql 执行sql语句
		ExecSql(ctx context.Context, sql string) (sql.Result, error)
		// GetCronTaskByBizCodeAndType 通过类型和业务code获取定时任务
		GetCronTaskByBizCodeAndType(ctx context.Context, taskType int64, bizCode string) (*TCronCycleTasks, error)
		// GetCronTaskById 通过id从数据库中获取任务信息
		GetCronTaskById(ctx context.Context, id string) (*TCronCycleTasks, error)
		// GetCronTasks 从数据库中获取所有任务信息
		GetCronTasks(ctx context.Context) ([]*TCronCycleTasks, error)
		// GetTableName 获取表名
		GetTableName() string
	}

	customTCronCycleTasksModel struct {
		*defaultTCronCycleTasksModel
	}
)

// NewTCronCycleTasksModel returns a model for the database table.
func NewTCronCycleTasksModel(conn sqlx.SqlConn) TCronCycleTasksModel {
	return &customTCronCycleTasksModel{
		defaultTCronCycleTasksModel: newTCronCycleTasksModel(conn),
	}
}

func (m *customTCronCycleTasksModel) withSession(session sqlx.Session) TCronCycleTasksModel {
	return NewTCronCycleTasksModel(sqlx.NewSqlConnFromSession(session))
}

// TransactCtx 添加事务的实现
func (m *customTCronCycleTasksModel) TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// ExecSql 执行sql语句
func (m *customTCronCycleTasksModel) ExecSql(ctx context.Context, sql string) (sql.Result, error) {
	return m.conn.ExecCtx(ctx, sql)
}

// GetCronTaskByBizCodeAndType 通过类型和业务code获取定时任务
func (m *customTCronCycleTasksModel) GetCronTaskByBizCodeAndType(ctx context.Context, taskType int64, bizCode string) (*TCronCycleTasks, error) {
	query := fmt.Sprintf("select %s from %s where `type` = ? and biz_code = ? limit 1", tCronCycleTasksRows, m.table)
	var resp TCronCycleTasks
	err := m.conn.QueryRowCtx(ctx, &resp, query, taskType, bizCode)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetCronTaskById 通过id从数据库中获取任务信息
func (m *customTCronCycleTasksModel) GetCronTaskById(ctx context.Context, id string) (*TCronCycleTasks, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tCronCycleTasksRows, m.table)
	var resp TCronCycleTasks
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetCronTasks 从数据库中获取所有任务信息
func (m *customTCronCycleTasksModel) GetCronTasks(ctx context.Context) ([]*TCronCycleTasks, error) {
	query := fmt.Sprintf("select %s from %s order by update_time asc", tCronCycleTasksRows, m.table)
	var resp []*TCronCycleTasks
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
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

// GetTableName 获取表名
func (m *customTCronCycleTasksModel) GetTableName() string {
	return m.tableName()
}
