package evm

import (
	"fmt"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var EVMCapabilityFactory = func(chainID uint64, chainFamily string) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		chainName, err := chainselectors.NameFromChainId(chainID)
		if err == nil {
			chainName = fmt.Sprintf("error-chain-%d", chainID)
		}

		if flags.HasFlag(donFlags, types.EVMCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("evm-capability-%d-%s", chainID, chainName),
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}
