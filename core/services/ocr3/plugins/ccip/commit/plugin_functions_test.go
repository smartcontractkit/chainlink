package commit

import (
	"context"
	"math/big"
	"reflect"
	"slices"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/mocks"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
		expMsgs            []model.CCIPMsg
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
			expMsgs: []model.CCIPMsg{},
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
			expMsgs: []model.CCIPMsg{
				{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{1}, SourceChain: 1, SeqNum: 11}},
				{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
				{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
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
			expMsgs: []model.CCIPMsg{
				{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{2}, SourceChain: 2, SeqNum: 21}},
				{CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{ID: [32]byte{3}, SourceChain: 2, SeqNum: 22}},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := mocks.NewCCIPReader()
			msgHasher := mocks.NewNopMessageHasher()
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

	readChains := make([]model.ChainSelector, numChains)
	maxSeqNumsPerChain := make([]model.SeqNumChain, numChains)
	for i := 0; i < numChains; i++ {
		readChains[i] = model.ChainSelector(i + 1)
		maxSeqNumsPerChain[i] = model.SeqNumChain{ChainSel: model.ChainSelector(i + 1), SeqNum: model.SeqNum(1)}
	}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		lggr, _ := logger.New()
		ccipReader := mocks.NewCCIPReader()
		msgHasher := mocks.NewNopMessageHasher()

		expNewMsgs := make([]model.CCIPMsg, 0, newMsgsPerChain*numChains)
		for _, seqNumChain := range maxSeqNumsPerChain {
			newMsgs := make([]model.CCIPMsg, 0, newMsgsPerChain)
			for msgSeqNum := 1; msgSeqNum <= newMsgsPerChain; msgSeqNum++ {
				newMsgs = append(newMsgs, model.CCIPMsg{
					CCIPMsgBaseDetails: model.CCIPMsgBaseDetails{
						ID:          model.Bytes32{byte(msgSeqNum)},
						SourceChain: seqNumChain.ChainSel,
						SeqNum:      model.SeqNum(msgSeqNum),
					},
				})
			}

			ccipReader.On(
				"MsgsBetweenSeqNums",
				ctx,
				[]model.ChainSelector{seqNumChain.ChainSel},
				model.NewSeqNumRange(
					seqNumChain.SeqNum+1,
					seqNumChain.SeqNum+model.SeqNum(1+newMsgsPerChain),
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

func Test_observeGasPrices(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		mockReader := mocks.NewCCIPReader()
		chains := []model.ChainSelector{1, 2, 3}
		mockGasPrices := []model.BigInt{
			{Int: big.NewInt(10)},
			{Int: big.NewInt(20)},
			{Int: big.NewInt(30)},
		}
		mockReader.On("GasPrices", ctx, chains).Return(mockGasPrices, nil)
		gasPrices, err := observeGasPrices(ctx, mockReader, chains)
		assert.NoError(t, err)
		assert.Equal(t, []model.GasPriceChain{
			model.NewGasPriceChain(mockGasPrices[0].Int, chains[0]),
			model.NewGasPriceChain(mockGasPrices[1].Int, chains[1]),
			model.NewGasPriceChain(mockGasPrices[2].Int, chains[2]),
		}, gasPrices)
	})

	t.Run("gas reader internal issue", func(t *testing.T) {
		mockReader := mocks.NewCCIPReader()
		chains := []model.ChainSelector{1, 2, 3}
		mockGasPrices := []model.BigInt{
			{Int: big.NewInt(10)},
			{Int: big.NewInt(20)},
		} // return 2 prices for 3 chains
		mockReader.On("GasPrices", ctx, chains).Return(mockGasPrices, nil)
		_, err := observeGasPrices(ctx, mockReader, chains)
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

func Test_validateObservedGasPrices(t *testing.T) {
	testCases := []struct {
		name      string
		gasPrices []model.GasPriceChain
		expErr    bool
	}{
		{
			name:      "empty is valid",
			gasPrices: []model.GasPriceChain{},
			expErr:    false,
		},
		{
			name: "all valid",
			gasPrices: []model.GasPriceChain{
				model.NewGasPriceChain(big.NewInt(10), 1),
				model.NewGasPriceChain(big.NewInt(20), 2),
				model.NewGasPriceChain(big.NewInt(1312), 3),
			},
			expErr: false,
		},
		{
			name: "duplicate gas price",
			gasPrices: []model.GasPriceChain{
				model.NewGasPriceChain(big.NewInt(10), 1),
				model.NewGasPriceChain(big.NewInt(20), 2),
				model.NewGasPriceChain(big.NewInt(1312), 1), // notice we already have a gas price for chain 1
			},
			expErr: true,
		},
		{
			name: "empty gas price",
			gasPrices: []model.GasPriceChain{
				model.NewGasPriceChain(big.NewInt(10), 1),
				model.NewGasPriceChain(big.NewInt(20), 2),
				model.NewGasPriceChain(nil, 3), // nil
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
		{
			name: "one message seq num with multiple reported ids",
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

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{10}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{10}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{111}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{111}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{3}, SourceChain: 1, SeqNum: 11}}},

				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
				{NewMsgs: []model.CCIPMsgBaseDetails{{ID: [32]byte{2}, SourceChain: 1, SeqNum: 11}}},
			},
			expMerkleRoots: []model.MerkleRootChain{
				{
					ChainSel:     1,
					SeqNumsRange: model.NewSeqNumRange(11, 11),
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
		fChain       int
		expSeqNums   []model.SeqNumChain
		expErr       bool
	}{
		{
			name:         "empty observations",
			observations: []model.CommitPluginObservation{},
			fChain:       2,
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
			fChain: 2,
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
			fChain:     3,
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
			fChain: 3,
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
			fChain: 2,
			expSeqNums: []model.SeqNumChain{
				{ChainSel: 2, SeqNum: 20},
				{ChainSel: 3, SeqNum: 30},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			seqNums, err := maxSeqNumsConsensus(lggr, tc.fChain, tc.observations)
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

func Test_gasPricesConsensus(t *testing.T) {
	testCases := []struct {
		name         string
		observations []model.CommitPluginObservation
		fChain       int
		expPrices    []model.GasPriceChain
		expErr       bool
	}{
		{
			name:         "empty",
			observations: make([]model.CommitPluginObservation, 0),
			fChain:       2,
			expPrices:    make([]model.GasPriceChain, 0),
			expErr:       false,
		},
		{
			name: "one chain happy path",
			observations: []model.CommitPluginObservation{
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
			},
			fChain: 2,
			expPrices: []model.GasPriceChain{
				model.NewGasPriceChain(big.NewInt(10), 1),
			},
			expErr: false,
		},
		{
			name: "one chain no consensus",
			observations: []model.CommitPluginObservation{
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
			},
			fChain:    3, // notice fChain is 3, means we need at least 2*3+1=7 observations
			expPrices: []model.GasPriceChain{},
			expErr:    false,
		},
		{
			name: "two chains determinism check",
			observations: []model.CommitPluginObservation{
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(20), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(11), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(10), 1)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(200), 10)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(100), 10)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(100), 10)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(110), 10)}},
				{GasPrices: []model.GasPriceChain{model.NewGasPriceChain(big.NewInt(100), 10)}},
			},
			fChain: 2,
			expPrices: []model.GasPriceChain{
				model.NewGasPriceChain(big.NewInt(10), 1),
				model.NewGasPriceChain(big.NewInt(100), 10),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)
			prices, err := gasPricesConsensus(lggr, tc.observations, tc.fChain)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expPrices, prices)
		})
	}
}
