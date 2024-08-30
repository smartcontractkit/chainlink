package keystone

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var (
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	KeystoneForwarder    deployment.ContractType = "KeystoneForwarder"
	OCR3Capability       deployment.ContractType = "OCR3Capability"
)

type deployable interface {
	// deploy deploys the contract and returns the address and transaction
	// implementors must confirm the transaction and return any error
	deploy(deployRequest) (deployResponse, error)
}

type deployResponse struct {
	Address common.Address
	Tx      common.Hash // todo: chain agnostic
	Tv      deployment.TypeAndVersion
}

type deployRequest struct {
	Chain deployment.Chain
}

type DonNode struct {
	Don  string
	Node string // not unique across environments
}

type CapabilityHost struct {
	NodeID       string // gloablly unique
	Capabilities []capabilities_registry.CapabilitiesRegistryCapability
}

type Nop struct {
	capabilities_registry.CapabilitiesRegistryNodeOperator
	NodeIDs []string // nodes run by this operator
}
