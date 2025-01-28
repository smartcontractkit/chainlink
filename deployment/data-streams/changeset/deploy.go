package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams"
)

func DeployChannelConfigStoreChangeSet(env deployment.Environment, cc data_streams.DeployContractConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := data_streams.DeployChannelConfigStore(env, ab, cc)
	if err != nil {
		env.Logger.Errorw("Failed to deploy ChannelConfigStore", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}
