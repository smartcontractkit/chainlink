package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPluginPortManager(t *testing.T) {
	// register one
	m := NewLoopRegistry(logger.TestLogger(t), nil, nil)
	pFoo, err := m.Register("foo")
	require.NoError(t, err)
	require.Equal(t, "foo", pFoo.Name)
	require.Greater(t, pFoo.EnvCfg.PrometheusPort, 0)
	// test duplicate
	pNil, err := m.Register("foo")
	require.ErrorIs(t, err, ErrExists)
	require.Nil(t, pNil)
	// ensure increasing port assignment
	pBar, err := m.Register("bar")
	require.NoError(t, err)
	require.Equal(t, "bar", pBar.Name)
	require.Equal(t, pFoo.EnvCfg.PrometheusPort+1, pBar.EnvCfg.PrometheusPort)
}

type mockCfgTracing struct{}

func (m *mockCfgTracing) Attributes() map[string]string {
	return map[string]string{"attribute": "value"}
}
func (m *mockCfgTracing) Enabled() bool           { return true }
func (m *mockCfgTracing) NodeID() string          { return "" }
func (m *mockCfgTracing) CollectorTarget() string { return "http://localhost:9000" }
func (m *mockCfgTracing) SamplingRatio() float64  { return 0.1 }
func (m *mockCfgTracing) TLSCertPath() string     { return "/path/to/cert.pem" }
func (m *mockCfgTracing) Mode() string            { return "tls" }

type mockCfgTelemetry struct{}

func (m mockCfgTelemetry) Enabled() bool { return true }

func (m mockCfgTelemetry) InsecureConnection() bool { return true }

func (m mockCfgTelemetry) CACertFile() string { return "path/to/cert.pem" }

func (m mockCfgTelemetry) OtelExporterGRPCEndpoint() string { return "http://localhost:9001" }

func (m mockCfgTelemetry) ResourceAttributes() map[string]string {
	return map[string]string{"foo": "bar"}
}

func (m mockCfgTelemetry) TraceSampleRatio() float64 { return 0.42 }

func TestLoopRegistry_Register(t *testing.T) {
	mockCfgTracing := &mockCfgTracing{}
	mockCfgTelemetry := &mockCfgTelemetry{}
	registry := make(map[string]*RegisteredLoop)

	// Create a LoopRegistry instance with mockCfgTracing
	loopRegistry := &LoopRegistry{
		lggr:         logger.TestLogger(t),
		registry:     registry,
		cfgTracing:   mockCfgTracing,
		cfgTelemetry: mockCfgTelemetry,
	}

	// Test case 1: Register new loop
	registeredLoop, err := loopRegistry.Register("testID")
	require.Nil(t, err)
	require.Equal(t, "testID", registeredLoop.Name)

	envCfg := registeredLoop.EnvCfg
	require.True(t, envCfg.TracingEnabled)
	require.Equal(t, "http://localhost:9000", envCfg.TracingCollectorTarget)
	require.Equal(t, map[string]string{"attribute": "value"}, envCfg.TracingAttributes)
	require.Equal(t, 0.1, envCfg.TracingSamplingRatio)
	require.Equal(t, "/path/to/cert.pem", envCfg.TracingTLSCertPath)

	require.True(t, envCfg.TelemetryEnabled)
	require.True(t, envCfg.TelemetryInsecureConnection)
	require.Equal(t, "path/to/cert.pem", envCfg.TelemetryCACertFile)
	require.Equal(t, "http://localhost:9001", envCfg.TelemetryEndpoint)
	require.Equal(t, loop.OtelAttributes{"foo": "bar"}, envCfg.TelemetryAttributes)
	require.Equal(t, 0.42, envCfg.TelemetryTraceSampleRatio)
}
