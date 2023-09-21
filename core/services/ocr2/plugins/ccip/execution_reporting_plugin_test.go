package ccip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/custom_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestExecutionReportingPlugin_Observation(t *testing.T) {
	testCases := []struct {
		name              string
		commitStorePaused bool
		inflightReports   []InflightInternalExecutionReport
		unexpiredReports  []ccipdata.Event[commit_store.CommitStoreReportAccepted]
		sendRequests      []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]
		executedSeqNums   []uint64
		blessedRoots      map[[32]byte]bool
		senderNonce       uint64
		rateLimiterState  evm_2_evm_offramp.RateLimiterTokenBucket
		expErr            bool
	}{
		{
			name:              "commit store is down",
			commitStorePaused: true,
			expErr:            true,
		},
		{
			name:              "happy flow",
			commitStorePaused: false,
			inflightReports:   []InflightInternalExecutionReport{},
			unexpiredReports: []ccipdata.Event[commit_store.CommitStoreReportAccepted]{
				{
					Data: commit_store.CommitStoreReportAccepted{
						Report: commit_store.CommitStoreCommitReport{
							PriceUpdates: commit_store.InternalPriceUpdates{},
							Interval:     commit_store.CommitStoreInterval{Min: 10, Max: 12},
							MerkleRoot:   [32]byte{123},
						},
					},
				},
			},
			blessedRoots: map[[32]byte]bool{
				[32]byte{123}: true,
			},
			rateLimiterState: evm_2_evm_offramp.RateLimiterTokenBucket{
				IsEnabled: false,
			},
			senderNonce: 9,
			sendRequests: []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]{
				{
					Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
						Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 10},
					},
				},
				{
					Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
						Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 11},
					},
				},
				{
					Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
						Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 12},
					},
				},
			},
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{}
			p.inflightReports = newInflightExecReportsContainer(time.Minute)
			p.inflightReports.reports = tc.inflightReports
			p.lggr = logger.TestLogger(t)

			commitStore, commitStoreAddr := testhelpers.NewFakeCommitStore(t, 1)
			commitStore.SetPaused(tc.commitStorePaused)
			commitStore.SetBlessedRoots(tc.blessedRoots)
			p.config.commitStore = commitStore

			offRamp, offRampAddr := testhelpers.NewFakeOffRamp(t)
			offRamp.SetRateLimiterState(tc.rateLimiterState)
			p.config.offRamp = offRamp

			destReader := ccipdata.NewMockReader(t)
			destReader.On("GetAcceptedCommitReportsGteTimestamp", ctx, commitStoreAddr, mock.Anything, 0).
				Return(tc.unexpiredReports, nil).Maybe()
			destReader.On("LatestBlock", ctx).Return(int64(1234), nil).Maybe()
			var executionEvents []ccipdata.Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged]
			for _, seqNum := range tc.executedSeqNums {
				executionEvents = append(executionEvents, ccipdata.Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged]{
					Data: evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{SequenceNumber: seqNum},
				})
			}
			destReader.On("GetExecutionStateChangesBetweenSeqNums", ctx, offRampAddr, mock.Anything, mock.Anything, 0).
				Return(executionEvents, nil).Maybe()
			p.config.destReader = destReader

			onRamp, onRampAddr := testhelpers.NewFakeOnRamp(t)
			p.config.onRamp = onRamp

			sourceReader := ccipdata.NewMockReader(t)
			sourceReader.On("GetSendRequestsBetweenSeqNums", ctx, onRampAddr, mock.Anything, mock.Anything, 0).
				Return(tc.sendRequests, nil).Maybe()
			p.config.sourceReader = sourceReader

			cachedDestTokens := cache.NewMockAutoSync[cache.CachedTokens](t)
			cachedDestTokens.On("Get", ctx).Return(cache.CachedTokens{
				SupportedTokens: map[common.Address]common.Address{},
				FeeTokens:       []common.Address{},
			}, nil).Maybe()
			p.cachedDestTokens = cachedDestTokens

			priceRegistry, _ := testhelpers.NewFakePriceRegistry(t)
			priceRegistry.SetTokenPrices([]price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(123), Timestamp: uint32(time.Now().Unix())},
			})
			p.destPriceRegistry = priceRegistry
			p.config.sourcePriceRegistry = priceRegistry

			sourceFeeTokens := cache.NewMockAutoSync[[]common.Address](t)
			sourceFeeTokens.On("Get", ctx).Return([]common.Address{}, nil).Maybe()
			p.cachedSourceFeeTokens = sourceFeeTokens

			p.snoozedRoots = cache.NewSnoozedRoots(time.Minute, time.Minute)

			_, err := p.Observation(ctx, types.ReportTimestamp{}, types.Query{})
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestExecutionReportingPlugin_Report(t *testing.T) {
	testCases := []struct {
		name            string
		f               int
		committedSeqNum uint64
		observations    []ExecutionObservation

		expectingSomeReport bool
		expectedReport      evm_2_evm_offramp.InternalExecutionReport
		expectingSomeErr    bool
	}{
		{
			name:            "not enough observations to form consensus",
			f:               5,
			committedSeqNum: 5,
			observations: []ExecutionObservation{
				{Messages: map[uint64]MsgData{3: {}, 4: {}}},
				{Messages: map[uint64]MsgData{3: {}, 4: {}}},
			},
			expectingSomeErr:    false,
			expectingSomeReport: false,
		},
		{
			name:                "zero observations",
			f:                   0,
			committedSeqNum:     5,
			observations:        []ExecutionObservation{},
			expectingSomeErr:    false,
			expectingSomeReport: false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := ExecutionReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.F = tc.f

			commitStore, _ := testhelpers.NewFakeCommitStore(t, tc.committedSeqNum)

			p.config.commitStore = commitStore

			observations := make([]types.AttributedObservation, len(tc.observations))
			for i := range observations {
				b, err := json.Marshal(tc.observations[i])
				assert.NoError(t, err)
				observations[i] = types.AttributedObservation{Observation: b, Observer: commontypes.OracleID(i + 1)}
			}

			_, _, err := p.Report(ctx, types.ReportTimestamp{}, types.Query{}, observations)
			if tc.expectingSomeErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}

}

func TestExecutionReportingPlugin_ShouldAcceptFinalizedReport(t *testing.T) {
	msg := evm_2_evm_offramp.InternalEVM2EVMMessage{
		SequenceNumber: 12,
		FeeTokenAmount: big.NewInt(1e9),
		Sender:         common.Address{},
		Nonce:          1,
		GasLimit:       big.NewInt(1),
		Strict:         false,
		Receiver:       common.Address{},
		Data:           nil,
		TokenAmounts:   nil,
		FeeToken:       common.Address{},
		MessageId:      [32]byte{},
	}
	report := evm_2_evm_offramp.InternalExecutionReport{
		Messages:          []evm_2_evm_offramp.InternalEVM2EVMMessage{msg},
		OffchainTokenData: [][][]byte{{}},
		Proofs:            [][32]byte{{}},
		ProofFlagBits:     big.NewInt(1),
	}
	encodedReport, err := abihelpers.EncodeExecutionReport(report)
	require.NoError(t, err)

	mockOffRamp, _ := testhelpers.NewFakeOffRamp(t)
	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			offRamp: mockOffRamp,
		},
		lggr:            logger.TestLogger(t),
		inflightReports: newInflightExecReportsContainer(models.MustMakeDuration(1 * time.Hour).Duration()),
	}

	mockedExecState := mockOffRamp.On("GetExecutionState", mock.Anything, uint64(12)).Return(uint8(abihelpers.ExecutionStateUntouched), nil).Once()

	should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, true, should)

	mockedExecState.Return(uint8(abihelpers.ExecutionStateSuccess), nil).Once()

	should, err = plugin.ShouldAcceptFinalizedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, false, should)
}

func TestExecutionReportingPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	msg := evm_2_evm_offramp.InternalEVM2EVMMessage{
		SequenceNumber: 12,
		FeeTokenAmount: big.NewInt(1e9),
		Sender:         common.Address{},
		Nonce:          1,
		GasLimit:       big.NewInt(1),
		Strict:         false,
		Receiver:       common.Address{},
		Data:           nil,
		TokenAmounts:   nil,
		FeeToken:       common.Address{},
		MessageId:      [32]byte{},
	}
	report := evm_2_evm_offramp.InternalExecutionReport{
		Messages:          []evm_2_evm_offramp.InternalEVM2EVMMessage{msg},
		OffchainTokenData: [][][]byte{{}},
		Proofs:            [][32]byte{{}},
		ProofFlagBits:     big.NewInt(1),
	}
	encodedReport, err := abihelpers.EncodeExecutionReport(report)
	require.NoError(t, err)

	mockOffRamp := &mock_contracts.EVM2EVMOffRampInterface{}
	mockCommitStore := &mock_contracts.CommitStoreInterface{}

	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			offRamp:     mockOffRamp,
			commitStore: mockCommitStore,
		},
		lggr:            logger.TestLogger(t),
		inflightReports: newInflightExecReportsContainer(models.MustMakeDuration(1 * time.Hour).Duration()),
	}

	mockedExecState := mockOffRamp.On("GetExecutionState", mock.Anything, uint64(12)).Return(uint8(abihelpers.ExecutionStateUntouched), nil).Once()

	should, err := plugin.ShouldTransmitAcceptedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, true, should)

	mockedExecState.Return(uint8(abihelpers.ExecutionStateFailure), nil).Once()
	should, err = plugin.ShouldTransmitAcceptedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, false, should)
}

func TestExecutionReportingPlugin_buildReport(t *testing.T) {
	ctx := testutils.Context(t)

	const numMessages = 100
	const tokensPerMessage = 20
	const bytesPerMessage = 1000

	executionReport := generateExecutionReport(t, numMessages, tokensPerMessage, bytesPerMessage)
	encodedReport, err := abihelpers.EncodeExecutionReport(executionReport)
	assert.NoError(t, err)
	// ensure "naive" full report would be bigger than limit
	assert.Greater(t, len(encodedReport), MaxExecutionReportLength, "full execution report length")

	observations := make([]ObservedMessage, len(executionReport.Messages))
	for i, msg := range executionReport.Messages {
		observations[i] = NewObservedMessage(msg.SequenceNumber, executionReport.OffchainTokenData[i])
	}

	// ensure that buildReport should cap the built report to fit in MaxExecutionReportLength
	p := &ExecutionReportingPlugin{}
	p.lggr = logger.TestLogger(t)

	commitStore, commitStoreAddress := testhelpers.NewFakeCommitStore(t, executionReport.Messages[len(executionReport.Messages)-1].SequenceNumber+1)
	commitStore.On("Verify", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(math.MaxInt64), nil)
	p.config.commitStore = commitStore

	destReader := ccipdata.NewMockReader(t)
	destReader.On("GetAcceptedCommitReportsGteSeqNum", ctx, commitStoreAddress, observations[0].SeqNr, 0).
		Return([]ccipdata.Event[commit_store.CommitStoreReportAccepted]{
			{
				Data: commit_store.CommitStoreReportAccepted{
					Report: commit_store.CommitStoreCommitReport{
						Interval: commit_store.CommitStoreInterval{
							Min: observations[0].SeqNr,
							Max: observations[len(observations)-1].SeqNr,
						},
					},
				},
			},
		}, nil)
	p.config.destReader = destReader

	p.config.leafHasher = leafHasher123{}

	onRamp, onRampAddr := testhelpers.NewFakeOnRamp(t)
	p.config.onRamp = onRamp

	sendReqs := make([]ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested], len(observations))
	for i := range observations {
		sendReqs[i] = ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]{
			Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{Message: evm_2_evm_onramp.InternalEVM2EVMMessage{
				SourceChainSelector: math.MaxUint64,
				SequenceNumber:      uint64(i + 1),
				FeeTokenAmount:      big.NewInt(math.MaxInt64),
				Sender:              utils.RandomAddress(),
				Nonce:               math.MaxUint64,
				GasLimit:            big.NewInt(math.MaxInt64),
				Strict:              false,
				Receiver:            utils.RandomAddress(),
				Data:                bytes.Repeat([]byte{0}, bytesPerMessage),
				TokenAmounts:        nil,
				FeeToken:            utils.RandomAddress(),
				MessageId:           [32]byte{12},
			}},
		}
	}
	sourceReader := ccipdata.NewMockReader(t)
	sourceReader.On("GetSendRequestsBetweenSeqNums",
		ctx, onRampAddr, observations[0].SeqNr, observations[len(observations)-1].SeqNr, 0).Return(sendReqs, nil)
	p.config.sourceReader = sourceReader

	execReport, err := p.buildReport(ctx, p.lggr, observations)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(execReport), MaxExecutionReportLength, "built execution report length")
}

