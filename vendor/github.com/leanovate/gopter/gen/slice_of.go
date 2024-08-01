package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// SliceOf generates an arbitrary slice of generated elements
// genParams.MaxSize sets an (exclusive) upper limit on the size of the slice
// genParams.MinSize sets an (inclusive) lower limit on the size of the slice
func SliceOf(elementGen gopter.Gen, typeOverrides ...reflect.Type) gopter.Gen {
	var typeOverride reflect.Type
	if len(typeOverrides) > 1 {
		panic("too many type overrides specified, at most 1 may be provided.")
	} else if len(typeOverrides) == 1 {
		typeOverride = typeOverrides[0]
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		len := 0
		if genParams.MaxSize > 0 || genParams.MinSize > 0 {
			if genParams.MinSize > genParams.MaxSize {
				panic("GenParameters.MinSize must be <= GenParameters.MaxSize")
			}

			if genParams.MaxSize == genParams.MinSize {
				len = genParams.MaxSize
			} else {
				len = genParams.Rng.Intn(genParams.MaxSize-genParams.MinSize) + genParams.MinSize
			}
		}
		result, elementSieve, elementShrinker := genSlice(elementGen, genParams, len, typeOverride)

		genResult := gopter.NewGenResult(result.Interface(), SliceShrinker(elementShrinker))
		if elementSieve != nil {
			genResult.Sieve = forAllSieve(elementSieve)
		}
		return genResult
	}
}

// SliceOfN generates a slice of generated elements with definied length
func SliceOfN(desiredlen int, elementGen gopter.Gen, typeOverrides ...reflect.Type) gopter.Gen {
	var typeOverride reflect.Type
	if len(typeOverrides) > 1 {
		panic("too many type overrides specified, at most 1 may be provided.")
	} else if len(typeOverrides) == 1 {
		typeOverride = typeOverrides[0]
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		result, elementSieve, elementShrinker := genSlice(elementGen, genParams, desiredlen, typeOverride)

		genResult := gopter.NewGenResult(result.Interface(), SliceShrinkerOne(elementShrinker))
		if elementSieve != nil {
			genResult.Sieve = func(v interface{}) bool {
				rv := reflect.ValueOf(v)
				return rv.Len() == desiredlen && forAllSieve(elementSieve)(v)
			}
		} else {
			genResult.Sieve = func(v interface{}) bool {
				return reflect.ValueOf(v).Len() == desiredlen
			}
		}
		return genResult
	}
}

func genSlice(elementGen gopter.Gen, genParams *gopter.GenParameters, desiredlen int, typeOverride reflect.Type) (reflect.Value, func(interface{}) bool, gopter.Shrinker) {
	element := elementGen(genParams)
	elementSieve := element.Sieve
	elementShrinker := element.Shrinker

	sliceType := typeOverride
	if sliceType == nil {
		sliceType = element.ResultType
	}

	result := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, desiredlen)

	for i := 0; i < desiredlen; i++ {
		value, ok := element.Retrieve()

		if ok {
			if value == nil {
				result = reflect.Append(result, reflect.Zero(sliceType))
			} else {
				result = reflect.Append(result, reflect.ValueOf(value))
			}
		}
		element = elementGen(genParams)
	}

	return result, elementSieve, elementShrinker
}

func forAllSieve(elementSieve func(interface{}) bool) func(interface{}) bool {
	return func(v interface{}) bool {
		rv := reflect.ValueOf(v)
		for i := rv.Len() - 1; i >= 0; i-- {
			if !elementSieve(rv.Index(i).Interface()) {
				return false
			}
		}
		return true
	}
}
