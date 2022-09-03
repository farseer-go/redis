package redis

import (
	"encoding/json"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"reflect"
)

// 二级缓存-本地缓存操作
type cacheInRedis struct {
}

func (r cacheInRedis) Get(cacheKey cache.CacheKey) collections.ListAny {
	// 从redis hash中读取到slice
	redisClient := NewClient(cacheKey.RedisConfigName)
	lst, err := redisClient.Hash.ToListAny(cacheKey.Key, cacheKey.ItemType)
	if err != nil {
		flog.Error(err)
	}
	return lst
}

func (r cacheInRedis) GetItem(cacheKey cache.CacheKey, cacheId string) any {
	// 动态创建实体
	entityPtr := reflect.New(cacheKey.ItemType).Interface()

	// hash get
	redisClient := NewClient(cacheKey.RedisConfigName)
	err := redisClient.Hash.ToEntity(cacheKey.Key, cacheId, entityPtr)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil
		}
		flog.Error(err)
	}
	return reflect.ValueOf(entityPtr).Elem().Interface()
}

func (r cacheInRedis) Set(cacheKey cache.CacheKey, val collections.ListAny) {
	// 将ListAny转成map
	values := make(map[string]any)
	for _, item := range val.ToArray() {
		id := cacheKey.GetUniqueId(item)
		jsonContent, _ := json.Marshal(item)
		values[id] = string(jsonContent)
	}

	redisClient := NewClient(cacheKey.RedisConfigName)
	err := redisClient.Hash.Set(cacheKey.Key, values)
	if err != nil {
		flog.Error(err)
	}
	if cacheKey.RedisExpiry > 0 {
		_, _ = redisClient.Key.SetTTL(cacheKey.Key, cacheKey.RedisExpiry)
	}
}

func (r cacheInRedis) SaveItem(cacheKey cache.CacheKey, newVal any) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	newValDataKey := cacheKey.GetUniqueId(newVal)
	err := redisClient.Hash.SetEntity(cacheKey.Key, newValDataKey, newVal)
	if err != nil {
		flog.Error(err)
	}
}

func (r cacheInRedis) Remove(cacheKey cache.CacheKey, cacheId string) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	_, err := redisClient.Hash.Del(cacheKey.Key, cacheId)
	if err != nil {
		flog.Error(err)
	}
}

func (r cacheInRedis) Clear(cacheKey cache.CacheKey) {
	redisClient := NewClient(cacheKey.RedisConfigName)
	_, err := redisClient.Key.Del(cacheKey.Key)
	if err != nil {
		flog.Error(err)
	}
}

func (r cacheInRedis) Count(cacheKey cache.CacheKey) int {
	redisClient := NewClient(cacheKey.RedisConfigName)
	return redisClient.Hash.Count(cacheKey.Key)
}

func (r cacheInRedis) ExistsItem(cacheKey cache.CacheKey, cacheId string) bool {
	redisClient := NewClient(cacheKey.RedisConfigName)
	exists, err := redisClient.Hash.Exists(cacheKey.Key, cacheId)
	if err != nil {
		flog.Error(err)
	}
	return exists
}

func (r cacheInRedis) ExistsKey(cacheKey cache.CacheKey) bool {
	redisClient := NewClient(cacheKey.RedisConfigName)
	exists, err := redisClient.Key.Exists(cacheKey.Key)
	if err != nil {
		flog.Error(err)
	}
	return exists
}
