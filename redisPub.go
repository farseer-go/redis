package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
)

type redisPub struct {
	rdb *redis.Client
}

func (receiver *redisPub) Publish(channel string, message any) (int64, error) {
	return receiver.rdb.Publish(fs.Context, channel, message).Result()
}

func (receiver *redisPub) Subscribe(channels ...string) <-chan *redis.Message {
	return receiver.rdb.Subscribe(fs.Context, channels...).Channel()
}
