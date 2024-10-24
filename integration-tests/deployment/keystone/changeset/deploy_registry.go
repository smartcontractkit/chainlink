package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func DeployCapabilityRegistry(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	registrySelector, ok := config.(uint64)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	err := kslib.DeployCapabilitiesRegistry(env.Logger, chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
