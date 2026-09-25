package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployCacheChangeset deploys the DataFeedsCache contract to the specified chains
// Returns a new DataStore with deployed DataFeedsCache contracts
var DeployCacheChangeset = cldf.CreateChangeSet(deployCacheLogic, deployCachePrecondition)

func deployCacheLogic(env cldf.Environment, c types.DeployConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ds := datastore.NewMemoryDataStore()
	ab := cldf.NewMemoryAddressBook()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.EVMChains()[chainSelector]
		cacheResponse, err := DeployCache(chain, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", cacheResponse.Tv.String(), chain.Selector, cacheResponse.Address.String())

		if err = ds.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       cacheResponse.Address.String(),
				Type:          datastore.ContractType(cacheResponse.Tv.Type),
				Version:       &cacheResponse.Tv.Version,
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(c.Labels...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
		err = ab.Save(chain.Selector, cacheResponse.Address.String(), cacheResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save DataFeedsCache: %w", err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds, AddressBook: ab}, nil
}

func deployCachePrecondition(env cldf.Environment, c types.DeployConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
