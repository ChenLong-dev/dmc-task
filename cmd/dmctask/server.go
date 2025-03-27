package main

import (
	"fmt"
	"os"
	"time"

	internal "dmc-task/cmd/dmctask/api"
	"dmc-task/cmd/dmctask/grpc"
	core "dmc-task/core"
	"dmc-task/core/cron"
	"dmc-task/core/gracefulstop"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/server"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
)

var cfgPath string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "用于启动服务",
	Long:  `dmc-task server的启动服务命令`,

	// 命令执行前，读取环境变量，读取配置文件
	PreRun: func(cmd *cobra.Command, args []string) {
		err := core.ConfigInit(cfgPath, &core.Cfg)
		if err != nil {
			fmt.Printf("init config file is failed! err:%v\n", err)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		// 初始化日志（logx）
		logxInit()

		// 初始化数据库
		_ = server.NewServiceContext(core.Cfg)
		model.InitMysql()

		// 初始化分布式锁
		_ = model.Reset()
		if model.Lock() {
			logx.Infof("this server is master, source:%s", server.SvrCtx.IsMasterSource)
			// 启动时间轮（只有master才启动）
			timewheel.Start()
		} else {
			logx.Info("this server is slave!")
		}

		// 初始化定时任务
		cron.Start()
		logx.Debugf("cfg:%+v", core.Cfg)

		// 初始化服务
		stop := make(chan bool, 1)
		gracefulstop.GracefulShutdown(stop)
		go grpc.Run(core.Cfg)
		go internal.Run(core.Cfg)
		select {
		case <-stop:
			logx.Info("dmc-task is closed!")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	serverCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "the path of the config file")
}

func logxInit() {
	c := logx.LogConf{
		ServiceName: core.Cfg.App.Name,
		Mode:        core.Cfg.Logx.Mode,
		Encoding:    core.Cfg.Logx.Encoding,
		TimeFormat:  time.RFC3339,
		Path:        core.Cfg.Logx.Path,
		Level:       core.Cfg.Logx.Level,
		KeepDays:    core.Cfg.Logx.KeepDays,
		MaxBackups:  core.Cfg.Logx.MaxBackups,
		MaxSize:     core.Cfg.Logx.MaxSize,
		Rotation:    core.Cfg.Logx.Rotation,
	}
	logx.MustSetup(c)
	logx.Infof("logx init is success! logx:%+v", core.Cfg.Logx)
	return
}
