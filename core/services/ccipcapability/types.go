package ccipcapability

import (
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
)

type RegistryState struct {
	DONs         []keystone_capability_registry.CapabilityRegistryDONInfo
	Capabilities []keystone_capability_registry.CapabilityRegistryCapability
	Nodes        []keystone_capability_registry.CapabilityRegistryNodeInfo
}

type CapabilityRegistry interface {
	// LatestState returns the latest state of the on-chain capability registry.
	LatestState() (RegistryState, error)
}

type ChainConfig interface {
	Readers() [][32]byte
	FChain() uint8
}

type PluginType uint8

const (
	PluginTypeCCIPCommit PluginType = 0
	PluginTypeCCIPExec   PluginType = 1
)

type OCRConfig interface {
	PluginType() PluginType
	ChainSelector() uint64
	F() uint8
	OffchainConfigVersion() uint64
	OfframpAddress() string
	Signers() [][2][32]byte
	Transmitters() [][2][32]byte
	OffchainConfig() []byte
}
