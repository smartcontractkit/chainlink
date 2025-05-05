package aptos

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployPlatformChangeset deploys the ChainlinkPlatform package to Aptos chain
// Returns a new addressbook with deployed forwarder/storage contracts
var DeployPlatformChangeset = deployment.CreateChangeSet(deployPlatformLogic, deployPlatformPrecondition)

func deployPlatformLogic(env deployment.Environment, c types.DeployConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.AptosChains[chainSelector]
		cacheResponse, err := DeployPlatform(chain, chain.DeployerSigner.AccountAddress(), c.Labels)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy ChainlinkPlatform: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", cacheResponse.Tv.String(), chain.Selector, cacheResponse.Address.String())

		err = ab.Save(chain.Selector, cacheResponse.Address.String(), cacheResponse.Tv)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save ChainlinkPlatform: %w", err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func deployPlatformPrecondition(env deployment.Environment, c types.DeployConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.AptosChains[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
