package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_lockResult_Lock(t *testing.T) {
	//
	configure.SetDefault("Redis.default", "Server=192.168.1.8:6379,DB=15,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	for i := 0; i < 100; i++ {
		control01(t, client)
	}
}

func control01(t *testing.T, client redis.IClient) {
	sw := stopwatch.StartNew()
	local := client.LockNew("key_local", "", 2*time.Second)
	defer local.ReleaseLock()
	defer func() {
		flog.Infof("LockNew，耗时：%s", sw.GetMillisecondsText())
	}()
	assert.True(t, local.TryLock())
	assert.False(t, local.TryLock())
}
