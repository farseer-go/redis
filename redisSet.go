package redis

import "github.com/go-redis/redis/v8"

type redisSet struct {
	rdb *redis.Client
}

func (redisSet *redisSet) SetAdd(key string, members ...any) (bool, error) {
	result, err := redisSet.rdb.SAdd(ctx, key, members...).Result()
	return result > 0, err
}

func (redisSet *redisSet) SetCount(key string) (int64, error) {
	return redisSet.rdb.SCard(ctx, key).Result()
}

func (redisSet *redisSet) SetRemove(key string, members ...any) (bool, error) {
	result, err := redisSet.rdb.SRem(ctx, key, members...).Result()
	return result > 0, err
}

func (redisSet *redisSet) SetGet(key string) ([]string, error) {
	return redisSet.rdb.SMembers(ctx, key).Result()
}

func (redisSet *redisSet) SetIsMember(key string, member any) (bool, error) {
	return redisSet.rdb.SIsMember(ctx, key, member).Result()
}

func (redisSet *redisSet) SetDiff(keys ...string) ([]string, error) {
	return redisSet.rdb.SDiff(ctx, keys...).Result()
}

func (redisSet *redisSet) SetDiffStore(destination string, keys ...string) (bool, error) {
	result, err := redisSet.rdb.SDiffStore(ctx, destination, keys...).Result()
	return result > 0, err
}

func (redisSet *redisSet) SetInter(keys ...string) ([]string, error) {
	return redisSet.rdb.SInter(ctx, keys...).Result()
}

func (redisSet *redisSet) SetInterStore(destination string, keys ...string) (bool, error) {
	result, err := redisSet.rdb.SInterStore(ctx, destination, keys...).Result()
	return result > 0, err
}

func (redisSet *redisSet) SetUnion(keys ...string) ([]string, error) {
	return redisSet.rdb.SUnion(ctx, keys...).Result()
}

func (redisSet *redisSet) SetUnionStore(destination string, keys ...string) (bool, error) {
	result, err := redisSet.rdb.SUnionStore(ctx, destination, keys...).Result()
	return result > 0, err
}
