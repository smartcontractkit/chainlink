package flags

import (
	"slices"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func DonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) []*cre.DonMetadata {
	var result []*cre.DonMetadata

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

func OneDonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) (*cre.DonMetadata, error) {
	donTopologies = DonMetadataWithFlag(donTopologies, flag)
	if len(donTopologies) != 1 {
		return nil, errors.Errorf("expected exactly one DON topology with flag %s, got %d", flag, len(donTopologies))
	}

	return donTopologies[0], nil
}

func NodeSetFlags(nodeSet *cre.CapabilitiesAwareNodeSet) ([]string, error) {
	var stringCaps []string

	stringCaps = append(stringCaps, append(nodeSet.Capabilities, nodeSet.DONTypes...)...)
	return stringCaps, nil
}
