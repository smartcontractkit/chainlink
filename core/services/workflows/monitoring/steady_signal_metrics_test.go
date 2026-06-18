package monitoring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSteadySignalStateConstants(t *testing.T) {
	require.Equal(t, int64(0), SteadySignalStateUnobserved)
	require.Equal(t, int64(1), SteadySignalStateTransition)
	require.Equal(t, int64(2), SteadySignalStateSteady)
}
