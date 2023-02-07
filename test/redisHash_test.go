package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

//
//// String 测试
//func TestClientHash(t *testing.T) {
//	client := newClient("default")
//
//	defer func() {
//		_, _ = client.Del("key_has1")
//		_, _ = client.Del("key_has2")
//	}()
//
//	err := client.HashSet("key_has1", "name", "小丽")
//	errV2 := client.HashSet("key_has1", "age", 40, "address", "上海")
//
//	if err == nil {
//		flog.Info("设置key_has1值成功.\n")
//	} else {
//		flog.Info("设置key_has1值错误:%v\n", err)
//	}
//
//	if errV2 == nil {
//		flog.Info("设置key_has1 v2 值成功.\n")
//	} else {
//		flog.Info("设置key_has1 v2 值错误:%v\n", errV2)
//	}
//
//	get, _ := client.HashGet("key_has1", "name")
//	flog.Info("获取key_has1  单个 name 值成功:%v\n", get)
//
//	all, _ := client.HashGetAll("key_has1")
//	flog.Info("获取key_has1  所有 值成功:%v\n", all)
//
//	exists, _ := client.HashExists("key_has1", "age")
//	flog.Info("age值是否存在:%v\n", exists)
//
//	get2, _ := client.HashGet("key_has1", "age")
//	flog.Info("获取key_has2  单个 age 值成功:%v\n", get2)
//
//	remove, _ := client.HashDel("key_has1", "age")
//	flog.Info("移出age成员:%v\n", remove)
//
//	err2 := client.HashSet("key_has2", "key1", "value1", "key2", 222)
//	if err2 == nil {
//		flog.Info("设置key_has2值成功.\n")
//	} else {
//		flog.Info("设置key_has2值错误:%v\n", err2)
//	}
//	all2, _ := client.HashGetAll("key_has2")
//	flog.Info("获取key_has2  所有 值成功:%v\n", all2)
//
//	//SetMap
//	//umap := map[string]string{"user": "harlen", "city": "河南", "age": "30"}
//	//err3 := client.HashSetMap("key_has3", umap)
//	//if err3 == nil {
//	//	flog.Info("设置key_has3值成功.\n")
//	//}
//	//all3, _ := client.HashGetAll("key_has3")
//	//flog.Info("获取key_has3  所有 值成功:%v\n", all3)
//}

func Test_redisHash_Count(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "address", "上海")
	count := client.HashCount("key_client")
	assert.Equal(t, count, 1)
}

func Test_redisHash_Del(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	del, _ := client.HashDel("key_client", "age")
	assert.Equal(t, del, true)
}

func Test_redisHash_Exists(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	exists, _ := client.HashExists("key_client", "age")
	assert.Equal(t, exists, true)
}

func Test_redisHash_Get(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	get, _ := client.HashGet("key_client", "age")
	atoi, _ := strconv.Atoi(get) //类型转换 string  转 int
	assert.Equal(t, atoi, 40)
}

func Test_redisHash_GetAll(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	client.HashSet("key_client", "age", 40, "address", "上海")
	all, _ := client.HashGetAll("key_client")
	assert.Equal(t, len(all), 2)
}

func Test_redisHash_Set(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	err := client.HashSet("key_client", "age", 40, "address", "上海")
	assert.Equal(t, err, nil)
}

func Test_redisHash_SetEntity(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	type user struct {
		Name string
		Age  int
		Sex  string
	}
	entity := user{Name: "小吴", Age: 20, Sex: "男"}
	client.HashSetEntity("key_client", "json", entity)
	var entity2 user
	client.HashToEntity("key_client", "json", &entity2)
	assert.Equal(t, entity2.Age, 20)
}

func Test_redisHash_ToArray(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	type user struct {
		Name string
		Age  int
		Sex  string
	}
	entity := user{Name: "小吴", Age: 20, Sex: "男"}
	client.HashSetEntity("key_client", "json", entity)
	var arrVal []user
	client.HashToArray("key_client", &arrVal)
	assert.Equal(t, arrVal[0].Age, 20)
}

func Test_redisHash_ToEntity(t *testing.T) {

	fs.Initialize[redis.Module]("unit test")
	client := container.Resolve[redis.IClient]("default")
	defer client.Del("key_client")
	type user struct {
		Name string
		Age  int
		Sex  string
	}
	var entity2 user
	entity := user{Name: "小吴", Age: 20, Sex: "男"}
	client.HashSetEntity("key_client", "json", entity)
	client.HashToEntity("key_client", "json", &entity2)
	assert.Equal(t, entity2.Sex, "男")
}
