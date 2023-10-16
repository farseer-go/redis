package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisPub struct {
	rdb *redis.Client
}

func (receiver *redisPub) Publish(channel string, message any) (int64, error) {
	trace := linkTrace.TraceRedis("Publish", channel, "")
	result, err := receiver.rdb.Publish(fs.Context, channel, message).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (receiver *redisPub) Subscribe(channels ...string) <-chan *redis.Message {
	trace := linkTrace.TraceRedis("Subscribe", strings.Join(channels, ","), "")
	defer func() { trace.End(nil) }()
	return receiver.rdb.Subscribe(fs.Context, channels...).Channel()
}
