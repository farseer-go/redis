package redis

import (
	"encoding/json"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"reflect"
	"time"
)

// 二级缓存-本地缓存操作
type cacheInRedis struct {
	expiry          time.Duration // 设置Memory缓存过期时间
	uniqueField     string        // hash中的主键（唯一ID的字段名称）
	itemType        reflect.Type  // itemType
	key             string        // 缓存KEY
	redisConfigName string
}

// 创建实例
func newCache(key string, uniqueField string, itemType reflect.Type, expiry time.Duration, redisConfigName string) cache.ICache {
	return &cacheInRedis{
		expiry:          expiry,
		uniqueField:     uniqueField,
		itemType:        itemType,
		key:             key,
		redisConfigName: redisConfigName,
	}
}

func (r *cacheInRedis) Get() collections.ListAny {
	// 从redis hash中读取到slice
	redisClient := NewClient(r.redisConfigName)
	lst, err := redisClient.Hash.ToListAny(r.key, r.itemType)
	if err != nil {
		_ = flog.Error(err)
	}
	return lst
}

func (r *cacheInRedis) GetItem(cacheId string) any {
	// 动态创建实体
	entityPtr := reflect.New(r.itemType).Interface()

	// hash get
	redisClient := NewClient(r.redisConfigName)
	exists, err := redisClient.Hash.ToEntity(r.key, cacheId, entityPtr)
	if err != nil {
		_ = flog.Error(err)
	}
	if !exists {
		return nil
	}
	return reflect.ValueOf(entityPtr).Elem().Interface()
}

func (r *cacheInRedis) Set(val collections.ListAny) {
	if val.Count() == 0 {
		return
	}

	// 将ListAny转成map
	values := make(map[string]any)
	for _, item := range val.ToArray() {
		id := r.GetUniqueId(item)
		jsonContent, _ := json.Marshal(item)
		values[id] = string(jsonContent)
	}

	redisClient := NewClient(r.redisConfigName)
	err := redisClient.Hash.Set(r.key, values)
	if err != nil {
		_ = flog.Error(err)
	}
	if r.expiry > 0 {
		_, _ = redisClient.Key.SetTTL(r.key, r.expiry)
	}
}

func (r *cacheInRedis) SaveItem(newVal any) {
	redisClient := NewClient(r.redisConfigName)
	newValDataKey := r.GetUniqueId(newVal)
	err := redisClient.Hash.SetEntity(r.key, newValDataKey, newVal)
	if err != nil {
		_ = flog.Error(err)
	}
}

func (r *cacheInRedis) Remove(cacheId string) {
	redisClient := NewClient(r.redisConfigName)
	_, err := redisClient.Hash.Del(r.key, cacheId)
	if err != nil {
		_ = flog.Error(err)
	}
}

func (r *cacheInRedis) Clear() {
	redisClient := NewClient(r.redisConfigName)
	_, err := redisClient.Key.Del(r.key)
	if err != nil {
		_ = flog.Error(err)
	}
}

func (r *cacheInRedis) Count() int {
	redisClient := NewClient(r.redisConfigName)
	return redisClient.Hash.Count(r.key)
}

func (r *cacheInRedis) ExistsItem(cacheId string) bool {
	redisClient := NewClient(r.redisConfigName)
	exists, err := redisClient.Hash.Exists(r.key, cacheId)
	if err != nil {
		_ = flog.Error(err)
	}
	return exists
}

func (r *cacheInRedis) ExistsKey() bool {
	redisClient := NewClient(r.redisConfigName)
	exists, err := redisClient.Key.Exists(r.key)
	if err != nil {
		_ = flog.Error(err)
	}
	return exists
}

// GetUniqueId 获取唯一字段数据
func (r *cacheInRedis) GetUniqueId(item any) string {
	val := reflect.ValueOf(item).FieldByName(r.uniqueField).Interface()
	return parse.Convert(val, "")
}
