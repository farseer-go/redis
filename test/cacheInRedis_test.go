package test

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type po struct {
	Name string
	Age  int
}

func init() {
	fs.Initialize[redis.Module]("unit test")
}

func TestCacheInRedis_Set(t *testing.T) {
	assert.Panics(t, func() {
		redis.SetProfiles[po]("TestCacheInRedis_Set", "ClientName", "default")
	})

	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_Set", "Name", "default")
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
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_GetItem", "Name", "default")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 18)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInRedis_GetItems(t *testing.T) {
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_GetItem", "Name", "default")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	lst := cacheManage.GetItems("steden", "steden2")
	lst = lst.OrderBy(func(item po) any {
		return item.Age
	}).ToList()

	assert.Equal(t, lst.Index(0).Name, "steden")
	assert.Equal(t, lst.Index(0).Age, 18)

	assert.Equal(t, lst.Index(1).Name, "steden2")
	assert.Equal(t, lst.Index(1).Age, 19)
}

func TestCacheInRedis_SaveItem(t *testing.T) {
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_SaveItem", "Name", "default")
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
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_Remove", "Name", "default")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Remove("steden")

	_, exists := cacheManage.GetItem("steden")
	assert.False(t, exists)

	item2, _ := cacheManage.GetItem("steden2")
	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInRedis_Clear(t *testing.T) {
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_Clear", "Name", "default")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.Equal(t, cacheManage.Count(), 2)
	cacheManage.Clear()
	assert.Equal(t, cacheManage.Count(), 0)
}

func TestCacheInRedis_Exists(t *testing.T) {
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_Exists", "Name", "default", func(op *cache.Op) {
		op.AbsoluteExpiration(10 * time.Millisecond)
	})
	assert.False(t, cacheManage.ExistsKey())
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.True(t, cacheManage.ExistsKey())
}

func TestCacheInRedis_Ttl(t *testing.T) {
	cacheManage := redis.SetProfiles[po]("TestCacheInRedis_Ttl", "Name", "default", func(op *cache.Op) {
		op.AbsoluteExpiration(1 * time.Second)
	})
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
