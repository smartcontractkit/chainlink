package changeset

import "github.com/smartcontractkit/chainlink/deployment"

var (
	CapabilitiesRegistry      deployment.ContractType = "CapabilitiesRegistry"      // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L392
	WorkflowRegistry          deployment.ContractType = "WorkflowRegistry"          // https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/workflow/WorkflowRegistry.sol
	KeystoneForwarder         deployment.ContractType = "KeystoneForwarder"         // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
	OCR3Capability            deployment.ContractType = "OCR3Capability"            // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
	FeedConsumer              deployment.ContractType = "FeedConsumer"              // no type and a version in contract https://github.com/smartcontractkit/chainlink/blob/89183a8a5d22b1aeca0ade3b76d16aa84067aa57/contracts/src/v0.8/keystone/KeystoneFeedsConsumer.sol#L1
	RBACTimelock              deployment.ContractType = "RBACTimelock"              // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/RBACTimelock.sol
	ProposerManyChainMultiSig deployment.ContractType = "ProposerManyChainMultiSig" // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/ManyChainMultiSig.sol
)