func TestExecutionReportingPlugin_buildBatch(t *testing.T) {
	c, _ := testhelpers.SetupChain(t)
	offRamp, _ := testhelpers.NewFakeOffRamp(t)
	// We do this just to have the parsing available.
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress("0x1"), c)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)

	sender1 := common.HexToAddress("0xa")
	destNative := common.HexToAddress("0xb")
	srcNative := common.HexToAddress("0xc")
	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			offRamp: offRamp,
			onRamp:  onRamp,
		},
		destWrappedNative: destNative,
		offchainConfig: ccipconfig.ExecOffchainConfig{
			SourceFinalityDepth:         5,
			DestOptimisticConfirmations: 1,
			DestFinalityDepth:           5,
			BatchGasLimit:               300_000,
			RelativeBoostPerWaitHour:    1,
			MaxGasPrice:                 1,
		},
		lggr: logger.TestLogger(t),
	}

	msg1 := internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
			SequenceNumber: 1,
			FeeTokenAmount: big.NewInt(1e9),
			Sender:         sender1,
			Nonce:          1,
			GasLimit:       big.NewInt(1),
			Strict:         false,
			Receiver:       common.Address{},
			Data:           nil,
			TokenAmounts:   nil,
			FeeToken:       srcNative,
			MessageId:      [32]byte{},
		},
		BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
	}

	msg2 := msg1
	msg2.Executed = true

	msg3 := msg1
	msg3.Executed = true
	msg3.Finalized = true

	msg4 := msg1
	msg4.TokenAmounts = []evm_2_evm_offramp.ClientEVMTokenAmount{
		{Token: srcNative, Amount: big.NewInt(100)},
	}

	msg5 := msg4
	msg5.SequenceNumber = msg5.SequenceNumber + 1
	msg5.Nonce = msg5.Nonce + 1

	var tt = []struct {
		name                     string
		reqs                     []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta
		inflight                 []InflightInternalExecutionReport
		tokenLimit, destGasPrice *big.Int
		srcPrices, dstPrices     map[common.Address]*big.Int
		offRampNoncesBySender    map[common.Address]uint64
		destRateLimits           map[common.Address]*big.Int
		srcToDestTokens          map[common.Address]common.Address
		expectedSeqNrs           []ObservedMessage
	}{
		{
			name:                  "single message no tokens",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        []ObservedMessage{{SeqNr: uint64(1)}},
		},
		{
			name:                  "executed non finalized messages should be skipped",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg2},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name:                  "finalized executed log",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg3},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name:                  "dst token price does not exist",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg2},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name:                  "src token price does not exist",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg2},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name:                  "rate limit hit",
			reqs:                  []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg4},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			destRateLimits: map[common.Address]*big.Int{
				destNative: big.NewInt(99),
			},
			srcToDestTokens: map[common.Address]common.Address{
				srcNative: destNative,
			},
			expectedSeqNrs: nil,
		},
		{
			name:         "message with tokens is not executed if limit is reached",
			reqs:         []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg4},
			inflight:     []InflightInternalExecutionReport{},
			tokenLimit:   big.NewInt(2),
			destGasPrice: big.NewInt(10),
			srcPrices:    map[common.Address]*big.Int{srcNative: big.NewInt(1e18)},
			dstPrices:    map[common.Address]*big.Int{destNative: big.NewInt(1e18)},
			srcToDestTokens: map[common.Address]common.Address{
				srcNative: destNative,
			},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name: "message with tokens is not executed if limit is reached when inflight is full",
			reqs: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg5},
			inflight: []InflightInternalExecutionReport{
				{
					createdAt: time.Now(),
					messages:  []evm_2_evm_offramp.InternalEVM2EVMMessage{msg4.InternalEVM2EVMMessage},
				},
			},
			tokenLimit:   big.NewInt(19),
			destGasPrice: big.NewInt(10),
			srcPrices:    map[common.Address]*big.Int{srcNative: big.NewInt(1e18)},
			dstPrices:    map[common.Address]*big.Int{destNative: big.NewInt(1e18)},
			srcToDestTokens: map[common.Address]common.Address{
				srcNative: destNative,
			},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        nil,
		},
		{
			name: "some messages skipped after hitting max batch data len",
			reqs: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          1,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageId:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
						SequenceNumber: 11,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          2,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, MaxDataLenPerBatch-500), // skipped from batch
						FeeToken:       srcNative,
						MessageId:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
						SequenceNumber: 12,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          2,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageId:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:              []InflightInternalExecutionReport{},
			tokenLimit:            big.NewInt(0),
			destGasPrice:          big.NewInt(10),
			srcPrices:             map[common.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:             map[common.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender: map[common.Address]uint64{sender1: 0},
			expectedSeqNrs:        []ObservedMessage{{SeqNr: uint64(10)}, {SeqNr: uint64(12)}},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			offRamp.SetSenderNonces(tc.offRampNoncesBySender)
			seqNrs := plugin.buildBatch(
				context.Background(),
				lggr,
				commitReportWithSendRequests{sendRequestsWithMeta: tc.reqs},
				tc.inflight,
				tc.tokenLimit,
				tc.srcPrices,
				tc.dstPrices,
				func() (*big.Int, error) { return tc.destGasPrice, nil },
				tc.srcToDestTokens,
				tc.destRateLimits,
			)
			assert.Equal(t, tc.expectedSeqNrs, seqNrs)
		})
	}
}

