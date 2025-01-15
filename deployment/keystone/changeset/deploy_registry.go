package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[uint64] = DeployCapabilityRegistry

func DeployCapabilityRegistry(env deployment.Environment, registrySelector uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	capabilitiesRegistryResp, err := kslib.DeployCapabilitiesRegistry(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", capabilitiesRegistryResp.Tv.String(), chain.Selector, capabilitiesRegistryResp.Address.String())

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
