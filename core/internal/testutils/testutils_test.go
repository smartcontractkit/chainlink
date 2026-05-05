package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitTimeoutWithinDeadlineCap(t *testing.T) {
	got := WaitTimeout(t)

	require.Greater(t, got, time.Duration(0))
	require.LessOrEqual(t, got, DefaultWaitTimeout)
}

func TestWaitTimeoutCustom(t *testing.T) {
	got := WaitTimeoutCustom(t, 10*time.Second)

	require.LessOrEqual(t, got, 10*time.Second)
	require.Greater(t, got, time.Duration(0))
}
