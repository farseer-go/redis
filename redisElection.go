package redis

import (
	"context"
	"time"

	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
)

// 选举
type redisElection struct {
	*redisManager
}

// Election 选举成功后，会自动续约。
// 未拿到master，会持续尝试获取master
func (receiver *redisElection) Election(ctx context.Context, key string, fn func()) {
	for {
		cmd := receiver.rdb.SetNX(ctx, key, core.AppId, 20*time.Second)
		// 拿到锁了
		if result, err := cmd.Result(); result && err == nil {
			ctx2, cancel := context.WithCancel(ctx)
			defer cancel()

			// 给锁续租约
			go receiver.leaseRenewal(key, ctx2)
			fn()
			return
		}

		// 没有拿到master的节点，需获取当前租约剩余时间，到期后，尝试获取
		duration, err := receiver.rdb.TTL(context.Background(), key).Result()
		if err != nil || duration <= 0 {
			duration = time.Second // 兜底，避免忙循环
		}

		select {
		case <-ctx.Done(): // 立即退出
			return
		case <-time.After(duration):
		}
	}
}

// GetLeaderId 获取当前LeaderId
func (receiver *redisElection) GetLeaderId(key string) int64 {
	traceDetail := receiver.traceManager.TraceRedis("GetLeaderId", key, "")
	result, err := receiver.rdb.Get(context.Background(), key).Result()
	defer func() { traceDetail.End(err) }()
	return parse.Convert(result, int64(0))
}

// 续约
func (receiver *redisElection) leaseRenewal(key string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			_, _ = receiver.rdb.Del(context.Background(), key).Result()
			return
		case <-time.After(10 * time.Second):
			_, _ = receiver.rdb.Expire(ctx, key, 20*time.Second).Result()
		}
	}
}
