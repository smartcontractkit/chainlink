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

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipevents"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	plugintesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/plugins"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	MaxTokensPerMessage = 5
	MaxPayloadLength    = 100_000
)

type execTestHarness = struct {
	plugintesthelpers.CCIPPluginTestHarness
	plugin *ExecutionReportingPlugin
}

func setupExecTestHarness(t *testing.T) execTestHarness {
	th := plugintesthelpers.SetupCCIPTestHarness(t)

	lggr := logger.TestLogger(t)
	destFeeEstimator := mocks.NewEvmFeeEstimator(t)

	destFeeEstimator.On(
		"GetFee",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Return(gas.EvmFee{Legacy: assets.NewWei(defaultGasPrice)}, uint32(200e3), nil)

	offchainConfig := ccipconfig.ExecOffchainConfig{
		SourceFinalityDepth:         0,
		DestOptimisticConfirmations: 0,
		MaxGasPrice:                 200e9,
		BatchGasLimit:               5e6,
		RootSnoozeTime:              models.MustMakeDuration(10 * time.Minute),
		InflightCacheExpiry:         models.MustMakeDuration(3 * time.Minute),
		RelativeBoostPerWaitHour:    0.07,
	}
	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			lggr:                     th.Lggr,
			sourceLP:                 th.SourceLP,
			destLP:                   th.DestLP,
			sourceEvents:             ccipevents.NewLogPollerClient(th.SourceLP, lggr, th.SourceClient),
			destEvents:               ccipevents.NewLogPollerClient(th.DestLP, lggr, th.DestClient),
			sourcePriceRegistry:      th.Source.PriceRegistry,
			onRamp:                   th.Source.OnRamp,
			commitStore:              th.Dest.CommitStore,
			offRamp:                  th.Dest.OffRamp,
			destClient:               th.DestClient,
			sourceClient:             th.SourceClient,
			sourceWrappedNativeToken: th.Source.WrappedNative.Address(),
			leafHasher:               hashlib.NewLeafHasher(th.Source.ChainSelector, th.Dest.ChainSelector, th.Source.OnRamp.Address(), hashlib.NewKeccakCtx()),
			destGasEstimator:         destFeeEstimator,
		},
		onchainConfig:         th.ExecOnchainConfig,
		offchainConfig:        offchainConfig,
		lggr:                  th.Lggr.Named("ExecutionReportingPlugin"),
		snoozedRoots:          cache.NewSnoozedRoots(th.ExecOnchainConfig.PermissionLessExecutionThresholdDuration(), offchainConfig.RootSnoozeTime.Duration()),
		inflightReports:       newInflightExecReportsContainer(offchainConfig.InflightCacheExpiry.Duration()),
		destPriceRegistry:     th.Dest.PriceRegistry,
		destWrappedNative:     th.Dest.WrappedNative.Address(),
		cachedSourceFeeTokens: cache.NewCachedFeeTokens(th.SourceLP, th.Source.PriceRegistry, int64(offchainConfig.SourceFinalityDepth)),
		cachedDestTokens:      cache.NewCachedSupportedTokens(th.DestLP, th.Dest.OffRamp, th.Dest.PriceRegistry, int64(offchainConfig.DestOptimisticConfirmations)),
	}
	return execTestHarness{
		CCIPPluginTestHarness: th,
		plugin:                &plugin,
	}
}

