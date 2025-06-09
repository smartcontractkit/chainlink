package v1_2_0

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	evm_2_evm_offramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib/rpclibmocks"
)

func TestOffRampGetDestinationTokensFromSourceTokens(t *testing.T) {
	ctx := testutils.Context(t)
	const numSrcTokens = 20

	testCases := []struct {
		name           string
		outputChangeFn func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr
		expErr         bool
	}{
		{
			name:           "happy path",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr { return outputs },
			expErr:         false,
		},
		{
			name: "rpc error",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[2].Err = errors.New("some error")
				return outputs
			},
			expErr: true,
		},
		{
			name: "unexpected outputs length should be fine if the type is correct",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = append(outputs[0].Outputs, "unexpected", 123)
				return outputs
			},
			expErr: false,
		},
		{
			name: "different compatible type",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = []any{outputs[0].Outputs[0].(common.Address)}
				return outputs
			},
			expErr: false,
		},
		{
			name: "different incompatible type",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = []any{outputs[0].Outputs[0].(common.Address).Bytes()}
				return outputs
			},
			expErr: true,
		},
	}

	lp := mocks.NewLogPoller(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchCaller := rpclibmocks.NewEvmBatchCaller(t)
			o := &OffRamp{evmBatchCaller: batchCaller, lp: lp}
			srcTks, dstTks, outputs := generateTokensAndOutputs(numSrcTokens)
			outputs = tc.outputChangeFn(outputs)
			batchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).Return(outputs, nil)
			genericAddrs := ccipcalc.EvmAddrsToGeneric(srcTks...)
			actualDstTokens, err := o.getDestinationTokensFromSourceTokens(ctx, genericAddrs)

			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, ccipcalc.EvmAddrsToGeneric(dstTks...), actualDstTokens)
		})
	}
}

func TestCachedOffRampTokens(t *testing.T) {
	// Test data.
	srcTks, dstTks, _ := generateTokensAndOutputs(3)

	// Mock contract wrapper.
	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("GetDestinationTokens", mock.Anything).Return(dstTks, nil)
	mockOffRamp.On("GetSupportedTokens", mock.Anything).Return(srcTks, nil)
	mockOffRamp.On("Address").Return(utils.RandomAddress())

	lp := mocks.NewLogPoller(t)
	lp.On("LatestBlock", mock.Anything).Return(logpoller.Block{BlockNumber: rand.Int63()}, nil)

	offRamp := OffRamp{
		offRampV120:    mockOffRamp,
		lp:             lp,
		Logger:         logger.Test(t),
		Client:         clienttest.NewClient(t),
		evmBatchCaller: rpclibmocks.NewEvmBatchCaller(t),
		cachedOffRampTokens: cache.NewLogpollerEventsBased[cciptypes.OffRampTokens](
			lp,
			offrampPoolAddedPoolRemovedEvents,
			mockOffRamp.Address(),
		),
	}

	ctx := testutils.Context(t)
	tokens, err := offRamp.GetTokens(ctx)
	require.NoError(t, err)

	// Verify data is properly loaded in the cache.
	expectedPools := make(map[cciptypes.Address]cciptypes.Address)
	for i := range dstTks {
		expectedPools[cciptypes.Address(dstTks[i].String())] = cciptypes.Address(dstTks[i].String())
	}
	require.Equal(t, cciptypes.OffRampTokens{
		DestinationTokens: ccipcalc.EvmAddrsToGeneric(dstTks...),
		SourceTokens:      ccipcalc.EvmAddrsToGeneric(srcTks...),
	}, tokens)
}

func generateTokensAndOutputs(nbTokens uint) ([]common.Address, []common.Address, []rpclib.DataAndErr) {
	srcTks := make([]common.Address, nbTokens)
	dstTks := make([]common.Address, nbTokens)
	outputs := make([]rpclib.DataAndErr, nbTokens)
	for i := range srcTks {
		srcTks[i] = utils.RandomAddress()
		dstTks[i] = utils.RandomAddress()
		outputs[i] = rpclib.DataAndErr{
			Outputs: []any{dstTks[i]}, Err: nil,
		}
	}
	return srcTks, dstTks, outputs
}

