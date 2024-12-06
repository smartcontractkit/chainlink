package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = DeployForwarder

// DeployForwarder deploys the KeystoneForwarder contract to all chains in the environment
// callers must merge the output addressbook with the existing one
func DeployForwarder(env deployment.Environment, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	for _, chain := range env.Chains {
		lggr.Infow("deploying forwarder", "chainSelector", chain.Selector)
		forwarderResp, err := kslib.DeployForwarder(chain, ab)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", chain.Selector, err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResp.Tv.String(), chain.Selector, forwarderResp.Address.String())
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
