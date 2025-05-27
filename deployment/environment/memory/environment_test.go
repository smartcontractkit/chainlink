package memory

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawNodesForChain(t *testing.T) {
	tests := []struct {
		name      string
		fChain    int
		numNodes  int
		seed      int64
		expectedN int
	}{
		{
			name:      "basic case: fChain=1, numNodes=10",
			fChain:    1, // 3*1+1 = 4. numNodes(10) >= 4. Valid.
			numNodes:  10,
			seed:      123,
			expectedN: 4,
		},
		{
			name:      "fChain=2, numNodes=7 (all nodes selected)",
			fChain:    2, // 3*2+1 = 7. numNodes(7) >= 7. Valid.
			numNodes:  7,
			seed:      42,
			expectedN: 7,
		},
		{
			name:      "numNodes exactly 3f+1: fChain=1, numNodes=4 (all nodes selected)",
			fChain:    1, // 3*1+1 = 4. numNodes(4) >= 4. Valid.
			numNodes:  4,
			seed:      789,
			expectedN: 4,
		},
		{
			name:      "large fChain: fChain=5, numNodes=50",
			fChain:    5, // 3*5+1 = 16. numNodes(50) >= 16. Valid.
			numNodes:  50,
			seed:      101,
			expectedN: 16,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nodeIdxs := drawNodesForChain(t, tc.fChain, tc.numNodes, tc.seed)

			// 1. Check the number of nodes drawn
			assert.Len(t, nodeIdxs, tc.expectedN, "Number of drawn nodes mismatch")

			// 2. Verify that all drawn node indices are valid and unique
			seen := make(map[int]bool)
			for _, idx := range nodeIdxs {
				assert.GreaterOrEqual(t, idx, 0, "Node index should be non-negative")
				assert.Less(t, idx, tc.numNodes, "Node index should be less than numNodes")
				assert.False(t, seen[idx], "Node index %d was drawn more than once", idx)
				seen[idx] = true
			}

			// 3. Idempotency: Same parameters should yield same result
			nodeIdxs2 := drawNodesForChain(t, tc.fChain, tc.numNodes, tc.seed)
			sort.Ints(nodeIdxs) // Sort for consistent comparison as order isn't guaranteed beyond seed determinism
			sort.Ints(nodeIdxs2)
			assert.Equal(t, nodeIdxs, nodeIdxs2, "Drawing with the same seed should produce the same set of nodes")

			// 4. Different seed behavior
			nodeIdxs3 := drawNodesForChain(t, tc.fChain, tc.numNodes, tc.seed+1)
			sort.Ints(nodeIdxs3)

			if tc.expectedN < tc.numNodes {
				// If we are selecting a SUBSET of available nodes, different seeds should produce different subsets.
				assert.NotEqual(t, nodeIdxs, nodeIdxs3,
					"Drawing with a different seed (%d vs %d) should produce different nodes when drawing a subset. fChain:%d, numNodes:%d, expectedN:%d. Got %v vs %v",
					tc.seed, tc.seed+1, tc.fChain, tc.numNodes, tc.expectedN, nodeIdxs, nodeIdxs3)
			} else {
				// If tc.expectedN == tc.numNodes (i.e., all available nodes are selected),
				// then different seeds will still result in all nodes being selected.
				assert.Equal(t, nodeIdxs, nodeIdxs3,
					"Drawing with a different seed (%d vs %d) but selecting all nodes should result in the same set of all nodes. fChain:%d, numNodes:%d, expectedN:%d. Got %v vs %v",
					tc.seed, tc.seed+1, tc.fChain, tc.numNodes, tc.expectedN, nodeIdxs, nodeIdxs3)
			}
		})
	}
}
