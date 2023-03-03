package gen

import (
	"github.com/leanovate/gopter"
)

type int64Shrink struct {
	original int64
	half     int64
}

func (s *int64Shrink) Next() (interface{}, bool) {
	if s.half == 0 {
		return nil, false
	}
	value := s.original - s.half
	s.half /= 2
	return value, true
}

type uint64Shrink struct {
	original uint64
	half     uint64
}

func (s *uint64Shrink) Next() (interface{}, bool) {
	if s.half == 0 {
		return nil, false
	}
	value := s.original - s.half
	s.half >>= 1
	return value, true
}

// Int64Shrinker is a shrinker for int64 numbers
func Int64Shrinker(v interface{}) gopter.Shrink {
	negShrink := int64Shrink{
		original: -v.(int64),
		half:     -v.(int64),
	}
	posShrink := int64Shrink{
		original: v.(int64),
		half:     v.(int64) / 2,
	}
	return gopter.Shrink(negShrink.Next).Interleave(gopter.Shrink(posShrink.Next))
}

// UInt64Shrinker is a shrinker for uint64 numbers
func UInt64Shrinker(v interface{}) gopter.Shrink {
	shrink := uint64Shrink{
		original: v.(uint64),
		half:     v.(uint64),
	}
	return shrink.Next
}

// Int32Shrinker is a shrinker for int32 numbers
func Int32Shrinker(v interface{}) gopter.Shrink {
	return Int64Shrinker(int64(v.(int32))).Map(int64To32)
}

// UInt32Shrinker is a shrinker for uint32 numbers
func UInt32Shrinker(v interface{}) gopter.Shrink {
	return UInt64Shrinker(uint64(v.(uint32))).Map(uint64To32)
}

// Int16Shrinker is a shrinker for int16 numbers
func Int16Shrinker(v interface{}) gopter.Shrink {
	return Int64Shrinker(int64(v.(int16))).Map(int64To16)
}

// UInt16Shrinker is a shrinker for uint16 numbers
func UInt16Shrinker(v interface{}) gopter.Shrink {
	return UInt64Shrinker(uint64(v.(uint16))).Map(uint64To16)
}

// Int8Shrinker is a shrinker for int8 numbers
func Int8Shrinker(v interface{}) gopter.Shrink {
	return Int64Shrinker(int64(v.(int8))).Map(int64To8)
}

// UInt8Shrinker is a shrinker for uint8 numbers
func UInt8Shrinker(v interface{}) gopter.Shrink {
	return UInt64Shrinker(uint64(v.(uint8))).Map(uint64To8)
}

// IntShrinker is a shrinker for int numbers
func IntShrinker(v interface{}) gopter.Shrink {
	return Int64Shrinker(int64(v.(int))).Map(int64ToInt)
}

// UIntShrinker is a shrinker for uint numbers
func UIntShrinker(v interface{}) gopter.Shrink {
	return UInt64Shrinker(uint64(v.(uint))).Map(uint64ToUint)
}
