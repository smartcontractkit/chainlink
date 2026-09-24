package sequences

import (
	"fmt"
	"slices"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

// ValidateNoDuplicateCapabilitiesAcrossDONs ensures a capability is not assigned to two
// different DONs within the same DON family. DONs in disjoint families (e.g. zone-a and
// zone-b) may hold the same capability.
func ValidateNoDuplicateCapabilitiesAcrossDONs(
	donCapabilityConfigs map[string][]contracts.CapabilityConfig,
	existingDONs []capabilities_registry_v2.CapabilitiesRegistryDONInfo,
) error {
	existingByName := make(map[string]capabilities_registry_v2.CapabilitiesRegistryDONInfo, len(existingDONs))
	for _, don := range existingDONs {
		existingByName[don.Name] = don
	}

	// capID -> DON name it is claimed by, across the input map itself.
	claimedByInput := make(map[string]string, len(donCapabilityConfigs))
	for donName, configs := range donCapabilityConfigs {
		for _, cfg := range configs {
			capID := cfg.Capability.CapabilityID
			otherDON, ok := claimedByInput[capID]
			if !ok {
				claimedByInput[capID] = donName
				continue
			}
			if otherDON == donName {
				continue
			}
			if !donsShareFamily(otherDON, donName, existingByName) {
				continue
			}
			return fmt.Errorf(
				"capability %q is assigned to both DON %q and DON %q, which share a DON family; "+
					"a capability can only be assigned to one DON per family",
				capID, otherDON, donName)
		}
	}

	for _, don := range existingDONs {
		for _, cfg := range don.CapabilityConfigurations {
			donName, ok := claimedByInput[cfg.CapabilityId]
			if !ok {
				continue
			}
			if don.Name == donName {
				// Re-assigning a capability to the DON that already holds it is not a conflict.
				continue
			}
			if !donsShareFamily(don.Name, donName, existingByName) {
				continue
			}
			return fmt.Errorf(
				"capability %q is already assigned to on-chain DON %q, which shares a DON family with DON %q; "+
					"a capability can only be assigned to one DON per family",
				cfg.CapabilityId, don.Name, donName)
		}
	}

	return nil
}

// donsShareFamily reports whether two DONs share at least one DON family. A DON not found
// on-chain is conservatively treated as sharing a family with every other DON.
func donsShareFamily(
	donA, donB string,
	existingByName map[string]capabilities_registry_v2.CapabilitiesRegistryDONInfo,
) bool {
	a, okA := existingByName[donA]
	b, okB := existingByName[donB]
	if !okA || !okB {
		return true
	}
	for _, family := range a.DonFamilies {
		if slices.Contains(b.DonFamilies, family) {
			return true
		}
	}
	return false
}