func TestMaxExecutionReportSize(t *testing.T) {
	// Ensure that given max payload size and max num tokens,
	// Our report size is under the tx size limit.
	th := setupExecTestHarness(t)
	th.plugin.F = 1
	mb := th.GenerateAndSendMessageBatch(t, 50, MaxPayloadLength, MaxTokensPerMessage)

	// commit root
	encoded, err := abihelpers.EncodeCommitReport(commit_store.CommitStoreCommitReport{
		Interval:   mb.Interval,
		MerkleRoot: mb.Root,
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
			DestChainSelector: 0,
			UsdPerUnitGas:     big.NewInt(0),
		},
	})
	require.NoError(t, err)
	latestEpocAndRound, err := th.Dest.CommitStoreHelper.GetLatestPriceEpochAndRound(nil)
	require.NoError(t, err)
	_, err = th.Dest.CommitStoreHelper.Report(th.Dest.User, encoded, big.NewInt(int64(latestEpocAndRound+1)))
	require.NoError(t, err)
	// double commit to ensure enough confirmations
	th.CommitAndPollLogs(t)
	th.CommitAndPollLogs(t)

	fullReport, err := abihelpers.EncodeExecutionReport(evm_2_evm_offramp.InternalExecutionReport{
		Messages:          mb.Messages,
		OffchainTokenData: mb.TokenData,
		Proofs:            mb.Proof.Hashes,
		ProofFlagBits:     mb.ProofBits,
	})
	require.NoError(t, err)
	// ensure "naive" full report would be bigger than limit
	require.Greater(t, len(fullReport), MaxExecutionReportLength, "full execution report length")

	observations := make([]ObservedMessage, len(mb.Messages))
	for i, msg := range mb.Messages {
		observations[i] = NewObservedMessage(msg.SequenceNumber, mb.TokenData[i])
	}

	// buildReport should cap the built report to fit in MaxExecutionReportLength
	execReport, err := th.plugin.buildReport(testutils.Context(t), th.Lggr, observations)
	require.NoError(t, err)
	require.LessOrEqual(t, len(execReport), MaxExecutionReportLength, "built execution report length")
}

func TestExecutionReportToEthTxMetadata(t *testing.T) {
	c := plugintesthelpers.SetupCCIPTestHarness(t)
	tests := []struct {
		name     string
		msgBatch plugintesthelpers.MessageBatch
		err      error
	}{
		{
			"happy flow",
			c.GenerateAndSendMessageBatch(t, 5, MaxPayloadLength, MaxTokensPerMessage),
			nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encExecReport, err := abihelpers.EncodeExecutionReport(evm_2_evm_offramp.InternalExecutionReport{
				Messages:          tc.msgBatch.Messages,
				OffchainTokenData: tc.msgBatch.TokenData,
				Proofs:            tc.msgBatch.Proof.Hashes,
				ProofFlagBits:     tc.msgBatch.ProofBits,
			})
			require.NoError(t, err)
			txMeta, err := ExecutionReportToEthTxMeta(encExecReport)
			if tc.err != nil {
				require.Equal(t, tc.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, txMeta)
			require.Len(t, txMeta.MessageIDs, len(tc.msgBatch.Messages))
		})
	}
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

func TestExecObservation(t *testing.T) {
	th := setupExecTestHarness(t)
	th.plugin.F = 1
	mb := th.GenerateAndSendMessageBatch(t, 2, 10, 1)

	// commit root
	encoded, err := abihelpers.EncodeCommitReport(commit_store.CommitStoreCommitReport{
		Interval:   mb.Interval,
		MerkleRoot: mb.Root,
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
			DestChainSelector: 0,
			UsdPerUnitGas:     big.NewInt(0),
		},
	})
	require.NoError(t, err)
	latestEpocAndRound, err := th.Dest.CommitStoreHelper.GetLatestPriceEpochAndRound(nil)
	require.NoError(t, err)
	_, err = th.Dest.CommitStoreHelper.Report(th.Dest.User, encoded, big.NewInt(int64(latestEpocAndRound+1)))
	require.NoError(t, err)
	// double commit to ensure enough confirmations
	th.CommitAndPollLogs(t)
	th.CommitAndPollLogs(t)

	expectedObservations := NewExecutionObservation([]ObservedMessage{
		{SeqNr: 1, MsgData: MsgData{TokenData: [][]byte{{}}}},
		{SeqNr: 2, MsgData: MsgData{TokenData: [][]byte{{}}}},
	})
	tests := []struct {
		name            string
		commitStoreDown bool
		expected        *ExecutionObservation
		expectedError   bool
	}{
		{
			"base",
			false,
			&expectedObservations,
			false,
		},
		{
			"commitStore down",
			true,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.commitStoreDown && !isCommitStoreDownNow(testutils.Context(t), th.Lggr, th.Dest.CommitStore) {
				_, err := th.Dest.CommitStore.Pause(th.Dest.User)
				require.NoError(t, err)
				th.CommitAndPollLogs(t)
			} else if !tt.commitStoreDown && isCommitStoreDownNow(testutils.Context(t), th.Lggr, th.Dest.CommitStore) {
				_, err := th.Dest.CommitStore.Unpause(th.Dest.User)
				require.NoError(t, err)
				th.CommitAndPollLogs(t)
			}

			gotObs, err := th.plugin.Observation(testutils.Context(t), ocrtypes.ReportTimestamp{}, ocrtypes.Query{})

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			var decodedObservation *ExecutionObservation
			if gotObs != nil {
				decodedObservation = new(ExecutionObservation)
				err = json.Unmarshal(gotObs, decodedObservation)
				require.NoError(t, err)

			}
			assert.Equal(t, tt.expected, decodedObservation)
		})
	}
}

