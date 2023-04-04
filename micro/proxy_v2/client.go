package proxy_v2

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"time"
)

/*
 有一个测试思路：有些方法比较长，依赖比较多，比如 InitClientProxy 依赖 proxy
 因此可以抽出一个子方法也就是 setFuncField 将 proxy 通过参数的方式传入
 这样子我们就可以为 setFuncField 编写单元测试
*/

func InitClientProxy(addr string, service Service) error {
	p := NewClient(addr)
	return setFuncField(service, p)
}

func setFuncField(service Service, p Proxy) error {
	if service == nil {
		return errors.New("rpc：不支持 nil")
	}
	val := reflect.ValueOf(service)
	if !(val.Kind() == reflect.Pointer && val.Elem().Kind() == reflect.Struct) {
		return errors.New("rpc：只支持指向结构体的指针")
	}

	val = val.Elem()
	typ := val.Type()

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if fieldTyp.Type.Kind() != reflect.Func {
			continue
		}

		// 捕获本地调用，然后调用 set 方法篡改它，改成发起 RPC 调用
		funcVal := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
			ctx := args[0].Interface().(context.Context)

			reqData, _ := json.Marshal(args[1].Interface())
			req := &Request{
				// 如果直接使用结构体名会有冲突的情况 有一种方案是使用 包名+结构体名
				// 这里则使用 Name 接口，由用户自定义
				ServiceName: service.Name(),
				MethodName:  fieldTyp.Name,
				Arg:         reqData,
			}

			retVal := reflect.New(fieldTyp.Type.Out(0).Elem())

			// 正式发起请求
			resp, err := p.Invoke(ctx, req)
			if err != nil {
				return []reflect.Value{
					// 第一个返回值的类型：fieldTyp.Type.Out(0)
					retVal,
					reflect.ValueOf(err),
				}
			}

			// 将结果反序列化
			err = json.Unmarshal(resp.Data, retVal.Interface())
			if err != nil {
				return []reflect.Value{
					// 第一个返回值的类型：fieldTyp.Type.Out(0)
					retVal,
					reflect.ValueOf(err),
				}
			}

			return []reflect.Value{
				retVal,
				reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()),
			}
		})
		fieldVal.Set(funcVal)
	}

	return nil
}

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

const numOfLengthBytes = 8

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	// 这里先使用 json
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendData(data)
	if err != nil {
		return nil, err
	}
	return &Response{Data: resp}, err
}

// SendData 发送数据
// 两段读
func (c *Client) SendData(bs []byte) ([]byte, error) {
	// 发送数据
	conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}

	contentBs := make([]byte, numOfLengthBytes, numOfLengthBytes+len(bs))
	binary.BigEndian.PutUint64(contentBs, uint64(len(bs)))
	contentBs = append(contentBs, bs...)
	_, err = conn.Write(contentBs)
	if err != nil {
		return nil, err
	}

	// 读取响应：读取内容长度
	lenBs := make([]byte, numOfLengthBytes)
	_, err = conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	contentLen := binary.BigEndian.Uint64(lenBs)

	// 读取响应：读取内容
	contentBs = make([]byte, contentLen)
	_, err = conn.Read(contentBs)
	if err != nil {
		return nil, err
	}
	return contentBs, nil
}
