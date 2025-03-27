//go:build linux || darwin
// +build linux darwin

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
	logx.WithContext(ctx).Debugf("[execCommand] timeout:%d, exec command:%s, params:%v", timeout, commandName, params)
	data = data[0:0]
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandName, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go func() {
		select {
		case <-ctx.Done():
			_ = kill(cmd)
		}
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}
	cmdStr := fmt.Sprintf("%s %s", commandName, strings.Join(cmd.Args, " "))
	logx.WithContext(ctx).Debugf("pid:%d, command:%s, output:%s", cmd.Process.Pid, cmdStr, string(output))
	data = strings.Split(string(output), "\n")
	return data, nil
}

func kill(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		pid := cmd.Process.Pid
		//logx.Infof("[kill cmd for unix] kill process pid:%d, cmd:%s", -pid, cmd.String())
		err := syscall.Kill(-pid, syscall.SIGKILL)
		if err != nil {
			//logx.Error(err)
			return err
		}
	}
	return nil
}
