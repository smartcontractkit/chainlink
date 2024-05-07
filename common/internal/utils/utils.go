package utils

import (
	"cmp"
	"slices"
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

// MinFunc returns the minimum value of the given element array with respect
// to the given key function. In the event U is not a compound type (e.g a
// struct) an identity function can be provided.
func MinFunc[U any, T constraints.Ordered](elems []U, f func(U) T) T {
	var min T
	if len(elems) == 0 {
		return min
	}

	e := slices.MinFunc(elems, func(a, b U) int {
		return cmp.Compare(f(a), f(b))
	})
	return f(e)
}
