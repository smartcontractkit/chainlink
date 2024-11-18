package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfig] = DeployMCMSWithTimelock

func DeployMCMSWithTimelock(e deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := internal.DeployMCMSWithTimelockContractsBatch(
		e.Logger, e.Chains, newAddresses, cfgByChain,
	)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: newAddresses}, err
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
