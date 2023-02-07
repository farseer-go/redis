package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type redisString struct {
	rdb *redis.Client
}

// StringSet 设置缓存
func (redisString *redisString) StringSet(key string, value any) error {
	return redisString.rdb.Set(ctx, key, value, 0).Err()
}

// StringGet 获取缓存
func (redisString *redisString) StringGet(key string) (string, error) {
	return redisString.rdb.Get(ctx, key).Result()
}

// StringSetNX 设置过期时间
func (redisString *redisString) StringSetNX(key string, value any, expiration time.Duration) (bool, error) {
	return redisString.rdb.SetNX(ctx, key, value, expiration).Result()
}