func TestExecReport(t *testing.T) {
	th := setupExecTestHarness(t)
	th.plugin.F = 1
	mb := th.GenerateAndSendMessageBatch(t, 2, 10, 1)

	// commit root
	encoded, err := abihelpers.EncodeCommitReport(commit_store.CommitStoreCommitReport{
		Interval:   mb.Interval,
		MerkleRoot: mb.Root,
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
			DestChainSelector: 0,
			UsdPerUnitGas:     big.NewInt(0),
		},
	})
	require.NoError(t, err)
	execReport := mb.ToExecutionReport()

	latestEpocAndRound, err := th.Dest.CommitStoreHelper.GetLatestPriceEpochAndRound(nil)
	require.NoError(t, err)
	_, err = th.Dest.CommitStoreHelper.Report(th.Dest.User, encoded, big.NewInt(int64(latestEpocAndRound+1)))
	require.NoError(t, err)
	// double commit to ensure enough confirmations
	th.CommitAndPollLogs(t)
	th.CommitAndPollLogs(t)

	tests := []struct {
		name                 string
		commitStoreDown      bool
		observations         [][]ObservedMessage
		expectedShouldReport bool
		expectedReport       *evm_2_evm_offramp.InternalExecutionReport
		expectedError        bool
	}{
		{
			"base",
			false,
			[][]ObservedMessage{
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}})},
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}})},
			},
			true,
			&execReport,
			false,
		},
		{
			"partial observation",
			false,
			[][]ObservedMessage{
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}})},
				{NewObservedMessage(1, [][]byte{{}})},
			},
			true,
			func() *evm_2_evm_offramp.InternalExecutionReport {
				mb2 := mb
				mb2.Messages = mb.Messages[:1]
				mb2.Messages = mb.Messages[:1]
				mb2.TokenData = mb.TokenData[:1]
				mb2.Interval = commit_store.CommitStoreInterval{Min: 1, Max: 1}
				mb2.Proof, err = mb2.Tree.Prove([]int{0})
				assert.NoError(t, err)
				mb2.ProofBits = abihelpers.ProofFlagsToBits(mb2.Proof.SourceFlags)
				report := mb2.ToExecutionReport()
				return &report
			}(),
			false,
		},
		{
			"empty",
			false,
			[][]ObservedMessage{
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}})},
				{},
			},
			false,
			nil,
			false,
		},
		{
			"unknown seqNr",
			false,
			[][]ObservedMessage{
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}}), NewObservedMessage(3, [][]byte{{}})},
				{NewObservedMessage(1, [][]byte{{}}), NewObservedMessage(2, [][]byte{{}}), NewObservedMessage(3, [][]byte{{}})},
			},
			false,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obs []ocrtypes.AttributedObservation
			for _, o := range tt.observations {
				encoded, err := NewExecutionObservation(o).Marshal()
				require.NoError(t, err)
				obs = append(obs, ocrtypes.AttributedObservation{Observation: encoded})
			}
			gotShouldReport, gotReport, err := th.plugin.Report(testutils.Context(t), ocrtypes.ReportTimestamp{}, ocrtypes.Query{}, obs)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedShouldReport, gotShouldReport)

			var encodedReport ocrtypes.Report
			if tt.expectedReport != nil {
				encodedReport, err = abihelpers.EncodeExecutionReport(*tt.expectedReport)
				require.NoError(t, err)
			}
			assert.Equal(t, encodedReport, gotReport)
		})
	}
}

