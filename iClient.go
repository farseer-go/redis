package redis

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

type IClient interface {
	// Original 获取原生的客户端
	Original() *redis.Client
	// RegisterEvent 注册core.IEvent实现
	RegisterEvent(eventName string, fns ...core.ConsumerFunc)

	// SetTTL 设置过期时间
	SetTTL(key string, d time.Duration) (bool, error)
	// TTL 获取过期时间
	TTL(key string) (time.Duration, error)
	// Del 删除
	Del(keys ...string) (bool, error)
	// Exists key值是否存在
	Exists(keys ...string) (bool, error)

	// StringSet 设置缓存
	StringSet(key string, value any) error
	// StringGet 获取缓存
	StringGet(key string) (string, error)
	// StringSetNX 设置过期时间
	StringSetNX(key string, value any, expiration time.Duration) (bool, error)

	// HashSetEntity 添加并序列化成json
	HashSetEntity(key string, field string, entity any) error
	// HashSet 添加
	//   - HashSet("myHashKey", "field1", "value1", "field2", "value2")
	//   - HashSet("myHashKey", []string{"field1", "value1", "field2", "value2"})
	//   - HashSet("myHashKey", map[string]any{"field1": "value1", "field2": "value2"})
	HashSet(key string, fieldValues ...any) error
	// HashGet 获取
	HashGet(key string, field string) (string, error)
	// HashGetAll 获取所有集合数据
	HashGetAll(key string) (map[string]string, error)
	// HashToEntity 获取单个对象
	//	var client DomainObject
	//	HashToEntity("redisKey", "field", &client)
	HashToEntity(key string, field string, entity any) (bool, error)
	// HashToArray 将hash.value反序列化成切片对象
	//	type record struct {
	//		ClientId int `json:"id"`
	//	}
	//	var records []record
	//	HashToArray("test", &records)
	HashToArray(key string, arrSlice any) error
	// HashToListAny 将hash的数据转成collections.ListAny
	HashToListAny(key string, itemType reflect.Type) (collections.ListAny, error)
	// HashExists 成员是否存在
	HashExists(key string, field string) (bool, error)
	// HashDel 移除指定成员
	HashDel(key string, fields ...string) (bool, error)
	// HashCount 获取hash的数量
	HashCount(key string) int

	// HashIncrInt Hash对int加减
	HashIncrInt(key string, field string, value int) (int, error)
	// HashIncrInt64 Hash对int64加减
	HashIncrInt64(key string, field string, value int64) (int64, error)
	// HashIncrFloat32 Hash对float32加减
	HashIncrFloat32(key string, field string, value float32) (float32, error)
	// HashIncrFloat64 Hash对float64加减
	HashIncrFloat64(key string, field string, value float64) (float64, error)

	// ListPushRight 向末层推送数据
	ListPushRight(key string, values ...any) (bool, error)
	// ListPushLeft 向头部推送数据
	ListPushLeft(key string, values ...any) (bool, error)
	// ListSet 指定索引设置值
	ListSet(key string, index int64, value any) (bool, error)
	// ListRemove 移除指定数量的value，
	// count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
	// count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
	// count = 0 : 移除表中所有与 VALUE 相等的值。
	ListRemove(key string, count int64, value any) (bool, error)
	// ListCount 获取长度
	ListCount(key string) (int64, error)
	// ListRange 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。
	// 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，以此类推。
	// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
	ListRange(key string, start int64, stop int64) ([]string, error)
	// ListLeftPop 命令移出并获取列表的第一个元素（头部）， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
	ListLeftPop(timeout time.Duration, keys ...string) ([]string, error)
	// ListRightPop 命令移出并获取列表的最后一个元素（尾部）， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
	ListRightPop(timeout time.Duration, keys ...string) ([]string, error)
	// ListRightPopPush 取出最后一个元素（尾部）并Push到destination
	ListRightPopPush(source, destination string, timeout time.Duration) (string, error)

	// SetAdd 添加
	SetAdd(key string, members ...any) (bool, error)
	// SetCount 获取数量
	SetCount(key string) (int64, error)
	// SetRemove 移除指定成员
	SetRemove(key string, members ...any) (bool, error)
	// SetGet 获取所有成员
	SetGet(key string) ([]string, error)
	// SetIsMember 判断指定成员是否存在
	SetIsMember(key string, member any) (bool, error)
	// SetDiff 获取差集
	SetDiff(keys ...string) ([]string, error)
	// SetDiffStore 将差集，保存在指定集合中
	SetDiffStore(destination string, keys ...string) (bool, error)
	// SetInter 获取交集
	SetInter(keys ...string) ([]string, error)
	// SetInterStore 将交集，保存在指定集合中
	SetInterStore(destination string, keys ...string) (bool, error)
	// SetUnion 获取并集
	SetUnion(keys ...string) ([]string, error)
	// SetUnionStore 将并集，保存在指定集合中
	SetUnionStore(destination string, keys ...string) (bool, error)

	// ZSetAdd 添加
	ZSetAdd(key string, members ...*redisZ) (bool, error)
	// ZSetScore 获取指定成员score
	ZSetScore(key string, member string) (float64, error)
	// ZSetRange 获取有序集合指定区间内的成员
	ZSetRange(key string, start int64, stop int64) ([]string, error)
	// ZSetRevRange 获取有序集合指定区间内的成员分数从高到低
	ZSetRevRange(key string, start int64, stop int64) ([]string, error)
	// ZSetRangeByScore 获取指定分数区间的成员列表
	ZSetRangeByScore(key string, opt *redisZRangeBy) ([]string, error)

	// LockNew 获得一个锁
	LockNew(key, val string, expiration time.Duration) core.ILock

	// Publish 发布消息
	Publish(channel string, message any) (int64, error)
	// Subscribe 订阅消息
	Subscribe(channels ...string) <-chan *redis.Message
	// Election 选举
	Election(key string, fn func())
	// GetLeaderId 获取当前LeaderId
	GetLeaderId(key string) int64
}
