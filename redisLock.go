package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"time"
)

// 分布式锁
type redisLock struct {
	rdb *redis.Client
}

type lockResult struct {
	key          string // 锁名称
	val          string // 锁值
	expiration   time.Duration
	rdb          *redis.Client
	traceManager trace.IManager
}

// LockNew 创建锁
func (r redisLock) LockNew(key, val string, expiration time.Duration) core.ILock {
	return &lockResult{
		rdb:          r.rdb,
		key:          key,
		val:          val,
		expiration:   expiration,
		traceManager: container.Resolve[trace.IManager](),
	}
}

// TryLock 尝试加锁
func (receiver *lockResult) TryLock() bool {
	traceDetail := receiver.traceManager.TraceRedis("TryLock", receiver.key, "")

	result, err := receiver.rdb.SetNX(fs.Context, receiver.key, receiver.val, receiver.expiration).Result()
	defer func() { traceDetail.End(err) }()
	if err != nil {
		_ = flog.Errorf("redis加锁异常：%s", err.Error())
	}
	return result
}

// TryLockRun 尝试加锁，执行完后，自动释放锁
func (receiver *lockResult) TryLockRun(fn func()) bool {
	traceDetail := receiver.traceManager.TraceRedis("TryLockRun", receiver.key, "")

	sw := stopwatch.StartNew()
	result, err := receiver.rdb.SetNX(fs.Context, receiver.key, receiver.val, receiver.expiration).Result()
	defer func() { traceDetail.End(err) }()
	flog.Debugf("获取Redis锁，耗时：%s", sw.GetMicrosecondsText())
	if err != nil {
		_ = flog.Errorf("redis加锁异常：%s", err.Error())
	}
	if result {
		defer receiver.ReleaseLock()
		fn()
	}
	return result
}

// GetLock 获取锁，直到获取成功
func (receiver *lockResult) GetLock() {
	traceDetail := receiver.traceManager.TraceRedis("GetLock", receiver.key, "")
	var err error
	defer func() { traceDetail.End(err) }()

	for {
		var result bool
		result, err = receiver.rdb.SetNX(fs.Context, receiver.key, receiver.val, receiver.expiration).Result()

		if err != nil {
			_ = flog.Errorf("redis加锁异常：%s", err.Error())
		}
		if result {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// GetLockRun 获取锁，直到获取成功，执行完后，自动释放锁
func (receiver *lockResult) GetLockRun(fn func()) {
	for {
		if receiver.TryLockRun(fn) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// ReleaseLock 锁放锁
func (receiver *lockResult) ReleaseLock() {
	traceDetail := receiver.traceManager.TraceRedis("ReleaseLock", receiver.key, "")
	_, err := receiver.rdb.Del(fs.Context, receiver.key).Result()
	defer func() { traceDetail.End(err) }()

}
