package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitTimeoutUsesDeadlineBudgetNotDefaultCap(t *testing.T) {
	if d, ok := t.Deadline(); ok {
		expectBefore := time.Until(d) * 9 / 10
		got := WaitTimeout(t)
		expectAfter := time.Until(d) * 9 / 10
		require.Greater(t, got, time.Duration(0))
		// got uses time.Until(deadline) inside WaitTimeout between these snapshots
		require.GreaterOrEqual(t, got, expectAfter)
		require.LessOrEqual(t, got, expectBefore)
	} else {
		require.Equal(t, DefaultWaitTimeout, WaitTimeout(t))
	}
}

func TestWaitTimeoutCustom(t *testing.T) {
	requested := 10 * time.Second

	if d, ok := t.Deadline(); ok {
		expectBefore := time.Until(d) * 9 / 10
		got := WaitTimeoutCustom(t, requested)
		expectAfter := time.Until(d) * 9 / 10
		require.Greater(t, got, time.Duration(0))
		require.LessOrEqual(t, got, requested)
		require.GreaterOrEqual(t, got, min(expectAfter, requested))
		require.LessOrEqual(t, got, min(expectBefore, requested))
	} else {
		require.Equal(t, requested, WaitTimeoutCustom(t, requested))
	}
}
