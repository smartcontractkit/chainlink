package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployAggregatorProxyChangeset deploys an AggregatorProxy contract on the given chains. It uses the address of DataFeedsCache contract
// from addressbook to set it in the AggregatorProxy constructor. Returns a new addressbook with deploy AggregatorProxy contract addresses.
var DeployAggregatorProxyChangeset = cldf.CreateChangeSet(deployAggregatorProxyLogic, deployAggregatorProxyPrecondition)

func deployAggregatorProxyLogic(env cldf.Environment, c types.DeployAggregatorProxyConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()

	for index, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.EVMChains()[chainSelector]

		dataFeedsCacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, chainSelector, nil)
		if dataFeedsCacheAddress == "" {
			return cldf.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d", chainSelector)
		}

		proxyResponse, err := DeployAggregatorProxy(chain, common.HexToAddress(dataFeedsCacheAddress), c.AccessController[index], c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy AggregatorProxy: %w", err)
		}

		lggr.Infof("Deployed %s chain selector %d addr %s", proxyResponse.Tv.String(), chain.Selector, proxyResponse.Address.String())

		err = ab.Save(chain.Selector, proxyResponse.Address.String(), proxyResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save AggregatorProxy: %w", err)
		}
	}
	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func deployAggregatorProxyPrecondition(env cldf.Environment, c types.DeployAggregatorProxyConfig) error {
	if len(c.AccessController) != len(c.ChainsToDeploy) {
		return errors.New("AccessController addresses must be provided for each chain to deploy")
	}

	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.EVMChains()[chainSelector]
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
