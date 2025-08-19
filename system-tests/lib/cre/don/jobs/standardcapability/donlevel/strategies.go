package donlevel

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func JobName(chainID uint64, flag cre.CapabilityFlag) string {
	return flag
}

func IsEnabled(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// Check if this DON has the capability enabled
	return flags.HasFlag(donWithMetadata.Flags, flag)
}

func EnabledChains(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	return []uint64{donTopology.HomeChainSelector}
}

func ConfigResolver(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, _ uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
	if nodeSetInput == nil {
		return false, nil, errors.New("node set input is nil")
	}

	return true, envconfig.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSetInput.CapabilityOverrides), nil
}
