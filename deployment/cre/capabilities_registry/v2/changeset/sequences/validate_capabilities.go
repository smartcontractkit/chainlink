package sequences

import (
	"fmt"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

func ValidateNoDuplicateCapabilitiesAcrossDONs(
	donCapabilityConfigs map[string][]contracts.CapabilityConfig,
	existingDONs []capabilities_registry_v2.CapabilitiesRegistryDONInfo,
	excludeDONName string,
) error {
	// capID -> DON name it is claimed by, across the input map itself.
	claimedByInput := make(map[string]string, len(donCapabilityConfigs))
	for donName, configs := range donCapabilityConfigs {
		for _, cfg := range configs {
			capID := cfg.Capability.CapabilityID
			if otherDON, ok := claimedByInput[capID]; ok && otherDON != donName {
				return fmt.Errorf(
					"capability %q is assigned to both DON %q and DON %q; a capability can only be assigned to one DON",
					capID, otherDON, donName)
			}
			claimedByInput[capID] = donName
		}
	}

	for _, don := range existingDONs {
		if don.Name == excludeDONName {
			continue
		}
		for _, cfg := range don.CapabilityConfigurations {
			donName, ok := claimedByInput[cfg.CapabilityId]
			if !ok {
				continue
			}
			return fmt.Errorf(
				"capability %q is already assigned to on-chain DON %q; cannot also assign it to DON %q. "+
					"A capability can only be assigned to one DON",
				cfg.CapabilityId, don.Name, donName)
		}
	}

	return nil
}
