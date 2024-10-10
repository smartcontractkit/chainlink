package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func DeployCapabilityRegistry(lggr logger.Logger, env deployment.Environment, ab deployment.AddressBook, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	err := kslib.DeployCapabilitiesRegistry(lggr, c, ab)

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
