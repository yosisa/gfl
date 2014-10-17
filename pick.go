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
		subfs, ok := fs[name]
		if !ok {
			continue
		}
		if omitempty && isEmptyValue(rv.Field(i)) {
			continue
		}
		result[name] = pick(rv.Field(i), subfs)
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

// isEmptyValue returns true if given v is zero value.
// Almost all code of this function is borrowed from encoding/json package.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
