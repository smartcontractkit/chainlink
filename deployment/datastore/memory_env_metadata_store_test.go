package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryEnvMetadataStore_Get(t *testing.T) {
	var (
		recordOne = EnvMetadata[DefaultMetadata]{
			Domain:      "example.com",
			Environment: "test",
			Metadata:    DefaultMetadata{Data: "data1"},
		}
	)

	tests := []struct {
		name              string
		givenState        *EnvMetadata[DefaultMetadata]
		domain            string
		recordShouldExist bool
		expectedRecord    EnvMetadata[DefaultMetadata]
		expectedError     error
	}{
		{
			name:              "env metadata set",
			givenState:        &recordOne,
			domain:            "example.com",
			recordShouldExist: true,
			expectedRecord:    recordOne,
		},
		{
			name:              "env metadata not set",
			domain:            "nonexistent.com",
			recordShouldExist: false,
			expectedRecord:    EnvMetadata[DefaultMetadata]{},
			expectedError:     ErrEnvMetadataNotSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryEnvMetadataStore[DefaultMetadata]{Record: tt.givenState}

			record, err := store.Get()
			if tt.recordShouldExist {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRecord, record)
			} else {
				require.Equal(t, tt.expectedError, err)
				require.Equal(t, tt.expectedRecord, record)
			}
		})
	}
}

func TestMemoryEnvMetadataStore_Set(t *testing.T) {
	var (
		recordOne = EnvMetadata[DefaultMetadata]{
			Domain:      "example.com",
			Environment: "test",
			Metadata:    DefaultMetadata{Data: "data1"},
		}
		recordTwo = EnvMetadata[DefaultMetadata]{
			Domain:      "example2.com",
			Environment: "test2",
			Metadata:    DefaultMetadata{Data: "data2"},
		}
	)

	tests := []struct {
		name           string
		initialState   *EnvMetadata[DefaultMetadata]
		updateRecord   EnvMetadata[DefaultMetadata]
		expectedRecord EnvMetadata[DefaultMetadata]
	}{
		{
			name:           "update existing record",
			initialState:   &recordOne,
			updateRecord:   recordTwo,
			expectedRecord: recordTwo,
		},
		{
			name:           "add new record",
			updateRecord:   recordOne,
			expectedRecord: recordOne,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := MemoryEnvMetadataStore[DefaultMetadata]{Record: tt.initialState}

			err := store.Set(tt.updateRecord)
			require.NoError(t, err)

			record, err := store.Get()
			require.NoError(t, err)
			require.Equal(t, tt.expectedRecord, record)
		})
	}
}
