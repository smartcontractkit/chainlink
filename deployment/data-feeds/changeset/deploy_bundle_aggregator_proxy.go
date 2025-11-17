package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployBundleAggregatorProxyChangeset deploys a BundleAggregatorProxy contract on the given chains. It uses the address of DataFeedsCache contract
// from DataStore to set it in the BundleAggregatorProxy constructor. It uses the provided owner address to set it in the BundleAggregatorProxy constructor.
// Returns a new DataStore with deployed BundleAggregatorProxy contract addresses.
var DeployBundleAggregatorProxyChangeset = cldf.CreateChangeSet(deployBundleAggregatorProxyLogic, deployBundleAggregatorProxyPrecondition)

func deployBundleAggregatorProxyLogic(env cldf.Environment, c types.DeployBundleAggregatorProxyConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ds := datastore.NewMemoryDataStore()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.EVMChains()[chainSelector]

		dataFeedsCacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, env.DataStore.Addresses(), chainSelector, &c.CacheLabel)
		if dataFeedsCacheAddress == "" {
			return cldf.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d", chainSelector)
		}

		bundleProxyResponse, err := DeployBundleAggregatorProxy(chain, common.HexToAddress(dataFeedsCacheAddress), c.Owners[chainSelector], c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BundleAggregatorProxy: %w", err)
		}

		lggr.Infof("Deployed %s chain selector %d addr %s", bundleProxyResponse.Tv.String(), chain.Selector, bundleProxyResponse.Address.String())

		if err = ds.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       bundleProxyResponse.Address.String(),
				Type:          datastore.ContractType(bundleProxyResponse.Tv.Type),
				Version:       &bundleProxyResponse.Tv.Version,
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(c.Labels...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}
	return cldf.ChangesetOutput{DataStore: ds}, nil
}

func deployBundleAggregatorProxyPrecondition(env cldf.Environment, c types.DeployBundleAggregatorProxyConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
		if !common.IsHexAddress(c.Owners[chainSelector].String()) {
			return fmt.Errorf("owner %s is not a valid address for chain %d", c.Owners[chainSelector].String(), chainSelector)
		}
	}

	return nil
}
