package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type redisKey struct {
	rdb *redis.Client
}

func (redisKey *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	return redisKey.rdb.Expire(ctx, key, d).Result()
}

func (redisKey *redisKey) TTL(key string) (time.Duration, error) {
	return redisKey.rdb.TTL(ctx, key).Result()
}

func (redisKey *redisKey) Del(keys ...string) (bool, error) {
	result, err := redisKey.rdb.Del(ctx, keys...).Result()
	return result > 0, err
}

func (redisKey *redisKey) Exists(keys ...string) (bool, error) {
	result, err := redisKey.rdb.Exists(ctx, keys...).Result()
	return result > 0, err
}
