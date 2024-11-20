package changeset

import (
	"fmt"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[DeployCCIPContractConfig] = InitialAddChain

// InitialAddChain enables new chains as destination for CCIP
// It performs the following steps:
// - AddChainConfig + AddDON (candidate->primary promotion i.e. init) on the home chain
// - SetOCR3Config on the remote chain
// ConfigureChain assumes that the home chain is already enabled and all CCIP contracts are already deployed.
func InitialAddChain(env deployment.Environment, c DeployCCIPContractConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid InitialAddChainConfig: %w", err)
	}
	err := ConfigureChain(env, c)
	if err != nil {
		env.Logger.Errorw("Failed to configure chain", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nil,
	}, nil
}
