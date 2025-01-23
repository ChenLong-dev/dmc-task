package command

import (
	"context"
	"testing"
	"time"
)

func TestExecCommand(t *testing.T) {
	// 测试用例1: 正常执行命令
	t.Run("正常执行命令", func(t *testing.T) {
		commandName := "echo"
		params := []string{"hello", "world"}
		var timeout int64 = 5

		data, err := ExecCommand(context.Background(), timeout, commandName, params)
		if err != nil {
			t.Errorf("执行命令失败: %v", err)
		}

		expected := []string{"hello world", ""}
		if len(data) != len(expected) || data[0] != expected[0] {
			t.Errorf("输出结果不匹配, 期望: %v, 实际: %v", expected, data)
		}
	})

	// 测试用例2: 命令执行超时
	t.Run("命令执行超时", func(t *testing.T) {
		commandName := "sleep"
		params := []string{"10"}
		var timeout int64 = 3

		_, err := ExecCommand(context.Background(), timeout, commandName, params)
		if err == nil {
			t.Error("预期命令不存在应返回错误，但未返回错误")
		}
	})

	// 测试用例3: 命令不存在
	t.Run("命令不存在", func(t *testing.T) {
		commandName := "nonexistent_command"
		params := []string{}
		var timeout int64 = 5

		_, err := ExecCommand(context.Background(), timeout, commandName, params)
		if err == nil {
			t.Error("预期命令不存在应返回错误，但未返回错误")
		}
	})

	// 测试用例4: 命令执行成功但无输出
	t.Run("命令执行成功但无输出", func(t *testing.T) {
		commandName := "true"
		params := []string{}
		var timeout int64 = 5

		data, err := ExecCommand(context.Background(), timeout, commandName, params)
		if err != nil {
			t.Errorf("执行命令失败: %v", err)
		}

		if len(data) == 0 {
			t.Errorf("预期无输出，但实际输出: %v", data)
		}
	})

	// 测试用例5: 命令执行失败
	t.Run("命令执行失败", func(t *testing.T) {
		commandName := "false"
		params := []string{}
		var timeout int64 = 5

		_, err := ExecCommand(context.Background(), timeout, commandName, params)
		if err == nil {
			t.Error("预期命令执行失败应返回错误，但未返回错误")
		}
	})

	// 测试用例6: 命令执行时间小于超时时间
	t.Run("命令执行时间小于超时时间", func(t *testing.T) {
		commandName := "sleep"
		params := []string{"1"}
		var timeout int64 = 5

		startTime := time.Now()
		_, err := ExecCommand(context.Background(), timeout, commandName, params)
		elapsedTime := time.Since(startTime)

		if err != nil {
			t.Errorf("执行命令失败: %v", err)
		}

		if elapsedTime >= time.Duration(timeout)*time.Second {
			t.Errorf("命令应在超时前完成，但实际耗时: %v", elapsedTime)
		}
	})
}

func TestExecMultiCommand(t *testing.T) {
	//commandName := "../../scripts/cmd.sh"
	commandName := "../../cmd/dmctask/dmctask.exe"
	//params := []string{"-b", "world", "xxxxx", "{\"a\":1, \"b\":2}"}
	params := []string{"version"}
	timeout := int64(3)
	_, err := ExecCommand(context.Background(), timeout, commandName, params)
	if err != nil {
		t.Errorf("执行命令失败: %v", err)

	}
}
