package changeset

import (
	"fmt"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[DeployChainContractsConfig] = DeployChainContracts

func DeployChainContracts(env deployment.Environment, c DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := DeployChainContractsForChains(env, newAddresses, c.HomeChainSelector, c.ChainSelectors)
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

type DeployChainContractsConfig struct {
	ChainSelectors    []uint64
	HomeChainSelector uint64
}

func (c DeployChainContractsConfig) Validate() error {
	for _, cs := range c.ChainSelectors {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	return nil
}
