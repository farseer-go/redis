package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisSet struct {
	rdb *redis.Client
}

func (redisSet *redisSet) SetAdd(key string, members ...any) (bool, error) {
	trace := linkTrace.TraceRedis("SetAdd", key, "")
	result, err := redisSet.rdb.SAdd(fs.Context, key, members...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisSet *redisSet) SetCount(key string) (int64, error) {
	trace := linkTrace.TraceRedis("SetCount", key, "")
	result, err := redisSet.rdb.SCard(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetRemove(key string, members ...any) (bool, error) {
	trace := linkTrace.TraceRedis("SetRemove", key, "")
	result, err := redisSet.rdb.SRem(fs.Context, key, members...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisSet *redisSet) SetGet(key string) ([]string, error) {
	trace := linkTrace.TraceRedis("SetGet", key, "")
	result, err := redisSet.rdb.SMembers(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetIsMember(key string, member any) (bool, error) {
	trace := linkTrace.TraceRedis("SetIsMember", key, "")
	result, err := redisSet.rdb.SIsMember(fs.Context, key, member).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetDiff(keys ...string) ([]string, error) {
	trace := linkTrace.TraceRedis("SetDiff", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SDiff(fs.Context, keys...).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetDiffStore(destination string, keys ...string) (bool, error) {
	trace := linkTrace.TraceRedis("SetDiffStore", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SDiffStore(fs.Context, destination, keys...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisSet *redisSet) SetInter(keys ...string) ([]string, error) {
	trace := linkTrace.TraceRedis("SetInter", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SInter(fs.Context, keys...).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetInterStore(destination string, keys ...string) (bool, error) {
	trace := linkTrace.TraceRedis("SetInterStore", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SInterStore(fs.Context, destination, keys...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisSet *redisSet) SetUnion(keys ...string) ([]string, error) {
	trace := linkTrace.TraceRedis("SetUnion", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SUnion(fs.Context, keys...).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisSet *redisSet) SetUnionStore(destination string, keys ...string) (bool, error) {
	trace := linkTrace.TraceRedis("SetUnionStore", strings.Join(keys, ","), "")
	result, err := redisSet.rdb.SUnionStore(fs.Context, destination, keys...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}
