package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Context(t *testing.T) context.Context {
	ctx := context.Background()
	var cancel func()

	if d, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, d)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	t.Cleanup(cancel)
	return ctx
}

// DefaultWaitTimeout is the default wait timeout. If you have a *testing.T, use WaitTimeout instead.
const DefaultWaitTimeout = 30 * time.Second

// WaitTimeout returns a timeout based on the test's Deadline, if available.
// Especially important to use in parallel tests, as their individual execution
// can get paused for arbitrary amounts of time.
func WaitTimeout(t *testing.T) time.Duration {
	if d, ok := t.Deadline(); ok {
		// 10% buffer for cleanup and scheduling delay
		return time.Until(d) * 9 / 10
	}
	return DefaultWaitTimeout
}

// TestInterval is just a sensible poll interval that gives fast tests without
// risk of spamming
const TestInterval = 100 * time.Millisecond

// AssertEventually waits for f to return true
func AssertEventually(t *testing.T, f func() bool) {
	assert.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireSignal waits for the channel to close (or receive anything) and
// fatals the test if the default wait timeout is exceeded
func RequireSignal(t *testing.T, ch <-chan struct{}, failMsg string) {
	select {
	case <-ch:
	case <-time.After(WaitTimeout(t)):
		t.Fatal(failMsg)
	}
}
