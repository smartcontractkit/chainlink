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

var _ deployment.ChangeSet[types.AcceptOwnershipConfig] = AcceptOwnershipChangeset

// AcceptOwnershipChangeset is a changeset that will create a proposal to accept the ownership of a contract
func AcceptOwnershipChangeset(env deployment.Environment, c types.AcceptOwnershipConfig) (deployment.ChangesetOutput, error) {
	if c.McmsConfig == nil {
		return deployment.ChangesetOutput{}, errors.New("mcms config is required")
	}

	if _, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, commonTypes.RBACTimelock); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("mcms proposer not present on the chain %w", err)
	}

	chain := env.Chains[c.ChainSelector]

	_, contract, err := commonChangesets.LoadOwnableContract(c.ContractAddress, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load the contract %w", err)
	}

	tx, err := contract.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create accept transfer ownership tx %w", err)
	}

	proposal, err := BuildMCMProposal(env, "accept ownership to timelock", c.ChainSelector, c.ContractAddress.Hex(), tx, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}
