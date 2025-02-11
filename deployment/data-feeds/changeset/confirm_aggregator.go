package changeset

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	proxy "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/aggregator_proxy"
)

var _ deployment.ChangeSet[types.ProposeConfirmAggregatorConfig] = ProposeAggregatorChangeset

func ConfirmAggregatorChangeset(env deployment.Environment, c types.ProposeConfirmAggregatorConfig) (deployment.ChangesetOutput, error) {
	chain, ok := env.Chains[c.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	aggregatorProxy, err := proxy.NewAggregatorProxy(c.Proxy, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load AggregatorProxy: %w", err)
	}

	tx, err := aggregatorProxy.ConfirmAggregator(chain.DeployerKey, c.NewAggregator)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute ConfirmAggregator: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{}, nil
}