func TestExecutionReportingPlugin_isRateLimitEnoughForTokenPool(t *testing.T) {
	testCases := []struct {
		name                    string
		destTokenPoolRateLimits map[common.Address]*big.Int
		tokenAmounts            []evm_2_evm_offramp.ClientEVMTokenAmount
		inflightTokenAmounts    map[common.Address]*big.Int
		srcToDestToken          map[common.Address]common.Address
		exp                     bool
	}{
		{
			name: "base",
			destTokenPoolRateLimits: map[common.Address]*big.Int{
				common.HexToAddress("10"): big.NewInt(100),
				common.HexToAddress("20"): big.NewInt(50),
			},
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: common.HexToAddress("1"), Amount: big.NewInt(50)},
				{Token: common.HexToAddress("2"), Amount: big.NewInt(20)},
			},
			srcToDestToken: map[common.Address]common.Address{
				common.HexToAddress("1"): common.HexToAddress("10"),
				common.HexToAddress("2"): common.HexToAddress("20"),
			},
			inflightTokenAmounts: map[common.Address]*big.Int{
				common.HexToAddress("1"): big.NewInt(20),
				common.HexToAddress("2"): big.NewInt(30),
			},
			exp: true,
		},
		{
			name: "rate limit hit",
			destTokenPoolRateLimits: map[common.Address]*big.Int{
				common.HexToAddress("10"): big.NewInt(100),
				common.HexToAddress("20"): big.NewInt(50),
			},
			srcToDestToken: map[common.Address]common.Address{
				common.HexToAddress("1"): common.HexToAddress("10"),
				common.HexToAddress("2"): common.HexToAddress("20"),
			},
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: common.HexToAddress("1"), Amount: big.NewInt(50)},
				{Token: common.HexToAddress("2"), Amount: big.NewInt(51)},
			},
			exp: true,
		},
		{
			name: "rate limit hit, inflight included",
			destTokenPoolRateLimits: map[common.Address]*big.Int{
				common.HexToAddress("10"): big.NewInt(100),
				common.HexToAddress("20"): big.NewInt(50),
			},
			srcToDestToken: map[common.Address]common.Address{
				common.HexToAddress("1"): common.HexToAddress("10"),
				common.HexToAddress("2"): common.HexToAddress("20"),
			},
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: common.HexToAddress("1"), Amount: big.NewInt(50)},
				{Token: common.HexToAddress("2"), Amount: big.NewInt(20)},
			},
			inflightTokenAmounts: map[common.Address]*big.Int{
				common.HexToAddress("1"): big.NewInt(51),
				common.HexToAddress("2"): big.NewInt(30),
			},
			exp: true,
		},
		{
			destTokenPoolRateLimits: map[common.Address]*big.Int{},
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: common.HexToAddress("1"), Amount: big.NewInt(50)},
				{Token: common.HexToAddress("2"), Amount: big.NewInt(20)},
			},
			srcToDestToken: map[common.Address]common.Address{
				common.HexToAddress("1"): common.HexToAddress("10"),
				common.HexToAddress("2"): common.HexToAddress("20"),
			},
			inflightTokenAmounts: map[common.Address]*big.Int{
				common.HexToAddress("1"): big.NewInt(20),
				common.HexToAddress("2"): big.NewInt(30),
			},
			name: "rate limit not applied to token",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{lggr: logger.TestLogger(t)}
			p.isRateLimitEnoughForTokenPool(tc.destTokenPoolRateLimits, tc.tokenAmounts, tc.inflightTokenAmounts, tc.srcToDestToken)
		})
	}
}

func TestExecutionReportingPlugin_destPoolRateLimits(t *testing.T) {
	tk1 := utils.RandomAddress()
	tk1dest := utils.RandomAddress()
	tk1pool := utils.RandomAddress()

	tk2 := utils.RandomAddress()
	tk2dest := utils.RandomAddress()
	tk2pool := utils.RandomAddress()

	testCases := []struct {
		name              string
		tokenAmounts      []evm_2_evm_offramp.ClientEVMTokenAmount
		sourceToDestToken map[common.Address]common.Address
		destPools         map[common.Address]common.Address
		poolRateLimits    map[common.Address]custom_token_pool.RateLimiterTokenBucket

		expRateLimits map[common.Address]*big.Int
		expErr        bool
	}{
		{
			name: "happy flow",
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: tk1},
				{Token: tk2},
				{Token: tk1},
				{Token: tk1},
			},
			sourceToDestToken: map[common.Address]common.Address{
				tk1: tk1dest,
				tk2: tk2dest,
			},
			destPools: map[common.Address]common.Address{
				tk1dest: tk1pool,
				tk2dest: tk2pool,
			},
			poolRateLimits: map[common.Address]custom_token_pool.RateLimiterTokenBucket{
				tk1pool: {Tokens: big.NewInt(1000), IsEnabled: true},
				tk2pool: {Tokens: big.NewInt(2000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
				tk2dest: big.NewInt(2000),
			},
			expErr: false,
		},
		{
			name: "token missing from source to dest mapping",
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: tk1},
				{Token: tk2}, // <-- missing form sourceToDestToken
			},
			sourceToDestToken: map[common.Address]common.Address{
				tk1: tk1dest,
			},
			destPools: map[common.Address]common.Address{
				tk1dest: tk1pool,
			},
			poolRateLimits: map[common.Address]custom_token_pool.RateLimiterTokenBucket{
				tk1pool: {Tokens: big.NewInt(1000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			expErr: false,
		},
		{
			name: "pool is disabled",
			tokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
				{Token: tk1},
				{Token: tk2},
			},
			sourceToDestToken: map[common.Address]common.Address{
				tk1: tk1dest,
				tk2: tk2dest,
			},
			destPools: map[common.Address]common.Address{
				tk1dest: tk1pool,
				tk2dest: tk2pool,
			},
			poolRateLimits: map[common.Address]custom_token_pool.RateLimiterTokenBucket{
				tk1pool: {Tokens: big.NewInt(1000), IsEnabled: true},
				tk2pool: {Tokens: big.NewInt(2000), IsEnabled: false}, // <--- pool disabled
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{}
			p.lggr = lggr

			offRamp, offRampAddr := testhelpers.NewFakeOffRamp(t)
			offRamp.SetTokenPools(tc.destPools)
			p.config.offRamp = offRamp

			p.customTokenPoolFactory = func(ctx context.Context, poolAddress common.Address, _ bind.ContractBackend) (custom_token_pool.CustomTokenPoolInterface, error) {
				mp := &mockPool{}
				mp.On("CurrentOffRampRateLimiterState", mock.Anything, offRampAddr).Return(tc.poolRateLimits[poolAddress], nil)
				return mp, nil
			}

			rateLimits, err := p.destPoolRateLimits(ctx, []commitReportWithSendRequests{
				{
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
								TokenAmounts: tc.tokenAmounts,
							},
						},
					},
				},
			}, tc.sourceToDestToken)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRateLimits, rateLimits)
		})
	}
}

