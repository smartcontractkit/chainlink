package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryDomainMetadataStore_Get(t *testing.T) {
	var (
		recordOne = DomainMetadata[DefaultMetadata]{
			Domain:      "example.com",
			Environment: "test",
			Metadata:    DefaultMetadata{Data: "data1"},
		}
	)

	tests := []struct {
		name              string
		givenState        []DomainMetadata[DefaultMetadata]
		domain            string
		recordShouldExist bool
		expectedRecord    DomainMetadata[DefaultMetadata]
	}{
		{
			name: "domain exists",
			givenState: []DomainMetadata[DefaultMetadata]{
				recordOne,
			},
			domain:            "example.com",
			recordShouldExist: true,
			expectedRecord:    recordOne,
		},
		{
			name:              "domain does not exist",
			givenState:        []DomainMetadata[DefaultMetadata]{},
			domain:            "nonexistent.com",
			recordShouldExist: false,
			expectedRecord:    DomainMetadata[DefaultMetadata]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryDomainMetadataStore[DefaultMetadata]{Records: tt.givenState}
			record, ok, err := store.Get()

			if tt.recordShouldExist {
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, tt.expectedRecord, record)
			} else {
				require.NoError(t, err)
				require.False(t, ok)
				require.Equal(t, tt.expectedRecord, record)
			}
		})
	}
}

func TestMemoryDomainMetadataStore_Update(t *testing.T) {
	var (
		recordOne = DomainMetadata[DefaultMetadata]{
			Domain:      "example.com",
			Environment: "test",
			Metadata:    DefaultMetadata{Data: "data1"},
		}
		recordTwo = DomainMetadata[DefaultMetadata]{
			Domain:      "example2.com",
			Environment: "test2",
			Metadata:    DefaultMetadata{Data: "data2"},
		}
	)

	tests := []struct {
		name           string
		initialState   []DomainMetadata[DefaultMetadata]
		updateRecord   DomainMetadata[DefaultMetadata]
		expectedRecord DomainMetadata[DefaultMetadata]
	}{
		{
			name:           "update existing record",
			initialState:   []DomainMetadata[DefaultMetadata]{recordOne},
			updateRecord:   recordTwo,
			expectedRecord: recordTwo,
		},
		{
			name:           "add new record",
			initialState:   []DomainMetadata[DefaultMetadata]{},
			updateRecord:   recordOne,
			expectedRecord: recordOne,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryDomainMetadataStore[DefaultMetadata]()
			for _, record := range tt.initialState {
				err := store.Update(record)
				require.NoError(t, err)
			}

			err := store.Update(tt.updateRecord)
			require.NoError(t, err)

			record, ok, err := store.Get()
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, tt.expectedRecord, record)
		})
	}
}
