package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"

	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var _ deployment.ChangeSet[ccipdeployment.DeployCCIPContractConfig] = InitialDeploy

func InitialDeploy(env deployment.Environment, c ccipdeployment.DeployCCIPContractConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := ccipdeployment.DeployCCIPContracts(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	js, err := ccipdeployment.NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: newAddresses}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		// Mapping of which nodes get which jobs.
		JobSpecs: js,
	}, nil
}
