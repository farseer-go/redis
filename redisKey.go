package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type redisKey struct {
	rdb *redis.Client
}

func (redisKey *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	trace := linkTrace.TraceRedis("SetTTL", key, "")

	result, err := redisKey.rdb.Expire(fs.Context, key, d).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisKey *redisKey) TTL(key string) (time.Duration, error) {
	trace := linkTrace.TraceRedis("TTL", key, "")

	result, err := redisKey.rdb.TTL(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisKey *redisKey) Del(keys ...string) (bool, error) {
	trace := linkTrace.TraceRedis("Del", strings.Join(keys, ","), "")

	result, err := redisKey.rdb.Del(fs.Context, keys...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisKey *redisKey) Exists(keys ...string) (bool, error) {
	trace := linkTrace.TraceRedis("Exists", strings.Join(keys, ","), "")

	result, err := redisKey.rdb.Exists(fs.Context, keys...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}
