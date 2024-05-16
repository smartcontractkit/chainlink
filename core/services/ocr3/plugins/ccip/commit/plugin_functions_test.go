package commit

import (
	"context"
	"math/big"
	"reflect"
	"slices"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/mocks"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_observeMaxSeqNumsPerChain(t *testing.T) {
	testCases := []struct {
		name           string
		prevOutcome    model.CommitPluginOutcome
		onChainSeqNums map[model.ChainSelector]model.SeqNum
		readChains     []model.ChainSelector
		destChain      model.ChainSelector
		expErr         bool
		expMaxSeqNums  []model.SeqNumChain
	}{
		{
			name:        "report on chain seq num when no previous outcome and can read dest",
			prevOutcome: model.CommitPluginOutcome{},
			onChainSeqNums: map[model.ChainSelector]model.SeqNum{
				1: 10,
				2: 20,
			},
			readChains: []model.ChainSelector{1, 2, 3},
			destChain:  3,
			expErr:     false,
			expMaxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
		},
		{
			name:        "nothing to report when there is no previous outcome and cannot read dest",
			prevOutcome: model.CommitPluginOutcome{},
			onChainSeqNums: map[model.ChainSelector]model.SeqNum{
				1: 10,
				2: 20,
			},
			readChains:    []model.ChainSelector{1, 2},
			destChain:     3,
			expErr:        false,
			expMaxSeqNums: []model.SeqNumChain{},
		},
		{
			name: "report previous outcome seq nums and override when on chain is higher if can read dest",
			prevOutcome: model.CommitPluginOutcome{
				MaxSeqNums: []model.SeqNumChain{
					{ChainSel: 1, SeqNum: 11}, // for chain 1 previous outcome is higher than on-chain state
					{ChainSel: 2, SeqNum: 19}, // for chain 2 previous outcome is behind on-chain state
				},
			},
			onChainSeqNums: map[model.ChainSelector]model.SeqNum{
				1: 10,
				2: 20,
			},
			readChains: []model.ChainSelector{1, 2, 3},
			destChain:  3,
			expErr:     false,
			expMaxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 11},
				{ChainSel: 2, SeqNum: 20},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := mocks.NewCCIPReader()
			knownSourceChains := slicelib.Filter(tc.readChains, func(ch model.ChainSelector) bool { return ch != tc.destChain })
			lggr := logger.Test(t)

			var encodedPrevOutcome []byte
			var err error
			if !reflect.DeepEqual(tc.prevOutcome, model.CommitPluginOutcome{}) {
				encodedPrevOutcome, err = tc.prevOutcome.Encode()
				assert.NoError(t, err)
			}

			onChainSeqNums := make([]model.SeqNum, 0)
			for _, chain := range knownSourceChains {
				if v, ok := tc.onChainSeqNums[chain]; !ok {
					t.Fatalf("invalid test case missing on chain seq num expectation for %d", chain)
				} else {
					onChainSeqNums = append(onChainSeqNums, v)
				}
			}
			mockReader.On("NextSeqNum", ctx, knownSourceChains).Return(onChainSeqNums, nil)

			seqNums, err := observeMaxSeqNums(
				ctx,
				lggr,
				mockReader,
				encodedPrevOutcome,
				mapset.NewSet(tc.readChains...),
				tc.destChain,
				knownSourceChains,
			)

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expMaxSeqNums, seqNums)
		})
	}
}

