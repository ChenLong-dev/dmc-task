package model

import (
	"context"
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
		"127.0.0.1",
		"dmc_task",
		23306)
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
