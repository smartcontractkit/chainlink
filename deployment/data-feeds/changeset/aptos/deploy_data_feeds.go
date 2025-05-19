package aptos

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployDataFeedsChangeset deploys the ChainlinkDataFeeds package to Aptos chain. Uses platform package address from existing datastore.
// Returns a new addressbook and datastore with deployed router/registry contracts addresses.
var DeployDataFeedsChangeset = cldf.CreateChangeSet(deployDataFeedsLogic, deployDataFeedsPrecondition)

func deployDataFeedsLogic(env cldf.Environment, c types.DeployAptosConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.AptosChains[chainSelector]

		// Use the owner address if provided, otherwise use the deployer signer address
		ownerAddress := chain.DeployerSigner.AccountAddress()
		if c.Owner != "" {
			ownerAddress = aptos.AccountAddress{}
			err := ownerAddress.ParseStringRelaxed(c.Owner)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}

		// Get the platform address from the datastore
		record, _ := env.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(chainSelector, "ChainlinkPlatform", semver.MustParse("1.0.0"), "aptos"),
		)

		platformAccountAddress := aptos.AccountAddress{}
		err := platformAccountAddress.ParseStringRelaxed(record.Address)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		dataFeedsResponse, err := DeployDataFeeds(chain, ownerAddress, platformAccountAddress, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ChainlinkDataFeeds: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", dataFeedsResponse.Tv.String(), chain.Selector, dataFeedsResponse.Address.String())

		err = ab.Save(chain.Selector, dataFeedsResponse.Address.String(), dataFeedsResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ChainlinkDataFeeds: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       dataFeedsResponse.Address.String(),
				Type:          "ChainlinkDataFeeds",
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     "aptos",
				Labels:        datastore.NewLabelSet(c.Labels...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}

func deployDataFeedsPrecondition(env cldf.Environment, c types.DeployAptosConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.AptosChains[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}

		_, err := env.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(chainSelector, "ChainlinkPlatform", semver.MustParse("1.0.0"), "aptos"),
		)
		if err != nil {
			return errors.New("ChainlinkPlatform not found in data store")
		}
	}

	return nil
}
