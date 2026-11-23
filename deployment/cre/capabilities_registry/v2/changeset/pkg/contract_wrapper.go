package pkg

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

// PageSize is the number of entries fetched per CapabilitiesRegistry view call.
const PageSize = 30

// paginate calls fetch with consecutive (start, limit) pages until a page comes back short.
func paginate[T any](fetch func(start, limit *big.Int) ([]T, error)) ([]T, error) {
	var all []T
	limit := big.NewInt(PageSize)
	for start := int64(0); ; start += PageSize {
		page, err := fetch(big.NewInt(start), limit)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if len(page) < PageSize {
			return all, nil
		}
	}
}

func GetCapabilities(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityInfo, error) {
	caps, err := paginate(func(start, limit *big.Int) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityInfo, error) {
		return capReg.GetCapabilities(opts, start, limit)
	})
	return caps, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodeOperators(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorInfo, error) {
	nops, err := paginate(func(start, limit *big.Int) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorInfo, error) {
		return capReg.GetNodeOperators(opts, start, limit)
	})
	return nops, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodes(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.INodeInfoProviderNodeInfo, error) {
	nodes, err := paginate(func(start, limit *big.Int) ([]capabilities_registry_v2.INodeInfoProviderNodeInfo, error) {
		return capReg.GetNodes(opts, start, limit)
	})
	return nodes, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONs(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
	dons, err := paginate(func(start, limit *big.Int) ([]capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
		return capReg.GetDONs(opts, start, limit)
	})
	return dons, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONsInFamily(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry, family string) ([]*big.Int, error) {
	donIDs, err := paginate(func(start, limit *big.Int) ([]*big.Int, error) {
		return capReg.GetDONsInFamily(opts, family, start, limit)
	})
	return donIDs, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}
