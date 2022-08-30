package redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//
//// String 测试
//func TestClientString(t *testing.T) {
//	client := NewClient("default")
//	err := client.String.Set("key1", "...3456")
//	if err == nil {
//		flog.Info("设置值:%v\n", "...3456")
//	}
//	get, _ := client.String.Get("key1")
//	flog.Info("获取值：%v\n", get)
//
//	//如果key值存在，设置这个会返回false
//	nx, _ := client.String.SetNX("key2", "1231", 100*time.Second)
//	flog.Info("设置过期时间：%v\n", nx)
//
//	get2, _ := client.String.Get("key2")
//	flog.Info("获取值：%v\n", get2)
//
//	ttl, _ := client.Key.TTL("key2")
//	flog.Info("获取过期时间：%v\n", ttl)
//}

func Test_redisString_Get(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key1")
	client.String.Set("key1", "9883")
	get, _ := client.String.Get("key1")
	assert.Equal(t, get, "9883")
}

func Test_redisString_Set(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key1")
	err := client.String.Set("key1", "...3456")
	assert.Equal(t, err, nil)
}

func Test_redisString_SetNX(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key1")
	nx, _ := client.String.SetNX("key1", "...3456", 100*time.Second)
	assert.Equal(t, nx, true)
}
