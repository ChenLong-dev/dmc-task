package main

import (
	internal "dmc-task/cmd/dmctask/api"
	core "dmc-task/core"
	"dmc-task/core/cron"
	"dmc-task/core/gracefulstop"
	"dmc-task/core/timewheel"
	"dmc-task/model"
	"dmc-task/server"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

var cfgPath string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "用于启动服务",
	Long:  `dmc-task server的启动服务命令`,

	// 命令执行前，读取环境变量，读取配置文件
	PreRun: func(cmd *cobra.Command, args []string) {
		err := core.ConfigInit(cfgPath)
		if err != nil {
			fmt.Printf("init config file is failed! err:%v\n", err)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		// 初始化日志（logx）
		logxInit()

		// 初始化服务
		_ = server.NewServiceContext(core.Cfg)
		model.InitMysql()
		_ = model.Reset()
		if model.Lock() {
			logx.Infof("this server is master, source:%s", server.SvrCtx.IsMasterSource)
			timewheel.Start() // 启动时间轮（只有master才启动）
		} else {
			logx.Info("this server is slave!")
		}
		cron.Start() // 初始化定时任务
		logx.Debugf("cfg:%+v", core.Cfg)

		// 开始服务
		gracefulstop.GracefulShutdown()
		internal.Run(core.Cfg)
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
		TimeFormat:  time.DateTime,
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
