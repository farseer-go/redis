package redis

import (
	"github.com/farseer-go/fs"
	"strings"
	"time"
)

type redisKey struct {
	*redisManager
}

func (receiver *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetTTL", key, "")

	result, err := receiver.GetClient().Expire(fs.Context, key, d).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisKey) TTL(key string) (time.Duration, error) {
	traceDetail := receiver.traceManager.TraceRedis("TTL", key, "")

	result, err := receiver.GetClient().TTL(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisKey) Del(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Del", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().Del(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisKey) Exists(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Exists", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().Exists(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
