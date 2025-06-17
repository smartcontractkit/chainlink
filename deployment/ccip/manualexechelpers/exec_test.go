package manualexechelpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/manualexechelpers"
)

func TestRootMap_AddAndGet(t *testing.T) {
	// Mock data
	root1 := manualexechelpers.RootCacheEntry{Root: offramp.InternalMerkleRoot{MinSeqNr: 1, MaxSeqNr: 10}, BlockNumber: 1}
	root2 := manualexechelpers.RootCacheEntry{Root: offramp.InternalMerkleRoot{MinSeqNr: 11, MaxSeqNr: 20}, BlockNumber: 2}
	root3 := manualexechelpers.RootCacheEntry{Root: offramp.InternalMerkleRoot{MinSeqNr: 21, MaxSeqNr: 30}, BlockNumber: 3}

	// Initialize RootMap
	rm := manualexechelpers.NewRootCache()

	// Add roots
	rm.Add(root1)
	rm.Add(root2)
	rm.Add(root3)

	// Build the RootMap
	rm.Build()

	// Test cases
	tests := []struct {
		key      uint64
		expected manualexechelpers.RootCacheEntry
		found    bool
	}{
		{key: 5, expected: root1, found: true},
		{key: 15, expected: root2, found: true},
		{key: 25, expected: root3, found: true},
		{key: 31, expected: manualexechelpers.RootCacheEntry{}, found: false},
		{key: 0, expected: manualexechelpers.RootCacheEntry{}, found: false},
	}

	for _, test := range tests {
		result, found := rm.Get(test.key)
		require.Equal(t, test.found, found)
		if found {
			require.Equal(t, test.expected.Root.MinSeqNr, result.MinSeqNr)
			require.Equal(t, test.expected.Root.MaxSeqNr, result.MaxSeqNr)
		}
	}
}
