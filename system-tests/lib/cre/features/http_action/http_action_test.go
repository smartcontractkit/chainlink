package httpaction

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
		"IncomingGlobalBurst":    10,
		"IncomingGlobalRPS":      20,
		"IncomingPerSenderBurst": 30,
		"IncomingPerSenderRPS":   40,
		"OutgoingGlobalBurst":    50,
		"OutgoingGlobalRPS":      60,
		"OutgoingPerSenderBurst": 70,
		"OutgoingPerSenderRPS":   80,
		"ProxyMode":              "direct",
		"GatewayConnection":      map[string]any{"initialIntervalMs": 100},
		"HTTPClient":             map[string]any{"allowedPorts": []int{443}},
	}

	tmpl, err := template.New("http-action-config").Parse(configTemplate)
	require.NoError(t, err)

	var rendered bytes.Buffer
	require.NoError(t, tmpl.Execute(&rendered, values))

	var config map[string]any
	require.NoError(t, json.Unmarshal(rendered.Bytes(), &config))
	assert.Contains(t, config, "incomingRateLimiter")
	assert.Contains(t, config, "outgoingRateLimiter")
	assert.NotContains(t, config, "proxyMode")
	assert.NotContains(t, config, "gatewayConnection")
	assert.NotContains(t, config, "httpClient")
}
