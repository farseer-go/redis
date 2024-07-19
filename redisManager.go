package redis

import (
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
)

// 实现同一个协程下的事务作用域
var routineRedisClient = asyncLocal.New[redis.Cmdable]()

type redisManager struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

func (receiver *redisManager) GetClient() redis.Cmdable {
	cmdClient := routineRedisClient.Get()
	if cmdClient != nil {
		return cmdClient
	}
	return receiver.rdb
}
