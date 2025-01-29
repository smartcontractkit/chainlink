package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	return tests.Context(tb)
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

// AssertEventually calls assert.Eventually with default wait and tick durations.
func AssertEventually(t *testing.T, f func() bool) bool {
	return assert.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireEventually calls assert.Eventually with default wait and tick durations.
func RequireEventually(t *testing.T, f func() bool) {
	require.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}
