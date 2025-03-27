package lock

import (
	"context"
	"database/sql"
	"dmc-task/utils"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ TDistributedLocksModel = (*customTDistributedLocksModel)(nil)

type (
	// TDistributedLocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTDistributedLocksModel.
	TDistributedLocksModel interface {
		tDistributedLocksModel
		withSession(session sqlx.Session) TDistributedLocksModel

		// Lock 获得分布式锁
		Lock(source string) bool
		// Unlock 释放分布式锁
		Unlock(source string)
		// Renew 续约分布式锁
		Renew(source string) error
		// Reset 重置分布式锁
		Reset(source string) error
	}

	customTDistributedLocksModel struct {
		*defaultTDistributedLocksModel
	}
)

// NewTDistributedLocksModel returns a model for the database table.
func NewTDistributedLocksModel(conn sqlx.SqlConn) TDistributedLocksModel {
	return &customTDistributedLocksModel{
		defaultTDistributedLocksModel: newTDistributedLocksModel(conn),
	}
}

func (m *customTDistributedLocksModel) withSession(session sqlx.Session) TDistributedLocksModel {
	return NewTDistributedLocksModel(sqlx.NewSqlConnFromSession(session))
}

// ExecSql 执行sql语句
func (m *customTDistributedLocksModel) ExecSql(ctx context.Context, sql string) (sql.Result, error) {
	return m.conn.ExecCtx(ctx, sql)
}

// Lock 获得分布式锁
func (m *customTDistributedLocksModel) Lock(source string) bool {
	data := &TDistributedLocks{
		LockName:   LockName,
		Source:     source,
		LockValue:  utils.GetUTCTime().Format(time.DateTime),
		ExpireTime: utils.GetUTCTime2(time.Second * time.Duration(LockExpire)),
	}
	if _, err := m.Insert(context.Background(), data); err != nil {
		return false
	}
	logx.Debugf("[Lock] +++ unlock lock success, source:%s", source)
	return true
}

// Unlock 释放分布式锁
func (m *customTDistributedLocksModel) Unlock(source string) {
	_ = m.Delete(context.Background(), LockName)
	logx.Debugf("[Unlock] --- unlock lock success, source:%s", source)
	return
}

// Renew 续约分布式锁
func (m *customTDistributedLocksModel) Renew(source string) error {
	updateSql := fmt.Sprintf("UPDATE %s SET expire_time = DATE_ADD(NOW(), INTERVAL %d SECOND) WHERE lock_name = ? AND source = ?",
		m.tableName(), LockExpire)
	result, err := m.conn.ExecCtx(context.Background(), updateSql, LockName, source)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		err = fmt.Errorf("renew lock failed, source:%s", source)
		return err
	}
	logx.Debug("[Renew] renew lock success, source:", source)
	return nil
}

// Reset 重置分布式锁
func (m *customTDistributedLocksModel) Reset(source string) error {
	resp, err := m.FindOne(context.Background(), LockName)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		logx.Error(err)
		return err
	}
	curr := utils.GetUTCTime()
	interval := resp.ExpireTime.Sub(curr).Seconds()
	logx.Debug("[Reset] interval:", interval)
	if interval <= 0 {
		err = m.Delete(context.Background(), LockName)
		if err != nil {
			logx.Error(err)
			return err
		}
		logx.Debugf("[Reset] reset lock success, source:%s", resp.Source)
	}
	return nil
}
