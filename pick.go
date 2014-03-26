package gofl

import (
	"reflect"
	"strings"
)

type fieldSet map[string]fieldSet

func newFieldSet(fields ...string) fieldSet {
	if len(fields) == 0 {
		return nil
	}

	fieldList := make(map[string][]string)
	for _, field := range fields {
		parts := strings.SplitN(field, ".", 2)
		if len(parts) == 2 {
			fieldList[parts[0]] = append(fieldList[parts[0]], parts[1])
		} else if _, ok := fieldList[parts[0]]; !ok {
			fieldList[parts[0]] = nil
		}
	}

	fs := make(fieldSet)
	for name, subFields := range fieldList {
		if subFields == nil {
			fs[name] = nil
		} else {
			fs[name] = newFieldSet(subFields...)
		}
	}
	return fs
}

func Pick(v interface{}, fields ...string) interface{} {
	rv := reflect.ValueOf(v)
	fs := newFieldSet(fields...)
	return pick(rv, fs)
}

func pick(v reflect.Value, fs fieldSet) interface{} {
	if fs == nil {
		return v.Interface()
	}

	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Struct:
		return pickStruct(v, fs)
	case reflect.Slice:
		return pickSlice(v, fs)
	default:
		return v.Interface()
	}
}

func pickStruct(rv reflect.Value, fs fieldSet) interface{} {
	rt := rv.Type()
	result := make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous {
			base := pickStruct(rv.Field(i), fs)
			for key, val := range base.(map[string]interface{}) {
				if _, ok := result[key]; !ok {
					result[key] = val
				}
			}
			continue
		}

		name, omitempty := getTag(field)
		if name == "" {
			name = field.Name
		}
		if omitempty && isEmpty(rv.Field(i)) {
			continue
		}

		if subfs, ok := fs[name]; ok {
			result[name] = pick(rv.Field(i), subfs)
		}
	}
	return result
}

func pickSlice(v reflect.Value, fs fieldSet) []interface{} {
	size := v.Len()
	result := make([]interface{}, size)
	for i := 0; i < size; i++ {
		result[i] = pick(v.Index(i), fs)
	}
	return result
}

func getTag(field reflect.StructField) (name string, omitempty bool) {
	items := strings.Split(field.Tag.Get("json"), ",")
	name = items[0]
	for _, item := range items[1:] {
		if item == "omitempty" {
			omitempty = true
		}
	}
	return
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice:
		return v.IsNil()
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}