func TestExecutionReportingPlugin_estimateDestinationGasPrice(t *testing.T) {
	testCases := []struct {
		name      string
		evmFee    gas.EvmFee
		evmFeeErr error

		expRes *big.Int
		expErr bool
	}{
		{
			name: "dynamic fee cap has precedence over legacy",
			evmFee: gas.EvmFee{
				Legacy:        assets.NewWei(big.NewInt(1000)),
				DynamicFeeCap: assets.NewWei(big.NewInt(2000)),
			},
			expRes: big.NewInt(2000),
		},
		{
			name: "legacy is used if dynamic fee cap is not provided",
			evmFee: gas.EvmFee{
				Legacy: assets.NewWei(big.NewInt(1000)),
			},
			expRes: big.NewInt(1000),
		},
		{
			name:      "stop on error",
			evmFeeErr: errors.New("some error"),
			expErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{}
			mockEstimator := mocks.NewEvmFeeEstimator(t)
			mockEstimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.evmFee, uint32(0), tc.evmFeeErr)
			p.config.destGasEstimator = mockEstimator

			res, err := p.estimateDestinationGasPrice(testutils.Context(t))
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestExecutionReportingPlugin_getReportsWithSendRequests(t *testing.T) {
	testCases := []struct {
		name                string
		reports             []commit_store.CommitStoreCommitReport
		expQueryMin         uint64 // expected min/max used in the query to get ccipevents
		expQueryMax         uint64
		onchainEvents       []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]
		destLatestBlock     int64
		destExecutedSeqNums []uint64

		expReports []commitReportWithSendRequests
		expErr     bool
	}{
		{
			name:       "no reports",
			reports:    nil,
			expReports: nil,
			expErr:     false,
		},
		{
			name: "two reports happy flow",
			reports: []commit_store.CommitStoreCommitReport{
				{
					Interval:   commit_store.CommitStoreInterval{Min: 1, Max: 2},
					MerkleRoot: [32]byte{100},
				},
				{
					Interval:   commit_store.CommitStoreInterval{Min: 3, Max: 3},
					MerkleRoot: [32]byte{200},
				},
			},
			expQueryMin: 1,
			expQueryMax: 3,
			onchainEvents: []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]{
				{Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
					Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 1},
				}},
				{Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
					Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 2},
				}},
				{Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
					Message: evm_2_evm_onramp.InternalEVM2EVMMessage{SequenceNumber: 3},
				}},
			},
			destLatestBlock:     10_000,
			destExecutedSeqNums: []uint64{1},
			expReports: []commitReportWithSendRequests{
				{
					commitReport: commit_store.CommitStoreCommitReport{
						Interval:   commit_store.CommitStoreInterval{Min: 1, Max: 2},
						MerkleRoot: [32]byte{100},
					},
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 1},
							Executed:               true,
							Finalized:              true,
						},
						{
							InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 2},
							Executed:               false,
							Finalized:              false,
						},
					},
				},
				{
					commitReport: commit_store.CommitStoreCommitReport{
						Interval:   commit_store.CommitStoreInterval{Min: 3, Max: 3},
						MerkleRoot: [32]byte{200},
					},
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 3},
							Executed:               false,
							Finalized:              false,
						},
					},
				},
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{}
			p.lggr = lggr

			onRamp, onRampAddr := testhelpers.NewFakeOnRamp(t)
			p.config.onRamp = onRamp

			offRamp, offRampAddr := testhelpers.NewFakeOffRamp(t)
			p.config.offRamp = offRamp

			sourceReader := ccipdata.NewMockReader(t)
			sourceReader.On("GetSendRequestsBetweenSeqNums", ctx, onRampAddr, tc.expQueryMin, tc.expQueryMax, 0).
				Return(tc.onchainEvents, nil).Maybe()
			p.config.sourceReader = sourceReader

			destReader := ccipdata.NewMockReader(t)
			destReader.On("LatestBlock", ctx).Return(tc.destLatestBlock, nil).Maybe()
			var executedEvents []ccipdata.Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged]
			for _, executedSeqNum := range tc.destExecutedSeqNums {
				executedEvents = append(executedEvents, ccipdata.Event[evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged]{
					Data:      evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{SequenceNumber: executedSeqNum},
					BlockMeta: ccipdata.BlockMeta{BlockNumber: tc.destLatestBlock - 10},
				})
			}
			destReader.On("GetExecutionStateChangesBetweenSeqNums", ctx, offRampAddr, tc.expQueryMin, tc.expQueryMax, 0).Return(executedEvents, nil).Maybe()
			p.config.destReader = destReader

			populatedReports, err := p.getReportsWithSendRequests(ctx, tc.reports)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expReports), len(populatedReports))
			for i, expReport := range tc.expReports {
				assert.Equal(t, len(expReport.sendRequestsWithMeta), len(populatedReports[i].sendRequestsWithMeta))
				for j, expReq := range expReport.sendRequestsWithMeta {
					assert.Equal(t, expReq.Executed, populatedReports[i].sendRequestsWithMeta[j].Executed)
					assert.Equal(t, expReq.Finalized, populatedReports[i].sendRequestsWithMeta[j].Finalized)
					assert.Equal(t, expReq.SequenceNumber, populatedReports[i].sendRequestsWithMeta[j].SequenceNumber)
				}
			}
		})
	}
}

