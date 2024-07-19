package redis

import (
	"context"
	"errors"
	"github.com/farseer-go/fs/flog"
	"github.com/go-redis/redis/v8"
	"github.com/timandy/routine"
)

type redisPipeline struct {
	*redisManager
}

func (receiver *redisPipeline) Transaction(executeFn func()) error {
	var err error
	traceDetail := receiver.traceManager.TraceRedis("TxPipeline", "", "")
	defer func() { traceDetail.End(err) }()

	// 开启事务
	txPipeline := receiver.GetClient().TxPipeline()
	flog.Infof("使用事务：%d", routine.Goid())
	routineRedisClient.Set(txPipeline)
	defer func() {
		routineRedisClient.Remove()
		flog.Infof("移除事务：%d", routine.Goid())
	}()
	executeFn()
	_, err = txPipeline.Exec(context.Background())
	return err
}

func (receiver *redisPipeline) Pipeline(executeFn func()) (PipelineCmder, error) {
	var err error
	traceDetail := receiver.traceManager.TraceRedis("Pipeline", "", "")
	defer func() { traceDetail.End(err) }()

	// 开启管道
	txPipeline := receiver.GetClient().Pipeline()
	flog.Infof("使用管道：%d", routine.Goid())
	routineRedisClient.Set(txPipeline)
	defer func() {
		routineRedisClient.Remove()
		flog.Infof("移除管道：%d", routine.Goid())
	}()

	executeFn()
	result, err := txPipeline.Exec(context.Background())
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	return PipelineCmder{cmder: result}, err
}
