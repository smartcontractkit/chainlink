package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployCacheChangeset deploys the DataFeedsCache contract to the specified chains
// Returns a new addressbook with deployed DataFeedsCache contracts
var DeployCacheChangeset = cldf.CreateChangeSet(deployCacheLogic, deployCachePrecondition)

func deployCacheLogic(env deployment.Environment, c types.DeployConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.Chains[chainSelector]
		cacheResponse, err := DeployCache(chain, c.Labels)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", cacheResponse.Tv.String(), chain.Selector, cacheResponse.Address.String())

		err = ab.Save(chain.Selector, cacheResponse.Address.String(), cacheResponse.Tv)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save DataFeedsCache: %w", err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func deployCachePrecondition(env deployment.Environment, c types.DeployConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.Chains[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
