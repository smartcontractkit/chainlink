package shared

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"
)

var CCIPCapabilityID = utils.Keccak256Fixed(MustABIEncode(`[{"type": "string"}, {"type": "string"}]`, CapabilityLabelledName, CapabilityVersion))

func GetCCIPDonsFromCapRegistry(ctx context.Context, capRegistry *capabilities_registry.CapabilitiesRegistry) ([]capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	if capRegistry == nil {
		return nil, nil
	}
	// Get the all Dons from the capabilities registry
	allDons, err := capRegistry.GetDONs(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get all Dons from capabilities registry: %w", err)
	}
	ccipDons := make([]capabilities_registry.CapabilitiesRegistryDONInfo, 0, len(allDons))
	for _, don := range allDons {
		for _, capConfig := range don.CapabilityConfigurations {
			if capConfig.CapabilityId == CCIPCapabilityID {
				ccipDons = append(ccipDons, don)
				break
			}
		}
	}

	return ccipDons, nil
}

func MustABIEncode(abiString string, args ...interface{}) []byte {
	encoded, err := utils.ABIEncode(abiString, args...)
	if err != nil {
		panic(err)
	}
	return encoded
}
