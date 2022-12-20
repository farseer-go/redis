package test

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type po struct {
	Name string
	Age  int
}

func TestCacheInRedis_Set(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	assert.Panics(t, func() {
		cache.SetProfilesInRedis[po]("test", "default", "ClientName", 0)
	})

	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Set(lst.ToArray()...)

	lst2 := cacheManage.Get()

	assert.Equal(t, lst.Count(), lst2.Count())

	for i := 0; i < lst.Count(); i++ {
		lst.Index(i)
		item := lst2.Where(func(item po) bool {
			return item.Name == lst.Index(i).Name
		}).First()

		assert.Equal(t, lst.Index(i).Name, item.Name)
		assert.Equal(t, lst.Index(i).Age, item.Age)
	}
}

func TestCacheInRedis_GetItem(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 18)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInRedis_SaveItem(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.SaveItem(po{Name: "steden", Age: 99})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 99)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInRedis_Remove(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Remove("steden")

	_, exists := cacheManage.GetItem("steden")
	assert.False(t, exists)

	item2, _ := cacheManage.GetItem("steden2")
	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInRedis_Clear(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.Equal(t, cacheManage.Count(), 2)
	cacheManage.Clear()
	assert.Equal(t, cacheManage.Count(), 0)
}

func TestCacheInRedis_Exists(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	assert.False(t, cacheManage.ExistsKey())
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.True(t, cacheManage.ExistsKey())
}

func TestCacheInRedis_Ttl(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	configure.SetDefault("redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	cache.SetProfilesInRedis[po]("test", "default", "Name", 1*time.Second)
	cacheManage := cache.GetCacheManage[po]("test")
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Set(lst.ToArray()...)
	lst2 := cacheManage.Get()
	assert.Equal(t, lst.Count(), lst2.Count())
	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lst2.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lst2.Index(i).Age)
	}
	time.Sleep(1 * time.Second)
	assert.False(t, cacheManage.ExistsKey())
}
