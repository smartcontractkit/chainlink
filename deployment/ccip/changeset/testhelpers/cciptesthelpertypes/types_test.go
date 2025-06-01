package cciptesthelpertypes

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/cciptesthelpertypes/mocks"
)

func chainSelectorSlice(n int) []cciptypes.ChainSelector {
	selectors := make([]cciptypes.ChainSelector, n)
	for i := range selectors {
		//nolint:gosec // no overflow risk here
		selectors[i] = cciptypes.ChainSelector(uint64(i) + 1000)
	}
	return selectors
}

func TestRandomTopology_getChainToFChainMapping(t *testing.T) {
	type args struct {
		chainSelectors []cciptypes.ChainSelector
	}
	tests := []struct {
		name                   string
		fields                 RandomTopologyArgs
		args                   args
		wantFChainDistribution map[int]int
		wantErr                bool
		wantErrMsgSubstring    string
	}{
		{
			name: "basic case with mixed fChain values",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 2, 2: 1}, // 2 chains with fChain 1, 1 chain with fChain 2
				Seed:              12345,
			},
			args: args{
				chainSelectors: chainSelectorSlice(3),
			},
			wantFChainDistribution: map[int]int{1: 2, 2: 1},
			wantErr:                false,
		},
		{
			name: "empty chain selectors",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1}, // This FChainToNumChains is for a scenario with chains, but test provides none.
				Seed:              1,
			},
			args: args{
				chainSelectors: []cciptypes.ChainSelector{},
			},
			wantFChainDistribution: map[int]int{},
			wantErr:                true,
			wantErrMsgSubstring:    "the number of fChains must be equal to the number of chainSelectors, len(fChains) = 1, len(chainSelectors) = 0",
		},
		{
			name: "all chains with same fChain value",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 3},
				Seed:              42,
			},
			args: args{
				chainSelectors: chainSelectorSlice(3),
			},
			wantFChainDistribution: map[int]int{1: 3},
			wantErr:                false,
		},
		{
			name: "more complex fChain distribution",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1, 2: 2, 3: 1}, // 1 f1, 2 f2, 1 f3 = 4 chains
				Seed:              987,
			},
			args: args{
				chainSelectors: chainSelectorSlice(4),
			},
			wantFChainDistribution: map[int]int{1: 1, 2: 2, 3: 1},
			wantErr:                false,
		},
		{
			name: "more skewed fChain distribution",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{3: 1, 4: 100}, // 1 chain with fChain 3, 100 chains with fChain 4
				Seed:              987,
			},
			args: args{
				chainSelectors: chainSelectorSlice(101), // 101 chains in total
			},
			wantFChainDistribution: map[int]int{3: 1, 4: 100},
			wantErr:                false,
		},
		{
			name: "error - not enough fChain assignments for chain selectors",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1}, // Only 1 fChain assignment available
				Seed:              777,
			},
			args: args{
				// 2 chain selectors, but only 1 assignment configured above
				chainSelectors: chainSelectorSlice(2),
			},
			wantErr:             true,
			wantErrMsgSubstring: "the number of fChains must be equal to the number of chainSelectors, len(fChains) = 1, len(chainSelectors) = 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &randomTopology{
				RandomTopologyArgs: tt.fields,
			}
			gen := rand.New(rand.NewSource(tt.fields.Seed))
			got, err := r.getChainToFChainMapping(gen, tt.args.chainSelectors)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsgSubstring != "" {
					require.Contains(t, err.Error(), tt.wantErrMsgSubstring)
				}
				return // Do not proceed with other checks if an error is expected
			}
			require.NoError(t, err) // Expect no error for non-wantErr cases

			// Verify the number of returned mappings
			require.Len(t, got, len(tt.args.chainSelectors), "RandomTopology.getChainToFChainMapping() returned map with incorrect number of keys")

			// Verify all original chain selectors are present as keys
			// For an empty chainSelectors slice, this loop won't run, which is correct.
			for _, cs := range tt.args.chainSelectors {
				require.Contains(t, got, cs, "RandomTopology.getChainToFChainMapping() response missing chain selector %d", cs)
			}

			// Verify the distribution of fChain values
			// For an empty chainSelectors slice, wantFChainDistribution should be an empty map.
			gotFChainDistribution := make(map[int]int)
			for _, fChain := range got {
				gotFChainDistribution[fChain]++
			}
			require.Equal(t, tt.wantFChainDistribution, gotFChainDistribution, "RandomTopology.getChainToFChainMapping() fChain distribution mismatch")
		})
	}
}

