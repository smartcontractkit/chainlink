package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// UpdateDataIDProxyChangeset is a changeset that updates the proxy-dataId mapping on DataFeedsCache contract.
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var UpdateDataIDProxyChangeset = cldf.CreateChangeSet(updateDataIDProxyLogic, updateDataIDProxyPrecondition)

func updateDataIDProxyLogic(env cldf.Environment, c types.UpdateDataIDProxyConfig) (cldf.ChangesetOutput, error) {
	state, _ := LoadOnchainState(env)
	chain := env.BlockChains.EVMChains()[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	dataIDs, _ := FeedIDsToBytes16(c.DataIDs)
	tx, err := contract.UpdateDataIdMappingsForProxies(txOpt, c.ProxyAddresses, dataIDs)

	if c.McmsConfig != nil {
		proposals := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: contract.Address().Hex(),
					tx:       tx,
				},
			},
		}

		proposal, err := BuildMultiChainProposals(env, "proposal to update proxy-dataId mapping on a cache", proposals, c.McmsConfig.MinDelay)
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

func updateDataIDProxyPrecondition(env cldf.Environment, c types.UpdateDataIDProxyConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if len(c.ProxyAddresses) == 0 || len(c.DataIDs) == 0 {
		return errors.New("empty proxies or dataIds")
	}
	if len(c.DataIDs) != len(c.ProxyAddresses) {
		return errors.New("dataIds and proxies length mismatch")
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
