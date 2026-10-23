package toml

import "reflect"

// deepCopy returns a copy of v that shares no pointers, slices, or maps with v.
// Config structs are trees of pointer fields, so a plain SetFrom would alias
// the source; writes through either side would then leak into the other.
func deepCopy[T any](v T) T {
	src := reflect.ValueOf(&v).Elem()
	dst := reflect.New(src.Type()).Elem()
	copyValue(dst, src)
	return dst.Interface().(T)
}

func copyValue(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Pointer:
		if src.IsNil() {
			return
		}
		p := reflect.New(src.Type().Elem())
		copyValue(p.Elem(), src.Elem())
		dst.Set(p)
	case reflect.Struct:
		if !allFieldsExported(src.Type()) {
			// Opaque value types (e.g. Duration, big.Int, url.Userinfo) cannot be
			// walked via reflection; copy them by value.
			dst.Set(src)
			return
		}
		for i := 0; i < src.NumField(); i++ {
			copyValue(dst.Field(i), src.Field(i))
		}
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		s := reflect.MakeSlice(src.Type(), src.Len(), src.Len())
		for i := 0; i < src.Len(); i++ {
			copyValue(s.Index(i), src.Index(i))
		}
		dst.Set(s)
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}
		m := reflect.MakeMapWithSize(src.Type(), src.Len())
		iter := src.MapRange()
		for iter.Next() {
			val := reflect.New(src.Type().Elem()).Elem()
			copyValue(val, iter.Value())
			m.SetMapIndex(iter.Key(), val)
		}
		dst.Set(m)
	case reflect.Interface:
		if src.IsNil() {
			return
		}
		val := reflect.New(src.Elem().Type()).Elem()
		copyValue(val, src.Elem())
		dst.Set(val)
	default:
		dst.Set(src)
	}
}

func allFieldsExported(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			return false
		}
	}
	return true
}
