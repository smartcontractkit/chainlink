package deployment

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// Implementations have been migrated to the chainlink-deployments-framework
// Using type alias here to avoid updating all the references in the codebase.
// This file will be removed in the future once we migrate all the code
type (
	ContractType     = deployment.ContractType
	TypeAndVersion   = deployment.TypeAndVersion
	AddressBook      = deployment.AddressBook
	AddressesByChain = deployment.AddressesByChain
	AddressBookMap   = deployment.AddressBookMap
	LabelSet         = deployment.LabelSet
)

var (
	TypeAndVersionFromString     = deployment.TypeAndVersionFromString
	NewTypeAndVersion            = deployment.NewTypeAndVersion
	MustTypeAndVersionFromString = deployment.MustTypeAndVersionFromString
	NewMemoryAddressBook         = deployment.NewMemoryAddressBook
	NewMemoryAddressBookFromMap  = deployment.NewMemoryAddressBookFromMap
	SearchAddressBook            = deployment.SearchAddressBook
	AddressBookContains          = deployment.AddressBookContains
	EnsureDeduped                = deployment.EnsureDeduped
	NewLabelSet                  = deployment.NewLabelSet

	ErrInvalidChainSelector = deployment.ErrInvalidChainSelector
	ErrInvalidAddress       = deployment.ErrInvalidAddress
	ErrChainNotFound        = deployment.ErrChainNotFound
)

func MigrateAddressBook(addrBook deployment.AddressBook) (datastore.MutableDataStore[datastore.DefaultMetadata,
	datastore.DefaultMetadata], error) {
	addrs, err := addrBook.Addresses()
	if err != nil {
		return nil, err
	}

	ds := datastore.NewMemoryDataStore[
		datastore.DefaultMetadata,
		datastore.DefaultMetadata,
	]()

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
