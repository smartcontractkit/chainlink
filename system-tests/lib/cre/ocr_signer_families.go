package cre

import (
	"slices"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

// NonEVMOCRSignerFamilies returns the non-EVM OCR signer families present in the
// environment. EVM is always included separately by OCR config generation.
func NonEVMOCRSignerFamilies(blockchains []blockchains.Blockchain) []string {
	familiesSet := make(map[string]struct{})
	for _, blockchain := range blockchains {
		switch {
		case blockchain.IsFamily(chainselectors.FamilyAptos):
			familiesSet[chainselectors.FamilyAptos] = struct{}{}
		case blockchain.IsFamily(chainselectors.FamilySolana):
			familiesSet[chainselectors.FamilySolana] = struct{}{}
		}
	}

	families := make([]string, 0, len(familiesSet))
	for family := range familiesSet {
		families = append(families, family)
	}
	slices.Sort(families)

	return families
}

func CapabilityToExtraSignerFamilies(families []string, labelledNames ...string) map[string][]string {
	if len(families) == 0 || len(labelledNames) == 0 {
		return nil
	}

	capabilityToFamilies := make(map[string][]string, len(labelledNames))
	for _, labelledName := range labelledNames {
		capabilityToFamilies[labelledName] = append([]string(nil), families...)
	}

	return capabilityToFamilies
}