// Helper function to generate unique P2P IDs for testing
func generateTestP2PIDs(count int) [][32]byte {
	p2pIDs := make([][32]byte, count)
	for i := range count {
		p2pIDs[i][0] = byte(i + 1) // Simple unique IDs
	}
	return p2pIDs
}

func TestRandomTopology_getNodesForChain(t *testing.T) {
	type args struct {
		fChain             int
		nonBootstrapP2pIDs [][32]byte
	}
	tests := []struct {
		name          string
		fields        RandomTopologyArgs
		args          args
		wantNodeCount int
		wantErr       bool
		wantErrMsg    string
	}{
		{
			name:   "happy path, fChain 1, enough nodes",
			fields: RandomTopologyArgs{Seed: 1},
			args: args{
				fChain:             1,
				nonBootstrapP2pIDs: generateTestP2PIDs(10),
			},
			wantNodeCount: NChain(1), // 3*1 + 1 = 4
			wantErr:       false,
		},
		{
			name:   "fChain 0, enough nodes",
			fields: RandomTopologyArgs{Seed: 2},
			args: args{
				fChain:             0,
				nonBootstrapP2pIDs: generateTestP2PIDs(5),
			},
			wantNodeCount: NChain(0), // 3*0 + 1 = 1
			wantErr:       false,
		},
		{
			name:   "fChain 2, exact number of nodes",
			fields: RandomTopologyArgs{Seed: 3},
			args: args{
				fChain:             2,
				nonBootstrapP2pIDs: generateTestP2PIDs(NChain(2)), // 3*2 + 1 = 7 nodes
			},
			wantNodeCount: NChain(2),
			wantErr:       false,
		},
		{
			name:   "insufficient nodes",
			fields: RandomTopologyArgs{Seed: 4},
			args: args{
				fChain:             1,                     // needs 4 nodes
				nonBootstrapP2pIDs: generateTestP2PIDs(3), // only 3 available
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the number of non-bootstrap ccip nodes must be at least %d, got %d", NChain(1), 3),
		},
		{
			name:   "no nodes available, fChain 0", // needs 1 node
			fields: RandomTopologyArgs{Seed: 5},
			args: args{
				fChain:             0,
				nonBootstrapP2pIDs: generateTestP2PIDs(0),
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the number of non-bootstrap ccip nodes must be at least %d, got %d", NChain(0), 0),
		},
		{
			name:   "fChain 3, more nodes than default MinRoleDONSize",
			fields: RandomTopologyArgs{Seed: 6},
			args: args{
				fChain:             3, // needs 3*3 + 1 = 10 nodes
				nonBootstrapP2pIDs: generateTestP2PIDs(15),
			},
			wantNodeCount: NChain(3),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &randomTopology{
				RandomTopologyArgs: tt.fields,
			}
			gen := rand.New(rand.NewSource(tt.fields.Seed))

			gotNodes, err := r.getNodesForChain(gen, tt.args.fChain, tt.args.nonBootstrapP2pIDs)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
				return // No further checks if an error is expected
			}

			require.NoError(t, err)
			require.NotNil(t, gotNodes)
			require.Len(t, gotNodes, tt.wantNodeCount, "Unexpected number of nodes returned")

			// Check for uniqueness
			seenNodes := make(map[[32]byte]bool)
			for _, nodeID := range gotNodes {
				require.False(t, seenNodes[nodeID], "Duplicate node ID found: %v", nodeID)
				seenNodes[nodeID] = true
			}

			// Check if all returned nodes are from the original set
			inputNodesSet := make(map[[32]byte]bool)
			for _, nodeID := range tt.args.nonBootstrapP2pIDs {
				inputNodesSet[nodeID] = true
			}
			for _, nodeID := range gotNodes {
				require.True(t, inputNodesSet[nodeID], "Returned node ID not found in input set: %v", nodeID)
			}

			// Test determinism: run again with the same seed and expect same result
			// (Important: re-initialize gen with the same seed)
			if len(tt.args.nonBootstrapP2pIDs) > 0 && tt.wantNodeCount > 0 { // Only if nodes are expected and available
				genDeterministic := rand.New(rand.NewSource(tt.fields.Seed))
				gotNodesDeterministic, detErr := r.getNodesForChain(genDeterministic, tt.args.fChain, tt.args.nonBootstrapP2pIDs)
				require.NoError(t, detErr)
				require.ElementsMatch(t, gotNodes, gotNodesDeterministic, "Node selection is not deterministic for the same seed")
			}
		})
	}
}

