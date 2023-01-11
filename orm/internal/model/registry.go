package model

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
)

// RegistryInterface 编程注册的方式实现自定义，通过 ModelOpt 来实现自定义
// 本案例采用 标签+接口 的方式实现自定义
type RegistryInterface interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...Opt) (*Model, error)
}

// TagSplit 标签元素的分隔符
const TagSplit = ","

type Registry struct {
	models map[reflect.Type]*Model
	lock   sync.RWMutex
}

func NewRegister() *Registry {
	return &Registry{
		models: make(map[reflect.Type]*Model),
	}
}

// Get 锁保护，double-check
func (r *Registry) Get(val any) (*Model, error) {
	key := reflect.TypeOf(val)
	r.lock.RLock()
	model, ok := r.models[key]
	r.lock.RUnlock()
	if ok {
		return model, nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	model, ok = r.models[key]
	if ok {
		return model, nil
	}

	model, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	if r.models == nil {
		r.models = make(map[reflect.Type]*Model)
	}
	r.models[key] = model
	return model, nil
}

// parseModel 按照目前的设计其实不需要 error
// 但是根据经验随着功能迭代后续会需要 error，那时候再修改调用方会比较麻烦，所以现在就直接添加上
func (r *Registry) parseModel(val any) (*Model, error) {
	if val == nil {
		return nil, errs.ErrParseModelValType
	}

	typ := reflect.TypeOf(val)

	// 处理多级指针的情况
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrParseModelValType
	}
	fieldCnt := typ.NumField()
	fieldMap := make(map[string]*Field, fieldCnt)
	columnMap := make(map[string]*Field, fieldCnt)
	var columns = make([]*Field, 0, fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		f := typ.Field(i)
		ormTagValues := parseTag(f.Tag)

		// 列名，如果指定了则取指定的，没有指定则取字段名的下划线形式
		colName, ok := ormTagValues["column"]
		if !ok || colName == "" {
			colName = strcase.ToSnake(f.Name)
		}

		fdMeta := &Field{
			ColName: colName,
			Typ:     f.Type,
			Name:    f.Name,
			Offset:  f.Offset,
			Index:   f.Index,
		}

		fieldMap[f.Name] = fdMeta
		columnMap[colName] = fdMeta
		columns = append(columns, fdMeta)

	}
	return &Model{
		TableName: parseTableName(val, typ),
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
		Columns:   columns,
	}, nil
}

// parseTableName 解析表名，需要判断是否实现了 TableName 接口
func parseTableName(val any, typ reflect.Type) string {
	var tableName string

	// 查看结构体是否实现 TableName 接口，实现的话直接调用获取表名
	if v, ok := val.(TableName); ok {
		tableName = v.TableName()
	}
	if tableName == "" {
		tableName = strcase.ToSnake(pluralize.NewClient().Plural(typ.Name()))
	}
	return tableName
}

// parseTag 解析 orm 标签
// orm:"column=id,primary_key"
func parseTag(tag reflect.StructTag) map[string]string {
	tagString := tag.Get("orm")
	kvs := strings.Split(tagString, TagSplit)
	res := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		segments := strings.SplitN(kv, "=", 2)
		v := ""
		if len(segments) > 1 {
			v = segments[1]
		}
		res[segments[0]] = v
	}
	return res
}
