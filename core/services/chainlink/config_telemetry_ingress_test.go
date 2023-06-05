package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryIngressConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ticfg := cfg.TelemetryIngress()
	assert.True(t, ticfg.Logging())
	assert.True(t, ticfg.UniConn())
	assert.Equal(t, "test-pub-key", ticfg.ServerPubKey())
	assert.Equal(t, "https://prom.test", ticfg.URL().String())
	assert.Equal(t, uint(1234), ticfg.BufferSize())
	assert.Equal(t, uint(4321), ticfg.MaxBatchSize())
	assert.Equal(t, time.Minute, ticfg.SendInterval())
	assert.Equal(t, 5*time.Second, ticfg.SendTimeout())
	assert.True(t, ticfg.UseBatchSend())
}
