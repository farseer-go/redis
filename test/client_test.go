package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 创建客户端 测试
func TestNewClient(t *testing.T) {
	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.StringSet("key_client", "...12312")
	exists, _ := client.Exists("key_client")
	assert.Equal(t, exists, true)
}
