package changeset

import (
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var _ deployment.ChangeSet[types.SetFeedDecimalConfig] = SetFeedConfigChangeset

func SetFeedConfigChangeset(env deployment.Environment, c types.SetFeedDecimalConfig) (deployment.ChangesetOutput, error) {
	err := ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate cache for chain %w", err)
	}

	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = deployment.SimTransactOpts()
	}

	tx, err := contract.SetDecimalFeedConfigs(txOpt, c.DataIDs, c.Descriptions, c.WorkflowMetadata)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed config %w", err)
	}

	if c.McmsConfig != nil {
		proposal, err := BuildMCMProposal(env, "proposal to set feed config on a cache", c.ChainSelector, contract.Address().Hex(), tx, c.McmsConfig.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{}, nil
}