func TestRandomTopology_validate(t *testing.T) {
	const homeChainSelector = 1
	type args struct {
		chainSelectors    []cciptypes.ChainSelector
		homeChainSelector cciptypes.ChainSelector
	}
	tests := []struct {
		name       string
		fields     RandomTopologyArgs
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "valid: sum of FChainToNumChains equals len(chainSelectors)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 2, 2: 1}, // Total 3 chains
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100, 200, 300}, // 3 selectors
				homeChainSelector: homeChainSelector,
			},
			wantErr: false,
		},
		{
			name: "invalid: sum of FChainToNumChains greater than len(chainSelectors)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 3, 2: 1}, // Total 4 chains
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100, 200, 300}, // 3 selectors
				homeChainSelector: homeChainSelector,
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the sum of the number of chains in the chain topology must be equal to the number of chains provided in the config, got %d, expected %d", 4, 3),
		},
		{
			name: "invalid: sum of FChainToNumChains less than len(chainSelectors)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1, 2: 1}, // Total 2 chains
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100, 200, 300}, // 3 selectors
				homeChainSelector: homeChainSelector,
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the sum of the number of chains in the chain topology must be equal to the number of chains provided in the config, got %d, expected %d", 2, 3),
		},
		{
			name: "valid: FChainToNumChains empty, chainSelectors empty",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{},
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{},
				homeChainSelector: homeChainSelector,
			},
			wantErr: false,
		},
		{
			name: "invalid: FChainToNumChains not empty, chainSelectors empty",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1}, // Total 1 chain
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{}, // 0 selectors
				homeChainSelector: homeChainSelector,
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the sum of the number of chains in the chain topology must be equal to the number of chains provided in the config, got %d, expected %d", 1, 0),
		},
		{
			name: "invalid: FChainToNumChains empty, chainSelectors not empty",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{},
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100}, // 1 selector
				homeChainSelector: homeChainSelector,
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the sum of the number of chains in the chain topology must be equal to the number of chains provided in the config, got %d, expected %d", 0, 1),
		},
		{
			name: "valid: FChainToNumChains contains zero counts",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 2, 2: 0, 3: 1}, // Total 3 chains (0 is ignored in sum)
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100, 200, 300}, // 3 selectors
				homeChainSelector: homeChainSelector,
			},
			wantErr: false, // The method sums values, so 0 doesn't change the sum, this is valid.
		},
		{
			name: "invalid: home chain selector included in chainSelectors",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1, 2: 1},
			},
			args: args{
				chainSelectors:    []cciptypes.ChainSelector{100, homeChainSelector},
				homeChainSelector: homeChainSelector,
			},
			wantErr:    true,
			wantErrMsg: fmt.Sprintf("the home chain selector %d is included in the chainSelectors, please remove it from the chainSelectors", homeChainSelector),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &randomTopology{
				RandomTopologyArgs: tt.fields,
			}
			err := r.validate(tt.args.chainSelectors, tt.args.homeChainSelector)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRandomTopology_ChainToNodeMapping(t *testing.T) {
	const homeChainSelector = 1
	type args struct {
		nonBootstrapP2pIDs [][32]byte
		chainSelectors     []cciptypes.ChainSelector
		homeChainSelector  cciptypes.ChainSelector
	}
	tests := []struct {
		name    string
		fields  RandomTopologyArgs
		args    args
		wantErr bool
		// wantErrMsg is a substring if specific, otherwise general error check
		wantErrMsgSubstring string
		// checkNodeCounts if true, verifies the distribution of node counts per chain
		checkNodeCounts bool
	}{
		{
			name: "happy path - mixed fChains, sufficient nodes",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 1, 2: 1}, // one chain f=1 (4 nodes), one chain f=2 (7 nodes)
				Seed:              101,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(10), // Needs max 7 nodes for one chain, plus others
				chainSelectors:     chainSelectorSlice(2),
				homeChainSelector:  homeChainSelector,
			},
			wantErr:         false,
			checkNodeCounts: true,
		},
		{
			name: "error - insufficient total nodes (MinRoleDONSize)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{0: 1}, // f=0 needs 1 node
				Seed:              102,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(MinRoleDONSize - 1), // e.g., 3 nodes
				chainSelectors:     chainSelectorSlice(1),
				homeChainSelector:  homeChainSelector,
			},
			wantErr:             true,
			wantErrMsgSubstring: fmt.Sprintf("number of non-bootstrap ccip nodes must be at least %d", MinRoleDONSize),
		},
		{
			name: "error - validation failure (chain count mismatch)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 2}, // expects 2 chains configured
				Seed:              103,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(5),
				chainSelectors:     chainSelectorSlice(1), // only 1 provided
				homeChainSelector:  homeChainSelector,
			},
			wantErr:             true,
			wantErrMsgSubstring: "the sum of the number of chains in the chain topology must be equal",
		},
		{
			name: "error - insufficient nodes for a specific fChain",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{2: 1}, // f=2 needs 7 nodes
				Seed:              104,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(6), // only 6 nodes, MinRoleDONSize is met (4)
				chainSelectors:     chainSelectorSlice(1),
				homeChainSelector:  homeChainSelector,
			},
			wantErr:             true,
			wantErrMsgSubstring: fmt.Sprintf("failed to get nodes for chain %s: the number of non-bootstrap ccip nodes must be at least %d, got %d", cciptypes.ChainSelector(1000), NChain(2), 6),
		},
		{
			name: "empty chain selectors, valid FChainToNumChains (empty)",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{},
				Seed:              105,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(5),
				chainSelectors:     []cciptypes.ChainSelector{},
				homeChainSelector:  homeChainSelector,
			},
			wantErr:         false,
			checkNodeCounts: true, // Expected distribution will be empty
		},
		{
			name: "multiple chains, same fChain, sufficient nodes",
			fields: RandomTopologyArgs{
				FChainToNumChains: map[int]int{1: 2}, // two chains, f=1 (4 nodes each)
				Seed:              106,
			},
			args: args{
				nonBootstrapP2pIDs: generateTestP2PIDs(8),
				chainSelectors:     chainSelectorSlice(2),
				homeChainSelector:  homeChainSelector,
			},
			wantErr:         false,
			checkNodeCounts: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &randomTopology{
				RandomTopologyArgs: tt.fields,
			}

			gotMapping, err := r.ChainToNodeMapping(tt.args.nonBootstrapP2pIDs, tt.args.chainSelectors, tt.args.homeChainSelector)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsgSubstring != "" {
					require.Contains(t, err.Error(), tt.wantErrMsgSubstring)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, gotMapping)
			require.Len(t, gotMapping, len(tt.args.chainSelectors)+1)

			for _, cs := range tt.args.chainSelectors {
				require.Contains(t, gotMapping, cs)
			}

			require.Contains(t, gotMapping, tt.args.homeChainSelector)
			require.Equal(t, tt.args.nonBootstrapP2pIDs, gotMapping[tt.args.homeChainSelector])

			if tt.checkNodeCounts {
				// Calculate expected distribution of node counts
				expectedNodeCountsDistribution := make(map[int]int) // map[numNodesForChain]countOfChains
				for fChainVal, numChainsWithThisF := range tt.fields.FChainToNumChains {
					expectedNodeCountsDistribution[NChain(fChainVal)] += numChainsWithThisF
				}

				// Calculate actual distribution of node counts
				actualNodeCountsDistribution := make(map[int]int)
				for chainSel, assignedNodes := range gotMapping {
					if chainSel == tt.args.homeChainSelector {
						continue // Skip the home chain for this specific distribution check
					}
					actualNodeCountsDistribution[len(assignedNodes)]++
				}
				require.Equal(t, expectedNodeCountsDistribution, actualNodeCountsDistribution, "Node count distribution mismatch")
			}

			// Per-chain node checks
			inputNodesSet := make(map[[32]byte]bool)
			for _, nodeID := range tt.args.nonBootstrapP2pIDs {
				inputNodesSet[nodeID] = true
			}

			for cs, assignedNodes := range gotMapping {
				// Check uniqueness of nodes for this chain
				seenChainNodes := make(map[[32]byte]bool)
				for _, nodeID := range assignedNodes {
					require.False(t, seenChainNodes[nodeID], "Duplicate node ID %v for chain %d", nodeID, cs)
					seenChainNodes[nodeID] = true
					// Check if node is from original pool
					require.True(t, inputNodesSet[nodeID], "Node ID %v for chain %d not from input pool", nodeID, cs)
				}
			}

			// Determinism check for successful cases
			if len(tt.args.chainSelectors) > 0 { // Only if there's something to compare
				rDeterministic := &randomTopology{
					RandomTopologyArgs: tt.fields,
				}
				gotMappingDeterministic, detErr := rDeterministic.ChainToNodeMapping(tt.args.nonBootstrapP2pIDs, tt.args.chainSelectors, tt.args.homeChainSelector)
				require.NoError(t, detErr, "Determinism check failed on second run")
				require.Len(t, gotMapping, len(gotMappingDeterministic), "Determinism check: map lengths differ")

				for cs, nodes1 := range gotMapping {
					nodes2, ok := gotMappingDeterministic[cs]
					require.True(t, ok, "Determinism check: chain selector %d missing in second run map", cs)
					// The order of nodes within the slice from getNodesForChain should be deterministic due to rand.Seed
					require.Equal(t, nodes1, nodes2, "Determinism check: node assignments for chain %d differ", cs)
				}
			}
		})
	}
}

