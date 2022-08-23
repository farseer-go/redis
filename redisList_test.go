package redis

import (
	"fmt"
	"github.com/farseer-go/fs/configure"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redisList(t *testing.T) {

	client := NewClient("default")
	defer func() {
		_, _ = client.Key.Del("key_list")
	}()

	//测试push
	push, err := client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	fmt.Println("添加返回结果", push, err)

	set, err := client.List.Set("key_list", 0, "深圳")
	fmt.Println("设置指定值返回结果：", set, err)

	rem, err := client.List.Rem("key_list", 0, "上海")
	fmt.Println("移除指定值返回结果：", rem, err)

	i, err := client.List.Len("key_list")
	fmt.Println("获取指定长度返回结果：", i, err)

	strings, err := client.List.Range("key_list", 0, i-1)
	fmt.Println("遍历key下所有数据：", strings, err)

	pop, err := client.List.BLPop(3*time.Second, "key_list")
	fmt.Println("没有阻塞的情况下返回结果：", pop, err)
}

func Test_redisList_BLPop(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	pop, _ := client.List.BLPop(3*time.Second, "key_list")
	assert.Equal(t, pop[1], "北京")
}

func Test_redisList_Len(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	le, _ := client.List.Len("key_list")
	assert.Equal(t, le, int64(5))
}

func Test_redisList_Push(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	le, _ := client.List.Len("key_list")
	assert.Equal(t, le, int64(5))
}

func Test_redisList_Range(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	strings, _ := client.List.Range("key_list", 1, 2)
	assert.Equal(t, strings[0], "上海")
}

func Test_redisList_Rem(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	rem, _ := client.List.Rem("key_list", 0, "上海")
	assert.Equal(t, rem, true)
}

func Test_redisList_Set(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_list")
	client.List.Push("key_list", "北京", "上海", "广州", "海南", "河南")
	set, _ := client.List.Set("key_list", 2, "河南")
	assert.Equal(t, set, true)
}
