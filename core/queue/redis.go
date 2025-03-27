package queue

import (
	"context"
	"dmc-task/server"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultStream   = "mystream"
	DefaultGroup    = "mygroup"
	DefaultConsumer = "myconsumer"
)

// CreateConsumerGroup 通过指定的stream和group创建Redis消费者组
func CreateConsumerGroup(ctx context.Context, stream, group string) error {
	return createGroup(ctx, stream, group)
}

// ProductMsgToQueue 向指定的stream中生产消息msg
func ProductMsgToQueue(ctx context.Context, stream string, msg interface{}) error {
	return productMsg(ctx, stream, msg)
}

// ReadMsgFromQueue 从指定的stream中读取消息
func ReadMsgFromQueue(ctx context.Context, stream, group, consumer string) (*redis.XMessage, error) {
	return readMsg(ctx, stream, group, consumer)
}

func AckMsgToQueue(ctx context.Context, stream, group, id string) error {
	return ackMsg(ctx, stream, group, id)
}

func createGroup(ctx context.Context, stream, group string) error {
	// TODO: Implement Redis consumer group
	_, err := server.SvrCtx.RedisClient.XGroupCreate(ctx, stream, group, "$").Result()
	if err != nil {
		logx.Error(err)
		return err
	}
	return nil
}

func productMsg(ctx context.Context, stream string, msg interface{}) error {
	// TODO: Implement Redis product message to queue
	_, err := server.SvrCtx.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: msg,
		MaxLen: 100,
		Approx: true,
	}).Result()
	if err != nil {
		logx.Errorf("Failed to add message to stream: %v", err)
		return err
	}
	return nil
}

func readMsg(ctx context.Context, stream, group, consumer string) (*redis.XMessage, error) {
	// TODO: Implement Redis consumer group read
	msgs, err := server.SvrCtx.RedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    1,
		Block:    1000, // 阻塞1000毫秒
	}).Result()
	if err != nil {
		if err == redis.Nil {
			// 超时，没有新消息
			return nil, nil
		}
		logx.Errorf("Failed to read from stream: %v", err)
		return nil, err
	}
	logx.Debugf("Read message from stream: %+v", msgs)
	for _, msg := range msgs[0].Messages {
		logx.Debugf("Read message from stream, id:%s, value:%+v", msg.ID, msg.Values)
		return &msg, nil
	}
	return nil, nil
}

func ackMsg(ctx context.Context, stream, group, id string) error {
	// TODO: Implement Redis consumer group ack message
	_, err := server.SvrCtx.RedisClient.XAck(ctx, stream, group, id).Result()
	if err != nil {
		logx.Errorf("Failed to ack message from stream: %v", err)
		return err
	}
	return nil
}
