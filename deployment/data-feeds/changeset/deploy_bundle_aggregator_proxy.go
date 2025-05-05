package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployBundleAggregatorProxyChangeset deploys a BundleAggregatorProxy contract on the given chains. It uses the address of DataFeedsCache contract
// from addressbook to set it in the BundleAggregatorProxy constructor. It uses the provided owner address to set it in the BundleAggregatorProxy constructor.
// Returns a new addressbook with deploy BundleAggregatorProxy contract addresses.
var DeployBundleAggregatorProxyChangeset = cldf.CreateChangeSet(deployBundleAggregatorProxyLogic, deployBundleAggregatorProxyPrecondition)

func deployBundleAggregatorProxyLogic(env deployment.Environment, c types.DeployBundleAggregatorProxyConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.Chains[chainSelector]

		dataFeedsCacheAddress := GetDataFeedsCacheAddress(env.ExistingAddresses, chainSelector, &c.CacheLabel) //nolint:staticcheck // TODO: replace with DataStore when ready
		if dataFeedsCacheAddress == "" {
			return deployment.ChangesetOutput{}, fmt.Errorf("DataFeedsCache contract address not found in addressbook for chain %d", chainSelector)
		}

		bundleProxyResponse, err := DeployBundleAggregatorProxy(chain, common.HexToAddress(dataFeedsCacheAddress), c.Owners[chainSelector], c.Labels)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy BundleAggregatorProxy: %w", err)
		}

		lggr.Infof("Deployed %s chain selector %d addr %s", bundleProxyResponse.Tv.String(), chain.Selector, bundleProxyResponse.Address.String())

		err = ab.Save(chain.Selector, bundleProxyResponse.Address.String(), bundleProxyResponse.Tv)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save BundleAggregatorProxy: %w", err)
		}
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func deployBundleAggregatorProxyPrecondition(env deployment.Environment, c types.DeployBundleAggregatorProxyConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.Chains[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
		_, err := env.ExistingAddresses.AddressesForChain(chainSelector) //nolint:staticcheck // TODO: replace with DataStore when ready
		if err != nil {
			return fmt.Errorf("failed to get addessbook for chain %d: %w", chainSelector, err)
		}
		if !common.IsHexAddress(c.Owners[chainSelector].String()) {
			return fmt.Errorf("owner %s is not a valid address for chain %d", c.Owners[chainSelector].String(), chainSelector)
		}
	}

	return nil
}
