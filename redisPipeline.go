package redis

import (
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
)

type redisPipeline struct {
	rdb *redis.Client
}

func (receiver *redisPipeline) Transaction(executeFn func(tx redis.Pipeliner)) ([]redis.Cmder, error) {
	var err error
	trace := linkTrace.TraceRedis("TxPipeline", "", "")
	defer func() { trace.End(err) }()

	// 开启事务
	txPipeline := receiver.rdb.TxPipeline()
	executeFn(txPipeline)
	var exec []redis.Cmder
	exec, err = txPipeline.Exec(nil)

	return exec, err
}

func (receiver *redisPipeline) Pipeline(executeFn func(tx redis.Pipeliner)) ([]redis.Cmder, error) {
	var err error
	trace := linkTrace.TraceRedis("TxPipeline", "", "")
	defer func() { trace.End(err) }()

	// 开启管道
	txPipeline := receiver.rdb.Pipeline()
	executeFn(txPipeline)
	var exec []redis.Cmder
	exec, err = txPipeline.Exec(nil)

	return exec, err
}
