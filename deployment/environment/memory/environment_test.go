package memory

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		})
	}
}

func TestMapChainsToNodes(t *testing.T) {
	tests := []struct {
		name             string
		chainIdxToFChain map[int]int // map[chainGlobalIndex]fChainValue
		numNodes         int
		seed             int64
	}{
		{
			name:             "one chain, all nodes selected",
			chainIdxToFChain: map[int]int{0: 1}, // f=1, 3f+1=4
			numNodes:         4,                 // numNodes == 3f+1
			seed:             123,
		},
		{
			name:             "two chains, same f-value, subset of nodes",
			chainIdxToFChain: map[int]int{0: 1, 1: 1}, // f=1, 3f+1=4
			numNodes:         10,                      // numNodes > 3f+1
			seed:             456,
		},
		{
			name: "multiple chains, different f-values",
			chainIdxToFChain: map[int]int{
				0: 1, // chainIdx 0, f=1 (draws 4 nodes)
				1: 2, // chainIdx 1, f=2 (draws 7 nodes)
			},
			numNodes: 10, // For f=1 (3*1+1=4) valid; For f=2 (3*2+1=7) valid.
			seed:     789,
		},
		{
			name: "complex case with overlapping node assignments",
			chainIdxToFChain: map[int]int{
				0: 1, // chain 0, f=1 -> 4 nodes
				1: 1, // chain 1, f=1 -> 4 nodes (same as for chain 0)
				2: 2, // chain 2, f=2 -> 7 nodes
			},
			numNodes: 8, // f=1: numNodes(8) >= 4. f=2: numNodes(8) >= 7. Valid.
			seed:     101,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Pre-calculate what drawNodesForChain would return for each unique fChain value
			// using the test case's t, numNodes, and seed. This ensures consistency because
			// mapChainsToNodes internally calls drawNodesForChain with the same seed for each fChain lookup.
			drawnNodesPerFChain := make(map[int][]int) // fChainValue -> sorted list of node indices
			uniqueFChainsInTest := make(map[int]struct{})
			for _, fChainVal := range tc.chainIdxToFChain {
				uniqueFChainsInTest[fChainVal] = struct{}{}
			}

			for fChainVal := range uniqueFChainsInTest {
				nodes := drawNodesForChain(t, fChainVal, tc.numNodes, tc.seed)
				sort.Ints(nodes)
				drawnNodesPerFChain[fChainVal] = nodes
			}

			actualResult := mapChainsToNodes(t, tc.chainIdxToFChain, tc.numNodes, tc.seed)

			// Construct expected result based on pre-calculated drawnNodesPerFChain
			expectedResult := make(map[int][]int)
			for chainIdx, fChainVal := range tc.chainIdxToFChain {
				nodesForThisFChain, ok := drawnNodesPerFChain[fChainVal]
				require.True(t, ok, "fChainVal %d was not pre-calculated, something is wrong in test setup", fChainVal)
				for _, nodeIdx := range nodesForThisFChain {
					expectedResult[nodeIdx] = append(expectedResult[nodeIdx], chainIdx)
				}
			}

			// Sort lists in both actual and expected results for consistent comparison
			for nodeIdx := range actualResult {
				sort.Ints(actualResult[nodeIdx])
				// Also assert that node indices are valid
				assert.GreaterOrEqual(t, nodeIdx, 0, "Node index in result map key should be non-negative")
				assert.Less(t, nodeIdx, tc.numNodes, "Node index in result map key should be less than numNodes")
			}
			for nodeIdx := range expectedResult {
				sort.Ints(expectedResult[nodeIdx])
			}

			assert.Equal(t, expectedResult, actualResult, "mapChainsToNodes result mismatch")

			// Idempotency Check: Calling again with same params should yield same result
			actualResult2 := mapChainsToNodes(t, tc.chainIdxToFChain, tc.numNodes, tc.seed)
			for nodeIdx := range actualResult2 {
				sort.Ints(actualResult2[nodeIdx])
			}
			assert.Equal(t, expectedResult, actualResult2, "mapChainsToNodes is not idempotent")
		})
	}
}
