package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"github.com/go-redis/redis/v8"
	"time"
)

// 选举
type redisElection struct {
	rdb          *redis.Client
	traceManager trace.IManager
}

// Election 选举成功后，会自动续约。
// 未拿到master，会持续尝试获取master
func (receiver *redisElection) Election(key string, fn func()) {
	for {
		cmd := receiver.rdb.SetNX(fs.Context, key, fs.AppId, 20*time.Second)
		result, _ := cmd.Result()
		// 拿到锁了
		if result {
			// 给锁续租约
			go receiver.leaseRenewal(key)
			fn()
			<-fs.Context.Done()
			return
		}

		// 没有拿到master的节点，需获取当前租约剩余时间，到期后，尝试获取
		duration, _ := receiver.rdb.TTL(fs.Context, key).Result()
		<-time.After(duration)
	}
}

// GetLeaderId 获取当前LeaderId
func (receiver *redisElection) GetLeaderId(key string) int64 {
	traceDetail := receiver.traceManager.TraceRedis("GetLeaderId", key, "")
	result, err := receiver.rdb.Get(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return parse.Convert(result, int64(0))
}

// 续约
func (receiver *redisElection) leaseRenewal(key string) {
	for {
		<-time.After(10 * time.Second)
		_, _ = receiver.rdb.Expire(fs.Context, key, 20*time.Second).Result()
	}
}
