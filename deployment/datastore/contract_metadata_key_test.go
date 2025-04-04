package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContractMetadataKey_Equals(t *testing.T) {
	tests := []struct {
		name     string
		key1     ContractMetadataKey
		key2     ContractMetadataKey
		expected bool
	}{
		{
			name:     "Equal keys",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			expected: true,
		},
		{
			name:     "Different chain selector",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(2, "0x1234567890abcdef"),
			expected: false,
		},
		{
			name:     "Different address",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(1, "0xabcdef1234567890"),
			expected: false,
		},
		{
			name:     "Completely different keys",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(2, "0xabcdef1234567890"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.key1.Equals(tt.key2))
		})
	}
}

func TestContractMetadataKey(t *testing.T) {
	chainSelector := uint64(1)
	address := "0x1234567890abcdef"

	key := NewContractMetadataKey(chainSelector, address)

	require.Equal(t, chainSelector, key.ChainSelector(), "ChainSelector should return the correct chain selector")
	require.Equal(t, address, key.Address(), "Address should return the correct address")
}
