package redis

import (
	"encoding/json"
	"errors"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
)

type redisHash struct {
	*redisManager
}

func (receiver *redisHash) HashSetEntity(key string, field string, entity any) error {
	traceDetail := receiver.traceManager.TraceRedis("HashSetEntity", key, field)
	jsonContent, err := json.Marshal(entity)
	defer func() { traceDetail.End(err) }()

	if err != nil {
		return err
	}

	err = receiver.GetClient().HSet(fs.Context, key, field, string(jsonContent)).Err()
	return err
}

func (receiver *redisHash) HashSet(key string, fieldValues ...any) error {
	var fields []string
	for i := 0; i < len(fieldValues); i += 2 {
		fields = append(fields, parse.ToString(fieldValues[i]))
	}

	traceDetail := receiver.traceManager.TraceRedis("HashSet", key, strings.Join(fields, ","))
	err := receiver.GetClient().HSet(fs.Context, key, fieldValues...).Err()
	defer func() { traceDetail.End(err) }()
	return err
}

func (receiver *redisHash) HashGet(key string, field string) (string, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashGet", key, field)

	result, err := receiver.GetClient().HGet(fs.Context, key, field).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisHash) HashGetAll(key string) (map[string]string, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashGetAll", key, "")

	result, err := receiver.GetClient().HGetAll(fs.Context, key).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisHash) HashToEntity(key string, field string, entity any) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashToEntity", key, field)

	jsonContent, err := receiver.GetClient().HGet(fs.Context, key, field).Result()
	defer func() { traceDetail.End(err) }()
	if err == redis.Nil {
		err = nil
		return false, err
	}
	if err != nil {
		return false, err
	}
	// 反序列
	return true, json.Unmarshal([]byte(jsonContent), entity)
}

func (receiver *redisHash) HashToArray(key string, arrSlice any) error {
	traceDetail := receiver.traceManager.TraceRedis("HashToArray", key, "")
	arrVal := reflect.ValueOf(arrSlice).Elem()
	arrType, isSlice := types.IsSlice(arrVal)
	if !isSlice {
		panic("arr入参必须为切片类型")
	}

	result, err := receiver.GetClient().HGetAll(fs.Context, key).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	if err != nil {
		return flog.Error(err)
	}

	newArr := reflect.MakeSlice(arrVal.Type(), 0, 0)
	for _, vJson := range result {
		item := reflect.New(arrType.Elem()).Interface()
		_ = json.Unmarshal([]byte(vJson), item)
		newArr = reflect.Append(newArr, reflect.ValueOf(item).Elem())
	}
	arrVal.Set(newArr)
	return nil
}

func (receiver *redisHash) HashToListAny(key string, itemType reflect.Type) (collections.ListAny, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashToListAny", key, "")

	lst := collections.NewListAny()
	result, err := receiver.GetClient().HGetAll(fs.Context, key).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	if err != nil {
		_ = flog.Error(err)
		return lst, err
	}
	for _, vJson := range result {
		item := reflect.New(itemType).Interface()
		_ = json.Unmarshal([]byte(vJson), item)
		lst.Add(reflect.ValueOf(item).Elem().Interface())
	}
	return lst, nil
}

func (receiver *redisHash) HashExists(key string, field string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashExists", key, field)

	result, err := receiver.GetClient().HExists(fs.Context, key, field).Result()
	if err == redis.Nil {
		err = nil
	}
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisHash) HashDel(key string, fields ...string) (bool, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashDel", key, strings.Join(fields, ","))

	result, err := receiver.GetClient().HDel(fs.Context, key, fields...).Result()
	defer func() { traceDetail.End(err) }()
	return result > 0, err
}

func (receiver *redisHash) HashCount(key string) int {
	traceDetail := receiver.traceManager.TraceRedis("HashCount", key, "")

	result, err := receiver.GetClient().HLen(fs.Context, key).Uint64()
	defer func() { traceDetail.End(err) }()
	return int(result)
}

func (receiver *redisHash) HashIncrInt(key string, field string, value int) (int, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashIncrInt", key, field)

	result, err := receiver.GetClient().HIncrBy(fs.Context, key, field, parse.Convert(value, int64(value))).Result()
	defer func() { traceDetail.End(err) }()
	return parse.ToInt(result), err
}

func (receiver *redisHash) HashIncrInt64(key string, field string, value int64) (int64, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashIncrInt64", key, field)

	result, err := receiver.GetClient().HIncrBy(fs.Context, key, field, value).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}

func (receiver *redisHash) HashIncrFloat32(key string, field string, value float32) (float32, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashIncrFloat32", key, field)

	result, err := receiver.GetClient().HIncrByFloat(fs.Context, key, field, float64(value)).Result()
	defer func() { traceDetail.End(err) }()
	return parse.ToFloat32(result), err
}

func (receiver *redisHash) HashIncrFloat64(key string, field string, value float64) (float64, error) {
	traceDetail := receiver.traceManager.TraceRedis("HashIncrFloat64", key, field)

	result, err := receiver.GetClient().HIncrByFloat(fs.Context, key, field, value).Result()
	defer func() { traceDetail.End(err) }()
	return result, err
}