func TestExecShouldAcceptFinalizedReport(t *testing.T) {
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

func TestExecShouldTransmitAcceptedReport(t *testing.T) {
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

func TestBuildBatch(t *testing.T) {
	c, _ := testhelpers.SetupChain(t)
	mockOffRamp := mock_contracts.EVM2EVMOffRampInterface{}
	// We do this just to have the parsing available.
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress("0x1"), c)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)

	sender1 := common.HexToAddress("0xa")
	destNative := common.HexToAddress("0xb")
	srcNative := common.HexToAddress("0xc")
	plugin := ExecutionReportingPlugin{
		config: ExecutionPluginConfig{
			offRamp: &mockOffRamp,
			// We use a real onRamp for parsing
			onRamp: onRamp,
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

	msg1 := evm2EVMOnRampCCIPSendRequestedWithMeta{
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
		blockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
	}

	msg2 := msg1
	msg2.executed = true

	msg3 := msg1
	msg3.executed = true
	msg3.finalized = true

	msg4 := msg1
	msg4.TokenAmounts = []evm_2_evm_offramp.ClientEVMTokenAmount{
		{Token: srcNative, Amount: big.NewInt(100)},
	}

	msg5 := msg4
	msg5.SequenceNumber = msg5.SequenceNumber + 1
	msg5.Nonce = msg5.Nonce + 1

	var tt = []struct {
		name                     string
		reqs                     []evm2EVMOnRampCCIPSendRequestedWithMeta
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg1},
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg2},
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg3},
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg2},
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg2},
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
			reqs:                  []evm2EVMOnRampCCIPSendRequestedWithMeta{msg4},
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
			reqs:         []evm2EVMOnRampCCIPSendRequestedWithMeta{msg4},
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
			reqs: []evm2EVMOnRampCCIPSendRequestedWithMeta{msg5},
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
			reqs: []evm2EVMOnRampCCIPSendRequestedWithMeta{
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
					blockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
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
					blockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
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
					blockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
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
			for sender, nonce := range tc.offRampNoncesBySender {
				mockOffRamp.On("GetSenderNonce", mock.Anything, sender).Return(nonce, nil)
			}

			seqNrs := plugin.buildBatch(
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
			priceReg := mock_contracts.NewPriceRegistryInterface(t)

			priceReg.On("GetTokenPrices", mock.Anything, append(tc.feeTokens, tc.tokens...)).
				Return(tc.retPrices, nil)
			priceReg.On("Address").Return(common.HexToAddress("1234"), nil)

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

	onRamp := mock_contracts.NewEVM2EVMOnRampInterface(t)
	onRamp.On("Address").Return(utils.RandomAddress(), nil)

	sourcePriceRegistry := mock_contracts.NewPriceRegistryInterface(t)
	sourcePriceRegistry.On("Address").Return(utils.RandomAddress(), nil)

	commitStore := mock_contracts.NewCommitStoreInterface(t)
	commitStore.On("Address").Return(utils.RandomAddress(), nil)

	offRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	offRamp.On("Address").Return(utils.RandomAddress(), nil)

	destPriceRegistryAddr := utils.RandomAddress()

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
		},
	}

	for _, f := range getExecutionPluginSourceLpChainFilters(onRamp.Address(), sourcePriceRegistry.Address()) {
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
