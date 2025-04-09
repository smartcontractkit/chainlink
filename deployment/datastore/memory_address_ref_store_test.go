package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryAddressRefStore_indexOf(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressRef
		giveKey       AddressRefKey
		expectedIndex int
	}{
		{
			name: "success: returns index of record",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: 1,
		},
		{
			name: "success: returns -1 if record not found",
			givenState: []AddressRef{
				recordOne,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			idx := store.indexOf(tt.giveKey)
			assert.Equal(t, tt.expectedIndex, idx)
		})
	}
}

func TestMemoryAddressRefStore_Add(t *testing.T) {
	var (
		record = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressRef
		giveRecord    AddressRef
		expectedState []AddressRef
		expectedError error
	}{
		{
			name:       "success: adds new record",
			givenState: []AddressRef{},
			giveRecord: record,
			expectedState: []AddressRef{
				record,
			},
		},
		{
			name: "error: already existing record",
			givenState: []AddressRef{
				record,
			},
			giveRecord:    record,
			expectedError: ErrAddressRefExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			err := store.Add(tt.giveRecord)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.Records)
			}
		})
	}
}

func TestMemoryAddressRefStore_AddOrUpdate(t *testing.T) {
	var (
		oldRecord = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}
		newRecord = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressRef
		expectedState []AddressRef
		giveRecord    AddressRef
	}{
		{
			name:       "success: adds new record",
			givenState: []AddressRef{},
			giveRecord: oldRecord,
			expectedState: []AddressRef{
				oldRecord,
			},
		},
		{
			name: "success: updates existing record",
			givenState: []AddressRef{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []AddressRef{
				newRecord,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			// Check the error, which will always be nil for the
			// in memory implementation, to satisfy the linter
			err := store.AddOrUpdate(tt.giveRecord)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedState, store.Records)
		})
	}
}

func TestMemoryAddressRefStore_Update(t *testing.T) {
	var (
		oldRecord = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}
		newRecord = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressRef
		expectedState []AddressRef
		giveRecord    AddressRef
		expectedError error
	}{
		{
			name: "success: updates existing record",
			givenState: []AddressRef{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []AddressRef{
				newRecord,
			},
		},
		{
			name:          "error: record not found",
			givenState:    []AddressRef{},
			giveRecord:    newRecord,
			expectedError: ErrAddressRefNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			err := store.Update(tt.giveRecord)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.Records)
			}
		})
	}
}

func TestMemoryAddressRefStore_Delete(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}

		recordThree = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 3,
			Type:          "typeZ",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressRef
		expectedState []AddressRef
		giveKey       AddressRefKey
		expectedError error
	}{
		{
			name: "success: deletes given record",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveKey: recordTwo.Key(),
			expectedState: []AddressRef{
				recordOne,
				recordThree,
			},
		},
		{
			name: "error: record not found",
			givenState: []AddressRef{
				recordOne,
				recordThree,
			},
			giveKey:       recordTwo.Key(),
			expectedError: ErrAddressRefNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			err := store.Delete(tt.giveKey)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.Records)
			}
		})
	}
}

func TestMemoryAddressRefStore_Fetch(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name            string
		givenState      []AddressRef
		expectedRecords []AddressRef
	}{
		{
			name: "success: fetches all records",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			expectedRecords: []AddressRef{
				recordOne,
				recordTwo,
			},
		},
		{
			name:            "success: fetches no records",
			givenState:      []AddressRef{},
			expectedRecords: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			// Check the error, which will always be nil for the
			// in memory implementation, to satisfy the linter
			records, err := store.Fetch()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedRecords, records)
		})
	}
}

func TestMemoryAddressRefStore_Get(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveKey        AddressRefKey
		expectedRecord AddressRef
		expectedError  error
	}{
		{
			name: "success: record exists",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveKey:        recordTwo.Key(),
			expectedRecord: recordTwo,
			expectedError:  nil,
		},
		{
			name:           "error: record not found",
			givenState:     []AddressRef{},
			giveKey:        recordTwo.Key(),
			expectedRecord: AddressRef{},
			expectedError:  ErrAddressRefNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			record, err := store.Get(tt.giveKey)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRecord, record)
			}
		})
	}
}

func TestMemoryAddressRefStore_Filter(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}

		recordThree = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 3,
			Type:          "typeZ",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveFilters    []FilterFunc[AddressRefKey, AddressRef]
		expectedResult []AddressRef
	}{
		{
			name: "success: no filters returns all records",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters:    []FilterFunc[AddressRefKey, AddressRef]{},
			expectedResult: []AddressRef{recordOne, recordTwo, recordThree},
		},
		{
			name: "success: returns record with given chain and type",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[AddressRefKey, AddressRef]{
				AddressRefByChainSelector(2),
				AddressRefByType("typeX"),
			},
			expectedResult: []AddressRef{recordTwo},
		},
		{
			name: "success: returns no record with given chain and type",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[AddressRefKey, AddressRef]{
				AddressRefByChainSelector(4),
				AddressRefByType("typeX"),
			},
			expectedResult: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryAddressRefStore{Records: tt.givenState}
			filteredRecords := store.Filter(tt.giveFilters...)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
