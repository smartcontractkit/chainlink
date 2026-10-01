package job

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// SetCapRegMetricsForTest routes capabilities_registry validation metrics to m for the duration
// of t. Tests using it must not run in parallel.
func SetCapRegMetricsForTest(t *testing.T, m *globalconfig.Metrics) {
	prev := capRegMetrics
	capRegMetrics = func() *globalconfig.Metrics { return m }
	t.Cleanup(func() { capRegMetrics = prev })
}
