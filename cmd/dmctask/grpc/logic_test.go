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
	addr = "127.0.0.1:8889"
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
		Id:    "",
		Limit: 1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("resp:%+v", resp)
}
