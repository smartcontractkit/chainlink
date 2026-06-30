//go:build !race

// Package synctest is a wrapper around the testing/synctest package to handle a known bug in Go's race detection when dealing with `synctest.Test`.
package synctest

import (
	"testing"
	"testing/synctest"
	"time"
)

// Test wraps `testing/synctest.Test`.
// In !race builds we enable the synctest bubble; under -race this is a no-op via synctest_race.go
// due to https://github.com/golang/go/issues/76691 (fixed in Go 1.27).
func Test(t *testing.T, f func(*testing.T)) {
	synctest.Test(t, f)
}

// Wait wraps `testing/synctest.Wait` (no-op under -race; see synctest_race.go).
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
