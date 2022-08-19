package redis

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/memoryCache"
	"reflect"
)

// 二级缓存-本地缓存操作
type cacheInRedis struct {
}

func newCacheInRedis() memoryCache.ICache {
	return cacheInRedis{}
}

func (r cacheInRedis) Get(cacheKey memoryCache.CacheKey) collections.ListAny {
	// 动态创建切片
	arrType := reflect.SliceOf(cacheKey.ItemType)
	arr := reflect.MakeSlice(arrType, 0, 0).Interface()

	// 从redis hash中读取到slice
	redisClient := NewClient(cacheKey.RedisConfigName)
	_ = redisClient.Hash.ToArray(cacheKey.Key, &arr)

	// 切片转ListAny
	return mapper.ToListAny(arr)
}

func (r cacheInRedis) GetItem(cacheKey memoryCache.CacheKey, cacheId string) any {
	// 动态创建实体
	entity := reflect.New(cacheKey.ItemType).Elem().Interface()

	// hash get
	redisClient := NewClient(cacheKey.RedisConfigName)
	_ = redisClient.Hash.ToEntity(cacheKey.Key, cacheId, &entity)
	return entity
}

func (r cacheInRedis) Set(cacheKey memoryCache.CacheKey, val collections.ListAny) {
	// 将ListAny转成map
	values := make(map[string]any)
	for _, item := range val.ToArray() {
		id := cacheKey.GetUniqueId(item)
		jsonContent, _ := json.Marshal(item)
		values[id] = string(jsonContent)
	}

	redisClient := NewClient(cacheKey.RedisConfigName)
	_ = redisClient.Hash.Set(cacheKey.Key, values)
}

func (r cacheInRedis) SaveItem(cacheKey memoryCache.CacheKey, newVal any) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	_ = redisClient.Hash.SetEntity(cacheKey.Key, cacheKey.UniqueField, newVal)
}

func (r cacheInRedis) Remove(cacheKey memoryCache.CacheKey, cacheId string) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	_, _ = redisClient.Hash.Del(cacheKey.Key, cacheId)
}

func (r cacheInRedis) Clear(cacheKey memoryCache.CacheKey) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	_, _ = redisClient.Key.Del(cacheKey.Key)
}

func (r cacheInRedis) Count(cacheKey memoryCache.CacheKey) int {
	redisClient := NewClient(cacheKey.RedisConfigName)
	return redisClient.Hash.Count(cacheKey.Key)
}

func (r cacheInRedis) ExistsItem(cacheKey memoryCache.CacheKey, cacheId string) bool {
	redisClient := NewClient(cacheKey.RedisConfigName)
	exists, _ := redisClient.Hash.Exists(cacheKey.Key, cacheId)
	return exists
}

func (r cacheInRedis) ExistsKey(cacheKey memoryCache.CacheKey) bool {
	redisClient := NewClient(cacheKey.RedisConfigName)
	exists, _ := redisClient.Key.Exists(cacheKey.Key)
	return exists
}
