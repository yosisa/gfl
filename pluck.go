package gfl

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

func Pluck(v interface{}, fields ...string) interface{} {
	rv := reflect.ValueOf(v)
	fs := newFieldSet(fields...)
	return pluck(rv, fs)
}

func pluck(v reflect.Value, fs fieldSet) interface{} {
	if fs == nil {
		return v.Interface()
	}

	switch v.Kind() {
	case reflect.Struct:
		return pluckStruct(v, fs)
	case reflect.Slice:
		return pluckSlice(v, fs)
	default:
		return v.Interface()
	}
}

func pluckStruct(rv reflect.Value, fs fieldSet) interface{} {
	rt := rv.Type()
	result := make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		name, omitempty := getTag(field)
		if name == "" {
			name = field.Name
		}
		if omitempty && isEmpty(rv.Field(i)) {
			continue
		}

		if subfs, ok := fs[name]; ok {
			result[name] = pluck(rv.Field(i), subfs)
		}
	}
	return result
}

func pluckSlice(v reflect.Value, fs fieldSet) []interface{} {
	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		result = append(result, pluck(v.Index(i), fs))
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
