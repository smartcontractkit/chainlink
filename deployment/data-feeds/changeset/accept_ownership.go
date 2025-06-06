package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// AcceptOwnershipChangeset is a changeset that will create an MCM proposal to accept the ownership of contracts.
// Returns an MSM proposal to accept the ownership of contracts. Doesn't return a new addressbook.
// Once proposal is executed, new owned contracts can be imported into the addressbook.
var AcceptOwnershipChangeset = cldf.CreateChangeSet(acceptOwnershipLogic, acceptOwnershipPrecondition)

func acceptOwnershipLogic(env cldf.Environment, c types.AcceptOwnershipConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.EVMChains()[c.ChainSelector]

	var mcmsProposals []ProposalData
	for _, contractAddress := range c.ContractAddresses {
		_, contract, err := commonChangesets.LoadOwnableContract(contractAddress, chain.Client)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load the contract %w", err)
		}

		tx, err := contract.AcceptOwnership(cldf.SimTransactOpts())
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create accept transfer ownership tx %w", err)
		}
		mcmsProposals = append(mcmsProposals, ProposalData{
			contract: contract.Address().Hex(),
			tx:       tx,
		})
	}

	proposalConfig := MultiChainProposalConfig{c.ChainSelector: mcmsProposals}
	proposal, err := BuildMultiChainProposals(env, "accept ownership to timelock", proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func acceptOwnershipPrecondition(env cldf.Environment, c types.AcceptOwnershipConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}

	return ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector)
}
