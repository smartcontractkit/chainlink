package ccipexec

import (
	"bytes"
	"context"
	"encoding/binary"
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	mockstatuschecker "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker/mocks"
)

type testCase struct {
	name                                             string
	reqs                                             []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
	inflight                                         []InflightInternalExecutionReport
	tokenLimit, destGasPrice, inflightAggregateValue *big.Int
	srcPrices, dstPrices                             map[cciptypes.Address]*big.Int
	offRampNoncesBySender                            map[cciptypes.Address]uint64
	srcToDestTokens                                  map[cciptypes.Address]cciptypes.Address
	expectedSeqNrs                                   []ccip.ObservedMessage
	expectedStates                                   []messageExecStatus
	statuschecker                                    func(m *mockstatuschecker.CCIPTransactionStatusChecker)
	skipGasPriceEstimator                            bool
}

func Test_NewBatchingStrategy(t *testing.T) {
	t.Parallel()

	mockStatusChecker := mockstatuschecker.NewCCIPTransactionStatusChecker(t)

	testCases := []int{0, 1, 2}

	for _, batchingStrategyId := range testCases {
		factory, err := NewBatchingStrategy(uint32(batchingStrategyId), mockStatusChecker)
		if batchingStrategyId == 2 {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, factory)
			assert.NoError(t, err)
		}
	}
}

