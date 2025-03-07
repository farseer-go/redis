package redis

import (
	"context"
	"strings"
	"time"

	"github.com/farseer-go/fs/parse"
	"github.com/go-redis/redis/v8"
)

type redisList struct {
	*redisManager
}

func (receiver *redisList) ListPushRight(key string, values ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListPushRight", key, "")

	result, err := receiver.GetClient().RPush(context.Background(), key, values).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result > 0, err
}

func (receiver *redisList) ListPushLeft(key string, values ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListPushLeft", key, "")

	result, err := receiver.GetClient().LPush(context.Background(), key, values).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result > 0, err
}

func (receiver *redisList) ListSet(key string, index int64, value any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListSet", key, "")

	result, err := receiver.GetClient().LSet(context.Background(), key, index, value).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result == "OK", err
}

func (receiver *redisList) ListRemove(key string, count int64, value any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRemove", key, "")

	result, err := receiver.GetClient().LRem(context.Background(), key, count, value).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result > 0, err
}

func (receiver *redisList) ListCount(key string) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListCount", key, "")

	result, err := receiver.GetClient().LLen(context.Background(), key).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result, err
}

func (receiver *redisList) ListRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRange", key, "")

	result, err := receiver.GetClient().LRange(context.Background(), key, start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}

func (receiver *redisList) ListLeftPop(timeout time.Duration, keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListLeftPop", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().BLPop(context.Background(), timeout, keys...).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}

func (receiver *redisList) ListRightPop(timeout time.Duration, keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRightPop", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().BRPop(context.Background(), timeout, keys...).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}

func (receiver *redisList) ListRightPopPush(source, destination string, timeout time.Duration) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRightPopPush", source, "")

	result, err := receiver.GetClient().BRPopLPush(context.Background(), source, destination, timeout).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result, err
}
