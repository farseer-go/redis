package redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redisKey_Del(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	del, _ := client.Key.Del("key_client")
	assert.Equal(t, del, true)
}

func Test_redisKey_Exists(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	exist, _ := client.Key.Exists("key_client")
	assert.Equal(t, exist, true)
}

func Test_redisKey_TTL(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	ttl, _ := client.Key.TTL("key_client")
	assert.Equal(t, ttl, time.Duration(-1))

}
