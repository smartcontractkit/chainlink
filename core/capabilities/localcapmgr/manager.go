package localcapmgr

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sync"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities/conversions"
)

// LocalCapabilityManager handles the full lifecycle of capabilities this node hosts.
// It starts, stops, and reconfigures local capabilities based on the on-chain registry state.
// It depends on standardcapabilities/delegate.go to build all dependencies for the capabilities.
type LocalCapabilityManager interface {
	services.Service

	// Called by Launcher.OnNewRegistry() for each registry update.
	Reconcile(ctx context.Context, allMyDONs []registry.DON) error
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
	capID  string
	donID  uint32
	config registry.CapabilityConfiguration
	// offchainOverrides is the offchain spec_config for (capID, donID), taken from the
	// registry snapshot of the reconcile that built this info. Nil when the cutover is off or
	// the offchain payload has no entry for this capability on this DON.
	offchainOverrides map[string]any
	// configJSON is the effective config the capability is launched with (see buildConfigJSON)
	// and configErr the error building it, if any.
	configJSON string
	configErr  error
	// configHash identifies the effective launch inputs (on-chain config bytes + effective
	// config JSON). A change in either, including an offchain-only change, restarts the
	// capability.
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
	// configProvider yields node-local (TOML) capability config overrides. It is the base layer
	// in buildConfigJSON, below the on-chain SpecConfig. Binary-path/allowlist still read
	// directly from localCfg.
	configProvider CapabilityConfigProvider
	// useOffchainRegistry is the cutover gate. When true, the offchain spec_config is applied
	// LAST in buildConfigJSON (offchain-wins over both TOML and on-chain). When false the
	// offchain layer is skipped entirely, and any key absent offchain falls back to the on-chain
	// or TOML value — the cutover is backwards compatible by construction.
	useOffchainRegistry bool
	// offchainRegistry is the offchain capabilities registry delivered via the cresettings job.
	// It is cross-validated against the on-chain registry every Reconcile (telemetry), regardless
	// of the cutover gate. Nil when the feature is not wired.
	offchainRegistry *globalconfig.GlobalConfig

	// reconcileMu serializes whole reconciles (registry-driven and offchain-update-driven), so
	// each one builds its desired state from a single offchain snapshot and the latest DON set.
	reconcileMu sync.Mutex
	// lastDONs is the DON set from the most recent registry-driven Reconcile; an offchain-only
	// update re-reconciles against it. Guarded by reconcileMu.
	lastDONs      []registry.DON
	hasReconciled bool

	stopCh services.StopChan
	wg     sync.WaitGroup

	runningCapabilities map[string]*runningCapability
	mu                  sync.RWMutex

	metrics *metrics
}

// Wraps standardcapabilities.Delegate.NewServices to avoid direct dependency on the Delegate.
// donID is the authoritative on-chain DON ID this plugin process is being spawned for; it is
// known here because Reconcile keys desired state by (capID, donID).
// ocr3Config is the on-chain OCR3 config parsed from the capability configuration, or nil when
// none is present; it lets the delegate align the node's signer/transmitter with the registry.
type NewServicesFn func(ctx context.Context, capID string, donID uint32, command string, configJSON string, ocr3Config *ocrtypes.ContractConfig) ([]job.ServiceCtx, error)

// offchainRegistry may be nil, in which case the offchain cross-validation pass is skipped.
//
// useOffchainRegistry is the cutover gate. When false (default), capability config is sourced
// from TOML/on-chain and the offchain registry is used for cross-validation telemetry only. When
// true (and offchainRegistry is non-nil), the offchain config is applied on top of the on-chain
// SpecConfig (offchain-wins), falling back to on-chain/TOML for any value the offchain payload
// does not carry. The gate makes the cutover per-node, reversible, and backwards compatible.
func NewLocalCapabilityManager(lggr logger.Logger, localCfg config.LocalCapabilities, newServicesFn NewServicesFn, offchainRegistry *globalconfig.GlobalConfig, useOffchainRegistry bool) (LocalCapabilityManager, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create local capability manager metrics: %w", err)
	}
	named := logger.Named(lggr, "LocalCapabilityManager")

	useOffchain := useOffchainRegistry && offchainRegistry != nil
	if useOffchain {
		named.Info("Offchain capabilities registry cutover ENABLED: offchain spec_config wins over on-chain and TOML")
	}

	return &localCapabilityManager{
		lggr:                named,
		localCfg:            localCfg,
		newServicesFn:       newServicesFn,
		configProvider:      tomlCapabilityConfigProvider{localCfg: localCfg},
		useOffchainRegistry: useOffchain,
		offchainRegistry:    offchainRegistry,
		stopCh:              make(services.StopChan),
		runningCapabilities: make(map[string]*runningCapability),
		metrics:             metrics,
	}, nil
}

