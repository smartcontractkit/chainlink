package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	proxy "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/aggregator_proxy"
)

// ConfirmAggregatorChangeset is a changeset that confirms a proposed aggregator on deployed AggregatorProxy contract
var ConfirmAggregatorChangeset = deployment.CreateChangeSet(confirmAggregatorLogic, confirmAggregatorPrecondition)

func confirmAggregatorLogic(env deployment.Environment, c types.ProposeConfirmAggregatorConfig) (deployment.ChangesetOutput, error) {
	chain := env.Chains[c.ChainSelector]

	aggregatorProxy, err := proxy.NewAggregatorProxy(c.ProxyAddress, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load AggregatorProxy: %w", err)
	}

	tx, err := aggregatorProxy.ConfirmAggregator(chain.DeployerKey, c.NewAggregatorAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute ConfirmAggregator: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{}, nil
}

func confirmAggregatorPrecondition(env deployment.Environment, c types.ProposeConfirmAggregatorConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	return nil
}
