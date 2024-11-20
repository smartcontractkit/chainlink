package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[DeployCCIPContractConfig] = InitialDeploy

func InitialDeploy(env deployment.Environment, c DeployCCIPContractConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := DeployCCIPContracts(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	js, err := NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
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
