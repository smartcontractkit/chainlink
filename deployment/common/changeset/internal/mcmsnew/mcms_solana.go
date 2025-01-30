package mcmsnew

import (
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// MCMSWithTimelockSolanaDeploy holds a bundle of MCMS contract deploys.
type MCMSWithTimelockSolanaDeploy struct {
	McmProgram           *solana.PublicKey
	CancellerSeed        [32]byte
	BypasserSeed         [32]byte
	ProposerSeed         [32]byte
	Timelock             *solana.PublicKey
	TimelockInstanceSeed [32]byte
	// TODO: not sure if this is needed
	CallProxy *solana.PublicKey
}

func deployMCMSWithConfigSolana(
	contractType deployment.ContractType,
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	mcmConfig MCMSWithTimelockConfig,
) (any, error) {
	panic("implement me")
}

// deployMCMSWithTimelockContractsSolana deploys an MCMS program f
// and initializes 3 instances for each of the timelock roles: Bypasser, ProposerMcm, Canceller on an Solana chain.
// as well as the timelock program. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func deployMCMSWithTimelockContractsSolana(
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	config types.MCMSWithTimelockConfig,
) (*MCMSWithTimelockEVMDeploy, error) {
	panic("implement me")
}
