package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestMemoryDataStore_Merge(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*MemoryDataStore[DefaultMetadata], *MemoryDataStore[DefaultMetadata])
		expectedCount int
		expectedLabel string
	}{
		{
			name: "Merge single address",
			setup: func() (*MemoryDataStore[DefaultMetadata], *MemoryDataStore[DefaultMetadata]) {
				dataStore1 := NewMemoryDataStore[DefaultMetadata]()
				dataStore2 := NewMemoryDataStore[DefaultMetadata]()
				err := dataStore2.Addresses().AddOrUpdate(AddressRef{
					Address:   "0x123",
					Type:      "type1",
					Version:   semver.MustParse("1.0.0"),
					Qualifier: "qualifier1",
				})
				require.NoError(t, err, "Adding data to dataStore2 should not fail")
				return dataStore1, dataStore2
			},
			expectedCount: 1,
		},
		{
			name: "Match existing address with labels",
			setup: func() (*MemoryDataStore[DefaultMetadata], *MemoryDataStore[DefaultMetadata]) {
				dataStore1 := NewMemoryDataStore[DefaultMetadata]()
				dataStore2 := NewMemoryDataStore[DefaultMetadata]()

				// Add initial data to dataStore1
				err := dataStore1.Addresses().AddOrUpdate(AddressRef{
					Address:   "0x123",
					Type:      "type1",
					Version:   semver.MustParse("1.0.0"),
					Qualifier: "qualifier1",
					Labels:    NewLabelSet("label1"),
				})
				require.NoError(t, err, "Adding initial data to dataStore1 should not fail")

				// Add matching data to dataStore2
				err = dataStore2.Addresses().AddOrUpdate(AddressRef{
					Address:   "0x123",
					Type:      "type1",
					Version:   semver.MustParse("1.0.0"),
					Qualifier: "qualifier1",
					Labels:    NewLabelSet("label2"),
				})
				require.NoError(t, err, "Adding matching data to dataStore2 should not fail")

				return dataStore1, dataStore2
			},
			expectedCount: 1,
			expectedLabel: "label2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataStore1, dataStore2 := tt.setup()

			// Merge dataStore2 into dataStore1
			err := dataStore1.Merge(dataStore2.Seal())
			require.NoError(t, err, "Merging dataStore2 into dataStore1 should not fail")

			// Verify that dataStore1 contains the merged data
			addressRefs, err := dataStore1.Addresses().Fetch()
			require.NoError(t, err, "Fetching addresses from dataStore1 should not fail")
			require.Len(t, addressRefs, tt.expectedCount, "dataStore1 should contain the expected number of addresses after merge")
			require.Equal(t, tt.expectedLabel, addressRefs[0].Labels.String(), "Labels should be updated correctly after merge")
		})
	}
}
