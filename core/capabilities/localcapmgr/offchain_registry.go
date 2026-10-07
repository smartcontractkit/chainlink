package localcapmgr

import (
	"context"
	"reflect"
	"sort"

	capabilities "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// Divergence kinds recorded when the offchain capabilities registry disagrees with the
// on-chain registry. They are emitted as telemetry only and never affect which config is
// applied (that is decided by the UseOffchainRegistry gate). "missing" means the on-chain registry has something the offchain payload lacks;
// "extra" means the offchain payload describes something not present on-chain for this node.
const (
	divergenceMissingDON        = "missing_don"
	divergenceMissingCapability = "missing_capability"
	divergenceExtraDON          = "extra_don"
	divergenceExtraCapability   = "extra_capability"
	// divergenceConfigMismatch means the capability is present in both registries but its
	// offchain spec_config would change the config the capability is launched with, relative
	// to the on-chain config (i.e. turning the gate on changes this capability).
	divergenceConfigMismatch = "config_mismatch"
)

// offchainCrossCheck is the outcome of comparing the parsed offchain registry against the
// on-chain DON set this node belongs to.
type offchainCrossCheck struct {
	version       uint64
	divergences   map[string]int64 // kind -> count
	matchedCaps   int64
	comparedDONs  int64
	offchainEmpty bool // true when no offchain payload has been applied yet
}

// crossValidateOffchain compares a snapshot of the applied offchain capabilities registry
// against the on-chain DON set and emits telemetry.
//
// This is detection only, not enforcement: it never mutates desired state and never blocks or
// gates a capability. It compares presence (DON/capability keys) and, for capabilities present
// in both, whether the offchain spec_config changes the effective launch config
// (config_mismatch). Whether the offchain spec_config is applied is decided solely by the
// UseOffchainRegistry gate; a divergence recorded here does not prevent it. Malformed payloads
// are rejected earlier, at ingestion (globalconfig.Validate in the cresettings job).
//
// Only allowlisted capabilities are compared, matching buildDesiredState: those are the only
// capabilities this node would run, so they are the only ones whose config matters here.
func (m *localCapabilityManager) crossValidateOffchain(ctx context.Context, allMyDONs []registry.DON) {
	if m.offchainRegistry == nil {
		return // offchain registry not wired (e.g. feature off, or tests)
	}
	reg, version := m.offchainRegistry.LoadParsed()
	m.recordOffchainCheck(ctx, m.computeOffchainCrossCheck(reg, version, allMyDONs))
}

// computeOffchainCrossCheck compares a parsed offchain registry against the on-chain DON set
// and returns the divergence summary. It is pure (no metrics, no mutation) so it can be unit
// tested; crossValidateOffchain wires it to telemetry.
func (m *localCapabilityManager) computeOffchainCrossCheck(reg *capabilitiespb.OffchainCapabilitiesRegistry, version uint64, allMyDONs []registry.DON) offchainCrossCheck {
	check := offchainCrossCheck{version: version, divergences: map[string]int64{}}

	if reg == nil {
		check.offchainEmpty = true
		return check
	}

	offchainDONs := reg.GetDons()
	donNames := OffchainDONNames(allMyDONs)

	// on-chain -> offchain: every allowlisted (DON, capability) on-chain should be present
	// offchain.
	for _, don := range allMyDONs {
		allowlisted := m.allowlistedCapIDs(don)
		if len(allowlisted) == 0 {
			continue
		}
		check.comparedDONs++

		donName, named := donNames[don.ID]
		offDON, ok := offchainDONs[donName]
		if !named || !ok {
			check.divergences[divergenceMissingDON]++
			m.lggr.Warnw("Offchain registry missing on-chain DON", "donID", don.ID, "donName", don.Name, "offchainVersion", version)
			continue
		}

		offCaps := offDON.GetCapabilities()
		for _, capID := range allowlisted {
			if offCap, ok := offCaps[capID]; ok {
				check.matchedCaps++
				if m.offchainChangesConfig(don.ID, capID, don.CapabilityConfigurations[capID], offCap) {
					check.divergences[divergenceConfigMismatch]++
					m.lggr.Debugw("Offchain spec_config differs from on-chain config",
						"donID", don.ID, "capID", capID, "offchainVersion", version)
				}
			} else {
				check.divergences[divergenceMissingCapability]++
				m.lggr.Warnw("Offchain registry missing on-chain capability",
					"donID", don.ID, "donName", donName, "capID", capID, "offchainVersion", version)
			}
		}
	}

	// offchain -> on-chain: offchain DONs/capabilities not present on-chain for this node.
	onchainDONs := map[string]map[string]struct{}{}
	for _, don := range allMyDONs {
		name, ok := donNames[don.ID]
		if !ok {
			continue
		}
		caps := map[string]struct{}{}
		for capID := range don.CapabilityConfigurations {
			caps[capID] = struct{}{}
		}
		onchainDONs[name] = caps
	}
	for donName, offDON := range offchainDONs {
		onCaps, ok := onchainDONs[donName]
		if !ok {
			check.divergences[divergenceExtraDON]++
			continue
		}
		for capID := range offDON.GetCapabilities() {
			if _, ok := onCaps[capID]; !ok {
				check.divergences[divergenceExtraCapability]++
			}
		}
	}

	return check
}

// offchainChangesConfig reports whether applying offCap on top of the on-chain
// config would change the config the capability is launched with. It covers both config
// surfaces the offchain registry carries: spec_config (compared through the effective launch
// config) and method_configs (compared structurally against the on-chain method configs).
// Keys the offchain payload sets to the legacy value, and keys it omits, are not differences.
func (m *localCapabilityManager) offchainChangesConfig(donID uint32, capID string, onchain registry.CapabilityConfiguration, offCap *capabilitiespb.CapabilityConfig) bool {
	if m.offchainChangesSpecConfig(donID, capID, onchain, offCap) {
		return true
	}
	return m.offchainChangesMethodConfigs(donID, capID, onchain, offCap)
}

// offchainChangesSpecConfig reports whether applying offCap's spec_config on top of the legacy
// on-chain config would change the config the capability is launched with. Keys the
// offchain payload sets to the legacy value, and keys it omits, are not differences.
func (m *localCapabilityManager) offchainChangesSpecConfig(donID uint32, capID string, onchain registry.CapabilityConfiguration, offCap *capabilitiespb.CapabilityConfig) bool {
	overrides, err := globalconfig.SpecConfigMap(offCap.GetSpecConfig())
	if err != nil || len(overrides) == 0 {
		return false
	}
	info := &capabilityInfo{capID: capID, donID: donID, config: onchain}
	legacy, err := m.buildConfigJSON(info)
	if err != nil {
		return false
	}
	info.offchainOverrides = overrides
	cutover, err := m.buildConfigJSON(info)
	if err != nil {
		return false
	}
	return legacy != cutover
}

// offchainChangesMethodConfigs reports whether the offchain method_configs would change the
// don2don method configs the launcher wires for this capability, relative to the on-chain
// method configs. A method present in both is a difference only when its converted config
// differs; methods the offchain payload omits keep their on-chain config (not differences);
// methods only present offchain are additions (differences).
func (m *localCapabilityManager) offchainChangesMethodConfigs(donID uint32, capID string, onchain registry.CapabilityConfiguration, offCap *capabilitiespb.CapabilityConfig) bool {
	offMethodConfigs, err := globalconfig.MethodConfigsFromProto(offCap.GetMethodConfigs())
	if err != nil || len(offMethodConfigs) == 0 {
		return false
	}
	onchainCfg, err := onchain.Unmarshal()
	if err != nil {
		return false
	}
	onMethodConfigs := onchainCfg.CapabilityMethodConfig
	for method, offCfg := range offMethodConfigs {
		if onCfg, ok := onMethodConfigs[method]; !ok || !methodConfigEqual(onCfg, offCfg) {
			return true
		}
	}
	return false
}

// methodConfigEqual compares two method configs for equality on the fields the launcher
// consumes (remote trigger/executable config and aggregator config).
func methodConfigEqual(a, b capabilities.CapabilityMethodConfig) bool {
	if !reflect.DeepEqual(a.RemoteTriggerConfig, b.RemoteTriggerConfig) {
		return false
	}
	if !reflect.DeepEqual(a.RemoteExecutableConfig, b.RemoteExecutableConfig) {
		return false
	}
	return reflect.DeepEqual(a.AggregatorConfig, b.AggregatorConfig)
}

// OffchainDONNames maps this node's on-chain DON IDs to the DON names that key the offchain
// registry. A DON without a name (e.g. from a registry version that does not record names), or
// whose name is shared with another of this node's DONs, is left out: its offchain config
// cannot be attributed unambiguously, so it keeps its on-chain config and is
// reported as missing_don by the cross-check.
//
// Exported for the capabilities launcher, which resolves offchain method_configs by DON name
// the same way (see the offchain method_configs design note).
func OffchainDONNames(dons []registry.DON) map[uint32]string {
	byName := make(map[string][]uint32, len(dons))
	for _, don := range dons {
		if don.Name == "" {
			continue
		}
		byName[don.Name] = append(byName[don.Name], don.ID)
	}
	out := make(map[uint32]string, len(dons))
	for name, ids := range byName {
		if len(ids) > 1 {
			continue
		}
		out[ids[0]] = name
	}
	return out
}

// allowlistedCapIDs returns the capability IDs configured on a DON that this node is
// allowlisted to launch, mirroring buildDesiredState's filter.
func (m *localCapabilityManager) allowlistedCapIDs(don registry.DON) []string {
	var out []string
	for capID := range don.CapabilityConfigurations {
		if m.localCfg != nil && m.localCfg.IsAllowlisted(capID) {
			out = append(out, capID)
		}
	}
	sort.Strings(out)
	return out
}
