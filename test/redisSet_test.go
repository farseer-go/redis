package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//func Test_redisSet(t *testing.T) {
//
//	client := newClient("default")
//	defer func() {
//		_, _ = client.Del("key_set")
//		_, _ = client.Del("key_set2")
//		_, _ = client.Del("key_set_diff")
//		_, _ = client.Del("key_set_inter")
//		_, _ = client.Del("key_set_union")
//	}()
//
//	add2, err2 := client.SetSetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
//	flog.Info("添加2返回结果：", add2, err2)
//	//添加
//	add, err := client.SetSetAdd("key_set", "小丽", "小王", "小白", "小小")
//	flog.Info("添加返回结果：", add, err)
//
//	card, err := client.SetCount("key_set")
//	flog.Info("获取数量：", card, err)
//
//	rem, err := client.SetSetRem("key_set", "小王")
//	flog.Info("移除指定成员返回结果：", rem, err)
//
//	members, err := client.SetMembers("key_set")
//	flog.Info("获取所有成员：", members, err)
//
//	member, err := client.SetIsMember("key_set", "小白")
//	flog.Info("判断指定成员是否存在：", member, err)
//
//	diff, err2 := client.SetDiff("key_set", "key_set2")
//	flog.Info("获取差集：", diff, err2)
//
//	store, err2 := client.SetDiffStore("key_set_diff", "key_set", "key_set2")
//	flog.Info("存储差集到指定集合：", store, err2)
//
//	inter, err2 := client.SetInter("key_set", "key_set2")
//	flog.Info("获取交集：", inter, err2)
//
//	interStore, err2 := client.SetInterStore("key_set_inter", "key_set", "key_set2")
//	flog.Info("存储交集到指定集合：", interStore, err2)
//
//	union, err2 := client.SetUnion("key_set", "key_set2")
//	flog.Info("获取并集：", union, err2)
//
//	unionStore, err2 := client.SetUnionStore("key_set_union", "key_set", "key_set2")
//	flog.Info("存储并集到指定集合：", unionStore, err2)
//
//}

func Test_redisSet_Add(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	add, _ := client.SetAdd("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	assert.Equal(t, add, true)
}

func Test_redisSet_Card(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	card, _ := client.SetCount("key_set")
	assert.Equal(t, card, int64(6))
}

func Test_redisSet_Diff(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	diff, _ := client.SetDiff("key_set", "key_set2")
	assert.Equal(t, len(diff), 1)

}

func Test_redisSet_DiffStore(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	defer client.Del("key_set3")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.SetDiffStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}

func Test_redisSet_Inter(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	inter, _ := client.SetInter("key_set", "key_set2")
	assert.Equal(t, len(inter), 3)
}

func Test_redisSet_InterStore(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	defer client.Del("key_set3")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.SetInterStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}

func Test_redisSet_IsMember(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	client.SetAdd("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	member, _ := client.SetIsMember("key_set", "小丽")
	assert.Equal(t, member, true)

}

func Test_redisSet_Members(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	client.SetAdd("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	member, _ := client.SetGet("key_set")
	assert.Equal(t, len(member), 6)
}

func Test_redisSet_Rem(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	client.SetAdd("key_set", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	rem, _ := client.SetRemove("key_set", "小丽")
	assert.Equal(t, rem, true)
}

func Test_redisSet_Union(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	union, _ := client.SetUnion("key_set", "key_set2")
	assert.Equal(t, len(union), 7)
}

func Test_redisSet_UnionStore(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_set")
	defer client.Del("key_set2")
	defer client.Del("key_set3")
	client.SetAdd("key_set2", "小丽", "小王", "小白", "小赵", "小钱", "小孙")
	//添加
	client.SetAdd("key_set", "小丽", "小王", "小白", "小小")
	store, _ := client.SetUnionStore("key_set3", "key_set", "key_set2")
	assert.Equal(t, store, true)
}
