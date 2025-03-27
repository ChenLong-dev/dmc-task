package model

import (
	"context"
	"dmc-task/core/common"
	"dmc-task/model/crontasks"
	"dmc-task/model/jobsflow"
	"dmc-task/model/lock"
	"dmc-task/server"
	"dmc-task/utils/gopool"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
	"sync"
	"testing"
	"time"
)

func newMysqlConn() *sqlx.SqlConn {
	// 创建Mysql客户端
	sqlConn, err := NewMysql(
		"root",
		"Shanhai*123",
		"10.30.4.229",
		"dmc_task",
		3306)
	if err != nil {
		logx.Errorf("new mysql conn is failed! err:%v", err)
		os.Exit(1)
	}
	server.NewServiceContext(nil)
	server.SvrCtx.MysqlConn = sqlConn
	return sqlConn
}

func cleanCronTasks(conn *sqlx.SqlConn) {
	m := crontasks.NewTCronTasksModel(*conn)
	cleanCronTasksSql := fmt.Sprintf("DELETE FROM %s", m.GetTableName())
	_, _ = m.ExecSql(context.Background(), cleanCronTasksSql)
}

func cleanJobsFlow(conn *sqlx.SqlConn) {
	m := jobsflow.NewTJobsFlowModel(*conn)
	cleanCronTasksSql := fmt.Sprintf("DELETE FROM %s", m.GetTableName())
	_, _ = m.ExecSql(context.Background(), cleanCronTasksSql)
}

func cleanAll(conn *sqlx.SqlConn) {
	cleanCronTasks(conn)
	cleanJobsFlow(conn)
}

func TestClearAll(t *testing.T) {
	conn := newMysqlConn()
	cleanAll(conn)
	t.Log("clean all success!")
}

func TestGoPool(t *testing.T) {
	p := gopool.NewPool("xxx", 100, gopool.NewConfig())
	type Result struct {
		n    int32
		Name string
	}
	rs := make([]Result, 100)
	for i := 0; i < 100; i++ {
		rs[i] = Result{Name: fmt.Sprintf("name-%d", i), n: int32(i)}
	}
	var wg sync.WaitGroup
	for i, r := range rs {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			fmt.Println(p.WorkerCount(), r, i)
		})
	}
	fmt.Println(p.WorkerCount())
	wg.Wait()
}

func TestDistributedLocks(t *testing.T) {
	conn := newMysqlConn()
	var wg sync.WaitGroup
	lockName := "xxx"
	m := lock.NewTDistributedLocksModel(*conn)
	m.Reset(lockName)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				if ok := m.Lock(lockName); ok {
					t.Logf("[%d] lock success!", idx)
					time.Sleep(1 * time.Second)
					if err := m.Renew(lockName); err != nil {
						t.Logf("[%d] lock renew success!", idx)
					}
					time.Sleep(1 * time.Second)
					m.Unlock(lockName)
					t.Logf("[%d] lock release success!", idx)
					break
				}
				t.Logf("[%d] lock failed, retry after 1 second", idx)
				time.Sleep(1 * time.Second)
			}
		}(i)
	}
	wg.Wait()
}

func TestDistributedLocksRest(t *testing.T) {
	lockName := "xxxxx"
	conn := newMysqlConn()
	m := lock.NewTDistributedLocksModel(*conn)
	if ok := m.Lock(lockName); ok {
		t.Log("lock success!")
	} else {
		t.Log("lock failed!")
	}
	//m.Unlock()
	//t.Log("lock release success!")
	m.Reset(lockName)
	if ok := m.Lock(lockName); ok {
		t.Log("lock success!")
	} else {
		t.Log("lock failed!")
	}
}

func TestPaginate(t *testing.T) {
	conn := newMysqlConn()
	m := jobsflow.NewTJobsFlowModel(*conn)
	queryReq := &PaginationRequest{
		Ctx:      context.Background(),
		Conn:     *conn,
		Table:    m.GetTableName(),
		Where:    "status = 3",
		OrderBy:  "create_time DESC",
		Args:     []interface{}{},
		Page:     3,
		PageSize: 100,
	}
	result, err := Paginate[jobsflow.TJobsFlow](queryReq)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(result.Count, len(result.Data))
	for i, v := range result.Data {
		t.Logf("[%d-%d] %d, %+v", result.Count, i, v.Status, v.CreateTime)
	}
}

func printData[T any](result *PaginationResult[T]) {
	fmt.Printf("count:%d, size:%d\n", result.Count, len(result.Data))
	for i, v := range result.Data {
		fmt.Printf("[%d-%d] %+v\n", result.Count, i, v)
	}
}

func TestQueryForJobsFlow(t *testing.T) {
	filter := common.FilterBase{
		Id:         "",
		BizCode:    "",
		BizId:      "",
		CronTaskId: "f8a3a960-86a8-4b9c-b565-6c6d5387a7b8",
		Status:     3,
		TimeType:   "finish_time",
		Start:      "2025-02-28 06:23:38",
		End:        "2025-02-28 06:24:38",
	}
	page := common.PageBase{
		Page:     1,
		PageSize: 100,
	}

	table := jobsflow.NewTJobsFlowModel(*newMysqlConn()).GetTableName()
	result, err := Query[jobsflow.TJobsFlow](context.Background(), table, filter, page)
	if err != nil {
		t.Log(err)
		return
	}
	//printData[jobsflow.TJobsFlow](result)
	t.Log(result.Count, len(result.Data))
	for i, v := range result.Data {
		t.Logf("[%d-%d] status:%d, create:%+v, update:%+v, start:%+v, finish:%+v",
			result.Count, i, v.Status, v.CreateTime, v.UpdateTime, v.StartTime, v.FinishTime)
	}
}
