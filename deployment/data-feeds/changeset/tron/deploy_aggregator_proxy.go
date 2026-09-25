package tron

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployAggregatorProxyChangeset deploys an AggregatorProxy contract on the given chains. It uses the address of DataFeedsCache contract
// from addressbook to set it in the AggregatorProxy constructor. Returns a new addressbook with deploy AggregatorProxy contract addresses.
var DeployAggregatorProxyChangeset = cldf.CreateChangeSet(deployAggregatorProxyLogic, deployAggregatorProxyPrecondition)

func deployAggregatorProxyLogic(env cldf.Environment, c types.DeployAggregatorProxyTronConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore()

	for index, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.TronChains()[chainSelector]

		cacheAddressStr := changeset.GetDataFeedsCacheAddress(env.ExistingAddresses, env.DataStore.Addresses(), chainSelector, nil)
		if cacheAddressStr == "" {
			return cldf.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d", chainSelector)
		}

		cacheAddress, err := address.Base58ToAddress(cacheAddressStr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to parse DataFeedsCache contract address %s", cacheAddressStr)
		}

		proxyResponse, err := DeployAggregatorProxy(chain, cacheAddress, c.AccessController[index], c.DeployOptions, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy AggregatorProxy: %w", err)
		}

		lggr.Infof("Deployed %s chain selector %d addr %s", proxyResponse.Tv.String(), chain.Selector, proxyResponse.Address.String())

		err = ab.Save(chain.Selector, proxyResponse.Address.String(), proxyResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save AggregatorProxy: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       proxyResponse.Address.String(),
				Type:          changeset.DataFeedsCache,
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(proxyResponse.Tv.Labels.List()...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}
	return cldf.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}

func deployAggregatorProxyPrecondition(env cldf.Environment, c types.DeployAggregatorProxyTronConfig) error {
	if len(c.AccessController) != len(c.ChainsToDeploy) {
		return errors.New("AccessController addresses must be provided for each chain to deploy")
	}

	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.TronChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
		_, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get addessbook for chain %d: %w", chainSelector, err)
		}
	}

	return nil
}
