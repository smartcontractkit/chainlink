package ccip

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
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
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
			unexpiredReports: []ccipdata.Event[ccipdata.CommitStoreReport]{
				{
					Data: ccipdata.CommitStoreReport{
						Interval:   ccipdata.CommitStoreInterval{Min: 10, Max: 12},
						MerkleRoot: [32]byte{123},
					},
				},
			},
			blessedRoots: map[[32]byte]bool{
				[32]byte{123}: true,
			},
			rateLimiterState: evm_2_evm_offramp.RateLimiterTokenBucket{
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

			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("IsDown", mock.Anything).Return(tc.commitStorePaused, nil)
			// Blessed roots return true
			for root, blessed := range tc.blessedRoots {
				commitStoreReader.On("IsBlessed", mock.Anything, root).Return(blessed, nil).Maybe()
			}
			commitStoreReader.On("GetAcceptedCommitReportsGteTimestamp", ctx, mock.Anything, 0).
				Return(tc.unexpiredReports, nil).Maybe()
			p.config.commitStoreReader = commitStoreReader

			destReader := ccipdatamocks.NewReader(t)
			destReader.On("LatestBlock", ctx).Return(logpoller.LogPollerBlock{BlockNumber: 1234}, nil).Maybe()
			p.config.destReader = destReader

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
			p.config.offRampReader = mockOffRampReader

			mockOnRampReader := ccipdatamocks.NewOnRampReader(t)
			mockOnRampReader.On("GetSendRequestsBetweenSeqNums", ctx, mock.Anything, mock.Anything, 0).
				Return(tc.sendRequests, nil).Maybe()
			p.config.onRampReader = mockOnRampReader

			cachedDestTokens := cache.NewMockAutoSync[cache.CachedTokens](t)
			cachedDestTokens.On("Get", ctx).Return(cache.CachedTokens{
				SupportedTokens: map[common.Address]common.Address{},
				FeeTokens:       []common.Address{},
			}, nil).Maybe()
			p.cachedDestTokens = cachedDestTokens

			destPriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceRegReader.On("GetTokenPrices", ctx, mock.Anything).Return(
				[]ccipdata.TokenPriceUpdate{{TokenPrice: ccipdata.TokenPrice{Token: common.HexToAddress("0x1"), Value: big.NewInt(123)}, TimestampUnixSec: big.NewInt(time.Now().Unix())}}, nil).Maybe()
			destPriceRegReader.On("Address").Return(utils.RandomAddress()).Maybe()
			sourcePriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
			sourcePriceRegReader.On("Address").Return(utils.RandomAddress()).Maybe()
			sourcePriceRegReader.On("GetTokenPrices", ctx, mock.Anything).Return(
				[]ccipdata.TokenPriceUpdate{{TokenPrice: ccipdata.TokenPrice{Token: common.HexToAddress("0x1"), Value: big.NewInt(123)}, TimestampUnixSec: big.NewInt(time.Now().Unix())}}, nil).Maybe()
			p.destPriceRegistry = destPriceRegReader
			p.config.sourcePriceRegistry = sourcePriceRegReader

			cachedTokenPools := cache.NewMockAutoSync[map[common.Address]common.Address](t)
			cachedTokenPools.On("Get", ctx).Return(tc.tokenPoolsMapping, nil).Maybe()
			p.cachedTokenPools = cachedTokenPools

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
		expectedReport      ccipdata.ExecReport
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

			//commitStore, _ := testhelpers.NewFakeCommitStore(t, tc.committedSeqNum)

			p.config.commitStoreReader = ccipdatamocks.NewCommitStoreReader(t)

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
	encodedReport, err := ccipdata.EncodeExecutionReport(report)
	require.NoError(t, err)

	mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
	mockOffRampReader.On("DecodeExecutionReport", encodedReport).Return(report, nil)

	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginStaticConfig{
			offRampReader: mockOffRampReader,
		},
		lggr:            logger.TestLogger(t),
		inflightReports: newInflightExecReportsContainer(models.MustMakeDuration(1 * time.Hour).Duration()),
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
	encodedReport, err := ccipdata.EncodeExecutionReport(report)
	require.NoError(t, err)

	mockCommitStoreReader := ccipdatamocks.NewCommitStoreReader(t)

	mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
	mockOffRampReader.On("DecodeExecutionReport", encodedReport).Return(report, nil)
	mockedExecState := mockOffRampReader.On("GetExecutionState", mock.Anything, uint64(12)).Return(uint8(ccipdata.ExecutionStateUntouched), nil).Once()

	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginStaticConfig{
			commitStoreReader: mockCommitStoreReader,
			offRampReader:     mockOffRampReader,
		},
		lggr:            logger.TestLogger(t),
		inflightReports: newInflightExecReportsContainer(models.MustMakeDuration(1 * time.Hour).Duration()),
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
	encodedReport, err := ccipdata.EncodeExecutionReport(executionReport)
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

	commitStore := ccipdatamocks.NewCommitStoreReader(t)
	commitStore.On("VerifyExecutionReport", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	commitStore.On("GetExpectedNextSequenceNumber", mock.Anything).
		Return(executionReport.Messages[len(executionReport.Messages)-1].SequenceNumber+1, nil)
	commitStore.On("GetAcceptedCommitReportsGteSeqNum", ctx, observations[0].SeqNr, 0).
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
	p.config.commitStoreReader = commitStore

	lp := lpMocks.NewLogPoller(t)
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	offRampReader, err := ccipdata.NewOffRampV1_0_0(logger.TestLogger(t), utils.RandomAddress(), nil, lp, nil)
	assert.NoError(t, err)
	p.config.offRampReader = offRampReader

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
		ctx, observations[0].SeqNr, observations[len(observations)-1].SeqNr, 0).Return(sendReqs, nil)
	p.config.onRampReader = sourceReader

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
			expectedSeqNrs:        []ObservedMessage{{SeqNr: uint64(10)}, {SeqNr: uint64(12)}},
		},
	}

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
				config: ExecutionPluginStaticConfig{
					offRampReader: mockOffRampReader,
				},
				destWrappedNative: destNative,
				offchainConfig: ccipdata.ExecOffchainConfig{
					SourceFinalityDepth:         5,
					DestOptimisticConfirmations: 1,
					DestFinalityDepth:           5,
					BatchGasLimit:               300_000,
					RelativeBoostPerWaitHour:    1,
					MaxGasPrice:                 1,
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
				func() (prices.GasPrice, error) { return tc.destGasPrice, nil },
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

			tokenPoolsCache := cache.NewMockAutoSync[map[common.Address]common.Address](t)
			tokenPoolsCache.On("Get", ctx).Return(poolsMapping, tc.destPoolsCacheErr).Maybe()
			p.cachedTokenPools = tokenPoolsCache

			offRampAddr := utils.RandomAddress()
			mockOffRampReader := ccipdatamocks.NewOffRampReader(t)
			mockOffRampReader.On("Address").Return(offRampAddr, nil).Maybe()
			mockOffRampReader.On("GetTokenPoolsRateLimits", ctx, tc.destPools).
				Return(tc.poolRateLimits, nil).
				Maybe()
			p.config.offRampReader = mockOffRampReader

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
			destLatestBlock:     10_000,
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
			p.config.offRampReader = offRampReader

			sourceReader := ccipdatamocks.NewOnRampReader(t)
			sourceReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.expQueryMin, tc.expQueryMax, 0).
				Return(tc.onchainEvents, nil).Maybe()
			p.config.onRampReader = sourceReader

			destReader := ccipdatamocks.NewReader(t)
			destReader.On("LatestBlock", ctx).Return(logpoller.LogPollerBlock{BlockNumber: tc.destLatestBlock}, nil).Maybe()
			var executedEvents []ccipdata.Event[ccipdata.ExecutionStateChanged]
			for _, executedSeqNum := range tc.destExecutedSeqNums {
				executedEvents = append(executedEvents, ccipdata.Event[ccipdata.ExecutionStateChanged]{
					Data: ccipdata.ExecutionStateChanged{SequenceNumber: executedSeqNum},
					Meta: ccipdata.Meta{BlockNumber: tc.destLatestBlock - 10},
				})
			}
			offRampReader.On("GetExecutionStateChangesBetweenSeqNums", ctx, tc.expQueryMin, tc.expQueryMax, 0).Return(executedEvents, nil).Maybe()
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

