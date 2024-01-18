package ccipexec

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/rand"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestExecutionReportingPlugin_Observation(t *testing.T) {
	testCases := []struct {
		name              string
		commitStorePaused bool
		inflightReports   []InflightInternalExecutionReport
		unexpiredReports  []ccipdata.Event[ccipdata.CommitStoreReport]
		sendRequests      []ccipdata.Event[internal.EVM2EVMMessage]
		executedSeqNums   []uint64
		tokenPoolsMapping map[common.Address]common.Address
		blessedRoots      map[[32]byte]bool
		senderNonce       uint64
		rateLimiterState  ccipdata.TokenBucketRateLimit
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
			unexpiredReports: []ccipdata.Event[ccipdata.CommitStoreReport]{
				{
					Data: ccipdata.CommitStoreReport{
						Interval:   ccipdata.CommitStoreInterval{Min: 10, Max: 12},
						MerkleRoot: [32]byte{123},
					},
				},
			},
			blessedRoots: map[[32]byte]bool{
				{123}: true,
			},
			rateLimiterState: ccipdata.TokenBucketRateLimit{
				IsEnabled: false,
			},
			tokenPoolsMapping: map[common.Address]common.Address{},
			senderNonce:       9,
			sendRequests: []ccipdata.Event[internal.EVM2EVMMessage]{
				{
					Data: internal.EVM2EVMMessage{SequenceNumber: 10},
				},
				{
					Data: internal.EVM2EVMMessage{SequenceNumber: 11},
				},
				{
					Data: internal.EVM2EVMMessage{SequenceNumber: 12},
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
			p.tokenDataWorker = tokendata.NewBackgroundWorker(
				ctx, make(map[common.Address]tokendata.Reader), 10, 5*time.Second, time.Hour)
			p.metricsCollector = ccip.NoopMetricsCollector

			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("IsDown", mock.Anything).Return(tc.commitStorePaused, nil)
			// Blessed roots return true
			for root, blessed := range tc.blessedRoots {
				commitStoreReader.On("IsBlessed", mock.Anything, root).Return(blessed, nil).Maybe()
			}
			commitStoreReader.On("GetAcceptedCommitReportsGteTimestamp", ctx, mock.Anything, 0).
				Return(tc.unexpiredReports, nil).Maybe()
			p.commitStoreReader = commitStoreReader

			var executionEvents []ccipdata.Event[ccipdata.ExecutionStateChanged]
			for _, seqNum := range tc.executedSeqNums {
				executionEvents = append(executionEvents, ccipdata.Event[ccipdata.ExecutionStateChanged]{
					Data: ccipdata.ExecutionStateChanged{SequenceNumber: seqNum},
				})
			}

			offRamp, _ := testhelpers.NewFakeOffRamp(t)
			offRamp.SetRateLimiterState(tc.rateLimiterState)

			mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
			mockOffRampReader.On("GetExecutionStateChangesBetweenSeqNums", ctx, mock.Anything, mock.Anything, 0).
				Return(executionEvents, nil).Maybe()
			mockOffRampReader.On("CurrentRateLimiterState", mock.Anything).Return(tc.rateLimiterState, nil).Maybe()
			mockOffRampReader.On("Address").Return(offRamp.Address()).Maybe()
			mockOffRampReader.On("GetSenderNonce", mock.Anything, mock.Anything).Return(offRamp.GetSenderNonce(nil, utils.RandomAddress())).Maybe()
			mockOffRampReader.On("GetTokenPoolsRateLimits", ctx, []common.Address{}).
				Return([]ccipdata.TokenBucketRateLimit{}, nil).Maybe()
			mockOffRampReader.On("GetSourceToDestTokensMapping", ctx).Return(nil, nil).Maybe()
			mockOffRampReader.On("GetTokens", ctx).Return(ccipdata.OffRampTokens{
				DestinationTokens: []common.Address{},
				SourceTokens:      []common.Address{},
			}, nil).Maybe()
			p.offRampReader = mockOffRampReader

			mockOnRampReader := ccipdatamocks.NewOnRampReader(t)
			mockOnRampReader.On("GetSendRequestsBetweenSeqNums", ctx, mock.Anything, mock.Anything, false).
				Return(tc.sendRequests, nil).Maybe()
			p.onRampReader = mockOnRampReader

			mockGasPriceEstimator := prices.NewMockGasPriceEstimatorExec(t)
			mockGasPriceEstimator.On("GetGasPrice", ctx).Return(prices.GasPrice(big.NewInt(1)), nil).Maybe()
			p.gasPriceEstimator = mockGasPriceEstimator

			destPriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceRegReader.On("GetTokenPrices", ctx, mock.Anything).Return(
				[]ccipdata.TokenPriceUpdate{{TokenPrice: ccipdata.TokenPrice{Token: common.HexToAddress("0x1"), Value: big.NewInt(123)}, TimestampUnixSec: big.NewInt(time.Now().Unix())}}, nil).Maybe()
			destPriceRegReader.On("Address").Return(utils.RandomAddress()).Maybe()
			destPriceRegReader.On("GetFeeTokens", ctx).Return([]common.Address{}, nil).Maybe()
			sourcePriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
			sourcePriceRegReader.On("Address").Return(utils.RandomAddress()).Maybe()
			sourcePriceRegReader.On("GetFeeTokens", ctx).Return([]common.Address{}, nil).Maybe()
			sourcePriceRegReader.On("GetTokenPrices", ctx, mock.Anything).Return(
				[]ccipdata.TokenPriceUpdate{{TokenPrice: ccipdata.TokenPrice{Token: common.HexToAddress("0x1"), Value: big.NewInt(123)}, TimestampUnixSec: big.NewInt(time.Now().Unix())}}, nil).Maybe()
			p.destPriceRegistry = destPriceRegReader
			p.sourcePriceRegistry = sourcePriceRegReader

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
		observations    []ccip.ExecutionObservation

		expectingSomeReport bool
		expectedReport      ccipdata.ExecReport
		expectingSomeErr    bool
	}{
		{
			name:            "not enough observations to form consensus",
			f:               5,
			committedSeqNum: 5,
			observations: []ccip.ExecutionObservation{
				{Messages: map[uint64]ccip.MsgData{3: {}, 4: {}}},
				{Messages: map[uint64]ccip.MsgData{3: {}, 4: {}}},
			},
			expectingSomeErr:    false,
			expectingSomeReport: false,
		},
		{
			name:                "zero observations",
			f:                   0,
			committedSeqNum:     5,
			observations:        []ccip.ExecutionObservation{},
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

			p.commitStoreReader = ccipdatamocks.NewCommitStoreReader(t)

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
	msg := internal.EVM2EVMMessage{
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
	report := ccipdata.ExecReport{
		Messages:          []internal.EVM2EVMMessage{msg},
		OffchainTokenData: [][][]byte{{}},
		Proofs:            [][32]byte{{}},
		ProofFlagBits:     big.NewInt(1),
	}

	encodedReport := encodeExecutionReport(t, report)
	mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
	mockOffRampReader.On("DecodeExecutionReport", encodedReport).Return(report, nil)

	plugin := ExecutionReportingPlugin{
		offRampReader:   mockOffRampReader,
		lggr:            logger.TestLogger(t),
		inflightReports: newInflightExecReportsContainer(1 * time.Hour),
	}

	mockedExecState := mockOffRampReader.On("GetExecutionState", mock.Anything, uint64(12)).Return(uint8(ccipdata.ExecutionStateUntouched), nil).Once()

	should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, true, should)

	mockedExecState.Return(uint8(ccipdata.ExecutionStateSuccess), nil).Once()

	should, err = plugin.ShouldAcceptFinalizedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, false, should)
}

func TestExecutionReportingPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	msg := internal.EVM2EVMMessage{
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
	report := ccipdata.ExecReport{
		Messages:          []internal.EVM2EVMMessage{msg},
		OffchainTokenData: [][][]byte{{}},
		Proofs:            [][32]byte{{}},
		ProofFlagBits:     big.NewInt(1),
	}
	encodedReport := encodeExecutionReport(t, report)

	mockCommitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
	mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
	mockOffRampReader.On("DecodeExecutionReport", encodedReport).Return(report, nil)
	mockedExecState := mockOffRampReader.On("GetExecutionState", mock.Anything, uint64(12)).Return(uint8(ccipdata.ExecutionStateUntouched), nil).Once()

	plugin := ExecutionReportingPlugin{
		commitStoreReader: mockCommitStoreReader,
		offRampReader:     mockOffRampReader,
		lggr:              logger.TestLogger(t),
		inflightReports:   newInflightExecReportsContainer(1 * time.Hour),
	}

	should, err := plugin.ShouldTransmitAcceptedReport(testutils.Context(t), ocrtypes.ReportTimestamp{}, encodedReport)
	require.NoError(t, err)
	assert.Equal(t, true, should)

	mockedExecState.Return(uint8(ccipdata.ExecutionStateFailure), nil).Once()
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
	encodedReport := encodeExecutionReport(t, executionReport)
	// ensure "naive" full report would be bigger than limit
	assert.Greater(t, len(encodedReport), MaxExecutionReportLength, "full execution report length")

	observations := make([]ccip.ObservedMessage, len(executionReport.Messages))
	for i, msg := range executionReport.Messages {
		observations[i] = ccip.NewObservedMessage(msg.SequenceNumber, executionReport.OffchainTokenData[i])
	}

	// ensure that buildReport should cap the built report to fit in MaxExecutionReportLength
	p := &ExecutionReportingPlugin{}
	p.lggr = logger.TestLogger(t)

	commitStore := ccipdatamocks.NewCommitStoreReader(t)
	commitStore.On("VerifyExecutionReport", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	commitStore.On("GetExpectedNextSequenceNumber", mock.Anything).
		Return(executionReport.Messages[len(executionReport.Messages)-1].SequenceNumber+1, nil)
	commitStore.On("GetCommitReportMatchingSeqNum", ctx, observations[0].SeqNr, 0).
		Return([]ccipdata.Event[ccipdata.CommitStoreReport]{
			{
				Data: ccipdata.CommitStoreReport{
					Interval: ccipdata.CommitStoreInterval{
						Min: observations[0].SeqNr,
						Max: observations[len(observations)-1].SeqNr,
					},
				},
			},
		}, nil)
	p.metricsCollector = ccip.NoopMetricsCollector
	p.commitStoreReader = commitStore

	lp := lpMocks.NewLogPoller(t)
	offRampReader, err := v1_0_0.NewOffRamp(logger.TestLogger(t), utils.RandomAddress(), nil, lp, nil)
	assert.NoError(t, err)
	p.offRampReader = offRampReader

	sendReqs := make([]ccipdata.Event[internal.EVM2EVMMessage], len(observations))
	sourceReader := ccipdatamocks.NewOnRampReader(t)
	for i := range observations {
		msg := internal.EVM2EVMMessage{
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
		}
		sendReqs[i] = ccipdata.Event[internal.EVM2EVMMessage]{Data: msg}
	}
	sourceReader.On("GetSendRequestsBetweenSeqNums",
		ctx, observations[0].SeqNr, observations[len(observations)-1].SeqNr, false).Return(sendReqs, nil)
	p.onRampReader = sourceReader

	execReport, err := p.buildReport(ctx, p.lggr, observations)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(execReport), MaxExecutionReportLength, "built execution report length")
}

func TestExecutionReportingPlugin_buildBatch(t *testing.T) {
	//_, _ := testhelpers.SetupChain(t)
	offRamp, _ := testhelpers.NewFakeOffRamp(t)
	// We do this just to have the parsing available.
	//onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress("0x1"), c)
	//require.NoError(t, err)
	lggr := logger.TestLogger(t)

	sender1 := common.HexToAddress("0xa")
	destNative := common.HexToAddress("0xb")
	srcNative := common.HexToAddress("0xc")

	msg1 := internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: internal.EVM2EVMMessage{
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
	msg4.TokenAmounts = []internal.TokenAmount{
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
		expectedSeqNrs           []ccip.ObservedMessage
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
			expectedSeqNrs:        []ccip.ObservedMessage{{SeqNr: uint64(1)}},
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
					messages:  []internal.EVM2EVMMessage{msg4.EVM2EVMMessage},
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
					EVM2EVMMessage: internal.EVM2EVMMessage{
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
					EVM2EVMMessage: internal.EVM2EVMMessage{
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
					EVM2EVMMessage: internal.EVM2EVMMessage{
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
			expectedSeqNrs:        []ccip.ObservedMessage{{SeqNr: uint64(10)}, {SeqNr: uint64(12)}},
		},
	}

	ctx := testutils.Context(t)

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			offRamp.SetSenderNonces(tc.offRampNoncesBySender)

			gasPriceEstimator := prices.NewMockGasPriceEstimatorExec(t)
			if tc.expectedSeqNrs != nil {
				gasPriceEstimator.On("EstimateMsgCostUSD", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(0), nil)
			}

			// Mock calls to reader.
			mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
			mockOffRampReader.On("GetSenderNonce", mock.Anything, sender1).Return(uint64(0), nil).Maybe()

			plugin := ExecutionReportingPlugin{
				tokenDataWorker:   tokendata.NewBackgroundWorker(ctx, map[common.Address]tokendata.Reader{}, 10, 5*time.Second, time.Hour),
				offRampReader:     mockOffRampReader,
				destWrappedNative: destNative,
				offchainConfig: ccipdata.ExecOffchainConfig{
					DestOptimisticConfirmations: 1,
					BatchGasLimit:               300_000,
					RelativeBoostPerWaitHour:    1,
				},
				lggr:              logger.TestLogger(t),
				gasPriceEstimator: gasPriceEstimator,
			}

			seqNrs := plugin.buildBatch(
				context.Background(),
				lggr,
				commitReportWithSendRequests{sendRequestsWithMeta: tc.reqs},
				tc.inflight,
				tc.tokenLimit,
				tc.srcPrices,
				tc.dstPrices,
				tc.destGasPrice,
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
		tokenAmounts            []internal.TokenAmount
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
			tokenAmounts: []internal.TokenAmount{
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
			tokenAmounts: []internal.TokenAmount{
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
			tokenAmounts: []internal.TokenAmount{
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
			tokenAmounts: []internal.TokenAmount{
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
		name         string
		tokenAmounts []internal.TokenAmount
		// the order of the following fields: sourceTokens, destTokens and poolRateLimits
		// should follow the order of the tokenAmounts
		sourceTokens      []common.Address
		destTokens        []common.Address
		destPools         []common.Address
		poolRateLimits    []ccipdata.TokenBucketRateLimit
		destPoolsCacheErr error

		expRateLimits map[common.Address]*big.Int
		expErr        bool
	}{
		{
			name: "happy flow",
			tokenAmounts: []internal.TokenAmount{
				{Token: tk1},
				{Token: tk2},
				{Token: tk1},
				{Token: tk1},
			},
			sourceTokens: []common.Address{tk1, tk2},
			destTokens:   []common.Address{tk1dest, tk2dest},
			destPools:    []common.Address{tk1pool, tk2pool},
			poolRateLimits: []ccipdata.TokenBucketRateLimit{
				{Tokens: big.NewInt(1000), IsEnabled: true},
				{Tokens: big.NewInt(2000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
				tk2dest: big.NewInt(2000),
			},
			expErr: false,
		},
		{
			name: "missing from source to dest mapping should not return error",
			tokenAmounts: []internal.TokenAmount{
				{Token: tk1},
				{Token: tk2}, // <- missing
			},
			sourceTokens: []common.Address{tk1},
			destTokens:   []common.Address{tk1dest},
			destPools:    []common.Address{tk1pool},
			poolRateLimits: []ccipdata.TokenBucketRateLimit{
				{Tokens: big.NewInt(1000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			expErr: false,
		},
		{
			name: "pool is disabled",
			tokenAmounts: []internal.TokenAmount{
				{Token: tk1},
				{Token: tk2},
			},
			sourceTokens: []common.Address{tk1, tk2},
			destTokens:   []common.Address{tk1dest, tk2dest},
			destPools:    []common.Address{tk1pool, tk2pool},
			poolRateLimits: []ccipdata.TokenBucketRateLimit{
				{Tokens: big.NewInt(1000), IsEnabled: true},
				{Tokens: big.NewInt(2000), IsEnabled: false},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			expErr: false,
		},
		{
			name: "dest pool cache error",
			tokenAmounts: []internal.TokenAmount{
				{Token: tk1},
			},
			sourceTokens: []common.Address{tk1},
			destTokens:   []common.Address{tk1dest},
			destPools:    []common.Address{tk1pool},
			poolRateLimits: []ccipdata.TokenBucketRateLimit{
				{Tokens: big.NewInt(1000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			destPoolsCacheErr: errors.New("some err"),
			expErr:            true,
		},
		{
			name: "pool for token not found",
			tokenAmounts: []internal.TokenAmount{
				{Token: tk1}, {Token: tk2}, {Token: tk1}, {Token: tk2},
			},
			sourceTokens: []common.Address{tk1, tk2},
			destTokens:   []common.Address{tk1dest, tk2dest},
			destPools:    []common.Address{tk1pool}, // <-- pool2 not found
			poolRateLimits: []ccipdata.TokenBucketRateLimit{
				{Tokens: big.NewInt(1000), IsEnabled: true},
			},
			expRateLimits: map[common.Address]*big.Int{
				tk1dest: big.NewInt(1000),
			},
			expErr: true,
		},
	}

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceToDestMapping := make(map[common.Address]common.Address)
			for i, srcTk := range tc.sourceTokens {
				sourceToDestMapping[srcTk] = tc.destTokens[i]
			}

			poolsMapping := make(map[common.Address]common.Address)
			for i, poolAddr := range tc.destPools {
				poolsMapping[tc.destTokens[i]] = poolAddr
			}

			p := &ExecutionReportingPlugin{}
			p.lggr = lggr

			offRampAddr := utils.RandomAddress()
			mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
			mockOffRampReader.On("Address").Return(offRampAddr, nil).Maybe()
			mockOffRampReader.On("GetTokens", ctx).Return(ccipdata.OffRampTokens{
				DestinationPool: poolsMapping,
			}, tc.destPoolsCacheErr).Maybe()
			mockOffRampReader.On("GetTokenPoolsRateLimits", ctx, tc.destPools).
				Return(tc.poolRateLimits, nil).
				Maybe()
			p.offRampReader = mockOffRampReader

			rateLimits, err := p.destPoolRateLimits(ctx, []commitReportWithSendRequests{
				{
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							EVM2EVMMessage: internal.EVM2EVMMessage{
								TokenAmounts: tc.tokenAmounts,
							},
						},
					},
				},
			}, sourceToDestMapping)

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRateLimits, rateLimits)
		})
	}
}

func TestExecutionReportingPlugin_getReportsWithSendRequests(t *testing.T) {
	testCases := []struct {
		name                string
		reports             []ccipdata.CommitStoreReport
		expQueryMin         uint64 // expected min/max used in the query to get ccipevents
		expQueryMax         uint64
		onchainEvents       []ccipdata.Event[internal.EVM2EVMMessage]
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
			reports: []ccipdata.CommitStoreReport{
				{
					Interval:   ccipdata.CommitStoreInterval{Min: 1, Max: 2},
					MerkleRoot: [32]byte{100},
				},
				{
					Interval:   ccipdata.CommitStoreInterval{Min: 3, Max: 3},
					MerkleRoot: [32]byte{200},
				},
			},
			expQueryMin: 1,
			expQueryMax: 3,
			onchainEvents: []ccipdata.Event[internal.EVM2EVMMessage]{
				{Data: internal.EVM2EVMMessage{SequenceNumber: 1}},
				{Data: internal.EVM2EVMMessage{SequenceNumber: 2}},
				{Data: internal.EVM2EVMMessage{SequenceNumber: 3}},
			},
			destExecutedSeqNums: []uint64{1},
			expReports: []commitReportWithSendRequests{
				{
					commitReport: ccipdata.CommitStoreReport{
						Interval:   ccipdata.CommitStoreInterval{Min: 1, Max: 2},
						MerkleRoot: [32]byte{100},
					},
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 1},
							Executed:       true,
							Finalized:      true,
						},
						{
							EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 2},
							Executed:       false,
							Finalized:      false,
						},
					},
				},
				{
					commitReport: ccipdata.CommitStoreReport{
						Interval:   ccipdata.CommitStoreInterval{Min: 3, Max: 3},
						MerkleRoot: [32]byte{200},
					},
					sendRequestsWithMeta: []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						{
							EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 3},
							Executed:       false,
							Finalized:      false,
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

			offRampReader := ccipdatamocks.NewOffRampReader(t)
			p.offRampReader = offRampReader

			sourceReader := ccipdatamocks.NewOnRampReader(t)
			sourceReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.expQueryMin, tc.expQueryMax, false).
				Return(tc.onchainEvents, nil).Maybe()
			p.onRampReader = sourceReader

			finalized := make(map[uint64]bool)
			for _, r := range tc.expReports {
				for _, s := range r.sendRequestsWithMeta {
					finalized[s.SequenceNumber] = s.Finalized
				}
			}

			var executedEvents []ccipdata.Event[ccipdata.ExecutionStateChanged]
			for _, executedSeqNum := range tc.destExecutedSeqNums {
				executedEvents = append(executedEvents, ccipdata.Event[ccipdata.ExecutionStateChanged]{
					Data: ccipdata.ExecutionStateChanged{
						SequenceNumber: executedSeqNum,
						Finalized:      finalized[executedSeqNum],
					},
				})
			}
			offRampReader.On("GetExecutionStateChangesBetweenSeqNums", ctx, tc.expQueryMin, tc.expQueryMax, 0).Return(executedEvents, nil).Maybe()

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

type delayedTokenDataWorker struct {
	delay time.Duration
	tokendata.Worker
}

func (m delayedTokenDataWorker) GetMsgTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([][]byte, error) {
	time.Sleep(m.delay)
	return nil, ctx.Err()
}

func TestExecutionReportingPlugin_getTokenDataWithCappedLatency(t *testing.T) {
	testCases := []struct {
		name               string
		allowedWaitingTime time.Duration
		workerLatency      time.Duration
		expErr             bool
	}{
		{
			name:               "happy flow",
			allowedWaitingTime: 10 * time.Millisecond,
			workerLatency:      time.Nanosecond,
			expErr:             false,
		},
		{
			name:               "worker takes long to reply",
			allowedWaitingTime: 10 * time.Millisecond,
			workerLatency:      20 * time.Millisecond,
			expErr:             true,
		},
	}

	ctx := testutils.Context(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ExecutionReportingPlugin{}
			p.tokenDataWorker = delayedTokenDataWorker{delay: tc.workerLatency}

			msg := internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: internal.EVM2EVMMessage{TokenAmounts: make([]internal.TokenAmount, 1)},
			}

			_, _, err := p.getTokenDataWithTimeout(ctx, msg, tc.allowedWaitingTime)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_calculateObservedMessagesConsensus(t *testing.T) {
	type args struct {
		observations []ccip.ExecutionObservation
		f            int
	}
	tests := []struct {
		name string
		args args
		want []ccip.ObservedMessage
	}{
		{
			name: "no observations",
			args: args{
				observations: nil,
				f:            0,
			},
			want: []ccip.ObservedMessage{},
		},
		{
			name: "common path",
			args: args{
				observations: []ccip.ExecutionObservation{
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
						},
					},
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0xff}}}, // different token data - should not be picked
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
							3: {TokenData: [][]byte{{0x3}, {0x3}, {0x3}}},
						},
					},
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
							2: {TokenData: [][]byte{{0x2}, {0x2}, {0x2}}},
						},
					},
				},
				f: 1,
			},
			want: []ccip.ObservedMessage{
				{SeqNr: 1, MsgData: ccip.MsgData{TokenData: [][]byte{{0x1}, {0x1}, {0x1}}}},
				{SeqNr: 2, MsgData: ccip.MsgData{TokenData: [][]byte{{0x2}, {0x2}, {0x2}}}},
			},
		},
		{
			name: "similar token data",
			args: args{
				observations: []ccip.ExecutionObservation{
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1}, {0x1}}},
						},
					},
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1, 0x1}}},
						},
					},
					{
						Messages: map[uint64]ccip.MsgData{
							1: {TokenData: [][]byte{{0x1}, {0x1, 0x1}}},
						},
					},
				},
				f: 1,
			},
			want: []ccip.ObservedMessage{
				{SeqNr: 1, MsgData: ccip.MsgData{TokenData: [][]byte{{0x1}, {0x1, 0x1}}}},
			},
		},
		{
			name: "results should be deterministic",
			args: args{
				observations: []ccip.ExecutionObservation{
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x2}}}}},
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x2}}}}},
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x1}}}}},
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x3}}}}},
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x3}}}}},
					{Messages: map[uint64]ccip.MsgData{1: {TokenData: [][]byte{{0x1}}}}},
				},
				f: 1,
			},
			want: []ccip.ObservedMessage{
				{SeqNr: 1, MsgData: ccip.MsgData{TokenData: [][]byte{{0x3}}}},
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
		retPrices []ccipdata.TokenPriceUpdate
		expPrices map[common.Address]*big.Int
		expErr    bool
	}{
		{
			name:      "base",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []ccipdata.TokenPriceUpdate{
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(20)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(30)}},
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
			retPrices: []ccipdata.TokenPriceUpdate{
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(20)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(30)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
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
			retPrices: []ccipdata.TokenPriceUpdate{
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(20)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(30)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(1000)}},
			},
			expErr: true,
		},
		{
			name:      "zero price should lead to an error",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []ccipdata.TokenPriceUpdate{
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(0)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(30)}},
			},
			expErr: true,
		},
		{
			name:      "contract returns less prices than requested",
			feeTokens: []common.Address{tk1, tk2},
			tokens:    []common.Address{tk3},
			retPrices: []ccipdata.TokenPriceUpdate{
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(10)}},
				{TokenPrice: ccipdata.TokenPrice{Value: big.NewInt(20)}},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceReg := ccipdatamocks.NewPriceRegistryReader(t)
			priceReg.On("GetTokenPrices", mock.Anything, mock.Anything).Return(tc.retPrices, nil)
			priceReg.On("Address").Return(utils.RandomAddress(), nil).Maybe()

			tokenPrices, err := getTokensPrices(context.Background(), priceReg, append(tc.feeTokens, tc.tokens...))
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			for tk, price := range tc.expPrices {
				assert.Equal(t, price, tokenPrices[tk])
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

		expInflightSeqNrs          mapset.Set[uint64]
		expInflightAggrVal         *big.Int
		expMaxInflightSenderNonces map[common.Address]uint64
		expInflightTokenAmounts    map[common.Address]*big.Int
		expErr                     bool
	}{
		{
			name: "base",
			inflight: []InflightInternalExecutionReport{
				{
					messages: []internal.EVM2EVMMessage{
						{
							Sender:         addrs[0],
							SequenceNumber: 100,
							Nonce:          2,
							TokenAmounts: []internal.TokenAmount{
								{Token: tokenAddrs[0], Amount: big.NewInt(1e18)},
								{Token: tokenAddrs[0], Amount: big.NewInt(2e18)},
							},
						},
						{
							Sender:         addrs[0],
							SequenceNumber: 106,
							Nonce:          4,
							TokenAmounts: []internal.TokenAmount{
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
			expInflightSeqNrs:  mapset.NewSet[uint64](100, 106),
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
					messages: []internal.EVM2EVMMessage{
						{
							Sender:         addrs[0],
							SequenceNumber: 100,
							Nonce:          2,
							TokenAmounts: []internal.TokenAmount{
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
			expInflightSeqNrs:          mapset.NewSet[uint64](),
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
			assert.True(t, tc.expInflightSeqNrs.Equal(inflightSeqNrs))
			assert.True(t, reflect.DeepEqual(tc.expInflightAggrVal, inflightAggrVal))
			assert.True(t, reflect.DeepEqual(tc.expMaxInflightSenderNonces, maxInflightSenderNonces))
			assert.True(t, reflect.DeepEqual(tc.expInflightTokenAmounts, inflightTokenAmounts))
		})
	}
}

func Test_commitReportWithSendRequests_validate(t *testing.T) {
	testCases := []struct {
		name           string
		reportInterval ccipdata.CommitStoreInterval
		numReqs        int
		expValid       bool
	}{
		{
			name:           "valid report",
			reportInterval: ccipdata.CommitStoreInterval{Min: 10, Max: 20},
			numReqs:        11,
			expValid:       true,
		},
		{
			name:           "report with one request",
			reportInterval: ccipdata.CommitStoreInterval{Min: 1234, Max: 1234},
			numReqs:        1,
			expValid:       true,
		},
		{
			name:           "request is missing",
			reportInterval: ccipdata.CommitStoreInterval{Min: 1234, Max: 1234},
			numReqs:        0,
			expValid:       false,
		},
		{
			name:           "requests are missing",
			reportInterval: ccipdata.CommitStoreInterval{Min: 1, Max: 10},
			numReqs:        5,
			expValid:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rep := commitReportWithSendRequests{
				commitReport: ccipdata.CommitStoreReport{
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
		report ccipdata.CommitStoreReport
		expRes bool
	}{
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 1},
			},
			report: ccipdata.CommitStoreReport{
				Interval: ccipdata.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: true,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 10},
			},
			report: ccipdata.CommitStoreReport{
				Interval: ccipdata.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: true,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 11},
			},
			report: ccipdata.CommitStoreReport{
				Interval: ccipdata.CommitStoreInterval{Min: 1, Max: 10},
			},
			expRes: false,
		},
		{
			name: "all requests executed and finalized",
			req: internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 10},
			},
			report: ccipdata.CommitStoreReport{
				Interval: ccipdata.CommitStoreInterval{Min: 10, Max: 10},
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
func generateExecutionReport(t *testing.T, numMsgs, tokensPerMsg, bytesPerMsg int) ccipdata.ExecReport {
	messages := make([]internal.EVM2EVMMessage, numMsgs)

	offChainTokenData := make([][][]byte, numMsgs)
	for i := range messages {
		tokenAmounts := make([]internal.TokenAmount, tokensPerMsg)
		for j := range tokenAmounts {
			tokenAmounts[j] = internal.TokenAmount{
				Token:  utils.RandomAddress(),
				Amount: big.NewInt(math.MaxInt64),
			}
		}

		messages[i] = internal.EVM2EVMMessage{
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

	return ccipdata.ExecReport{
		Messages:          messages,
		OffchainTokenData: offChainTokenData,
		Proofs:            make([][32]byte, numMsgs),
		ProofFlagBits:     big.NewInt(rand.Int64()),
	}
}

func Test_selectReportsToFillBatch(t *testing.T) {
	reports := []ccipdata.CommitStoreReport{
		{Interval: ccipdata.CommitStoreInterval{Min: 1, Max: 10}},
		{Interval: ccipdata.CommitStoreInterval{Min: 11, Max: 20}},
		{Interval: ccipdata.CommitStoreInterval{Min: 21, Max: 25}},
		{Interval: ccipdata.CommitStoreInterval{Min: 26, Max: math.MaxUint64}},
	}

	tests := []struct {
		name            string
		step            uint64
		numberOfBatches int
	}{
		{
			name:            "pick all at once when step size is high",
			step:            100,
			numberOfBatches: 1,
		},
		{
			name:            "pick one by one when step size is 1",
			step:            1,
			numberOfBatches: 4,
		},
		{
			name:            "pick two when step size doesn't match report",
			step:            15,
			numberOfBatches: 2,
		},
		{
			name:            "pick one by one when step size is smaller then reports",
			step:            4,
			numberOfBatches: 4,
		},
		{
			name:            "batch some reports together",
			step:            7,
			numberOfBatches: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unexpiredReportsBatches [][]ccipdata.CommitStoreReport
			for i := 0; i < len(reports); {
				unexpiredReports, step := selectReportsToFillBatch(reports[i:], tt.step)
				unexpiredReportsBatches = append(unexpiredReportsBatches, unexpiredReports)
				i += step
			}
			assert.Len(t, unexpiredReportsBatches, tt.numberOfBatches)

			var flatten []ccipdata.CommitStoreReport
			for _, r := range unexpiredReportsBatches {
				flatten = append(flatten, r...)
			}
			assert.Len(t, flatten, len(reports))
			assert.Equal(t, reports, flatten)
		})
	}
}

func Test_prepareTokenExecData(t *testing.T) {
	ctx := testutils.Context(t)

	weth := utils.RandomAddress()
	wavax := utils.RandomAddress()
	link := utils.RandomAddress()
	usdc := utils.RandomAddress()

	wethPriceUpdate := ccipdata.TokenPriceUpdate{TokenPrice: ccipdata.TokenPrice{Token: weth, Value: big.NewInt(2e18)}}
	wavaxPriceUpdate := ccipdata.TokenPriceUpdate{TokenPrice: ccipdata.TokenPrice{Token: wavax, Value: big.NewInt(3e18)}}
	linkPriceUpdate := ccipdata.TokenPriceUpdate{TokenPrice: ccipdata.TokenPrice{Token: link, Value: big.NewInt(4e18)}}
	usdcPriceUpdate := ccipdata.TokenPriceUpdate{TokenPrice: ccipdata.TokenPrice{Token: usdc, Value: big.NewInt(5e18)}}

	tokenPrices := map[common.Address]ccipdata.TokenPriceUpdate{weth: wethPriceUpdate, wavax: wavaxPriceUpdate, link: linkPriceUpdate, usdc: usdcPriceUpdate}

	tests := []struct {
		name               string
		sourceFeeTokens    []common.Address
		sourceFeeTokensErr error
		destTokens         []common.Address
		destTokensErr      error
		destFeeTokens      []common.Address
		destFeeTokensErr   error
		sourcePrices       []ccipdata.TokenPriceUpdate
		destPrices         []ccipdata.TokenPriceUpdate
	}{
		{
			name:         "only native token",
			sourcePrices: []ccipdata.TokenPriceUpdate{wethPriceUpdate},
			destPrices:   []ccipdata.TokenPriceUpdate{wavaxPriceUpdate},
		},
		{
			name:          "additional dest fee token",
			destFeeTokens: []common.Address{link},
			sourcePrices:  []ccipdata.TokenPriceUpdate{wethPriceUpdate},
			destPrices:    []ccipdata.TokenPriceUpdate{linkPriceUpdate, wavaxPriceUpdate},
		},
		{
			name:         "dest tokens",
			destTokens:   []common.Address{link, usdc},
			sourcePrices: []ccipdata.TokenPriceUpdate{wethPriceUpdate},
			destPrices:   []ccipdata.TokenPriceUpdate{linkPriceUpdate, usdcPriceUpdate, wavaxPriceUpdate},
		},
		{
			name:            "source fee tokens",
			sourceFeeTokens: []common.Address{usdc},
			sourcePrices:    []ccipdata.TokenPriceUpdate{usdcPriceUpdate, wethPriceUpdate},
			destPrices:      []ccipdata.TokenPriceUpdate{wavaxPriceUpdate},
		},
		{
			name:            "source, dest and fee tokens",
			sourceFeeTokens: []common.Address{usdc},
			destTokens:      []common.Address{link},
			destFeeTokens:   []common.Address{usdc},
			sourcePrices:    []ccipdata.TokenPriceUpdate{usdcPriceUpdate, wethPriceUpdate},
			destPrices:      []ccipdata.TokenPriceUpdate{usdcPriceUpdate, linkPriceUpdate, wavaxPriceUpdate},
		},
		{
			name:            "source, dest and fee tokens with duplicates",
			sourceFeeTokens: []common.Address{link, weth},
			destTokens:      []common.Address{link, wavax},
			destFeeTokens:   []common.Address{link, wavax},
			sourcePrices:    []ccipdata.TokenPriceUpdate{linkPriceUpdate, wethPriceUpdate},
			destPrices:      []ccipdata.TokenPriceUpdate{linkPriceUpdate, wavaxPriceUpdate},
		},
		{
			name:               "everything fails when source fails",
			sourceFeeTokensErr: errors.New("source error"),
		},
		{
			name:             "everything fails when dest fee fails",
			destFeeTokensErr: errors.New("dest fee error"),
		},
		{
			name:          "everything fails when dest  fails",
			destTokensErr: errors.New("dest error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offrampReader := ccipdatamocks.NewOffRampReader(t)
			sourcePriceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			gasPriceEstimator := prices.NewMockGasPriceEstimatorExec(t)

			offrampReader.On("CurrentRateLimiterState", ctx).Return(ccipdata.TokenBucketRateLimit{}, nil).Maybe()
			offrampReader.On("GetSourceToDestTokensMapping", ctx).Return(map[common.Address]common.Address{}, nil).Maybe()
			gasPriceEstimator.On("GetGasPrice", ctx).Return(prices.GasPrice(big.NewInt(1e9)), nil).Maybe()

			offrampReader.On("GetTokens", ctx).Return(ccipdata.OffRampTokens{DestinationTokens: tt.destTokens}, tt.destTokensErr).Maybe()
			sourcePriceRegistry.On("GetFeeTokens", ctx).Return(tt.sourceFeeTokens, tt.sourceFeeTokensErr).Maybe()
			sourcePriceRegistry.On("GetTokenPrices", ctx, mock.Anything).Return(tt.sourcePrices, nil).Maybe()
			destPriceRegistry.On("GetFeeTokens", ctx).Return(tt.destFeeTokens, tt.destFeeTokensErr).Maybe()
			destPriceRegistry.On("GetTokenPrices", ctx, mock.Anything).Return(tt.destPrices, nil).Maybe()

			reportingPlugin := ExecutionReportingPlugin{
				offRampReader:            offrampReader,
				sourcePriceRegistry:      sourcePriceRegistry,
				destPriceRegistry:        destPriceRegistry,
				gasPriceEstimator:        gasPriceEstimator,
				sourceWrappedNativeToken: weth,
				destWrappedNative:        wavax,
			}

			tokenData, err := reportingPlugin.prepareTokenExecData(ctx)
			if tt.destFeeTokensErr != nil || tt.sourceFeeTokensErr != nil || tt.destTokensErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, tokenData.sourceTokenPrices, len(tt.sourcePrices))
			assert.Len(t, tokenData.destTokenPrices, len(tt.destPrices))

			for token, price := range tokenData.sourceTokenPrices {
				assert.Equal(t, tokenPrices[token].Value, price)
			}

			for token, price := range tokenData.destTokenPrices {
				assert.Equal(t, tokenPrices[token].Value, price)
			}
		})
	}
}

func encodeExecutionReport(t *testing.T, report ccipdata.ExecReport) []byte {
	reader, err := v1_2_0.NewOffRamp(logger.TestLogger(t), utils.RandomAddress(), nil, nil, nil)
	require.NoError(t, err)
	encodedReport, err := reader.EncodeExecutionReport(report)
	require.NoError(t, err)
	return encodedReport
}
