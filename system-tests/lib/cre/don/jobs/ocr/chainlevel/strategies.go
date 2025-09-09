package chainlevel

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

var ConfigMerger = func(flag cre.CapabilityFlag, nodeSetInput *cre.CapabilitiesAwareNodeSet, chainIDUint64 uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error) {
	// Build user configuration from defaults + chain overrides
	enabled, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(flag, nodeSetInput.ChainCapabilities, capabilityConfig.Config, chainIDUint64)
	if rErr != nil {
		return nil, false, errors.Wrap(rErr, "failed to resolve capability config for chain")
	}
	if !enabled {
		return nil, false, nil
	}

	return mergedConfig, true, nil
}

var CapabilityEnabler = func(nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	if nodeSetInput == nil || nodeSetInput.ChainCapabilities == nil {
		return false
	}
	if cc, ok := nodeSetInput.ChainCapabilities[flag]; !ok || cc == nil || len(cc.EnabledChains) == 0 {
		return false
	}

	return true
}

var EnabledChainsProvider = func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) ([]uint64, error) {
	if nodeSetInput == nil || nodeSetInput.ChainCapabilities == nil {
		return []uint64{}, nil
	}

	return nodeSetInput.ChainCapabilities[flag].EnabledChains, nil
}
