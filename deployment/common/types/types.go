package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	BypasserManyChainMultisig  deployment.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig deployment.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  deployment.ContractType = "ProposerManyChainMultiSig"
	RBACTimelock               deployment.ContractType = "RBACTimelock"
)

type MCMSWithTimelockConfig struct {
	Canceller         config.Config
	Bypasser          config.Config
	Proposer          config.Config
	TimelockExecutors []common.Address
	TimelockMinDelay  *big.Int
}
