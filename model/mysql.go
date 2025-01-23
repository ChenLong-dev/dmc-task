package model

import (
	"context"
	"dmc-task/model/crontasks"
	"dmc-task/model/jobsflow"
	"dmc-task/model/lock"
	"dmc-task/server"
	"dmc-task/utils"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

func InitMysql() {
	logx.Debug("init mysql... ", server.SvrCtx.Config.MySQL)
	config := server.SvrCtx.Config.MySQL
	conn, err := NewMysql(config.Username, config.Password, config.Host, config.Database, config.Port)
	if err != nil {
		panic(err)
	}
	server.SvrCtx.MysqlConn = conn

	sqlx.DisableStmtLog() // 禁用sql语句日志
	logx.Infof("[mysql] init success! config:%+v", config)
}

func NewMysql(username, password, host, database string, port int) (*sqlx.SqlConn, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		username, password, host, port, database)
	logx.Debugf("mysql connecting! dataSource:%s", dataSource)
	mysql := sqlx.NewMysql(dataSource)
	db, err := mysql.RawDB()
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = db.PingContext(cxt)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return &mysql, nil
}

func Lock() bool {
	source := fmt.Sprintf("%s-%s:%d", utils.GetLocalIP(), server.SvrCtx.Config.Server.Host, server.SvrCtx.Config.Server.Port)
	if ok := lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Lock(source); !ok {
		return false
	}
	server.SvrCtx.IsMasterSource = source
	return true
}

func Unlock() {
	if server.SvrCtx.IsMasterSource == "" {
		return
	}
	lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Unlock(server.SvrCtx.IsMasterSource)
}

func Renew() error {
	if server.SvrCtx.IsMasterSource == "" {
		err := fmt.Errorf("not found master source")
		return err
	}
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Renew(server.SvrCtx.IsMasterSource)
}

func Reset() error {
	return lock.NewTDistributedLocksModel(*server.SvrCtx.MysqlConn).Reset(server.SvrCtx.IsMasterSource)
}

func generateCronTask() *crontasks.TCronTasks {
	//rand := utils.GetRandInt(10, 20)
	//status := core.Pending
	//bid := fmt.Sprintf("bid_%d", rand)
	//return &crontasks.TCronTasks{
	//	Id:            uuid.New().String(),
	//	Type:          int64(core.FixedTimeSingleTask),
	//	BizCode:           bid,
	//
	//	Param: "{}",
	//	ExecParams:    "{}",
	//	StartTime:     sql.NullTime{},
	//	FinishTime:    sql.NullTime{},
	//	ExecTime:      utils.GetUTCTime2(time.Duration(rand) * time.Second),
	//	ExecInterval:  0,
	//	Status:        int64(status),
	//	ResultMsg:     core.GetResult(core.Success.Code, bid, "待执行....", status),
	//}
	return nil
}

func generateJobFlow(cronTask *crontasks.TCronTasks) *jobsflow.TJobsFlow {
	//if cronTask == nil {
	//	status := core.Pending
	//	result := core.GetResult(core.Success.Code, "", "待执行....", status)
	//	return &jobsflow.TJobsFlow{
	//		Id:            uuid.New().String(),
	//		Type:          int64(core.RealTimeSingleTask),
	//		CronTaskId:    "",
	//		Bid:           "",
	//		CallbackParam: "{}",
	//		ExecParams:    "{}",
	//		StartTime:     utils.GetUTCTime(),
	//		Status:        int64(status),
	//		ResultMsg:     result,
	//	}
	//} else {
	//	status := core.Running
	//	result := core.GetResult(core.Success.Code, cronTask.BizId, "进行中....", status)
	//	return &jobsflow.TJobsFlow{
	//		Id:            uuid.New().String(),
	//		Type:          cronTask.Type,
	//		CronTaskId:    cronTask.Id,
	//		Bid:           cronTask.BizCode,
	//		CallbackParam: cronTask.CallbackParam,
	//		ExecParams:    cronTask.ExecParams,
	//		StartTime:     utils.GetUTCTime(),
	//		Status:        int64(status),
	//		ResultMsg:     result,
	//	}
	//}
	return nil
}

