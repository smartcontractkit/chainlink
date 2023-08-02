package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFluxMonitorConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	fm := cfg.FluxMonitor()

	assert.Equal(t, uint32(100), fm.DefaultTransactionQueueDepth())
	assert.Equal(t, true, fm.SimulateTransactions())
}
