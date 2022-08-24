package redis

import (
	"fmt"
	"github.com/farseer-go/fs/configure"
	"testing"
	"time"
)

func Test_lockResult_Lock(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	local := client.Lock.GetLocker("key_local", time.Duration(1000))
	go control01(local)
	go control02(local)
	time.Sleep(time.Duration(100))
}

func control01(local lockResult) {
	for i := 0; i < 10; i++ {
		if local.TryLock() {
			fmt.Println("-----值：a", i)
			local.ReleaseLock()
		}
	}
}
func control02(local lockResult) {
	for i := 0; i < 10; i++ {
		if local.TryLock() {
			fmt.Println("-----值：b", i)
			local.ReleaseLock()
		}
	}
}
