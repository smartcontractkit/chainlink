package changeset

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

// We expect the change set input to be unique per change set.
// TODO: Maybe there's a generics approach here?
// Note if the change set is a deployment and it fails we have 2 options:
// - Just throw away the addresses, fix issue and try again (potentially expensive on mainnet)
// - Roll forward with another change set completing the deployment
func Apply0001(env deployment.Environment, c ccipdeployment.DeployCCIPContractConfig) (deployment.ChangeSetOutput, error) {
	ab, err := ccipdeployment.DeployCCIPContracts(env, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "addresses", ab)
		return deployment.ChangeSetOutput{}, err
	}
	js, err := ccipdeployment.NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangeSetOutput{}, err
	}
	proposal, err := ccipdeployment.GenerateAcceptOwnershipProposal(env, env.AllChainSelectors(), ab)
	if err != nil {
		return deployment.ChangeSetOutput{}, err
	}
	return deployment.ChangeSetOutput{
		Proposals:   []deployment.Proposal{proposal},
		AddressBook: ab,
		// Mapping of which nodes get which jobs.
		JobSpecs: js,
	}, nil
}
