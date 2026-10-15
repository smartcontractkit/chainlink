package localcapmgr

import (
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// Divergence kinds, per allowlisted capability of this node, between the offchain registry and
// the on-chain one. Telemetry only; a non-zero count means turning the gate on changes config.
const (
	// divergenceMissing: the offchain registry has no entry for the capability on its DON.
	divergenceMissing = "missing"
	// divergenceMismatch: the offchain spec_config differs from the on-chain one.
	divergenceMismatch = "mismatch"
)

// offchainDivergences counts this node's allowlisted capabilities per divergence kind.
func (m *localCapabilityManager) offchainDivergences(offchainDONs map[uint32]*capabilitiespb.DONConfig, allMyDONs []registry.DON) map[string]int64 {
	out := map[string]int64{divergenceMissing: 0, divergenceMismatch: 0}
	for _, don := range allMyDONs {
		offDON := offchainDONs[don.ID]
		for capID, onchain := range don.CapabilityConfigurations {
			if m.localCfg == nil || !m.localCfg.IsAllowlisted(capID) {
				continue
			}
			offCap, ok := offDON.GetCapabilities()[capID]
			if !ok {
				out[divergenceMissing]++
				m.lggr.Debugw("Offchain registry has no config for capability", "donID", don.ID, "capID", capID)
				continue
			}
			info := &capabilityInfo{capID: capID, donID: don.ID, config: onchain}
			legacy, err := m.buildConfigJSON(info)
			if err != nil {
				continue
			}
			info.useOffchain = true
			info.offchainSpecConfig = m.offchainSpecConfig(offDON, capID, don.ID)
			cutover, _ := m.buildConfigJSON(info)
			// A mismatch means turning the gate on changes this capability: either the launch
			// config (spec_config) or the don2don method configs would differ.
			methodsDiffer := offchainMethodConfigsDiffer(onchain, offCap)
			if cutover != legacy || methodsDiffer {
				out[divergenceMismatch]++
				m.lggr.Debugw("Offchain config differs from on-chain", "donID", don.ID, "capID", capID, "specConfigDiffers", cutover != legacy, "methodConfigsDiffer", methodsDiffer)
			}
		}
	}
	return out
}

// offchainMethodConfigsDiffer reports whether the offchain method_configs would change the
// don2don method configs the launcher wires for this capability, relative to the on-chain
// method configs. A method present in both differs only when its converted config differs;
// methods the offchain payload omits keep their on-chain config (not a difference); methods
// only present offchain are additions (a difference).
func offchainMethodConfigsDiffer(onchain registry.CapabilityConfiguration, offCap *capabilitiespb.CapabilityConfig) bool {
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

// offchainSpecConfig returns the offchain spec_config of capID on don, or nil when there is none.
// A spec_config that cannot be converted is rejected at ingestion (globalconfig.Validate); if one
// gets here anyway, it is treated as empty.
func (m *localCapabilityManager) offchainSpecConfig(don *capabilitiespb.DONConfig, capID string, donID uint32) map[string]any {
	out, err := globalconfig.SpecConfigMap(don.GetCapabilities()[capID].GetSpecConfig())
	if err != nil {
		m.lggr.Warnw("Invalid offchain spec_config, ignoring offchain override", "capID", capID, "donID", donID, "error", err)
		return nil
	}

	return out
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

// offchainDONNames maps this node's on-chain DON IDs to the DON names that key the offchain
// registry. A DON without a name (e.g. from a registry version that does not record names), or
// whose name is shared with another of this node's DONs, is left out: its offchain config
// cannot be attributed unambiguously, so with the gate on it gets an empty spec config, and it
// is reported as missing by the cross-check.
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
