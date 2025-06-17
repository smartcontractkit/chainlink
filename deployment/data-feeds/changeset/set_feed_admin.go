package changeset

import (
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedAdminChangeset is a changeset that sets/removes an admin on DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var SetFeedAdminChangeset = cldf.CreateChangeSet(setFeedAdminLogic, setFeedAdminPrecondition)

func setFeedAdminLogic(env cldf.Environment, c types.SetFeedAdminConfig) (cldf.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.BlockChains.EVMChains()[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := contract.SetFeedAdmin(txOpt, c.AdminAddress, c.IsAdmin)

	if c.McmsConfig != nil {
		proposalConfig := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: contract.Address().Hex(),
					tx:       tx,
				},
			},
		}

		proposal, err := BuildMultiChainProposals(env, "proposal to set feed admin on a cache", proposalConfig, c.McmsConfig.MinDelay)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return cldf.ChangesetOutput{}, nil
}

func setFeedAdminPrecondition(env cldf.Environment, c types.SetFeedAdminConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil {
			return err
		}
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
