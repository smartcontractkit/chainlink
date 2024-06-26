package execute

import (
	"context"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/ccipocr3/internal/mocks"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestSomethingCool(t *testing.T) {

	foo := map[cciptypes.ChainSelector]int{
		cciptypes.ChainSelector(1):              1,
		cciptypes.ChainSelector(math.MaxUint64): 1,
	}

	js, _ := json.Marshal(foo)
	t.Log(string(js))

	b := []byte(`{"1":1,"18446744073709551615":1}`)
	var bar map[cciptypes.ChainSelector]int
	assert.NoError(t, json.Unmarshal(b, &bar))
	t.Log(bar)
}

func Test_getPendingExecutedReports(t *testing.T) {
	tests := []struct {
		name    string
		reports []cciptypes.CommitPluginReportWithMeta
		ranges  map[cciptypes.ChainSelector][]cciptypes.SeqNumRange
		want    cciptypes.ExecutePluginCommitObservations
		want1   time.Time
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "empty",
			reports: nil,
			ranges:  nil,
			want:    cciptypes.ExecutePluginCommitObservations{},
			want1:   time.Time{},
			wantErr: assert.NoError,
		},
		{
			name: "single non-executed report",
			reports: []cciptypes.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report: cciptypes.CommitPluginReport{
						MerkleRoots: []cciptypes.MerkleRootChain{
							{
								ChainSel:     1,
								SeqNumsRange: cciptypes.NewSeqNumRange(1, 10),
							},
						},
					},
				},
			},
			ranges: map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				1: nil,
			},
			want: cciptypes.ExecutePluginCommitObservations{
				1: []cciptypes.ExecutePluginCommitDataWithMessages{
					{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SourceChain:         1,
						SequenceNumberRange: cciptypes.NewSeqNumRange(1, 10),
						ExecutedMessages:    nil,
						Timestamp:           time.UnixMilli(10101010101),
						BlockNum:            999,
					}},
				},
			},
			want1:   time.UnixMilli(10101010101),
			wantErr: assert.NoError,
		},
		{
			name: "single half-executed report",
			reports: []cciptypes.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report: cciptypes.CommitPluginReport{
						MerkleRoots: []cciptypes.MerkleRootChain{
							{
								ChainSel:     1,
								SeqNumsRange: cciptypes.NewSeqNumRange(1, 10),
							},
						},
					},
				},
			},
			ranges: map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				1: {
					cciptypes.NewSeqNumRange(1, 3),
					cciptypes.NewSeqNumRange(7, 8),
				},
			},
			want: cciptypes.ExecutePluginCommitObservations{
				1: []cciptypes.ExecutePluginCommitDataWithMessages{
					{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SourceChain:         1,
						SequenceNumberRange: cciptypes.NewSeqNumRange(1, 10),
						Timestamp:           time.UnixMilli(10101010101),
						BlockNum:            999,
						ExecutedMessages:    []cciptypes.SeqNum{1, 2, 3, 7, 8},
					}},
				},
			},
			want1:   time.UnixMilli(10101010101),
			wantErr: assert.NoError,
		},
		{
			name: "last timestamp",
			reports: []cciptypes.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report:    cciptypes.CommitPluginReport{},
				},
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(9999999999999999),
					Report:    cciptypes.CommitPluginReport{},
				},
			},
			ranges:  map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{},
			want:    cciptypes.ExecutePluginCommitObservations{},
			want1:   time.UnixMilli(9999999999999999),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockReader := mocks.NewCCIPReader()
			mockReader.On(
				"CommitReportsGTETimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			).Return(tt.reports, nil)
			for k, v := range tt.ranges {
				mockReader.On("ExecutedMessageRanges", mock.Anything, k, mock.Anything, mock.Anything).Return(v, nil)
			}

			// CCIP Reader mocks:
			// once:
			//      CommitReportsGTETimestamp(ctx, dest, ts, 1000) -> ([]cciptypes.CommitPluginReportWithMeta, error)
			// for each chain selector:
			//      ExecutedMessageRanges(ctx, selector, dest, seqRange) -> ([]cciptypes.SeqNumRange, error)

			got, got1, err := getPendingExecutedReports(context.Background(), mockReader, 123, time.Now())
			if !tt.wantErr(t, err, "getPendingExecutedReports(...)") {
				return
			}
			assert.Equalf(t, tt.want, got, "getPendingExecutedReports(...)")
			assert.Equalf(t, tt.want1, got1, "getPendingExecutedReports(...)")
		})
	}
}
