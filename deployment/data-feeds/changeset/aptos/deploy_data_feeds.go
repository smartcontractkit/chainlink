package aptos

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployDataFeedsChangeset deploys the ChainlinkDataFeeds package to Aptos chain.
// Returns a new addressbook and datastore with deployed router/registry contracts addresses.
var DeployDataFeedsChangeset = cldf.CreateChangeSet(deployDataFeedsLogic, deployDataFeedsPrecondition)

func deployDataFeedsLogic(env cldf.Environment, c types.DeployAptosConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	dataStore := datastore.NewMemoryDataStore()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.AptosChains()[chainSelector]

		// Use the owner address if provided, otherwise use the deployer signer address
		ownerAddress := chain.DeployerSigner.AccountAddress()
		if c.OwnerAddress != "" {
			ownerAddress = aptos.AccountAddress{}
			err := ownerAddress.ParseStringRelaxed(c.OwnerAddress)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}

		platformAccountAddress := aptos.AccountAddress{}
		_ = platformAccountAddress.ParseStringRelaxed(c.PlatformAddress)

		secondaryPlatformAccountAddress := aptos.AccountAddress{}
		_ = secondaryPlatformAccountAddress.ParseStringRelaxed(c.SecondaryPlatformAddress)

		dataFeedsResponse, err := DeployDataFeeds(chain, ownerAddress, platformAccountAddress, secondaryPlatformAccountAddress, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ChainlinkDataFeeds: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", dataFeedsResponse.Tv.String(), chain.Selector, dataFeedsResponse.Address.String())

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       dataFeedsResponse.Address.String(),
				Type:          cs.DataFeedsCache,
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(c.Labels...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{DataStore: dataStore}, nil
}

func deployDataFeedsPrecondition(env cldf.Environment, c types.DeployAptosConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.AptosChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}

		platformAccountAddress := aptos.AccountAddress{}
		err := platformAccountAddress.ParseStringRelaxed(c.PlatformAddress)
		if err != nil {
			return err
		}

		secondaryPlatformAccountAddress := aptos.AccountAddress{}
		err = secondaryPlatformAccountAddress.ParseStringRelaxed(c.SecondaryPlatformAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
