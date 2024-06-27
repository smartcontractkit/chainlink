package functions_test

import (
	"encoding/hex"
	"testing"

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
)

func TestLogPollerWrapper_BasicEmptyEvents(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	client := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	config := config.PluginConfig{
		ContractUpdateCheckFrequencySec: 1000000, // only once
		ContractVersion:                 1,
	}
	lpWrapper, err := functions.NewLogPollerWrapper(gethcommon.Address{}, config, client, lp, lggr)
	require.NoError(t, err)

	lp.On("LatestBlock").Return(int64(100), nil)
	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{}, nil)

	contractAddr, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(contractAddr, nil)
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	require.NoError(t, lpWrapper.Start(testutils.Context(t)))
	reqs, resps, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	require.Equal(t, 0, len(reqs))
	require.Equal(t, 0, len(resps))
	lpWrapper.Close()
}
