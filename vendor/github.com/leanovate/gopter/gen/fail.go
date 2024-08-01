package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// Fail is a generator that always fails to generate a value
// Useful as fallback
func Fail(resultType reflect.Type) gopter.Gen {
	return func(*gopter.GenParameters) *gopter.GenResult {
		return gopter.NewEmptyResult(resultType)
	}
}
