package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func DeployForwarder(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	// expect OCR3 to be deployed & capabilities registry
	regAddrs, err := ab.AddressesForChain(registryChainSel)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", registryChainSel, err)
	}
	if len(regAddrs) != 2 {
		return deployment.ChangesetOutput{}, fmt.Errorf("expected 2 addresses for chain %d, got %d", registryChainSel, len(regAddrs))
	}
	for _, c := range env.Chains {
		lggr.Infow("deploying forwarder", "chainSelector", c.Selector)
		err := kslib.DeployForwarder(lggr, c, ab)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", c.Selector, err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
