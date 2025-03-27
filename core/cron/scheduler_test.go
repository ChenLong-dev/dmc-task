package cron

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/server"
	"dmc-task/utils"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func newMysqlConn() *sqlx.SqlConn {
	// 创建Mysql客户端
	sqlConn, err := model.NewMysql(
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

func TestAddCronTask(t *testing.T) {
	newMysqlConn()
	Start()
	taskParam := common.CronCycleTask{
		Type:     int64(core.CronCycleTask),
		Cron:     "*/5 * * * * *",
		BizCode:  "test",
		ExecPath: "echo hello",
		Param:    "",
		Timeout:  5,
	}
	id1, err := AddTask(taskParam)
	if err != nil {
		t.Errorf("add cron task failed! err:%v", err)
	}
	taskParam.Cron = "@every 5m"
	id2, err := AddTask(taskParam)
	if err != nil {
		t.Errorf("add cron task failed! err:%v", err)
	}
	t.Logf("id1:%d, id2:%d", id1, id2)

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	// 停止Cron实例
	Stop()
	t.Log("程序已退出")
}

func TestRemoveCronTask(t *testing.T) {
	newMysqlConn()
	Start()
	taskParam := common.CronCycleTask{
		Type:     int64(core.CronCycleTask),
		Cron:     "*/5 * * * * *",
		BizCode:  "test",
		ExecPath: "echo hello",
		Param:    "",
		Timeout:  5,
	}
	id1, err := AddTask(taskParam)
	if err != nil {
		t.Errorf("add cron task failed! id1:%d, err:%v", id1, err)
	}
	taskParam.Cron = "@every 5m"
	id2, err := AddTask(taskParam)
	if err != nil {
		t.Errorf("add cron task failed! id1:%d, err:%v", id1, err)
	}
	t.Logf("id1:%d, id2:%d", id1, id2)
	time.Sleep(time.Second * 5)
	RemoveTask(context.Background(), id1)
	taskParam.Cron = "*/2 * * * * *"
	id3, err := AddTask(taskParam)
	if err != nil {
		t.Errorf("add cron task failed! id1:%d, err:%v", id1, err)
	}
	t.Logf("id3:%d", id3)

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	// 停止Cron实例
	Stop()
	t.Log("程序已退出")
}

// ---------------------------------------------------------
func TestAddFixedTimeTask(t *testing.T) {
	newMysqlConn()
	timewheel.Start()
	req := common.AddFixedTimeSingleTaskReq{
		FixedTimeSingleTask: common.FixedTimeSingleTask{
			Type:     int64(core.FixedTimeSingleTask),
			BizCode:  "test",
			BizId:    "1234567890",
			ExecPath: "sleep",
			Param:    "5",
			ExecTime: utils.GetUTCTime().Add(time.Second * 65).Unix(),
			Timeout:  10,
		},
	}
	taskId, err := AddDataToCronTasks(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 5)
	t.Logf("taskId:%s", taskId)
}

// ---------------------------------------------------------

func TestAddRealTimeTask(t *testing.T) {
	newMysqlConn()
	ctx := context.Background()
	taskParam := common.RealTimeSingleTask{
		Type:     int64(core.CronCycleTask),
		BizCode:  "test",
		ExecPath: "sleep",
		Param:    "3",
		Timeout:  5,
	}
	s, err := addRealTimeTask(ctx, taskParam)
	if err != nil {
		t.Errorf("add real time task failed! err:%v", err)
	}
	t.Logf("s:%s", s)
	time.Sleep(time.Second * 15)
}

// ---------------------------------------------------------
func operation(ctx context.Context, duration time.Duration) {
	now := time.Now()
	select {
	case <-time.After(duration): // duration时间后开始执行后面操作
		for i := 0; i < 5; i++ {
			fmt.Println("Operation running - ", i)
			time.Sleep(time.Millisecond * 500)
		}
		fmt.Println("Operation finished")
	case <-ctx.Done(): // 接收到取消信号后开始执行后面操作
		fmt.Println("Operation cancelled: ", time.Now().Sub(now))
	}
}

/*
WithCancel函数返回一个新的上下文对象和一个取消函数。调用这个取消函数将取消这个上下文对象，以及从它派生的所有上下文对象。
WithCancel用于创建可以被手动取消的上下文。这对于告知goroutine停止当前工作并及时退出非常有用。
*/
func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		operation(ctx, 5*time.Second)
	}()

	time.Sleep(3 * time.Second) // 模拟在操作完成前进行取消
	cancel()                    // 取消操作

	// 等待足够长的时间以确保goroutine可以响应取消事件
	time.Sleep(1 * time.Second)
	t.Log("程序已退出")
}

/*
WithDeadline函数返回一个新的上下文对象，这个对象会在指定的截止时间自动取消。
WithDeadline用于创建具有明确截止时间的上下文。当达到截止时间时，上下文会自动取消。这对于设置任务的最长执行时间非常有用。
*/
func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Operation finished")
	case <-ctx.Done():
		fmt.Println("Operation cancelled due to deadline")
	}
	t.Log("程序已退出")
}

func TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Operation finished")
	case <-ctx.Done():
		fmt.Println("Operation cancelled due to timeout")
	}
	t.Log("程序已退出")
}
