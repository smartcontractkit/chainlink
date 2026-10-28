package chainlink

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

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

func TestTelemetryIngressConfig_BeholderEnabled(t *testing.T) {
	trueVal, falseVal := true, false
	assert.True(t, (&telemetryIngressConfig{c: toml.TelemetryIngress{BeholderEnabled: &trueVal}}).BeholderEnabled())
	assert.False(t, (&telemetryIngressConfig{c: toml.TelemetryIngress{BeholderEnabled: &falseVal}}).BeholderEnabled())
}

func TestCoreValidateConfig_TelemetryIngressBeholderEnabled(t *testing.T) {
	const enabled = `
[TelemetryIngress]
BeholderEnabled = true
`
	const beholderChip = `
[Telemetry]
Enabled = true
ChipIngressEndpoint = 'chip.example.com:443'
`
	for _, tt := range []struct {
		name    string
		toml    string
		wantErr string // empty: the BeholderEnabled rule must not fire
	}{
		{name: "off by default", toml: ``},
		{name: "requires Beholder telemetry", toml: enabled,
			wantErr: "TelemetryIngress.BeholderEnabled: invalid value (true): requires Telemetry.Enabled and Telemetry.ChipIngressEndpoint"},
		{name: "requires a chip ingress endpoint", toml: enabled + "\n[Telemetry]\nEnabled = true\n",
			wantErr: "TelemetryIngress.BeholderEnabled: invalid value (true): requires Telemetry.Enabled and Telemetry.ChipIngressEndpoint"},
		{name: "cannot combine with ChipIngressEnabled", toml: "[TelemetryIngress]\nBeholderEnabled = true\nChipIngressEnabled = true\n" + beholderChip,
			wantErr: "TelemetryIngress.BeholderEnabled: invalid value (true): cannot be combined with TelemetryIngress.ChipIngressEnabled"},
		{name: "valid with Beholder chip ingress", toml: enabled + beholderChip},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var c Config
			require.NoError(t, commonconfig.DecodeTOML(strings.NewReader(tt.toml), &c))
			c.setDefaults()
			err := c.Core.ValidateConfig()
			if tt.wantErr == "" {
				if err != nil {
					assert.NotContains(t, err.Error(), "TelemetryIngress.BeholderEnabled")
				}
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
