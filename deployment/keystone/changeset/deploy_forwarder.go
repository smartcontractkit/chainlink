package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet = DeployForwarder

func DeployForwarder(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	registryChainSel, ok := config.(uint64)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	lggr := env.Logger
	// expect OCR3 to be deployed & capabilities registry
	regAddrs, err := env.ExistingAddresses.AddressesForChain(registryChainSel)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", registryChainSel, err)
	}
	if len(regAddrs) != 2 {
		return deployment.ChangesetOutput{}, fmt.Errorf("expected 2 addresses for chain %d, got %d", registryChainSel, len(regAddrs))
	}
	ab := deployment.NewMemoryAddressBook()
	for _, chain := range env.Chains {
		lggr.Infow("deploying forwarder", "chainSelector", chain.Selector)
		err := kslib.DeployForwarder(lggr, chain, ab)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", chain.Selector, err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
