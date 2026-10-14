package localcapmgr

import (
	"sort"

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
	// to the legacy TOML + on-chain config (i.e. turning the gate on changes this capability).
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

// computeOffchainCrossCheck compares a parsed offchain registry against the on-chain DON set
// and returns the divergence summary. It is pure (no metrics, no mutation); Reconcile records it.
//
// This is detection only: it never affects desired state. Only allowlisted capabilities are
// compared, matching buildDesiredState.
func (m *localCapabilityManager) computeOffchainCrossCheck(reg *capabilitiespb.OffchainCapabilitiesRegistry, version uint64, allMyDONs []registry.DON) offchainCrossCheck {
	check := offchainCrossCheck{version: version, divergences: map[string]int64{}}

	if reg == nil {
		check.offchainEmpty = true
		return check
	}

	offchainDONs := reg.GetDons()
	donNames := m.offchainDONNames(allMyDONs)

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
					m.lggr.Debugw("Offchain spec_config differs from legacy (TOML + on-chain) config",
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

// offchainChangesConfig reports whether applying offCap's spec_config on top of the legacy
// (TOML + on-chain) config would change the config the capability is launched with. Keys the
// offchain payload sets to the legacy value, and keys it omits, are not differences.
func (m *localCapabilityManager) offchainChangesConfig(donID uint32, capID string, onchain registry.CapabilityConfiguration, offCap *capabilitiespb.CapabilityConfig) bool {
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

// offchainDONsByID maps this node's on-chain DON IDs to their offchain config. DONs without a
// usable name (see offchainDONNames) or without an offchain entry are left out.
func (m *localCapabilityManager) offchainDONsByID(reg *capabilitiespb.OffchainCapabilitiesRegistry, dons []registry.DON) map[uint32]*capabilitiespb.DONConfig {
	if reg == nil {
		return nil
	}
	out := map[uint32]*capabilitiespb.DONConfig{}
	for id, name := range m.offchainDONNames(dons) {
		if d, ok := reg.GetDons()[name]; ok {
			out[id] = d
		}
	}
	return out
}

// offchainOverrides returns the offchain spec_config of capID on don, or nil when there is none.
// A spec_config that cannot be converted is rejected at ingestion (globalconfig.Validate); if one
// gets here anyway, it is ignored so the capability keeps its on-chain config.
func (m *localCapabilityManager) offchainOverrides(don *capabilitiespb.DONConfig, capID string, donID uint32) map[string]any {
	out, err := globalconfig.SpecConfigMap(don.GetCapabilities()[capID].GetSpecConfig())
	if err != nil {
		m.lggr.Warnw("Invalid offchain spec_config, ignoring offchain override", "capID", capID, "donID", donID, "error", err)
		return nil
	}
	return out
}

// offchainDONNames maps this node's on-chain DON IDs to the DON names that key the offchain
// registry. A DON without a name (e.g. from a registry version that does not record names), or
// whose name is shared with another of this node's DONs, is left out: its offchain config
// cannot be attributed unambiguously, so it keeps its legacy (on-chain/TOML) config and is
// reported as missing_don by the cross-check.
func (m *localCapabilityManager) offchainDONNames(dons []registry.DON) map[uint32]string {
	byName := make(map[string][]uint32, len(dons))
	for _, don := range dons {
		if don.Name == "" {
			m.lggr.Debugw("On-chain DON has no name; offchain config cannot apply to it", "donID", don.ID)
			continue
		}
		byName[don.Name] = append(byName[don.Name], don.ID)
	}
	out := make(map[uint32]string, len(dons))
	for name, ids := range byName {
		if len(ids) > 1 {
			m.lggr.Warnw("On-chain DON name is not unique among this node's DONs; offchain config cannot apply to them",
				"donName", name, "donIDs", ids)
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
