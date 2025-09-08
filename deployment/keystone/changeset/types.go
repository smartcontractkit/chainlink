package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

var (
	CapabilitiesRegistry      = contracts.CapabilitiesRegistry
	WorkflowRegistry          = contracts.WorkflowRegistry
	KeystoneForwarder         = contracts.KeystoneForwarder
	OCR3Capability            = contracts.OCR3Capability
	BalanceReader             = contracts.BalanceReader
	FeedConsumer              = contracts.FeedConsumer
	RBACTimelock              = contracts.RBACTimelock
	ProposerManyChainMultiSig = contracts.ProposerManyChainMultiSig
)
