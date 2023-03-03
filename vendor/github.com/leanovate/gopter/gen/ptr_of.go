package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// PtrOf generates either a pointer to a generated element or a nil pointer
func PtrOf(elementGen gopter.Gen) gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		element := elementGen(genParams)
		elementShrinker := element.Shrinker
		elementSieve := element.Sieve
		value, ok := element.Retrieve()
		if !ok || genParams.NextBool() {
			result := gopter.NewEmptyResult(reflect.PtrTo(element.ResultType))
			result.Sieve = func(v interface{}) bool {
				if elementSieve == nil {
					return true
				}
				r := reflect.ValueOf(v)
				return !r.IsValid() || r.IsNil() || elementSieve(r.Elem().Interface())
			}
			return result
		}
		// To get the right pointer type we have to create a slice with one element
		slice := reflect.MakeSlice(reflect.SliceOf(element.ResultType), 0, 1)
		slice = reflect.Append(slice, reflect.ValueOf(value))

		result := gopter.NewGenResult(slice.Index(0).Addr().Interface(), PtrShrinker(elementShrinker))
		result.Sieve = func(v interface{}) bool {
			if elementSieve == nil {
				return true
			}
			r := reflect.ValueOf(v)
			return !r.IsValid() || r.IsNil() || elementSieve(r.Elem().Interface())
		}
		return result
	}
}
