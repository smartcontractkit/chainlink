package contracts

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	CapabilitiesRegistry      cldf.ContractType = "CapabilitiesRegistry"      // https://github.com/smartcontractkit/chainlink-evm/blob/f190212bab15e84fe49db88f495ad026e6c1d520/contracts/src/v0.8/workflow/dev/v2/CapabilitiesRegistry.sol#L450
	WorkflowRegistry          cldf.ContractType = "WorkflowRegistry"          // https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/workflow/WorkflowRegistry.sol
	KeystoneForwarder         cldf.ContractType = "KeystoneForwarder"         // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
	OCR3Capability            cldf.ContractType = "OCR3Capability"            // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
	FeedConsumer              cldf.ContractType = "FeedConsumer"              // no type and a version in contract https://github.com/smartcontractkit/chainlink/blob/89183a8a5d22b1aeca0ade3b76d16aa84067aa57/contracts/src/v0.8/keystone/KeystoneFeedsConsumer.sol#L1
	RBACTimelock              cldf.ContractType = "RBACTimelock"              // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/RBACTimelock.sol
	ProposerManyChainMultiSig cldf.ContractType = "ProposerManyChainMultiSig" // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/ManyChainMultiSig.sol
)

type RegisteredDonConfig struct {
	Name             string
	NodeIDs          []string // ids in the offchain client
	RegistryChainSel uint64
	Registry         *capabilities_registry_v2.CapabilitiesRegistry
}

type DonNodeSet struct {
	Name    string
	NodeIDs []string
}

// RegisteredDon is a representation of a don that exists in the capabilities registry all with the enriched node data
type RegisteredDon struct {
	Name  string
	Info  capabilities_registry_v2.CapabilitiesRegistryDONInfo
	Nodes []deployment.Node
}
