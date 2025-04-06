package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContractMetadata_Clone(t *testing.T) {
	original := ContractMetadata[DefaultMetadata]{
		ChainSelector: 1,
		Address:       "0x123",
		Metadata:      DefaultMetadata{Data: "test data"},
	}

	cloned := original.Clone()

	require.Equal(t, original.ChainSelector, cloned.ChainSelector)
	require.Equal(t, original.Address, cloned.Address)
	require.Equal(t, original.Metadata, cloned.Metadata)

	// Modify the original and ensure the cloned remains unchanged
	original.ChainSelector = 2
	original.Address = "0x456"
	original.Metadata = DefaultMetadata{Data: "updated data"}

	require.NotEqual(t, original.ChainSelector, cloned.ChainSelector)
	require.NotEqual(t, original.Address, cloned.Address)
	require.NotEqual(t, original.Metadata, cloned.Metadata)
}

func TestContractMetadata_Key(t *testing.T) {
	metadata := ContractMetadata[DefaultMetadata]{
		ChainSelector: 1,
		Address:       "0x123",
		Metadata:      DefaultMetadata{Data: "test data"},
	}

	key := metadata.Key()
	expectedKey := NewContractMetadataKey(1, "0x123")

	require.Equal(t, expectedKey, key)
}
