package commit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/ccipocr3/internal/mocks"
	"github.com/smartcontractkit/ccipocr3/internal/model"
)

func Test_getPendingExecutedReports(t *testing.T) {
	tests := []struct {
		name    string
		reports []model.CommitPluginReportWithMeta
		ranges  map[model.ChainSelector][]model.SeqNumRange
		want    model.ExecutePluginCommitObservations
		want1   time.Time
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "empty",
			reports: nil,
			ranges:  nil,
			want:    model.ExecutePluginCommitObservations{},
			want1:   time.Time{},
			wantErr: assert.NoError,
		},
		{
			name: "single non-executed report",
			reports: []model.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report: model.CommitPluginReport{
						MerkleRoots: []model.MerkleRootChain{
							{
								ChainSel:     1,
								SeqNumsRange: model.NewSeqNumRange(1, 10),
							},
						},
					},
				},
			},
			ranges: map[model.ChainSelector][]model.SeqNumRange{
				1: nil,
			},
			want: model.ExecutePluginCommitObservations{
				1: []model.ExecutePluginCommitData{
					{
						SequenceNumberRange: model.NewSeqNumRange(1, 10),
						ExecutedMessages:    nil,
						Timestamp:           time.UnixMilli(10101010101),
						BlockNum:            999,
					},
				},
			},
			want1:   time.UnixMilli(10101010101),
			wantErr: assert.NoError,
		},
		{
			name: "single half-executed report",
			reports: []model.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report: model.CommitPluginReport{
						MerkleRoots: []model.MerkleRootChain{
							{
								ChainSel:     1,
								SeqNumsRange: model.NewSeqNumRange(1, 10),
							},
						},
					},
				},
			},
			ranges: map[model.ChainSelector][]model.SeqNumRange{
				1: {
					model.NewSeqNumRange(1, 3),
					model.NewSeqNumRange(7, 8),
				},
			},
			want: model.ExecutePluginCommitObservations{
				1: []model.ExecutePluginCommitData{
					{
						SequenceNumberRange: model.NewSeqNumRange(1, 10),
						Timestamp:           time.UnixMilli(10101010101),
						BlockNum:            999,
						ExecutedMessages:    []model.SeqNum{1, 2, 3, 7, 8},
					},
				},
			},
			want1:   time.UnixMilli(10101010101),
			wantErr: assert.NoError,
		},
		{
			name: "last timestamp",
			reports: []model.CommitPluginReportWithMeta{
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(10101010101),
					Report:    model.CommitPluginReport{},
				},
				{
					BlockNum:  999,
					Timestamp: time.UnixMilli(9999999999999999),
					Report:    model.CommitPluginReport{},
				},
			},
			ranges:  map[model.ChainSelector][]model.SeqNumRange{},
			want:    model.ExecutePluginCommitObservations{},
			want1:   time.UnixMilli(9999999999999999),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockReader := mocks.NewCCIPReader()
			mockReader.On("CommitReportsGTETimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.reports, nil)
			for k, v := range tt.ranges {
				mockReader.On("ExecutedMessageRanges", mock.Anything, k, mock.Anything, mock.Anything).Return(v, nil)
			}

			// CCIP Reader mocks:
			// once:
			//      CommitReportsGTETimestamp(ctx, dest, ts, 1000) -> ([]model.CommitPluginReportWithMeta, error)
			// for each chain selector:
			//      ExecutedMessageRanges(ctx, selector, dest, seqRange) -> ([]model.SeqNumRange, error)

			got, got1, err := getPendingExecutedReports(context.Background(), mockReader, 123, time.Now())
			if !tt.wantErr(t, err, "getPendingExecutedReports(...)") {
				return
			}
			assert.Equalf(t, tt.want, got, "getPendingExecutedReports(...)")
			assert.Equalf(t, tt.want1, got1, "getPendingExecutedReports(...)")
		})
	}
}
