package plugins

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

func TestPluginPortManager(t *testing.T) {
	// register one
	m := NewTestLoopRegistry(logger.Test(t))
	pFoo, err := m.Register("foo")
	require.NoError(t, err)
	require.Equal(t, "foo", pFoo.Name)
	require.Positive(t, pFoo.EnvCfg.PrometheusPort)
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

func (m mockCfgTelemetry) EmitterBatchProcessor() bool { return true }

func (m mockCfgTelemetry) EmitterExportTimeout() time.Duration { return 1 * time.Second }

func (m mockCfgTelemetry) ChipIngressEndpoint() string { return "example.com/chip-ingress" }

func (m mockCfgTelemetry) ChipIngressInsecureConnection() bool { return false }

func (m mockCfgTelemetry) HeartbeatInterval() time.Duration {
	return 5 * time.Second
}

func (m mockCfgTelemetry) LogStreamingEnabled() bool { return false }

type mockCfgDatabase struct{}

func (m mockCfgDatabase) Backup() config.Backup { panic("unimplemented") }

func (m mockCfgDatabase) Listener() config.Listener { return mockCfgListener{} }

func (m mockCfgDatabase) Lock() config.Lock { panic("unimplemented") }

func (m mockCfgDatabase) DefaultIdleInTxSessionTimeout() time.Duration { return time.Hour }

func (m mockCfgDatabase) DefaultLockTimeout() time.Duration { return time.Minute }

func (m mockCfgDatabase) DefaultQueryTimeout() time.Duration { return time.Second }

func (m mockCfgDatabase) DriverName() string { panic("unimplemented") }

func (m mockCfgDatabase) LogSQL() bool { return true }

func (m mockCfgDatabase) MaxIdleConns() int { return 99 }

func (m mockCfgDatabase) MaxOpenConns() int { return 42 }

func (m mockCfgDatabase) MigrateDatabase() bool { panic("unimplemented") }

func (m mockCfgDatabase) URL() url.URL {
	return url.URL{Scheme: "fake", Host: "database.url"}
}

type mockCfgListener struct{}

func (l mockCfgListener) MaxReconnectDuration() time.Duration { panic("unimplemented") }

func (l mockCfgListener) MinReconnectInterval() time.Duration { panic("unimplemented") }

func (l mockCfgListener) FallbackPollInterval() time.Duration { return time.Microsecond }

type mockCfgMercury struct{}

func (m mockCfgMercury) Credentials(credName string) *types.MercuryCredentials { panic("implement me") }

func (m mockCfgMercury) Cache() config.MercuryCache {
	return mockCfgCache{}
}

func (m mockCfgMercury) TLS() config.MercuryTLS { panic("implement me") }

func (m mockCfgMercury) Transmitter() config.MercuryTransmitter {
	return mockCfgTransmitter{}
}

func (m mockCfgMercury) VerboseLogging() bool {
	return true
}

type mockCfgCache struct{}

func (m mockCfgCache) LatestReportTTL() time.Duration {
	return time.Microsecond
}

func (m mockCfgCache) MaxStaleAge() time.Duration {
	return time.Nanosecond
}

func (m mockCfgCache) LatestReportDeadline() time.Duration {
	return time.Millisecond
}

type mockCfgTransmitter struct{}

func (t mockCfgTransmitter) Protocol() config.MercuryTransmitterProtocol { return "foo" }

func (t mockCfgTransmitter) TransmitQueueMaxSize() uint32 { return 42 }

func (t mockCfgTransmitter) TransmitTimeout() time.Duration { return time.Second }

func (t mockCfgTransmitter) TransmitConcurrency() uint32 { return 13 }

func (t mockCfgTransmitter) ReaperFrequency() time.Duration { return time.Hour }

func (t mockCfgTransmitter) ReaperMaxAge() time.Duration { return time.Minute }

func TestLoopRegistry_Register(t *testing.T) {
	mockCfgDatabase := &mockCfgDatabase{}
	mockCfgMercury := &mockCfgMercury{}
	mockCfgTracing := &mockCfgTracing{}
	mockCfgTelemetry := &mockCfgTelemetry{}
	registry := make(map[string]*RegisteredLoop)

	// Create a LoopRegistry instance with mockCfgTracing
	loopRegistry := &LoopRegistry{
		registry:         registry,
		lggr:             logger.Test(t),
		appID:            "app-id",
		featureLogPoller: true,
		cfgDatabase:      mockCfgDatabase,
		cfgMercury:       mockCfgMercury,
		cfgTracing:       mockCfgTracing,
		cfgTelemetry:     mockCfgTelemetry,
	}

	// Test case 1: Register new loop
	registeredLoop, err := loopRegistry.Register("testID")
	require.NoError(t, err)
	require.Equal(t, "testID", registeredLoop.Name)

	envCfg := registeredLoop.EnvCfg

	require.Equal(t, "app-id", envCfg.AppID)
	require.Equal(t, (*commonconfig.SecretURL)(&url.URL{Scheme: "fake", Host: "database.url"}), envCfg.DatabaseURL)
	require.Equal(t, time.Hour, envCfg.DatabaseIdleInTxSessionTimeout)
	require.Equal(t, time.Minute, envCfg.DatabaseLockTimeout)
	require.Equal(t, time.Second, envCfg.DatabaseQueryTimeout)
	require.Equal(t, time.Microsecond, envCfg.DatabaseListenerFallbackPollInterval)
	require.True(t, envCfg.DatabaseLogSQL)
	require.Equal(t, 42, envCfg.DatabaseMaxOpenConns)
	require.Equal(t, 99, envCfg.DatabaseMaxIdleConns)

	require.True(t, envCfg.FeatureLogPoller)

	require.Equal(t, time.Millisecond, envCfg.MercuryCacheLatestReportDeadline)
	require.Equal(t, time.Microsecond, envCfg.MercuryCacheLatestReportTTL)
	require.Equal(t, time.Nanosecond, envCfg.MercuryCacheMaxStaleAge)
	require.Equal(t, "foo", envCfg.MercuryTransmitterProtocol)
	require.Equal(t, uint32(42), envCfg.MercuryTransmitterTransmitQueueMaxSize)
	require.Equal(t, time.Second, envCfg.MercuryTransmitterTransmitTimeout)
	require.Equal(t, uint32(13), envCfg.MercuryTransmitterTransmitConcurrency)
	require.Equal(t, time.Hour, envCfg.MercuryTransmitterReaperFrequency)
	require.Equal(t, time.Minute, envCfg.MercuryTransmitterReaperMaxAge)
	require.True(t, envCfg.MercuryVerboseLogging)

	require.True(t, envCfg.TracingEnabled)
	require.Equal(t, "http://localhost:9000", envCfg.TracingCollectorTarget)
	require.Equal(t, map[string]string{"attribute": "value"}, envCfg.TracingAttributes)
	require.Equal(t, 0.1, envCfg.TracingSamplingRatio)
	require.Equal(t, "/path/to/cert.pem", envCfg.TracingTLSCertPath)

	require.True(t, envCfg.TelemetryEnabled)
	require.True(t, envCfg.TelemetryInsecureConnection)
	require.Equal(t, "path/to/cert.pem", envCfg.TelemetryCACertFile)
	require.Equal(t, "http://localhost:9001", envCfg.TelemetryEndpoint)
	require.Equal(t, beholder.OtelAttributes{"foo": "bar"}, envCfg.TelemetryAttributes)
	require.Equal(t, 0.42, envCfg.TelemetryTraceSampleRatio)
	require.True(t, envCfg.TelemetryEmitterBatchProcessor)
	require.Equal(t, 1*time.Second, envCfg.TelemetryEmitterExportTimeout)
	require.False(t, envCfg.TelemetryLogStreamingEnabled)

	require.Equal(t, "example.com/chip-ingress", envCfg.ChipIngressEndpoint)
}
