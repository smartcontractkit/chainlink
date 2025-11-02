package donlevel

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

var ConfigMerger = func(flag cre.CapabilityFlag, nodeSet cre.NodeSetWithCapabilityConfigs, _ uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error) {
	// Merge global defaults with DON-specific overrides
	if nodeSet == nil {
		return nil, false, nil
	}

	return config.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSet.GetCapabilityConfigOverrides()), true, nil
}

var CapabilityEnabler = func(don *cre.Don, flag cre.CapabilityFlag) bool {
	if don == nil {
		return false
	}
	return don.HasFlag(flag)
}

var EnabledChainsProvider = func(registryChainSelector uint64, _ cre.NodeSetWithCapabilityConfigs, _ cre.CapabilityFlag) ([]uint64, error) {
	chain, ok := chainselectors.ChainBySelector(registryChainSelector)
	if !ok {
		return nil, fmt.Errorf("chain for selector '%d' not found", registryChainSelector)
	}

	return []uint64{chain.EvmChainID}, nil
}
