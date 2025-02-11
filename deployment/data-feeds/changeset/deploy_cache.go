package changeset

import (
	"errors"
	"fmt"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var _ deployment.ChangeSet[types.DeployConfig] = DeployCacheChangeset

func DeployCacheChangeset(env deployment.Environment, c types.DeployConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	for _, chainSelector := range c.ChainsToDeploy {
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
		}
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
