package redis

import (
	"github.com/farseer-go/fs"
	"strings"
)

type redisSet struct {
	*redisManager
}

func (receiver *redisSet) SetAdd(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetAdd", key, "")
	result, err := receiver.GetClient().SAdd(fs.Context, key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetCount(key string) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetCount", key, "")
	result, err := receiver.GetClient().SCard(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetRemove(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetRemove", key, "")
	result, err := receiver.GetClient().SRem(fs.Context, key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetGet(key string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetGet", key, "")
	result, err := receiver.GetClient().SMembers(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetIsMember(key string, member any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetIsMember", key, "")
	result, err := receiver.GetClient().SIsMember(fs.Context, key, member).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiff(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiff", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SDiff(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiffStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiffStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SDiffStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetInter(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInter", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SInter(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetInterStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInterStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SInterStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetUnion(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnion", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SUnion(fs.Context, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetUnionStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnionStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SUnionStore(fs.Context, destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
