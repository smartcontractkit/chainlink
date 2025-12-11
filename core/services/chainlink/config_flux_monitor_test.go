package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestFluxMonitorConfig(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	fm := cfg.FluxMonitor()

	assert.Equal(t, uint32(100), fm.DefaultTransactionQueueDepth())
	assert.True(t, fm.SimulateTransactions())
}
