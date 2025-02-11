package changeset

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var _ deployment.ChangeSet[types.DeployAggregatorProxyConfig] = DeployAggregatorProxyChangeset

func DeployAggregatorProxyChangeset(env deployment.Environment, c types.DeployAggregatorProxyConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	if len(c.AccessController) != len(c.ChainsToDeploy) {
		return deployment.ChangesetOutput{}, errors.New("AccessController addresses must be provided for each chain to deploy")
	}

	for index, chainSelector := range c.ChainsToDeploy {
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
		}

		addressMap, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get addessbook for chain %d: %w", chainSelector, err)
		}

		var dataFeedsCacheAddress string
		cacheTV := deployment.NewTypeAndVersion(DataFeedsCache, deployment.Version1_0_0)
		cacheTV.Labels.Add("data-feeds")
		for addr, tv := range addressMap {
			if tv.String() == cacheTV.String() {
				dataFeedsCacheAddress = addr
			}
		}

		if dataFeedsCacheAddress == "" {
			return deployment.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d: %w", chainSelector, err)
		}

		proxyResponse, err := DeployAggregatorProxy(chain, common.HexToAddress(dataFeedsCacheAddress), c.AccessController[index], c.Labels)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy AggregatorProxy: %w", err)
		}

		lggr.Infof("Deployed %s chain selector %d addr %s", proxyResponse.Tv.String(), chain.Selector, proxyResponse.Address.String())

		err = ab.Save(chain.Selector, proxyResponse.Address.String(), proxyResponse.Tv)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save AggregatorProxy: %w", err)
		}
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
