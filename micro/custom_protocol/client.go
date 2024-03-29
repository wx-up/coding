package proxy_v2

import (
	"context"
	"errors"
	"net"
	"reflect"
	"time"

	"github.com/wx-up/coding/micro/custom_protocol/serialize/json"

	"github.com/wx-up/coding/micro/custom_protocol/serialize"

	"github.com/wx-up/coding/micro/custom_protocol/message"
)

/*
 有一个测试思路：有些方法比较长，依赖比较多，比如 InitService 依赖 proxy
 因此可以抽出一个子方法也就是 setFuncField 将 proxy 通过参数的方式传入
 这样子我们就可以为 setFuncField 编写单元测试
*/

func (c *Client) InitService(service Service) error {
	return setFuncField(service, c, c.serializer)
}

func setFuncField(service Service, p Proxy, s serialize.Serializer) error {
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

			reqData, _ := s.Encode(args[1].Interface())
			req := &message.Request{
				// 如果直接使用结构体名会有冲突的情况 有一种方案是使用 包名+结构体名
				// 这里则使用 Name 接口，由用户自定义
				ServiceName: service.Name(),
				MethodName:  fieldTyp.Name,
				Serializer:  uint8(s.Code()),
				Data:        reqData,
			}
			req.CalculateBodyLength()
			req.CalculateHeaderLength()

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

			// 处理业务错误
			if len(resp.Error) > 0 {
				return []reflect.Value{retVal, reflect.ValueOf(errors.New(string(resp.Error)))}
			}

			// 将结果反序列化
			if len(resp.Data) > 0 {
				err = s.Decode(resp.Data, retVal.Interface())
				if err != nil {
					return []reflect.Value{
						// 第一个返回值的类型：fieldTyp.Type.Out(0)
						retVal,
						reflect.ValueOf(err),
					}
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
	addr       string
	serializer serialize.Serializer
}

type Option func(*Client)

func WithSerializer(serializer serialize.Serializer) Option {
	return func(client *Client) {
		client.serializer = serializer
	}
}

func NewClient(addr string, opts ...Option) *Client {
	c := &Client{
		addr:       addr,
		serializer: json.New(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	data := message.EncodeReq(req)
	resp, err := c.SendData(data)
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(resp), nil
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

	_, err = conn.Write(bs)
	if err != nil {
		return nil, err
	}

	return Read(conn)
}
