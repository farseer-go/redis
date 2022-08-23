package redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

//
//// String 测试
//func TestClientHash(t *testing.T) {
//	client := NewClient("default")
//
//	defer func() {
//		_, _ = client.Key.Del("key_has1")
//		_, _ = client.Key.Del("key_has2")
//	}()
//
//	err := client.Hash.Set("key_has1", "name", "小丽")
//	errV2 := client.Hash.Set("key_has1", "age", 40, "address", "上海")
//
//	if err == nil {
//		fmt.Printf("设置key_has1值成功.\n")
//	} else {
//		fmt.Printf("设置key_has1值错误:%v\n", err)
//	}
//
//	if errV2 == nil {
//		fmt.Printf("设置key_has1 v2 值成功.\n")
//	} else {
//		fmt.Printf("设置key_has1 v2 值错误:%v\n", errV2)
//	}
//
//	get, _ := client.Hash.Get("key_has1", "name")
//	fmt.Printf("获取key_has1  单个 name 值成功:%v\n", get)
//
//	all, _ := client.Hash.GetAll("key_has1")
//	fmt.Printf("获取key_has1  所有 值成功:%v\n", all)
//
//	exists, _ := client.Hash.Exists("key_has1", "age")
//	fmt.Printf("age值是否存在:%v\n", exists)
//
//	get2, _ := client.Hash.Get("key_has1", "age")
//	fmt.Printf("获取key_has2  单个 age 值成功:%v\n", get2)
//
//	remove, _ := client.Hash.Del("key_has1", "age")
//	fmt.Printf("移出age成员:%v\n", remove)
//
//	err2 := client.Hash.Set("key_has2", "key1", "value1", "key2", 222)
//	if err2 == nil {
//		fmt.Printf("设置key_has2值成功.\n")
//	} else {
//		fmt.Printf("设置key_has2值错误:%v\n", err2)
//	}
//	all2, _ := client.Hash.GetAll("key_has2")
//	fmt.Printf("获取key_has2  所有 值成功:%v\n", all2)
//
//	//SetMap
//	//umap := map[string]string{"user": "harlen", "city": "河南", "age": "30"}
//	//err3 := client.Hash.SetMap("key_has3", umap)
//	//if err3 == nil {
//	//	fmt.Printf("设置key_has3值成功.\n")
//	//}
//	//all3, _ := client.Hash.GetAll("key_has3")
//	//fmt.Printf("获取key_has3  所有 值成功:%v\n", all3)
//}

func Test_redisHash_Count(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "address", "上海")
	count := client.Hash.Count("key_client")
	assert.Equal(t, count, 1)
}

func Test_redisHash_Del(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	del, _ := client.Hash.Del("key_client", "age")
	assert.Equal(t, del, true)
}

func Test_redisHash_Exists(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	exists, _ := client.Hash.Exists("key_client", "age")
	assert.Equal(t, exists, true)
}

func Test_redisHash_Get(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	get, _ := client.Hash.Get("key_client", "age")
	atoi, _ := strconv.Atoi(get) //类型转换 string  转 int
	assert.Equal(t, atoi, 40)
}

func Test_redisHash_GetAll(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	client.Hash.Set("key_client", "age", 40, "address", "上海")
	all, _ := client.Hash.GetAll("key_client")
	assert.Equal(t, len(all), 2)
}

func Test_redisHash_Set(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	err := client.Hash.Set("key_client", "age", 40, "address", "上海")
	assert.Equal(t, err, nil)
}

func Test_redisHash_SetEntity(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	type user struct {
		Name string
		Age  int
		Sex  string
	}
	entity := user{Name: "小吴", Age: 20, Sex: "男"}
	client.Hash.SetEntity("key_client", "json", entity)
	var entity2 user
	client.Hash.ToEntity("key_client", "json", &entity2)
	assert.Equal(t, entity2.Age, 20)
}

// array 方法待验证
func Test_redisHash_ToArray(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	var arrVal []string
	client.Hash.Set("key_client", "Name", "小强", "Age", 40, "Address", "上海")
	client.Hash.ToArray("key_client", &arrVal)
	assert.Equal(t, len(arrVal), 3)
}

func Test_redisHash_ToEntity(t *testing.T) {
	configure.SetDefault("Redis.default", "Server=localhost:6379,DB=15,Password=redis123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000")
	client := NewClient("default")
	defer client.Key.Del("key_client")
	type user struct {
		Name string
		Age  int
		Sex  string
	}
	var entity2 user
	entity := user{Name: "小吴", Age: 20, Sex: "男"}
	client.Hash.SetEntity("key_client", "json", entity)
	client.Hash.ToEntity("key_client", "json", &entity2)
	assert.Equal(t, entity2.Sex, "男")
}
