package redis

import (
	"github.com/farseer-go/fs"
	"github.com/go-redis/redis/v8"
)

type redisZSet struct {
	rdb *redis.Client
}

type redisZ struct {
	Score  float64
	Member any
}
type redisZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func (redisZSet *redisZSet) ZSetAdd(key string, members ...*redisZ) (bool, error) {
	var redisZZ []*redis.Z
	for _, member := range members {
		redisZZ = append(redisZZ, &redis.Z{Score: member.Score, Member: member.Member})
	}
	result, err := redisZSet.rdb.ZAdd(fs.Context, key, redisZZ...).Result()
	return result > 0, err
}

func (redisZSet *redisZSet) ZSetScore(key string, member string) (float64, error) {
	return redisZSet.rdb.ZScore(fs.Context, key, member).Result()
}

func (redisZSet *redisZSet) ZSetRange(key string, start int64, stop int64) ([]string, error) {
	return redisZSet.rdb.ZRange(fs.Context, key, start, stop).Result()
}

func (redisZSet *redisZSet) ZSetRevRange(key string, start int64, stop int64) ([]string, error) {
	return redisZSet.rdb.ZRevRange(fs.Context, key, start, stop).Result()
}

func (redisZSet *redisZSet) ZSetRangeByScore(key string, opt *redisZRangeBy) ([]string, error) {
	rby := redis.ZRangeBy{Min: opt.Min, Max: opt.Max, Offset: opt.Offset, Count: opt.Count}
	return redisZSet.rdb.ZRangeByScore(fs.Context, key, &rby).Result()
}
