package donlevel

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

var ConfigMerger = func(flag cre.CapabilityFlag, nodeSetInput *cre.CapabilitiesAwareNodeSet, _ uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error) {
	// Merge global defaults with DON-specific overrides
	if nodeSetInput == nil {
		return nil, false, nil
	}

	return config.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSetInput.CapabilityOverrides), true, nil
}

var CapabilityEnabler = func(don *cre.DON, flag cre.CapabilityFlag) bool {
	if don == nil {
		return false
	}
	return don.HasFlag(flag)
}

var EnabledChainsProvider = func(donTopology *cre.DonTopology, _ *cre.CapabilitiesAwareNodeSet, _ cre.CapabilityFlag) ([]uint64, error) {

	chain, ok := chainselectors.ChainBySelector(donTopology.HomeChainSelector)
	if !ok {
		return nil, fmt.Errorf("chain for selector '%d' not found", donTopology.HomeChainSelector)
	}

	return []uint64{chain.EvmChainID}, nil
}
