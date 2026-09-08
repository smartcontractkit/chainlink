package standardcapabilities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// stubHTTPTriggerCapability implements coreconfig.HTTPTriggerCapability with
// fixed values; zero values mean "unset in TOML".
type stubHTTPTriggerCapability struct {
	metadataBatchSize            uint16
	sendChannelBufferSize        uint16
	maxAuthorizedKeysPerWorkflow uint16
	requestCacheTTL              uint32
	gatewayConnection            stubHTTPTriggerGatewayConnection
}

type stubHTTPTriggerGatewayConnection struct {
	retryConfig               stubHTTPTriggerRetryConfig
	maxPushMetadataDurationMs uint32
	maxPullMetadataDurationMs uint32
}

type stubHTTPTriggerRetryConfig struct {
	initialIntervalMs int
	maxIntervalTimeMs int
	multiplier        float64
}

func (s *stubHTTPTriggerCapability) MetadataBatchSize() uint16     { return s.metadataBatchSize }
func (s *stubHTTPTriggerCapability) SendChannelBufferSize() uint16 { return s.sendChannelBufferSize }
func (s *stubHTTPTriggerCapability) MaxAuthorizedKeysPerWorkflow() uint16 {
	return s.maxAuthorizedKeysPerWorkflow
}
func (s *stubHTTPTriggerCapability) RequestCacheTTL() uint32 { return s.requestCacheTTL }
func (s *stubHTTPTriggerCapability) GatewayConnection() coreconfig.HTTPTriggerGatewayConnection {
	return &s.gatewayConnection
}

func (g *stubHTTPTriggerGatewayConnection) RetryConfig() coreconfig.HTTPTriggerRetryConfig {
	return &g.retryConfig
}
func (g *stubHTTPTriggerGatewayConnection) MaxPushMetadataDurationMs() uint32 {
	return g.maxPushMetadataDurationMs
}
func (g *stubHTTPTriggerGatewayConnection) MaxPullMetadataDurationMs() uint32 {
	return g.maxPullMetadataDurationMs
}

func (r *stubHTTPTriggerRetryConfig) InitialIntervalMs() int { return r.initialIntervalMs }
func (r *stubHTTPTriggerRetryConfig) MaxIntervalTimeMs() int { return r.maxIntervalTimeMs }
func (r *stubHTTPTriggerRetryConfig) Multiplier() float64    { return r.multiplier }

// stubHTTPActionCapability implements coreconfig.HTTPActionCapability with
// fixed values; zero values mean "unset in TOML".
type stubHTTPActionCapability struct {
	proxyMode         string
	gatewayConnection stubHTTPActionGatewayConnection
	httpClient        stubHTTPActionHTTPClient
}

type stubHTTPActionGatewayConnection struct {
	initialIntervalMs uint32
	maxElapsedTimeMs  uint32
	multiplier        float64
}

type stubHTTPActionHTTPClient struct {
	blockedIPs     []string
	blockedIPsCIDR []string
	allowedPorts   []int
	allowedSchemes []string
	allowedIPs     []string
	allowedIPsCIDR []string
}

func (s *stubHTTPActionCapability) ProxyMode() string { return s.proxyMode }
func (s *stubHTTPActionCapability) GatewayConnection() coreconfig.HTTPActionGatewayConnection {
	return &s.gatewayConnection
}
func (s *stubHTTPActionCapability) HTTPClient() coreconfig.HTTPActionHTTPClient { return &s.httpClient }

func (g *stubHTTPActionGatewayConnection) InitialIntervalMs() uint32 { return g.initialIntervalMs }
func (g *stubHTTPActionGatewayConnection) MaxElapsedTimeMs() uint32  { return g.maxElapsedTimeMs }
func (g *stubHTTPActionGatewayConnection) Multiplier() float64       { return g.multiplier }

