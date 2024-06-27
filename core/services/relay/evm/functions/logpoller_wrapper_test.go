package functions_test

import (
	"encoding/hex"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type subscriber struct {
	updates       sync.WaitGroup
	expectedCalls int
}

func (s *subscriber) UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error {
	if s.expectedCalls == 0 {
		panic("unexpected call to UpdateRoutes")
	}
	if activeCoordinator == (common.Address{}) {
		panic("activeCoordinator should not be zero")
	}
	s.expectedCalls--
	s.updates.Done()
	return nil
}

func newSubscriber(expectedCalls int) *subscriber {
	sub := &subscriber{expectedCalls: expectedCalls}
	sub.updates.Add(expectedCalls)
	return sub
}

func addr(t *testing.T, lastByte string) []byte {
	contractAddr, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000" + lastByte)
	require.NoError(t, err)
	return contractAddr
}

func setUp(t *testing.T, updateFrequencySec uint32) (*lpmocks.LogPoller, types.LogPollerWrapper, *evmclimocks.Client) {
	lggr := logger.TestLogger(t)
	client := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	config := config.PluginConfig{
		ContractUpdateCheckFrequencySec: updateFrequencySec,
		ContractVersion:                 1,
	}
	lpWrapper, err := functions.NewLogPollerWrapper(gethcommon.Address{}, config, client, lp, lggr)
	require.NoError(t, err)

	lp.On("LatestBlock").Return(int64(100), nil)

	return lp, lpWrapper, client
}

func TestLogPollerWrapper_SingleSubscriberEmptyEvents(t *testing.T) {
	t.Parallel()
	lp, lpWrapper, client := setUp(t, 100_000) // check only once

	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{}, nil)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(addr(t, "01"), nil)
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	subscriber := newSubscriber(1)
	lpWrapper.SubscribeToUpdates("mock_subscriber", subscriber)

	require.NoError(t, lpWrapper.Start(testutils.Context(t)))
	subscriber.updates.Wait()
	reqs, resps, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	require.Equal(t, 0, len(reqs))
	require.Equal(t, 0, len(resps))
	lpWrapper.Close()
}

func TestLogPollerWrapper_ErrorOnZeroAddresses(t *testing.T) {
	t.Parallel()
	_, lpWrapper, client := setUp(t, 100_000) // check only once

	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(addr(t, "00"), nil)

	require.NoError(t, lpWrapper.Start(testutils.Context(t)))
	_, _, err := lpWrapper.LatestEvents()
	require.Error(t, err)
	lpWrapper.Close()
}
