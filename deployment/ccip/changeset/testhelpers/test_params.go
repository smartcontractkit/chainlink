package testhelpers

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

const (
	TestNodeOperator = "NodeOperator"
)

func NewTestRMNStaticConfig() rmn_home.RMNHomeStaticConfig {
	return rmn_home.RMNHomeStaticConfig{
		Nodes:          []rmn_home.RMNHomeNode{},
		OffchainConfig: []byte("static config"),
	}
}

func NewTestRMNDynamicConfig() rmn_home.RMNHomeDynamicConfig {
	return rmn_home.RMNHomeDynamicConfig{
		SourceChains:   []rmn_home.RMNHomeSourceChain{},
		OffchainConfig: []byte("dynamic config"),
	}
}

func NewTestNodeOperator(admin common.Address) []capabilities_registry.CapabilitiesRegistryNodeOperator {
	return []capabilities_registry.CapabilitiesRegistryNodeOperator{
		{
			Admin: admin,
			Name:  TestNodeOperator,
		},
	}
}
