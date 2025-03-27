package queue

import (
	"context"
	"dmc-task/core"
	"dmc-task/server"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

const (
	stream   = "teststream"
	group    = "testgroup"
	consumer = "testconsumer"
)

func NewRedisClient() {
	// 创建Redis客户端
	_ = core.ConfigInit("../../cmd/dmctask/conf/conf.yaml", &core.Cfg)
	server.NewServiceContext(core.Cfg)
	return
}

// 测试队列功能
func TestQueue(t *testing.T) {
	NewRedisClient()
	ctx := context.Background()
	// 1. 测试创建消费者组 - 正常情况
	err := CreateConsumerGroup(ctx, stream, group)
	assert.NoError(t, err, "创建消费者组失败")

	// 2. 测试生产消息 - 正常情况
	msg := map[string]interface{}{"field": "value"}
	err = ProductMsgToQueue(ctx, stream, msg)
	assert.NoError(t, err, "生产消息失败")

	// 3. 测试读取消息 - 正常情况
	readMsg, err := ReadMsgFromQueue(ctx, stream, group, consumer)
	assert.NoError(t, err, "读取消息失败")
	assert.NotNil(t, readMsg, "读取到的消息为空")

	// 4. 测试确认消息 - 正常情况
	err = AckMsgToQueue(ctx, stream, group, readMsg.ID)
	assert.NoError(t, err, "确认消息失败")

	// 5. 测试重复确认消息 - 边界情况
	err = AckMsgToQueue(ctx, stream, group, readMsg.ID)
	assert.NoError(t, err, "重复确认消息不应该发生错误")

	// 6. 测试读取没有新消息 - 边界情况
	_, err = ReadMsgFromQueue(ctx, stream, group, consumer)
	assert.NoError(t, err, "读取没有新消息时发生错误")

	// 7. 测试创建同名消费者组 - 边界情况
	err = CreateConsumerGroup(ctx, stream, group)
	assert.Error(t, err, "不应该允许创建同名消费者组")

	// 清理测试数据
	_, err = server.SvrCtx.RedisClient.Del(ctx, stream).Result()
	assert.NoError(t, err, "清理测试数据失败")
	_, err = server.SvrCtx.RedisClient.XGroupDestroy(ctx, stream, group).Result()
	assert.NoError(t, err, "销毁消费者组失败")
}

func TestRedisExample(t *testing.T) {
	// 连接到Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "shanhai888888",  // 密码（如果有的话）
		DB:       0,                // 使用默认DB
	})

	ctx := context.Background()

	// 创建Stream
	t.Logf("开始创建Stream")
	//_, err := rdb.XAdd(ctx, &redis.XAddArgs{
	//	Stream: "mystream",
	//	Values: map[string]interface{}{
	//		"field1": "value1",
	//		"field2": "value2",
	//	},
	//}).Result()
	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}).Result()
	if err != nil {
		t.Logf("Failed to add message to stream: %v", err)
	}

	// 创建Consumer Group
	t.Logf("开始创建Consumer Group")
	_, err = rdb.XGroupCreate(ctx, "mystream", "mygroup", "$").Result()
	_, err = rdb.XGroupCreate(ctx, stream, group, "$").Result()
	if err != nil && err != redis.Nil {
		t.Logf("Failed to create consumer group: %v", err)
	}
	t.Logf("[created] Consumer group created 成功！")

	// 消费者读取消息
	go func() {
		for {
			t.Logf("开始消费消息")
			msgs, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "mygroup",
				Consumer: "myconsumer",
				Streams:  []string{"mystream", ">"},
				Count:    1,
				Block:    1000, // 阻塞1000毫秒
			}).Result()

			//msgs, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			//	Group:    group,
			//	Consumer: consumer,
			//	Streams:  []string{"mystream", ">"},
			//	Count:    1,
			//	Block:    1000, // 阻塞1000毫秒
			//}).Result()
			if err != nil {
				if err == redis.Nil {
					// 超时，没有新消息
					continue
				}
				t.Logf("Failed to read from stream: %v", err)
			}

			for _, msg := range msgs[0].Messages {
				t.Logf("[received] Received: %s %s", msg.ID, msg.Values)

				// 确认消息
				//_, err = rdb.XAck(ctx, "mystream", "mygroup", msg.ID).Result()
				_, err = rdb.XAck(ctx, stream, group, msg.ID).Result()
				if err != nil {
					t.Logf("Failed to ack message: %v", err)
				}
				t.Logf("[acked] Message %s acked 成功！", msg.ID)
			}
		}
	}()

	// 模拟生产者继续发送消息
	for i := 0; i < 5; i++ {
		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "mystream",
			Values: map[string]interface{}{
				"field1": fmt.Sprintf("value%d", i+1),
				"field2": "another value",
			},
			MaxLen: 100,
			Approx: true,
		}).Result()
		if err != nil {
			log.Fatalf("Failed to add message to stream: %v", err)
		}
		t.Logf("[produced] Message produced 成功！i:%d", i)
		time.Sleep(2 * time.Second) // 模拟生产间隔
	}

	//time.Sleep(10 * time.Second)
	//// 清理测试数据
	//_, err := core.SvrCtx.RedisClient.Del(ctx, stream).Result()
	//if err != nil {
	//	t.Logf("清理测试数据失败")
	//}
	//t.Logf("清理数据成功！ %s", stream)
	//_, err = core.SvrCtx.RedisClient.XGroupDestroy(ctx, stream, group).Result()
	//if err != nil {
	//	t.Logf("销毁消费者组失败")
	//}
	//t.Logf("销毁消费者组成功！ %s %s", stream, group)
}
