package reflect

import (
	"errors"
	"reflect"
)

// iterateFunc 只支持 结构体或者结构体的指针
func iterateFunc(val any) (map[string]*FuncInfo, error) {
	if val == nil {
		return nil, errors.New("val 不能为 nil")
	}

	refV := reflect.ValueOf(val)
	refT := refV.Type()

	// val只能为结构体或者为指向结构体的指针
	if !(refT.Kind() == reflect.Struct || (refT.Kind() == reflect.Pointer && refT.Elem().Kind() == reflect.Struct)) {
		return nil, errors.New("不支持的 val 类型")
	}

	methodNum := refT.NumMethod()
	res := make(map[string]*FuncInfo)
	for i := 0; i < methodNum; i++ {
		info := &FuncInfo{}
		method := refT.Method(i)

		// 入参数
		inNum := method.Type.NumIn()
		// 参数构建
		ps := make([]reflect.Value, 0, inNum)
		// 方法的第一个参数是接收者
		ps = append(ps, reflect.ValueOf(val))
		for j := 0; j < inNum; j++ {
			info.In = append(info.In, method.Type.In(j))

			// 如果方法还有其他参数，就设置为对应类型的零值
			if j > 0 {
				ps = append(ps, reflect.Zero(method.Type.In(j)))
			}
		}

		// 出参
		outNum := method.Type.NumOut()
		for j := 0; j < outNum; j++ {
			info.Out = append(info.Out, method.Type.Out(j))
		}

		// 函数调用
		resVal := method.Func.Call(ps)

		// 将结果转为 any 类型
		for _, v := range resVal {
			info.Result = append(info.Result, v.Interface())
		}

		info.Name = method.Name
		res[method.Name] = info
	}

	return res, nil
}

type FuncInfo struct {
	Name string
	In   []reflect.Type
	Out  []reflect.Type

	Result []any
}