func (m *localCapabilityManager) Start(ctx context.Context) error {
	return m.StartOnce("LocalCapabilityManager", func() error {
		if m.stopCh == nil {
			m.stopCh = make(services.StopChan)
		}
		if m.useOffchainRegistry && m.offchainRegistry != nil {
			updates, unsubscribe := m.offchainRegistry.Subscribe()
			m.wg.Go(func() {
				defer unsubscribe()
				m.watchOffchainUpdates(updates)
			})
		}
		m.lggr.Info("LocalCapabilityManager started")
		return nil
	})
}

// watchOffchainUpdates re-reconciles whenever a new offchain payload is applied, so an
// offchain-only change takes effect without waiting for the next registry sync.
func (m *localCapabilityManager) watchOffchainUpdates(updates <-chan struct{}) {
	ctx, cancel := m.stopCh.NewCtx()
	defer cancel()
	for {
		select {
		case <-m.stopCh:
			return
		case <-updates:
			m.reconcileOnOffchainUpdate(ctx)
		}
	}
}

// reconcileOnOffchainUpdate re-runs the reconcile against the last known on-chain DON set.
// Before the first registry-driven Reconcile there is no DON set to reconcile against; the
// first Reconcile reads the latest offchain snapshot anyway.
func (m *localCapabilityManager) reconcileOnOffchainUpdate(ctx context.Context) {
	m.reconcileMu.Lock()
	defer m.reconcileMu.Unlock()
	if !m.hasReconciled {
		return
	}
	m.lggr.Infow("Offchain capabilities registry updated, reconciling local capabilities")
	m.reconcileLocked(ctx, m.lastDONs)
}

