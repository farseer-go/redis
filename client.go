package redis

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/go-redis/redis/v8"
	"time"
)

type client struct {
	redisKey
	redisString
	redisHash
	redisList
	redisSet
	redisZSet
	redisLock
	redisPub
	original *redis.Client
}

// newClient 初始化
func newClient(redisConfig redisConfig) IClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Server,                                            //localhost:6379
		Password:     redisConfig.Password,                                          // no password Set
		DB:           redisConfig.DB,                                                // use default DB
		DialTimeout:  time.Duration(redisConfig.ConnectTimeout) * time.Millisecond,  //链接超时时间设置
		WriteTimeout: time.Duration(redisConfig.SyncTimeout) * time.Millisecond,     //同步超时时间设置
		ReadTimeout:  time.Duration(redisConfig.ResponseTimeout) * time.Millisecond, //响应超时时间设置
	})

	return &client{
		original:    rdb,
		redisKey:    redisKey{rdb: rdb},
		redisString: redisString{rdb: rdb},
		redisHash:   redisHash{rdb: rdb},
		redisList:   redisList{rdb: rdb},
		redisSet:    redisSet{rdb: rdb},
		redisZSet:   redisZSet{rdb: rdb},
		redisLock:   redisLock{rdb: rdb},
		redisPub:    redisPub{rdb: rdb},
	}
}

func (receiver *client) RegisterEvent(eventName string, fns ...core.ConsumerFunc) {
	// 注册仓储
	container.Register(func() core.IEvent {
		return &registerEvent{
			eventName: eventName,
			client:    receiver,
		}
	}, eventName)

	go subscribe(receiver, eventName, fns)
}

// Original 获取原生的客户端
func (receiver *client) Original() *redis.Client {
	return receiver.original
}
