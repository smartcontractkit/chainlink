package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitTimeoutCapsLongPackageDeadline(t *testing.T) {
	got := WaitTimeout(t)

	require.LessOrEqual(t, got, DefaultWaitTimeout)
	require.Greater(t, got, time.Duration(0))
}
