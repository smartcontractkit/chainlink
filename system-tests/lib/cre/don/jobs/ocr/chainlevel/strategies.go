package chainlevel

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

var ConfigMerger = func(flag cre.CapabilityFlag, nodeSet cre.NodeSetWithChainCapabilities, chainIDUint64 uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error) {
	// Build user configuration from defaults + chain overrides
	enabled, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilities(), capabilityConfig.Config, chainIDUint64)
	if rErr != nil {
		return nil, false, errors.Wrap(rErr, "failed to resolve capability config for chain")
	}
	if !enabled {
		return nil, false, nil
	}

	return mergedConfig, true, nil
}

var CapabilityEnabler = func(don *cre.DON, flag cre.CapabilityFlag) bool {
	if don == nil || don.GetChainCapabilities() == nil {
		return false
	}
	if cc, ok := don.GetChainCapabilities()[flag]; !ok || cc == nil || len(cc.EnabledChains) == 0 {
		return false
	}

	return true
}

var EnabledChainsProvider = func(donTopology *cre.DonTopology, nodeSet cre.NodeSetWithChainCapabilities, flag cre.CapabilityFlag) ([]uint64, error) {
	if nodeSet == nil || nodeSet.GetChainCapabilities() == nil {
		return []uint64{}, nil
	}

	return nodeSet.GetChainCapabilities()[flag].EnabledChains, nil
}
