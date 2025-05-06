package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// RemoveFeedChangeset is a changeset that removes a feed configuration and aggregator proxy mapping from DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transactions with the deployer key.
var RemoveFeedChangeset = cldf.CreateChangeSet(removeFeedLogic, removeFeedPrecondition)

func removeFeedLogic(env deployment.Environment, c types.RemoveFeedConfig) (deployment.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = deployment.SimTransactOpts()
	}
	dataIDs, _ := FeedIDsToBytes16(c.DataIDs)

	// remove the feed config
	removeConfigTx, err := contract.RemoveFeedConfigs(txOpt, dataIDs)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to remove feed config %w", err)
	}

	if c.McmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, removeConfigTx, err); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", removeConfigTx.Hash().String(), err)
		}
	}

	// remove from proxy mapping
	removeProxyMappingTx, err := contract.RemoveDataIdMappingsForProxies(txOpt, c.ProxyAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to remove proxy mapping %w", err)
	}

	if c.McmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, removeProxyMappingTx, err); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", removeProxyMappingTx.Hash().String(), err)
		}
		return deployment.ChangesetOutput{}, nil
	}

	proposalConfig := MultiChainProposalConfig{
		c.ChainSelector: []ProposalData{
			{
				contract: contract.Address().Hex(),
				tx:       removeConfigTx,
			},
			{
				contract: contract.Address().Hex(),
				tx:       removeProxyMappingTx,
			},
		},
	}

	proposal, err := BuildMultiChainProposals(env, "proposal to remove a feed from cache", proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func removeFeedPrecondition(env deployment.Environment, c types.RemoveFeedConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if (len(c.DataIDs) == 0) || (len(c.ProxyAddresses) == 0) {
		return errors.New("dataIDs and proxy addresses must not be empty")
	}
	if len(c.DataIDs) != len(c.ProxyAddresses) {
		return errors.New("dataIDs and proxy addresses must have the same length")
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
