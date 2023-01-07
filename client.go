package redis

import (
	"context"
	"github.com/farseer-go/fs/configure"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client struct {
	Key    *redisKey
	String *redisString
	Hash   *redisHash
	List   *redisList
	Set    *redisSet
	ZSet   *redisZSet
	Lock   *redisLock
}

// 上下文定义
var ctx = context.Background()

// NewClient 初始化
func NewClient(redisName string) *Client {
	configString := configure.GetString("Redis." + redisName)
	if configString == "" {
		panic("[farseer.yaml]找不到相应的配置：Redis." + redisName)
	}
	redisConfig := configure.ParseString[redisConfig](configString)
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Server,                                            //localhost:6379
		Password:     redisConfig.Password,                                          // no password Set
		DB:           redisConfig.DB,                                                // use default DB
		DialTimeout:  time.Duration(redisConfig.ConnectTimeout) * time.Millisecond,  //链接超时时间设置
		WriteTimeout: time.Duration(redisConfig.SyncTimeout) * time.Millisecond,     //同步超时时间设置
		ReadTimeout:  time.Duration(redisConfig.ResponseTimeout) * time.Millisecond, //响应超时时间设置
	})
	key := &redisKey{rdb: rdb}
	str := &redisString{rdb: rdb}
	hash := &redisHash{rdb: rdb}
	list := &redisList{rdb: rdb}
	set := &redisSet{rdb: rdb}
	zset := &redisZSet{rdb: rdb}
	lock := &redisLock{rdb: rdb}
	return &Client{Key: key, String: str, Hash: hash, List: list, Set: set, ZSet: zset, Lock: lock}
}
