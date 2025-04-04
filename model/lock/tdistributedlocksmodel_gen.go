// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package lock

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tDistributedLocksFieldNames          = builder.RawFieldNames(&TDistributedLocks{})
	tDistributedLocksRows                = strings.Join(tDistributedLocksFieldNames, ",")
	tDistributedLocksRowsExpectAutoSet   = strings.Join(stringx.Remove(tDistributedLocksFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tDistributedLocksRowsWithPlaceHolder = strings.Join(stringx.Remove(tDistributedLocksFieldNames, "`lock_name`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tDistributedLocksModel interface {
		Insert(ctx context.Context, data *TDistributedLocks) (sql.Result, error)
		FindOne(ctx context.Context, lockName string) (*TDistributedLocks, error)
		Update(ctx context.Context, data *TDistributedLocks) error
		Delete(ctx context.Context, lockName string) error
	}

	defaultTDistributedLocksModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TDistributedLocks struct {
		LockName   string    `db:"lock_name"`
		Source     string    `db:"source"`
		LockValue  string    `db:"lock_value"`
		ExpireTime time.Time `db:"expire_time"`
	}
)

func newTDistributedLocksModel(conn sqlx.SqlConn) *defaultTDistributedLocksModel {
	return &defaultTDistributedLocksModel{
		conn:  conn,
		table: "`t_distributed_locks`",
	}
}

func (m *defaultTDistributedLocksModel) Delete(ctx context.Context, lockName string) error {
	query := fmt.Sprintf("delete from %s where `lock_name` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, lockName)
	return err
}

func (m *defaultTDistributedLocksModel) FindOne(ctx context.Context, lockName string) (*TDistributedLocks, error) {
	query := fmt.Sprintf("select %s from %s where `lock_name` = ? limit 1", tDistributedLocksRows, m.table)
	var resp TDistributedLocks
	err := m.conn.QueryRowCtx(ctx, &resp, query, lockName)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTDistributedLocksModel) Insert(ctx context.Context, data *TDistributedLocks) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, tDistributedLocksRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.LockName, data.Source, data.LockValue, data.ExpireTime)
	return ret, err
}

func (m *defaultTDistributedLocksModel) Update(ctx context.Context, data *TDistributedLocks) error {
	query := fmt.Sprintf("update %s set %s where `lock_name` = ?", m.table, tDistributedLocksRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Source, data.LockValue, data.ExpireTime, data.LockName)
	return err
}

func (m *defaultTDistributedLocksModel) tableName() string {
	return m.table
}
