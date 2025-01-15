package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type DeployFeedsConsumerRequest struct {
	ChainSelector uint64
}

var _ deployment.ChangeSet[*DeployFeedsConsumerRequest] = DeployFeedsConsumer

// DeployFeedsConsumer deploys the FeedsConsumer contract to the chain with the given chainSelector.
func DeployFeedsConsumer(env deployment.Environment, req *DeployFeedsConsumerRequest) (deployment.ChangesetOutput, error) {
	chainSelector := req.ChainSelector
	lggr := env.Logger
	chain, ok := env.Chains[chainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	deployResp, err := kslib.DeployFeedsConsumer(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy FeedsConsumer: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", deployResp.Tv.String(), chain.Selector, deployResp.Address.String())
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
