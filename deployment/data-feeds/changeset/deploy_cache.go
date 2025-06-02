package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployCacheChangeset deploys the DataFeedsCache contract to the specified chains
// Returns a new addressbook with deployed DataFeedsCache contracts
var DeployCacheChangeset = cldf.CreateChangeSet(deployCacheLogic, deployCachePrecondition)

func deployCacheLogic(env cldf.Environment, c types.DeployConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.EVMChains()[chainSelector]
		cacheResponse, err := DeployCache(chain, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", cacheResponse.Tv.String(), chain.Selector, cacheResponse.Address.String())

		err = ab.Save(chain.Selector, cacheResponse.Address.String(), cacheResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save DataFeedsCache: %w", err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab}, nil
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
