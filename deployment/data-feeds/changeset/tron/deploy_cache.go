package tron

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var DeployCacheChangeset = cldf.CreateChangeSet(deployCacheLogic, deployCachePrecondition)

// DeployCacheChangeset deploys the DataFeedsCache contract to the specified chains
// Returns a new addressbook with deployed DataFeedsCache contracts
func deployCacheLogic(env cldf.Environment, c types.DeployTronConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.TronChains()[chainSelector]
		cacheResponse, err := DeployCache(chain, c.DeployOptions, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", cacheResponse.Tv.String(), chain.Selector, cacheResponse.Address.String())

		err = ab.Save(chain.Selector, cacheResponse.Address.String(), cacheResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save DataFeedsCache: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       cacheResponse.Address.String(),
				Type:          cs.DataFeedsCache,
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(cacheResponse.Tv.Labels.List()...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}

func deployCachePrecondition(env cldf.Environment, c types.DeployTronConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.TronChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
