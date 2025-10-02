package changeset

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// ImportAddressToDataStoreChangeset is a changeset that reads already deployed contract addresses from input file
// and saves them to the data store. Returns a new datatore with the imported addresses.
var ImportAddressToDataStoreChangeset = cldf.CreateChangeSet(importAddressToDatastoreLogic, importAddressToDatastorePrecondition)

func importAddressToDatastoreLogic(env cldf.Environment, c types.ImportAddressesConfig) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	addresses := c.Addresses

	for _, address := range addresses {
		labels := datastore.NewLabelSet()
		for _, label := range address.Labels {
			labels.Add(label)
		}
		if err := ds.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: c.ChainSelector,
				Address:       address.Address,
				Type:          address.Type,
				Version:       semver.MustParse(address.Version),
				Qualifier:     address.Qualifier,
				Labels:        labels,
			},
		); err != nil {
			return cldf.ChangesetOutput{DataStore: ds},
				fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds}, nil
}

func importAddressToDatastorePrecondition(env cldf.Environment, c types.ImportAddressesConfig) error {
	if !env.BlockChains.Exists(c.ChainSelector) {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if len(c.Addresses) == 0 {
		return errors.New("no addresses to import")
	}

	return nil
}
