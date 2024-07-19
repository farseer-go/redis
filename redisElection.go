package redis

import (
	"context"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
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
		// 拿到锁了
		if result, _ := cmd.Result(); result {
			ctx, cancel := context.WithCancel(fs.Context)
			// 给锁续租约
			go receiver.leaseRenewal(key, ctx)
			fn()
			cancel()
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
func (receiver *redisElection) leaseRenewal(key string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			receiver.GetClient().Del(fs.Context, key).Result()
			flog.Infof(key + "退出")
			return
		case <-time.After(10 * time.Second):
			for {
				result, _ := receiver.GetClient().Expire(ctx, key, 20*time.Second).Result()
				if result {
					flog.Infof(key + "续约成功")
					break
				}
				flog.Infof(key + "续约失败")
				time.Sleep(time.Second)
			}
		}
	}
}
