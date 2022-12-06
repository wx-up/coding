package unsafe

import (
	"fmt"
	"reflect"
)

// PrintFieldOffset 打印字段地址相对于起始地址的偏移量
func PrintFieldOffset(entity any) {
	res, err := printFieldOffset(entity)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range res {
		fmt.Printf("字段名：%s，偏移量：%d\n", k, v)
	}
}

func printFieldOffset(entity any) (map[string]uintptr, error) {
	relV := reflect.ValueOf(entity)
	relT := relV.Type()
	filedNum := relT.NumField()
	res := make(map[string]uintptr)
	for i := 0; i < filedNum; i++ {
		field := relT.Field(i)
		res[field.Name] = field.Offset
	}
	return res, nil
}
