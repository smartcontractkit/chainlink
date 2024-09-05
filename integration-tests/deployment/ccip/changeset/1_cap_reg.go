package changeset

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

// Separate migration because cap reg is an env var for CL nodes.
func Apply0001(env deployment.Environment, homeChainSel uint64) (deployment.ChangesetOutput, error) {
	// Note we also deploy the cap reg.
	ab, _, err := ccipdeployment.DeployCapReg(env.Logger, env.Chains, homeChainSel)
	if err != nil {
		env.Logger.Errorw("Failed to deploy cap reg", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []deployment.Proposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}
