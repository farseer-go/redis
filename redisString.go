package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisString struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

// StringSet 设置缓存
func (receiver *redisString) StringSet(key string, value any) error {
	traceDetail := receiver.traceManager.TraceRedis("StringSet", key, "")
	err := receiver.rdb.Set(fs.Context, key, value, 0).Err()
	defer func() { traceDetail.End(err) }()
	return err
}

// StringGet 获取缓存
func (receiver *redisString) StringGet(key string) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringGet", key, "")
	result, err := receiver.rdb.Get(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

// StringSetNX 设置过期时间
func (receiver *redisString) StringSetNX(key string, value any, expiration time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringSetNX", key, "")
	result, err := receiver.rdb.SetNX(fs.Context, key, value, expiration).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}
