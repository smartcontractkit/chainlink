package flags

import (
	"slices"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func DonMetadataWithFlag(donTopologies []*types.DonMetadata, flag string) []*types.DonMetadata {
	var result []*types.DonMetadata

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

func HasOnlyOneFlag(values []string, flag string) bool {
	return slices.Contains(values, flag) && len(values) == 1
}

func OneDonMetadataWithFlag(donTopologies []*types.DonMetadata, flag string) (*types.DonMetadata, error) {
	donTopologies = DonMetadataWithFlag(donTopologies, flag)
	if len(donTopologies) != 1 {
		return nil, errors.Errorf("expected exactly one DON topology with flag %s, got %d", flag, len(donTopologies))
	}

	return donTopologies[0], nil
}

func NodeSetFlags(nodeSet *types.CapabilitiesAwareNodeSet) ([]string, error) {
	var stringCaps []string

	stringCaps = append(stringCaps, append(nodeSet.Capabilities, nodeSet.DONTypes...)...)
	return stringCaps, nil
}
