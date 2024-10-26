package redis

import (
	"context"
	"strings"
	"time"

	"github.com/farseer-go/fs/parse"
)

type redisKey struct {
	*redisManager
}

func (receiver *redisKey) SetTTL(key string, d time.Duration) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetTTL", key, "")

	result, err := receiver.GetClient().Expire(context.Background(), key, d).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result, err
}

func (receiver *redisKey) TTL(key string) (time.Duration, error) {
	traceDetail := receiver.traceManager.TraceRedis("TTL", key, "")

	result, err := receiver.GetClient().TTL(context.Background(), key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisKey) Del(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Del", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().Del(context.Background(), keys...).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result > 0, err
}

func (receiver *redisKey) Exists(keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("Exists", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().Exists(context.Background(), keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
