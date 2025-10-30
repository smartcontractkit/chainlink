package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestTelemetryIngressConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ticfg := cfg.TelemetryIngress()
	assert.True(t, ticfg.Logging())
	assert.False(t, ticfg.UniConn())
	assert.Equal(t, uint(1234), ticfg.BufferSize())
	assert.Equal(t, uint(4321), ticfg.MaxBatchSize())
	assert.Equal(t, time.Minute, ticfg.SendInterval())
	assert.Equal(t, 5*time.Second, ticfg.SendTimeout())
	assert.True(t, ticfg.UseBatchSend())

	tec := cfg.TelemetryIngress().Endpoints()

	assert.Len(t, tec, 1)
	assert.Equal(t, "EVM", tec[0].Network())
	assert.Equal(t, "1", tec[0].ChainID())
	assert.Equal(t, "prom.test", tec[0].URL().String())
	assert.Equal(t, "test-pub-key", tec[0].ServerPubKey())
}

func TestTelemetryIngressConfig_ChipIngressEnabled(t *testing.T) {
	t.Run("returns false when ChipIngressEnabled is explicitly false", func(t *testing.T) {
		falseVal := false
		config := &telemetryIngressConfig{
			c: toml.TelemetryIngress{
				ChipIngressEnabled: &falseVal,
			},
		}
		assert.False(t, config.ChipIngressEnabled())
	})

	t.Run("returns true when ChipIngressEnabled is true", func(t *testing.T) {
		trueVal := true
		config := &telemetryIngressConfig{
			c: toml.TelemetryIngress{
				ChipIngressEnabled: &trueVal,
			},
		}
		assert.True(t, config.ChipIngressEnabled())
	})
}
