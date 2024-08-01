package gen

import (
	"math"
	"reflect"

	"github.com/leanovate/gopter"
)

// Float64Range generates float64 numbers within a given range
func Float64Range(min, max float64) gopter.Gen {
	d := max - min
	if d < 0 || d > math.MaxFloat64 {
		return Fail(reflect.TypeOf(float64(0)))
	}

	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		genResult := gopter.NewGenResult(min+genParams.Rng.Float64()*d, Float64Shrinker)
		genResult.Sieve = func(v interface{}) bool {
			return v.(float64) >= min && v.(float64) <= max
		}
		return genResult
	}
}

// Float64 generates arbitrary float64 numbers that do not contain NaN or Inf
func Float64() gopter.Gen {
	return gopter.CombineGens(
		Int64Range(0, 1),
		Int64Range(0, 0x7fe),
		Int64Range(0, 0xfffffffffffff),
	).Map(func(values []interface{}) float64 {
		sign := uint64(values[0].(int64))
		exponent := uint64(values[1].(int64))
		mantissa := uint64(values[2].(int64))

		return math.Float64frombits((sign << 63) | (exponent << 52) | mantissa)
	}).WithShrinker(Float64Shrinker)
}

// Float32Range generates float32 numbers within a given range
func Float32Range(min, max float32) gopter.Gen {
	d := max - min
	if d < 0 || d > math.MaxFloat32 {
		return Fail(reflect.TypeOf(float32(0)))
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		genResult := gopter.NewGenResult(min+genParams.Rng.Float32()*d, Float32Shrinker)
		genResult.Sieve = func(v interface{}) bool {
			return v.(float32) >= min && v.(float32) <= max
		}
		return genResult
	}
}

// Float32 generates arbitrary float32 numbers that do not contain NaN or Inf
func Float32() gopter.Gen {
	return gopter.CombineGens(
		Int32Range(0, 1),
		Int32Range(0, 0xfe),
		Int32Range(0, 0x7fffff),
	).Map(func(values []interface{}) float32 {
		sign := uint32(values[0].(int32))
		exponent := uint32(values[1].(int32))
		mantissa := uint32(values[2].(int32))

		return math.Float32frombits((sign << 31) | (exponent << 23) | mantissa)
	}).WithShrinker(Float32Shrinker)
}
