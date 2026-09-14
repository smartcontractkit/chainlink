package httptrigger

import (
	"bytes"
	"encoding/json"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigTemplateOmitsNodeConfig(t *testing.T) {
	values := map[string]any{
		"IncomingGlobalBurst":          10,
		"IncomingGlobalRPS":            20,
		"IncomingPerSenderBurst":       30,
		"IncomingPerSenderRPS":         40,
		"OutgoingGlobalBurst":          50,
		"OutgoingGlobalRPS":            60,
		"OutgoingPerSenderBurst":       70,
		"OutgoingPerSenderRPS":         80,
		"MetadataBatchSize":            50,
		"SendChannelBufferSize":        1000,
		"MaxAuthorizedKeysPerWorkflow": 100,
		"RequestCacheTTL":              86400,
		"GatewayConnection":            map[string]any{"maxPushMetadataDurationMs": 30000},
		"MaxPushMetadataDurationMs":    30000,
		"MaxPullMetadataDurationMs":    30000,
	}

	tmpl, err := template.New("http-trigger-config").Parse(configTemplate)
	require.NoError(t, err)

	var rendered bytes.Buffer
	require.NoError(t, tmpl.Execute(&rendered, values))

	var config map[string]any
	require.NoError(t, json.Unmarshal(rendered.Bytes(), &config))
	assert.Contains(t, config, "incomingRateLimiter")
	assert.Contains(t, config, "outgoingRateLimiter")
	assert.NotContains(t, config, "metadataBatchSize")
	assert.NotContains(t, config, "sendChannelBufferSize")
	assert.NotContains(t, config, "maxAuthorizedKeysPerWorkflow")
	assert.NotContains(t, config, "requestCacheTTL")
	assert.NotContains(t, config, "gatewayConnection")
	assert.NotContains(t, config, "maxPushMetadataDurationMs")
	assert.NotContains(t, config, "maxPullMetadataDurationMs")
}
