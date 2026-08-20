package crelimits_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crelimits"
)

// errGate is a GateLimiter whose AllowErr always fails with a given error, standing in
// for a settings read failure (as opposed to limits.ErrorNotAllowed, which means the
// gate evaluated fine and is closed).
type errGate struct{ err error }

func (g errGate) Limit(context.Context) (bool, error) { return false, g.err }
func (g errGate) AllowErr(context.Context) error      { return g.err }
func (g errGate) Close() error                        { return nil }

type recordingLogger struct{ calls int }

func (l *recordingLogger) Errorw(string, ...any) { l.calls++ }

func TestGateOpen(t *testing.T) {
	t.Parallel()

	readErr := errors.New("settings unavailable")

	testCases := []struct {
		name     string
		gate     limits.GateLimiter
		wantOpen bool
		wantErr  error
	}{
		{name: "open gate", gate: limits.NewGateLimiter(true), wantOpen: true},
		{name: "closed gate is not an error", gate: limits.NewGateLimiter(false), wantOpen: false},
		{name: "read failure surfaces", gate: errGate{err: readErr}, wantOpen: false, wantErr: readErr},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			open, err := crelimits.GateOpen(t.Context(), tc.gate)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantOpen, open)
		})
	}
}

// TestGateOpen_ClosedIsDistinctFromUnevaluatable is the property the fail-closed call
// sites depend on: a closed gate must not look like a read failure, and vice versa.
func TestGateOpen_ClosedIsDistinctFromUnevaluatable(t *testing.T) {
	t.Parallel()

	closedOpen, closedErr := crelimits.GateOpen(t.Context(), limits.NewGateLimiter(false))
	require.NoError(t, closedErr, "a closed gate is a normal outcome, not an error")
	assert.False(t, closedOpen)

	failedOpen, failedErr := crelimits.GateOpen(t.Context(), errGate{err: errors.New("boom")})
	require.Error(t, failedErr, "an unevaluatable gate must be distinguishable so callers can fail closed")
	assert.False(t, failedOpen)
}

func TestGateAllows(t *testing.T) {
	t.Parallel()

	t.Run("open gate allows without logging", func(t *testing.T) {
		t.Parallel()
		lggr := &recordingLogger{}
		assert.True(t, crelimits.GateAllows(t.Context(), lggr, limits.NewGateLimiter(true), "TestGate"))
		assert.Zero(t, lggr.calls)
	})

	t.Run("closed gate denies without logging", func(t *testing.T) {
		t.Parallel()
		lggr := &recordingLogger{}
		assert.False(t, crelimits.GateAllows(t.Context(), lggr, limits.NewGateLimiter(false), "TestGate"))
		assert.Zero(t, lggr.calls, "a closed gate is expected, not an unexpected error")
	})

	t.Run("read failure denies and logs", func(t *testing.T) {
		t.Parallel()
		lggr := &recordingLogger{}
		assert.False(t, crelimits.GateAllows(t.Context(), lggr, errGate{err: errors.New("boom")}, "TestGate"))
		assert.Equal(t, 1, lggr.calls)
	})
}
