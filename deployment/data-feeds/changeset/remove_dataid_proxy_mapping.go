package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// RemoveFeedProxyMappingChangeset is a changeset that only removes a feed-aggregator proxy mapping from DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var RemoveFeedProxyMappingChangeset = cldf.CreateChangeSet(removeFeedProxyMappingLogic, removeFeedFeedProxyMappingPrecondition)

func removeFeedProxyMappingLogic(env cldf.Environment, c types.RemoveFeedProxyConfig) (cldf.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.BlockChains.EVMChains()[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := contract.RemoveDataIdMappingsForProxies(txOpt, c.ProxyAddresses)

	if c.McmsConfig != nil {
		proposalConfig := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: contract.Address().Hex(),
					tx:       tx,
				},
			},
		}

		proposal, err := BuildMultiChainProposals(env, "proposal to remove a feed proxy mapping from cache", proposalConfig, c.McmsConfig.MinDelay)
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

func removeFeedFeedProxyMappingPrecondition(env cldf.Environment, c types.RemoveFeedProxyConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if len(c.ProxyAddresses) == 0 {
		return errors.New("proxy addresses must not be empty")
	}
	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil {
			return err
		}
	}

	return ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
}
