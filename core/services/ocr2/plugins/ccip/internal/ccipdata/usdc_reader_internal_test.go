package ccipdata

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogPollerClient_GetLastUSDCMessagePriorToLogIndexInTx(t *testing.T) {
	addr := utils.RandomAddress()
	txHash := common.BytesToHash(addr[:])
	ccipLogIndex := int64(100)

	expectedData := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f80000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc08610000000000000000"
	expectedPostParse := "0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861"
	lggr := logger.TestLogger(t)

	t.Run("multiple found", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)
		lp.On("IndexedLogsByTxHash",
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
			mock.Anything,
		).Return([]logpoller.Log{
			{LogIndex: ccipLogIndex - 2, Data: []byte("-2")},
			{LogIndex: ccipLogIndex - 1, Data: hexutil.MustDecode(expectedData)},
			{LogIndex: ccipLogIndex, Data: []byte("0")},
			{LogIndex: ccipLogIndex + 1, Data: []byte("1")},
		}, nil)

		usdcMessageData, err := u.GetLastUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, txHash.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedPostParse, hexutil.Encode(usdcMessageData))
		lp.AssertExpectations(t)
	})

	t.Run("none found", func(t *testing.T) {
		lp := lpmocks.NewLogPoller(t)
		u, _ := NewUSDCReader(lggr, "job_123", utils.RandomAddress(), lp, false)
		lp.On("IndexedLogsByTxHash",
			u.usdcMessageSent,
			u.transmitterAddress,
			txHash,
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		usdcMessageData, err := u.GetLastUSDCMessagePriorToLogIndexInTx(context.Background(), ccipLogIndex, txHash.String())
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
		lp := lpmocks.NewLogPoller(t)
		lggr := logger.TestLogger(t)

		jobID1 := "job-1"
		jobID2 := "job-2"
		transmitter := utils.RandomAddress()

		var registeredFilters []string   // filters registered by job-1
		var deregisteredFilters []string // should not exist in filters deregistered by job-2

		// job-1 registers filters
		lp.On("RegisterFilter", mock.Anything).Run(func(args mock.Arguments) {
			f := args.Get(0).(logpoller.Filter)
			registeredFilters = append(registeredFilters, f.Name)
		}).Return(nil)
		_, err := NewUSDCReader(lggr, jobID1, transmitter, lp, true)
		assert.NoError(t, err)

		// job-2 de-registers filters
		lp.On("UnregisterFilter", mock.Anything).Run(func(args mock.Arguments) {
			name := args.Get(0).(string)
			deregisteredFilters = append(deregisteredFilters, name)
		}).Return(nil)
		err = CloseUSDCReader(lggr, jobID2, transmitter, lp)
		assert.NoError(t, err)

		assert.NotEmpty(t, deregisteredFilters)
		assert.NotEmpty(t, registeredFilters)
		for _, f := range registeredFilters {
			assert.False(t, slices.Contains(deregisteredFilters, f))
		}
	})

}
