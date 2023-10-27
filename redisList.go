package redis

import (
	"github.com/farseer-go/fs"
	"strings"
	"time"
)

type redisList struct {
	*redisManager
}

func (receiver *redisList) ListPushRight(key string, values ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListPushRight", key, "")

	result, err := receiver.GetClient().RPush(fs.Context, key, values).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisList) ListPushLeft(key string, values ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListPushLeft", key, "")

	result, err := receiver.GetClient().LPush(fs.Context, key, values).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisList) ListSet(key string, index int64, value any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListSet", key, "")

	result, err := receiver.GetClient().LSet(fs.Context, key, index, value).Result()
	defer func() { traceDetail.End(err) }()
	if result == "OK" {
		return true, err
	}
	return false, err
}

func (receiver *redisList) ListRemove(key string, count int64, value any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRemove", key, "")

	result, err := receiver.GetClient().LRem(fs.Context, key, count, value).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisList) ListCount(key string) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListCount", key, "")

	result, err := receiver.GetClient().LLen(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisList) ListRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRange", key, "")

	result, err := receiver.GetClient().LRange(fs.Context, key, start, stop).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisList) ListLeftPop(timeout time.Duration, keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListLeftPop", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().BLPop(fs.Context, timeout, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisList) ListRightPop(timeout time.Duration, keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRightPop", strings.Join(keys, ","), "")

	result, err := receiver.GetClient().BRPop(fs.Context, timeout, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisList) ListRightPopPush(source, destination string, timeout time.Duration) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ListRightPopPush", source, "")

	result, err := receiver.GetClient().BRPopLPush(fs.Context, source, destination, timeout).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}
