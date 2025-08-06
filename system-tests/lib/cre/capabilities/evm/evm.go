package evm

import (
	"fmt"
	"strconv"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var EVMCapabilityFactory = func(chainID uint64, chainFamily string) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		if flags.HasFlag(donFlags, cre.EVMCapability) {
			selector, err := chainselectors.SelectorFromChainId(chainID)
			if err != nil {
				fmt.Printf("Error getting selector from chainID: %d, err: %s\n", chainID, err.Error())
				selector = 0
			}
			labelledName := "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10)
			fmt.Println("EVM Capability Labelled Name:", labelledName)
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   labelledName,
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}
