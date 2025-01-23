package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	// 执行该命令时，会执行的函数
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("this is test cobra example")
		version()
	},
}

func version() {
	fmt.Println("sleep 5 seconds ...")
	time.Sleep(5 * time.Second)
}