// GetAndUpdateCronTasksAndJobsFlow 获取并更新定时任务和Job流水记录表的状态的事务
func GetAndUpdateCronTasksAndJobsFlow(ctx context.Context) ([]*jobsflow.TJobsFlow, error) {
	var jobFlows []*jobsflow.TJobsFlow
	//mt := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	//mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	//// 带Pending状态查询t_cron_tasks表中待执行的任务
	//cronTasks, err := mt.GetCronTasksByStatus(ctx, int64(core.Pending))
	//if err != nil {
	//	return jobFlows, err
	//}
	//if len(cronTasks) == 0 {
	//	return nil, nil
	//}
	//logx.Debugf("[1] select cronTask:%+v", len(cronTasks))

	//for i, cronTask := range cronTasks {
	//	jobFlow := generateJobFlow(cronTask)
	//	err = mt.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
	//		// 生成Job流水记录表
	//		cronTask.StartTime = sql.NullTime{utils.GetUTCTime(), true}
	//		cronTask.Status = int64(core.Running)
	//		cronTask.ResultMsg = core.GetResult(core.Success.Code, cronTask.Bid, "进行中....", core.Running)
	//
	//		// 更新t_cron_tasks表状态为Running
	//		queryTask := fmt.Sprintf("update %s set %s where `id` = ?", mt.GetTableName(), mt.GetCronTasksRowsWithPlaceHolder())
	//		resTask, err := session.ExecCtx(ctx, queryTask, cronTask.Type, cronTask.Bid, cronTask.CallbackParam,
	//			cronTask.ExecParams, cronTask.StartTime, cronTask.FinishTime, cronTask.ExecTime, cronTask.ExecInterval,
	//			cronTask.Status, cronTask.ResultMsg, cronTask.Id)
	//		if err != nil {
	//			return err
	//		}
	//		lastID, _ := resTask.LastInsertId()
	//		logx.Debugf("[2-%d] update t_cron_tasks table lastID:%d", i, lastID)
	//
	//		// 插入t_jobs_flow表
	//		queryJob := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", mj.GetTableName(),
	//			mj.GetJobsFlowRowsExpectAutoSet())
	//		retJob, err := session.ExecCtx(ctx, queryJob, jobFlow.Id, jobFlow.Type, jobFlow.CronTaskId, jobFlow.Bid,
	//			jobFlow.CallbackParam, jobFlow.ExecParams, jobFlow.StartTime, jobFlow.FinishTime, jobFlow.ExecInterval,
	//			jobFlow.Status, jobFlow.ResultMsg)
	//		if err != nil {
	//			return err
	//		}
	//		lastID, _ = retJob.LastInsertId()
	//		logx.Debugf("[3-%d] insert t_jobs_flow table lastID:%d", i, lastID)
	//		return nil
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	jobFlows = append(jobFlows, jobFlow)
	//}
	return jobFlows, nil
}

// GetAndUpdateJobFlow 获取并更新Job流水记录表的状态的事务
func GetAndUpdateJobFlow(ctx context.Context) (*jobsflow.TJobsFlow, error) {
	//var jobFlow *jobsflow.TJobsFlow
	//var err error
	//mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	//err = mj.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
	//	// 带Pending状态查询t_cron_tasks表中待执行的任务
	//	jobFlow, err = mj.GetJobFlowByStatus(ctx, int64(core.Pending))
	//	if err != nil {
	//		return err
	//	}
	//	logx.Debugf("[1] select jobFlow:%+v", jobFlow)
	//
	//	// 生成Job流水记录表
	//	jobFlow.StartTime = utils.GetUTCTime()
	//	jobFlow.Status = int64(core.Running)
	//	jobFlow.ResultMsg = core.GetResult(core.Success.Code, jobFlow.Bid, "进行中....", core.Running)
	//
	//	// 更新t_jobs_flow表
	//	queryJob := fmt.Sprintf("update %s set %s where `id` = ?", mj.GetTableName(), mj.GetJobsFlowRowsWithPlaceHolder())
	//	resJob, err := session.ExecCtx(ctx, queryJob, jobFlow.Type, jobFlow.CronTaskId, jobFlow.Bid, jobFlow.CallbackParam,
	//		jobFlow.ExecParams, jobFlow.StartTime, jobFlow.FinishTime, jobFlow.ExecInterval, jobFlow.Status,
	//		jobFlow.ResultMsg, jobFlow.Id)
	//	if err != nil {
	//		return err
	//	}
	//	lastID, _ := resJob.LastInsertId()
	//	logx.Debugf("[2] update t_jobs_flow table lastID:%d", lastID)
	//	return nil
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return jobFlow, nil
	return nil, nil
}

