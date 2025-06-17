package chainwriter

import (
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// Deprecated: use capabilities.writeevm.WriteEVMCapabilityFactory instead
var ChainWriterCapabilityFactory = func(chainID uint64) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		fullName := corevm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		if flags.HasFlag(donFlags, types.WriteEVMCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   splitName[0],
					Version:        splitName[1],
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}
