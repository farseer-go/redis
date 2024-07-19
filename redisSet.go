package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisSet struct {
	*redisManager
}

func (receiver *redisSet) SetAdd(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetAdd", key, "")
	result, err := receiver.GetClient().SAdd(context.Background(), key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetCount(key string) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetCount", key, "")
	result, err := receiver.GetClient().SCard(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetRemove(key string, members ...any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetRemove", key, "")
	result, err := receiver.GetClient().SRem(context.Background(), key, members...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetGet(key string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetGet", key, "")
	result, err := receiver.GetClient().SMembers(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetIsMember(key string, member any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetIsMember", key, "")
	result, err := receiver.GetClient().SIsMember(context.Background(), key, member).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiff(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiff", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SDiff(context.Background(), keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetDiffStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetDiffStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SDiffStore(context.Background(), destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetInter(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInter", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SInter(context.Background(), keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetInterStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetInterStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SInterStore(context.Background(), destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisSet) SetUnion(keys ...string) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnion", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SUnion(context.Background(), keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisSet) SetUnionStore(destination string, keys ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("SetUnionStore", strings.Join(keys, ","), "")
	result, err := receiver.GetClient().SUnionStore(context.Background(), destination, keys...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}
