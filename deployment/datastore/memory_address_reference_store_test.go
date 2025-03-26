package datastore

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
)

// MutableStore interface methods tests
func TestInMemoryAddressReferenceStore_indexOf(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressReferenceRecord
		giveKey       AddressReferenceKey
		expectedIndex int
	}{
		{
			name: "success: returns index of record",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: 1,
		},
		{
			name: "success: returns -1 if record not found",
			givenState: []AddressReferenceRecord{
				recordOne,
			},
			giveKey:       recordTwo.Key(),
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
			idx := store.indexOf(tt.giveKey)
			assert.Equal(t, tt.expectedIndex, idx)
		})
	}
}

func TestInMemoryAddressReferenceStore_Add(t *testing.T) {
	var (
		record = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressReferenceRecord
		giveRecord    AddressReferenceRecord
		expectedState []AddressReferenceRecord
		expectedError error
	}{
		{
			name:       "success: adds new record",
			givenState: []AddressReferenceRecord{},
			giveRecord: record,
			expectedState: []AddressReferenceRecord{
				record,
			},
		},
		{
			name: "error: already existing record",
			givenState: []AddressReferenceRecord{
				record,
			},
			giveRecord:    record,
			expectedError: ErrAddressReferenceRecordExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
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

func TestInMemoryAddressReferenceStore_AddOrUpdate(t *testing.T) {
	var (
		oldRecord = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}
		newRecord = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressReferenceRecord
		expectedState []AddressReferenceRecord
		giveRecord    AddressReferenceRecord
	}{
		{
			name:       "success: adds new record",
			givenState: []AddressReferenceRecord{},
			giveRecord: oldRecord,
			expectedState: []AddressReferenceRecord{
				oldRecord,
			},
		},
		{
			name: "success: updates existing record",
			givenState: []AddressReferenceRecord{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []AddressReferenceRecord{
				newRecord,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
			// Check the error for the in-memory store, which will always be nil for the
			// in memory implementation, to satisfy the linter
			err := store.AddOrUpdate(tt.giveRecord)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedState, store.records)
		})
	}
}

func TestInMemoryAddressReferenceStore_Update(t *testing.T) {
	var (
		oldRecord = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}
		newRecord = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressReferenceRecord
		expectedState []AddressReferenceRecord
		giveRecord    AddressReferenceRecord
		expectedError error
	}{
		{
			name: "success: updates existing record",
			givenState: []AddressReferenceRecord{
				oldRecord,
			},
			giveRecord: newRecord,
			expectedState: []AddressReferenceRecord{
				newRecord,
			},
		},
		{
			name:          "error: record not found",
			givenState:    []AddressReferenceRecord{},
			giveRecord:    newRecord,
			expectedError: ErrAddressReferenceRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
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

func TestInMemoryAddressReferenceStore_Delete(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}

		recordThree = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     3,
			Type:      "typeZ",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name          string
		givenState    []AddressReferenceRecord
		expectedState []AddressReferenceRecord
		giveKey       AddressReferenceKey
		expectedError error
	}{
		{
			name: "success: deletes given record",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveKey: recordTwo.Key(),
			expectedState: []AddressReferenceRecord{
				recordOne,
				recordThree,
			},
		},
		{
			name: "error: record not found",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordThree,
			},
			giveKey:       recordTwo.Key(),
			expectedError: ErrAddressReferenceRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
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
func TestInMemoryAddressReferenceStore_Fetch(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name            string
		givenState      []AddressReferenceRecord
		expectedRecords []AddressReferenceRecord
		expectedError   error
	}{
		{
			name: "success: fetches all records",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			expectedRecords: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			expectedError: nil,
		},
		{
			name:            "success: fetches no records",
			givenState:      []AddressReferenceRecord{},
			expectedRecords: []AddressReferenceRecord{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
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

func TestInMemoryAddressReferenceStore_Get(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveKey        AddressReferenceKey
		expectedRecord AddressReferenceRecord
		expectedError  error
	}{
		{
			name: "success: record exists",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveKey:        recordTwo.Key(),
			expectedRecord: recordTwo,
			expectedError:  nil,
		},
		{
			name:           "error: record not found",
			givenState:     []AddressReferenceRecord{},
			giveKey:        recordTwo.Key(),
			expectedRecord: AddressReferenceRecord{},
			expectedError:  ErrAddressReferenceRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
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

func TestInMemoryAddressReferenceStore_Filter(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}

		recordThree = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     3,
			Type:      "typeZ",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveFilters    []FilterFunc[AddressReferenceKey, AddressReferenceRecord]
		expectedResult []AddressReferenceRecord
	}{{
		name: "success: no filters returns all records",
		givenState: []AddressReferenceRecord{
			recordOne,
			recordTwo,
			recordThree,
		},
		giveFilters:    []FilterFunc[AddressReferenceKey, AddressReferenceRecord]{},
		expectedResult: []AddressReferenceRecord{recordOne, recordTwo, recordThree},
	},
		{
			name: "success: returns record with given chain and type",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[AddressReferenceKey, AddressReferenceRecord]{
				AddressReferenceRecordByChain(2),
				AddressReferenceRecordByType("typeX"),
			},
			expectedResult: []AddressReferenceRecord{recordTwo},
		},
		{
			name: "success: returns no record with given chain and type",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
				recordThree,
			},
			giveFilters: []FilterFunc[AddressReferenceKey, AddressReferenceRecord]{
				AddressReferenceRecordByChain(4),
				AddressReferenceRecordByType("typeX"),
			},
			expectedResult: []AddressReferenceRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := InMemoryAddressReferenceStore{records: tt.givenState}
			filteredRecords := store.Filter(tt.giveFilters...)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

// Default Filters tests
func TestByChain(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveChain      uint64
		expectedResult []AddressReferenceRecord
	}{
		{
			name: "success: returns record with given chain",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveChain:      2,
			expectedResult: []AddressReferenceRecord{recordTwo},
		},
		{
			name: "success: returns no record with given chain",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveChain:      5,
			expectedResult: []AddressReferenceRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressReferenceRecordByChain(tt.giveChain)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestByType(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveType       deployment.ContractType
		expectedResult []AddressReferenceRecord
	}{
		{
			name: "success: returns record with given type",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveType: "typeX",
			expectedResult: []AddressReferenceRecord{
				recordTwo,
			},
		},
		{
			name: "success: returns no record with given type",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveType:       "typeL",
			expectedResult: []AddressReferenceRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressReferenceRecordByType(tt.giveType)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestByVersion(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveVersion    semver.Version
		expectedResult []AddressReferenceRecord
	}{
		{
			name: "success: returns record with given version",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveVersion: semver.MustParse("0.5.0"),
			expectedResult: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
		},
		{
			name: "success: returns no record with given version",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveVersion:    semver.MustParse("0.6.0"),
			expectedResult: []AddressReferenceRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressReferenceRecordByVersion(tt.giveVersion)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestByQualifier(t *testing.T) {
	var (
		recordOne = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     1,
			Type:      "type1",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressReferenceRecord{
			Address:   "0x2324224",
			Chain:     2,
			Type:      "typeX",
			Version:   semver.MustParse("0.5.0"),
			Qualifier: "qual2",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressReferenceRecord
		giveQualifier  string
		expectedResult []AddressReferenceRecord
	}{
		{
			name: "success: returns record with given qualifier",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveQualifier: "qual1",
			expectedResult: []AddressReferenceRecord{
				recordOne,
			},
		},
		{
			name: "success: returns no record with given qualifier",
			givenState: []AddressReferenceRecord{
				recordOne,
				recordTwo,
			},
			giveQualifier:  "qual32",
			expectedResult: []AddressReferenceRecord(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressReferenceRecordByQualifier(tt.giveQualifier)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
