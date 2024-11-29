package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

// DeployFeedsConsumer deploys the FeedsConsumer contract to the chain with the given chainSelector.
func DeployFeedsConsumer(env deployment.Environment, chainSelector uint64) (deployment.ChangesetOutput, error) {
	chain, ok := env.Chains[chainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	err := kslib.DeployFeedsConsumer(env.Logger, chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy FeedsConsumer: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