func (h *stubHTTPActionHTTPClient) BlockedIPs() []string     { return h.blockedIPs }
func (h *stubHTTPActionHTTPClient) BlockedIPsCIDR() []string { return h.blockedIPsCIDR }
func (h *stubHTTPActionHTTPClient) AllowedPorts() []int      { return h.allowedPorts }
func (h *stubHTTPActionHTTPClient) AllowedSchemes() []string { return h.allowedSchemes }
func (h *stubHTTPActionHTTPClient) AllowedIPs() []string     { return h.allowedIPs }
func (h *stubHTTPActionHTTPClient) AllowedIPsCIDR() []string { return h.allowedIPsCIDR }

func TestInjectHTTPTriggerConfig_TOMLOverridesJobSpec(t *testing.T) {
	t.Parallel()

	cfg := &stubHTTPTriggerCapability{
		metadataBatchSize:            25,
		sendChannelBufferSize:        500,
		maxAuthorizedKeysPerWorkflow: 10,
		requestCacheTTL:              3600,
		gatewayConnection: stubHTTPTriggerGatewayConnection{
			retryConfig:               stubHTTPTriggerRetryConfig{initialIntervalMs: 200, maxIntervalTimeMs: 60000, multiplier: 3.0},
			maxPushMetadataDurationMs: 45000,
			maxPullMetadataDurationMs: 46000,
		},
	}

	// Job spec carries different values; TOML must win.
	jobSpecConfig := `{"metadataBatchSize":999,"sendChannelBufferSize":888,"requestCacheTTL":777,"gatewayConnection":{"retryConfig":{"initialIntervalMs":111}}}`

	got := injectHTTPTriggerConfig(logger.TestLogger(t), cfg, jobSpecConfig)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))

	assert.Equal(t, float64(25), m["metadataBatchSize"])
	assert.Equal(t, float64(500), m["sendChannelBufferSize"])
	assert.Equal(t, float64(10), m["maxAuthorizedKeysPerWorkflow"])
	assert.Equal(t, float64(3600), m["requestCacheTTL"])

	gw := m["gatewayConnection"].(map[string]any)
	assert.Equal(t, float64(45000), gw["maxPushMetadataDurationMs"])
	assert.Equal(t, float64(46000), gw["maxPullMetadataDurationMs"])

	retry := gw["retryConfig"].(map[string]any)
	assert.Equal(t, float64(200), retry["initialIntervalMs"])
	assert.Equal(t, float64(60000), retry["maxIntervalTimeMs"])
	assert.Equal(t, float64(3.0), retry["multiplier"])
}

func TestInjectHTTPTriggerConfig_UnsetTOMLFallsBackToJobSpec(t *testing.T) {
	t.Parallel()

	// All TOML values unset (zero): job-spec values must survive.
	cfg := &stubHTTPTriggerCapability{}
	jobSpecConfig := `{"metadataBatchSize":999,"sendChannelBufferSize":888,"requestCacheTTL":777,"gatewayConnection":{"retryConfig":{"initialIntervalMs":111},"maxPushMetadataDurationMs":555}}`

	got := injectHTTPTriggerConfig(logger.TestLogger(t), cfg, jobSpecConfig)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))

	assert.Equal(t, float64(999), m["metadataBatchSize"])
	assert.Equal(t, float64(888), m["sendChannelBufferSize"])
	assert.Equal(t, float64(777), m["requestCacheTTL"])

	gw := m["gatewayConnection"].(map[string]any)
	assert.Equal(t, float64(555), gw["maxPushMetadataDurationMs"])
	retry := gw["retryConfig"].(map[string]any)
	assert.Equal(t, float64(111), retry["initialIntervalMs"])
}

func TestInjectHTTPTriggerConfig_EmptyConfigProducesEmptyJSON(t *testing.T) {
	t.Parallel()

	// No TOML values, no job-spec config: result is empty JSON so the
	// capability binary applies all of its built-in defaults.
	cfg := &stubHTTPTriggerCapability{}
	got := injectHTTPTriggerConfig(logger.TestLogger(t), cfg, "")
	assert.Equal(t, "{}", got)
}

func TestInjectHTTPTriggerConfig_MalformedJobSpecConfigStillInjects(t *testing.T) {
	t.Parallel()

	cfg := &stubHTTPTriggerCapability{metadataBatchSize: 25}
	got := injectHTTPTriggerConfig(logger.TestLogger(t), cfg, "{not valid json")

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))
	assert.Equal(t, float64(25), m["metadataBatchSize"])
}

