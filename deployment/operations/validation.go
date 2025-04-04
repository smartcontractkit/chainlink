package operations

import (
	"encoding/json"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// IsSerializable returns true if the value can be marshaled and unmarshaled without losing information, false otherwise.
// For idempotency and reporting purposes, we need to ensure that the value can be marshaled and unmarshaled
// without losing information.
// If the value implements json.Marshaler and json.Unmarshaler, it is assumed to be serializable.
func IsSerializable(lggr logger.Logger, v any) bool {
	if !isValueSerializable(lggr, reflect.ValueOf(v)) {
		return false
	}

	_, err := json.Marshal(v)
	if err != nil {
		lggr.Errorw("Failed to marshal value", "err", err)
		return false
	}

	return true
}

func isValueSerializable(lggr logger.Logger, v reflect.Value) bool {
	// Handle nil values
	if !v.IsValid() {
		return true
	}

	// Check if type implements json.Marshaler and json.Unmarshaler
	t := v.Type()
	marshalType := reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	unmarshalType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

	// If it implements both interfaces, assume it's serializable
	if t.Implements(marshalType) && t.Implements(unmarshalType) {
		return true
	}

	// Dereference pointers to get to the underlying value
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		return isValueSerializable(lggr, v.Elem())
	}

	// Rest of implementation for other types...
	switch v.Kind() {
	case reflect.Struct:
		// Check if each field is serializable
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)

			// Check for json:"-" tag and skip those fields
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			// Check if the field is exported (public)
			if !v.Type().Field(i).IsExported() {
				lggr.Errorw("Struct contains unexported field",
					"struct", v.Type().Name(),
					"field", v.Type().Field(i).Name)
				return false
			}
			if !isValueSerializable(lggr, field) {
				return false
			}
		}
		return true

	case reflect.Map:
		// Check all values in the map
		for _, key := range v.MapKeys() {
			if !isValueSerializable(lggr, v.MapIndex(key)) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.Array:
		// Check each element
		for i := 0; i < v.Len(); i++ {
			if !isValueSerializable(lggr, v.Index(i)) {
				return false
			}
		}
		return true

	default:
		// Basic types (string, int, float, etc.) are serializable
		return true
	}
}
