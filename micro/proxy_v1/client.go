package proxy_v1

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

// func InitClientProxy[T Service](service T) error {}

func InitClientProxy(service Service, proxy Proxy) error {
	refVal := reflect.ValueOf(service)
	refTyp := refVal.Type()
	if !(refTyp.Kind() == reflect.Pointer && refTyp.Elem().Kind() == reflect.Struct) {
		return errors.New("类型错误：必须是指向结构体的指针")
	}
	refVal = refVal.Elem()
	refTyp = refVal.Type()

	numField := refTyp.NumField()

	for i := 0; i < numField; i++ {
		// 字段类型
		fieldTyp := refTyp.Field(i)
		fieldVal := refVal.Field(i)

		// 不是函数类型就跳过
		if fieldTyp.Type.Kind() != reflect.Func {
			continue
		}

		if !fieldVal.CanSet() {
			continue
		}

		// 替换成一个新的方法实现
		newFn := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
			// 入参检测
			if err := checkIn(args); err != nil {
				panic(err)
			}
			// 返回值检测
			if err := checkReturn(fieldTyp.Type); err != nil {
				panic(err)
			}

			// 构建调用信息
			// 单元测试的话实际是断言 req
			// 这里需要一个探针，尝试着将 req 传递出去，方便外面断言
			// 这里可以采用 error 传递，也可以通过 context 传递
			req := &Request{
				// 客户端往往会命名一个和服务端服务名称不一致的结构体方法，比如叫做 UserClient 等等
				// 因此这里不能直接使用结构体的名称
				// 这里的解决方案：使用 接口
				ServerName: service.Name(),
				// 客户端和服务端也有可能存在方法名不一致的情况，正常来说，客户端叫 GetById 服务端也会叫 GetById
				// 不一致处理方法：引入 tag，像 orm 的字段和数据库字段映射那样处理
				MethodName: fieldTyp.Name,
				Arg:        args[1].Interface(),
			}

			// 第一个返回值的类型
			firstRetType := fieldTyp.Type.Out(0)

			// 将 req 发送出去，并获得响应
			// 这里需要一个抽象屏蔽发送的逻辑，因为你可以使用 TCP、也可以使用 http
			resp, err := proxy.Invoke(args[0].Interface().(context.Context), req)
			if err != nil {
				results = append(results, reflect.Zero(firstRetType))
				results = append(results, reflect.ValueOf(err))
				return
			}

			// 为第一个返回值分配内存
			// 注意：firstRetType 是一个指针类型， reflect.New 创建得到的是一个二级指针
			// 所以需要 Elem() 一下
			val := reflect.New(firstRetType.Elem()).Interface()

			// 需要将 resp 的结果填充到第一个返回值中
			// 假设这里采用的 json 序列化
			err = json.Unmarshal(resp.Data, val)

			results = append(results, reflect.ValueOf(val))

			if err != nil {
				results = append(results, reflect.ValueOf(err))
			} else {
				// 直接 reflect.ValueOf(nil) 是不行的，需要携带类型，记住就好
				results = append(results, reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()))
			}

			//outNum := fieldTyp.Type.NumOut()
			//for j := 0; j < outNum; j++ {
			//  填入类型的零值
			//	results = append(results, reflect.Zero(fieldTyp.Type.Out(j)))
			//}
			return results
		})
		fieldVal.Set(newFn)
	}
	return nil
}

// checkReturn 检测返回值
// 返回值，第一个参数为指向结构体的指针，第二个参数为 error
func checkReturn(refType reflect.Type) error {
	if refType.NumOut() != 2 {
		return errors.New("只能有两个返回值")
	}
	firstRet := refType.Out(0)
	if !(firstRet.Kind() == reflect.Pointer && firstRet.Elem().Kind() == reflect.Struct) {
		return errors.New("第一个参数必须是指向结构体的指针")
	}
	if !refType.Out(1).Implements(reflect.TypeOf(new(error)).Elem()) {
		return errors.New("第二个参数必须实现了 error 接口")
	}
	return nil
}

// checkIn 检测入参
// 入参，第一个参数为 context.Context 第二个参数为指向结构体的指针
func checkIn(args []reflect.Value) error {
	if len(args) != 2 {
		return errors.New("只能有两个参数")
	}
	_, ok := args[0].Interface().(context.Context)
	if !ok {
		return errors.New("第1个参数类型必须是 context.Context")
	}
	refTyp := args[1].Type()
	if !(refTyp.Kind() == reflect.Pointer && refTyp.Elem().Kind() == reflect.Struct) {
		return errors.New("第2个参数类型必须是指向结构体的指针")
	}
	return nil
}

type Service interface {
	Name() string
}