func (m *localCapabilityManager) Close() error {
	return m.StopOnce("LocalCapabilityManager", func() error {
		if m.stopCh != nil {
			close(m.stopCh)
		}
		m.wg.Wait()

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
	allMyDONs []registry.DON,
) error {
	m.reconcileMu.Lock()
	defer m.reconcileMu.Unlock()
	m.lastDONs = allMyDONs
	m.hasReconciled = true
	m.reconcileLocked(ctx, allMyDONs)
	return nil
}

// reconcileLocked performs one reconcile. Callers must hold reconcileMu.
func (m *localCapabilityManager) reconcileLocked(ctx context.Context, allMyDONs []registry.DON) {
	// Take ONE snapshot of the offchain registry for the whole pass, so cross-validation and
	// every capability's effective config are computed from the same payload even if a new one
	// is stored concurrently. LoadParsed returns a deep copy, so nothing here can mutate the
	// shared applied config.
	var offchainReg *capabilitiespb.OffchainCapabilitiesRegistry
	var offchainVersion uint64
	if m.offchainRegistry != nil {
		offchainReg, offchainVersion = m.offchainRegistry.LoadParsed()
	}

	var offchain CapabilityConfigProvider
	if m.useOffchainRegistry && offchainReg != nil {
		offchain = offchainCapabilityConfigProvider{reg: offchainReg, version: offchainVersion, lggr: m.lggr}
	}
	desired := m.buildDesiredState(allMyDONs, offchain)

	// Cross-validate the offchain registry against the on-chain DON set. Telemetry only — it
	// does not influence desired state (and does not gate or block anything).
	if m.offchainRegistry != nil {
		m.recordOffchainCheck(ctx, m.computeOffchainCrossCheck(offchainReg, offchainVersion, allMyDONs))
	}

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
}

// buildDesiredState extracts capabilities that should be running from DON configs, and
// computes each one's effective config. Only includes capabilities that are in the
// RegistryBasedLaunchAllowlist. offchain is the offchain spec_config layer for this reconcile,
// or nil when the cutover is off or no offchain payload is applied.
func (m *localCapabilityManager) buildDesiredState(myCapabilityDONs []registry.DON, offchain CapabilityConfigProvider) map[string]*capabilityInfo {
	desired := make(map[string]*capabilityInfo)
	for _, don := range myCapabilityDONs {
		for capID, capCfg := range don.CapabilityConfigurations {
			if m.localCfg == nil || !m.localCfg.IsAllowlisted(capID) {
				continue
			}

			info := &capabilityInfo{
				capID:  capID,
				donID:  don.ID,
				config: capCfg,
			}
			if offchain != nil {
				info.offchainOverrides = offchain.LocalConfigOverrides(capID, don.ID)
			}
			info.configJSON, info.configErr = m.buildConfigJSON(info)
			info.configHash = effectiveConfigHash(capCfg.Config, info.configJSON)
			desired[runningKey(capID, don.ID)] = info
		}
	}
	return desired
}

func (m *localCapabilityManager) startCapability(ctx context.Context, info *capabilityInfo) (*runningCapability, error) {
	start := time.Now()

	// command is only meaningful for standard capabilities that launch a plugin
	// binary (e.g. consensus, cron). OCR2-based capabilities (e.g. dontime) run
	// in-process and do not need a binary, so an empty command is allowed there;
	// the newServicesFn routes on capability ID and ignores it.
	command := m.resolveCapabilityBinary(info.capID)
	if info.configErr != nil {
		return nil, fmt.Errorf("build config for %s: %w", info.capID, info.configErr)
	}
	configJSON := info.configJSON
	if len(info.config.Config) > 0 {
		if _, err := info.config.Unmarshal(); err != nil {
			m.lggr.Warnw("Failed to unmarshal onchain config, launching without on-chain spec config",
				"capID", info.capID, "donID", info.donID, "error", err)
		}
	}

	// TODO(CRE-1775): also derive and pass OracleFactoryConfigs if present onchain.
	ocr3Config := extractDefaultOCR3Config(info.config)
	svcs, err := m.newServicesFn(ctx, info.capID, info.donID, command, configJSON, ocr3Config)
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

// overridesFor returns the node-local (TOML) capability config base via the config provider.
// It falls back to reading TOML directly when no provider is set (e.g. managers built as struct
// literals in tests); the constructor always installs a provider in production. The offchain
// layer is applied separately (and last) in buildConfigJSON.
func (m *localCapabilityManager) overridesFor(capID string, donID uint32) map[string]any {
	if m.configProvider != nil {
		return m.configProvider.LocalConfigOverrides(capID, donID)
	}
	if m.localCfg == nil {
		return nil
	}
	if capCfg := m.localCfg.GetCapabilityConfig(capID); capCfg != nil {
		return toAnyMap(capCfg.Config())
	}
	return nil
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

// buildConfigJSON merges capability config into a flat JSON object, in increasing order of
// precedence:
//  1. node-local TOML overrides (base),
//  2. on-chain SpecConfig,
//  3. offchain SpecConfig (only when the cutover is enabled).
//
// The offchain layer is applied last so it wins over on-chain, which is what makes it a real
// cutover; because it is skipped entirely when the cutover is off, and any key the offchain
// payload omits keeps its on-chain/TOML value, the layering is backwards compatible.
func (m *localCapabilityManager) buildConfigJSON(info *capabilityInfo) (string, error) {
	merged := make(map[string]any)

	// 1. node-local TOML base.
	for k, v := range m.overridesFor(info.capID, info.donID) {
		merged[k] = v
	}

	// 2. on-chain SpecConfig.
	if len(info.config.Config) > 0 {
		capCfg, err := info.config.Unmarshal()
		if err != nil {
			// Logged at launch (startCapability); this runs on every reconcile.
			m.lggr.Debugw("Failed to unmarshal onchain config, using local config only",
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

	// 3. offchain SpecConfig wins (cutover on only; nil otherwise). Applied last, keys absent
	// offchain retain their on-chain/TOML value.
	maps.Copy(merged, info.offchainOverrides)

	if len(merged) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(merged)
	if err != nil {
		return "", fmt.Errorf("marshal merged config for %s: %w", info.capID, err)
	}
	return string(b), nil
}

// extractDefaultOCR3Config returns the `default` on-chain OCR3 config parsed from the
// capability configuration, or nil when the configuration is empty, cannot be parsed,
// or carries no OCR3 config. The delegate uses it to align the node's signer and
// transmitter with the registry.
// By `default,` we mean the config from the registry stored under the "default" key.
func extractDefaultOCR3Config(cc registry.CapabilityConfiguration) *ocrtypes.ContractConfig {
	if len(cc.Config) == 0 {
		return nil
	}
	parsed, err := cc.Unmarshal()
	if err != nil {
		return nil
	}
	cfg, ok := parsed.Ocr3Configs[capabilitiespb.OCR3ConfigDefaultKey]
	if !ok {
		return nil
	}
	return &cfg
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

// effectiveConfigHash hashes everything a capability is launched with that can change at
// runtime: the raw on-chain capability config (which also carries OCR3 configs) and the
// effective merged config JSON (TOML < on-chain < offchain). Hashing the effective JSON rather
// than only the on-chain bytes is what makes an offchain-only change restart the capability.
// buildConfigJSON marshals a map, so the JSON has sorted keys and is deterministic.
func effectiveConfigHash(onchainConfig []byte, configJSON string) string {
	h := sha256.New()
	var n [8]byte
	binary.BigEndian.PutUint64(n[:], uint64(len(onchainConfig)))
	h.Write(n[:])
	h.Write(onchainConfig)
	h.Write([]byte(configJSON))
	return hex.EncodeToString(h.Sum(nil))
}
