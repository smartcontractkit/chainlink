package httpaction

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

func TestBuildCapabilityConfig(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"IncomingGlobalBurst":    10,
		"IncomingGlobalRPS":      20,
		"IncomingPerSenderBurst": 30,
		"IncomingPerSenderRPS":   40,
		"OutgoingGlobalBurst":    50,
		"OutgoingGlobalRPS":      60,
		"OutgoingPerSenderBurst": 70,
		"OutgoingPerSenderRPS":   80,
		"ProxyMode":              "direct",
		"GatewayConnection":      map[string]any{"initialIntervalMs": 200},
		"HTTPClient":             map[string]any{"allowedPorts": []int{8443}},
	}
	config, err := BuildCapabilityConfig(input, true)
	require.NoError(t, err)
	require.True(t, config.LocalOnly)
	spec, err := values.FromMapValueProto(config.SpecConfig)
	require.NoError(t, err)
	raw, err := spec.Unwrap()
	require.NoError(t, err)
	decoded := raw.(map[string]any)
	assert.Contains(t, decoded, "incomingRateLimiter")
	assert.Contains(t, decoded, "outgoingRateLimiter")
	assert.Equal(t, "direct", decoded["proxyMode"])
	assert.Contains(t, decoded, "gatewayConnection")
	assert.Contains(t, decoded, "httpClient")
}

func TestPostEnvStartupDoesNotCreateJobs(t *testing.T) {
	t.Parallel()
	require.NoError(t, (&HTTPAction{}).PostEnvStartup(t.Context(), zerolog.Nop(), nil, nil, nil))
}

func TestBuildCapabilityConfigDefaultsToGateway(t *testing.T) {
	t.Parallel()
	config, err := BuildCapabilityConfig(nil, false)
	require.NoError(t, err)
	spec, err := values.FromMapValueProto(config.SpecConfig)
	require.NoError(t, err)
	raw, err := spec.Unwrap()
	require.NoError(t, err)
	assert.Equal(t, "gateway", raw.(map[string]any)["proxyMode"])
	assert.False(t, config.LocalOnly)
}

func TestBuildCapabilityConfigNativeFieldsTakePrecedence(t *testing.T) {
	t.Parallel()
	config, err := BuildCapabilityConfig(map[string]any{
		"ProxyMode":           "gateway",
		"proxyMode":           "direct",
		"IncomingGlobalBurst": 10,
		"incomingRateLimiter": map[string]any{"globalBurst": 20},
		"gatewayConnection":   map[string]any{"initialIntervalMs": 200},
		"MetadataBatchSize":   25,
	}, false)
	require.NoError(t, err)
	spec, err := values.FromMapValueProto(config.SpecConfig)
	require.NoError(t, err)
	raw, err := spec.Unwrap()
	require.NoError(t, err)
	got := raw.(map[string]any)
	assert.Equal(t, "direct", got["proxyMode"])
	assert.EqualValues(t, 20, got["incomingRateLimiter"].(map[string]any)["globalBurst"])
	assert.EqualValues(t, 200, got["gatewayConnection"].(map[string]any)["initialIntervalMs"])
	assert.NotContains(t, got, "metadataBatchSize", "trigger settings do not belong in the action config")
}

func TestBuildCapabilityConfigRejectsUnsupportedValues(t *testing.T) {
	t.Parallel()
	_, err := BuildCapabilityConfig(map[string]any{"HTTPClient": make(chan int)}, false)
	require.ErrorContains(t, err, "build HTTP action spec_config")
}
