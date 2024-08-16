package migrations

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

// We expect the migration input to be unique per migration.
// TODO: Maybe there's a generics approach here?
// Note if the migration is a deployment and it fails we have 2 options:
// - Just throw away the addresses, fix issue and try again (potentially expensive on mainnet)
// - Roll forward with another migration completing the deployment
func Apply0001(env deployment.Environment, c ccipdeployment.DeployCCIPContractConfig) (deployment.MigrationOutput, error) {
	ab, err := ccipdeployment.DeployCCIPContracts(env, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "addresses", ab)
		return deployment.MigrationOutput{}, err
	}
	js, err := ccipdeployment.NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.MigrationOutput{}, err
	}
	proposal, err := ccipdeployment.GenerateAcceptOwnershipProposal(env, env.AllChainSelectors(), ab)
	if err != nil {
		return deployment.MigrationOutput{}, err
	}
	return deployment.MigrationOutput{
		Proposals:   []deployment.Proposal{proposal},
		AddressBook: ab,
		// Mapping of which nodes get which jobs.
		JobSpecs: js,
	}, nil
}
