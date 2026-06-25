//go:build race

package synctest

import (
	"testing"
	"time"
)

// Test wraps `testing/synctest.Test` to disable the synctest bubble when `-race` is enabled.
// There is a known issue with `synctest` when run with the `-race` detector that causes false positives.
// This should be fixed in Go 1.27 and we can remove this package.
// https://github.com/golang/go/issues/76691
func Test(t *testing.T, f func(*testing.T)) {
	f(t)
}

// Wait wraps `testing/synctest.Wait` to disable the synctest bubble when `-race` is enabled.
func Wait() {
}

// SlowDelay returns a slow time duration for use in synctest tests.
func SlowDelay() time.Duration {
	return 500 * time.Millisecond
}

// FastDelay returns a fast time duration for use in synctest tests.
func FastDelay() time.Duration {
	return 10 * time.Millisecond
}
