package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetBundleFeedConfigChangeset is a changeset that sets a feed configuration on DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var SetBundleFeedConfigChangeset = cldf.CreateChangeSet(setBundleFeedConfigLogic, setBundleFeedConfigPrecondition)

func setBundleFeedConfigLogic(env cldf.Environment, c types.SetFeedBundleConfig) (cldf.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.BlockChains.EVMChains()[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	dataIDs, _ := FeedIDsToBytes16(c.DataIDs)
	tx, err := contract.SetBundleFeedConfigs(txOpt, dataIDs, c.Descriptions, c.DecimalsMatrix, c.WorkflowMetadata)

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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
		if tx != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
		}

		return cldf.ChangesetOutput{}, fmt.Errorf("failed to submit transaction: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func setBundleFeedConfigPrecondition(env cldf.Environment, c types.SetFeedBundleConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if (len(c.DataIDs) == 0) || (len(c.Descriptions) == 0) || (len(c.WorkflowMetadata) == 0) || (len(c.DecimalsMatrix) == 0) {
		return errors.New("dataIDs, descriptions, decimalsMatrix and workflowMetadata must not be empty")
	}
	if len(c.DataIDs) != len(c.Descriptions) || len(c.DataIDs) != len(c.DecimalsMatrix) {
		return errors.New("dataIDs, decimalsMatrix and descriptions must have the same length")
	}
	_, err := FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil { //nolint:staticcheck // TODO: replace with DataStore when ready
			return err
		}
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
