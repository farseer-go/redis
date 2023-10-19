package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type redisKey struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

func (receiver *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetTTL", key, "")

	result, err := receiver.rdb.Expire(fs.Context, key, d).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisKey) TTL(key string) (time.Duration, error) {
	traceDetail := receiver.traceManager.TraceRedis("TTL", key, "")

	result, err := receiver.rdb.TTL(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisKey) Del(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Del", strings.Join(keys, ","), "")

	result, err := receiver.rdb.Del(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisKey) Exists(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Exists", strings.Join(keys, ","), "")

	result, err := receiver.rdb.Exists(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
