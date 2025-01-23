//go:build windows

package command

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
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
		logx.Error(err)
		return
	}
	cmdStr := fmt.Sprintf("%s %s", commandName, strings.Join(cmd.Args, " "))
	logx.Debugf("pid:%d, command:%s, output:%s", cmd.Process.Pid, cmdStr, string(output))
	data = strings.Split(string(output), "\n")
	return data, nil
}

func kill(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		logx.Debugf("[kill cmd for windows] kill process pid:%d, cmd:%s", -cmd.Process.Pid, cmd.String())
		p, err := os.FindProcess(cmd.Process.Pid)
		if err != nil {
			logx.Error(err)
			return err
		}
		return p.Signal(syscall.SIGTERM)
		//return p.Kill()
	}
	return nil
}
