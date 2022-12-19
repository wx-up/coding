package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	pathParam map[string]string

	// cacheQueryParam 缓存请求参数
	// 这个缓存不存在所谓的失效不一致问题，对于 Web 框架来说，收到请求之后，值就是确切无疑的，不会再变
	cacheQueryParam url.Values

	// 匹配到的路由
	MatchPath string

	// 缓存响应，如果不缓存的话，用户调用 Resp.Write 写入的响应，后续是获取不到的
	// 比如场景：我需要用日志记录请求的响应，默认情况下实现不了
	// 解决方案：用户写入响应的时候，实际是写到内存，业务执行完成之后，再 flush
	RespStatusCode int
	RespData       []byte
}

var ErrBodyNotJsonType = errors.New("body 不是 json 格式")

// BindJSON 解析 body json 参数
func (c *Context) BindJSON(val any) error {
	if !strings.Contains(c.Req.Header.Get("Content-Type"), "application/json") {
		return ErrBodyNotJsonType
	}
	if c.Req.Body == nil {
		return errors.New("web：body 为空")
	}
	// 不需要自己从 c.Req.Body 中读取 bytes，NewDecoder 内部自己会读取
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

// FormValue 解析表单数据
func (c *Context) FormValue(key string) (string, error) {
	// 多次调用 ParseForm 并不会重复解析
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}

	return c.Req.FormValue(key), nil
}

// ParamValue 解析 query 参数
func (c *Context) ParamValue(key string) (string, error) {
	// URL.Query() 每次都会解析，所以需要自己实现缓存提高性能
	if c.cacheQueryParam == nil {
		c.cacheQueryParam = c.Req.URL.Query()
	}

	vs, ok := c.cacheQueryParam[key]
	if !ok || len(vs) <= 0 {
		return "", errors.New("key 不存在")
	}
	return vs[0], nil
}

// PathValue 解析路径参数
func (c *Context) PathValue(key string) (string, error) {
	if c.pathParam == nil {
		return "", errors.New("key 不存在")
	}
	v, ok := c.pathParam[key]
	if !ok {
		return "", errors.New("key 不存在")
	}
	return v, nil
}

// RespJSON 返回 json
func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(bs)
	return err
}

// SetCookie 设置 cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) PathValueV2(key string) StringValue {
	if c.pathParam == nil {
		return StringValue{err: errors.New("key 不存在")}
	}
	v, ok := c.pathParam[key]
	if !ok {
		return StringValue{err: errors.New("key 不存在")}
	}
	return StringValue{val: v}
}

type StringValue struct {
	val string
	err error
}

func (sv StringValue) ToInt64() (int64, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	return strconv.ParseInt(sv.val, 10, 64)
}
