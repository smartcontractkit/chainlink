package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestMemoryDataStore_Merge(t *testing.T) {
	// Create two mutable MemoryDataStores
	dataStore1 := NewMemoryDataStore()
	dataStore2 := NewMemoryDataStore()

	// Add some data to the second data store
	err := dataStore2.Addresses().AddOrUpdate(AddressRef{
		Address:   "0x123",
		Type:      "type1",
		Version:   semver.MustParse("1.0.0"),
		Qualifier: "qualifier1",
	})
	require.NoError(t, err, "Adding data to dataStore2 should not fail")

	// Merge dataStore2 into dataStore1
	err = dataStore1.Merge(dataStore2)
	require.NoError(t, err, "Merging dataStore2 into dataStore1 should not fail")

	// Verify that dataStore1 contains the merged data
	addressRefs, err := dataStore1.Addresses().Fetch()
	require.NoError(t, err, "Fetching addresses from dataStore1 should not fail")
	require.Len(t, addressRefs, 1, "dataStore1 should contain 1 address after merge")
	require.Equal(t, "0x123", addressRefs[0].Address, "Merged address should match")
}

func TestSealedMemoryDataStore_Merge(t *testing.T) {
	// Create a mutable MemoryDataStore and seal it
	mutableDataStore := NewMemoryDataStore()
	sealedDataStore := mutableDataStore.Seal()

	// Create another mutable MemoryDataStore with some data
	otherDataStore := NewMemoryDataStore()
	err := otherDataStore.Addresses().AddOrUpdate(AddressRef{
		Address:   "0x456",
		Type:      "type2",
		Version:   semver.MustParse("1.0.0"),
		Qualifier: "qualifier2",
	})
	require.NoError(t, err, "Adding data to otherDataStore should not fail")

	// Merge otherDataStore into the sealed data store
	err = sealedDataStore.Merge(otherDataStore.Seal())
	require.NoError(t, err, "Merging into a sealed data store should not fail")

	// Verify that the sealed data store contains the merged data
	addressRefs, err := sealedDataStore.Addresses().Fetch()
	require.NoError(t, err, "Fetching addresses from sealedDataStore should not fail")
	require.Len(t, addressRefs, 1, "sealedDataStore should contain 1 address after merge")
	require.Equal(t, "0x456", addressRefs[0].Address, "Merged address should match")
}
