package commit

import (
	"context"
	"math/big"
	"reflect"
	"slices"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/mocks"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func Test_observeMaxSeqNumsPerChain(t *testing.T) {
	testCases := []struct {
		name             string
		prevOutcome      cciptypes.CommitPluginOutcome
		onChainSeqNums   map[cciptypes.ChainSelector]cciptypes.SeqNum
		readChains       []cciptypes.ChainSelector
		destChain        cciptypes.ChainSelector
		expErr           bool
		expSeqNumsInSync bool
		expMaxSeqNums    []cciptypes.SeqNumChain
	}{
		{
			name:        "report on chain seq num when no previous outcome and can read dest",
			prevOutcome: cciptypes.CommitPluginOutcome{},
			onChainSeqNums: map[cciptypes.ChainSelector]cciptypes.SeqNum{
				1: 10,
				2: 20,
			},
			readChains:       []cciptypes.ChainSelector{1, 2, 3},
			destChain:        3,
			expErr:           false,
			expSeqNumsInSync: true,
			expMaxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
		},
		{
			name:        "nothing to report when there is no previous outcome and cannot read dest",
			prevOutcome: cciptypes.CommitPluginOutcome{},
			onChainSeqNums: map[cciptypes.ChainSelector]cciptypes.SeqNum{
				1: 10,
				2: 20,
			},
			readChains:       []cciptypes.ChainSelector{1, 2},
			destChain:        3,
			expErr:           false,
			expSeqNumsInSync: false,
			expMaxSeqNums:    []cciptypes.SeqNumChain{},
		},
		{
			name: "report previous outcome seq nums and override when on chain is higher if can read dest",
			prevOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: 1, SeqNum: 11}, // for chain 1 previous outcome is higher than on-chain state
					{ChainSel: 2, SeqNum: 19}, // for chain 2 previous outcome is behind on-chain state
				},
			},
			onChainSeqNums: map[cciptypes.ChainSelector]cciptypes.SeqNum{
				1: 10,
				2: 20,
			},
			readChains:       []cciptypes.ChainSelector{1, 2, 3},
			destChain:        3,
			expErr:           false,
			expSeqNumsInSync: true,
			expMaxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 11},
				{ChainSel: 2, SeqNum: 20},
			},
		},
		{
			name: "report previous outcome seq nums and mark as non synced if cannot read dest",
			prevOutcome: cciptypes.CommitPluginOutcome{
				MaxSeqNums: []cciptypes.SeqNumChain{
					{ChainSel: 1, SeqNum: 11}, // for chain 1 previous outcome is higher than on-chain state
					{ChainSel: 2, SeqNum: 19}, // for chain 2 previous outcome is behind on-chain state
				},
			},
			onChainSeqNums: map[cciptypes.ChainSelector]cciptypes.SeqNum{
				1: 10,
				2: 20,
			},
			readChains:       []cciptypes.ChainSelector{1, 2},
			destChain:        3,
			expErr:           false,
			expSeqNumsInSync: false,
			expMaxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 11},
				{ChainSel: 2, SeqNum: 19},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := mocks.NewCCIPReader()
			knownSourceChains := slicelib.Filter(tc.readChains, func(ch cciptypes.ChainSelector) bool { return ch != tc.destChain })
			lggr := logger.Test(t)

			var encodedPrevOutcome []byte
			var err error
			if !reflect.DeepEqual(tc.prevOutcome, cciptypes.CommitPluginOutcome{}) {
				encodedPrevOutcome, err = tc.prevOutcome.Encode()
				assert.NoError(t, err)
			}

			onChainSeqNums := make([]cciptypes.SeqNum, 0)
			for _, chain := range knownSourceChains {
				if v, ok := tc.onChainSeqNums[chain]; !ok {
					t.Fatalf("invalid test case missing on chain seq num expectation for %d", chain)
				} else {
					onChainSeqNums = append(onChainSeqNums, v)
				}
			}
			mockReader.On("NextSeqNum", ctx, knownSourceChains).Return(onChainSeqNums, nil)

			seqNums, synced, err := observeMaxSeqNums(
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
			assert.Equal(t, tc.expSeqNumsInSync, synced)
		})
	}
}