func TestExecutionReportingPluginFactory_UpdateLogPollerFilters(t *testing.T) {
	const numFilters = 10
	filters := make([]logpoller.Filter, numFilters)
	for i := range filters {
		filters[i] = logpoller.Filter{
			Name:      fmt.Sprintf("filter-%d", i),
			EventSigs: []common.Hash{common.HexToHash(fmt.Sprintf("%d", i))},
			Addresses: []common.Address{common.HexToAddress(fmt.Sprintf("%d", i))},
			Retention: time.Duration(i) * time.Second,
		}
	}

	destLP := lpMocks.NewLogPoller(t)
	sourceLP := lpMocks.NewLogPoller(t)

	onRamp, _ := testhelpers.NewFakeOnRamp(t)
	sourcePriceRegistry, _ := testhelpers.NewFakePriceRegistry(t)
	commitStore, _ := testhelpers.NewFakeCommitStore(t, 1)
	offRamp, _ := testhelpers.NewFakeOffRamp(t)

	destPriceRegistryAddr := utils.RandomAddress()

	tokenDataProviders := make(map[common.Address]tokendata.Reader)

	rf := &ExecutionReportingPluginFactory{
		filtersMu:          &sync.Mutex{},
		sourceChainFilters: filters[:5],
		destChainFilters:   filters[5:10],
		config: ExecutionPluginConfig{
			destLP:              destLP,
			sourceLP:            sourceLP,
			onRamp:              onRamp,
			commitStore:         commitStore,
			offRamp:             offRamp,
			sourcePriceRegistry: sourcePriceRegistry,
			tokenDataProviders:  tokenDataProviders,
		},
	}

	for _, f := range getExecutionPluginSourceLpChainFilters(onRamp.Address(), sourcePriceRegistry.Address(), tokenDataProviders) {
		sourceLP.On("RegisterFilter", f).Return(nil)
	}
	for _, f := range getExecutionPluginDestLpChainFilters(commitStore.Address(), offRamp.Address(), destPriceRegistryAddr) {
		destLP.On("RegisterFilter", f).Return(nil)
	}
	for _, f := range rf.sourceChainFilters[1:] { // zero address is skipped
		sourceLP.On("UnregisterFilter", f.Name, mock.Anything).Return(nil)
	}
	for _, f := range rf.destChainFilters {
		destLP.On("UnregisterFilter", f.Name, mock.Anything).Return(nil)
	}

	err := rf.UpdateLogPollerFilters(destPriceRegistryAddr)
	assert.NoError(t, err)
}

func TestExecutionReportToEthTxMeta(t *testing.T) {
	t.Run("happy flow", func(t *testing.T) {
		executionReport := generateExecutionReport(t, 10, 3, 1000)
		encExecReport, err := abihelpers.EncodeExecutionReport(executionReport)
		assert.NoError(t, err)
		txMeta, err := ExecutionReportToEthTxMeta(encExecReport)
		assert.NoError(t, err)
		assert.Len(t, txMeta.MessageIDs, len(executionReport.Messages))
	})

	t.Run("invalid report", func(t *testing.T) {
		_, err := ExecutionReportToEthTxMeta([]byte("whatever"))
		assert.Error(t, err)
	})
}

func TestUpdateSourceToDestTokenMapping(t *testing.T) {
	expectedNewBlockNumber := int64(10000)
	logs := []logpoller.Log{{BlockNumber: expectedNewBlockNumber}}
	mockDestLP := &lpMocks.LogPoller{}

	mockDestLP.On("LatestLogEventSigsAddrsWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(logs, nil)
	mockDestLP.On("LatestBlock", mock.Anything).Return(expectedNewBlockNumber, nil)

	sourceToken, destToken := common.HexToAddress("111111"), common.HexToAddress("222222")

	mockOffRamp := &mock_contracts.EVM2EVMOffRampInterface{}
	mockOffRamp.On("Address").Return(common.HexToAddress("0x01"))
	mockOffRamp.On("GetSupportedTokens", mock.Anything).Return([]common.Address{sourceToken}, nil)
	mockOffRamp.On("GetDestinationToken", mock.Anything, sourceToken).Return(destToken, nil)

	mockPriceRegistry := &mock_contracts.PriceRegistryInterface{}
	mockPriceRegistry.On("Address").Return(common.HexToAddress("0x02"))
	mockPriceRegistry.On("GetFeeTokens", mock.Anything).Return([]common.Address{}, nil)

	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			destLP:  mockDestLP,
			offRamp: mockOffRamp,
		},
		cachedDestTokens: cache.NewCachedSupportedTokens(mockDestLP, mockOffRamp, mockPriceRegistry, 0),
	}

	value, err := plugin.cachedDestTokens.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, destToken, value.SupportedTokens[sourceToken])
}

func Test_calculateObservedMessagesConsensus(t *testing.T) {
	type args struct {
		observations []ExecutionObservation
		f            int
	}
	tests := []struct {
		name string
		args args
		want []ObservedMessage
	}{
		{
			name: "no observations",
			args: args{
				observations: nil,
				f:            0,
			},
			want: []ObservedMessage{},
		},
		{
			name: "common path",
			args: args{
				observations: []ExecutionObservation{
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
						},
					},
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0xff}}}, // different token data - should not be picked
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
							3: {TokenData: [][]byte{{0x3}, {0x3}, {0x3}}},
						},
					},
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
						},
					},
				},
				f: 1,
			},
			want: []ObservedMessage{
				{SeqNr: 1, MsgData: MsgData{TokenData: [][]byte{{0x1}, {0x1}, {0x1}}}},
				{SeqNr: 2, MsgData: MsgData{TokenData: [][]byte{{0x2}, {0x2}, {0x2}}}},
			},
		},
		{
			name: "similar token data",
			args: args{
				observations: []ExecutionObservation{
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
						},
					},
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1, 0x1}}},
						},
					},
					{
						Messages: map[uint64]MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1, 0x1}}},
						},
					},
				},
				f: 1,
			},
			want: []ObservedMessage{
				{SeqNr: 1, MsgData: MsgData{TokenData: [][]byte{{0x1}, {0x1, 0x1}}}},
			},
		},
		{
			name: "results should be deterministic",
			args: args{
				observations: []ExecutionObservation{
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x2}}}}},
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x2}}}}},
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x1}}}}},
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x3}}}}},
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x3}}}}},
					{Messages: map[uint64]MsgData{1: {TokenData: [][]byte{{0x1}}}}},
				},
				f: 1,
			},
			want: []ObservedMessage{
				{SeqNr: 1, MsgData: MsgData{TokenData: [][]byte{{0x1}}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := calculateObservedMessagesConsensus(
				tt.args.observations,
				tt.args.f,
			)
			assert.NoError(t, err)
			sort.Slice(res, func(i, j int) bool {
				return res[i].SeqNr < res[j].SeqNr
			})
			assert.Equalf(t, tt.want, res, "calculateObservedMessagesConsensus(%v, %v)", tt.args.observations, tt.args.f)
		})
	}
}

