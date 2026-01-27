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
	resolved, err := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
	if err != nil {
		return nil, err
	}
	return resolved.Values, nil
}
