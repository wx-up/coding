package orm

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
	"sync"
)

type Registry struct {
	models map[reflect.Type]*Model
	lock   sync.RWMutex
}

// get 锁保护，double-check
func (r *Registry) get(val any) (*Model, error) {
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
	fieldMap := make(map[string]*field, fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		f := typ.Field(i)
		fieldMap[f.Name] = &field{
			colName: strcase.ToSnake(f.Name),
		}
	}
	return &Model{
		tableName: strcase.ToSnake(pluralize.NewClient().Plural(typ.Name())),
		fieldMap:  fieldMap,
	}, nil
}
