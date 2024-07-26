package test

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redisKey_Del(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	del, _ := client.Del("key_client")
	assert.Equal(t, del, true)
}

func Test_redisKey_Exists(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	exist, _ := client.Exists("key_client")
	assert.Equal(t, exist, true)
}

func Test_redisKey_TTL(t *testing.T) {
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	ttl, _ := client.TTL("key_client")
	assert.Equal(t, ttl, time.Duration(-1))

}
