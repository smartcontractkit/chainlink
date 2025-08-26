package donlevel

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func JobNamer(chainID uint64, flag cre.CapabilityFlag) string {
	return flag
}

func CapabilityEnabler(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// for DON-level capabilities, we only need to check if the DON has the capability enabled
	return flags.HasFlag(donWithMetadata.Flags, flag)
}

func EnabledChainsProvider(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	// Most DON-level capabilities do not operate on specific chains, so we return the home chain selector to satisfy the interface
	return []uint64{donTopology.HomeChainSelector}
}

func ConfigResolver(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, _ uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
	if nodeSetInput == nil {
		return false, nil, errors.New("node set input is nil")
	}

	return true, envconfig.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSetInput.CapabilityOverrides), nil
}
