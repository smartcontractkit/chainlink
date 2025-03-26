package datastore

import "github.com/smartcontractkit/chainlink/deployment"

// The following functions are a default set of filters that can be used with the Filter method of the
// AddressReferenceStore AddressReferenceStore interface. These filters are composable and can be combined
// to create more complex filters.
// For example, to filter records by chain and contract type, you can use the following:
//	```
//		records := store.Filter(
//			AddressReferenceRecordByChain(1),
//			AddressReferenceRecordByType(deployment.ContractType("type1")),
//			AddressReferenceRecordByVersion("my-qualifier"),
//		)
//	```
// This allows for a more flexible and reusable way to filter records. And opens the possibility for any user
// to create their own custom filters by implementing the FilterFunc type.

// AddressReferenceRecordByChain returns a filter that only includes records with the provided chain.
func AddressReferenceRecordByChain(chain uint64) FilterFunc[AddressReferenceKey, AddressReferenceRecord] {
	return func(records []AddressReferenceRecord) []AddressReferenceRecord {
		var filtered []AddressReferenceRecord
		for _, record := range records {
			if record.Chain == chain {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressReferenceRecordByType returns a filter that only includes records with the provided contract type.
func AddressReferenceRecordByType(contractType deployment.ContractType) FilterFunc[AddressReferenceKey, AddressReferenceRecord] {
	return func(records []AddressReferenceRecord) []AddressReferenceRecord {
		var filtered []AddressReferenceRecord
		for _, record := range records {
			if record.Type == contractType {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressReferenceRecordByVersion returns a filter that only includes records with the provided version.
func AddressReferenceRecordByVersion(version string) FilterFunc[AddressReferenceKey, AddressReferenceRecord] {
	return func(records []AddressReferenceRecord) []AddressReferenceRecord {
		var filtered []AddressReferenceRecord
		for _, record := range records {
			if record.Version.String() == version {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}

// AddressReferenceRecordByQualifier returns a filter that only includes records with the provided qualifier.
func AddressReferenceRecordByQualifier(qualifier string) FilterFunc[AddressReferenceKey, AddressReferenceRecord] {
	return func(records []AddressReferenceRecord) []AddressReferenceRecord {
		var filtered []AddressReferenceRecord
		for _, record := range records {
			if record.Qualifier == qualifier {
				filtered = append(filtered, record)
			}
		}
		return filtered
	}
}
