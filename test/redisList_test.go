package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redisList(t *testing.T) {
	//client := newClient("default")
	//defer func() {
	//	_, _ = client.Del("key_list")
	//}()
	//
	////测试push
	//push, err := client.ListPush("key_list", "北京", "上海", "广州", "海南", "河南")
	//flog.Info("添加返回结果", push, err)
	//
	//set, err := client.ListSet("key_list", 0, "深圳")
	//flog.Info("设置指定值返回结果：", set, err)
	//
	//rem, err := client.ListSetRem("key_list", 0, "上海")
	//flog.Info("移除指定值返回结果：", rem, err)
	//
	//i, err := client.ListLen("key_list")
	//flog.Info("获取指定长度返回结果：", i, err)
	//
	//strings, err := client.ListRange("key_list", 0, i-1)
	//flog.Info("遍历key下所有数据：", strings, err)
	//
	//pop, err := client.ListBLPop(3*time.Second, "key_list")
	//flog.Info("没有阻塞的情况下返回结果：", pop, err)
}

func Test_redisList_BLPop(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	pop, _ := client.ListLeftPop(3*time.Second, "key_list")
	assert.Equal(t, pop[1], "北京")
}

func Test_redisList_Len(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	le, _ := client.ListCount("key_list")
	assert.Equal(t, le, int64(5))
}

func Test_redisList_Push(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	le, _ := client.ListCount("key_list")
	assert.Equal(t, le, int64(5))
}

func Test_redisList_Range(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	strings, _ := client.ListRange("key_list", 1, 2)
	assert.Equal(t, strings[0], "上海")
}

func Test_redisList_Rem(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	rem, _ := client.ListRemove("key_list", 0, "上海")
	assert.Equal(t, rem, true)
}

func Test_redisList_Set(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_list")
	client.ListPushRight("key_list", "北京", "上海", "广州", "海南", "河南")
	set, _ := client.ListSet("key_list", 2, "河南")
	assert.Equal(t, set, true)
}
