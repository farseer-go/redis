package redis

import (
	"context"

	"github.com/farseer-go/fs/parse"
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
	result, err := receiver.GetClient().ZAdd(context.Background(), key, redisZZ...).Result()
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(parse.ToInt(result))
	}
	return result > 0, err
}

func (receiver *redisZSet) ZSetScore(key string, member string) (float64, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetScore", key, "")
	result, err := receiver.GetClient().ZScore(context.Background(), key, member).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(1)
	}
	return result, err
}

func (receiver *redisZSet) ZSetRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRange", key, "")
	result, err := receiver.GetClient().ZRange(context.Background(), key, start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}

func (receiver *redisZSet) ZSetRevRange(key string, start int64, stop int64) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRevRange", key, "")
	result, err := receiver.GetClient().ZRevRange(context.Background(), key, start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}

func (receiver *redisZSet) ZSetRangeByScore(key string, opt *redisZRangeBy) ([]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("ZSetRangeByScore", key, "")
	rby := redis.ZRangeBy{Min: opt.Min, Max: opt.Max, Offset: opt.Offset, Count: opt.Count}
	result, err := receiver.GetClient().ZRangeByScore(context.Background(), key, &rby).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()

	if err == nil {
		traceDetail.SetRows(len(result))
	}
	return result, err
}
