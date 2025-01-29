package oraclecreator

import (
	"bytes"
	"crypto/sha256"
	"sort"
	"testing"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "need to create candidate",
			currentDigests:  []cciptypes.Bytes32{{1}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{2},
			expectedActions: []syncAction{
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "both configs empty",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{},
			candidateDigest: cciptypes.Bytes32{},
			expectedActions: []syncAction{
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{1}},
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "replace both configs",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{3},
			candidateDigest: cciptypes.Bytes32{4},
			expectedActions: []syncAction{
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{1}},
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{3}},
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{4}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commitConfigDigest := cciptypes.Bytes32{1}

			currentDigests := make([]cciptypes.Bytes32, len(tt.currentDigests))
			for i, d := range tt.currentDigests {
				currentDigests[i] = writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
					sha256.Sum256(append(commitConfigDigest[:], d[:]...)))
			}

			actions := calculateSyncActions(
				commitConfigDigest,
				currentDigests,
				tt.activeDigest,
				tt.candidateDigest,
			)

			require.Equal(t, len(tt.expectedActions), len(actions))

			// Sort both slices to ensure consistent comparison
			sort.Slice(actions, func(i, j int) bool {
				if actions[i].actionType != actions[j].actionType {
					return actions[i].actionType < actions[j].actionType
				}
				return bytes.Compare(actions[i].endpointConfigDigest[:], actions[j].endpointConfigDigest[:]) < 0
			})
			sort.Slice(tt.expectedActions, func(i, j int) bool {
				if tt.expectedActions[i].actionType != tt.expectedActions[j].actionType {
					return tt.expectedActions[i].actionType < tt.expectedActions[j].actionType
				}
				return bytes.Compare(tt.expectedActions[i].endpointConfigDigest[:], tt.expectedActions[j].endpointConfigDigest[:]) < 0
			})

			for i := range actions {
				require.Equal(t, tt.expectedActions[i].actionType, actions[i].actionType)

				expEndpointConfigDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
					sha256.Sum256(append(commitConfigDigest[:], tt.expectedActions[i].endpointConfigDigest[:]...)))

				require.Equal(t, expEndpointConfigDigest, actions[i].endpointConfigDigest)
			}
		})
	}
}
