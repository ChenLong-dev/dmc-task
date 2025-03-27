package cron

import (
	"context"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/model/lock"
	"dmc-task/server"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var c *cron.Cron

func init() {
	c = cron.New(cron.WithSeconds())
}

func Start() {
	addTasks()

	c.Start()
}

func Stop() {
	c.Stop()
}

func addTasks() {
	if server.SvrCtx.IsMasterSource != "" { // 只有master才添加定时循环任务和定时任务扫描
		// 第一次启动，初次添加定时循环任务
		_ = addCronCycleInitTasks()
		// 添加数据库循环扫描任务（定时循环任务和固定循环任务）
		_ = addCronScanFromDBTask()
	}
	//添加分布式锁轮询任务
	time.Sleep(time.Millisecond * 500)
	_ = AddLockTask()

}

func execFunc(taskParam common.CronCycleTask) func() {
	// TODO: 执行任务
	return func() {
		logx.Debugf("[execTask] %+v", taskParam)
		execCronCycle(taskParam)
		entriesPrint(taskParam)
	}
}

func entriesPrint(taskParam common.CronCycleTask) {
	for _, entry := range c.Entries() {
		logx.Debugf(">>>> ID:%d, Delay:%+v, job:%+v, wrappedJob:%+v, BizCode:%+v", entry.ID, entry.Schedule, entry.Job,
			entry.WrappedJob, taskParam.BizCode)
	}
}

func addDynamicTask(taskParam common.CronCycleTask) (int64, error) {
	entryId, err := c.AddFunc(taskParam.Cron, execFunc(taskParam))
	logx.Infof("[addDynamicTask] entryId: %d", entryId)
	return int64(entryId), err
}

func removeDynamicTask(ctx context.Context, entryId cron.EntryID) {
	c.Remove(entryId)
	logx.WithContext(ctx).Infof("[removeDynamicTask] entryId: %d", entryId)
}

func AddTask(taskParam common.CronCycleTask) (int64, error) {
	return addDynamicTask(taskParam)
}

func RemoveTask(ctx context.Context, entryId int64) {
	removeDynamicTask(ctx, cron.EntryID(entryId))
}

func AddLockTask() error {
	if !server.SvrCtx.Config.App.IsDistributed { // 判断是否开启分布式锁
		return nil
	}
	spec := fmt.Sprintf("@every %ds", lock.LockExpire/2)
	entryId, err := c.AddFunc(spec, func() {
		logx.Debug("续期....")
		if err := model.Renew(); err == nil { // 续期成功
			return
		}
		logx.Debug("抢锁....")
		if model.Lock() { // 抢锁成功
			logx.Infof("this server is master, get source:%s", server.SvrCtx.IsMasterSource)
			logx.Debug("启动时间轮...")
			timewheel.Start()                       // 启动时间轮
			if server.SvrCtx.IsMasterSource != "" { // 只有master才添加定时循环任务和定时任务扫描
				// 添加数据库循环扫描任务（定时循环任务和固定循环任务）
				logx.Debug("添加数据库循环扫描任务cron ....")
				// 第一次启动，初次添加定时循环任务
				_ = addCronCycleInitTasks()
				// 添加数据库循环扫描任务（定时循环任务和固定循环任务）
				_ = addCronScanFromDBTask()
			}
			return
		}
		logx.Debug("重置锁....")
		if err := model.Reset(); err != nil {
			logx.Error(err)
		}
	})
	if err != nil {
		logx.Error(err)
		return err
	}
	logx.Infof("[AddLockTask] spec：%s, entryId: %d", spec, entryId)
	return nil
}

func addCronScanFromDBTask() error {
	addCronScanFromDB()
	entryId, err := c.AddFunc(core.FixCycleSpec, addCronScanFromDB)
	if err != nil {
		logx.Info("add CronScanFromDBTask: ", err)
		return err
	}

	logx.Infof("[add CronScanFromDBTask] spec:%s, entryId: %d", core.FixCycleSpec, entryId)
	return nil
}

func addCronScanFromDB() {
	// 添加/删除/更新定时循环任务scheduler
	addCronCycleTasks()
	// 添加固定时间任务scheduler
	addFixedTimeSingleTasksFromDB()
	return
}
