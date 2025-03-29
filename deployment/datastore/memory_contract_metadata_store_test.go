package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MutableStore interface methods tests
func TestInMemoryContractMetadataStore_indexOf(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata
		giveKey       ContractMetadataKey
		expectedIndex int
	}{
		{
			name: "success: returns index of record",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: 1,
		},
		{
			name: "success: returns -1 if record not found",
			givenState: []ContractMetadata{
				recordOne,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			idx := store.indexOf(tt.giveKey)
			assert.Equal(t, tt.expectedIndex, idx)
		})
	}
}

func TestInMemoryContractMetadataStore_Add(t *testing.T) {
	var (
		record = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata
		giveRecord    ContractMetadata
		expectedState []ContractMetadata
		expectedError error
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadata{},
			giveRecord: record,
			expectedState: []ContractMetadata{
				record,
			},
		},
		{
			name: "error: already existing record",
			givenState: []ContractMetadata{
				record,
			},
			giveRecord:    record,
			expectedError: ErrContractMetadataExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			err := store.Add(tt.giveRecord)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.records)
			}
		})
	}
}

func TestInMemoryContractMetadataStore_AddOrUpdate(t *testing.T) {
	var (
		oldRecord = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		newRecord = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata
		expectedState []ContractMetadata
		giveRecord    ContractMetadata
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadata{},
			giveRecord: oldRecord,
			expectedState: []ContractMetadata{
				oldRecord,
			},
		},
		{
			name: "success: updates existing record",
			givenState: []ContractMetadata{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadata{
				newRecord,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			// Check the error for the in-memory store, which will always be nil for the
			// in memory implementation, to satisfy the linter
			err := store.AddOrUpdate(tt.giveRecord)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedState, store.records)
		})
	}
}

func TestInMemoryContractMetadataStore_Update(t *testing.T) {
	var (
		oldRecord = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		newRecord = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata
		expectedState []ContractMetadata
		giveRecord    ContractMetadata
		expectedError error
	}{
		{
			name: "success: updates existing record",
			givenState: []ContractMetadata{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadata{
				newRecord,
			},
		},
		{
			name:          "error: record not found",
			givenState:    []ContractMetadata{},
			giveRecord:    newRecord,
			expectedError: ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			err := store.Update(tt.giveRecord)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.records)
			}
		})
	}
}

func TestInMemoryMemoryContractMetadataStore_Delete(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}

		recordThree = ContractMetadata{
			ChainSelector: 3,
			Address:       "0x2324224",
			Metadata:      "metadata3",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata
		expectedState []ContractMetadata
		giveKey       ContractMetadataKey
		expectedError error
	}{
		{
			name: "success: deletes given record",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveKey: recordTwo.Key(),
			expectedState: []ContractMetadata{
				recordOne,
				recordThree,
			},
		},
		{
			name: "error: record not found",
			givenState: []ContractMetadata{
				recordOne,
				recordThree,
			},
			giveKey:       recordTwo.Key(),
			expectedError: ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			err := store.Delete(tt.giveKey)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedState, store.records)
			}
		})
	}
}

// Store interface methods tests
func TestInMemoryContractMetadataStore_Fetch(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name            string
		givenState      []ContractMetadata
		expectedRecords []ContractMetadata
		expectedError   error
	}{
		{
			name: "success: fetches all records",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			expectedRecords: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			expectedError: nil,
		},
		{
			name:            "success: fetches no records",
			givenState:      []ContractMetadata{},
			expectedRecords: []ContractMetadata{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			records, err := store.Fetch()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRecords, records)
			}
		})
	}
}

func TestInMemoryContractMetadataStore_Get(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadata
		giveKey        ContractMetadataKey
		expectedRecord ContractMetadata
		expectedError  error
	}{
		{
			name: "success: record exists",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			giveKey:        recordTwo.Key(),
			expectedRecord: recordTwo,
		},
		{
			name:           "error: record not found",
			givenState:     []ContractMetadata{},
			giveKey:        recordTwo.Key(),
			expectedRecord: ContractMetadata{},
			expectedError:  ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
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

func TestInMemoryContractMetadataStore_Filter(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}

		recordThree = ContractMetadata{
			ChainSelector: 3,
			Address:       "0x2324224",
			Metadata:      "metadata3",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadata
		giveFilters    []FilterFunc[ContractMetadataKey, ContractMetadata]
		expectedResult []ContractMetadata
	}{{
		name: "success: no filters returns all records",
		givenState: []ContractMetadata{
			recordOne,
			recordTwo,
			recordThree,
		},
		giveFilters:    []FilterFunc[ContractMetadataKey, ContractMetadata]{},
		expectedResult: []ContractMetadata{recordOne, recordTwo, recordThree},
	},
		{
			name: "success: returns record with given chain and type",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadata]{
				ContractMetadataByChainSelector(2),
			},
			expectedResult: []ContractMetadata{recordTwo},
		},
		{
			name: "success: returns no record with given chain and type",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadata]{
				ContractMetadataByChainSelector(4),
			},
			expectedResult: []ContractMetadata(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore{records: tt.givenState}
			filteredRecords := store.Filter(tt.giveFilters...)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