func Test_validateSendRequests(t *testing.T) {
	testCases := []struct {
		name             string
		seqNums          []uint64
		providedInterval cciptypes.CommitStoreInterval
		expErr           bool
	}{
		{
			name:             "zero interval no seq nums",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
			expErr:           true,
		},
		{
			name:             "exp 1 seq num got none",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           true,
		},
		{
			name:             "exp 10 seq num got none",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 10},
			expErr:           true,
		},
		{
			name:             "got 1 seq num as expected",
			seqNums:          []uint64{1},
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           false,
		},
		{
			name:             "got 5 seq num as expected",
			seqNums:          []uint64{11, 12, 13, 14, 15},
			providedInterval: cciptypes.CommitStoreInterval{Min: 11, Max: 15},
			expErr:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sendReqs := make([]cciptypes.EVM2EVMMessageWithTxMeta, 0, len(tc.seqNums))
			for _, seqNum := range tc.seqNums {
				sendReqs = append(sendReqs, cciptypes.EVM2EVMMessageWithTxMeta{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: seqNum},
				})
			}
			err := validateSendRequests(sendReqs, tc.providedInterval)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

type delayedTokenDataWorker struct {
	delay time.Duration
	tokendata.Worker
}

func (m delayedTokenDataWorker) GetMsgTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([][]byte, error) {
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
			tokenDataWorker := delayedTokenDataWorker{delay: tc.workerLatency}

			msg := cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{TokenAmounts: make([]cciptypes.TokenAmount, 1)},
			}

			_, _, err := getTokenDataWithTimeout(ctx, msg, tc.allowedWaitingTime, tokenDataWorker)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestBatchingStrategies(t *testing.T) {
	sender1 := ccipcalc.HexToAddress("0xa")
	destNative := ccipcalc.HexToAddress("0xb")
	srcNative := ccipcalc.HexToAddress("0xc")

	msg1 := createTestMessage(1, sender1, 1, srcNative, big.NewInt(1e9), false, nil)

	msg2 := msg1
	msg2.Executed = true

	msg3 := msg1
	msg3.Executed = true
	msg3.Finalized = true

	msg4 := msg1
	msg4.TokenAmounts = []cciptypes.TokenAmount{
		{Token: srcNative, Amount: big.NewInt(100)},
	}

	msg5 := msg4
	msg5.SequenceNumber = msg5.SequenceNumber + 1
	msg5.Nonce = msg5.Nonce + 1

	zkMsg1 := createTestMessage(1, sender1, 0, srcNative, big.NewInt(1e9), false, nil)
	zkMsg2 := createTestMessage(2, sender1, 0, srcNative, big.NewInt(1e9), false, nil)
	zkMsg3 := createTestMessage(3, sender1, 0, srcNative, big.NewInt(1e9), false, nil)
	zkMsg4 := createTestMessage(4, sender1, 0, srcNative, big.NewInt(1e9), false, nil)

	testCases := []testCase{
		{
			name:                   "single message no tokens",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(1)}},
			expectedStates:         []messageExecStatus{newMessageExecState(msg1.SequenceNumber, msg1.MessageID, AddedToBatch)},
		},
		{
			name:                   "gasPriceEstimator returns error",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
		},
		{
			name:                   "executed non finalized messages should be skipped",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates:         []messageExecStatus{newMessageExecState(msg2.SequenceNumber, msg2.MessageID, AlreadyExecuted)},
			skipGasPriceEstimator:  true,
		},
		{
			name:                   "finalized executed log",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg3},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates:         []messageExecStatus{newMessageExecState(msg3.SequenceNumber, msg3.MessageID, AlreadyExecuted)},
			skipGasPriceEstimator:  true,
		},
		{
			name:                   "dst token price does not exist",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates:         []messageExecStatus{newMessageExecState(msg1.SequenceNumber, msg1.MessageID, TokenNotInDestTokenPrices)},
			skipGasPriceEstimator:  true,
		},
		{
			name:                   "src token price does not exist",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates:         []messageExecStatus{newMessageExecState(msg1.SequenceNumber, msg1.MessageID, TokenNotInSrcTokenPrices)},
			skipGasPriceEstimator:  true,
		},
		{
			name:                   "message with tokens is not executed if limit is reached",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg4},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(99),
			destGasPrice:           big.NewInt(1),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1e18)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1e18)},
			srcToDestTokens: map[cciptypes.Address]cciptypes.Address{
				srcNative: destNative,
			},
			offRampNoncesBySender: map[cciptypes.Address]uint64{sender1: 0},
			expectedStates:        []messageExecStatus{newMessageExecState(msg4.SequenceNumber, msg4.MessageID, AggregateTokenLimitExceeded)},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "message with tokens is not executed if limit is reached when inflight is full",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg5},
			inflight:               []InflightInternalExecutionReport{{createdAt: time.Now(), messages: []cciptypes.EVM2EVMMessage{msg4.EVM2EVMMessage}}},
			inflightAggregateValue: big.NewInt(100),
			tokenLimit:             big.NewInt(50),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1e18)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1e18)},
			srcToDestTokens: map[cciptypes.Address]cciptypes.Address{
				srcNative: destNative,
			},
			offRampNoncesBySender: map[cciptypes.Address]uint64{sender1: 1},
			expectedStates:        []messageExecStatus{newMessageExecState(msg5.SequenceNumber, msg5.MessageID, AggregateTokenLimitExceeded)},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "skip when nonce doesn't match chain value",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 123},
			expectedStates:         []messageExecStatus{newMessageExecState(msg1.SequenceNumber, msg1.MessageID, InvalidNonce)},
			skipGasPriceEstimator:  true,
		},
		{
			name:                   "skip when nonce not found",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{msg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{},
			expectedStates:         []messageExecStatus{newMessageExecState(msg1.SequenceNumber, msg1.MessageID, MissingNonce)},
			skipGasPriceEstimator:  true,
		},
		{
			name: "unordered messages",
			reqs: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          0,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(10)}},
			expectedStates: []messageExecStatus{
				newMessageExecState(10, [32]byte{}, AddedToBatch),
			},
		},
		{
			name: "unordered messages not blocked by nonce",
			reqs: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 9,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          5,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          0,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 3},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(10)}},
			expectedStates: []messageExecStatus{
				newMessageExecState(9, [32]byte{}, InvalidNonce),
				newMessageExecState(10, [32]byte{}, AddedToBatch),
			},
		},
	}

	bestEffortTestCases := []testCase{
		{
			name: "skip when batch gas limit is reached",
			reqs: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          1,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 11,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          2,
						GasLimit:       big.NewInt(math.MaxInt64),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 12,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          3,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(10)}},
			expectedStates: []messageExecStatus{
				newMessageExecState(10, [32]byte{}, AddedToBatch),
				newMessageExecState(11, [32]byte{}, InsufficientRemainingBatchGas),
				newMessageExecState(12, [32]byte{}, InvalidNonce),
			},
		},
		{
			name: "some messages skipped after hitting max batch data len",
			reqs: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          1,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 11,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          2,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, MaxDataLenPerBatch-500), // skipped from batch
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 12,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          3,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(10)}},
			expectedStates: []messageExecStatus{
				newMessageExecState(10, [32]byte{}, AddedToBatch),
				newMessageExecState(11, [32]byte{}, InsufficientRemainingBatchDataLength),
				newMessageExecState(12, [32]byte{}, InvalidNonce),
			},
		},
		{
			name: "unordered messages then ordered messages",
			reqs: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 9,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          0,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 10,
						FeeTokenAmount: big.NewInt(1e9),
						Sender:         sender1,
						Nonce:          5,
						GasLimit:       big.NewInt(1),
						Data:           bytes.Repeat([]byte{'a'}, 1000),
						FeeToken:       srcNative,
						MessageID:      [32]byte{},
					},
					BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
				},
			},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 4},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: uint64(9)}, {SeqNr: uint64(10)}},
			expectedStates: []messageExecStatus{
				newMessageExecState(9, [32]byte{}, AddedToBatch),
				newMessageExecState(10, [32]byte{}, AddedToBatch),
			},
		},
	}

	specificZkOverflowTestCases := []testCase{
		{
			name:                   "batch size is 1",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2, zkMsg3},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg1.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, mock.Anything).Return([]types.TransactionStatus{}, -1, nil)
			},
		},
		{
			name:                   "snooze fatal message and return empty batch",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMFatalStatus),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Fatal}, 0, nil)
			},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "snooze fatal message and add next message to batch",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg2.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMFatalStatus),
				newMessageExecState(zkMsg2.SequenceNumber, zkMsg2.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Fatal}, 0, nil)
				m.On("CheckMessageStatus", mock.Anything, zkMsg2.MessageID.String()).Return([]types.TransactionStatus{}, -1, nil)
			},
		},
		{
			name:                   "all messages are fatal and batch is empty",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMFatalStatus),
				newMessageExecState(zkMsg2.SequenceNumber, zkMsg2.MessageID, TXMFatalStatus),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Fatal}, 0, nil)
				m.On("CheckMessageStatus", mock.Anything, zkMsg2.MessageID.String()).Return([]types.TransactionStatus{types.Fatal}, 0, nil)
			},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "message batched when unconfirmed or failed",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg1.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Unconfirmed, types.Failed}, 1, nil)
			},
		},
		{
			name:                   "message snoozed when multiple statuses with fatal",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg2.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMFatalStatus),
				newMessageExecState(zkMsg2.SequenceNumber, zkMsg2.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Unconfirmed, types.Failed, types.Fatal}, 2, nil)
				m.On("CheckMessageStatus", mock.Anything, zkMsg2.MessageID.String()).Return([]types.TransactionStatus{}, -1, nil)
			},
		},
		{
			name:                   "txm return error for message",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg2.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMCheckError),
				newMessageExecState(zkMsg2.SequenceNumber, zkMsg2.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{}, -1, errors.New("dummy txm error"))
				m.On("CheckMessageStatus", mock.Anything, zkMsg2.MessageID.String()).Return([]types.TransactionStatus{}, -1, nil)
			},
		},
		{
			name:                   "snooze message when inflight",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1},
			inflight:               createInflight(zkMsg1),
			inflightAggregateValue: zkMsg1.FeeTokenAmount,
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, SkippedInflight),
			},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "snooze when not inflight but txm returns error",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMCheckError),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{}, -1, errors.New("dummy txm error"))
			},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "snooze when not inflight but txm returns fatal status",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1},
			inflight:               []InflightInternalExecutionReport{},
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, TXMFatalStatus),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg1.MessageID.String()).Return([]types.TransactionStatus{types.Unconfirmed, types.Failed, types.Fatal}, 2, nil)
			},
			skipGasPriceEstimator: true,
		},
		{
			name:                   "snooze messages when inflight but batch valid messages",
			reqs:                   []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{zkMsg1, zkMsg2, zkMsg3, zkMsg4},
			inflight:               createInflight(zkMsg1, zkMsg2),
			inflightAggregateValue: big.NewInt(0),
			tokenLimit:             big.NewInt(0),
			destGasPrice:           big.NewInt(10),
			srcPrices:              map[cciptypes.Address]*big.Int{srcNative: big.NewInt(1)},
			dstPrices:              map[cciptypes.Address]*big.Int{destNative: big.NewInt(1)},
			offRampNoncesBySender:  map[cciptypes.Address]uint64{sender1: 0},
			expectedSeqNrs:         []ccip.ObservedMessage{{SeqNr: zkMsg3.SequenceNumber}},
			expectedStates: []messageExecStatus{
				newMessageExecState(zkMsg1.SequenceNumber, zkMsg1.MessageID, SkippedInflight),
				newMessageExecState(zkMsg2.SequenceNumber, zkMsg2.MessageID, SkippedInflight),
				newMessageExecState(zkMsg3.SequenceNumber, zkMsg3.MessageID, AddedToBatch),
			},
			statuschecker: func(m *mockstatuschecker.CCIPTransactionStatusChecker) {
				m.Mock = mock.Mock{} // reset mock
				m.On("CheckMessageStatus", mock.Anything, zkMsg3.MessageID.String()).Return([]types.TransactionStatus{}, -1, nil)
			},
			skipGasPriceEstimator: false,
		},
	}

	t.Run("BestEffortBatchingStrategy", func(t *testing.T) {
		strategy := &BestEffortBatchingStrategy{}
		runBatchingStrategyTests(t, strategy, 1_000_000, append(testCases, bestEffortTestCases...))
	})

	t.Run("ZKOverflowBatchingStrategy", func(t *testing.T) {
		mockedStatusChecker := mockstatuschecker.NewCCIPTransactionStatusChecker(t)
		strategy := &ZKOverflowBatchingStrategy{
			statuschecker: mockedStatusChecker,
		}
		runBatchingStrategyTests(t, strategy, 1_000_000, append(testCases, specificZkOverflowTestCases...))
	})
}

