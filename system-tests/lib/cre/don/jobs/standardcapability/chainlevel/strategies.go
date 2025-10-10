package chainlevel

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func JobNamer(chainID uint64, flag cre.CapabilityFlag) string {
	return fmt.Sprintf("%s-%d", flag, chainID)
}

func CapabilityEnabler(_ []string, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// for chain-level capabilities, we need to check which chains the capability is enabled for
	if nodeSet == nil || nodeSet.ChainCapabilities == nil {
		return false
	}

	chainCapConfig, ok := nodeSet.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil || len(chainCapConfig.EnabledChains) == 0 {
		return false
	}

	return true
}

func EnabledChainsProvider(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	// for chain-level capabilities, we need to return the list of chains the capability is enabled for
	chainCapConfig, ok := nodeSetInput.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil {
		return []uint64{}
	}

	return chainCapConfig.EnabledChains
}

func ConfigResolver(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
	// chain-level capabilities can have per-chain configuration overrides, we need to resolve the config for the given chain
	enabled, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
		flag,
		nodeSetInput.ChainCapabilities,
		capabilityConfig.Config,
		chainID,
	)
	if rErr != nil {
		return false, nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
	}
	if !enabled {
		return false, nil, errors.New("capability not enabled for chain")
	}

	return true, mergedConfig, nil
}
