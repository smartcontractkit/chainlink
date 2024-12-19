package workflowregistry

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = Deploy

func Deploy(env deployment.Environment, registrySelector uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	wrResp, err := kslib.DeployWorkflowRegistry(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", wrResp.Tv.String(), chain.Selector, wrResp.Address.String())

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
