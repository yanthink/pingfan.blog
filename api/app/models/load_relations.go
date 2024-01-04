package models

import (
	"blog/app/helpers"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// IncludeParser 自动解析 query include
type IncludeParser interface {
	ParseInclude(relations any)
}

// Loader 延时加载关系接口
type Loader interface {
	Load(relations map[string][]any, missing ...bool)
}

// ParseRelations 解析关系名
// Orders.Goods -> [Orders, Orders.Goods]
func ParseRelations(relations []string) (parsed []string) {
	parsed = make([]string, 0)

	for _, preload := range relations {
		for _, s := range strings.Split(preload, ",") {
			nested := strings.Split(s, ".")

			part := nested[0]
			parsed = append(parsed, part)

			for i := 1; i < len(nested); i++ {
				part += "." + nested[i]
				parsed = append(parsed, part)
			}
		}
	}

	parsed = helpers.Unique(parsed)

	return
}

// FilterRelations 获取 relations 和 availableRelations 的交集
func FilterRelations(relations any, availableRelations any) (result map[string][]any) {
	result = make(map[string][]any, 0)

	preloadsKeys := ParseRelations(RelationsToKeys(relations))

	availablePreloadsMap := RelationsToMap(availableRelations)
	availablePreloadsKeys := ParseRelations(RelationsToKeys(availableRelations))

	for _, key := range helpers.Intersect(availablePreloadsKeys, preloadsKeys) {
		result[key] = nil

		if v, ok := availablePreloadsMap[key]; ok {
			result[key] = v
		}
	}

	return
}

// RelationsToMap 转换成 gorm 的预加载格式
func RelationsToMap(relations any) (result map[string][]any) {
	result = make(map[string][]any, 0)

	switch v := relations.(type) {
	case string:
		for _, s := range strings.Split(v, ",") {
			result[s] = nil
		}
	case []string:
		for _, s := range v {
			result[s] = nil
		}
	case map[string][]any:
		result = v
	}

	return
}

// RelationsToKeys 转换成关系名称 []string 格式
func RelationsToKeys(preloads any) (keys []string) {
	switch v := preloads.(type) {
	case string:
		keys = strings.Split(v, ",")
	case []string:
		keys = v
	case map[string][]any:
		keys = helpers.Keys(v)
	}

	keys = helpers.Map(keys, func(_ int, item string) string {
		return strings.SplitN(item, "=", 2)[0]
	})

	return
}

// RelationsToParams 获取 relations map[string]string 参数
func RelationsToParams(preloads any) (result map[string]string) {
	result = map[string]string{}

	var keys []string

	switch v := preloads.(type) {
	case string:
		keys = strings.Split(v, ",")
	case []string:
		keys = v
	case map[string][]any:
		keys = helpers.Keys(v)
	}

	for _, key := range keys {
		var params string
		values := strings.SplitN(key, "=", 2)

		if len(values) == 2 {
			params = values[1]
		}

		result[values[0]] = params
	}

	return
}

// LoadRelations 延时加载关系
func LoadRelations(db *gorm.DB, model any, preloads map[string][]any, missing ...bool) {
	if len(preloads) < 1 {
		return
	}

	modelsSlice := helpers.MakeSlice(model).Elem()

	if modelsSlice.Len() < 1 {
		return
	}

	dest := modelsSlice.Addr().Interface()

	tx := db.Model(dest)
	stmt := tx.Statement
	stmt.Dest = dest

	if err := stmt.Parse(stmt.Dest); err != nil {
		return
	}

	loaded := make(map[string]bool, 0)

	for name, value := range preloads {
		if loaded[name] {
			continue
		}
		loaded[name] = true

		preloadFields := strings.Split(name, ".")
		field := preloadFields[0]

		if rel := stmt.Schema.Relationships.Relations[field]; rel == nil {
			break
		}

		if len(missing) < 1 || missing[0] {
			var loadedModels []any

			for i := 0; i < modelsSlice.Len(); i++ {
				modelValue := modelsSlice.Index(i).Elem()
				relValue := modelValue.FieldByName(field)

				if !relValue.IsValid() || relValue.IsZero() {
					break
				}

				switch relValue.Kind() {
				case reflect.Slice, reflect.Array:
					for i := 0; i < relValue.Len(); i++ {
						relModelValue := relValue.Index(i)
						if relModelValue.Kind() != reflect.Ptr {
							relModelValue = relModelValue.Addr()
						}
						loadedModels = append(loadedModels, relModelValue.Interface())
					}
				default:
					if relValue.Kind() != reflect.Ptr {
						relValue = relValue.Addr()
					}

					loadedModels = append(loadedModels, relValue.Interface())
				}
			}

			if len(loadedModels) > 0 {
				if len(preloadFields) > 1 {
					LoadRelations(db, loadedModels, map[string][]any{strings.Join(preloadFields[1:], "."): value}, missing...)
				}
				continue
			}
		}

		tx = tx.Preload(name, value...)
	}

	stmt.ReflectValue = modelsSlice
	preload := tx.Callback().Query().Get("gorm:preload")
	preload(tx)
}
