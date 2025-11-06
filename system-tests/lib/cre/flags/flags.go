package flags

import (
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func HasFlag(values []string, flag string) bool {
	return HasFlagForAnyChain(values, flag)
}

func HasOnlyOneFlag(values []string, flag string) bool {
	return slices.Contains(values, flag) && len(values) == 1
}

func HasFlagForChain(values []string, capability string, chainID uint64) bool {
	return slices.Contains(values, capability+"-"+strconv.FormatUint(chainID, 10))
}

func HasFlagForAnyChain(values []string, capability string) bool {
	if slices.Contains(values, capability) {
		return true
	}

	for _, value := range values {
		if strings.HasPrefix(value, capability+"-") {
			return true
		}
	}

	return false
}

func RequiresForwarderContract(values []string, chainID uint64) bool {
	return HasFlagForChain(values, cre.EVMCapability, chainID) || HasFlagForChain(values, cre.WriteEVMCapability, chainID) || HasFlagForAnyChain(values, cre.WriteSolanaCapability)
}

func DonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) []*cre.DonMetadata {
	var result []*cre.DonMetadata

	for _, donTopology := range donTopologies {
		if HasFlagForAnyChain(donTopology.Flags, flag) {
			result = append(result, donTopology)
		}
	}

	return result
}

func OneDonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) (*cre.DonMetadata, error) {
	donTopologies = DonMetadataWithFlag(donTopologies, flag)
	if len(donTopologies) != 1 {
		return nil, errors.Errorf("expected exactly one DON topology with flag %s, got %d", flag, len(donTopologies))
	}

	return donTopologies[0], nil
}
