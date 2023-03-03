package gen

import (
	"time"

	"github.com/leanovate/gopter"
)

// Time generates an arbitrary time.Time within year [0, 9999]
func Time() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		sec := genParams.Rng.Int63n(253402214400) // Ensure year in [0, 9999]
		usec := genParams.Rng.Int63n(1000000000)

		return gopter.NewGenResult(time.Unix(sec, usec), TimeShrinker)
	}
}

// AnyTime generates an arbitrary time.Time struct (might be way out of bounds of any reason)
func AnyTime() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		sec := genParams.NextInt64()
		usec := genParams.NextInt64()

		return gopter.NewGenResult(time.Unix(sec, usec), TimeShrinker)
	}
}

// TimeRange generates an arbitrary time.Time with a range
// from defines the start of the time range
// duration defines the overall duration of the time range
func TimeRange(from time.Time, duration time.Duration) gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		v := from.Add(time.Duration(genParams.Rng.Int63n(int64(duration))))
		return gopter.NewGenResult(v, TimeShrinker)
	}
}
