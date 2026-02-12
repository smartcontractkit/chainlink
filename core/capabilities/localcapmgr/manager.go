package localcapmgr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// LocalCapabilityManager handles the full lifecycle of capabilities this node hosts.
// It starts, stops, and reconfigures local capabilities based on the on-chain registry state.
// This replaces logic currently in standardcapabilities/delegate.go for registry-based launching.
type LocalCapabilityManager interface {
	services.Service

	// Reconcile compares running capabilities against desired state from registry.
	// Starts new capabilities, stops removed ones, restarts those with changed config.
	// Called by Launcher.OnNewRegistry() for each registry update.
	Reconcile(ctx context.Context, myCapabilityDONs []registrysyncer.DON, localRegistry *registrysyncer.LocalRegistry) error
}

// ServiceBuilder creates job services for a capability. This wraps
// standardcapabilities.Delegate.NewServices to decouple localcapmgr from the delegate.
type ServiceBuilder func(ctx context.Context, capID string, command string, configJSON string) ([]job.ServiceCtx, error)

// runningCapability tracks a started capability.
type runningCapability struct {
	capID      string
	donID      uint32
	services   []job.ServiceCtx
	configHash string
}

// capabilityInfo describes a capability that should be running.
type capabilityInfo struct {
	capID      string
	donID      uint32
	config     registrysyncer.CapabilityConfiguration
	configHash string
}

// runningKey returns a deterministic key for tracking running capabilities.
func runningKey(capID string, donID uint32) string {
	return fmt.Sprintf("%s:%d", capID, donID)
}

// Params contains all dependencies for creating a LocalCapabilityManager.
type Params struct {
	Logger        corelogger.Logger
	LocalConfig   config.LocalCapabilities
	BuildServices ServiceBuilder
}

type localCapabilityManager struct {
	services.StateMachine
	lggr corelogger.Logger

	localCfg      config.LocalCapabilities
	buildServices ServiceBuilder

	runningCapabilities map[string]*runningCapability
	mu                  sync.RWMutex

	metrics *metrics
}

// New creates a LocalCapabilityManager.
func New(params Params) (LocalCapabilityManager, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create local capability manager metrics: %w", err)
	}
	return &localCapabilityManager{
		lggr:                params.Logger.Named("LocalCapabilityManager"),
		localCfg:            params.LocalConfig,
		buildServices:       params.BuildServices,
		runningCapabilities: make(map[string]*runningCapability),
		metrics:             metrics,
	}, nil
}

func (m *localCapabilityManager) Start(ctx context.Context) error {
	return m.StartOnce("LocalCapabilityManager", func() error {
		m.lggr.Info("LocalCapabilityManager started")
		return nil
	})
}

