package redis

import (
	"encoding/json"
	"github.com/farseer-go/fs/types"
	"github.com/go-redis/redis/v8"
	"reflect"
)

type PipelineCmder struct {
	cmder []redis.Cmder
}

func (receiver PipelineCmder) ToList(sliceOrList any) {
	sliceOrListVal := reflect.ValueOf(sliceOrList).Elem()

	// 切片类型
	if sliceOrListType, isSlice := types.IsSlice(sliceOrListVal); isSlice {
		value := reflect.MakeSlice(sliceOrListType, 0, 0)
		for _, cmder := range receiver.cmder {
			switch cmderType := cmder.(type) {
			case *redis.StringCmd:
				item := reflect.New(sliceOrListType.Elem()).Interface()
				if jsonVal := cmderType.Val(); jsonVal != "" {
					_ = json.Unmarshal([]byte(jsonVal), item)
					value = reflect.Append(value, reflect.ValueOf(item).Elem())
				}
			}
		}
		sliceOrListVal.Set(value)
		return
	}

	if sliceOrListType, isList := types.IsList(sliceOrListVal); isList {
		itemType := types.GetListItemType(sliceOrListType)
		// 初始化
		value := types.ListNew(sliceOrListType)
		for _, cmder := range receiver.cmder {
			switch cmderType := cmder.(type) {
			case *redis.StringCmd:
				item := reflect.New(itemType).Interface()
				if jsonVal := cmderType.Val(); jsonVal != "" {
					_ = json.Unmarshal([]byte(jsonVal), item)
					types.ListAdd(value, item)
				}
			}
		}
		sliceOrListVal.Set(value.Elem())
		return
	}

	panic("sliceOrList入参必须为切片或collections.List类型")
}
