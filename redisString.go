package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisString struct {
	rdb *redis.Client
}

// StringSet 设置缓存
func (redisString *redisString) StringSet(key string, value any) error {
	trace := linkTrace.TraceRedis("StringSet", key, "")
	err := redisString.rdb.Set(fs.Context, key, value, 0).Err()
	defer func() { trace.End(err) }()
	return err
}

// StringGet 获取缓存
func (redisString *redisString) StringGet(key string) (string, error) {
	trace := linkTrace.TraceRedis("StringGet", key, "")
	result, err := redisString.rdb.Get(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

// StringSetNX 设置过期时间
func (redisString *redisString) StringSetNX(key string, value any, expiration time.Duration) (bool, error) {
	trace := linkTrace.TraceRedis("StringSetNX", key, "")
	result, err := redisString.rdb.SetNX(fs.Context, key, value, expiration).Result()
	defer func() { trace.End(err) }()
	return result, err
}
