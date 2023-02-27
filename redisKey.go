package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisKey struct {
	rdb *redis.Client
}

func (redisKey *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	return redisKey.rdb.Expire(fs.Context, key, d).Result()
}

func (redisKey *redisKey) TTL(key string) (time.Duration, error) {
	return redisKey.rdb.TTL(fs.Context, key).Result()
}

func (redisKey *redisKey) Del(keys ...string) (bool, error) {
	result, err := redisKey.rdb.Del(fs.Context, keys...).Result()
	return result > 0, err
}

func (redisKey *redisKey) Exists(keys ...string) (bool, error) {
	result, err := redisKey.rdb.Exists(fs.Context, keys...).Result()
	return result > 0, err
}
