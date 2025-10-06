package pkg

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

// TODO: we have to support pagination eventually
var (
	MaxCapabilities = big.NewInt(128)
	MaxDONs         = big.NewInt(32)
	MaxNodes        = big.NewInt(256)
	MaxNOPs         = big.NewInt(128)
)

func GetCapabilities(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityInfo, error) {
	caps, err := capReg.GetCapabilities(opts, big.NewInt(0), MaxCapabilities)
	return caps, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodeOperators(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorInfo, error) {
	nops, err := capReg.GetNodeOperators(opts, big.NewInt(0), MaxNOPs)
	return nops, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodes(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.INodeInfoProviderNodeInfo, error) {
	nodes, err := capReg.GetNodes(opts, big.NewInt(0), MaxNodes)
	return nodes, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONs(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
	donsInfo, err := capReg.GetDONs(opts, big.NewInt(0), MaxDONs)
	return donsInfo, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONsInFamily(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry, family string) ([]*big.Int, error) {
	donsInfo, err := capReg.GetDONsInFamily(opts, family, big.NewInt(0), MaxDONs)
	return donsInfo, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}