func TestInjectHTTPActionConfig_TOMLOverridesJobSpec(t *testing.T) {
	t.Parallel()

	cfg := &stubHTTPActionCapability{
		proxyMode: "direct",
		gatewayConnection: stubHTTPActionGatewayConnection{
			initialIntervalMs: 200,
			maxElapsedTimeMs:  60000,
			multiplier:        3.0,
		},
		httpClient: stubHTTPActionHTTPClient{
			blockedIPs:     []string{"10.0.0.1"},
			blockedIPsCIDR: []string{"10.0.0.0/8"},
			allowedPorts:   []int{8443},
			allowedSchemes: []string{"http", "https"},
			allowedIPs:     []string{"1.2.3.4"},
			allowedIPsCIDR: []string{"1.2.3.0/24"},
		},
	}

	// Job spec carries different values; TOML must win.
	jobSpecConfig := `{"proxyMode":"gateway","gatewayConnection":{"initialIntervalMs":111},"httpClient":{"allowedPorts":[443]}}`

	got := injectHTTPActionConfig(logger.TestLogger(t), cfg, jobSpecConfig)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))

	assert.Equal(t, "direct", m["proxyMode"])

	gw := m["gatewayConnection"].(map[string]any)
	assert.Equal(t, float64(200), gw["initialIntervalMs"])
	assert.Equal(t, float64(60000), gw["maxElapsedTimeMs"])
	assert.Equal(t, float64(3.0), gw["multiplier"])

	hc := m["httpClient"].(map[string]any)
	assert.Equal(t, []any{"10.0.0.1"}, hc["blockedIPs"])
	assert.Equal(t, []any{"10.0.0.0/8"}, hc["blockedIPsCIDR"])
	assert.Equal(t, []any{float64(8443)}, hc["allowedPorts"])
	assert.Equal(t, []any{"http", "https"}, hc["allowedSchemes"])
	assert.Equal(t, []any{"1.2.3.4"}, hc["allowedIPs"])
	assert.Equal(t, []any{"1.2.3.0/24"}, hc["allowedIPsCIDR"])
}

func TestInjectHTTPActionConfig_UnsetTOMLFallsBackToJobSpec(t *testing.T) {
	t.Parallel()

	// All TOML values unset: job-spec values must survive.
	cfg := &stubHTTPActionCapability{}
	jobSpecConfig := `{"proxyMode":"direct","gatewayConnection":{"initialIntervalMs":111},"httpClient":{"allowedPorts":[443,8443]}}`

	got := injectHTTPActionConfig(logger.TestLogger(t), cfg, jobSpecConfig)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))

	assert.Equal(t, "direct", m["proxyMode"])
	gw := m["gatewayConnection"].(map[string]any)
	assert.Equal(t, float64(111), gw["initialIntervalMs"])
	hc := m["httpClient"].(map[string]any)
	assert.Equal(t, []any{float64(443), float64(8443)}, hc["allowedPorts"])
}

func TestInjectHTTPActionConfig_EmptyConfigProducesEmptyJSON(t *testing.T) {
	t.Parallel()

	cfg := &stubHTTPActionCapability{}
	got := injectHTTPActionConfig(logger.TestLogger(t), cfg, "")
	assert.Equal(t, "{}", got)
}

func TestInjectHTTPActionConfig_HTTPClientOnlyFromTOML(t *testing.T) {
	t.Parallel()

	// httpClient restrictions come only from TOML. When TOML sets them and the
	// job spec has none, they must still be injected (direct mode depends on it).
	cfg := &stubHTTPActionCapability{
		proxyMode:  "direct",
		httpClient: stubHTTPActionHTTPClient{allowedPorts: []int{9443}, allowedSchemes: []string{"https"}},
	}

	got := injectHTTPActionConfig(logger.TestLogger(t), cfg, `{"proxyMode":"direct"}`)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &m))

	hc := m["httpClient"].(map[string]any)
	assert.Equal(t, []any{float64(9443)}, hc["allowedPorts"])
	assert.Equal(t, []any{"https"}, hc["allowedSchemes"])
}
