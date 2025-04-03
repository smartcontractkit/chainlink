package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractMetadata_Clone(t *testing.T) {
	original := ContractMetadata[DefaultMetadata]{
		ChainSelector: 1,
		Address:       "0x123",
		Metadata:      DefaultMetadata("test data"),
	}

	cloned := original.Clone()

	assert.Equal(t, original.ChainSelector, cloned.ChainSelector)
	assert.Equal(t, original.Address, cloned.Address)
	assert.Equal(t, original.Metadata, cloned.Metadata)

	// Modify the original and ensure the cloned remains unchanged
	original.ChainSelector = 2
	original.Address = "0x456"
	original.Metadata = DefaultMetadata("updated data")

	assert.NotEqual(t, original.ChainSelector, cloned.ChainSelector)
	assert.NotEqual(t, original.Address, cloned.Address)
	assert.NotEqual(t, original.Metadata, cloned.Metadata)
}

func TestContractMetadata_Key(t *testing.T) {
	metadata := ContractMetadata[DefaultMetadata]{
		ChainSelector: 1,
		Address:       "0x123",
		Metadata:      DefaultMetadata("test data"),
	}

	key := metadata.Key()
	expectedKey := NewContractMetadataKey(1, "0x123")

	assert.Equal(t, expectedKey, key)
}
