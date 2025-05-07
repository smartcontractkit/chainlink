package aptos

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// DeployPlatformChangeset deploys the ChainlinkPlatform package to Aptos chain.
// Returns a new addressbook and data store with deployed forwarder/storage contracts addresses.
var DeployPlatformChangeset = cldf.CreateChangeSet(deployPlatformLogic, deployPlatformPrecondition)

func deployPlatformLogic(env deployment.Environment, c types.DeployAptosConfig) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.AptosChains[chainSelector]

		// Use the owner address if provided, otherwise use the deployer signer address
		ownerAddress := chain.DeployerSigner.AccountAddress()
		if c.Owner != "" {
			ownerAddress = aptos.AccountAddress{}
			err := ownerAddress.ParseStringRelaxed(c.Owner)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}

		platformResponse, err := DeployPlatform(chain, ownerAddress, c.Labels)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy ChainlinkPlatform: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", platformResponse.Tv.String(), chain.Selector, platformResponse.Address.String())

		err = ab.Save(chain.Selector, platformResponse.Address.String(), platformResponse.Tv)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save ChainlinkPlatform: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       platformResponse.Address.String(),
				Type:          "ChainlinkPlatform",
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     "aptos",
				Labels:        datastore.NewLabelSet(c.Labels...),
			},
		); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}

func deployPlatformPrecondition(env deployment.Environment, c types.DeployAptosConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.AptosChains[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