// Function to set up and run tests for a given batching strategy
func runBatchingStrategyTests(t *testing.T, strategy BatchingStrategy, availableGas uint64, testCases []testCase) {
	destNative := ccipcalc.HexToAddress("0xb")

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)

			gasPriceEstimator := prices.NewMockGasPriceEstimatorExec(t)
			if !tc.skipGasPriceEstimator {
				if tc.expectedSeqNrs != nil {
					gasPriceEstimator.On("EstimateMsgCostUSD", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(0), nil)
				} else {
					gasPriceEstimator.On("EstimateMsgCostUSD", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(0), errors.New("error"))
				}
			}

			// default case for ZKOverflowBatchingStrategy
			if strategyType := reflect.TypeOf(strategy); tc.statuschecker == nil && strategyType == reflect.TypeOf(&ZKOverflowBatchingStrategy{}) {
				strategy.(*ZKOverflowBatchingStrategy).statuschecker.(*mockstatuschecker.CCIPTransactionStatusChecker).On("CheckMessageStatus", mock.Anything, mock.Anything).Return([]types.TransactionStatus{}, -1, nil)
			}

			// Mock calls to TXM
			if tc.statuschecker != nil {
				tc.statuschecker(strategy.(*ZKOverflowBatchingStrategy).statuschecker.(*mockstatuschecker.CCIPTransactionStatusChecker))
			}

			batchContext := &BatchContext{
				report:                     commitReportWithSendRequests{sendRequestsWithMeta: tc.reqs},
				inflight:                   tc.inflight,
				inflightAggregateValue:     tc.inflightAggregateValue,
				lggr:                       lggr,
				availableDataLen:           MaxDataLenPerBatch,
				availableGas:               availableGas,
				expectedNonces:             make(map[cciptypes.Address]uint64),
				sendersNonce:               tc.offRampNoncesBySender,
				sourceTokenPricesUSD:       tc.srcPrices,
				destTokenPricesUSD:         tc.dstPrices,
				gasPrice:                   tc.destGasPrice,
				sourceToDestToken:          tc.srcToDestTokens,
				aggregateTokenLimit:        tc.tokenLimit,
				tokenDataRemainingDuration: 5 * time.Second,
				tokenDataWorker:            tokendata.NewBackgroundWorker(map[cciptypes.Address]tokendata.Reader{}, 10, 5*time.Second, time.Hour),
				gasPriceEstimator:          gasPriceEstimator,
				destWrappedNative:          destNative,
				offchainConfig: cciptypes.ExecOffchainConfig{
					DestOptimisticConfirmations: 1,
					BatchGasLimit:               300_000,
					RelativeBoostPerWaitHour:    1,
				},
			}

			seqNrs, execStates := strategy.BuildBatch(context.Background(), batchContext)

			runAssertions(t, tc, seqNrs, execStates, strategy)
		})
	}
}

