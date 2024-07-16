package ccipdata

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogPollerClient_GetUSDCMessagePriorToLogIndexInTx(t *testing.T) {
	addr := utils.RandomAddress()
	txHash := common.BytesToHash(addr[:])
	ccipLogIndex := int64(100)

	expectedData := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f80000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc08610000000000000000"
	expectedPostParse := "0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861"
	lggr := logger.TestLogger(t)

	t.Run("multiple found - selected last", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)

		lp.On("IndexedLogsByTxHash",
			mock.Anything,
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
		).Return([]logpoller.Log{
			{LogIndex: ccipLogIndex - 2, Data: []byte("-2")},
			{LogIndex: ccipLogIndex - 1, Data: hexutil.MustDecode(expectedData)},
			{LogIndex: ccipLogIndex, Data: []byte("0")},
			{LogIndex: ccipLogIndex + 1, Data: []byte("1")},
		}, nil)
		usdcMessageData, err := u.GetUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, 0, txHash.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedPostParse, hexutil.Encode(usdcMessageData))
		lp.AssertExpectations(t)
	})

	t.Run("multiple found - selected first", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)

		lp.On("IndexedLogsByTxHash",
			mock.Anything,
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
		).Return([]logpoller.Log{
			{LogIndex: ccipLogIndex - 2, Data: hexutil.MustDecode(expectedData)},
			{LogIndex: ccipLogIndex - 1, Data: []byte("-2")},
			{LogIndex: ccipLogIndex, Data: []byte("0")},
			{LogIndex: ccipLogIndex + 1, Data: []byte("1")},
		}, nil)
		usdcMessageData, err := u.GetUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, 1, txHash.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedPostParse, hexutil.Encode(usdcMessageData))
		lp.AssertExpectations(t)
	})

	t.Run("logs fetched from memory in subsequent calls", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)

		lp.On("IndexedLogsByTxHash",
			mock.Anything,
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
		).Return([]logpoller.Log{
			{LogIndex: ccipLogIndex - 2, Data: hexutil.MustDecode(expectedData)},
			{LogIndex: ccipLogIndex - 1, Data: []byte("-2")},
			{LogIndex: ccipLogIndex, Data: []byte("0")},
			{LogIndex: ccipLogIndex + 1, Data: []byte("1")},
		}, nil).Once()

		// first call logs must be fetched from lp
		usdcMessageData, err := u.GetUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, 1, txHash.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedPostParse, hexutil.Encode(usdcMessageData))

		// subsequent call, logs must be fetched from memory
		usdcMessageData, err = u.GetUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, 1, txHash.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedPostParse, hexutil.Encode(usdcMessageData))

		lp.AssertExpectations(t)
	})

	t.Run("none found", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)
		lp.On("IndexedLogsByTxHash",
			mock.Anything,
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
		).Return([]logpoller.Log{}, nil)

		usdcMessageData, err := u.GetUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, 0, txHash.String())
		assert.Errorf(t, err, fmt.Sprintf("no USDC message found prior to log index %d in tx %s", ccipLogIndex, txHash.Hex()))
		assert.Nil(t, usdcMessageData)

		lp.AssertExpectations(t)
	})
}

func TestParse(t *testing.T) {
	expectedBody, err := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f80000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc08610000000000000000")
	require.NoError(t, err)

	parsedBody, err := parseUSDCMessageSent(expectedBody)
	require.NoError(t, err)

	expectedPostParse := "0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861"

	require.Equal(t, expectedPostParse, hexutil.Encode(parsedBody))
}

func TestFilters(t *testing.T) {
	t.Run("filters of different jobs should be distinct", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		chainID := testutils.NewRandomEVMChainID()
		db := pgtest.NewSqlxDB(t)
		o := logpoller.NewORM(chainID, db, lggr)
		ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{}, 10e6)
		esc := client.NewSimulatedBackendClient(t, ec, chainID)
		lpOpts := logpoller.Opts{
			PollPeriod:               1 * time.Hour,
			FinalityDepth:            1,
			BackfillBatchSize:        1,
			RpcBatchSize:             1,
			KeepFinalizedBlocksDepth: 100,
		}
		headTracker := headtracker.NewSimulatedHeadTracker(esc, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
		if lpOpts.PollPeriod == 0 {
			lpOpts.PollPeriod = 1 * time.Hour
		}
		lp := logpoller.NewLogPoller(o, esc, lggr, headTracker, lpOpts)

		jobID1 := "job-1"
		jobID2 := "job-2"
		transmitter := utils.RandomAddress()

		f1 := logpoller.FilterName("USDC message sent", jobID1, transmitter.Hex())
		f2 := logpoller.FilterName("USDC message sent", jobID2, transmitter.Hex())

		_, err := NewUSDCReader(lggr, jobID1, transmitter, lp, true)
		assert.NoError(t, err)
		assert.True(t, lp.HasFilter(f1))

		_, err = NewUSDCReader(lggr, jobID2, transmitter, lp, true)
		assert.NoError(t, err)
		assert.True(t, lp.HasFilter(f2))

		err = CloseUSDCReader(lggr, jobID2, transmitter, lp)
		assert.NoError(t, err)
		assert.True(t, lp.HasFilter(f1))
		assert.False(t, lp.HasFilter(f2))
	})
}
