package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

func DeployCapabilityRegistry(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, *capabilities_registry.CapabilitiesRegistry, error) {
	registrySelector, ok := config.(uint64)
	if !ok {
		return deployment.ChangesetOutput{}, nil, deployment.ErrInvalidConfig
	}
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, nil, fmt.Errorf("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	capReg, err := kslib.DeployCapabilitiesRegistry(env.Logger, chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, capReg, nil
}