func Test_observeNewMsgs(t *testing.T) {
	testCases := []struct {
		name               string
		maxSeqNumsPerChain []model.SeqNumChain
		readChains         []model.ChainSelector
		destChain          model.ChainSelector
		msgScanBatchSize   int
		newMsgs            map[model.ChainSelector][]model.CCIPMsg
		expMsgs            []model.CCIPMsgBaseDetails
		expErr             bool
	}{
		{
			name: "no new messages",
			maxSeqNumsPerChain: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []model.ChainSelector{1, 2},
			msgScanBatchSize: 256,
			newMsgs: map[model.ChainSelector][]model.CCIPMsg{
				1: {},
				2: {},
			},
			expMsgs: []model.CCIPMsgBaseDetails{},
			expErr:  false,
		},
		{
			name: "new messages",
			maxSeqNumsPerChain: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []model.ChainSelector{1, 2},
			msgScanBatchSize: 256,
			newMsgs: map[model.ChainSelector][]model.CCIPMsg{
				1: {
					{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}},
				},
				2: {
					{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
					{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
				},
			},
			expMsgs: []model.CCIPMsgBaseDetails{
				{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11},
				{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21},
				{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22},
			},
			expErr: false,
		},
		{
			name: "new messages but one chain is not readable",
			maxSeqNumsPerChain: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []model.ChainSelector{2},
			msgScanBatchSize: 256,
			newMsgs: map[model.ChainSelector][]model.CCIPMsg{
				2: {
					{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
					{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
				},
			},
			expMsgs: []model.CCIPMsgBaseDetails{
				{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21},
				{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := mocks.NewCCIPReader()
			lggr := logger.Test(t)

			for _, seqNumChain := range tc.maxSeqNumsPerChain {
				if slices.Contains(tc.readChains, seqNumChain.ChainSel) {
					mockReader.On(
						"MsgsBetweenSeqNums",
						ctx,
						[]model.ChainSelector{seqNumChain.ChainSel},
						model.NewSeqNumRange(seqNumChain.SeqNum+1, seqNumChain.SeqNum+model.SeqNum(1+tc.msgScanBatchSize)),
					).Return(tc.newMsgs[seqNumChain.ChainSel], nil)
				}
			}

			msgs, err := observeNewMsgs(
				ctx,
				lggr,
				mockReader,
				mapset.NewSet(tc.readChains...),
				tc.maxSeqNumsPerChain,
				tc.msgScanBatchSize,
			)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expMsgs, msgs)
			mockReader.AssertExpectations(t)
		})
	}
}

func Test_observeTokenPrices(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		priceReader := mocks.NewTokenPricesReader()
		tokens := []types.Account{"0x1", "0x2", "0x3"}
		mockPrices := []*big.Int{big.NewInt(10), big.NewInt(20), big.NewInt(30)}
		priceReader.On("GetTokenPricesUSD", ctx, tokens).Return(mockPrices, nil)
		prices, err := observeTokenPrices(ctx, priceReader, tokens)
		assert.NoError(t, err)
		assert.Equal(t, []model.TokenPrice{
			model.NewTokenPrice("0x1", big.NewInt(10)),
			model.NewTokenPrice("0x2", big.NewInt(20)),
			model.NewTokenPrice("0x3", big.NewInt(30)),
		}, prices)
	})

	t.Run("price reader internal issue", func(t *testing.T) {
		priceReader := mocks.NewTokenPricesReader()
		tokens := []types.Account{"0x1", "0x2", "0x3"}
		mockPrices := []*big.Int{big.NewInt(10), big.NewInt(20)} // returned two prices for three tokens
		priceReader.On("GetTokenPricesUSD", ctx, tokens).Return(mockPrices, nil)
		_, err := observeTokenPrices(ctx, priceReader, tokens)
		assert.Error(t, err)
	})

}

func Test_validateObservedSequenceNumbers(t *testing.T) {
	testCases := []struct {
		name       string
		msgs       []model.CCIPMsgBaseDetails
		maxSeqNums []model.SeqNumChain
		expErr     bool
	}{
		{
			name:       "empty",
			msgs:       nil,
			maxSeqNums: nil,
			expErr:     false,
		},
		{
			name: "dup seq num observation",
			msgs: nil,
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
				{ChainSel: 1, SeqNum: 10},
			},
			expErr: true,
		},
		{
			name: "seq nums ok",
			msgs: nil,
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "dup msg seq num",
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 13}, // dup
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: true,
		},
		{
			name: "msg seq nums ok",
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: model.Bytes32{1}, SourceChain: 2, SeqNum: 21},
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "msg seq nums does not match observed max seq num",
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 10}, // max seq num is already 10
				{ID: model.Bytes32{1}, SourceChain: 2, SeqNum: 21},
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: true,
		},
		{
			name: "max seq num not found",
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: model.Bytes32{1}, SourceChain: 2, SeqNum: 21}, // max seq num not reported
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObservedSequenceNumbers(tc.msgs, tc.maxSeqNums)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_validateObserverReadingEligibility(t *testing.T) {
	testCases := []struct {
		name         string
		observer     commontypes.OracleID
		msgs         []model.CCIPMsgBaseDetails
		observerInfo map[commontypes.OracleID]model.ObserverInfo
		expErr       bool
	}{
		{
			name:     "observer can read all chains",
			observer: commontypes.OracleID(10),
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{3}, SourceChain: 2, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 3, SeqNum: 12},
				{ID: model.Bytes32{2}, SourceChain: 3, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]model.ObserverInfo{
				10: {Reads: []model.ChainSelector{1, 2, 3}},
			},
			expErr: false,
		},
		{
			name:     "observer cannot read one chain",
			observer: commontypes.OracleID(10),
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{3}, SourceChain: 2, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 3, SeqNum: 12},
				{ID: model.Bytes32{2}, SourceChain: 3, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]model.ObserverInfo{
				10: {Reads: []model.ChainSelector{1, 3}},
			},
			expErr: true,
		},
		{
			name:     "observer cfg not found",
			observer: commontypes.OracleID(10),
			msgs: []model.CCIPMsgBaseDetails{
				{ID: model.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: model.Bytes32{3}, SourceChain: 2, SeqNum: 12},
				{ID: model.Bytes32{1}, SourceChain: 3, SeqNum: 12},
				{ID: model.Bytes32{2}, SourceChain: 3, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]model.ObserverInfo{
				20: {Reads: []model.ChainSelector{1, 3}}, // observer 10 not found
			},
			expErr: true,
		},
		{
			name:         "no msgs",
			observer:     commontypes.OracleID(10),
			msgs:         []model.CCIPMsgBaseDetails{},
			observerInfo: map[commontypes.OracleID]model.ObserverInfo{},
			expErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObserverReadingEligibility(tc.observer, tc.msgs, tc.observerInfo)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_validateObservedTokenPrices(t *testing.T) {
	testCases := []struct {
		name        string
		tokenPrices []model.TokenPrice
		expErr      bool
	}{
		{
			name:        "empty is valid",
			tokenPrices: []model.TokenPrice{},
			expErr:      false,
		},
		{
			name: "all valid",
			tokenPrices: []model.TokenPrice{
				model.NewTokenPrice("0x1", big.NewInt(1)),
				model.NewTokenPrice("0x2", big.NewInt(1)),
				model.NewTokenPrice("0x3", big.NewInt(1)),
				model.NewTokenPrice("0xa", big.NewInt(1)),
			},
			expErr: false,
		},
		{
			name: "dup price",
			tokenPrices: []model.TokenPrice{
				model.NewTokenPrice("0x1", big.NewInt(1)),
				model.NewTokenPrice("0x2", big.NewInt(1)),
				model.NewTokenPrice("0x1", big.NewInt(1)), // dup
				model.NewTokenPrice("0xa", big.NewInt(1)),
			},
			expErr: true,
		},
		{
			name: "nil price",
			tokenPrices: []model.TokenPrice{
				model.NewTokenPrice("0x1", big.NewInt(1)),
				model.NewTokenPrice("0x2", big.NewInt(1)),
				model.NewTokenPrice("0x3", nil), // nil price
				model.NewTokenPrice("0xa", big.NewInt(1)),
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObservedTokenPrices(tc.tokenPrices)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})

	}
}

func Test_newMsgsConsensus(t *testing.T) {
	testCases := []struct {
		name           string
		maxSeqNums     []model.SeqNumChain
		observations   []model.CommitPluginObservation
		expMerkleRoots []model.MerkleRootChain
		fChain         map[model.ChainSelector]int
		expErr         bool
	}{
		{
			name:           "empty",
			maxSeqNums:     []model.SeqNumChain{},
			observations:   nil,
			expMerkleRoots: []model.MerkleRootChain{},
			expErr:         false,
		},
		{
			name: "one message but not reaching 2fChain+1 observations",
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []model.MerkleRootChain{},
			expErr:         false,
		},
		{
			name: "one message reaching 2fChain+1 observations",
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(11, 11),
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages all of them reaching 2fChain+1 observations",
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(11, 13),
				},
			},
			expErr: false,
		},
		{
			name: "one message sequence number is lower than consensus max seq num",
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(12, 13),
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages some of them not reaching 2fChain+1 observations",
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(11, 11), // we stop at 11 because there is a gap for going to 13
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages on different chains",
			fChain: map[model.ChainSelector]int{
				1: 2,
				2: 1,
			},
			maxSeqNums: []model.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			observations: []model.CommitPluginObservation{
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(11, 11), // we stop at 11 because there is a gap for going to 13
				},
				{
					ChainSel:     2,
					SeqNumsRange: model.NewSeqNumRange(21, 22), // we stop at 11 because there is a gap for going to 13
				},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			merkleRoots, err := newMsgsConsensus(lggr, tc.maxSeqNums, tc.observations, tc.fChain)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expMerkleRoots), len(merkleRoots))
			for i, exp := range tc.expMerkleRoots {
				assert.Equal(t, exp.ChainSel, merkleRoots[i].ChainSel)
				assert.Equal(t, exp.SeqNumsRange, merkleRoots[i].SeqNumsRange)
			}
		})
	}
}

