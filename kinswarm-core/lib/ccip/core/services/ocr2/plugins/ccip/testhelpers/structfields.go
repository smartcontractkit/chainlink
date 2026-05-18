package testhelpers

import (
	"fmt"
	"reflect"
	"strings"
)

// FindStructFieldsOfCertainType recursively iterates over struct fields and returns all the fields of the provided type.
func FindStructFieldsOfCertainType(targetType string, v any) []string {
	typesAndFields := TypesAndFields("", reflect.ValueOf(v))
	results := make([]string, 0)
	for _, field := range typesAndFields {
		if strings.Contains(field, targetType) {
			results = append(results, field)
		}
	}
	return results
}

// TypesAndFields will find and return all the fields and their types of the provided value.
// NOTE: This is not intended for production use, it's a helper method for tests.
func TypesAndFields(prefix string, v reflect.Value) []string {
	results := make([]string, 0)

	s := v
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		typeAndName := fmt.Sprintf("%s%s %v", prefix, f.Type(), typeOfT.Field(i).Name)
		results = append(results, typeAndName)

		if f.Kind().String() == "ptr" {
			results = append(results, TypesAndFields(typeOfT.Field(i).Name, f.Elem())...)
		}

		if f.Kind().String() == "struct" {
			x1 := reflect.ValueOf(f.Interface())
			results = append(results, TypesAndFields(typeOfT.Field(i).Name, x1)...)
		}
	}

	return results
}
