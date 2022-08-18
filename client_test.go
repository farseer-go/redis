package redis

import (
	"fmt"
	"github.com/farseer-go/fs/configure"
	"testing"
)

// 测试
func TestClient(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	err := client.String.Set("key1", "...12312")
	if err == nil {
		fmt.Printf("设置值:%v\n", "...12312")
	}
	get, err := client.Key.Exists("key1")
	if err != nil {
		fmt.Printf("错误信息：%v\n", err)
	} else {
		fmt.Printf("是否存在：%v\n", get)
	}

	remove, _ := client.Key.Del("key1")
	fmt.Printf("是否删除：%v\n", remove)

	s, _ := client.String.Get("key1")
	fmt.Printf("获取值：%v\n", s)

	get2, err2 := client.Key.Exists("key1")
	if err != nil {
		fmt.Printf("错误信息：%v\n", err2)
	} else {
		fmt.Printf("是否存在：%v\n", get2)
	}
}
