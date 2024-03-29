package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisString struct {
	*redisManager
}

// StringSet 设置缓存
func (receiver *redisString) StringSet(key string, value any) error {
	traceDetail := receiver.traceManager.TraceRedis("StringSet", key, "")
	err := receiver.GetClient().Set(fs.Context, key, value, 0).Err()
	defer func() { traceDetail.End(err) }()
	return err
}

// StringGet 获取缓存
func (receiver *redisString) StringGet(key string) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringGet", key, "")
	result, err := receiver.GetClient().Get(fs.Context, key).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

// StringSetNX 设置过期时间
func (receiver *redisString) StringSetNX(key string, value any, expiration time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringSetNX", key, "")
	result, err := receiver.GetClient().SetNX(fs.Context, key, value, expiration).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}