func (m *localCapabilityManager) Close() error {
	return m.StopOnce("LocalCapabilityManager", func() error {
		m.mu.Lock()
		defer m.mu.Unlock()

		var firstErr error
		for key, rc := range m.runningCapabilities {
			m.lggr.Infow("Stopping capability on shutdown", "capID", rc.capID, "donID", rc.donID)
			if err := m.stopRunningCapability(rc); err != nil {
				m.lggr.Errorw("Failed to stop capability on shutdown", "key", key, "error", err)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
		m.runningCapabilities = make(map[string]*runningCapability)
		m.lggr.Info("LocalCapabilityManager stopped")
		return firstErr
	})
}

func (m *localCapabilityManager) Ready() error {
	return m.StateMachine.Ready()
}

func (m *localCapabilityManager) HealthReport() map[string]error {
	return map[string]error{m.Name(): m.Ready()}
}

func (m *localCapabilityManager) Name() string {
	return m.lggr.Name()
}

// Reconcile compares running capabilities against the desired state from the registry.
// It starts new capabilities, stops removed ones, and restarts those with changed config.
func (m *localCapabilityManager) Reconcile(
	ctx context.Context,
	myCapabilityDONs []registrysyncer.DON,
	localRegistry *registrysyncer.LocalRegistry,
) error {
	desired := m.buildDesiredState(myCapabilityDONs)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop capabilities that should no longer be running.
	for key, rc := range m.runningCapabilities {
		if _, ok := desired[key]; !ok {
			m.lggr.Infow("Stopping removed capability", "capID", rc.capID, "donID", rc.donID)
			if err := m.stopRunningCapability(rc); err != nil {
				m.lggr.Errorw("Failed to stop removed capability", "capID", rc.capID, "donID", rc.donID, "error", err)
			}
			m.metrics.recordStop(ctx, rc.capID)
			delete(m.runningCapabilities, key)
		}
	}

	// Start new capabilities or restart those with changed config.
	for key, info := range desired {
		existing, ok := m.runningCapabilities[key]
		if ok && existing.configHash == info.configHash {
			continue // already running with same config
		}

		if ok {
			// Config changed - stop old instance first.
			m.lggr.Infow("Restarting capability due to config change",
				"capID", info.capID, "donID", info.donID,
				"oldHash", existing.configHash, "newHash", info.configHash)
			if err := m.stopRunningCapability(existing); err != nil {
				m.lggr.Errorw("Failed to stop capability for config update", "capID", info.capID, "error", err)
			}
			m.metrics.recordConfigUpdate(ctx, info.capID)
			delete(m.runningCapabilities, key)
		}

		// Start new capability.
		rc, err := m.startCapability(ctx, info)
		if err != nil {
			m.lggr.Errorw("Failed to start capability", "capID", info.capID, "donID", info.donID, "error", err)
			continue
		}
		m.runningCapabilities[key] = rc
	}

	m.metrics.recordRunning(ctx, int64(len(m.runningCapabilities)))
	return nil
}

// buildDesiredState extracts capabilities that should be running from DON configs.
// Only includes capabilities that are in the RegistryBasedLaunchAllowlist.
func (m *localCapabilityManager) buildDesiredState(myCapabilityDONs []registrysyncer.DON) map[string]*capabilityInfo {
	desired := make(map[string]*capabilityInfo)
	for _, don := range myCapabilityDONs {
		for capID, capCfg := range don.CapabilityConfigurations {
			if m.localCfg == nil || !m.localCfg.IsAllowlisted(capID) {
				continue
			}

			key := runningKey(capID, don.ID)
			desired[key] = &capabilityInfo{
				capID:      capID,
				donID:      don.ID,
				config:     capCfg,
				configHash: configHash(capCfg.Config),
			}
		}
	}
	return desired
}

// startCapability creates and starts services for a single capability via the delegate.
func (m *localCapabilityManager) startCapability(ctx context.Context, info *capabilityInfo) (*runningCapability, error) {
	start := time.Now()

	command := m.resolveCapabilityBinary(info.capID)
	configJSON := string(info.config.Config) // TODO: wrong, extract from localCfg, later merge with onchain one

	svcs, err := m.buildServices(ctx, info.capID, command, configJSON)
	if err != nil {
		return nil, fmt.Errorf("build services for %s: %w", info.capID, err)
	}

	for i, svc := range svcs {
		if err := svc.Start(ctx); err != nil {
			for j := i - 1; j >= 0; j-- {
				_ = svcs[j].Close()
			}
			return nil, fmt.Errorf("start service %d for %s: %w", i, info.capID, err)
		}
	}

	duration := time.Since(start)
	m.metrics.recordLaunch(ctx, info.capID, duration)
	m.lggr.Infow("Started capability",
		"capID", info.capID, "donID", info.donID,
		"duration", duration, "configHash", info.configHash)

	return &runningCapability{
		capID:      info.capID,
		donID:      info.donID,
		services:   svcs,
		configHash: info.configHash,
	}, nil
}

// resolveCapabilityBinary checks TOML config for a binary path override,
// falling back to the capability ID as the binary name.
func (m *localCapabilityManager) resolveCapabilityBinary(capID string) string {
	if m.localCfg != nil {
		capCfg := m.localCfg.GetCapabilityConfig(capID)
		if capCfg != nil && capCfg.BinaryPathOverride() != "" {
			m.lggr.Debugw("Using binary path override from TOML", "capID", capID, "path", capCfg.BinaryPathOverride())
			return capCfg.BinaryPathOverride()
		}
	}

	// Fall back to capability ID as binary name. This assumes the binary is
	// available on PATH or as a registered LOOP plugin.
	return capID // TODO: wrong
}

// stopRunningCapability stops all services for a running capability.
func (m *localCapabilityManager) stopRunningCapability(rc *runningCapability) error {
	var firstErr error
	for _, svc := range rc.services {
		if err := svc.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close service for %s: %w", rc.capID, err)
		}
	}
	return firstErr
}

// configHash computes a SHA-256 hash of the raw config bytes for change detection.
func configHash(configBytes []byte) string {
	if len(configBytes) == 0 {
		return ""
	}
	h := sha256.Sum256(configBytes)
	return hex.EncodeToString(h[:])
}
