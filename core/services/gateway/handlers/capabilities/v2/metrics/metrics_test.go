package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMetrics(t *testing.T) {
	t.Parallel()

	metrics, err := NewMetrics()
	require.NoError(t, err)
	require.NotNil(t, metrics)
	require.NotNil(t, metrics.Action)
	require.NotNil(t, metrics.Trigger)
}