func Test_getTokensPrices(t *testing.T) {
	tk1 := common.HexToAddress("1")
	tk2 := common.HexToAddress("2")
	tk3 := common.HexToAddress("3")

	testCases := []struct {
		name      string
		feeTokens []common.Address
		tokens    []common.Address
		retPrices []price_registry.InternalTimestampedPackedUint224
		expPrices map[common.Address]*big.Int
		expErr    bool
	}{
		{
			name:      "base",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(10)},
				{Value: big.NewInt(20)},
				{Value: big.NewInt(30)},
			},
			expPrices: map[common.Address]*big.Int{
				tk1: big.NewInt(10),
				tk2: big.NewInt(20),
				tk3: big.NewInt(30),
			},
			expErr: false,
		},
		{
			name:      "token is both fee token and normal token",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3, tk1},
			retPrices: []price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(10)},
				{Value: big.NewInt(20)},
				{Value: big.NewInt(30)},
				{Value: big.NewInt(10)},
			},
			expPrices: map[common.Address]*big.Int{
				tk1: big.NewInt(10),
				tk2: big.NewInt(20),
				tk3: big.NewInt(30),
			},
			expErr: false,
		},
		{
			name:      "token is both fee token and normal token and price registry gave different price",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3, tk1},
			retPrices: []price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(10)},
				{Value: big.NewInt(20)},
				{Value: big.NewInt(30)},
				{Value: big.NewInt(1000)}, // different price for same token
			},
			expErr: true,
		},
		{
			name:      "zero price should lead to an error",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(10)},
				{Value: big.NewInt(0)},
				{Value: big.NewInt(30)},
			},
			expErr: true,
		},
		{
			name:      "contract returns less prices than requested",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []price_registry.InternalTimestampedPackedUint224{
				{Value: big.NewInt(10)},
				{Value: big.NewInt(20)},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceReg, _ := testhelpers.NewFakePriceRegistry(t)
			priceReg.SetTokenPrices(tc.retPrices)

			prices, err := getTokensPrices(context.Background(), tc.feeTokens, priceReg, tc.tokens)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			for tk, price := range tc.expPrices {
				assert.Equal(t, price, prices[tk])
			}
		})
	}
}

func Test_calculateMessageMaxGas(t *testing.T) {
	type args struct {
		gasLimit    *big.Int
		numRequests int
		dataLen     int
		numTokens   int
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name:    "base",
			args:    args{gasLimit: big.NewInt(1000), numRequests: 5, dataLen: 5, numTokens: 2},
			want:    203836,
			wantErr: false,
		},
		{
			name:    "large",
			args:    args{gasLimit: big.NewInt(1000), numRequests: 1000, dataLen: 1000, numTokens: 1000},
			want:    36482676,
			wantErr: false,
		},
		{
			name:    "gas limit overflow",
			args:    args{gasLimit: big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))},
			want:    36391540,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateMessageMaxGas(tt.args.gasLimit, tt.args.numRequests, tt.args.dataLen, tt.args.numTokens)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "calculateMessageMaxGas(%v, %v, %v, %v)", tt.args.gasLimit, tt.args.numRequests, tt.args.dataLen, tt.args.numTokens)
		})
	}
}

func Test_inflightAggregates(t *testing.T) {
	const n = 10
	addrs := make([]common.Address, n)
	tokenAddrs := make([]common.Address, n)
	for i := range addrs {
		addrs[i] = utils.RandomAddress()
		tokenAddrs[i] = utils.RandomAddress()
	}

	testCases := []struct {
		name            string
		inflight        []InflightInternalExecutionReport
		destTokenPrices map[common.Address]*big.Int
		sourceToDest    map[common.Address]common.Address

		expInflightSeqNrs          map[uint64]struct{}
		expInflightAggrVal         *big.Int
		expMaxInflightSenderNonces map[common.Address]uint64
		expInflightTokenAmounts    map[common.Address]*big.Int
		expErr                     bool
	}{
		{
			name: "base",
			inflight: []InflightInternalExecutionReport{
				{
					messages: []evm_2_evm_offramp.InternalEVM2EVMMessage{
						{
							Sender:         addrs[0],
							SequenceNumber: 100,
							Nonce:          2,
							TokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
								{Token: tokenAddrs[0], Amount: big.NewInt(1e18)},
								{Token: tokenAddrs[0], Amount: big.NewInt(2e18)},
							},
						},
						{
							Sender:         addrs[0],
							SequenceNumber: 106,
							Nonce:          4,
							TokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
								{Token: tokenAddrs[0], Amount: big.NewInt(1e18)},
								{Token: tokenAddrs[0], Amount: big.NewInt(5e18)},
								{Token: tokenAddrs[2], Amount: big.NewInt(5e18)},
							},
						},
					},
				},
			},
			destTokenPrices: map[common.Address]*big.Int{
				tokenAddrs[1]: big.NewInt(1000),
				tokenAddrs[3]: big.NewInt(500),
			},
			sourceToDest: map[common.Address]common.Address{
				tokenAddrs[0]: tokenAddrs[1],
				tokenAddrs[2]: tokenAddrs[3],
			},
			expInflightSeqNrs: map[uint64]struct{}{
				100: {},
				106: {},
			},
			expInflightAggrVal: big.NewInt(9*1000 + 5*500),
			expMaxInflightSenderNonces: map[common.Address]uint64{
				addrs[0]: 4,
			},
			expInflightTokenAmounts: map[common.Address]*big.Int{
				tokenAddrs[0]: big.NewInt(9e18),
				tokenAddrs[2]: big.NewInt(5e18),
			},
			expErr: false,
		},
		{
			name: "missing price",
			inflight: []InflightInternalExecutionReport{
				{
					messages: []evm_2_evm_offramp.InternalEVM2EVMMessage{
						{
							Sender:         addrs[0],
							SequenceNumber: 100,
							Nonce:          2,
							TokenAmounts: []evm_2_evm_offramp.ClientEVMTokenAmount{
								{Token: tokenAddrs[0], Amount: big.NewInt(1e18)},
							},
						},
					},
				},
			},
			destTokenPrices: map[common.Address]*big.Int{
				tokenAddrs[3]: big.NewInt(500),
			},
			sourceToDest: map[common.Address]common.Address{
				tokenAddrs[2]: tokenAddrs[3],
			},
			expErr: true,
		},
		{
			name:                       "nothing inflight",
			inflight:                   []InflightInternalExecutionReport{},
			expInflightSeqNrs:          map[uint64]struct{}{},
			expInflightAggrVal:         big.NewInt(0),
			expMaxInflightSenderNonces: map[common.Address]uint64{},
			expInflightTokenAmounts:    map[common.Address]*big.Int{},
			expErr:                     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inflightSeqNrs, inflightAggrVal, maxInflightSenderNonces, inflightTokenAmounts, err := inflightAggregates(
				tc.inflight, tc.destTokenPrices, tc.sourceToDest)

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tc.expInflightSeqNrs, inflightSeqNrs))
			assert.True(t, reflect.DeepEqual(tc.expInflightAggrVal, inflightAggrVal))
			assert.True(t, reflect.DeepEqual(tc.expMaxInflightSenderNonces, maxInflightSenderNonces))
			assert.True(t, reflect.DeepEqual(tc.expInflightTokenAmounts, inflightTokenAmounts))
		})
	}
}

