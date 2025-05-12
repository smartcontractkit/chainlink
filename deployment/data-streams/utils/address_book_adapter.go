package utils

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// AddressBookToDataStore converts an AddressBook to a DataStore
// AddressBook is deprecated and will be removed in the future. You can use this function to migrate or interact with legacy code that uses AddressBook.
// You can call this with your custom metadata like this:
//
//	ds, err := cldf.AddressBookToDataStore[metadata.CustomContractMetaData, datastore.CustomEnvMetaData](addressBook)
func AddressBookToDataStore[CM datastore.Cloneable[CM], EM datastore.Cloneable[EM]](ab cldf.AddressBook) (datastore.DataStore[CM, EM], error) {
	ds := datastore.NewMemoryDataStore[CM, EM]()

	// Get all addresses from the AddressBook
	addresses, err := ab.Addresses()
	if err != nil {
		return nil, err
	}

	// For each address, create an AddressRef and add it to the DataStore
	for chainSelector, chainAddresses := range addresses {
		for address, tv := range chainAddresses {
			// Create an AddressRef with the chain selector stored in ChainSelector field
			addressRef := datastore.AddressRef{
				Address:       address,
				ChainSelector: chainSelector,
				Labels:        datastore.NewLabelSet(tv.Labels.List()...),
				Type:          datastore.ContractType(tv.Type),
				Version:       &tv.Version,
			}

			// Add the AddressRef to the DataStore
			err := ds.Addresses().Upsert(addressRef)
			if err != nil {
				return nil, err
			}
		}
	}

	return ds.Seal(), nil
}

// DataStoreToAddressBook converts a DataStore to an AddressBook
// DataStore ContractMetadata and EnvMetadata are not preserved in the AddressBook.
// AddressBook is deprecated and will be removed in the future. You can use this function to migrate or interact with legacy code that uses AddressBook.
func DataStoreToAddressBook[CM datastore.Cloneable[CM], EM datastore.Cloneable[EM]](ds datastore.DataStore[CM, EM]) (cldf.AddressBook, error) {
	ab := cldf.NewMemoryAddressBook()

	// Get all addresses from the DataStore
	addressRefs, err := ds.Addresses().Fetch()
	if err != nil {
		return nil, err
	}

	// For each address, create a TypeAndVersion and add it to the AddressBook
	for _, addressRef := range addressRefs {
		// Create a TypeAndVersion
		tv := cldf.TypeAndVersion{
			Type:    cldf.ContractType(addressRef.Type),
			Version: *addressRef.Version,
			Labels:  cldf.NewLabelSet(addressRef.Labels.List()...),
		}

		// Add the TypeAndVersion to the AddressBook
		err := ab.Save(addressRef.ChainSelector, addressRef.Address, tv)
		if err != nil {
			return nil, err
		}
	}

	return ab, nil
}

// AddressBookToNewDataStore converts an AddressBook to a new mutable DataStore
// AddressBook is deprecated and will be removed in the future. You can use this function to migrate or interact with legacy code that uses AddressBook.
func AddressBookToNewDataStore[CM datastore.Cloneable[CM], EM datastore.Cloneable[EM]](ab cldf.AddressBook) (*datastore.MemoryDataStore[CM, EM], error) {
	ds := datastore.NewMemoryDataStore[CM, EM]()

	// Get all addresses from the AddressBook
	addresses, err := ab.Addresses()
	if err != nil {
		return nil, err
	}

	// For each address, create an AddressRef and add it to the DataStore
	for chainSelector, chainAddresses := range addresses {
		for address, tv := range chainAddresses {
			addressRef := datastore.AddressRef{
				Address:       address,
				ChainSelector: chainSelector,
				Labels:        datastore.NewLabelSet(tv.Labels.List()...),
				Type:          datastore.ContractType(tv.Type),
				Version:       &tv.Version,
			}

			err := ds.Addresses().Upsert(addressRef)
			if err != nil {
				return nil, err
			}
		}
	}

	return ds, nil
}

func AddressRefsToAddressBook(addressRefs []datastore.AddressRef) (cldf.AddressBook, error) {
	ab := cldf.NewMemoryAddressBook()
	for _, addressRef := range addressRefs {
		tv := cldf.TypeAndVersion{
			Type:    cldf.ContractType(addressRef.Type),
			Version: *addressRef.Version,
			Labels:  cldf.NewLabelSet(addressRef.Labels.List()...),
		}

		err := ab.Save(addressRef.ChainSelector, addressRef.Address, tv)
		if err != nil {
			return nil, err
		}
	}

	return ab, nil
}

func AddressRefsToAddressByChain(addressRefs []datastore.AddressRef) cldf.AddressesByChain {
	addressesByChain := make(cldf.AddressesByChain)
	for _, addressRef := range addressRefs {
		if addressesByChain[addressRef.ChainSelector] == nil {
			addressesByChain[addressRef.ChainSelector] = make(map[string]cldf.TypeAndVersion)
		}
		addressesByChain[addressRef.ChainSelector][addressRef.Address] = cldf.TypeAndVersion{
			Type:    cldf.ContractType(addressRef.Type),
			Version: *addressRef.Version,
			Labels:  cldf.NewLabelSet(addressRef.Labels.List()...),
		}
	}
	return addressesByChain
}
