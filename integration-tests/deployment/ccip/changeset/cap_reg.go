package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

// Separated changset because cap reg is an env var for CL nodes.
func CapRegChangeSet(env deployment.Environment, homeChainSel uint64) (deployment.ChangesetOutput, error) {
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