func TestTopology_ChainToNodeMapping(t *testing.T) {
	type args struct {
		nonBootstrapP2pIDs    [][32]byte
		nonHomeChainSelectors []cciptypes.ChainSelector
		homeChainSelector     cciptypes.ChainSelector
	}
	tests := []struct {
		name                string
		args                args
		mockSetup           func(m *mocks.RoleDONTopology, args args) map[cciptypes.ChainSelector][][32]byte
		wantErr             bool
		wantErrMsgSubstring string
		wantMapping         map[cciptypes.ChainSelector][][32]byte
	}{
		{
			name: "happy path - mock returns valid mapping",
			args: args{
				nonBootstrapP2pIDs:    generateTestP2PIDs(5),
				nonHomeChainSelectors: []cciptypes.ChainSelector{100, 200},
				homeChainSelector:     300,
			},
			mockSetup: func(m *mocks.RoleDONTopology, args args) map[cciptypes.ChainSelector][][32]byte {
				mockReturnMapping := map[cciptypes.ChainSelector][][32]byte{
					100: {args.nonBootstrapP2pIDs[0], args.nonBootstrapP2pIDs[1]},
					200: {args.nonBootstrapP2pIDs[2], args.nonBootstrapP2pIDs[3]},
					// Home chain not included in mock's direct return, as topology wrapper handles it.
				}
				m.On("ChainToNodeMapping", args.nonBootstrapP2pIDs, args.nonHomeChainSelectors, args.homeChainSelector).
					Return(mockReturnMapping, nil)

				// Construct the expected final mapping for assertion
				expectedMapping := make(map[cciptypes.ChainSelector][][32]byte)
				maps.Copy(expectedMapping, mockReturnMapping)
				expectedMapping[args.homeChainSelector] = args.nonBootstrapP2pIDs
				return expectedMapping
			},
			wantErr: false,
		},
		{
			name: "error path - mock returns error",
			args: args{
				nonBootstrapP2pIDs:    generateTestP2PIDs(4),
				nonHomeChainSelectors: []cciptypes.ChainSelector{100},
				homeChainSelector:     200,
			},
			mockSetup: func(m *mocks.RoleDONTopology, args args) map[cciptypes.ChainSelector][][32]byte {
				m.On("ChainToNodeMapping", args.nonBootstrapP2pIDs, args.nonHomeChainSelectors, args.homeChainSelector).
					Return(nil, errors.New("mock error"))
				return nil
			},
			wantErr:             true,
			wantErrMsgSubstring: "mock error",
		},
		{
			name: "happy path - empty nonHomeChainSelectors",
			args: args{
				nonBootstrapP2pIDs:    generateTestP2PIDs(MinRoleDONSize),
				nonHomeChainSelectors: []cciptypes.ChainSelector{},
				homeChainSelector:     500,
			},
			mockSetup: func(m *mocks.RoleDONTopology, args args) map[cciptypes.ChainSelector][][32]byte {
				mockReturnMapping := map[cciptypes.ChainSelector][][32]byte{
					// No non-home chains in mock's direct return
				}
				m.On("ChainToNodeMapping", args.nonBootstrapP2pIDs, args.nonHomeChainSelectors, args.homeChainSelector).
					Return(mockReturnMapping, nil)

				expectedMapping := make(map[cciptypes.ChainSelector][][32]byte)
				expectedMapping[args.homeChainSelector] = args.nonBootstrapP2pIDs
				return expectedMapping
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopology := mocks.NewRoleDONTopology(t)

			// Set up expectations and get the final expected mapping
			tt.wantMapping = tt.mockSetup(mockTopology, tt.args)

			top := &topology{
				impl: mockTopology,
			}

			gotMapping, err := top.ChainToNodeMapping(tt.args.nonBootstrapP2pIDs, tt.args.nonHomeChainSelectors, tt.args.homeChainSelector)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsgSubstring != "" {
					require.Contains(t, err.Error(), tt.wantErrMsgSubstring)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantMapping, gotMapping)

			// Explicitly check that home chain is mapped to all nonBootstrapP2pIDs
			if !tt.wantErr {
				require.Contains(t, gotMapping, tt.args.homeChainSelector)
				require.ElementsMatch(t, tt.args.nonBootstrapP2pIDs, gotMapping[tt.args.homeChainSelector])
			}
		})
	}
}
