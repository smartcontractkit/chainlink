package httptrigger

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
		"IncomingGlobalBurst":          10,
		"IncomingGlobalRPS":            20,
		"IncomingPerSenderBurst":       30,
		"IncomingPerSenderRPS":         40,
		"OutgoingGlobalBurst":          50,
		"OutgoingGlobalRPS":            60,
		"OutgoingPerSenderBurst":       70,
		"OutgoingPerSenderRPS":         80,
		"MetadataBatchSize":            25,
		"SendChannelBufferSize":        500,
		"MaxAuthorizedKeysPerWorkflow": 10,
		"RequestCacheTTL":              3600,
		"GatewayConnection":            map[string]any{"retryConfig": map[string]any{"initialIntervalMs": 200}},
		"MaxPushMetadataDurationMs":    45000,
		"MaxPullMetadataDurationMs":    46000,
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
	assert.EqualValues(t, 25, decoded["metadataBatchSize"])
	assert.EqualValues(t, 500, decoded["sendChannelBufferSize"])
	assert.EqualValues(t, 10, decoded["maxAuthorizedKeysPerWorkflow"])
	assert.EqualValues(t, 3600, decoded["requestCacheTTL"])
	assert.Contains(t, decoded, "gatewayConnection")
	gateway := decoded["gatewayConnection"].(map[string]any)
	assert.EqualValues(t, 45000, gateway["maxPushMetadataDurationMs"])
	assert.EqualValues(t, 46000, gateway["maxPullMetadataDurationMs"])
	assert.NotContains(t, decoded, "maxPushMetadataDurationMs")
	assert.NotContains(t, decoded, "maxPullMetadataDurationMs")
}

func TestPostEnvStartupDoesNotCreateJobs(t *testing.T) {
	t.Parallel()
	require.NoError(t, (&HTTPTrigger{}).PostEnvStartup(t.Context(), zerolog.Nop(), nil, nil, nil))
}

func TestBuildCapabilityConfigNativeFieldsTakePrecedence(t *testing.T) {
	t.Parallel()
	gateway := map[string]any{
		"maxPushMetadataDurationMs": 500,
		"retryConfig":               map[string]any{"initialIntervalMs": 200},
	}
	config, err := BuildCapabilityConfig(map[string]any{
		"MetadataBatchSize":         25,
		"metadataBatchSize":         10,
		"IncomingGlobalBurst":       10,
		"incomingRateLimiter":       map[string]any{"globalBurst": 20},
		"gatewayConnection":         gateway,
		"MaxPushMetadataDurationMs": 45000,
		"MaxPullMetadataDurationMs": 46000,
		"ProxyMode":                 "direct",
		"HTTPClient":                map[string]any{"allowedPorts": []int{8443}},
	}, false)
	require.NoError(t, err)
	spec, err := values.FromMapValueProto(config.SpecConfig)
	require.NoError(t, err)
	raw, err := spec.Unwrap()
	require.NoError(t, err)
	got := raw.(map[string]any)
	assert.EqualValues(t, 10, got["metadataBatchSize"])
	assert.EqualValues(t, 20, got["incomingRateLimiter"].(map[string]any)["globalBurst"])
	gotGateway := got["gatewayConnection"].(map[string]any)
	assert.EqualValues(t, 200, gotGateway["retryConfig"].(map[string]any)["initialIntervalMs"])
	assert.EqualValues(t, 500, gotGateway["maxPushMetadataDurationMs"])
	assert.EqualValues(t, 46000, gotGateway["maxPullMetadataDurationMs"])
	assert.NotContains(t, gateway, "maxPullMetadataDurationMs", "building config must not mutate the input")
	assert.NotContains(t, got, "proxyMode", "action settings do not belong in the trigger config")
	assert.NotContains(t, got, "httpClient")
}

func TestBuildCapabilityConfigRejectsUnsupportedValues(t *testing.T) {
	t.Parallel()
	_, err := BuildCapabilityConfig(map[string]any{"GatewayConnection": make(chan int)}, false)
	require.ErrorContains(t, err, "build HTTP trigger spec_config")
}
