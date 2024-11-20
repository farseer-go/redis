package redis

import (
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"github.com/go-redis/redis/v8"
)

type PipelineCmder struct {
	cmder []redis.Cmder
}

func (receiver PipelineCmder) Fill(sliceOrList any) {
	sliceOrListVal := reflect.ValueOf(sliceOrList).Elem()

	// 切片类型
	if sliceType, isSlice := types.IsSlice(sliceOrListVal); isSlice {
		itemType := sliceType.Elem()
		value := reflect.MakeSlice(sliceType, 0, 0)

		for _, cmder := range receiver.cmder {
			switch cmderType := cmder.(type) {
			case *redis.StringCmd:
				item := reflect.New(itemType).Interface()
				if jsonVal := cmderType.Val(); jsonVal != "" {
					_ = sonic.Unmarshal([]byte(jsonVal), item)
					value = reflect.Append(value, reflect.ValueOf(item).Elem())
				}
			case *redis.StringStringMapCmd:
				jsonVal := cmderType.Val()
				for _, v := range jsonVal {
					value = reflect.Append(value, reflect.ValueOf(parse.ConvertValue(v, itemType)))
				}
			}
		}
		sliceOrListVal.Set(value)
		return
	}

	// List类型
	if listType, isList := types.IsList(sliceOrListVal); isList {
		itemType := types.GetListItemType(listType)
		value := types.ListNew(listType)

		for _, cmder := range receiver.cmder {
			switch cmderType := cmder.(type) {
			case *redis.StringCmd:
				item := reflect.New(itemType).Interface()
				if jsonVal := cmderType.Val(); jsonVal != "" {
					_ = sonic.Unmarshal([]byte(jsonVal), item)
					types.ListAdd(value, item)
				}
			case *redis.StringStringMapCmd:
				jsonVal := cmderType.Val()
				for _, v := range jsonVal {
					types.ListAdd(value, parse.ConvertValue(v, itemType))
				}
			}
		}
		sliceOrListVal.Set(value.Elem())
		return
	}

	// map类型
	if mapType, isList := types.IsMap(sliceOrListVal); isList {
		// make.... var m map[k]v = make(map[k]v)
		sliceOrListVal.Set(reflect.MakeMap(mapType))
		for _, cmder := range receiver.cmder {
			switch cmderType := cmder.(type) {
			case *redis.StringStringMapCmd:
				jsonVal := cmderType.Val()
				for k, v := range jsonVal {
					sliceOrListVal.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
				}
			}
		}
		return
	}
	panic("sliceOrList入参必须为切片、collections.List类型、map类型")
}
