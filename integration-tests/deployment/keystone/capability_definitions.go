package keystone

import kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

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

var DonToCapabilities = map[string][]kcr.CapabilitiesRegistryCapability{
	WFDonName:     []kcr.CapabilitiesRegistryCapability{OCR3Cap},
	TargetDonName: []kcr.CapabilitiesRegistryCapability{WriteChainCap},
	StreamDonName: []kcr.CapabilitiesRegistryCapability{StreamTriggerCap},
}

var (
	WFDonName     = "wf"
	TargetDonName = "target"
	StreamDonName = "streams"
)
