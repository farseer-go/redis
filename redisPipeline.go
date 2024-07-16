package redis

import (
	"context"
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
	routineRedisClient.Set(txPipeline)
	defer func() { routineRedisClient.Remove() }()
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
	routineRedisClient.Set(txPipeline)
	defer func() { routineRedisClient.Remove() }()

	executeFn()
	result, err := txPipeline.Exec(context.Background())
	return PipelineCmder{cmder: result}, err
}
