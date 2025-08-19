package chainlevel

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func JobName(chainID uint64, flag cre.CapabilityFlag) string {
	return fmt.Sprintf("%s-%d", flag, chainID)
}

func IsEnabled(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// Check if this capability is enabled for any chains on this DON
	if donWithMetadata == nil || nodeSet == nil || nodeSet.ChainCapabilities == nil {
		return false
	}

	chainCapConfig, ok := nodeSet.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil || len(chainCapConfig.EnabledChains) == 0 {
		return false
	}

	return true
}

func EnabledChains(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	chainCapConfig, ok := nodeSetInput.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil {
		return []uint64{}
	}

	return chainCapConfig.EnabledChains
}

func ConfigResolver(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
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
