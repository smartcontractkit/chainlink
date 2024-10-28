package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

// We expect the change set input to be unique per change set.
// TODO: Maybe there's a generics approach here?
// Note if the change set is a deployment and it fails we have 2 options:
// - Just throw away the addresses, fix issue and try again (potentially expensive on mainnet)
func InitialDeployChangeSet(ab deployment.AddressBook, env deployment.Environment, c ccipdeployment.DeployCCIPContractConfig) (deployment.ChangesetOutput, error) {
	err := ccipdeployment.DeployCCIPContracts(env, ab, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	js, err := ccipdeployment.NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		// Mapping of which nodes get which jobs.
		JobSpecs: js,
	}, nil
}
