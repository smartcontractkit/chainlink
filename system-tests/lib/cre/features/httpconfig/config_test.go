package httpconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

func TestBuildNativeFieldsTakePrecedence(t *testing.T) {
	t.Parallel()
	config, err := Build(map[string]any{
		"ProxyMode":           "gateway",
		"proxyMode":           "direct",
		"IncomingGlobalBurst": 10,
		"incomingRateLimiter": map[string]any{"globalBurst": 20},
		"gatewayConnection":   map[string]any{"retryConfig": map[string]any{"initialIntervalMs": 200}},
	}, false)
	require.NoError(t, err)
	spec, err := values.FromMapValueProto(config.SpecConfig)
	require.NoError(t, err)
	raw, err := spec.Unwrap()
	require.NoError(t, err)
	got := raw.(map[string]any)
	assert.Equal(t, "direct", got["proxyMode"])
	assert.EqualValues(t, 20, got["incomingRateLimiter"].(map[string]any)["globalBurst"])
	assert.EqualValues(t, 200, got["gatewayConnection"].(map[string]any)["retryConfig"].(map[string]any)["initialIntervalMs"])
}

func TestBuildRejectsUnsupportedValues(t *testing.T) {
	t.Parallel()
	_, err := Build(map[string]any{"HTTPClient": make(chan int)}, false)
	require.ErrorContains(t, err, "build HTTP spec_config")
}
