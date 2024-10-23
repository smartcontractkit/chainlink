package oraclecreator

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func TestCalculateSyncActions(t *testing.T) {
	tests := []struct {
		name            string
		currentDigests  []cciptypes.Bytes32
		activeDigest    cciptypes.Bytes32
		candidateDigest cciptypes.Bytes32
		expectedActions []syncAction
	}{
		{
			name:            "no changes needed",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{2},
			expectedActions: nil,
		},
		{
			name:            "need to close candidate",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{}, // empty
			expectedActions: []syncAction{
				{Type: ActionClose, configDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "need to create candidate",
			currentDigests:  []cciptypes.Bytes32{{1}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{2},
			expectedActions: []syncAction{
				{Type: ActionCreate, configDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "both configs empty",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{},
			candidateDigest: cciptypes.Bytes32{},
			expectedActions: []syncAction{
				{Type: ActionClose, configDigest: cciptypes.Bytes32{1}},
				{Type: ActionClose, configDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "replace both configs",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{3},
			candidateDigest: cciptypes.Bytes32{4},
			expectedActions: []syncAction{
				{Type: ActionClose, configDigest: cciptypes.Bytes32{1}},
				{Type: ActionClose, configDigest: cciptypes.Bytes32{2}},
				{Type: ActionCreate, configDigest: cciptypes.Bytes32{3}},
				{Type: ActionCreate, configDigest: cciptypes.Bytes32{4}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions := calculateSyncActions(
				tt.currentDigests,
				tt.activeDigest,
				tt.candidateDigest,
			)

			require.Equal(t, len(tt.expectedActions), len(actions))

			// Sort both slices to ensure consistent comparison
			sort.Slice(actions, func(i, j int) bool {
				if actions[i].Type != actions[j].Type {
					return actions[i].Type < actions[j].Type
				}
				return bytes.Compare(actions[i].configDigest[:], actions[j].configDigest[:]) < 0
			})
			sort.Slice(tt.expectedActions, func(i, j int) bool {
				if tt.expectedActions[i].Type != tt.expectedActions[j].Type {
					return tt.expectedActions[i].Type < tt.expectedActions[j].Type
				}
				return bytes.Compare(tt.expectedActions[i].configDigest[:], tt.expectedActions[j].configDigest[:]) < 0
			})

			for i := range actions {
				require.Equal(t, tt.expectedActions[i].Type, actions[i].Type)
				require.Equal(t, tt.expectedActions[i].configDigest, actions[i].configDigest)
			}
		})
	}
}