// GetAndUpdateJobsFlow 获取并更新Job流水记录表的状态的事务
func GetAndUpdateJobsFlow(ctx context.Context) ([]*jobsflow.TJobsFlow, error) {
	//var err error
	//mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	//// 带Pending状态查询t_cron_tasks表中待执行的任务
	//jobsFlow, err := mj.GetJobsFlowByStatus(ctx, int64(core.Pending))
	//if err != nil {
	//	if err != jobsflow.ErrNotFound {
	//		return nil, err
	//	}
	//}
	//if len(jobsFlow) == 0 {
	//	return nil, nil
	//}
	//logx.Debugf("[1] select jobsFlow[%d]", len(jobsFlow))
	//
	//for i, jobFlow := range jobsFlow {
	//	// 生成Job流水记录表
	//	jobFlow.StartTime = utils.GetUTCTime()
	//	jobFlow.Status = int64(core.Running)
	//	jobFlow.ResultMsg = core.GetResult(core.Success.Code, jobFlow.Bid, "进行中....", core.Running)
	//
	//	// 更新t_jobs_flow表
	//	err = mj.Update(ctx, jobFlow)
	//	if err != nil {
	//		return nil, err
	//	}
	//	logx.Debugf("[2-%d] update t_jobs_flow table, id:%s", i, jobFlow.Id)
	//}
	//return jobsFlow, nil
	return nil, nil
}

// UpdateCronTasksAndJobsFlow 更新定时任务和Job流水记录表的状态的事务
func UpdateCronTasksAndJobsFlow(ctx context.Context, cronTask *crontasks.TCronTasks, jobFlow *jobsflow.TJobsFlow) error {
	//mt := crontasks.NewTCronTasksModel(*server.SvrCtx.MysqlConn)
	//mj := jobsflow.NewTJobsFlowModel(*server.SvrCtx.MysqlConn)
	//err := mj.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
	//	// 更新t_jobs_flow表
	//	queryJob := fmt.Sprintf("update %s set %s where `id` = ?", mj.GetTableName(), mj.GetJobsFlowRowsWithPlaceHolder())
	//	resJob, err := session.ExecCtx(ctx, queryJob, jobFlow.Type, jobFlow.CronTaskId, jobFlow.Bid, jobFlow.CallbackParam,
	//		jobFlow.ExecParams, jobFlow.StartTime, jobFlow.FinishTime, jobFlow.ExecInterval, jobFlow.Status,
	//		jobFlow.ResultMsg, jobFlow.Id)
	//	if err != nil {
	//		return err
	//	}
	//	lastID, _ := resJob.LastInsertId()
	//	logx.Debugf("[update-1] update t_jobs_flow table lastID:%d, id:%s", lastID, jobFlow.Id)
	//
	//	// 更新t_cron_tasks表
	//	queryTask := fmt.Sprintf("update %s set %s where `id` = ?", mt.GetTableName(), mt.GetCronTasksRowsWithPlaceHolder())
	//	resTask, err := session.ExecCtx(ctx, queryTask, cronTask.Type, cronTask.Bid, cronTask.CallbackParam, cronTask.ExecParams,
	//		cronTask.StartTime, cronTask.FinishTime, cronTask.ExecTime, cronTask.ExecInterval, cronTask.Status,
	//		cronTask.ResultMsg, cronTask.Id)
	//	if err != nil {
	//		return err
	//	}
	//	lastID, _ = resTask.LastInsertId()
	//	logx.Debugf("[update-2] update t_jobs_flow table lastID:%d, id:%s", lastID, jobFlow.Id)
	//	return nil
	//})
	//if err != nil {
	//	logx.Error(err)
	//	return err
	//}
	return nil
}
