package execute

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func Test_validateObserverReadingEligibility(t *testing.T) {
	tests := []struct {
		name         string
		observer     commontypes.OracleID
		observerCfg  map[commontypes.OracleID]cciptypes.ObserverInfo
		observedMsgs cciptypes.ExecutePluginMessageObservations
		expErr       string
	}{
		{
			name:     "ValidObserverAndMessages",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				1: {Reads: []cciptypes.ChainSelector{1, 2}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{
				1: {1: {}, 2: {}},
				2: {},
			},
		},
		{
			name:     "ObserverNotFound",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				2: {Reads: []cciptypes.ChainSelector{1, 2}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{
				1: {1: {}, 2: {}},
			},
			expErr: "observer not found in config",
		},
		{
			name:     "ObserverNotAllowedToReadChain",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				1: {Reads: []cciptypes.ChainSelector{1}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{
				2: {1: {}},
			},
			expErr: "observer not allowed to read from chain 2",
		},
		{
			name:     "NoMessagesObserved",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				1: {Reads: []cciptypes.ChainSelector{1, 2}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{},
		},
		{
			name:     "EmptyMessagesInChain",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				1: {Reads: []cciptypes.ChainSelector{1, 2}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{
				1: {},
				2: {1: {}, 2: {}},
			},
		},
		{
			name:     "AllMessagesEmpty",
			observer: commontypes.OracleID(1),
			observerCfg: map[commontypes.OracleID]cciptypes.ObserverInfo{
				1: {Reads: []cciptypes.ChainSelector{1, 2}},
			},
			observedMsgs: cciptypes.ExecutePluginMessageObservations{
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
		observedData map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages
		expErr       bool
	}{
		{
			name: "ValidData",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{
				1: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.SeqNumRange{1, 10},
							ExecutedMessages:    []cciptypes.SeqNum{1, 2, 3},
						},
					},
				},
				2: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{2},
							SequenceNumberRange: cciptypes.SeqNumRange{11, 20},
							ExecutedMessages:    []cciptypes.SeqNum{11, 12, 13},
						},
					},
				},
			},
		},
		{
			name: "DuplicateMerkleRoot",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{
				1: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.SeqNumRange{1, 10},
							ExecutedMessages:    []cciptypes.SeqNum{1, 2, 3},
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.SeqNumRange{11, 20},
							ExecutedMessages:    []cciptypes.SeqNum{11, 12, 13},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "OverlappingSequenceNumberRange",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{
				1: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.SeqNumRange{1, 10},
							ExecutedMessages:    []cciptypes.SeqNum{1, 2, 3},
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{2},
							SequenceNumberRange: cciptypes.SeqNumRange{5, 15},
							ExecutedMessages:    []cciptypes.SeqNum{6, 7, 8},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "ExecutedMessageOutsideObservedRange",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{
				1: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.SeqNumRange{1, 10},
							ExecutedMessages:    []cciptypes.SeqNum{1, 2, 11},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "NoCommitData",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{
				1: {},
			},
		},
		{
			name:         "EmptyObservedData",
			observedData: map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages{},
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
		reports []cciptypes.ExecutePluginCommitDataWithMessages
	}

	tests := []struct {
		name string
		args args
		want []cciptypes.SeqNumRange
		err  error
	}{
		{
			name: "empty",
			args: args{reports: []cciptypes.ExecutePluginCommitDataWithMessages{}},
			want: nil,
		},
		{
			name: "overlapping ranges",
			args: args{reports: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(15, 25),
					},
				},
			},
			},
			err: errOverlappingRanges,
		},
		{
			name: "simple ranges collapsed",
			args: args{reports: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(21, 40),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(41, 60),
					},
				},
			},
			},
			want: []cciptypes.SeqNumRange{{10, 60}},
		},
		{
			name: "non-contiguous ranges",
			args: args{reports: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60)},
				},
			},
			},
			want: []cciptypes.SeqNumRange{{10, 20}, {30, 40}, {50, 60}},
		},
		{
			name: "contiguous and non-contiguous ranges",
			args: args{reports: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(21, 40),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
					},
				},
			},
			},
			want: []cciptypes.SeqNumRange{{10, 40}, {50, 60}},
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
		reports []cciptypes.CommitPluginReportWithMeta
	}
	tests := []struct {
		name string
		args args
		want cciptypes.ExecutePluginCommitObservations
	}{
		{
			name: "empty",
			args: args{reports: []cciptypes.CommitPluginReportWithMeta{}},
			want: cciptypes.ExecutePluginCommitObservations{},
		},
		{
			name: "reports",
			args: args{reports: []cciptypes.CommitPluginReportWithMeta{{
				Report: cciptypes.CommitPluginReport{
					MerkleRoots: []cciptypes.MerkleRootChain{
						{ChainSel: 1, SeqNumsRange: cciptypes.NewSeqNumRange(10, 20), MerkleRoot: cciptypes.Bytes32{1}},
						{ChainSel: 2, SeqNumsRange: cciptypes.NewSeqNumRange(30, 40), MerkleRoot: cciptypes.Bytes32{2}},
					}}}}},
			want: cciptypes.ExecutePluginCommitObservations{
				1: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SourceChain:         1,
							MerkleRoot:          cciptypes.Bytes32{1},
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
							ExecutedMessages:    nil,
						},
					},
				},
				2: {
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SourceChain:         2,
							MerkleRoot:          cciptypes.Bytes32{2},
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
							ExecutedMessages:    nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, groupByChainSelector(tt.args.reports), "groupByChainSelector(%v)", tt.args.reports)
		})
	}
}

