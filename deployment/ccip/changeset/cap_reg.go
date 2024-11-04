package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var _ deployment.ChangeSet = DeployCapReg

// DeployCapReg is a separate changeset because cap reg is an env var for CL nodes.
func DeployCapReg(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	homeChainSel, ok := config.(uint64)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	// Note we also deploy the cap reg.
	ab := deployment.NewMemoryAddressBook()
	_, err := ccipdeployment.DeployCapReg(env.Logger, ab, env.Chains[homeChainSel])
	if err != nil {
		env.Logger.Errorw("Failed to deploy cap reg", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}
