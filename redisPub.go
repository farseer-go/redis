package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisPub struct {
	*redisManager
}

func (receiver *redisPub) Publish(channel string, message any) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("Publish", channel, "")
	result, err := receiver.GetClient().Publish(context.Background(), channel, message).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisPub) Subscribe(channels ...string) <-chan *redis.Message {
	traceDetail := receiver.traceManager.TraceRedis("Subscribe", strings.Join(channels, ","), "")
	defer func() { traceDetail.End(nil) }()
	return receiver.rdb.Subscribe(context.Background(), channels...).Channel()
}
