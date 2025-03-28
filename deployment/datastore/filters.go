package datastore

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
)

// The following functions are a default set of filters that can be used with the Filter method of the
// AddressRefStore interface. These filters are composable and can be combined to create more complex filters.
// For example, to filter records by chain and contract type, you can use the following:
//	```
//		records := store.Filter(
//			AddressRefByChainSelector(1),
//			AddressRefByType(deployment.ContractType("type1")),
//			AddressRefByVersion("my-qualifier"),
//		)
//	```
// This allows for a more flexible and reusable way to filter records. And opens the possibility for any user
// to create their own custom filters by implementing the FilterFunc type.

// All the filters below are used to filter AddressRef records in the AddressRefStore.
// They all implement the FilterFunc type.
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByChainSelector(0)
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByType(deployment.ContractType(""))
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByVersion(nil)
var _ FilterFunc[AddressRefKey, AddressRef] = AddressRefByQualifier("")

// AddressRefByChainSelector returns a filter that only includes records with the provided chain.
func AddressRefByChainSelector(chainSelector uint64) FilterFunc[AddressRefKey, AddressRef] {
	return func(records []AddressRef) []AddressRef {
		filtered := []AddressRef{}
		for _, record := range records {
			if record.ChainSelector == chainSelector {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressRefByType returns a filter that only includes records with the provided contract type.
func AddressRefByType(contractType deployment.ContractType) FilterFunc[AddressRefKey, AddressRef] {
	return func(records []AddressRef) []AddressRef {
		filtered := []AddressRef{}
		for _, record := range records {
			if record.Type == contractType {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressRefByVersion returns a filter that only includes records with the provided version.
func AddressRefByVersion(version *semver.Version) FilterFunc[AddressRefKey, AddressRef] {
	return func(records []AddressRef) []AddressRef {
		filtered := []AddressRef{}
		for _, record := range records {
			if record.Version.Equal(version) {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressRefByQualifier returns a filter that only includes records with the provided qualifier.
func AddressRefByQualifier(qualifier string) FilterFunc[AddressRefKey, AddressRef] {
	return func(records []AddressRef) []AddressRef {
		filtered := []AddressRef{}
		for _, record := range records {
			if record.Qualifier == qualifier {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}
