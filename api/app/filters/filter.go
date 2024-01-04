package filters

import (
	"gorm.io/gorm"
	"reflect"
)

func New(filter any, data any) func(db *gorm.DB) *gorm.DB {
	filterValue := reflect.ValueOf(filter)
	dataValue := reflect.ValueOf(data).Elem()
	dataType := dataValue.Type()

	// 获取 filter 结构体中所有的方法并存储在 map 中
	methods := make(map[string]reflect.Value)
	for i := 0; i < filterValue.NumMethod(); i++ {
		method := filterValue.Type().Method(i)
		methods[method.Name] = filterValue.Method(i)
	}

	return func(db *gorm.DB) *gorm.DB {
		dbValue := reflect.ValueOf(db)

		for i := 0; i < dataValue.NumField(); i++ {
			if tag := dataType.Field(i).Tag.Get("filter"); tag == "-" {
				continue
			}

			fieldName := dataType.Field(i).Name

			if fieldValue := dataValue.Field(i); fieldValue.IsValid() && !fieldValue.IsZero() {
				if method, ok := methods[fieldName]; ok {
					args := []reflect.Value{
						dbValue,
						fieldValue,
					}

					dbValue = method.Call(args)[0]
				}
			}
		}

		return dbValue.Interface().(*gorm.DB)
	}
}
