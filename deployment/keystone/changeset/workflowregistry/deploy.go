package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[uint64] = Deploy

func Deploy(env deployment.Environment, registrySelector uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	wrResp, err := deployWorkflowRegistry(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", wrResp.Tv.String(), chain.Selector, wrResp.Address.String())

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
