package postgres

import "reflect"

func StructValues(kv map[string]reflect.Value) []interface{} {
	values := make([]interface{}, 0, len(kv))
	for _, v := range kv {
		values = append(values, v.Interface())
	}
	return values
}

func StructFieldsPtr(kv map[string]reflect.Value) []interface{} {
	ptrs := make([]interface{}, 0, len(kv))
	for _, v := range kv {
		ptrs = append(ptrs, v.Addr().Interface())
	}
	return ptrs
}

func StructFields(kv map[string]reflect.Value) []string {
	fields := make([]string, 0, len(kv))
	for k := range kv {
		fields = append(fields, k)
	}
	return fields
}

func StructMap(s interface{}, expect bool, keys ...string) map[string]reflect.Value {
	v := reflect.ValueOf(s).Elem()
	var fields func(v reflect.Value) map[string]reflect.Value
	fields = func(v reflect.Value) map[string]reflect.Value {
		t := v.Type()
		kv := make(map[string]reflect.Value)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("db")
			if tag == "-" || tag == "" {
				continue
			}
			var keyExists bool
			for _, k := range keys {
				if tag == k {
					keyExists = true
				}
			}

			if (expect && keyExists) || (!expect && !keyExists) {
				continue
			}

			if tag == "inline" {
				for k, v := range fields(v.Field(i)) {
					kv[k] = v
				}
			} else {
				kv[tag] = v.Field(i)
			}
		}
		return kv
	}

	return fields(v)
}
