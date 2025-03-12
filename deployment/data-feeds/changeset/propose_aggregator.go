package changeset

import (
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	proxy "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/aggregator_proxy"
)

// ProposeAggregatorChangeset is a changeset that proposes a new aggregator on existing AggregatorProxy contract
// This changeset may return a timelock proposal if the MCMS config is provided, otherwise it will execute the transaction with the deployer key.
var ProposeAggregatorChangeset = deployment.CreateChangeSet(proposeAggregatorLogic, proposeAggregatorPrecondition)

func proposeAggregatorLogic(env deployment.Environment, c types.ProposeConfirmAggregatorConfig) (deployment.ChangesetOutput, error) {
	chain := env.Chains[c.ChainSelector]

	aggregatorProxy, err := proxy.NewAggregatorProxy(c.ProxyAddress, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load AggregatorProxy: %w", err)
	}

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = deployment.SimTransactOpts()
	}

	tx, err := aggregatorProxy.ProposeAggregator(txOpt, c.NewAggregatorAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute ProposeAggregator: %w", err)
	}

	if c.McmsConfig != nil {
		proposal, err := BuildMCMProposals(env, "proposal to propose a new aggregator", c.ChainSelector, []ProposalData{
			{
				contract: aggregatorProxy.Address().Hex(),
				tx:       tx,
			},
		}, c.McmsConfig.MinDelay)
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

func proposeAggregatorPrecondition(env deployment.Environment, c types.ProposeConfirmAggregatorConfig) error {
	_, ok := env.Chains[c.ChainSelector]
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