func Test_LogsAreProperlyMarkedAsFinalized(t *testing.T) {
	minSeqNr := uint64(10)
	maxSeqNr := uint64(14)
	inputLogs := []logpoller.Log{
		CreateExecutionStateChangeEventLog(t, 10, 2, utils.RandomBytes32()),
		CreateExecutionStateChangeEventLog(t, 11, 3, utils.RandomBytes32()),
		CreateExecutionStateChangeEventLog(t, 12, 5, utils.RandomBytes32()),
		CreateExecutionStateChangeEventLog(t, 14, 7, utils.RandomBytes32()),
	}

	tests := []struct {
		name                        string
		lastFinalizedBlock          int64
		expectedFinalizedSequenceNr []uint64
	}{
		{
			"all logs are finalized",
			10,
			[]uint64{10, 11, 12, 14},
		},
		{
			"some logs are finalized",
			5,
			[]uint64{10, 11, 12},
		},
		{
			"no logs are finalized",
			1,
			[]uint64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offrampAddress := utils.RandomAddress()

			lp := mocks.NewLogPoller(t)
			lp.On("LatestBlock", mock.Anything).
				Return(logpoller.Block{FinalizedBlockNumber: tt.lastFinalizedBlock}, nil)
			lp.On("IndexedLogsTopicRange", mock.Anything, ExecutionStateChangedEvent, offrampAddress, 1, logpoller.EvmWord(minSeqNr), logpoller.EvmWord(maxSeqNr), evmtypes.Confirmations(0)).
				Return(inputLogs, nil)

			offRamp, err := NewOffRamp(logger.Test(t), offrampAddress, clienttest.NewClient(t), lp, nil, nil, nil)
			require.NoError(t, err)
			logs, err := offRamp.GetExecutionStateChangesBetweenSeqNums(testutils.Context(t), minSeqNr, maxSeqNr, 0)
			require.NoError(t, err)
			assert.Len(t, logs, len(inputLogs))

			for _, log := range logs {
				assert.Equal(t, slices.Contains(tt.expectedFinalizedSequenceNr, log.SequenceNumber), log.IsFinalized())
			}
		})
	}
}

func TestGetRouter(t *testing.T) {
	routerAddr := utils.RandomAddress()

	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("GetDynamicConfig", mock.Anything).Return(evm_2_evm_offramp_1_2_0.EVM2EVMOffRampDynamicConfig{
		Router: routerAddr,
	}, nil)

	offRamp := OffRamp{
		offRampV120: mockOffRamp,
	}

	ctx := testutils.Context(t)
	gotRouterAddr, err := offRamp.GetRouter(ctx)
	require.NoError(t, err)

	gotRouterEvmAddr, err := ccipcalc.GenericAddrToEvm(gotRouterAddr)
	require.NoError(t, err)
	assert.Equal(t, routerAddr, gotRouterEvmAddr)
}

func CreateExecutionStateChangeEventLog(t *testing.T, seqNr uint64, blockNumber int64, messageID common.Hash) logpoller.Log {
	tAbi, err := evm_2_evm_offramp.EVM2EVMOffRampMetaData.GetAbi()
	require.NoError(t, err)
	eseEvent, ok := tAbi.Events["ExecutionStateChanged"]
	require.True(t, ok)

	logData, err := eseEvent.Inputs.NonIndexed().Pack(uint8(1), []byte("some return data"))
	require.NoError(t, err)
	seqNrBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqNrBytes, seqNr)
	seqNrTopic := common.BytesToHash(seqNrBytes)
	topic0 := evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}.Topic()

	return logpoller.Log{
		Topics: [][]byte{
			topic0[:],
			seqNrTopic[:],
			messageID[:],
		},
		Data:        logData,
		LogIndex:    1,
		BlockHash:   utils.RandomBytes32(),
		BlockNumber: blockNumber,
		EventSig:    topic0,
		Address:     testutils.NewAddress(),
		TxHash:      utils.RandomBytes32(),
	}
}
