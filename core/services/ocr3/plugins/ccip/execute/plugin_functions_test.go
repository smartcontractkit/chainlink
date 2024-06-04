package commit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/libocr/commontypes"
)

func Test_validateObserverReadingEligibility(t *testing.T) {
	tests := []struct {
		name         string
		observer     commontypes.OracleID
		observerCfg  map[commontypes.OracleID]model.ObserverInfo
		observedMsgs model.ExecutePluginMessageObservations
		expErr       string
	}{
		{
			name:     "ValidObserverAndMessages",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				1: {Reads: []model.ChainSelector{1, 2}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{
				1: {1: {}, 2: {}},
				2: {},
			},
		},
		{
			name:     "ObserverNotFound",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				2: {Reads: []model.ChainSelector{1, 2}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{
				1: {1: {}, 2: {}},
			},
			expErr: "observer not found in config",
		},
		{
			name:     "ObserverNotAllowedToReadChain",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				1: {Reads: []model.ChainSelector{1}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{
				2: {1: {}},
			},
			expErr: "observer not allowed to read from chain 2",
		},
		{
			name:     "NoMessagesObserved",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				1: {Reads: []model.ChainSelector{1, 2}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{},
		},
		{
			name:     "EmptyMessagesInChain",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				1: {Reads: []model.ChainSelector{1, 2}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{
				1: {},
				2: {1: {}, 2: {}},
			},
		},
		{
			name:     "AllMessagesEmpty",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]model.ObserverInfo{
				1: {Reads: []model.ChainSelector{1, 2}},
			},
			observedMsgs: model.ExecutePluginMessageObservations{
				1: {},
				2: {},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObserverReadingEligibility(tc.observer, tc.observerCfg, tc.observedMsgs)
			if len(tc.expErr) != 0 {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_validateObservedSequenceNumbers(t *testing.T) {
	testCases := []struct {
		name         string
		observedData map[model.ChainSelector][]model.ExecutePluginCommitData
		expErr       bool
	}{
		{
			name: "ValidData",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{
				1: {
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.SeqNumRange{1, 10},
						ExecutedMessages:    []model.SeqNum{1, 2, 3},
					},
				},
				2: {
					{
						MerkleRoot:          model.Bytes32{2},
						SequenceNumberRange: model.SeqNumRange{11, 20},
						ExecutedMessages:    []model.SeqNum{11, 12, 13},
					},
				},
			},
		},
		{
			name: "DuplicateMerkleRoot",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{
				1: {
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.SeqNumRange{1, 10},
						ExecutedMessages:    []model.SeqNum{1, 2, 3},
					},
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.SeqNumRange{11, 20},
						ExecutedMessages:    []model.SeqNum{11, 12, 13},
					},
				},
			},
			expErr: true,
		},
		{
			name: "OverlappingSequenceNumberRange",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{
				1: {
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.SeqNumRange{1, 10},
						ExecutedMessages:    []model.SeqNum{1, 2, 3},
					},
					{
						MerkleRoot:          model.Bytes32{2},
						SequenceNumberRange: model.SeqNumRange{5, 15},
						ExecutedMessages:    []model.SeqNum{6, 7, 8},
					},
				},
			},
			expErr: true,
		},
		{
			name: "ExecutedMessageOutsideObservedRange",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{
				1: {
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.SeqNumRange{1, 10},
						ExecutedMessages:    []model.SeqNum{1, 2, 11},
					},
				},
			},
			expErr: true,
		},
		{
			name: "NoCommitData",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{
				1: {},
			},
		},
		{
			name:         "EmptyObservedData",
			observedData: map[model.ChainSelector][]model.ExecutePluginCommitData{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObservedSequenceNumbers(tc.observedData)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_computeRanges(t *testing.T) {
	type args struct {
		reports []model.ExecutePluginCommitData
	}

	tests := []struct {
		name string
		args args
		want []model.SeqNumRange
		err  error
	}{
		{
			name: "empty",
			args: args{reports: []model.ExecutePluginCommitData{}},
			want: nil,
		},
		{
			name: "overlapping ranges",
			args: args{reports: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(15, 25)}},
			},
			err: errOverlappingRanges,
		},
		{
			name: "simple ranges collapsed",
			args: args{reports: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(21, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(41, 60)}},
			},
			want: []model.SeqNumRange{{10, 60}},
		},
		{
			name: "non-contiguous ranges",
			args: args{reports: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)}},
			},
			want: []model.SeqNumRange{{10, 20}, {30, 40}, {50, 60}},
		},
		{
			name: "contiguous and non-contiguous ranges",
			args: args{reports: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(21, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)}},
			},
			want: []model.SeqNumRange{{10, 40}, {50, 60}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := computeRanges(tt.args.reports)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_groupByChainSelector(t *testing.T) {
	type args struct {
		reports []model.CommitPluginReportWithMeta
	}
	tests := []struct {
		name string
		args args
		want model.ExecutePluginCommitObservations
	}{
		{
			name: "empty",
			args: args{reports: []model.CommitPluginReportWithMeta{}},
			want: model.ExecutePluginCommitObservations{},
		},
		{
			name: "reports",
			args: args{reports: []model.CommitPluginReportWithMeta{{
				Report: model.CommitPluginReport{
					MerkleRoots: []model.MerkleRootChain{
						{ChainSel: 1, SeqNumsRange: model.NewSeqNumRange(10, 20), MerkleRoot: model.Bytes32{1}},
						{ChainSel: 2, SeqNumsRange: model.NewSeqNumRange(30, 40), MerkleRoot: model.Bytes32{2}},
					}}}}},
			want: model.ExecutePluginCommitObservations{
				1: {
					{
						MerkleRoot:          model.Bytes32{1},
						SequenceNumberRange: model.NewSeqNumRange(10, 20),
						ExecutedMessages:    nil,
					},
				},
				2: {
					{

						MerkleRoot:          model.Bytes32{2},
						SequenceNumberRange: model.NewSeqNumRange(30, 40),
						ExecutedMessages:    nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, groupByChainSelector(tt.args.reports), "groupByChainSelector(%v)", tt.args.reports)
		})
	}
}

func Test_filterOutFullyExecutedMessages(t *testing.T) {
	type args struct {
		reports          []model.ExecutePluginCommitData
		executedMessages []model.SeqNumRange
	}
	tests := []struct {
		name    string
		args    args
		want    []model.ExecutePluginCommitData
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty",
			args: args{
				reports:          nil,
				executedMessages: nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "empty2",
			args: args{
				reports:          []model.ExecutePluginCommitData{},
				executedMessages: nil,
			},
			want:    []model.ExecutePluginCommitData{},
			wantErr: assert.NoError,
		},
		{
			name: "no executed messages",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: nil,
			},
			want: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "executed messages",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(0, 100),
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "2 partially executed",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(15, 35),
				},
			},
			want: []model.ExecutePluginCommitData{
				{
					SequenceNumberRange: model.NewSeqNumRange(10, 20),
					ExecutedMessages:    []model.SeqNum{15, 16, 17, 18, 19, 20},
				},
				{
					SequenceNumberRange: model.NewSeqNumRange(30, 40),
					ExecutedMessages:    []model.SeqNum{30, 31, 32, 33, 34, 35},
				},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "2 partially executed 1 fully executed",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(15, 55),
				},
			},
			want: []model.ExecutePluginCommitData{
				{
					SequenceNumberRange: model.NewSeqNumRange(10, 20),
					ExecutedMessages:    []model.SeqNum{15, 16, 17, 18, 19, 20},
				},
				{
					SequenceNumberRange: model.NewSeqNumRange(50, 60),
					ExecutedMessages:    []model.SeqNum{50, 51, 52, 53, 54, 55},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "first report executed",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(10, 20),
				},
			},
			want: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "last report executed",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(50, 60),
				},
			},
			want: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "sort-report",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				},
				executedMessages: nil,
			},
			want: []model.ExecutePluginCommitData{
				{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
				{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
				{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "sort-executed",
			args: args{
				reports: []model.ExecutePluginCommitData{
					{SequenceNumberRange: model.NewSeqNumRange(10, 20)},
					{SequenceNumberRange: model.NewSeqNumRange(30, 40)},
					{SequenceNumberRange: model.NewSeqNumRange(50, 60)},
				},
				executedMessages: []model.SeqNumRange{
					model.NewSeqNumRange(50, 60),
					model.NewSeqNumRange(10, 20),
					model.NewSeqNumRange(30, 40),
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterOutExecutedMessages(tt.args.reports, tt.args.executedMessages)
			if !tt.wantErr(t, err, fmt.Sprintf("filterOutExecutedMessages(%v, %v)", tt.args.reports, tt.args.executedMessages)) {
				return
			}
			assert.Equalf(t, tt.want, got, "filterOutExecutedMessages(%v, %v)", tt.args.reports, tt.args.executedMessages)
		})
	}
}
