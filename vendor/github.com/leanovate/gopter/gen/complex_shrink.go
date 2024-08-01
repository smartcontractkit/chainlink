package gen

import "github.com/leanovate/gopter"

// Complex128Shrinker is a shrinker for complex128 numbers
func Complex128Shrinker(v interface{}) gopter.Shrink {
	c := v.(complex128)
	realShrink := Float64Shrinker(real(c)).Map(func(r float64) complex128 {
		return complex(r, imag(c))
	})
	imagShrink := Float64Shrinker(imag(c)).Map(func(i float64) complex128 {
		return complex(real(c), i)
	})
	return realShrink.Interleave(imagShrink)
}

// Complex64Shrinker is a shrinker for complex64 numbers
func Complex64Shrinker(v interface{}) gopter.Shrink {
	c := v.(complex64)
	realShrink := Float64Shrinker(float64(real(c))).Map(func(r float64) complex64 {
		return complex(float32(r), imag(c))
	})
	imagShrink := Float64Shrinker(float64(imag(c))).Map(func(i float64) complex64 {
		return complex(real(c), float32(i))
	})
	return realShrink.Interleave(imagShrink)
}
