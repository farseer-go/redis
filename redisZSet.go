package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
)

type redisZSet struct {
	*redisManager
}

type redisZ struct {
	Score  float64
	Member any
}
type redisZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func (receiver *redisZSet) ZSetAdd(key string, members ...*redisZ) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetAdd", key, "")
	var redisZZ []*redis.Z
	for _, member := range members {
		redisZZ = append(redisZZ, &redis.Z{Score: member.Score, Member: member.Member})
	}
	result, err := receiver.GetClient().ZAdd(fs.Context, key, redisZZ...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisZSet) ZSetScore(key string, member string) (float64, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetScore", key, "")
	result, err := receiver.GetClient().ZScore(fs.Context, key, member).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisZSet) ZSetRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRange", key, "")
	result, err := receiver.GetClient().ZRange(fs.Context, key, start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisZSet) ZSetRevRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRevRange", key, "")
	result, err := receiver.GetClient().ZRevRange(fs.Context, key, start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisZSet) ZSetRangeByScore(key string, opt *redisZRangeBy) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRangeByScore", key, "")
	rby := redis.ZRangeBy{Min: opt.Min, Max: opt.Max, Offset: opt.Offset, Count: opt.Count}
	result, err := receiver.GetClient().ZRangeByScore(fs.Context, key, &rby).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}