func Test_maxSeqNumsConsensus(t *testing.T) {
	testCases := []struct {
		name         string
		observations []model.CommitPluginObservation
		fChain       map[model.ChainSelector]int
		destChain    model.ChainSelector
		expSeqNums   []model.SeqNumChain
		expErr       bool
	}{
		{
			name:         "empty observations",
			observations: []model.CommitPluginObservation{},
			fChain:       map[model.ChainSelector]int{1: 2},
			destChain:    1,
			expSeqNums:   []model.SeqNumChain{},
			expErr:       false,
		},
		{
			name: "one chain all followers agree",
			observations: []model.CommitPluginObservation{
				{
					MaxSeqNums: []model.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
					},
				},
			},
			fChain:    map[model.ChainSelector]int{1: 2},
			destChain: 1,
			expSeqNums: []model.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "one chain all followers agree but not enough observations",
			observations: []model.CommitPluginObservation{
				{
					MaxSeqNums: []model.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
					},
				},
			},
			fChain:     map[model.ChainSelector]int{1: 3},
			destChain:  1,
			expSeqNums: []model.SeqNumChain{},
			expErr:     false,
		},
		{
			name: "one chain 3 followers not in sync, 4 in sync",
			observations: []model.CommitPluginObservation{
				{
					MaxSeqNums: []model.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 19},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 19},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 19},
						{ChainSel: 2, SeqNum: 20},
					},
				},
			},
			fChain:    map[model.ChainSelector]int{1: 3},
			destChain: 1,
			expSeqNums: []model.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "two chains",
			observations: []model.CommitPluginObservation{
				{
					MaxSeqNums: []model.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},

						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
					},
				},
			},
			fChain: map[model.ChainSelector]int{
				1: 2,
			},
			destChain: 1,
			expSeqNums: []model.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
				{ChainSel: 3, SeqNum: 30},
			},
			expErr: false,
		},
		{
			name: "two chains but f chain is not defined for dest",
			observations: []model.CommitPluginObservation{
				{
					MaxSeqNums: []model.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},

						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
						{ChainSel: 3, SeqNum: 30},
					},
				},
			},
			fChain:    map[model.ChainSelector]int{},
			destChain: 1,
			expErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			lggr := logger.Test(t)
			p := NewPlugin(
				ctx,
				commontypes.OracleID(123),
				model.CommitPluginConfig{
					FChain:    tc.fChain,
					DestChain: tc.destChain,
				},
				nil,
				nil,
				nil,
				lggr,
			)

			seqNums, err := p.maxSeqNumsConsensus(tc.observations)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expSeqNums, seqNums)
		})
	}
}

