package flags

import (
	"slices"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func DONTopologyWithFlag(donTopologies []*types.DonWithMetadata, flag string) []*types.DonWithMetadata {
	var result []*types.DonWithMetadata

	for _, donTopology := range donTopologies {
		if HasFlag(donTopology.Flags, flag) {
			result = append(result, donTopology)
		}
	}

	return result
}

func HasFlag(values []string, flag string) bool {
	return slices.Contains(values, flag)
}

func OneDONTopologyWithFlag(donTopologies []*types.DonWithMetadata, flag string) (*types.DonWithMetadata, error) {
	donTopologies = DONTopologyWithFlag(donTopologies, flag)
	if len(donTopologies) != 1 {
		return nil, errors.Errorf("expected exactly one DON topology with flag %s, got %d", flag, len(donTopologies))
	}

	return donTopologies[0], nil
}

func NodeSetFlags(nodeSet *types.CapabilitiesAwareNodeSet) ([]string, error) {
	var stringCaps []string
	if len(nodeSet.Capabilities) == 0 && nodeSet.DONType == "" {
		// if no flags are set, we assign all known capabilities to the DON
		return types.SingleDonFlags, nil
	}

	stringCaps = append(stringCaps, append(nodeSet.Capabilities, nodeSet.DONType)...)
	return stringCaps, nil
}
