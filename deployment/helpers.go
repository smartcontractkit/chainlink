package deployment

import (
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func ValidateSelectorsInEnvironment(e cldf.Environment, chains []uint64) error {
	for _, chain := range chains {
		if !e.BlockChains.Exists(chain) {
			return fmt.Errorf("chain %d not found in environment", chain)
		}
	}
	return nil
}

func IsAddressListUnique(addresses []common.Address) bool {
	addressSet := make(map[common.Address]struct{})
	for _, address := range addresses {
		if _, exists := addressSet[address]; exists {
			return false
		}
		addressSet[address] = struct{}{}
	}
	return true
}

func AddressListContainsEmptyAddress(addresses []common.Address) bool {
	return slices.Contains(addresses, (common.Address{}))
}

func MigrateAddressBook(addrBook cldf.AddressBook) (datastore.MutableDataStore, error) {
	addrs, err := addrBook.Addresses()
	if err != nil {
		return nil, err
	}

	ds := datastore.NewMemoryDataStore()

	for chainSelector, chainAddresses := range addrs {
		for addr, typever := range chainAddresses {
			ref := datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          datastore.ContractType(typever.Type),
				Version:       &typever.Version,
				// Since the address book does not have a qualifier, we use the address and type as a
				// unique identifier for the addressRef. Otherwise, we would have some clashes in the
				// between address refs.
				Qualifier: fmt.Sprintf("%s-%s", addr, typever.Type),
			}

			// If the address book has labels, we need to add them to the addressRef
			if !typever.Labels.IsEmpty() {
				ref.Labels = datastore.NewLabelSet(typever.Labels.List()...)
			}

			if err = ds.Addresses().Add(ref); err != nil {
				return nil, err
			}
		}
	}

	return ds, nil
}