func Test_tokenPricesConsensus(t *testing.T) {
	testCases := []struct {
		name         string
		observations []model.CommitPluginObservation
		fChain       int
		expPrices    []model.TokenPrice
		expErr       bool
	}{
		{
			name:         "empty",
			observations: make([]model.CommitPluginObservation, 0),
			fChain:       2,
			expPrices:    make([]model.TokenPrice, 0),
			expErr:       false,
		},
		{
			name: "happy flow",
			observations: []model.CommitPluginObservation{
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(10)),
						model.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(11)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(11)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(10)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(11)),
						model.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
			},
			fChain: 2,
			expPrices: []model.TokenPrice{
				model.NewTokenPrice("0x1", big.NewInt(11)),
				model.NewTokenPrice("0x2", big.NewInt(21)),
			},
			expErr: false,
		},
		{
			name: "not enough observations for some token",
			observations: []model.CommitPluginObservation{
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(11)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(11)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(10)),
						model.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []model.TokenPrice{
						model.NewTokenPrice("0x1", big.NewInt(10)),
						model.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
			},
			fChain: 2,
			expPrices: []model.TokenPrice{
				model.NewTokenPrice("0x2", big.NewInt(21)),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prices, err := tokenPricesConsensus(tc.observations, tc.fChain)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expPrices, prices)
		})
	}
}
