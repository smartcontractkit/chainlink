package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
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

	proposal, err := BuildMCMProposals(env, "accept ownership to timelock", c.ChainSelector, []ProposalData{
		{
			contract: c.ContractAddress.Hex(),
			tx:       tx,
		},
	}, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func acceptOwnershipPrecondition(env deployment.Environment, c types.AcceptOwnershipConfig) error {
	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}

	if _, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, commonTypes.RBACTimelock); err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}

	return nil
}
