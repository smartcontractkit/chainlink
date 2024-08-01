package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// MapOf generates an arbitrary map of generated kay values.
// genParams.MaxSize sets an (exclusive) upper limit on the size of the map
// genParams.MinSize sets an (inclusive) lower limit on the size of the map
func MapOf(keyGen, elementGen gopter.Gen) gopter.Gen {
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

		result, keySieve, keyShrinker, elementSieve, elementShrinker := genMap(keyGen, elementGen, genParams, len)

		genResult := gopter.NewGenResult(result.Interface(), MapShrinker(keyShrinker, elementShrinker))
		if keySieve != nil || elementSieve != nil {
			genResult.Sieve = forAllKeyValueSieve(keySieve, elementSieve)
		}
		return genResult
	}
}

func genMap(keyGen, elementGen gopter.Gen, genParams *gopter.GenParameters, len int) (reflect.Value, func(interface{}) bool, gopter.Shrinker, func(interface{}) bool, gopter.Shrinker) {
	element := elementGen(genParams)
	elementSieve := element.Sieve
	elementShrinker := element.Shrinker

	key := keyGen(genParams)
	keySieve := key.Sieve
	keyShrinker := key.Shrinker

	result := reflect.MakeMapWithSize(reflect.MapOf(key.ResultType, element.ResultType), len)

	for i := 0; i < len; i++ {
		keyValue, keyOk := key.Retrieve()
		elementValue, elementOk := element.Retrieve()

		if keyOk && elementOk {
			if key == nil {
				if elementValue == nil {
					result.SetMapIndex(reflect.Zero(key.ResultType), reflect.Zero(element.ResultType))
				} else {
					result.SetMapIndex(reflect.Zero(key.ResultType), reflect.ValueOf(elementValue))
				}
			} else {
				if elementValue == nil {
					result.SetMapIndex(reflect.ValueOf(keyValue), reflect.Zero(element.ResultType))
				} else {
					result.SetMapIndex(reflect.ValueOf(keyValue), reflect.ValueOf(elementValue))
				}
			}
		}
		key = keyGen(genParams)
		element = elementGen(genParams)
	}

	return result, keySieve, keyShrinker, elementSieve, elementShrinker
}

func forAllKeyValueSieve(keySieve, elementSieve func(interface{}) bool) func(interface{}) bool {
	return func(v interface{}) bool {
		rv := reflect.ValueOf(v)
		for _, key := range rv.MapKeys() {
			if keySieve != nil && !keySieve(key.Interface()) {
				return false
			}
			if elementSieve != nil && !elementSieve(rv.MapIndex(key).Interface()) {
				return false
			}
		}
		return true
	}
}
