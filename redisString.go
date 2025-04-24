package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisString struct {
	*redisManager
}

// StringSet 设置缓存
func (receiver *redisString) StringSet(key string, value any) error {
	traceDetail := receiver.traceManager.TraceRedis("StringSet", key, "")
	err := receiver.GetClient().Set(context.Background(), key, value, 0).Err()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return err
}

// StringGet 获取缓存
func (receiver *redisString) StringGet(key string) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringGet", key, "")
	result, err := receiver.GetClient().Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result, err
}

// StringSetEX 设置过期时间（如果KEY存在，则会覆盖）
func (receiver *redisString) StringSetEX(key string, value any, expiration time.Duration) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringSetNX", key, "")
	result, err := receiver.GetClient().SetEX(context.Background(), key, value, expiration).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result, err
}

// StringSetNX 设置过期时间（如果KEY存在，则会更新失败）
func (receiver *redisString) StringSetNX(key string, value any, expiration time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringSetNX", key, "")
	result, err := receiver.GetClient().SetNX(context.Background(), key, value, expiration).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil && result {
		traceDetail.SetRows(1)
	}
	return result, err
}

// StringIncrInt String对int加减
func (receiver *redisString) StringIncrInt(key string, value int) (int, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringIncrInt", key, "")
	result, err := receiver.GetClient().IncrBy(context.Background(), key, int64(value)).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil && result > 0 {
		traceDetail.SetRows(1)
	}
	return int(result), err
}

// StringIncrInt64 String对int64加减
func (receiver *redisString) StringIncrInt64(key string, value int64) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringIncrInt64", key, "")
	result, err := receiver.GetClient().IncrBy(context.Background(), key, value).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil && result > 0 {
		traceDetail.SetRows(1)
	}
	return result, err
}

// StringIncrFloat32 String对float32加减
func (receiver *redisString) StringIncrFloat32(key string, value float32) (float32, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringIncrFloat32", key, "")
	result, err := receiver.GetClient().IncrByFloat(context.Background(), key, float64(value)).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil && result > 0 {
		traceDetail.SetRows(1)
	}
	return float32(result), err
}

// StringIncrFloat64 String对float64加减
func (receiver *redisString) StringIncrFloat64(key string, value float64) (float64, error) {
	traceDetail := receiver.traceManager.TraceRedis("StringIncrFloat64", key, "")
	result, err := receiver.GetClient().IncrByFloat(context.Background(), key, value).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil && result > 0 {
		traceDetail.SetRows(1)
	}
	return result, err
}