// Utility function to run common assertions
func runAssertions(t *testing.T, tc testCase, seqNrs []ccip.ObservedMessage, execStates []messageExecStatus, strategy BatchingStrategy) {
	if tc.expectedSeqNrs == nil {
		assert.Len(t, seqNrs, 0)
	} else {
		assert.Equal(t, tc.expectedSeqNrs, seqNrs)
	}

	if tc.expectedStates == nil {
		assert.Len(t, execStates, 0)
	} else {
		assert.Equal(t, tc.expectedStates, execStates)
	}

	batchingStratID := strategy.GetBatchingStrategyID()
	if strategyType := reflect.TypeOf(strategy); strategyType == reflect.TypeOf(&BestEffortBatchingStrategy{}) {
		assert.Equal(t, batchingStratID, uint32(0))
	} else {
		assert.Equal(t, batchingStratID, uint32(1))
	}
}

func createTestMessage(seqNr uint64, sender cciptypes.Address, nonce uint64, feeToken cciptypes.Address, feeAmount *big.Int, executed bool, data []byte) cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta {
	return cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: cciptypes.EVM2EVMMessage{
			SequenceNumber: seqNr,
			FeeTokenAmount: feeAmount,
			Sender:         sender,
			Nonce:          nonce,
			GasLimit:       big.NewInt(1),
			Strict:         false,
			Receiver:       "",
			Data:           data,
			TokenAmounts:   nil,
			FeeToken:       feeToken,
			MessageID:      generateMessageIDFromInt(seqNr),
		},
		BlockTimestamp: time.Date(2010, 1, 1, 12, 12, 12, 0, time.UTC),
		Executed:       executed,
	}
}

func generateMessageIDFromInt(input uint64) [32]byte {
	var messageID [32]byte
	binary.LittleEndian.PutUint32(messageID[:], uint32(input))
	return messageID
}

func createInflight(msgs ...cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) []InflightInternalExecutionReport {
	reports := make([]InflightInternalExecutionReport, len(msgs))

	for i, msg := range msgs {
		reports[i] = InflightInternalExecutionReport{
			messages:  []cciptypes.EVM2EVMMessage{msg.EVM2EVMMessage},
			createdAt: msg.BlockTimestamp,
		}
	}

	return reports
}
