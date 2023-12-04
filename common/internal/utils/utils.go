package utils

import (
	"time"

	"github.com/jpillora/backoff"
	"golang.org/x/exp/constraints"
)

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}

}

// MinKey returns the minimum value of the given element array with respect
// to the given key function. In the event U is not a compound type (e.g a
// struct) an identity function can be provided.
func MinKey[U any, T constraints.Ordered](elems []U, key func(U) T) T {
	var min T
	if len(elems) == 0 {
		return min
	}

	min = key(elems[0])
	for i := 1; i < len(elems); i++ {
		v := key(elems[i])
		if v < min {
			min = v
		}
	}

	return min
}