func Test_filterOutFullyExecutedMessages(t *testing.T) {
	type args struct {
		reports          []cciptypes.ExecutePluginCommitDataWithMessages
		executedMessages []cciptypes.SeqNumRange
	}
	tests := []struct {
		name    string
		args    args
		want    []cciptypes.ExecutePluginCommitDataWithMessages
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
				reports:          []cciptypes.ExecutePluginCommitDataWithMessages{},
				executedMessages: nil,
			},
			want:    []cciptypes.ExecutePluginCommitDataWithMessages{},
			wantErr: assert.NoError,
		},
		{
			name: "no executed messages",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
				},
				executedMessages: nil,
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20)}},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40)}},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60)}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "executed messages",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20)}},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40)}},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60)}},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(0, 100),
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "2 partially executed",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20)},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40)},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60)},
					},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(15, 35),
				},
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						ExecutedMessages:    []cciptypes.SeqNum{15, 16, 17, 18, 19, 20},
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						ExecutedMessages:    []cciptypes.SeqNum{30, 31, 32, 33, 34, 35},
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "2 partially executed 1 fully executed",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(15, 55),
				},
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						ExecutedMessages:    []cciptypes.SeqNum{15, 16, 17, 18, 19, 20},
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						ExecutedMessages:    []cciptypes.SeqNum{50, 51, 52, 53, 54, 55},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "first report executed",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(10, 20),
				},
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "last report executed",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(50, 60),
				},
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "sort-report",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
				},
				executedMessages: nil,
			},
			want: []cciptypes.ExecutePluginCommitDataWithMessages{
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
					},
				},
				{
					ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
						SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "sort-executed",
			args: args{
				reports: []cciptypes.ExecutePluginCommitDataWithMessages{
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(10, 20),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(30, 40),
						},
					},
					{
						ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{
							SequenceNumberRange: cciptypes.NewSeqNumRange(50, 60),
						},
					},
				},
				executedMessages: []cciptypes.SeqNumRange{
					cciptypes.NewSeqNumRange(50, 60),
					cciptypes.NewSeqNumRange(10, 20),
					cciptypes.NewSeqNumRange(30, 40),
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

func Test_decodeAttributedObservations(t *testing.T) {
	mustEncode := func(obs cciptypes.ExecutePluginObservation) []byte {
		enc, err := obs.Encode()
		if err != nil {
			t.Fatal("Unable to encode")
		}
		return enc
	}
	tests := []struct {
		name    string
		args    []types.AttributedObservation
		want    []decodedAttributedObservation
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "empty",
			args:    nil,
			want:    []decodedAttributedObservation{},
			wantErr: assert.NoError,
		},
		{
			name: "one observation",
			args: []types.AttributedObservation{
				{
					Observer: commontypes.OracleID(1),
					Observation: mustEncode(cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							1: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{1}}}},
						},
					}),
				},
			},
			want: []decodedAttributedObservation{
				{
					Observer: commontypes.OracleID(1),
					Observation: cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							1: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{1}}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple observations",
			args: []types.AttributedObservation{
				{
					Observer: commontypes.OracleID(1),
					Observation: mustEncode(cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							1: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{1}}}},
						},
					}),
				},
				{
					Observer: commontypes.OracleID(2),
					Observation: mustEncode(cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							2: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{2}}}},
						},
					}),
				},
			},
			want: []decodedAttributedObservation{
				{
					Observer: commontypes.OracleID(1),
					Observation: cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							1: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{1}}}},
						},
					},
				},
				{
					Observer: commontypes.OracleID(2),
					Observation: cciptypes.ExecutePluginObservation{
						CommitReports: cciptypes.ExecutePluginCommitObservations{
							2: {{ExecutePluginCommitData: cciptypes.ExecutePluginCommitData{MerkleRoot: cciptypes.Bytes32{2}}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid observation",
			args: []types.AttributedObservation{
				{
					Observer:    commontypes.OracleID(1),
					Observation: []byte("invalid"),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeAttributedObservations(tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("decodeAttributedObservations(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "decodeAttributedObservations(%v)", tt.args)
		})
	}
}
