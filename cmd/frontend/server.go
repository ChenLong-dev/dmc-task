/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"dmc-task/cmd/frontend/biz"
	"dmc-task/cmd/frontend/internal"
	"dmc-task/core"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

var cfgPath string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,

	// 命令执行前，读取环境变量，读取配置文件
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("input config path:%s\n", cfgPath)
		err := core.FrontendConfigInit(cfgPath)
		if err != nil {
			fmt.Printf("init config file is failed! err:%v\n", err)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		// 初始化日志（logx）
		logxInit()

		logx.Debugf("cfg:%+v", core.FrontendCfg)

		// 初始化grpc客户端
		InitGRPCClient()

		// 初始化服务
		biz.Run(core.FrontendCfg)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	serverCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "the path of the config file")
}

func logxInit() {
	c := logx.LogConf{
		ServiceName: core.FrontendCfg.App.Name,
		Mode:        core.FrontendCfg.Logx.Mode,
		Encoding:    core.FrontendCfg.Logx.Encoding,
		TimeFormat:  time.DateTime,
		Path:        core.FrontendCfg.Logx.Path,
		Level:       core.FrontendCfg.Logx.Level,
		KeepDays:    core.FrontendCfg.Logx.KeepDays,
		MaxBackups:  core.FrontendCfg.Logx.MaxBackups,
		MaxSize:     core.FrontendCfg.Logx.MaxSize,
		Rotation:    core.FrontendCfg.Logx.Rotation,
	}
	logx.MustSetup(c)
	logx.Infof("logx init is success! logx:%+v", core.FrontendCfg.Logx)
	return
}

func InitGRPCClient() {
	addr := "127.0.0.1:7889"
	// 创建连接器
	var err error
	internal.GrpcClient, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil || internal.GrpcClient == nil {
		logx.Error(err)
		panic(err)
	}
	logx.Infof("grpc client init is success! addr:%s, state:%s", addr, internal.GrpcClient.GetState())
}
