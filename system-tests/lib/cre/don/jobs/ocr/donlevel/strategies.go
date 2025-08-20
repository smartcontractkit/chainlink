package donlevel

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
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

var EnabledChainsProvider = func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) ([]uint64, error) {

	chain, ok := chainsel.ChainBySelector(donTopology.HomeChainSelector)
	if !ok {
		return nil, fmt.Errorf("chain for selector '%d' not found", donTopology.HomeChainSelector)
	}

	return []uint64{chain.EvmChainID}, nil
}
