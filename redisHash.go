package redis

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/types"
	"github.com/go-redis/redis/v8"
	"reflect"
)

type redisHash struct {
	rdb *redis.Client
}

// SetEntity 添加并序列化成json
func (redisHash *redisHash) SetEntity(key string, field string, entity any) error {
	jsonContent, _ := json.Marshal(entity)
	return redisHash.rdb.HSet(ctx, key, field, string(jsonContent)).Err()
}

// Set 添加
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (redisHash *redisHash) Set(key string, values ...interface{}) error {
	return redisHash.rdb.HSet(ctx, key, values...).Err()
}

// Get 获取
func (redisHash *redisHash) Get(key string, field string) (string, error) {
	return redisHash.rdb.HGet(ctx, key, field).Result()
}

// ToEntity 获取单个对象
//
//	var client DomainObject
//	_ = repository.Client.Hash.ToEntity("redisKey", "field", &client)
func (redisHash *redisHash) ToEntity(key string, field string, entity any) error {
	jsonContent, err := redisHash.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}
	// 反序列
	return json.Unmarshal([]byte(jsonContent), entity)
}

// GetAll 获取所有集合数据
func (redisHash *redisHash) GetAll(key string) (map[string]string, error) {
	return redisHash.rdb.HGetAll(ctx, key).Result()
}

// ToArray 将hash.value反序列化成切片对象
//
//	type record struct {
//		ClientId int `json:"id"`
//	}
//	var records []record
//	ToArray("test", &records)
func (redisHash *redisHash) ToArray(key string, arrSlice any) error {
	arrVal := reflect.ValueOf(arrSlice).Elem()
	arrType, isSlice := types.IsSlice(arrVal)
	if !isSlice {
		panic("arr入参必须为切片类型")
	}

	result, err := redisHash.rdb.HGetAll(ctx, key).Result()
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

// ToListAny 将hash的数据转成collections.ListAny
func (redisHash *redisHash) ToListAny(key string, itemType reflect.Type) (collections.ListAny, error) {
	lst := collections.NewListAny()
	result, err := redisHash.rdb.HGetAll(ctx, key).Result()
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

// Exists 成员是否存在
func (redisHash *redisHash) Exists(key string, field string) (bool, error) {
	return redisHash.rdb.HExists(ctx, key, field).Result()
}

// Del 移除指定成员
func (redisHash *redisHash) Del(key string, fields ...string) (bool, error) {
	result, err := redisHash.rdb.HDel(ctx, key, fields...).Result()
	return result > 0, err
}

// Count 获取hash的数量
func (redisHash *redisHash) Count(key string) int {
	hLen := redisHash.rdb.HLen(ctx, key)
	count, _ := hLen.Uint64()
	return int(count)
}
