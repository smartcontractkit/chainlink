package plugins

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/smartcontractkit/freeport"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

var ErrExists = errors.New("plugin already registered")

type RegisteredLoop struct {
	Name   string
	EnvCfg loop.EnvConfig
}

// LoopRegistry is responsible for assigning ports to plugins that are to be used for the
// plugin's prometheus HTTP server, and for passing the tracing configuration to the plugin.
type LoopRegistry struct {
	mu       sync.Mutex
	registry map[string]*RegisteredLoop

	lggr                   logger.Logger
	cfgDatabase            config.Database
	cfgTracing             config.Tracing
	cfgTelemetry           config.Telemetry
	telemetryAuthHeaders   map[string]string
	telemetryAuthPubKeyHex string
}

func NewLoopRegistry(lggr logger.Logger, dbConfig config.Database, tracing config.Tracing, telemetry config.Telemetry, telemetryAuthHeaders map[string]string, telemetryAuthPubKeyHex string) *LoopRegistry {
	return &LoopRegistry{
		registry:               map[string]*RegisteredLoop{},
		lggr:                   logger.Named(lggr, "LoopRegistry"),
		cfgDatabase:            dbConfig,
		cfgTracing:             tracing,
		cfgTelemetry:           telemetry,
		telemetryAuthHeaders:   telemetryAuthHeaders,
		telemetryAuthPubKeyHex: telemetryAuthPubKeyHex,
	}
}

func NewTestLoopRegistry(lggr logger.Logger) *LoopRegistry {
	return &LoopRegistry{
		registry: map[string]*RegisteredLoop{},
		lggr:     logger.Named(lggr, "LoopRegistry"),
	}
}

// Register creates a port of the plugin. It is not idempotent. Duplicate calls to Register will return [ErrExists]
// Safe for concurrent use.
func (m *LoopRegistry) Register(id string) (*RegisteredLoop, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.registry[id]; exists {
		return nil, ErrExists
	}
	ports, err := freeport.Take(1)
	if err != nil {
		return nil, fmt.Errorf("failed to get free port: %w", err)
	}
	if len(ports) != 1 {
		return nil, errors.New("failed to get free port: no ports returned")
	}
	envCfg := loop.EnvConfig{PrometheusPort: ports[0]}

	if m.cfgDatabase != nil {
		dbURL := m.cfgDatabase.URL()
		envCfg.DatabaseURL = (*commonconfig.SecretURL)(&dbURL)
		envCfg.DatabaseIdleInTxSessionTimeout = m.cfgDatabase.DefaultIdleInTxSessionTimeout()
		envCfg.DatabaseLockTimeout = m.cfgDatabase.DefaultLockTimeout()
		envCfg.DatabaseQueryTimeout = m.cfgDatabase.DefaultQueryTimeout()
		envCfg.DatabaseLogSQL = m.cfgDatabase.LogSQL()
		envCfg.DatabaseMaxOpenConns = m.cfgDatabase.MaxOpenConns()
		envCfg.DatabaseMaxIdleConns = m.cfgDatabase.MaxIdleConns()
	}

	if m.cfgTracing != nil {
		envCfg.TracingEnabled = m.cfgTracing.Enabled()
		envCfg.TracingCollectorTarget = m.cfgTracing.CollectorTarget()
		envCfg.TracingSamplingRatio = m.cfgTracing.SamplingRatio()
		envCfg.TracingTLSCertPath = m.cfgTracing.TLSCertPath()
		envCfg.TracingAttributes = m.cfgTracing.Attributes()
	}

	if m.cfgTelemetry != nil {
		envCfg.TelemetryEnabled = m.cfgTelemetry.Enabled()
		envCfg.TelemetryEndpoint = m.cfgTelemetry.OtelExporterGRPCEndpoint()
		envCfg.TelemetryInsecureConnection = m.cfgTelemetry.InsecureConnection()
		envCfg.TelemetryCACertFile = m.cfgTelemetry.CACertFile()
		envCfg.TelemetryAttributes = m.cfgTelemetry.ResourceAttributes()
		envCfg.TelemetryTraceSampleRatio = m.cfgTelemetry.TraceSampleRatio()
		envCfg.TelemetryEmitterBatchProcessor = m.cfgTelemetry.EmitterBatchProcessor()
		envCfg.TelemetryEmitterExportTimeout = m.cfgTelemetry.EmitterExportTimeout()
		envCfg.TelemetryAuthPubKeyHex = m.telemetryAuthPubKeyHex
		envCfg.ChipIngressEndpoint = m.cfgTelemetry.ChipIngressEndpoint()
	}
	m.lggr.Debugf("Registered loopp %q with port %d", id, envCfg.PrometheusPort)

	// Add auth header after logging config
	if m.cfgTelemetry != nil {
		envCfg.TelemetryAuthHeaders = m.telemetryAuthHeaders
	}

	m.registry[id] = &RegisteredLoop{Name: id, EnvCfg: envCfg}
	return m.registry[id], nil
}

// Unregister remove a loop from the registry
// Safe for concurrent use.
func (m *LoopRegistry) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	loop, exists := m.registry[id]
	if !exists {
		m.lggr.Debugf("Trying to unregistered a loop that is not registered %q", id)
		return
	}

	freeport.Return([]int{loop.EnvCfg.PrometheusPort})
	delete(m.registry, id)
	m.lggr.Debugf("Unregistered loopp %q", id)
}

// Return slice sorted by plugin name. Safe for concurrent use.
func (m *LoopRegistry) List() []*RegisteredLoop {
	var registeredLoops []*RegisteredLoop
	m.mu.Lock()
	for _, known := range m.registry {
		registeredLoops = append(registeredLoops, known)
	}
	m.mu.Unlock()

	sort.Slice(registeredLoops, func(i, j int) bool {
		return registeredLoops[i].Name < registeredLoops[j].Name
	})
	return registeredLoops
}

// Get plugin by id. Safe for concurrent use.
func (m *LoopRegistry) Get(id string) (*RegisteredLoop, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.registry[id]
	return p, exists
}
