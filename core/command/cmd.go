package command

import (
	"context"
	"runtime"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultTimeout = 15
)

func ExecCommand(ctx context.Context, timeout int64, commandName string, params []string) (data []string, err error) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	switch runtime.GOOS {
	case "windows":
		data, err = execCommand(ctx, timeout, commandName, params)
	case "linux":
		data, err = execCommand(ctx, timeout, commandName, params)
	case "darwin":
		data, err = execCommand(ctx, timeout, commandName, params)
	default:
		logx.WithContext(ctx).Error("not support os")
	}
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return
	}

	return
}
