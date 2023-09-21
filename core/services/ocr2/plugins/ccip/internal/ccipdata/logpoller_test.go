package ccipdata

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestLogPollerClient_loadDependency(t *testing.T) {
	c := &LogPollerReader{}

	someAddr := utils.RandomAddress()

	onRamp, err := c.loadOnRamp(someAddr)
	assert.NoError(t, err)
	onRamp2, err := c.loadOnRamp(someAddr)
	assert.NoError(t, err)
	// the objects should have the same pointer
	// since the second time the dependency should've been loaded from cache instead of initializing a new instance.
	assert.True(t, onRamp == onRamp2)

	offRamp, err := c.loadOffRamp(someAddr)
	assert.NoError(t, err)
	offRamp2, err := c.loadOffRamp(someAddr)
	assert.NoError(t, err)
	assert.True(t, offRamp == offRamp2)

	priceReg, err := c.loadPriceRegistry(someAddr)
	assert.NoError(t, err)
	priceReg2, err := c.loadPriceRegistry(someAddr)
	assert.NoError(t, err)
	assert.True(t, priceReg == priceReg2)
}

func Test_parseLogs(t *testing.T) {
	// generate 100 logs
	logs := make([]logpoller.Log, 100)
	for i := range logs {
		logs[i].LogIndex = int64(i + 1)
		logs[i].BlockNumber = int64(i) * 1000
		logs[i].BlockTimestamp = time.Now()
	}

	parseFn := func(log types.Log) (*uint, error) {
		// Simulate some random error
		if log.Index == 100 {
			return nil, fmt.Errorf("some error")
		}
		return &log.Index, nil
	}

	parsedEvents, err := parseLogs[uint](logs, logger.TestLogger(t), parseFn)
	assert.NoError(t, err)
	assert.Len(t, parsedEvents, 99)

	// Make sure everything is parsed according to the parse func
	for i, ev := range parsedEvents {
		assert.Equal(t, i+1, int(ev.Data))
		assert.Equal(t, int(i)*1000, int(ev.BlockNumber))
		assert.Greater(t, ev.BlockTimestamp, time.Now().Add(-time.Minute))
	}
}

func TestLogPollerClient_GetSendRequestsGteSeqNum(t *testing.T) {
	onRampAddr := utils.RandomAddress()
	seqNum := uint64(100)
	confs := 4

	t.Run("using confs", func(t *testing.T) {
		lp := mocks.NewLogPoller(t)
		lp.On("LogsDataWordGreaterThan",
			abihelpers.EventSignatures.SendRequested,
			onRampAddr,
			abihelpers.EventSignatures.SendRequestedSequenceNumberWord,
			abihelpers.EvmWord(seqNum),
			confs,
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		c := &LogPollerReader{lp: lp}
		events, err := c.GetSendRequestsGteSeqNum(
			context.Background(),
			onRampAddr,
			seqNum,
			false,
			confs,
		)
		assert.NoError(t, err)
		assert.Empty(t, events)
		lp.AssertExpectations(t)
	})

	t.Run("using latest confirmed block", func(t *testing.T) {
		h := &types.Header{Number: big.NewInt(100000)}

		lp := mocks.NewLogPoller(t)
		lp.On("LogsUntilBlockHashDataWordGreaterThan",
			abihelpers.EventSignatures.SendRequested,
			onRampAddr,
			abihelpers.EventSignatures.SendRequestedSequenceNumberWord,
			abihelpers.EvmWord(seqNum),
			h.Hash(),
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		cl := evmClientMocks.NewClient(t)
		cl.On("HeaderByNumber", mock.Anything, mock.Anything).Return(h, nil)

		c := &LogPollerReader{lp: lp, client: cl}
		events, err := c.GetSendRequestsGteSeqNum(
			context.Background(),
			onRampAddr,
			seqNum,
			true,
			confs,
		)
		assert.NoError(t, err)
		assert.Empty(t, events)
		lp.AssertExpectations(t)
		cl.AssertExpectations(t)
	})
}

func TestLogPollerClient_GetLastUSDCMessagePriorToLogIndexInTx(t *testing.T) {
	txHash := utils.RandomAddress().Hash()
	ccipLogIndex := int64(100)

	expectedData := []byte("-1")

	t.Run("multiple found", func(t *testing.T) {
		lp := mocks.NewLogPoller(t)
		lp.On("IndexedLogsByTxHash",
			abihelpers.EventSignatures.USDCMessageSent,
			txHash,
			mock.Anything,
		).Return([]logpoller.Log{
			{LogIndex: ccipLogIndex - 2, Data: []byte("-2")},
			{LogIndex: ccipLogIndex - 1, Data: expectedData},
			{LogIndex: ccipLogIndex, Data: []byte("0")},
			{LogIndex: ccipLogIndex + 1, Data: []byte("1")},
		}, nil)

		c := &LogPollerReader{lp: lp}
		usdcMessageData, err := c.GetLastUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, txHash)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, usdcMessageData)

		lp.AssertExpectations(t)
	})

	t.Run("none found", func(t *testing.T) {
		lp := mocks.NewLogPoller(t)
		lp.On("IndexedLogsByTxHash",
			abihelpers.EventSignatures.USDCMessageSent,
			txHash,
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		c := &LogPollerReader{lp: lp}
		usdcMessageData, err := c.GetLastUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, txHash)
		assert.Errorf(t, err, fmt.Sprintf("no USDC message found prior to log index %d in tx %s", ccipLogIndex, txHash.Hex()))
		assert.Nil(t, usdcMessageData)

		lp.AssertExpectations(t)
	})
}
