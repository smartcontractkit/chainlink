package gen

import (
	"time"

	"github.com/leanovate/gopter"
)

// TimeShrinker is a shrinker for time.Time structs
func TimeShrinker(v interface{}) gopter.Shrink {
	t := v.(time.Time)
	sec := t.Unix()
	nsec := int64(t.Nanosecond())
	secShrink := uint64Shrink{
		original: uint64(sec),
		half:     uint64(sec),
	}
	nsecShrink := uint64Shrink{
		original: uint64(nsec),
		half:     uint64(nsec),
	}
	return gopter.Shrink(secShrink.Next).Map(func(v uint64) time.Time {
		return time.Unix(int64(v), nsec)
	}).Interleave(gopter.Shrink(nsecShrink.Next).Map(func(v uint64) time.Time {
		return time.Unix(sec, int64(v))
	}))
}
