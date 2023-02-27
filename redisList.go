package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisList struct {
	rdb *redis.Client
}

func (redisList *redisList) ListPushRight(key string, values ...any) (bool, error) {
	result, err := redisList.rdb.RPush(fs.Context, key, values).Result()
	return result > 0, err
}

func (redisList *redisList) ListPushLeft(key string, values ...any) (bool, error) {
	result, err := redisList.rdb.LPush(fs.Context, key, values).Result()
	return result > 0, err
}

func (redisList *redisList) ListSet(key string, index int64, value any) (bool, error) {
	result, err := redisList.rdb.LSet(fs.Context, key, index, value).Result()
	if result == "OK" {
		return true, err
	}
	return false, err
}

func (redisList *redisList) ListRemove(key string, count int64, value any) (bool, error) {
	result, err := redisList.rdb.LRem(fs.Context, key, count, value).Result()
	return result > 0, err

}

func (redisList *redisList) ListCount(key string) (int64, error) {
	return redisList.rdb.LLen(fs.Context, key).Result()
}

func (redisList *redisList) ListRange(key string, start int64, stop int64) ([]string, error) {
	return redisList.rdb.LRange(fs.Context, key, start, stop).Result()
}

func (redisList *redisList) ListLeftPop(timeout time.Duration, keys ...string) ([]string, error) {
	return redisList.rdb.BLPop(fs.Context, timeout, keys...).Result()
}

func (redisList *redisList) ListRightPop(timeout time.Duration, keys ...string) ([]string, error) {
	return redisList.rdb.BRPop(fs.Context, timeout, keys...).Result()
}

func (redisList *redisList) ListRightPopPush(source, destination string, timeout time.Duration) (string, error) {
	return redisList.rdb.BRPopLPush(fs.Context, source, destination, timeout).Result()
}
