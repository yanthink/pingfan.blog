package helpers

import (
	"reflect"
)

func ElemValue(s any) (value reflect.Value) {
	value = reflect.ValueOf(s)

	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	return
}

func ToPointer[T any](value T) *T {
	return &value
}
