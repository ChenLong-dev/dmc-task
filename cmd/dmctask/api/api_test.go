package internal

import (
	"bytes"
	"dmc-task/core"
	"dmc-task/core/common"
	"dmc-task/core/middleware"
	"dmc-task/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	addr = "127.0.0.1:7888"
)

// 封装HTTP POST请求的函数
func sendPostRequest(url string, headerData map[string]string, postData map[string]interface{}) (string, error) {
	// 将map[string]interface{}转换为JSON格式的字节数组
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}

	// 创建一个请求体
	requestBody := bytes.NewBuffer(jsonData)

	// 创建一个POST请求
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return "", err
	}

	// 设置Content-Type为application/json，表明发送的是JSON数据
	for k, v := range headerData {
		req.Header.Set(k, v)
	}
	//req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	// 返回响应体字符串和任何错误
	return bodyString, nil
}

func getHeader() (header map[string]string) {
	header = make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + middleware.GetTaskSecretKey()
	return header
}

func addCronTask(req common.AddFixedTimeSingleTaskReq) {
	postData := make(map[string]interface{})
	postData["type"] = req.Type
	postData["biz_code"] = req.BizCode
	postData["biz_id"] = req.BizId
	postData["exec_path"] = req.ExecPath
	postData["param"] = req.Param
	postData["exec_time"] = req.ExecTime
	postData["timeout"] = req.Timeout
	postData["ext_info"] = req.ExtInfo
	resp, err := sendPostRequest(fmt.Sprintf("http://%s/%s", addr, "v1/cron/add"), getHeader(), postData)
	if err != nil {
		fmt.Printf("--- [%s] biz_code:%s, biz_id:%s is error: %v\n", core.TaskTypeMap[core.TaskType(req.Type)],
			req.BizCode, req.BizId, err)
		return
	}
	fmt.Printf("+++ [%s] biz_code:%s, biz_id:%s, resp:%s\n", core.TaskTypeMap[core.TaskType(req.Type)],
		req.BizCode, req.BizId, resp)
}

func TestAddCronTasks(t *testing.T) {
	execTime := utils.GetUTCTime2(160*time.Second).Unix() - 8*60*60
	fmt.Printf("utc:%v,%d, timestamp:%d\n", time.Now().UTC(), time.Now().UTC().Unix(), execTime)
	fmt.Printf("xxxxxx: %d, %v\n", time.Now().Unix()-time.Now().UTC().Unix(), time.Now().UTC())
	req := common.AddFixedTimeSingleTaskReq{}
	req.Type = 2
	req.BizCode = "xxxx"
	req.BizId = "111111"
	req.ExecPath = "ls"
	req.Param = "-al"
	req.ExecTime = execTime
	req.Timeout = 5
	req.ExtInfo = "{}"
	addCronTask(req)
}

func addJobTask(req common.AddRealTimeSingleTaskReq) {
	postData := make(map[string]interface{})
	postData["type"] = req.Type
	postData["biz_code"] = req.BizCode
	postData["biz_id"] = req.BizId
	postData["exec_path"] = req.ExecPath
	postData["param"] = req.Param
	postData["timeout"] = req.Timeout
	postData["ext_info"] = req.ExtInfo
	_, err := sendPostRequest(fmt.Sprintf("http://%s/%s", addr, "v1/job/add"), getHeader(), postData)
	if err != nil {
		fmt.Printf("--- [%s] biz_code:%s, biz_id:%s is error: %v\n", core.TaskTypeMap[core.TaskType(req.Type)],
			req.BizCode, req.BizId, err)
		return
	}
	//fmt.Printf("+++ [%s] biz_code:%s, biz_id:%s, resp:%s\n", core.TaskTypeMap[core.TaskType(req.Type)],
	// req.BizCode, req.BizId, resp)
}

func TestAddJobTasks(t *testing.T) {
	bizCode := "343434343434"
	mod := 100
	for i := 0; i < 50000; i++ {
		req := common.AddRealTimeSingleTaskReq{}
		req.Type = 1
		req.BizCode = bizCode
		req.BizId = fmt.Sprintf("2222%d", i)
		req.ExecPath = "ls"
		req.Param = "-al"
		req.Timeout = 5
		req.ExtInfo = fmt.Sprintf("{\"i\": %d}", i)
		addJobTask(req)
		if i%mod == 0 {
			time.Sleep(2 * time.Second)
		}
	}
}

func BenchmarkAddJobTasks(b *testing.B) {
	bizCode := "BenchmarkAddJobTasks-1"
	for i := 0; i < b.N; i++ {
		req := common.AddRealTimeSingleTaskReq{}
		req.Type = 1
		req.BizCode = bizCode
		req.BizId = fmt.Sprintf("%d", i+1)
		req.ExecPath = "ls"
		req.Param = "-al"
		req.Timeout = 5
		req.ExtInfo = fmt.Sprintf("{\"i\": %d}", i+1)
		addJobTask(req)
	}
}
