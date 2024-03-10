package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeperConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	keeper := cfg.Keeper()

	assert.Equal(t, uint32(17), keeper.DefaultTransactionQueueDepth())
	assert.Equal(t, uint16(12), keeper.GasPriceBufferPercent())
	assert.Equal(t, uint16(43), keeper.GasTipCapBufferPercent())
	assert.Equal(t, uint16(89), keeper.BaseFeeBufferPercent())
	assert.Equal(t, int64(91), keeper.TurnLookBack())
	assert.Equal(t, int64(31), keeper.MaxGracePeriod())

	registry := keeper.Registry()
	assert.Equal(t, uint32(90), registry.CheckGasOverhead())
	assert.Equal(t, uint32(4294967295), registry.PerformGasOverhead())
	assert.Equal(t, uint32(5000), registry.MaxPerformDataSize())
	assert.Equal(t, 1*time.Hour, registry.SyncInterval())
	assert.Equal(t, uint32(31), registry.SyncUpkeepQueueSize())
}