func Test_observeNewMsgs(t *testing.T) {
	testCases := []struct {
		name               string
		maxSeqNumsPerChain []cciptypes.SeqNumChain
		readChains         []cciptypes.ChainSelector
		destChain          cciptypes.ChainSelector
		msgScanBatchSize   int
		newMsgs            map[cciptypes.ChainSelector][]cciptypes.CCIPMsg
		expMsgs            []cciptypes.CCIPMsg
		expErr             bool
	}{
		{
			name: "no new messages",
			maxSeqNumsPerChain: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []cciptypes.ChainSelector{1, 2},
			msgScanBatchSize: 256,
			newMsgs: map[cciptypes.ChainSelector][]cciptypes.CCIPMsg{
				1: {},
				2: {},
			},
			expMsgs: []cciptypes.CCIPMsg{},
			expErr:  false,
		},
		{
			name: "new messages",
			maxSeqNumsPerChain: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []cciptypes.ChainSelector{1, 2},
			msgScanBatchSize: 256,
			newMsgs: map[cciptypes.ChainSelector][]cciptypes.CCIPMsg{
				1: {
					{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}},
				},
				2: {
					{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
					{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
				},
			},
			expMsgs: []cciptypes.CCIPMsg{
				{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}},
				{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
				{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
			},
			expErr: false,
		},
		{
			name: "new messages but one chain is not readable",
			maxSeqNumsPerChain: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			readChains:       []cciptypes.ChainSelector{2},
			msgScanBatchSize: 256,
			newMsgs: map[cciptypes.ChainSelector][]cciptypes.CCIPMsg{
				2: {
					{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
					{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
				},
			},
			expMsgs: []cciptypes.CCIPMsg{
				{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
				{CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := mocks.NewCCIPReader()
			msgHasher := mocks.NewMessageHasher()
			lggr := logger.Test(t)

			for _, seqNumChain := range tc.maxSeqNumsPerChain {
				if slices.Contains(tc.readChains, seqNumChain.ChainSel) {
					mockReader.On(
						"MsgsBetweenSeqNums",
						ctx,
						seqNumChain.ChainSel,
						cciptypes.NewSeqNumRange(seqNumChain.SeqNum+1, seqNumChain.SeqNum+cciptypes.SeqNum(1+tc.msgScanBatchSize)),
					).Return(tc.newMsgs[seqNumChain.ChainSel], nil)
				}
			}

			msgs, err := observeNewMsgs(
				ctx,
				lggr,
				mockReader,
				msgHasher,
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

func Benchmark_observeNewMsgs(b *testing.B) {
	const (
		numChains       = 5
		readerDelayMS   = 100
		newMsgsPerChain = 256
	)

	readChains := make([]cciptypes.ChainSelector, numChains)
	maxSeqNumsPerChain := make([]cciptypes.SeqNumChain, numChains)
	for i := 0; i < numChains; i++ {
		readChains[i] = cciptypes.ChainSelector(i + 1)
		maxSeqNumsPerChain[i] = cciptypes.SeqNumChain{ChainSel: cciptypes.ChainSelector(i + 1), SeqNum: cciptypes.SeqNum(1)}
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		lggr, _ := logger.New()
		ccipReader := mocks.NewCCIPReader()
		msgHasher := mocks.NewMessageHasher()

		expNewMsgs := make([]cciptypes.CCIPMsg, 0, newMsgsPerChain*numChains)
		for _, seqNumChain := range maxSeqNumsPerChain {
			newMsgs := make([]cciptypes.CCIPMsg, 0, newMsgsPerChain)
			for msgSeqNum := 1; msgSeqNum <= newMsgsPerChain; msgSeqNum++ {
				newMsgs = append(newMsgs, cciptypes.CCIPMsg{
					CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
						ID:          cciptypes.Bytes32{byte(msgSeqNum)},
						SourceChain: seqNumChain.ChainSel,
						SeqNum:      cciptypes.SeqNum(msgSeqNum),
					},
				})
			}

			ccipReader.On(
				"MsgsBetweenSeqNums",
				ctx,
				[]cciptypes.ChainSelector{seqNumChain.ChainSel},
				cciptypes.NewSeqNumRange(
					seqNumChain.SeqNum+1,
					seqNumChain.SeqNum+cciptypes.SeqNum(1+newMsgsPerChain),
				),
			).Run(func(args mock.Arguments) {
				time.Sleep(time.Duration(readerDelayMS) * time.Millisecond)
			}).Return(newMsgs, nil)
			expNewMsgs = append(expNewMsgs, newMsgs...)
		}

		msgs, err := observeNewMsgs(
			ctx,
			lggr,
			ccipReader,
			msgHasher,
			mapset.NewSet(readChains...),
			maxSeqNumsPerChain,
			newMsgsPerChain,
		)
		assert.NoError(b, err)
		assert.Equal(b, expNewMsgs, msgs)

		// (old)     sequential: 509.345 ms/op   (numChains * readerDelayMS)
		// (current) parallel:   102.543 ms/op     (readerDelayMS)
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
		assert.Equal(t, []cciptypes.TokenPrice{
			cciptypes.NewTokenPrice("0x1", big.NewInt(10)),
			cciptypes.NewTokenPrice("0x2", big.NewInt(20)),
			cciptypes.NewTokenPrice("0x3", big.NewInt(30)),
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

func Test_observeGasPrices(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		mockReader := mocks.NewCCIPReader()
		chains := []cciptypes.ChainSelector{1, 2, 3}
		mockGasPrices := []cciptypes.BigInt{
			cciptypes.NewBigIntFromInt64(10),
			cciptypes.NewBigIntFromInt64(20),
			cciptypes.NewBigIntFromInt64(30),
		}
		mockReader.On("GasPrices", ctx, chains).Return(mockGasPrices, nil)
		gasPrices, err := observeGasPrices(ctx, mockReader, chains)
		assert.NoError(t, err)
		assert.Equal(t, []cciptypes.GasPriceChain{
			cciptypes.NewGasPriceChain(mockGasPrices[0].Int, chains[0]),
			cciptypes.NewGasPriceChain(mockGasPrices[1].Int, chains[1]),
			cciptypes.NewGasPriceChain(mockGasPrices[2].Int, chains[2]),
		}, gasPrices)
	})

	t.Run("gas reader internal issue", func(t *testing.T) {
		mockReader := mocks.NewCCIPReader()
		chains := []cciptypes.ChainSelector{1, 2, 3}
		mockGasPrices := []cciptypes.BigInt{
			cciptypes.NewBigIntFromInt64(10),
			cciptypes.NewBigIntFromInt64(20),
		} // return 2 prices for 3 chains
		mockReader.On("GasPrices", ctx, chains).Return(mockGasPrices, nil)
		_, err := observeGasPrices(ctx, mockReader, chains)
		assert.Error(t, err)
	})
}

func Test_validateObservedSequenceNumbers(t *testing.T) {
	testCases := []struct {
		name       string
		msgs       []cciptypes.CCIPMsgBaseDetails
		maxSeqNums []cciptypes.SeqNumChain
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
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
				{ChainSel: 1, SeqNum: 10},
			},
			expErr: true,
		},
		{
			name: "seq nums ok",
			msgs: nil,
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "dup msg seq num",
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 13}, // dup
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: true,
		},
		{
			name: "msg seq nums ok",
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: cciptypes.Bytes32{1}, SourceChain: 2, SeqNum: 21},
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: false,
		},
		{
			name: "msg seq nums does not match observed max seq num",
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 10}, // max seq num is already 10
				{ID: cciptypes.Bytes32{1}, SourceChain: 2, SeqNum: 21},
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			expErr: true,
		},
		{
			name: "max seq num not found",
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 13},
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 14},
				{ID: cciptypes.Bytes32{1}, SourceChain: 2, SeqNum: 21}, // max seq num not reported
			},
			maxSeqNums: []cciptypes.SeqNumChain{
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
		msgs         []cciptypes.CCIPMsgBaseDetails
		seqNums      []cciptypes.SeqNumChain
		observerInfo map[commontypes.OracleID]cciptypes.ObserverInfo
		expErr       bool
	}{
		{
			name:     "observer can read all chains",
			observer: commontypes.OracleID(10),
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{3}, SourceChain: 2, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 3, SeqNum: 12},
				{ID: cciptypes.Bytes32{2}, SourceChain: 3, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]cciptypes.ObserverInfo{
				10: {Reads: []cciptypes.ChainSelector{1, 2, 3}},
			},
			expErr: false,
		},
		{
			name:     "observer is a writer so can observe seq nums",
			observer: commontypes.OracleID(10),
			msgs:     []cciptypes.CCIPMsgBaseDetails{},
			seqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]cciptypes.ObserverInfo{
				10: {Reads: []cciptypes.ChainSelector{1, 3}, Writer: true},
			},
			expErr: false,
		},
		{
			name:     "observer is not a writer so cannot observe seq nums",
			observer: commontypes.OracleID(10),
			msgs:     []cciptypes.CCIPMsgBaseDetails{},
			seqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]cciptypes.ObserverInfo{
				10: {Reads: []cciptypes.ChainSelector{1, 3}, Writer: false},
			},
			expErr: true,
		},
		{
			name:     "observer cfg not found",
			observer: commontypes.OracleID(10),
			msgs: []cciptypes.CCIPMsgBaseDetails{
				{ID: cciptypes.Bytes32{1}, SourceChain: 1, SeqNum: 12},
				{ID: cciptypes.Bytes32{3}, SourceChain: 2, SeqNum: 12},
				{ID: cciptypes.Bytes32{1}, SourceChain: 3, SeqNum: 12},
				{ID: cciptypes.Bytes32{2}, SourceChain: 3, SeqNum: 12},
			},
			observerInfo: map[commontypes.OracleID]cciptypes.ObserverInfo{
				20: {Reads: []cciptypes.ChainSelector{1, 3}}, // observer 10 not found
			},
			expErr: true,
		},
		{
			name:     "no msgs",
			observer: commontypes.OracleID(10),
			msgs:     []cciptypes.CCIPMsgBaseDetails{},
			observerInfo: map[commontypes.OracleID]cciptypes.ObserverInfo{
				10: {Reads: []cciptypes.ChainSelector{1, 3}},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObserverReadingEligibility(tc.observer, tc.msgs, tc.seqNums, tc.observerInfo)
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
		tokenPrices []cciptypes.TokenPrice
		expErr      bool
	}{
		{
			name:        "empty is valid",
			tokenPrices: []cciptypes.TokenPrice{},
			expErr:      false,
		},
		{
			name: "all valid",
			tokenPrices: []cciptypes.TokenPrice{
				cciptypes.NewTokenPrice("0x1", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x2", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x3", big.NewInt(1)),
				cciptypes.NewTokenPrice("0xa", big.NewInt(1)),
			},
			expErr: false,
		},
		{
			name: "dup price",
			tokenPrices: []cciptypes.TokenPrice{
				cciptypes.NewTokenPrice("0x1", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x2", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x1", big.NewInt(1)), // dup
				cciptypes.NewTokenPrice("0xa", big.NewInt(1)),
			},
			expErr: true,
		},
		{
			name: "nil price",
			tokenPrices: []cciptypes.TokenPrice{
				cciptypes.NewTokenPrice("0x1", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x2", big.NewInt(1)),
				cciptypes.NewTokenPrice("0x3", nil), // nil price
				cciptypes.NewTokenPrice("0xa", big.NewInt(1)),
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

func Test_validateObservedGasPrices(t *testing.T) {
	testCases := []struct {
		name      string
		gasPrices []cciptypes.GasPriceChain
		expErr    bool
	}{
		{
			name:      "empty is valid",
			gasPrices: []cciptypes.GasPriceChain{},
			expErr:    false,
		},
		{
			name: "all valid",
			gasPrices: []cciptypes.GasPriceChain{
				cciptypes.NewGasPriceChain(big.NewInt(10), 1),
				cciptypes.NewGasPriceChain(big.NewInt(20), 2),
				cciptypes.NewGasPriceChain(big.NewInt(1312), 3),
			},
			expErr: false,
		},
		{
			name: "duplicate gas price",
			gasPrices: []cciptypes.GasPriceChain{
				cciptypes.NewGasPriceChain(big.NewInt(10), 1),
				cciptypes.NewGasPriceChain(big.NewInt(20), 2),
				cciptypes.NewGasPriceChain(big.NewInt(1312), 1), // notice we already have a gas price for chain 1
			},
			expErr: true,
		},
		{
			name: "empty gas price",
			gasPrices: []cciptypes.GasPriceChain{
				cciptypes.NewGasPriceChain(big.NewInt(10), 1),
				cciptypes.NewGasPriceChain(big.NewInt(20), 2),
				cciptypes.NewGasPriceChain(nil, 3), // nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateObservedGasPrices(tc.gasPrices)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_newMsgsConsensusForChain(t *testing.T) {
	testCases := []struct {
		name           string
		maxSeqNums     []cciptypes.SeqNumChain
		observations   []cciptypes.CommitPluginObservation
		expMerkleRoots []cciptypes.MerkleRootChain
		fChain         map[cciptypes.ChainSelector]int
		expErr         bool
	}{
		{
			name:           "empty",
			maxSeqNums:     []cciptypes.SeqNumChain{},
			observations:   nil,
			expMerkleRoots: []cciptypes.MerkleRootChain{},
			expErr:         false,
		},
		{
			name: "one message but not reaching 2fChain+1 observations",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{},
			expErr:         false,
		},
		{
			name: "one message reaching 2fChain+1 observations",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(11, 11),
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages all of them reaching 2fChain+1 observations",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(11, 13),
				},
			},
			expErr: false,
		},
		{
			name: "one message sequence number is lower than consensus max seq num",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 10}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(12, 13),
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages some of them not reaching 2fChain+1 observations",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 12}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 13}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(11, 11), // we stop at 11 because there is a gap for going to 13
				},
			},
			expErr: false,
		},
		{
			name: "multiple messages on different chains",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
				2: 1,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
				{ChainSel: 2, SeqNum: 20},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 2, SeqNum: 21}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{4}, SourceChain: 2, SeqNum: 22}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(11, 11), // we stop at 11 because there is a gap for going to 13
				},
				{
					ChainSel:     2,
					SeqNumsRange: cciptypes.NewSeqNumRange(21, 22), // we stop at 11 because there is a gap for going to 13
				},
			},
			expErr: false,
		},
		{
			name: "one message seq num with multiple reported ids",
			fChain: map[cciptypes.ChainSelector]int{
				1: 2,
			},
			maxSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 1, SeqNum: 10},
			},
			observations: []cciptypes.CommitPluginObservation{
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{10}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{10}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{111}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{111}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []cciptypes.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []cciptypes.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: cciptypes.NewSeqNumRange(11, 11),
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
		observations []cciptypes.CommitPluginObservation
		fChain       int
		expSeqNums   []cciptypes.SeqNumChain
	}{
		{
			name:         "empty observations",
			observations: []cciptypes.CommitPluginObservation{},
			fChain:       2,
			expSeqNums:   []cciptypes.SeqNumChain{},
		},
		{
			name: "one chain all followers agree",
			observations: []cciptypes.CommitPluginObservation{
				{
					MaxSeqNums: []cciptypes.SeqNumChain{
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
			fChain: 2,
			expSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
			},
		},
		{
			name: "one chain all followers agree but not enough observations",
			observations: []cciptypes.CommitPluginObservation{
				{
					MaxSeqNums: []cciptypes.SeqNumChain{
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
						{ChainSel: 2, SeqNum: 20},
					},
				},
			},
			fChain:     3,
			expSeqNums: []cciptypes.SeqNumChain{},
		},
		{
			name: "one chain 3 followers not in sync, 4 in sync",
			observations: []cciptypes.CommitPluginObservation{
				{
					MaxSeqNums: []cciptypes.SeqNumChain{
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
			fChain: 3,
			expSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
			},
		},
		{
			name: "two chains",
			observations: []cciptypes.CommitPluginObservation{
				{
					MaxSeqNums: []cciptypes.SeqNumChain{
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
			fChain: 2,
			expSeqNums: []cciptypes.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
				{ChainSel: 3, SeqNum: 30},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			seqNums := maxSeqNumsConsensus(lggr, tc.fChain, tc.observations)
			assert.Equal(t, tc.expSeqNums, seqNums)
		})
	}
}

func Test_tokenPricesConsensus(t *testing.T) {
	testCases := []struct {
		name         string
		observations []cciptypes.CommitPluginObservation
		fChain       int
		expPrices    []cciptypes.TokenPrice
		expErr       bool
	}{
		{
			name:         "empty",
			observations: make([]cciptypes.CommitPluginObservation, 0),
			fChain:       2,
			expPrices:    make([]cciptypes.TokenPrice, 0),
			expErr:       false,
		},
		{
			name: "happy flow",
			observations: []cciptypes.CommitPluginObservation{
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(10)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(10)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
			},
			fChain: 2,
			expPrices: []cciptypes.TokenPrice{
				cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
				cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
			},
			expErr: false,
		},
		{
			name: "not enough observations for some token",
			observations: []cciptypes.CommitPluginObservation{
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(11)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(10)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
					},
				},
				{
					TokenPrices: []cciptypes.TokenPrice{
						cciptypes.NewTokenPrice("0x1", big.NewInt(10)),
						cciptypes.NewTokenPrice("0x2", big.NewInt(20)),
					},
				},
			},
			fChain: 2,
			expPrices: []cciptypes.TokenPrice{
				cciptypes.NewTokenPrice("0x2", big.NewInt(21)),
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

func Test_gasPricesConsensus(t *testing.T) {
	testCases := []struct {
		name         string
		observations []cciptypes.CommitPluginObservation
		fChain       int
		expPrices    []cciptypes.GasPriceChain
	}{
		{
			name:         "empty",
			observations: make([]cciptypes.CommitPluginObservation, 0),
			fChain:       2,
			expPrices:    make([]cciptypes.GasPriceChain, 0),
		},
		{
			name: "one chain happy path",
			observations: []cciptypes.CommitPluginObservation{
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
			},
			fChain: 2,
			expPrices: []cciptypes.GasPriceChain{
				cciptypes.NewGasPriceChain(big.NewInt(10), 1),
			},
		},
		{
			name: "one chain no consensus",
			observations: []cciptypes.CommitPluginObservation{
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
			},
			fChain:    3, // notice fChain is 3, means we need at least 2*3+1=7 observations
			expPrices: []cciptypes.GasPriceChain{},
		},
		{
			name: "two chains determinism check",
			observations: []cciptypes.CommitPluginObservation{
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(200), 10)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(100), 10)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(100), 10)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(110), 10)}},
				{GasPrices: []cciptypes.GasPriceChain{cciptypes.NewGasPriceChain(big.NewInt(100), 10)}},
			},
			fChain: 2,
			expPrices: []cciptypes.GasPriceChain{
				cciptypes.NewGasPriceChain(big.NewInt(10), 1),
				cciptypes.NewGasPriceChain(big.NewInt(100), 10),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			prices := gasPricesConsensus(lggr, tc.observations, tc.fChain)
			assert.Equal(t, tc.expPrices, prices)
		})
	}
}
