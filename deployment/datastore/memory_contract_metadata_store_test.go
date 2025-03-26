package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MutableStore interface methods tests
func TestInMemoryContractMetadataStore_indexOf(t *testing.T) {
	var (
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadataRecord
		giveKey       ContractMetadataKey
		expectedIndex int
	}{
		{
			name: "success: returns index of record",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: 1,
		},
		{
			name: "success: returns -1 if record not found",
			givenState: []ContractMetadataRecord{
				recordOne,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
			idx := store.indexOf(tt.giveKey)
			assert.Equal(t, tt.expectedIndex, idx)
		})
	}
}

func TestInMemoryContractMetadataStore_Add(t *testing.T) {
	var (
		record = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadataRecord
		giveRecord    ContractMetadataRecord
		expectedState []ContractMetadataRecord
		expectedError error
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadataRecord{},
			giveRecord: record,
			expectedState: []ContractMetadataRecord{
				record,
			},
		},
		{
			name: "error: already existing record",
			givenState: []ContractMetadataRecord{
				record,
			},
			giveRecord:    record,
			expectedError: ErrContractMetadataRecordExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		oldRecord = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		newRecord = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadataRecord
		expectedState []ContractMetadataRecord
		giveRecord    ContractMetadataRecord
	}{
		{
			name:       "success: adds new record",
			givenState: []ContractMetadataRecord{},
			giveRecord: oldRecord,
			expectedState: []ContractMetadataRecord{
				oldRecord,
			},
		},
		{
			name: "success: updates existing record",
			givenState: []ContractMetadataRecord{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadataRecord{
				newRecord,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		oldRecord = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		newRecord = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadataRecord
		expectedState []ContractMetadataRecord
		giveRecord    ContractMetadataRecord
		expectedError error
	}{
		{
			name: "success: updates existing record",
			givenState: []ContractMetadataRecord{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []ContractMetadataRecord{
				newRecord,
			},
		},
		{
			name:          "error: record not found",
			givenState:    []ContractMetadataRecord{},
			giveRecord:    newRecord,
			expectedError: ErrContractMetadataRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}

		recordThree = ContractMetadataRecord{
			Chain:    3,
			Address:  "0x2324224",
			Metadata: "metadata3",
		}
	)

	tests := []struct {
		name          string
		givenState    []ContractMetadataRecord
		expectedState []ContractMetadataRecord
		giveKey       ContractMetadataKey
		expectedError error
	}{
		{
			name: "success: deletes given record",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveKey: recordTwo.Key(),
			expectedState: []ContractMetadataRecord{
				recordOne,
				recordThree,
			},
		},
		{
			name: "error: record not found",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordThree,
			},
			giveKey:       recordTwo.Key(),
			expectedError: ErrContractMetadataRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name            string
		givenState      []ContractMetadataRecord
		expectedRecords []ContractMetadataRecord
		expectedError   error
	}{
		{
			name: "success: fetches all records",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			expectedRecords: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			expectedError: nil,
		},
		{
			name:            "success: fetches no records",
			givenState:      []ContractMetadataRecord{},
			expectedRecords: []ContractMetadataRecord{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadataRecord
		giveKey        ContractMetadataKey
		expectedRecord ContractMetadataRecord
		expectedError  error
	}{
		{
			name: "success: record exists",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			giveKey:        recordTwo.Key(),
			expectedRecord: recordTwo,
		},
		{
			name:           "error: record not found",
			givenState:     []ContractMetadataRecord{},
			giveKey:        recordTwo.Key(),
			expectedRecord: ContractMetadataRecord{},
			expectedError:  ErrContractMetadataRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
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
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}

		recordThree = ContractMetadataRecord{
			Chain:    3,
			Address:  "0x2324224",
			Metadata: "metadata3",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadataRecord
		giveFilters    []FilterFunc[ContractMetadataKey, ContractMetadataRecord]
		expectedResult []ContractMetadataRecord
	}{{
		name: "success: no filters returns all records",
		givenState: []ContractMetadataRecord{
			recordOne,
			recordTwo,
			recordThree,
		},
		giveFilters:    []FilterFunc[ContractMetadataKey, ContractMetadataRecord]{},
		expectedResult: []ContractMetadataRecord{recordOne, recordTwo, recordThree},
	},
		{
			name: "success: returns record with given chain and type",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadataRecord]{
				ContractMetadataByChain(2),
			},
			expectedResult: []ContractMetadataRecord{recordTwo},
		},
		{
			name: "success: returns no record with given chain and type",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[ContractMetadataKey, ContractMetadataRecord]{
				ContractMetadataByChain(4),
			},
			expectedResult: []ContractMetadataRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryContractMetadataStore{records: tt.givenState}
			filteredRecords := store.Filter(tt.giveFilters...)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

// // Default Filters tests
func TestMetaByChain(t *testing.T) {
	var (
		recordOne = ContractMetadataRecord{
			Chain:    1,
			Address:  "0x2324224",
			Metadata: "metadata1",
		}

		recordTwo = ContractMetadataRecord{
			Chain:    2,
			Address:  "0x2324224",
			Metadata: "metadata2",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadataRecord
		giveChain      uint64
		expectedResult []ContractMetadataRecord
	}{
		{
			name: "success: returns record with given chain",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			giveChain:      2,
			expectedResult: []ContractMetadataRecord{recordTwo},
		},
		{
			name: "success: returns no record with given chain",
			givenState: []ContractMetadataRecord{
				recordOne,
				recordTwo,
			},
			giveChain:      5,
			expectedResult: []ContractMetadataRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ContractMetadataByChain(tt.giveChain)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
