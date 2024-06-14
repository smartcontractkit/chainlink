package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func Test_CommitReportValidator_ExecutePluginCommitData(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		reports []cciptypes.ExecutePluginCommitData
		valid   []cciptypes.ExecutePluginCommitData
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			valid:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "single report, enough observations",
			min:  1,
			reports: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}},
			},
			valid: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "single report, not enough observations",
			min:  2,
			reports: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}},
			},
			valid:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "multiple reports, partial observations",
			min:  2,
			reports: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{3}},
				{MerkleRoot: [32]byte{1}},
				{MerkleRoot: [32]byte{2}},
				{MerkleRoot: [32]byte{1}},
				{MerkleRoot: [32]byte{2}},
			},
			valid: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}},
				{MerkleRoot: [32]byte{2}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple reports for same root",
			min:  2,
			reports: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}, BlockNum: 1},
				{MerkleRoot: [32]byte{1}, BlockNum: 2},
				{MerkleRoot: [32]byte{1}, BlockNum: 3},
				{MerkleRoot: [32]byte{1}, BlockNum: 4},
				{MerkleRoot: [32]byte{1}, BlockNum: 1},
			},
			valid: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}, BlockNum: 1},
			},
			wantErr: assert.NoError,
		},
		{
			name: "different executed messages same root",
			min:  2,
			reports: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{1, 2}},
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{2, 3}},
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{3, 4}},
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{4, 5}},
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{5, 6}},
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{1, 2}},
			},
			valid: []cciptypes.ExecutePluginCommitData{
				{MerkleRoot: [32]byte{1}, ExecutedMessages: []cciptypes.SeqNum{1, 2}},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Initialize the minObservationValidator
			idFunc := func(data cciptypes.ExecutePluginCommitData) [32]byte {
				return sha3.Sum256([]byte(fmt.Sprintf("%v", data)))
			}
			validator := NewMinObservationValidator[cciptypes.ExecutePluginCommitData](tt.min, idFunc)
			for _, report := range tt.reports {
				err := validator.Add(report)
				require.NoError(t, err)
			}

			// Test the results
			got, err := validator.GetValid()
			if !tt.wantErr(t, err, "GetValid()") {
				return
			}
			if !assert.ElementsMatch(t, got, tt.valid) {
				t.Errorf("GetValid() = %v, valid %v", got, tt.valid)
			}
		})
	}
}

func Test_CommitReportValidator_Generics(t *testing.T) {
	type Generic struct {
		number int
	}

	// Initialize the minObservationValidator
	idFunc := func(data Generic) [32]byte {
		return sha3.Sum256([]byte(fmt.Sprintf("%v", data)))
	}
	validator := NewMinObservationValidator[Generic](2, idFunc)

	wantValue := Generic{number: 1}
	otherValue := Generic{number: 2}

	err := validator.Add(wantValue)
	require.NoError(t, err)
	err = validator.Add(wantValue)
	require.NoError(t, err)
	err = validator.Add(otherValue)
	require.NoError(t, err)

	// Test the results

	wantValid := []Generic{wantValue}
	got, err := validator.GetValid()
	require.NoError(t, err)
	if !assert.ElementsMatch(t, got, wantValid) {
		t.Errorf("GetValid() = %v, valid %v", got, wantValid)
	}
}
