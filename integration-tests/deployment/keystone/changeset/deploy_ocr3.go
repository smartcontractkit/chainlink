package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

var _ deployment.ChangeSet = DeployOCR3

func DeployOCR3(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	c, ok := config.(DeployRegistryConfig)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	lggr := env.Logger
	// must have capabilities registry deployed
	regAddrs, err := c.ExistingAddressBook.AddressesForChain(c.RegistryChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", c.RegistryChainSelector, err)
	}
	found := false
	for _, addr := range regAddrs {
		if addr.Type == kslib.CapabilityRegistryTypeVersion.Type {
			found = true
			break
		}
	}
	if !found {
		return deployment.ChangesetOutput{}, fmt.Errorf("no capabilities registry found for changeset %s", "0001_deploy_registry")
	}

	// ocr3 only deployed on registry chain
	registryChain, ok := env.Chains[c.RegistryChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	err = kslib.DeployOCR3(lggr, registryChain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
