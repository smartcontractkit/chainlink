package gen

import (
	"math"

	"github.com/leanovate/gopter"
)

type float64Shrink struct {
	original float64
	half     float64
}

func (s *float64Shrink) isZeroOrVeryClose() bool {
	if s.half == 0 {
		return true
	}
	muliple := s.half * 100000
	return math.Abs(muliple) < 1 && muliple != 0
}

func (s *float64Shrink) Next() (interface{}, bool) {
	if s.isZeroOrVeryClose() {
		return nil, false
	}
	value := s.original - s.half
	s.half /= 2
	return value, true
}

// Float64Shrinker is a shrinker for float64 numbers
func Float64Shrinker(v interface{}) gopter.Shrink {
	negShrink := float64Shrink{
		original: -v.(float64),
		half:     -v.(float64),
	}
	posShrink := float64Shrink{
		original: v.(float64),
		half:     v.(float64) / 2,
	}
	return gopter.Shrink(negShrink.Next).Interleave(gopter.Shrink(posShrink.Next))
}

// Float32Shrinker is a shrinker for float32 numbers
func Float32Shrinker(v interface{}) gopter.Shrink {
	return Float64Shrinker(float64(v.(float32))).Map(func(e float64) float32 {
		return float32(e)
	})
}
