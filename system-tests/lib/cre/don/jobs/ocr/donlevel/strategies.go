package donlevel

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var ConfigMerger = func(flag cre.CapabilityFlag, nodeSetInput *cre.CapabilitiesAwareNodeSet, chainIDUint64 uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error) {
	// Merge global defaults with DON-specific overrides
	if nodeSetInput == nil {
		return nil, false, nil
	}

	return envconfig.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSetInput.CapabilityOverrides), true, nil
}

var CapabilityEnabler = func(nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	if nodeSetInput == nil {
		return false
	}
	return flags.HasFlag(nodeSetInput.ComputedCapabilities, flag)
}

var EnabledChainsProvider = func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	return []uint64{donTopology.HomeChainSelector}
}
