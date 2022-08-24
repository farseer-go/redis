package redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//func Test_redisSet(t *testing.T) {
//
//	client := NewClient("default")
//	defer func() {
//		_, _ = client.Key.Del("key_set")
//		_, _ = client.Key.Del("key_set2")
//		_, _ = client.Key.Del("key_set_diff")
//		_, _ = client.Key.Del("key_set_inter")
//		_, _ = client.Key.Del("key_set_union")
//	}()
//
//	add2, err2 := client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
//	fmt.Println("添加2返回结果：", add2, err2)
//	//添加
//	add, err := client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
//	fmt.Println("添加返回结果：", add, err)
//
//	card, err := client.Set.Card("key_set")
//	fmt.Println("获取数量：", card, err)
//
//	rem, err := client.Set.Rem("key_set", "小王")
//	fmt.Println("移除指定成员返回结果：", rem, err)
//
//	members, err := client.Set.Members("key_set")
//	fmt.Println("获取所有成员：", members, err)
//
//	member, err := client.Set.IsMember("key_set", "小白")
//	fmt.Println("判断指定成员是否存在：", member, err)
//
//	diff, err2 := client.Set.Diff("key_set", "key_set2")
//	fmt.Println("获取差集：", diff, err2)
//
//	store, err2 := client.Set.DiffStore("key_set_diff", "key_set", "key_set2")
//	fmt.Println("存储差集到指定集合：", store, err2)
//
//	inter, err2 := client.Set.Inter("key_set", "key_set2")
//	fmt.Println("获取交集：", inter, err2)
//
//	interStore, err2 := client.Set.InterStore("key_set_inter", "key_set", "key_set2")
//	fmt.Println("存储交集到指定集合：", interStore, err2)
//
//	union, err2 := client.Set.Union("key_set", "key_set2")
//	fmt.Println("获取并集：", union, err2)
//
//	unionStore, err2 := client.Set.UnionStore("key_set_union", "key_set", "key_set2")
//	fmt.Println("存储并集到指定集合：", unionStore, err2)
//
//}

func Test_redisSet_Add(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	add, _ := client.Set.Add("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	assert.Equal(t, add, true)
}

func Test_redisSet_Card(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	card, _ := client.Set.Card("key_set")
	assert.Equal(t, card, int64(6))
}

func Test_redisSet_Diff(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	diff, _ := client.Set.Diff("key_set", "key_set2")
	assert.Equal(t, len(diff), 1)

}

func Test_redisSet_DiffStore(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	defer client.Key.Del("key_set3")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.Set.DiffStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}

func Test_redisSet_Inter(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	inter, _ := client.Set.Inter("key_set", "key_set2")
	assert.Equal(t, len(inter), 3)
}

func Test_redisSet_InterStore(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	defer client.Key.Del("key_set3")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.Set.InterStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}

func Test_redisSet_IsMember(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	client.Set.Add("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	member, _ := client.Set.IsMember("key_set", "小丽")
	assert.Equal(t, member, true)

}

func Test_redisSet_Members(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	client.Set.Add("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	member, _ := client.Set.Members("key_set")
	assert.Equal(t, len(member), 6)
}

func Test_redisSet_Rem(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	client.Set.Add("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	rem, _ := client.Set.Rem("key_set", "小丽")
	assert.Equal(t, rem, true)
}

func Test_redisSet_Union(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	union, _ := client.Set.Union("key_set", "key_set2")
	assert.Equal(t, len(union), 7)
}

func Test_redisSet_UnionStore(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_set")
	defer client.Key.Del("key_set2")
	defer client.Key.Del("key_set3")
	client.Set.Add("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.Set.Add("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.Set.UnionStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}
