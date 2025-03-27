package utils

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestWithRetry(t *testing.T) {
	tests := []struct {
		name       string
		maxRetries int
		operation  func() error
		wantErr    bool
		errMsg     string
		setup      func() context.Context
	}{
		{
			name:       "successful operation without retry",
			maxRetries: 3,
			operation: func() error {
				return nil
			},
			wantErr: false,
			setup: func() context.Context {
				return context.Background()
			},
		},
		{
			name:       "successful operation after retry",
			maxRetries: 3,
			operation: (func() func() error {
				attempts := 0
				return func() error {
					attempts++
					if attempts < 2 {
						return errors.New("temporary error")
					}
					return nil
				}
			})(),
			wantErr: false,
			setup: func() context.Context {
				return context.Background()
			},
		},
		{
			name:       "operation fails all retries",
			maxRetries: 2,
			operation: func() error {
				return errors.New("persistent error")
			},
			wantErr: true,
			errMsg:  "failed after 2 retries",
			setup: func() context.Context {
				return context.Background()
			},
		},
		{
			name:       "context cancellation",
			maxRetries: 3,
			operation: func() error {
				return errors.New("some error")
			},
			wantErr: true,
			errMsg:  "context canceled",
			setup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // 立即取消上下文
				return ctx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			err := WithRetry(ctx, tt.maxRetries, tt.operation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WithRetry() expected error, got nil")
				}
				if tt.errMsg != "" && !errors.Is(err, context.Canceled) {
					if got := err.Error(); !contains(got, tt.errMsg) {
						t.Errorf("WithRetry() error = %v, want %v", got, tt.errMsg)
					}
				}
			} else if err != nil {
				t.Errorf("WithRetry() unexpected error: %v", err)
			}
		})
	}
}

func TestSplitSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			size:     2,
			expected: [][]int{},
		},
		{
			name:     "size larger than slice",
			input:    []int{1, 2, 3},
			size:     5,
			expected: [][]int{{1, 2, 3}},
		},
		{
			name:     "even split",
			input:    []int{1, 2, 3, 4},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}},
		},
		{
			name:     "uneven split",
			input:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "size zero",
			input:    []int{1, 2, 3},
			size:     0,
			expected: nil,
		},
		{
			name:     "size negative",
			input:    []int{1, 2, 3},
			size:     -1,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitSlice(tt.input, tt.size)

			if tt.expected == nil && result != nil {
				t.Errorf("SplitSlice() = %v, want nil", result)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("SplitSlice() length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if !sliceEqual(result[i], tt.expected[i]) {
					t.Errorf("SplitSlice() chunk %d = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// 辅助函数

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func sliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
