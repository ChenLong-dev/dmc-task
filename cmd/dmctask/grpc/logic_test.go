package grpc

import (
	"context"
	protoc "dmc-task/cmd/dmctask/grpc/task"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

const (
	addr = "10.30.4.231:7889"
)

func getConn() (conn *grpc.ClientConn) {
	// 1.创建连接器
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return
}

func closeConn(conn *grpc.ClientConn) {
	if conn == nil {
		fmt.Println("conn == nil")
		return
	}
	if err := conn.Close(); err != nil {
		fmt.Printf("grpc conn closed is error, err:%+v", err)
		return
	}
	fmt.Println("grpc conn closed is successful!")
}

func TestAddRealTimeSingleTask(t *testing.T) {
	conn := getConn()
	defer closeConn(conn)

	client := protoc.NewTaskClient(conn)
	resp, err := client.AddRealTimeSingleTask(context.Background(), &protoc.AddRealTimeSingleTaskReq{
		Task: &protoc.RealTimeSingleTask{
			Type:     1,
			BizCode:  "222888",
			BizId:    "222222",
			ExecPath: "ls",
			Param:    "-al",
			Timeout:  15,
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("resp:%+v", resp)

}

func TestQueryRealTimeSingleTask(t *testing.T) {
	conn := getConn()
	defer closeConn(conn)

	client := protoc.NewTaskClient(conn)
	resp, err := client.QueryRealTimeSingleTask(context.Background(), &protoc.QueryRealTimeSingleTaskReq{
		Filter: &protoc.FilterBase{},
		Page:   &protoc.PageBase{},
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("resp:%+v", resp)
}

func TestStartOrStopCronCycleTask(t *testing.T) {
	conn := getConn()
	defer closeConn(conn)

	client := protoc.NewTaskClient(conn)
	resp, err := client.StartOrStopCronCycleTask(context.Background(), &protoc.StartOrStopCronCycleTaskReq{
		Id:      "bf8d34ae-9865-4965-ab9e-d186b42b1e9c",
		IsStart: false,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("resp:%+v", resp)
}
