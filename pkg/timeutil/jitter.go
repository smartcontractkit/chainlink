package timeutil

import (
	mrand "math/rand"
	"time"
)

// JitterPct is a percent by which to scale a duration up or down.
// For example, 0.1 will result in +/- 10%.
type JitterPct float64

func (p JitterPct) Apply(d time.Duration) time.Duration {
	// #nosec
	if d == 0 {
		return 0
	}
	// ensure non-zero arg to Intn to avoid panic
	ub := max(1, int(float64(d.Abs())*float64(p)))
	// #nosec - non critical randomness
	jitter := mrand.Intn(2*ub) - ub
	return time.Duration(int(d) + jitter)
}
