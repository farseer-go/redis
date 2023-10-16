package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"time"
)

// 分布式锁
type redisLock struct {
	rdb *redis.Client
}

type lockResult struct {
	key        string // 锁名称
	val        string // 锁值
	expiration time.Duration
	rdb        *redis.Client
}

// LockNew 创建锁
func (r redisLock) LockNew(key, val string, expiration time.Duration) core.ILock {
	return &lockResult{
		rdb:        r.rdb,
		key:        key,
		val:        val,
		expiration: expiration,
	}
}

// TryLock 尝试加锁
func (r *lockResult) TryLock() bool {
	trace := linkTrace.TraceRedis("TryLock", r.key, "")

	result, err := r.rdb.SetNX(fs.Context, r.key, r.val, r.expiration).Result()
	defer func() { trace.End(err) }()
	if err != nil {
		_ = flog.Errorf("redis加锁异常：%s", err.Error())
	}
	return result
}

// TryLockRun 尝试加锁，执行完后，自动释放锁
func (r *lockResult) TryLockRun(fn func()) bool {
	trace := linkTrace.TraceRedis("TryLockRun", r.key, "")

	sw := stopwatch.StartNew()
	result, err := r.rdb.SetNX(fs.Context, r.key, r.val, r.expiration).Result()
	defer func() { trace.End(err) }()
	flog.Debugf("获取Redis锁，耗时：%s", sw.GetMicrosecondsText())
	if err != nil {
		_ = flog.Errorf("redis加锁异常：%s", err.Error())
	}
	if result {
		defer r.ReleaseLock()
		fn()
	}
	return result
}

// GetLock 获取锁，直到获取成功
func (r *lockResult) GetLock() {
	trace := linkTrace.TraceRedis("GetLock", r.key, "")
	var err error
	defer func() { trace.End(err) }()

	for {
		var result bool
		result, err = r.rdb.SetNX(fs.Context, r.key, r.val, r.expiration).Result()

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
func (r *lockResult) GetLockRun(fn func()) {
	for {
		if r.TryLockRun(fn) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// ReleaseLock 锁放锁
func (r *lockResult) ReleaseLock() {
	trace := linkTrace.TraceRedis("ReleaseLock", r.key, "")
	_, err := r.rdb.Del(fs.Context, r.key).Result()
	defer func() { trace.End(err) }()

}
