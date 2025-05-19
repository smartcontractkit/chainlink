package internal

import kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

// TODO: KS-457 configuration management for capabilities from external sources
var StreamTriggerCap = kcr.CapabilitiesRegistryCapability{
	LabelledName:   "streams-trigger",
	Version:        "1.0.0",
	CapabilityType: uint8(0), // trigger
}

var WriteChainCap = kcr.CapabilitiesRegistryCapability{
	LabelledName:   "write_ethereum-testnet-sepolia",
	Version:        "1.0.0",
	CapabilityType: uint8(3), // target
}

var OCR3Cap = kcr.CapabilitiesRegistryCapability{
	LabelledName:   "offchain_reporting",
	Version:        "1.0.0",
	CapabilityType: uint8(2), // consensus
}
