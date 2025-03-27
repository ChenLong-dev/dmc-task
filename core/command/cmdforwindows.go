//go:build windows

package command

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func execCommand(ctx context.Context, timeout int64, commandName string, params []string) (data []string, err error) {
	data = data[0:0]
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandName, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = kill(cmd)
		}
	}()
	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	cmdStr := fmt.Sprintf("%s %s", commandName, strings.Join(cmd.Args, " "))
	logx.WithContext(ctx).Debugf("pid:%d, command:%s, output:%s", cmd.Process.Pid, cmdStr, string(output))
	data = strings.Split(string(output), "\n")
	return data, nil
}

//func kill(cmd *exec.Cmd) error {
//	if cmd != nil && cmd.Process != nil {
//		logx.Debugf("[kill cmd for windows] kill process pid:%d, cmd:%s", -cmd.Process.Pid, cmd.String())
//		p, err := os.FindProcess(cmd.Process.Pid)
//		if err != nil {
//			logx.Error(err)
//			return err
//		}
//		return p.Signal(syscall.SIGTERM)
//		//return p.Kill()
//	}
//	return nil
//}

func kill(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		pid := cmd.Process.Pid
		logx.Debugf("[kill cmd for windows] attempting to kill process pid:%d, cmd:%s", pid, cmd.String())

		// 首先检查进程是否存在
		checkCmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		checkOutput, _ := checkCmd.CombinedOutput()

		// 如果输出中不包含 PID，说明进程已经不存在
		if !strings.Contains(string(checkOutput), fmt.Sprintf("%d", pid)) {
			logx.Debugf("Process %d already terminated", pid)
			return nil
		}

		// 使用 taskkill 命令强制终止进程
		killCmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
		output, err := killCmd.CombinedOutput()
		if err != nil {
			// 如果进程不存在，不认为是错误
			if strings.Contains(string(output), "没有找到进程") ||
				strings.Contains(string(output), "not found") {
				logx.Debugf("Process %d already terminated", pid)
				return nil
			}
			logx.Errorf("Failed to kill process %d: %v, output: %s", pid, err, string(output))
			return err
		}

		logx.Debugf("Successfully terminated process pid:%d, output: %s", pid, string(output))
		return nil
	}
	return nil
}
