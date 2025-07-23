package cciptesthelpertypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestDistributedTopologyArgs_validate(t *testing.T) {
	selectorA := cciptypes.ChainSelector(1337)
	selectorB := cciptypes.ChainSelector(2337)

	cases := []struct {
		desc      string
		fVals     map[cciptypes.ChainSelector]uint8
		selectors []cciptypes.ChainSelector
		expectErr bool
	}{
		{
			desc:      "empty FValues map",
			fVals:     map[cciptypes.ChainSelector]uint8{},
			selectors: []cciptypes.ChainSelector{selectorA},
			expectErr: true,
		},
		{
			desc:      "missing selector in FValues",
			fVals:     map[cciptypes.ChainSelector]uint8{selectorA: 2},
			selectors: []cciptypes.ChainSelector{selectorA, selectorB},
			expectErr: true,
		},
		{
			desc:      "all selectors present in FValues",
			fVals:     map[cciptypes.ChainSelector]uint8{selectorA: 2, selectorB: 1},
			selectors: []cciptypes.ChainSelector{selectorA, selectorB},
			expectErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			args := &DistributedTopologyArgs{FValues: tc.fVals}
			err := args.validate(tc.selectors)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDistributedTopologyArgs_ChainToNodeMapping(t *testing.T) {
	selectorA := cciptypes.ChainSelector(2337)
	selectorB := cciptypes.ChainSelector(2338)
	selectorC := cciptypes.ChainSelector(2339)
	homeSelector := cciptypes.ChainSelector(1337)

	allNodes := make([][32]byte, 0, 21)
	for i := 0; i < 16; i++ {
		allNodes = append(allNodes, [32]byte{byte(i + 1)})
	}

	cases := []struct {
		desc           string
		fVals          map[cciptypes.ChainSelector]uint8
		nodeIDs        [][32]byte
		selectors      []cciptypes.ChainSelector
		home           cciptypes.ChainSelector
		expectedOutput map[cciptypes.ChainSelector][][32]byte
		expectErr      bool
	}{
		{
			desc:      "no nodes provided",
			fVals:     map[cciptypes.ChainSelector]uint8{selectorA: 1},
			nodeIDs:   [][32]byte{},
			selectors: []cciptypes.ChainSelector{selectorA},
			home:      homeSelector,
			expectErr: true,
		},
		{
			desc:      "missing fValue for selector",
			fVals:     map[cciptypes.ChainSelector]uint8{},
			nodeIDs:   allNodes,
			selectors: []cciptypes.ChainSelector{selectorA},
			home:      homeSelector,
			expectErr: true,
		},
		{
			desc:      "valid mapping with a minimal selectors and f values",
			fVals:     map[cciptypes.ChainSelector]uint8{selectorA: 1, selectorB: 2},
			nodeIDs:   allNodes,
			selectors: []cciptypes.ChainSelector{selectorA, selectorB},
			home:      homeSelector,
			expectErr: false,
			expectedOutput: map[cciptypes.ChainSelector][][32]byte{
				homeSelector: allNodes,
				selectorA:    {{1}, {2}, {3}, {4}},
				selectorB:    {{5}, {6}, {7}, {8}, {9}, {10}, {11}},
			},
		},
		{
			desc:      "valid mapping with wrap around",
			fVals:     map[cciptypes.ChainSelector]uint8{selectorA: 2, selectorB: 2, selectorC: 2},
			nodeIDs:   allNodes,
			selectors: []cciptypes.ChainSelector{selectorA, selectorB, selectorC},
			home:      homeSelector,
			expectErr: false,
			expectedOutput: map[cciptypes.ChainSelector][][32]byte{
				homeSelector: allNodes,
				selectorA:    {{1}, {2}, {3}, {4}, {5}, {6}, {7}},
				selectorB:    {{8}, {9}, {10}, {11}, {12}, {13}, {14}},
				selectorC:    {{15}, {16}, {1}, {2}, {3}, {4}, {5}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			args := &DistributedTopologyArgs{FValues: tc.fVals}
			mapping, err := args.ChainToNodeMapping(tc.nodeIDs, tc.selectors, tc.home)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, mapping)
				require.ElementsMatch(t, tc.expectedOutput[tc.home], mapping[tc.home])
				for _, selector := range tc.selectors {
					require.ElementsMatch(t, tc.expectedOutput[selector], mapping[selector],
						"unexpected mapping for selector %s", selector)
				}
			}
		})
	}
}
