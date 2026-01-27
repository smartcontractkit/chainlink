package chainlevel

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func JobNamer(chainID uint64, flag cre.CapabilityFlag) string {
	return fmt.Sprintf("%s-%d", flag, chainID)
}

func CapabilityEnabler(_ []string, nodeSet cre.NodeSetWithCapabilityConfigs, flag cre.CapabilityFlag) bool {
	enabledChains, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil || len(enabledChains) == 0 {
		return false
	}

	return true
}

func EnabledChainsProvider(_ uint64, nodeSet cre.NodeSetWithCapabilityConfigs, flag cre.CapabilityFlag) []uint64 {
	// for chain-level capabilities, we need to return the list of chains the capability is enabled for
	enabledChains, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil || len(enabledChains) == 0 {
		return []uint64{}
	}

	return enabledChains
}

func ConfigResolver(nodeSet cre.NodeSetWithCapabilityConfigs, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (map[string]any, error) {
	// chain-level capabilities can have per-chain configuration overrides, we need to resolve the config for the given chain
	config, ok := nodeSet.GetCapabilityConfig(cre.FlagWithChainID(flag, chainID))
	if !ok {
		return nil, fmt.Errorf("capability config not found for flag %s", cre.FlagWithChainID(flag, chainID))
	}
	return config.Values, nil
}
