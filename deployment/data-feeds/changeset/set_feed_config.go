package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedConfigChangeset is a changeset that sets a feed configuration on DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var SetFeedConfigChangeset = cldf.CreateChangeSet(setFeedConfigLogic, setFeedConfigPrecondition)

func setFeedConfigLogic(env deployment.Environment, c types.SetFeedDecimalConfig) (deployment.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = deployment.SimTransactOpts()
	}

	dataIDs, _ := FeedIDsToBytes16(c.DataIDs)
	tx, err := contract.SetDecimalFeedConfigs(txOpt, dataIDs, c.Descriptions, c.WorkflowMetadata)

	if c.McmsConfig != nil {
		proposals := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: contract.Address().Hex(),
					tx:       tx,
				},
			},
		}
		proposal, err := BuildMultiChainProposals(env, "proposal to set feed config on a cache", proposals, c.McmsConfig.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		if tx != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
		}

		return deployment.ChangesetOutput{}, fmt.Errorf("failed to submit transaction: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

func setFeedConfigPrecondition(env deployment.Environment, c types.SetFeedDecimalConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if (len(c.DataIDs) == 0) || (len(c.Descriptions) == 0) || (len(c.WorkflowMetadata) == 0) {
		return errors.New("dataIDs, descriptions and workflowMetadata must not be empty")
	}
	if len(c.DataIDs) != len(c.Descriptions) {
		return errors.New("dataIDs and descriptions must have the same length")
	}
	_, err := FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil {
			return err
		}
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
