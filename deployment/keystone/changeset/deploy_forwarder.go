package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = DeployForwarder

// DeployForwarder deploys the KeystoneForwarder contract to all chains in the environment
// callers must merge the output addressbook with the existing one
// TODO: add selectors to deploy only to specific chains
func DeployForwarder(env deployment.Environment, _ uint64) (deployment.ChangesetOutput, error) {
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

var _ deployment.ChangeSet[ConfigureForwardContractsRequest] = ConfigureForwardContracts

type ConfigureForwardContractsRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	UseMCMS bool
}

func (r ConfigureForwardContractsRequest) Validate() error {
	if len(r.WFNodeIDs) == 0 {
		return fmt.Errorf("WFNodeIDs must not be empty")
	}
	return nil
}

func ConfigureForwardContracts(env deployment.Environment, req ConfigureForwardContractsRequest) (deployment.ChangesetOutput, error) {
	wfDon, err := kslib.NewRegisteredDon(env, kslib.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}
	r, err := kslib.ConfigureForwardContracts(&env, kslib.ConfigureForwarderContractsRequest{
		Dons:    []kslib.RegisteredDon{*wfDon},
		UseMCMS: req.UseMCMS,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure forward contracts: %w", err)
	}
	return deployment.ChangesetOutput{
		Proposals: r.Proposals,
	}, nil
}
