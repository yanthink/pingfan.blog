package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func StructToMap(obj any) map[string]any {
	var result = make(map[string]any)

	rv := reflect.ValueOf(obj)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		vField := rv.Field(i)
		tField := rt.Field(i)

		for vField.Kind() == reflect.Ptr && !vField.IsNil() {
			vField = vField.Elem()
		}

		tagParam := strings.Split(tField.Tag.Get("json"), ",")
		tagName := tagParam[0]
		tagOption := tagParam[1:]

		if !tField.IsExported() || tagName == "-" {
			continue
		}

		if vField.IsZero() && Contains(tagOption, "omitempty") {
			continue
		}

		if tagName == "" {
			tagName = tField.Name
		}

		if tField.Anonymous && vField.Kind() == reflect.Struct {
			result = Merge(result, StructToMap(vField.Interface()))
			continue
		}

		if vField.Kind() == reflect.Struct {
			switch v := vField.Interface().(type) {
			case time.Time:
				result[tagName] = v
			default:
				result[tagName] = StructToMap(vField.Interface())
			}
			continue
		}

		result[tagName] = vField.Interface()

		if Contains(tagOption, "string") {
			result[tagName] = fmt.Sprintf("%v", result[tagName])
		}
	}

	return result
}
