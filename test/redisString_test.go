package test

import (
	"testing"
	"time"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
)

//
//// String 测试
//func TestClientString(t *testing.T) {
//	client := newClient("default")
//	err := client.StringSet("key1", "...3456")
//	if err == nil {
//		flog.Info("设置值:%v\n", "...3456")
//	}
//	get, _ := client.StringGet("key1")
//	flog.Info("获取值：%v\n", get)
//
//	//如果key值存在，设置这个会返回false
//	nx, _ := client.StringSetNX("key2", "1231", 100*time.Second)
//	flog.Info("设置过期时间：%v\n", nx)
//
//	get2, _ := client.StringGet("key2")
//	flog.Info("获取值：%v\n", get2)
//
//	ttl, _ := client.TTL("key2")
//	flog.Info("获取过期时间：%v\n", ttl)
//}

func Test_redisString_Get(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	client.StringSet("key1", "9883")
	get, _ := client.StringGet("key1")
	assert.Equal(t, get, "9883")
}

func Test_redisString_Set(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	err := client.StringSet("key1", "...3456")
	assert.Equal(t, err, nil)
}

func Test_redisString_SetNX(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	nx, _ := client.StringSetNX("key1", "...3456", 100*time.Second)
	assert.Equal(t, nx, true)
}

func Test_redis_StringIncrInt(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	vl, _ := client.StringIncrInt("key1", 1)
	assert.Equal(t, vl, int(1))
	vl, _ = client.StringIncrInt("key1", 2)
	assert.Equal(t, vl, int(3))
}
func Test_redis_StringIncrInt64(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	vl, _ := client.StringIncrInt64("key1", 1)
	assert.Equal(t, vl, int64(1))
	vl, _ = client.StringIncrInt64("key1", 2)
	assert.Equal(t, vl, int64(3))
}
func Test_redis_StringIncrFloat32(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	vl, _ := client.StringIncrFloat32("key1", 1.2)
	assert.Equal(t, vl, float32(1.2))
	vl, _ = client.StringIncrFloat32("key1", 2.3)
	assert.Equal(t, vl, float32(3.5))
}
func Test_redis_StringIncrFloat64(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key1")
	vl, _ := client.StringIncrFloat64("key1", 1.2)
	assert.Equal(t, vl, float64(1.2))
	vl, _ = client.StringIncrFloat64("key1", 2.3)
	assert.Equal(t, vl, float64(3.5))
}
