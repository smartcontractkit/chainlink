package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// AcceptOwnershipChangeset is a changeset that will create an MCM proposal to accept the ownership of a contract.
// Returns an MSM proposal to accept the ownership of a contract. Doesn't return a new addressbook.
// Once proposal is executed, new owned contract can be imported into the addressbook.
var AcceptOwnershipChangeset = deployment.CreateChangeSet(acceptOwnershipLogic, acceptOwnershipPrecondition)

func acceptOwnershipLogic(env deployment.Environment, c types.AcceptOwnershipConfig) (deployment.ChangesetOutput, error) {
	chain := env.Chains[c.ChainSelector]

	_, contract, err := commonChangesets.LoadOwnableContract(c.ContractAddress, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load the contract %w", err)
	}

	tx, err := contract.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create accept transfer ownership tx %w", err)
	}

	proposalConfig := MultiChainProposalConfig{
		c.ChainSelector: []ProposalData{
			{
				contract: contract.Address().Hex(),
				tx:       tx,
			},
		},
	}

	proposal, err := BuildMultiChainProposals(env, "accept ownership to timelock", proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func acceptOwnershipPrecondition(env deployment.Environment, c types.AcceptOwnershipConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}

	return ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector)
}
