package redis

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/linkTrace"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
)

type redisHash struct {
	rdb *redis.Client
}

func (redisHash *redisHash) HashSetEntity(key string, field string, entity any) error {
	trace := linkTrace.TraceRedis("HashSetEntity", key, field)
	jsonContent, err := json.Marshal(entity)
	defer func() { trace.End(err) }()

	if err != nil {
		return err
	}

	err = redisHash.rdb.HSet(fs.Context, key, field, string(jsonContent)).Err()
	return err
}

func (redisHash *redisHash) HashSet(key string, fieldValues ...any) error {
	var fields []string
	for i := 0; i < len(fieldValues); i += 2 {
		fields = append(fields, parse.ToString(fieldValues[i]))
	}

	trace := linkTrace.TraceRedis("HashSet", key, strings.Join(fields, ","))
	err := redisHash.rdb.HSet(fs.Context, key, fieldValues...).Err()
	defer func() { trace.End(err) }()
	return err
}

func (redisHash *redisHash) HashGet(key string, field string) (string, error) {
	trace := linkTrace.TraceRedis("HashGet", key, field)

	result, err := redisHash.rdb.HGet(fs.Context, key, field).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisHash *redisHash) HashGetAll(key string) (map[string]string, error) {
	trace := linkTrace.TraceRedis("HashGetAll", key, "")

	result, err := redisHash.rdb.HGetAll(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisHash *redisHash) HashToEntity(key string, field string, entity any) (bool, error) {
	trace := linkTrace.TraceRedis("HashToEntity", key, field)

	jsonContent, err := redisHash.rdb.HGet(fs.Context, key, field).Result()
	defer func() { trace.End(err) }()
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

func (redisHash *redisHash) HashToArray(key string, arrSlice any) error {
	trace := linkTrace.TraceRedis("HashToArray", key, "")
	arrVal := reflect.ValueOf(arrSlice).Elem()
	arrType, isSlice := types.IsSlice(arrVal)
	if !isSlice {
		panic("arr入参必须为切片类型")
	}

	result, err := redisHash.rdb.HGetAll(fs.Context, key).Result()
	defer func() { trace.End(err) }()
	if err != nil {
		return flog.Error(err)
	}

	lst := collections.NewListAny()
	for _, vJson := range result {
		item := reflect.New(arrType.Elem()).Interface()
		_ = json.Unmarshal([]byte(vJson), item)
		lst.Add(reflect.ValueOf(item).Elem().Interface())
	}

	lst.MapToArray(arrSlice)
	return nil
}

func (redisHash *redisHash) HashToListAny(key string, itemType reflect.Type) (collections.ListAny, error) {
	trace := linkTrace.TraceRedis("HashToListAny", key, "")

	lst := collections.NewListAny()
	result, err := redisHash.rdb.HGetAll(fs.Context, key).Result()
	defer func() { trace.End(err) }()
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

func (redisHash *redisHash) HashExists(key string, field string) (bool, error) {
	trace := linkTrace.TraceRedis("HashExists", key, field)

	result, err := redisHash.rdb.HExists(fs.Context, key, field).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisHash *redisHash) HashDel(key string, fields ...string) (bool, error) {
	trace := linkTrace.TraceRedis("HashDel", key, strings.Join(fields, ","))

	result, err := redisHash.rdb.HDel(fs.Context, key, fields...).Result()
	defer func() { trace.End(err) }()
	return result > 0, err
}

func (redisHash *redisHash) HashCount(key string) int {
	trace := linkTrace.TraceRedis("HashCount", key, "")

	result, err := redisHash.rdb.HLen(fs.Context, key).Uint64()
	defer func() { trace.End(err) }()
	return int(result)
}

func (redisHash *redisHash) HashIncrInt(key string, field string, value int) (int, error) {
	trace := linkTrace.TraceRedis("HashIncrInt", key, field)

	result, err := redisHash.rdb.HIncrBy(fs.Context, key, field, parse.Convert(value, int64(value))).Result()
	defer func() { trace.End(err) }()
	return parse.ToInt(result), err
}

func (redisHash *redisHash) HashIncrInt64(key string, field string, value int64) (int64, error) {
	trace := linkTrace.TraceRedis("HashIncrInt64", key, field)

	result, err := redisHash.rdb.HIncrBy(fs.Context, key, field, value).Result()
	defer func() { trace.End(err) }()
	return result, err
}

func (redisHash *redisHash) HashIncrFloat32(key string, field string, value float32) (float32, error) {
	trace := linkTrace.TraceRedis("HashIncrFloat32", key, field)

	result, err := redisHash.rdb.HIncrByFloat(fs.Context, key, field, float64(value)).Result()
	defer func() { trace.End(err) }()
	return parse.ToFloat32(result), err
}

func (redisHash *redisHash) HashIncrFloat64(key string, field string, value float64) (float64, error) {
	trace := linkTrace.TraceRedis("HashIncrFloat64", key, field)

	result, err := redisHash.rdb.HIncrByFloat(fs.Context, key, field, value).Result()
	defer func() { trace.End(err) }()
	return result, err
}
