package utils

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestGetUTCTime(t *testing.T) {
	// 测试UTC时间获取的正确性
	utcTime := GetUTCTime()
	if utcTime.Location() != time.UTC {
		t.Errorf("期望时间是UTC，实际获取到的时间是%v", utcTime.Location())
	}
}

func TestGetUTCTime2(t *testing.T) {
	// 测试添加时间间隔后UTC时间的正确性
	duration := 2 * time.Hour
	utcTime := GetUTCTime2(duration)
	expectedTime := GetUTCTime().Add(duration)
	if !utcTime.Equal(expectedTime) {
		t.Errorf("期望时间是%v，实际获取到的时间是%v", expectedTime, utcTime)
	}
}

func TestGetLocalTime(t *testing.T) {
	// 测试本地时间获取的正确性
	localTime := GetLocalTime()
	if localTime.Location() != time.Local {
		t.Errorf("期望时间是本地时间，实际获取到的时间是%v", localTime.Location())
	}
}

func TestGetLocalTime2(t *testing.T) {
	// 测试添加时间间隔后本地时间的正确性
	duration := 2 * time.Hour
	localTime := GetLocalTime2(duration)
	expectedTime := GetLocalTime().Add(duration)
	if !localTime.Equal(expectedTime) {
		t.Errorf("期望时间是%v，实际获取到的时间是%v", expectedTime, localTime)
	}
}

func TestGetTime(t *testing.T) {
	// 测试通过时间戳转换为时间的正确性
	timestamp := int64(1609459200) // 2021-01-01 00:00:00 UTC
	expectedTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	actualTime := GetTime(timestamp)
	if !actualTime.Equal(expectedTime) {
		t.Errorf("期望时间是%v，实际获取到的时间是%v", expectedTime, actualTime)
	}
}

func TestGetTimestamp(t *testing.T) {
	// 测试通过时间转换为时间戳的正确性
	expectedTimestamp := int64(1609459200) // 2021-01-01 00:00:00 UTC
	actualTimestamp := GetTimestamp(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	if actualTimestamp != expectedTimestamp {
		t.Errorf("期望时间戳是%v，实际获取到的时间戳是%v", expectedTimestamp, actualTimestamp)
	}
}

func TestGetTimeStr(t *testing.T) {
	// 测试时间格式化为字符串的正确性
	expectedTimeStr := "2021-01-01 00:00:00.000000000 +0000 UTC"
	actualTimeStr := GetTimeStr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	if actualTimeStr != expectedTimeStr {
		t.Errorf("期望时间字符串是%v，实际获取到的时间字符串是%v", expectedTimeStr, actualTimeStr)
	}
}

func TestGetRandInt(t *testing.T) {
	// 测试生成随机数的正确性
	min, max := 1, 10
	randInt := GetRandInt(min, max)
	if randInt < min || randInt > max {
		t.Errorf("期望随机数在%v到%v之间，实际获取到的随机数是%v", min, max, randInt)
	}
}

func TestMarshalByJson(t *testing.T) {
	// 测试结构体序列化为JSON的正确性
	testStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Alice",
		Age:  30,
	}
	expectedJson := `{"name":"Alice","age":30}`
	actualJson, err := MarshalByJson(testStruct)
	if err != nil {
		t.Errorf("序列化过程中发生错误: %v", err)
	}
	if string(actualJson) != expectedJson {
		t.Errorf("期望JSON字符串是%v，实际获取到的JSON字符串是%v", expectedJson, string(actualJson))
	}
}

func TestUnmarshalByJson(t *testing.T) {
	// 测试JSON字符串反序列化为结构体的正确性
	jsonStr := `{"name":"Bob","age":25}`
	expectedStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Bob",
		Age:  25,
	}
	var actualStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalByJson([]byte(jsonStr), &actualStruct)
	if err != nil {
		t.Errorf("反序列化过程中发生错误: %v", err)
	}
	if actualStruct != expectedStruct {
		t.Errorf("期望结构体是%v，实际获取到的结构体是%v", expectedStruct, actualStruct)
	}
}

func TestGetLocalIP(t *testing.T) {
	// 测试获取本地IP的正确性
	// 这个测试可能在不同的环境中得到不同的结果，因此这里我们只检查返回的字符串是否为空或格式是否正确
	ip := GetLocalIP()
	if ip == "" {
		t.Errorf("获取到的本地IP为空")
	}
	t.Log("获取到的本地IP是", ip)
	// 检查是否为有效的IPv4地址
	parts := strings.Split(ip, "-")
	for _, part := range parts {
		if net.ParseIP(part).To4() == nil {
			t.Errorf("获取到的IP地址%v不是有效的IPv4地址", part)
		}
	}
}
