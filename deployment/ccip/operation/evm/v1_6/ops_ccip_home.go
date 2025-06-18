package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type AddDONOpInput struct {
	Nodes                    [][32]byte
	CapabilityConfigurations []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
	IsPublic                 bool
	AcceptsWorkflows         bool
	F                        uint8
}

var AddDONOp = opsutil.NewEVMCallOperation(
	"AddDONOp",
	semver.MustParse("1.0.0"),
	"Adds DONs via the CapabilitiesRegistry",
	ccip_home.CCIPHomeABI,
	shared.CapabilitiesRegistry,
	capabilities_registry.NewCapabilitiesRegistry,
	func(capReg *capabilities_registry.CapabilitiesRegistry, opts *bind.TransactOpts, input AddDONOpInput) (*types.Transaction, error) {
		return capReg.AddDON(opts, input.Nodes, input.CapabilityConfigurations, input.IsPublic, input.AcceptsWorkflows, input.F)
	},
)

type UpdateDONOpInput struct {
	ID                       uint32
	Nodes                    [][32]byte
	CapabilityConfigurations []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
	IsPublic                 bool
	F                        uint8
}

var UpdateDONOp = opsutil.NewEVMCallOperation(
	"UpdateDONOp",
	semver.MustParse("1.0.0"),
	"Updates DONs via the CapabilitiesRegistry",
	ccip_home.CCIPHomeABI,
	shared.CapabilitiesRegistry,
	capabilities_registry.NewCapabilitiesRegistry,
	func(capReg *capabilities_registry.CapabilitiesRegistry, opts *bind.TransactOpts, input UpdateDONOpInput) (*types.Transaction, error) {
		return capReg.UpdateDON(opts, input.ID, input.Nodes, input.CapabilityConfigurations, input.IsPublic, input.F)
	},
)

type ApplyChainConfigUpdatesOpInput struct {
	RemoteChainRemoves []uint64
	RemoteChainAdds    []ccip_home.CCIPHomeChainConfigArgs
}

var ApplyChainConfigUpdatesOp = opsutil.NewEVMCallOperation(
	"ApplyChainConfigUpdatesOp",
	semver.MustParse("1.0.0"),
	"Updates chain configurations on CCIPHome",
	ccip_home.CCIPHomeABI,
	shared.CCIPHome,
	ccip_home.NewCCIPHome,
	func(ccipHome *ccip_home.CCIPHome, opts *bind.TransactOpts, input ApplyChainConfigUpdatesOpInput) (*types.Transaction, error) {
		return ccipHome.ApplyChainConfigUpdates(opts, input.RemoteChainRemoves, input.RemoteChainAdds)
	},
)
