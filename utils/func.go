package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"math"
	"time"
)

// WithRetry 添加一个通用的重试函数
func WithRetry(ctx context.Context, maxRetries int, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			// 指数退避，每次重试等待时间增加
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			logc.Info(ctx, "retry operation",
				logx.Field("attempt", i+1),
				logx.Field("wait_time", waitTime),
				logx.Field("err", err))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
			}
		}

		err = operation()
		if err == nil {
			return nil
		}

		// 判断是否是不需要重试的错误
		if errors.Is(err, context.Canceled) {
			return err
		}
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

// SplitSlice splits a slice into smaller slices of size num
func SplitSlice[T any](slice []T, num int) [][]T {
	if num <= 0 {
		return nil
	}

	var chunks [][]T
	length := len(slice)

	for i := 0; i < length; i += num {
		end := i + num
		if end > length {
			end = length
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
