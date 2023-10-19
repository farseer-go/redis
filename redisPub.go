package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisPub struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

func (receiver *redisPub) Publish(channel string, message any) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("Publish", channel, "")
	result, err := receiver.rdb.Publish(fs.Context, channel, message).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisPub) Subscribe(channels ...string) <-chan *redis.Message {
	traceDetail := receiver.traceManager.TraceRedis("Subscribe", strings.Join(channels, ","), "")
	defer func() { traceDetail.End(nil) }()
	return receiver.rdb.Subscribe(fs.Context, channels...).Channel()
}