func Test_commitReportWithSendRequests_validate(t *testing.T) {
	testCases := []struct {
		name           string
		reportInterval commit_store.CommitStoreInterval
		numReqs        int
		expValid       bool
	}{
		{
			name:           "valid report",
			reportInterval: commit_store.CommitStoreInterval{Min: 10, Max: 20},
			numReqs:        11,
			expValid:       true,
		},
		{
			name:           "report with one request",
			reportInterval: commit_store.CommitStoreInterval{Min: 1234, Max: 1234},
			numReqs:        1,
			expValid:       true,
		},
		{
			name:           "request is missing",
			reportInterval: commit_store.CommitStoreInterval{Min: 1234, Max: 1234},
			numReqs:        0,
			expValid:       false,
		},
		{
			name:           "requests are missing",
			reportInterval: commit_store.CommitStoreInterval{Min: 1, Max: 10},
			numReqs:        5,
			expValid:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rep := commitReportWithSendRequests{
				commitReport: commit_store.CommitStoreCommitReport{
					Interval: tc.reportInterval,
				},
				sendRequestsWithMeta: make([]internal.EVM2EVMOnRampCCIPSendRequestedWithMeta, tc.numReqs),
			}
			err := rep.validate()
			isValid := err == nil
			assert.Equal(t, tc.expValid, isValid)
		})
	}
}

func Test_commitReportWithSendRequests_allRequestsAreExecutedAndFinalized(t *testing.T) {
	testCases := []struct {
		name   string
		reqs   []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta
		expRes bool
	}{
		{
			name: "all requests executed and finalized",
			reqs: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{Executed: true, Finalized: true},
				{Executed: true, Finalized: true},
				{Executed: true, Finalized: true},
			},
			expRes: true,
		},
		{
			name:   "true when there are zero requests",
			reqs:   []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{},
			expRes: true,
		},
		{
			name: "some request not executed",
			reqs: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{Executed: true, Finalized: true},
				{Executed: true, Finalized: true},
				{Executed: false, Finalized: true},
			},
			expRes: false,
		},
		{
			name: "some request not finalized",
			reqs: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{Executed: true, Finalized: true},
				{Executed: true, Finalized: true},
				{Executed: true, Finalized: false},
			},
			expRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rep := commitReportWithSendRequests{sendRequestsWithMeta: tc.reqs}
			res := rep.allRequestsAreExecutedAndFinalized()
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func Test_commitReportWithSendRequests_sendReqFits(t *testing.T) {
	testCases := []struct {
		name   string
		req    internal.EVM2EVMOnRampCCIPSendRequestedWithMeta
		report commit_store.CommitStoreCommitReport
		expRes bool
	}{
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 1},
			},
			report: commit_store.CommitStoreCommitReport{
				Interval: commit_store.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: true,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 10},
			},
			report: commit_store.CommitStoreCommitReport{
				Interval: commit_store.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: true,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 11},
			},
			report: commit_store.CommitStoreCommitReport{
				Interval: commit_store.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: false,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{SequenceNumber: 10},
			},
			report: commit_store.CommitStoreCommitReport{
				Interval: commit_store.CommitStoreInterval{Min: 10, Max: 10},
			},
			expRes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &commitReportWithSendRequests{commitReport: tc.report}
			assert.Equal(t, tc.expRes, r.sendReqFits(tc.req))
		})
	}
}

// generateExecutionReport generates an execution report that can be used in tests
func generateExecutionReport(t *testing.T, numMsgs, tokensPerMsg, bytesPerMsg int) evm_2_evm_offramp.InternalExecutionReport {
	messages := make([]evm_2_evm_offramp.InternalEVM2EVMMessage, numMsgs)

	offChainTokenData := make([][][]byte, numMsgs)
	for i := range messages {
		tokenAmounts := make([]evm_2_evm_offramp.ClientEVMTokenAmount, tokensPerMsg)
		for j := range tokenAmounts {
			tokenAmounts[j] = evm_2_evm_offramp.ClientEVMTokenAmount{
				Token:  utils.RandomAddress(),
				Amount: big.NewInt(math.MaxInt64),
			}
		}

		messages[i] = evm_2_evm_offramp.InternalEVM2EVMMessage{
			SourceChainSelector: rand.Uint64(),
			SequenceNumber:      uint64(i + 1),
			FeeTokenAmount:      big.NewInt(rand.Int64()),
			Sender:              utils.RandomAddress(),
			Nonce:               rand.Uint64(),
			GasLimit:            big.NewInt(rand.Int64()),
			Strict:              false,
			Receiver:            utils.RandomAddress(),
			Data:                bytes.Repeat([]byte{1}, bytesPerMsg),
			TokenAmounts:        tokenAmounts,
			FeeToken:            utils.RandomAddress(),
			MessageId:           utils.RandomBytes32(),
		}

		data := []byte(`{"foo": "bar"}`)
		offChainTokenData[i] = [][]byte{data, data, data}
	}

	return evm_2_evm_offramp.InternalExecutionReport{
		Messages:          messages,
		OffchainTokenData: offChainTokenData,
		Proofs:            make([][32]byte, numMsgs),
		ProofFlagBits:     big.NewInt(rand.Int64()),
	}
}

type mockPool struct {
	custom_token_pool.CustomTokenPoolInterface
	mock.Mock
}

func (mp *mockPool) CurrentOffRampRateLimiterState(opts *bind.CallOpts, offRamp common.Address) (custom_token_pool.RateLimiterTokenBucket, error) {
	args := mp.Called(opts, offRamp)
	return args.Get(0).(custom_token_pool.RateLimiterTokenBucket), args.Error(1)
}
