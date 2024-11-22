package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = DeployCapabilityRegistry

func DeployCapabilityRegistry(env deployment.Environment, registrySelector uint64) (deployment.ChangesetOutput, error) {
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
