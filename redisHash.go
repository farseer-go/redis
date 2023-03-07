package redis

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"github.com/go-redis/redis/v8"
	"reflect"
)

type redisHash struct {
	rdb *redis.Client
}

func (redisHash *redisHash) HashSetEntity(key string, field string, entity any) error {
	jsonContent, _ := json.Marshal(entity)
	return redisHash.rdb.HSet(fs.Context, key, field, string(jsonContent)).Err()
}

func (redisHash *redisHash) HashSet(key string, fieldValues ...any) error {
	return redisHash.rdb.HSet(fs.Context, key, fieldValues...).Err()
}

func (redisHash *redisHash) HashGet(key string, field string) (string, error) {
	return redisHash.rdb.HGet(fs.Context, key, field).Result()
}

func (redisHash *redisHash) HashGetAll(key string) (map[string]string, error) {
	return redisHash.rdb.HGetAll(fs.Context, key).Result()
}

func (redisHash *redisHash) HashToEntity(key string, field string, entity any) (bool, error) {
	jsonContent, err := redisHash.rdb.HGet(fs.Context, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	// 反序列
	return true, json.Unmarshal([]byte(jsonContent), entity)
}

func (redisHash *redisHash) HashToArray(key string, arrSlice any) error {
	arrVal := reflect.ValueOf(arrSlice).Elem()
	arrType, isSlice := types.IsSlice(arrVal)
	if !isSlice {
		panic("arr入参必须为切片类型")
	}

	result, err := redisHash.rdb.HGetAll(fs.Context, key).Result()
	if err != nil {
		return flog.Error(err)
	}

	lst := collections.NewListAny()
	for _, vJson := range result {
		item := reflect.New(arrType.Elem()).Interface()
		_ = json.Unmarshal([]byte(vJson), item)
		lst.Add(reflect.ValueOf(item).Elem().Interface())
	}

	lst.MapToArray(arrSlice)
	return nil
}

func (redisHash *redisHash) HashToListAny(key string, itemType reflect.Type) (collections.ListAny, error) {
	lst := collections.NewListAny()
	result, err := redisHash.rdb.HGetAll(fs.Context, key).Result()
	if err != nil {
		_ = flog.Error(err)
		return lst, err
	}
	for _, vJson := range result {
		item := reflect.New(itemType).Interface()
		_ = json.Unmarshal([]byte(vJson), item)
		lst.Add(reflect.ValueOf(item).Elem().Interface())
	}
	return lst, nil
}

func (redisHash *redisHash) HashExists(key string, field string) (bool, error) {
	return redisHash.rdb.HExists(fs.Context, key, field).Result()
}

func (redisHash *redisHash) HashDel(key string, fields ...string) (bool, error) {
	result, err := redisHash.rdb.HDel(fs.Context, key, fields...).Result()
	return result > 0, err
}

func (redisHash *redisHash) HashCount(key string) int {
	hLen := redisHash.rdb.HLen(fs.Context, key)
	count, _ := hLen.Uint64()
	return int(count)
}

func (redisHash *redisHash) HashIncrInt(key string, field string, value int) (int, error) {
	val, err := redisHash.rdb.HIncrBy(fs.Context, key, field, parse.Convert(value, int64(value))).Result()
	return parse.Convert(val, 0), err
}

func (redisHash *redisHash) HashIncrInt64(key string, field string, value int64) (int64, error) {
	return redisHash.rdb.HIncrBy(fs.Context, key, field, value).Result()
}

func (redisHash *redisHash) HashIncrFloat32(key string, field string, value float32) (float32, error) {
	val, err := redisHash.rdb.HIncrByFloat(fs.Context, key, field, float64(value)).Result()
	return parse.Convert(val, float32(0)), err
}

func (redisHash *redisHash) HashIncrFloat64(key string, field string, value float64) (float64, error) {
	return redisHash.rdb.HIncrByFloat(fs.Context, key, field, value).Result()
}
