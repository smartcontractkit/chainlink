//go:build !race

// Package synctest is a wrapper around the testing/synctest package to handle a known bug in Go's race detection when dealing with `synctest.Test`.
package synctest

import (
	"testing"
	"testing/synctest"
	"time"
)

// Test wraps `testing/synctest.Test` to disable the synctest bubble when `-race` is enabled.
// There is a known issue with `synctest` when run with the `-race` detector that causes false positives.
// This should be fixed in Go 1.27 and we can remove this package.
// https://github.com/golang/go/issues/76691
func Test(t *testing.T, f func(*testing.T)) {
	synctest.Test(t, f)
}

// Wait wraps `testing/synctest.Wait` to disable the synctest bubble when `-race` is enabled.
func Wait() {
	synctest.Wait()
}

// SlowDelay returns a slow time duration for use in synctest tests.
func SlowDelay() time.Duration {
	return 1 * time.Minute
}

// FastDelay returns a fast time duration for use in synctest tests.
func FastDelay() time.Duration {
	return 1 * time.Second
}
