package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
	"time"
)

// 选举
type redisElection struct {
	*redisManager
}

// Election 选举成功后，会自动续约。
// 未拿到master，会持续尝试获取master
func (receiver *redisElection) Election(key string, fn func()) {
	for {
		cmd := receiver.GetClient().SetNX(fs.Context, key, core.AppId, 20*time.Second)
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
		duration, _ := receiver.GetClient().TTL(fs.Context, key).Result()
		<-time.After(duration)
	}
}

// GetLeaderId 获取当前LeaderId
func (receiver *redisElection) GetLeaderId(key string) int64 {
	traceDetail := receiver.traceManager.TraceRedis("GetLeaderId", key, "")
	result, err := receiver.GetClient().Get(fs.Context, key).Result()
	defer func() { traceDetail.End(err) }()
	return parse.Convert(result, int64(0))
}

// 续约
func (receiver *redisElection) leaseRenewal(key string) {
	for {
		<-time.After(10 * time.Second)
		_, _ = receiver.GetClient().Expire(fs.Context, key, 20*time.Second).Result()
	}
}
