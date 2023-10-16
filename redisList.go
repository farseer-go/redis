package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type redisList struct {
	rdb *redis.Client
}

func (redisList *redisList) ListPushRight(key string, values ...any) (bool, error) {
	trace := linkTrace.TraceRedis("ListPushRight", key, "")

	result, err := redisList.rdb.RPush(fs.Context, key, values).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisList *redisList) ListPushLeft(key string, values ...any) (bool, error) {
	trace := linkTrace.TraceRedis("ListPushLeft", key, "")

	result, err := redisList.rdb.LPush(fs.Context, key, values).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisList *redisList) ListSet(key string, index int64, value any) (bool, error) {
	trace := linkTrace.TraceRedis("ListSet", key, "")

	result, err := redisList.rdb.LSet(fs.Context, key, index, value).Result()
	defer func() { trace.End(err) }()
	if result == "OK" {
		return true, err
	}
	return false, err
}

func (redisList *redisList) ListRemove(key string, count int64, value any) (bool, error) {
	trace := linkTrace.TraceRedis("ListRemove", key, "")

	result, err := redisList.rdb.LRem(fs.Context, key, count, value).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisList *redisList) ListCount(key string) (int64, error) {
	trace := linkTrace.TraceRedis("ListCount", key, "")

	result, err := redisList.rdb.LLen(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisList *redisList) ListRange(key string, start int64, stop int64) ([]string, error) {
	trace := linkTrace.TraceRedis("ListRange", key, "")

	result, err := redisList.rdb.LRange(fs.Context, key, start, stop).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisList *redisList) ListLeftPop(timeout time.Duration, keys ...string) ([]string, error) {
	trace := linkTrace.TraceRedis("ListLeftPop", strings.Join(keys, ","), "")

	result, err := redisList.rdb.BLPop(fs.Context, timeout, keys...).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisList *redisList) ListRightPop(timeout time.Duration, keys ...string) ([]string, error) {
	trace := linkTrace.TraceRedis("ListRightPop", strings.Join(keys, ","), "")

	result, err := redisList.rdb.BRPop(fs.Context, timeout, keys...).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisList *redisList) ListRightPopPush(source, destination string, timeout time.Duration) (string, error) {
	trace := linkTrace.TraceRedis("ListRightPopPush", source, "")

	result, err := redisList.rdb.BRPopLPush(fs.Context, source, destination, timeout).Result()
	defer func() { trace.End(err) }()
	return result, err
}
