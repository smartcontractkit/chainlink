package mcmsnew

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/mcmsnew"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfig] = DeployMCMSWithTimelockSolana

func DeployMCMSWithTimelockSolana(e deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := mcmsnew.DeployMCMSWithTimelockContractsBatch(e.Logger, e.SolChains, newAddresses, cfgByChain)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: newAddresses}, err
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
