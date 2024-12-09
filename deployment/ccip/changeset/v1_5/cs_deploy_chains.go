package v1_5

import (
	"fmt"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[DeployChainContractsConfig] = DeployChainContracts

type DeployChainContractsConfig struct {
	ChainSelectors []uint64
}

func DeployChainContracts(env deployment.Environment, c DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := deployChainContractsForChains(env, newAddresses, c.ChainSelectors)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}

func (c DeployChainContractsConfig) Validate() error {
	for _, cs := range c.ChainSelectors {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

func deployChainContractsForChains(env deployment.Environment, ab deployment.AddressBook, selectors []uint64) error {

}
