package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisSet struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

func (receiver *redisSet) SetAdd(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetAdd", key, "")
	result, err := receiver.rdb.SAdd(fs.Context, key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetCount(key string) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetCount", key, "")
	result, err := receiver.rdb.SCard(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetRemove(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetRemove", key, "")
	result, err := receiver.rdb.SRem(fs.Context, key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetGet(key string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetGet", key, "")
	result, err := receiver.rdb.SMembers(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetIsMember(key string, member any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetIsMember", key, "")
	result, err := receiver.rdb.SIsMember(fs.Context, key, member).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiff(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiff", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SDiff(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiffStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiffStore", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SDiffStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetInter(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInter", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SInter(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetInterStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInterStore", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SInterStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetUnion(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnion", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SUnion(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetUnionStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnionStore", strings.Join(keys, ","), "")
	result, err := receiver.rdb.SUnionStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
