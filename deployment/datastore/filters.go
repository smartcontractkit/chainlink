package datastore

import (
	"github.com/Masterminds/semver/v3"
)

// The following functions are a default set of filters that can be used with the Filter method of the
// AddressRefStore interface. These filters are composable and can be combined to create more complex filters.
// For example, to filter records by chain and contract type, you can use the following:
//	```
//		records := store.Filter(
//			AddressRefByChainSelector(1),
//			AddressRefByType(ContractType("type1")),
//			AddressRefByVersion("my-qualifier"),
//		)
//	```
// This allows for a more flexible and reusable way to filter records. And opens the possibility for any user
// to create their own custom filters by implementing the FilterFunc type.

// All the filters below are used to filter AddressRef records in the AddressRefStore.
// They all implement the FilterFunc type.
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByChainSelector(0)
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByType(ContractType(""))
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByVersion(nil)
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByQualifier("")

// addressRefFilter returns a filter that includes records for which the predicate returns true.
// This is a generalized filter function that can be used to create custom filters.
func addressRefFilter(predicate func(record AddressRef) bool) FilterFunc[AddressRefKey, AddressRef] {
	return func(records []AddressRef) []AddressRef {
		filtered := make([]AddressRef, 0, len(records)) // Pre-allocate capacity
		for _, record := range records {
			if predicate(record) {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressRefByChainSelector returns a filter that only includes records with the provided chain.
func AddressRefByChainSelector(chainSelector uint64) FilterFunc[AddressRefKey, AddressRef] {
	return addressRefFilter(func(record AddressRef) bool {
		return record.ChainSelector == chainSelector
	})
}

// AddressRefByType returns a filter that only includes records with the provided contract type.
func AddressRefByType(contractType ContractType) FilterFunc[AddressRefKey, AddressRef] {
	return addressRefFilter(func(record AddressRef) bool {
		return record.Type == contractType
	})
}

// AddressRefByVersion returns a filter that only includes records with the provided version.
func AddressRefByVersion(version *semver.Version) FilterFunc[AddressRefKey, AddressRef] {
	return addressRefFilter(func(record AddressRef) bool {
		return record.Version.Equal(version)
	})
}

// AddressRefByQualifier returns a filter that only includes records with the provided qualifier.
func AddressRefByQualifier(qualifier string) FilterFunc[AddressRefKey, AddressRef] {
	return addressRefFilter(func(record AddressRef) bool {
		return record.Qualifier == qualifier
	})
}