/*
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
	commitStoreReader, _ := testhelpers.NewFakeCommitStore(t, 1)
	offRamp, _ := testhelpers.NewFakeOffRamp(t)

	destPriceRegistryAddr := utils.RandomAddress()

	tokenDataProviders := make(map[common.Address]tokendata.Reader)

	rf := &ExecutionReportingPluginFactory{
		filtersMu:          &sync.Mutex{},
		sourceChainFilters: filters[:5],
		destChainFilters:   filters[5:10],
		config: ExecutionPluginStaticConfig{
			destLP:              destLP,
			sourceLP:            sourceLP,
			onRamp:              onRamp,
			commitStoreReader:         commitStoreReader,
			offRamp:             offRamp,
			sourcePriceRegistry: sourcePriceRegistry,
			tokenDataProviders:  tokenDataProviders,
		},
	}

	for _, f := range getExecutionPluginSourceLpChainFilters(sourcePriceRegistry.Address()) {
		sourceLP.On("RegisterFilter", f).Return(nil)
	}
	for _, f := range getExecutionPluginDestLpChainFilters(commitStoreReader.Address(), offRamp.Address(), destPriceRegistryAddr) {
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
*/

/*
func TestExecutionReportToEthTxMeta(t *testing.T) {
	t.Run("happy flow", func(t *testing.T) {
		executionReport := generateExecutionReport(t, 10, 3, 1000)
		encExecReport, err := ccipdata.EncodeExecutionReport(executionReport)
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
*/

