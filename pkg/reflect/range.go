package reflect

import (
	"errors"
	"reflect"
)

// iterate 支持数组、切片、字符串的遍历
func iterate(val any) ([]any, error) {
	if val == nil {
		return nil, errors.New("val 不能为 nil")
	}
	refV := reflect.ValueOf(val)
	if !(refV.Kind() == reflect.Slice || refV.Kind() == reflect.Array || refV.Kind() == reflect.String) {
		return nil, errors.New("val 类型错误")
	}

	// 获取长度
	length := refV.Len()
	res := make([]any, 0, length)

	// 迭代
	for i := 0; i < length; i++ {
		ele := refV.Index(i)
		// 调用 Interface 获取具体的值
		res = append(res, ele.Interface())
	}
	return res, nil
}

// iterateMap 迭代map
func iterateMap(val any) ([]any, []any, error) {
	refV := reflect.ValueOf(val)
	refT := refV.Type()
	if refT.Kind() != reflect.Map {
		return nil, nil, errors.New("val 类型错误")
	}

	// 获取 key 的切片信息
	mapKeys := refV.MapKeys()
	keys, values := make([]any, 0, len(mapKeys)), make([]any, 0, len(mapKeys))

	// 遍历 key 获取 value
	for _, key := range mapKeys {
		keys = append(keys, key.Interface())
		values = append(values, refV.MapIndex(key).Interface())
	}
	return keys, values, nil
}

// iterateMapAnotherMethod 迭代 map 的另一种方式（ 迭代器 ）
func iterateMapAnotherMethod(val any) ([]any, []any, error) {
	return nil, nil, nil
}
