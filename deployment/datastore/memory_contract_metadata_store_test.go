package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryContractMetadataStore_indexOf(t *testing.T) {
	var (
		recordOne = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata1"),
		}

		recordTwo = ContractMetadata[DefaultMetadata]{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata2"),
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata[DefaultMetadata]
		giveKey       ContractMetadataKey
		expectedIndex int
	}{
		{
			name: "success: returns index of record",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: 1,
		},
		{
			name: "success: returns -1 if record not found",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
			idx := store.indexOf(tt.giveKey)
			assert.Equal(t, tt.expectedIndex, idx)
		})
	}
}

func TestMemoryContractMetadataStore_Add(t *testing.T) {
	var (
		record = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata[DefaultMetadata]
		giveRecord    ContractMetadata[DefaultMetadata]
		expectedState []ContractMetadata[DefaultMetadata]
		expectedError error
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadata[DefaultMetadata]{},
			giveRecord: record,
			expectedState: []ContractMetadata[DefaultMetadata]{
				record,
			},
		},
		{
			name: "error: already existing record",
			givenState: []ContractMetadata[DefaultMetadata]{
				record,
			},
			giveRecord:    record,
			expectedError: ErrContractMetadataExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
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

func TestMemoryContractMetadataStore_AddOrUpdate(t *testing.T) {
	var (
		oldRecord = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata1"),
		}

		newRecord = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata2"),
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata[DefaultMetadata]
		expectedState []ContractMetadata[DefaultMetadata]
		giveRecord    ContractMetadata[DefaultMetadata]
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadata[DefaultMetadata]{},
			giveRecord: oldRecord,
			expectedState: []ContractMetadata[DefaultMetadata]{
				oldRecord,
			},
		},
		{
			name: "success: updates existing record",
			givenState: []ContractMetadata[DefaultMetadata]{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadata[DefaultMetadata]{
				newRecord,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
			// Check the error for the in-memory store, which will always be nil for the
			// in memory implementation, to satisfy the linter
			err := store.AddOrUpdate(tt.giveRecord)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedState, store.Records)
		})
	}
}

func TestMemoryContractMetadataStore_Update(t *testing.T) {
	var (
		oldRecord = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata1"),
		}

		newRecord = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata2"),
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata[DefaultMetadata]
		expectedState []ContractMetadata[DefaultMetadata]
		giveRecord    ContractMetadata[DefaultMetadata]
		expectedError error
	}{
		{
			name: "success: updates existing record",
			givenState: []ContractMetadata[DefaultMetadata]{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadata[DefaultMetadata]{
				newRecord,
			},
		},
		{
			name:          "error: record not found",
			givenState:    []ContractMetadata[DefaultMetadata]{},
			giveRecord:    newRecord,
			expectedError: ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
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

func TestMemoryMemoryContractMetadataStore_Delete(t *testing.T) {
	var (
		recordOne = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata[DefaultMetadata]{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}

		recordThree = ContractMetadata[DefaultMetadata]{
			ChainSelector: 3,
			Address:       "0x2324224",
			Metadata:      "metadata3",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadata[DefaultMetadata]
		expectedState []ContractMetadata[DefaultMetadata]
		giveKey       ContractMetadataKey
		expectedError error
	}{
		{
			name: "success: deletes given record",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveKey: recordTwo.Key(),
			expectedState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordThree,
			},
		},
		{
			name: "error: record not found",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordThree,
			},
			giveKey:       recordTwo.Key(),
			expectedError: ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
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

func TestMemoryContractMetadataStore_Fetch(t *testing.T) {
	var (
		recordOne = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata[DefaultMetadata]{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name            string
		givenState      []ContractMetadata[DefaultMetadata]
		expectedRecords []ContractMetadata[DefaultMetadata]
		expectedError   error
	}{
		{
			name: "success: fetches all records",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
			},
			expectedRecords: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
			},
			expectedError: nil,
		},
		{
			name:            "success: fetches no records",
			givenState:      []ContractMetadata[DefaultMetadata]{},
			expectedRecords: []ContractMetadata[DefaultMetadata]{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
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

func TestMemoryContractMetadataStore_Get(t *testing.T) {
	var (
		recordOne = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata[DefaultMetadata]{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadata[DefaultMetadata]
		giveKey        ContractMetadataKey
		expectedRecord ContractMetadata[DefaultMetadata]
		expectedError  error
	}{
		{
			name: "success: record exists",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
			},
			giveKey:        recordTwo.Key(),
			expectedRecord: recordTwo,
		},
		{
			name:           "error: record not found",
			givenState:     []ContractMetadata[DefaultMetadata]{},
			giveKey:        recordTwo.Key(),
			expectedRecord: ContractMetadata[DefaultMetadata]{},
			expectedError:  ErrContractMetadataNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
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

func TestMemoryContractMetadataStore_Filter(t *testing.T) {
	var (
		recordOne = ContractMetadata[DefaultMetadata]{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata1"),
		}

		recordTwo = ContractMetadata[DefaultMetadata]{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata2"),
		}

		recordThree = ContractMetadata[DefaultMetadata]{
			ChainSelector: 3,
			Address:       "0x2324224",
			Metadata:      DefaultMetadata("metadata3"),
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadata[DefaultMetadata]
		giveFilters    []FilterFunc[ContractMetadataKey, ContractMetadata[DefaultMetadata]]
		expectedResult []ContractMetadata[DefaultMetadata]
	}{{
		name: "success: no filters returns all records",
		givenState: []ContractMetadata[DefaultMetadata]{
			recordOne,
			recordTwo,
			recordThree,
		},
		giveFilters:    []FilterFunc[ContractMetadataKey, ContractMetadata[DefaultMetadata]]{},
		expectedResult: []ContractMetadata[DefaultMetadata]{recordOne, recordTwo, recordThree},
	},
		{
			name: "success: returns record with given chain and type",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadata[DefaultMetadata]]{
				ContractMetadataByChainSelector[DefaultMetadata](2),
			},
			expectedResult: []ContractMetadata[DefaultMetadata]{recordTwo},
		},
		{
			name: "success: returns no record with given chain and type",
			givenState: []ContractMetadata[DefaultMetadata]{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadata[DefaultMetadata]]{
				ContractMetadataByChainSelector[DefaultMetadata](4),
			},
			expectedResult: []ContractMetadata[DefaultMetadata]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryContractMetadataStore[DefaultMetadata]{Records: tt.givenState}
			filteredRecords := store.Filter(tt.giveFilters...)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
