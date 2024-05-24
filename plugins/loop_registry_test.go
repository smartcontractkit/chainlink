package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPluginPortManager(t *testing.T) {
	// register one
	m := NewLoopRegistry(logger.TestLogger(t), nil)
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

// Mock tracing config
type MockCfgTracing struct{}

func (m *MockCfgTracing) Attributes() map[string]string {
	return map[string]string{"attribute": "value"}
}
func (m *MockCfgTracing) Enabled() bool           { return true }
func (m *MockCfgTracing) NodeID() string          { return "" }
func (m *MockCfgTracing) CollectorTarget() string { return "http://localhost:9000" }
func (m *MockCfgTracing) SamplingRatio() float64  { return 0.1 }
func (m *MockCfgTracing) TLSCertPath() string     { return "/path/to/cert.pem" }
func (m *MockCfgTracing) Mode() string            { return "tls" }

func TestLoopRegistry_Register(t *testing.T) {
	mockCfgTracing := &MockCfgTracing{}
	registry := make(map[string]*RegisteredLoop)

	// Create a LoopRegistry instance with mockCfgTracing
	loopRegistry := &LoopRegistry{
		lggr:       logger.TestLogger(t),
		registry:   registry,
		cfgTracing: mockCfgTracing,
	}

	// Test case 1: Register new loop
	registeredLoop, err := loopRegistry.Register("testID")
	require.Nil(t, err)
	require.Equal(t, "testID", registeredLoop.Name)
	require.True(t, registeredLoop.EnvCfg.TracingEnabled)
	require.Equal(t, "http://localhost:9000", registeredLoop.EnvCfg.TracingCollectorTarget)
	require.Equal(t, map[string]string{"attribute": "value"}, registeredLoop.EnvCfg.TracingAttributes)
	require.Equal(t, 0.1, registeredLoop.EnvCfg.TracingSamplingRatio)
	require.Equal(t, "/path/to/cert.pem", registeredLoop.EnvCfg.TracingTLSCertPath)
}
