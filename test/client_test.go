package test

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 创建客户端 测试
func TestNewClient(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := redis.NewClient("default")
	defer client.Key.Del("key_client")
	client.String.Set("key_client", "...12312")
	exists, _ := client.Key.Exists("key_client")
	assert.Equal(t, exists, true)
}
