package localcapmgr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities/conversions"
)

// LocalCapabilityManager handles the full lifecycle of capabilities this node hosts.
// It starts, stops, and reconfigures local capabilities based on the on-chain registry state.
// It depends on standardcapabilities/delegate.go to build all dependencies for the capabilities.
type LocalCapabilityManager interface {
	services.Service

	// Called by Launcher.OnNewRegistry() for each registry update.
	Reconcile(ctx context.Context, allMyDONs []registrysyncer.DON) error
}

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

func runningKey(capID string, donID uint32) string {
	return fmt.Sprintf("%s:%d", capID, donID)
}

type localCapabilityManager struct {
	services.StateMachine
	lggr logger.Logger

	localCfg      config.LocalCapabilities
	newServicesFn NewServicesFn

	runningCapabilities map[string]*runningCapability
	mu                  sync.RWMutex

	metrics *metrics
}

// Wraps standardcapabilities.Delegate.NewServices to avoid direct dependency on the Delegate.
type NewServicesFn func(ctx context.Context, capID string, command string, configJSON string) ([]job.ServiceCtx, error)

func NewLocalCapabilityManager(lggr logger.Logger, localCfg config.LocalCapabilities, newServicesFn NewServicesFn) (LocalCapabilityManager, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create local capability manager metrics: %w", err)
	}
	return &localCapabilityManager{
		lggr:                logger.Named(lggr, "LocalCapabilityManager"),
		localCfg:            localCfg,
		newServicesFn:       newServicesFn,
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

		var errs []error
		for key, rc := range m.runningCapabilities {
			m.lggr.Infow("Stopping capability on shutdown", "capID", rc.capID, "donID", rc.donID)
			if err := m.closeServices(rc); err != nil {
				m.lggr.Errorw("Failed to stop capability on shutdown", "key", key, "error", err)
				errs = append(errs, err)
			}
		}
		m.runningCapabilities = make(map[string]*runningCapability)
		m.lggr.Info("LocalCapabilityManager stopped")
		return errors.Join(errs...)
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
	allMyDONs []registrysyncer.DON,
) error {
	desired := m.buildDesiredState(allMyDONs)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop capabilities that should no longer be running.
	for key, rc := range m.runningCapabilities {
		if _, ok := desired[key]; !ok {
			m.lggr.Infow("Stopping removed capability", "capID", rc.capID, "donID", rc.donID)
			if err := m.closeServices(rc); err != nil {
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
			if err := m.closeServices(existing); err != nil {
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

func (m *localCapabilityManager) startCapability(ctx context.Context, info *capabilityInfo) (*runningCapability, error) {
	start := time.Now()

	command := m.resolveCapabilityBinary(info.capID)
	if command == "" {
		return nil, fmt.Errorf("could not resolve capability binary for %s", info.capID)
	}
	configJSON, err := m.buildConfigJSON(info)
	if err != nil {
		return nil, fmt.Errorf("build config for %s: %w", info.capID, err)
	}

	// TODO(CRE-1775): pass also Ocr3Configs and OracleFactoryConfigs if present onchain
	svcs, err := m.newServicesFn(ctx, info.capID, command, configJSON)
	if err != nil {
		return nil, fmt.Errorf("build services for %s: %w", info.capID, err)
	}

	for i, svc := range svcs {
		if err := svc.Start(ctx); err != nil {
			for j := i - 1; j >= 0; j-- {
				if closeErr := svcs[j].Close(); closeErr != nil {
					m.lggr.Errorw("Failed to close service during rollback", "capID", info.capID, "serviceIndex", j, "error", closeErr)
				}
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

func (m *localCapabilityManager) resolveCapabilityBinary(capID string) string {
	if m.localCfg != nil {
		capCfg := m.localCfg.GetCapabilityConfig(capID)
		if capCfg != nil && capCfg.BinaryPathOverride() != "" {
			m.lggr.Debugw("Using binary path override from TOML", "capID", capID, "path", capCfg.BinaryPathOverride())
			return capCfg.BinaryPathOverride()
		}
	}

	// fall back to default command based on capability ID
	return conversions.GetCommandFromCapabilityID(capID)
}

// buildConfigJSON merges the node-local TOML config with the onchain SpecConfig
// into a flat JSON object. Onchain values take precedence over local ones.
func (m *localCapabilityManager) buildConfigJSON(info *capabilityInfo) (string, error) {
	merged := make(map[string]any)

	if m.localCfg != nil {
		capCfg := m.localCfg.GetCapabilityConfig(info.capID)
		if capCfg != nil {
			for k, v := range capCfg.Config() {
				merged[k] = v
			}
		}
	}

	if len(info.config.Config) > 0 {
		capCfg, err := info.config.Unmarshal()
		if err != nil {
			m.lggr.Warnw("Failed to unmarshal onchain config, using local config only",
				"capID", info.capID, "error", err)
		} else if capCfg.SpecConfig != nil {
			unwrapped, err := capCfg.SpecConfig.Unwrap()
			if err != nil {
				return "", fmt.Errorf("unwrap onchain spec config for %s: %w", info.capID, err)
			}
			if onchain, ok := unwrapped.(map[string]any); ok {
				maps.Copy(merged, onchain)
			}
		}
	}

	if len(merged) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(merged)
	if err != nil {
		return "", fmt.Errorf("marshal merged config for %s: %w", info.capID, err)
	}
	return string(b), nil
}

func (m *localCapabilityManager) closeServices(rc *runningCapability) error {
	return services.MultiCloser(rc.services).Close()
}

func configHash(configBytes []byte) string {
	if len(configBytes) == 0 {
		return ""
	}
	h := sha256.Sum256(configBytes)
	return hex.EncodeToString(h[:])
}
