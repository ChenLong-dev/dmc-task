// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package jobsflow

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
	tJobsFlowFieldNames          = builder.RawFieldNames(&TJobsFlow{})
	tJobsFlowRows                = strings.Join(tJobsFlowFieldNames, ",")
	tJobsFlowRowsExpectAutoSet   = strings.Join(stringx.Remove(tJobsFlowFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tJobsFlowRowsWithPlaceHolder = strings.Join(stringx.Remove(tJobsFlowFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tJobsFlowModel interface {
		Insert(ctx context.Context, data *TJobsFlow) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*TJobsFlow, error)
		Update(ctx context.Context, data *TJobsFlow) error
		Delete(ctx context.Context, id string) error
	}

	defaultTJobsFlowModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TJobsFlow struct {
		Id           string       `db:"id"`            // job的任务ID
		Type         int64        `db:"type"`          // 任务类型
		CronTaskId   string       `db:"cron_task_id"`  // 定时任务ID
		BizCode      string       `db:"biz_code"`      // 业务Code
		BizId        string       `db:"biz_id"`        // 业务ID
		ExecPath     string       `db:"exec_path"`     // 执行路径
		Param        string       `db:"param"`         // 任务的执行参数
		Timeout      int64        `db:"timeout"`       // 任务超时时间，单位秒
		StartTime    sql.NullTime `db:"start_time"`    // 定时任务执行的实际开始时间
		FinishTime   sql.NullTime `db:"finish_time"`   // 定时任务执行的实际结束时间
		ExecInterval int64        `db:"exec_interval"` // job的执行时间（finish_time-start_time）
		Status       int64        `db:"status"`        // 执行状态
		ResultMsg    string       `db:"result_msg"`    // Job执行结果的描述
		ExtInfo      string       `db:"ext_info"`      // 扩展信息
		UpdateTime   time.Time    `db:"update_time"`   // job flow的更新时间
		CreateTime   time.Time    `db:"create_time"`   // job flow的创建时间
	}
)

func newTJobsFlowModel(conn sqlx.SqlConn) *defaultTJobsFlowModel {
	return &defaultTJobsFlowModel{
		conn:  conn,
		table: "`t_jobs_flow`",
	}
}

func (m *defaultTJobsFlowModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTJobsFlowModel) FindOne(ctx context.Context, id string) (*TJobsFlow, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tJobsFlowRows, m.table)
	var resp TJobsFlow
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

func (m *defaultTJobsFlowModel) Insert(ctx context.Context, data *TJobsFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tJobsFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Type, data.CronTaskId, data.BizCode, data.BizId, data.ExecPath, data.Param, data.Timeout, data.StartTime, data.FinishTime, data.ExecInterval, data.Status, data.ResultMsg, data.ExtInfo)
	return ret, err
}

func (m *defaultTJobsFlowModel) Update(ctx context.Context, data *TJobsFlow) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tJobsFlowRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Type, data.CronTaskId, data.BizCode, data.BizId, data.ExecPath, data.Param, data.Timeout, data.StartTime, data.FinishTime, data.ExecInterval, data.Status, data.ResultMsg, data.ExtInfo, data.Id)
	return err
}

func (m *defaultTJobsFlowModel) tableName() string {
	return m.table
}
