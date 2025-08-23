package redis

import (
	"reflect"
	"time"

	"github.com/farseer-go/cache"
	"github.com/farseer-go/cache/eumExpiryType"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
)

// 二级缓存-本地缓存操作
type cacheInRedis struct {
	expiry      time.Duration      // 设置Redis缓存过期时间
	expiryType  eumExpiryType.Enum // 过期策略
	uniqueField string             // hash中的主键（唯一ID的字段名称）
	itemType    reflect.Type       // itemType
	key         string             // 缓存KEY
	lastVisitAt time.Time          // 最后一次访问时间
	redisClient IClient            // redis client
	isSetTtl    bool               // 用于记录，是否需要设置TTL
}

// 创建实例
func newCache(key string, uniqueField string, itemType reflect.Type, redisConfigName string, ops ...cache.Option) cache.ICache {
	op := &cache.Op{}
	for _, option := range ops {
		option(op)
	}

	r := &cacheInRedis{
		expiry:      op.Expiry,
		expiryType:  op.ExpiryType,
		uniqueField: uniqueField,
		itemType:    itemType,
		key:         key,
		redisClient: container.Resolve[IClient](redisConfigName),
		lastVisitAt: time.Now(),
		isSetTtl:    false,
	}

	if r.expiry > 0 {
		go r.updateTtl()
	}
	return r
}

func (r *cacheInRedis) Get() collections.ListAny {
	// 从redis hash中读取到slice
	lst, err := r.redisClient.HashToListAny(r.key, r.itemType)
	flog.ErrorIfExists(err)

	r.updateExpiry()
	r.trySetTTL()
	return lst
}

func (r *cacheInRedis) GetItem(cacheId any) any {
	// 动态创建实体
	entityPtr := reflect.New(r.itemType).Interface()

	// hash get
	exists, err := r.redisClient.HashToEntity(r.key, parse.ToString(cacheId), entityPtr)
	flog.ErrorIfExists(err)
	if !exists {
		return nil
	}

	r.updateExpiry()
	r.trySetTTL()

	return reflect.ValueOf(entityPtr).Elem().Interface()
}

func (r *cacheInRedis) GetItems(cacheIds []any) collections.ListAny {
	var keys []string
	for _, cacheId := range cacheIds {
		keys = append(keys, parse.ToString(cacheId))
	}

	items := collections.NewListAny()

	// 批量获取
	results, err := r.redisClient.HashGets(r.key, keys...)
	flog.ErrorIfExists(err)
	if err != nil {
		return items
	}

	// 转成数组
	for _, jsonContent := range results {
		entity := reflect.New(r.itemType).Interface()
		if err = snc.Unmarshal([]byte(jsonContent), entity); err == nil {
			items.Add(entity)
		}
	}

	r.updateExpiry()
	r.trySetTTL()
	return items
}

func (r *cacheInRedis) Set(val collections.ListAny) {
	if val.Count() == 0 {
		return
	}

	// 将ListAny转成map
	values := make(map[string]any)
	for _, item := range val.ToArray() {
		id := r.GetUniqueId(item)
		jsonContent, _ := snc.Marshal(item)
		values[id] = string(jsonContent)
	}

	err := r.redisClient.HashSet(r.key, values)
	flog.ErrorIfExists(err)

	// 设置缓存失效时间
	if r.expiry > 0 {
		r.updateExpiry()
		r.setTTL()
	}
}

func (r *cacheInRedis) SaveItem(newVal any) {
	newValDataKey := r.GetUniqueId(newVal)
	err := r.redisClient.HashSetEntity(r.key, newValDataKey, newVal)
	flog.ErrorIfExists(err)

	r.updateExpiry()
	r.trySetTTL()
}

func (r *cacheInRedis) Remove(cacheId any) {
	_, err := r.redisClient.HashDel(r.key, parse.Convert(cacheId, ""))
	flog.ErrorIfExists(err)
}

func (r *cacheInRedis) Clear() {
	_, err := r.redisClient.Del(r.key)
	flog.ErrorIfExists(err)
	r.isSetTtl = false
}

func (r *cacheInRedis) Count() int {
	r.updateExpiry()
	r.trySetTTL()
	return r.redisClient.HashCount(r.key)
}

func (r *cacheInRedis) ExistsItem(cacheId any) bool {
	exists, err := r.redisClient.HashExists(r.key, parse.Convert(cacheId, ""))
	flog.ErrorIfExists(err)

	r.updateExpiry()
	r.trySetTTL()
	return exists
}

func (r *cacheInRedis) ExistsKey() bool {
	exists, err := r.redisClient.Exists(r.key)
	flog.ErrorIfExists(err)

	r.updateExpiry()
	r.trySetTTL()
	return exists
}

// GetUniqueId 获取唯一字段数据
func (r *cacheInRedis) GetUniqueId(item any) string {
	val := reflect.ValueOf(item).FieldByName(r.uniqueField).Interface()
	return parse.Convert(val, "")
}

// 更新缓存过期时间
func (r *cacheInRedis) updateExpiry() {
	if r.expiry > 0 && r.expiryType == eumExpiryType.SlidingExpiration {
		r.lastVisitAt = time.Now()
	}
}

// 续期
func (r *cacheInRedis) updateTtl() {
	expiry := r.expiry
	if r.expiry >= 2*time.Second {
		expiry = r.expiry - time.Second
	} else if r.expiry >= time.Second {
		expiry = r.expiry - 500*time.Millisecond
	} else if r.expiry >= 500*time.Millisecond {
		expiry = r.expiry - 100*time.Millisecond
	}

	for {
		r.isSetTtl = false
		// 绝对时间的情况下，这时redis已自动过期了，所以恢复成未设置，在其它地方需要根据这个值来决定是否要重新设置TTL
		if r.expiryType == eumExpiryType.AbsoluteExpiration {
			// 如果设置了过期时间，需要检查下TTL
			// 如果过期时间比当前设置的过期时间还大，则需要重新设置
			if ttl, _ := r.redisClient.TTL(r.key); ttl == -1 || ttl > r.expiry {
				r.setTTL()
			}
		} else if r.expiryType == eumExpiryType.SlidingExpiration && !r.lastVisitAt.IsZero() {
			r.setTTL()
		}

		time.Sleep(expiry)
	}
}

func (r *cacheInRedis) setTTL() {
	r.isSetTtl = true
	r.lastVisitAt = time.Time{}
	_, _ = r.redisClient.SetTTL(r.key, r.expiry)
}

// 如果未设置ttl，则要设置
func (r *cacheInRedis) trySetTTL() {
	if r.expiry > 0 && r.expiryType == eumExpiryType.AbsoluteExpiration && !r.isSetTtl {
		r.setTTL()
	}
}
