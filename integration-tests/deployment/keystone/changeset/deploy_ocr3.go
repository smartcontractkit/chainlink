package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func DeployOCR3(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	// must have capabilities registry deployed
	regAddrs, err := ab.AddressesForChain(registryChainSel)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", registryChainSel, err)
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
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	err = kslib.DeployOCR3(lggr, c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil

}