/* this is a test related to the cache, should not be here
func TestUpdateSourceToDestTokenMapping(t *testing.T) {
	expectedNewBlockNumber := int64(10000)
	logs := []logpoller.Log{{BlockNumber: expectedNewBlockNumber}}
	mockDestLP := &lpMocks.LogPoller{}

	mockDestLP.On("LatestLogEventSigsAddrsWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(logs, nil)
	mockDestLP.On("LatestBlock", mock.Anything).Return(expectedNewBlockNumber, nil)

	sourceToken, destToken := common.HexToAddress("111111"), common.HexToAddress("222222")

	mockOffRamp := ccipdata.NewMockOffRampReader(t)
	mockOffRamp.On("Address").Return(common.HexToAddress("0x01"))
	mockOffRamp.On("GetSupportedTokens", mock.Anything).Return([]common.Address{sourceToken}, nil)
	mockOffRamp.On("GetDestinationToken", mock.Anything, sourceToken).Return(destToken, nil)

	mockPriceRegistry := ccipdata.NewMockPriceRegistryReader(t)
	mockPriceRegistry.On("Address").Return(common.HexToAddress("0x02"))
	mockPriceRegistry.On("GetFeeTokens", mock.Anything).Return([]common.Address{}, nil)

	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginStaticConfig{
			destLP:        mockDestLP,
			offRampReader: mockOffRamp,
		},
		cachedDestTokens: cache.NewCachedSupportedTokens(mockDestLP, mockOffRamp, mockPriceRegistry, 0),
	}

	value, err := plugin.cachedDestTokens.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, destToken, value.SupportedTokens[sourceToken])
}
*/

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
				{SeqNr: 1, MsgData: MsgData{TokenData: [][]byte{{0x3}}}},
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
			priceReg.On("Address").Return(utils.RandomAddress(), nil)

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
