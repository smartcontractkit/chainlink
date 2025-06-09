package changeset

import (
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// ProposeAggregatorChangeset is a changeset that proposes a new aggregator on existing AggregatorProxy contract
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var ProposeAggregatorChangeset = cldf.CreateChangeSet(proposeAggregatorLogic, proposeAggregatorPrecondition)

func proposeAggregatorLogic(env cldf.Environment, c types.ProposeConfirmAggregatorConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.EVMChains()[c.ChainSelector]

	aggregatorProxy, err := proxy.NewAggregatorProxy(c.ProxyAddress, chain.Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load AggregatorProxy: %w", err)
	}

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := aggregatorProxy.ProposeAggregator(txOpt, c.NewAggregatorAddress)

	if c.McmsConfig != nil {
		proposalConfig := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: aggregatorProxy.Address().Hex(),
					tx:       tx,
				},
			},
		}
		proposal, err := BuildMultiChainProposals(env, "proposal to propose a new aggregator", proposalConfig, c.McmsConfig.MinDelay)
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

func proposeAggregatorPrecondition(env cldf.Environment, c types.ProposeConfirmAggregatorConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil {
			return err
		}
	}

	return nil
}
