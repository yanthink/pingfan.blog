package helpers

import (
	"fmt"
	"math/rand"
	"reflect"
)

// MakeSlice 将变量转换成具体指针类型的切片
// []any{*models.User, *models.User} -> *[]*models.User{*models.User, *models.User}
// []any{models.User, models.User} -> *[]*models.User{*models.User, *models.User}
// *models.User -> *[]*models.User{*models.User}
// models.User -> *[]*models.User{*models.User}
func MakeSlice(slice any) reflect.Value {
	var newSlice reflect.Value

	sliceValue := ElemValue(&slice)

	switch sliceValue.Kind() {
	case reflect.Slice, reflect.Array:
		size := sliceValue.Len()

		if size < 1 {
			if !sliceValue.CanAddr() {
				ptr := reflect.New(sliceValue.Type())
				ptr.Elem().Set(sliceValue)
				sliceValue = ptr.Elem()
			}

			return sliceValue.Addr()
		}

		itemValue := ElemValue(sliceValue.Index(0).Addr().Interface())
		sliceType := reflect.PtrTo(itemValue.Type())

		newSlice = reflect.MakeSlice(reflect.SliceOf(sliceType), size, size)

		for i := 0; i < size; i++ {
			itemValue = ElemValue(sliceValue.Index(i).Addr().Interface())
			newSlice.Index(i).Set(itemValue.Addr().Convert(sliceType))
		}
	default:
		sliceType := reflect.PtrTo(sliceValue.Type())
		newSlice = reflect.MakeSlice(reflect.SliceOf(sliceType), 1, 1)

		if !sliceValue.CanAddr() {
			ptr := reflect.New(sliceValue.Type())
			ptr.Elem().Set(sliceValue)
			sliceValue = ptr.Elem()
		}

		newSlice.Index(0).Set(sliceValue.Addr().Convert(sliceType))
	}

	results := reflect.New(newSlice.Type())
	results.Elem().Set(newSlice)

	return results
}

// Slice 将任意类型切片转化成具体类型切片
func Slice[T any](slice any) (result []T) {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Invalid slice type, value of type %T", slice))
	}

	if sliceValue.IsNil() {
		return nil
	}

	result = make([]T, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		result[i] = sliceValue.Index(i).Interface().(T)
	}

	return
}

// RandomByWeight 根据 weight 获取随机一项
func RandomByWeight[T comparable](items []T, weightFn func(item T) int64) (result T, index int) {
	var totalWeight int64

	for _, item := range items {
		totalWeight += weightFn(item)
	}

	// 生成 1 和 totalWeight 之间的伪随机int64类型整数
	r := rand.Int63n(totalWeight) + 1

	for index, result = range items {
		if r <= weightFn(result) {
			return
		}

		r -= weightFn(result)
	}

	index = rand.Intn(len(items))
	result = items[index]

	return
}

// Unique 去重
func Unique[T comparable](slice []T) []T {
	visited := make(map[T]bool)
	var result []T

	for _, v := range slice {
		if !visited[v] {
			visited[v] = true
			result = append(result, v)
		}
	}

	return result
}

// Intersect 交集
func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	var out []T
	result := map[T]bool{}

	for _, val := range slices[0] {
		result[val] = true
	}

	for i := 1; i < len(slices); i++ {
		out = []T{}
		nextResult := map[T]bool{}

		for _, val := range slices[i] {
			if result[val] {
				nextResult[val] = true
				out = append(out, val)
			}
		}

		result = nextResult
	}

	return out
}

// Map 遍历返回一个新切片
func Map[T any, U any](slice []T, iteratee func(index int, item T) U) []U {
	result := make([]U, len(slice), cap(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = iteratee(i, slice[i])
	}

	return result
}

func KeyBy[T any, U comparable](slice []T, iteratee func(item T) U) map[U]T {
	result := make(map[U]T, len(slice))

	for _, v := range slice {
		k := iteratee(v)
		result[k] = v
	}

	return result
}

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

func ContainsBy[T any](slice []T, predicate func(item T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

func Filter[T any](slice []T, predicate func(index int, item T) bool) []T {
	result := make([]T, 0)

	for i, v := range slice {
		if predicate(i, v) {
			result = append(result, v)
		}
	}

	return result
}
