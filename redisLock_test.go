package redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
	"testing"
	"time"
)

func Test_lockResult_Lock(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	local := client.Lock.GetLocker("key_local", time.Duration(1000))
	control01(local)
	go control02(local)
	time.Sleep(time.Duration(100))
}

func control01(local lockResult) {
	if !local.TryLock() {
		flog.Info("-----01加锁失败")
	}
	defer local.ReleaseLock()
	for i := 0; i < 10; i++ {
		flog.Info("-----值：a", i)
	}
}
func control02(local lockResult) {
	if !local.TryLock() {
		flog.Info("-----02加锁失败")
	}
	defer local.ReleaseLock()
	for i := 0; i < 10; i++ {
		flog.Info("-----值：b", i)
	}
}