package functions_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

const (
	addr1 = "9ed925d8206a4f88a2f643b28b3035b315753cd6"
	addr2 = "ea6721ac65bced841b8ec3fc5fedea6141a0ade4"
	addr3 = "84689acc87ff22841b8ec378300da5e141a99911"
)

func sampleEncodedAllowlist(t *testing.T) []byte {
	abiEncodedAddresses :=
		"0000000000000000000000000000000000000000000000000000000000000020" +
			"0000000000000000000000000000000000000000000000000000000000000002" +
			"000000000000000000000000" + addr1 +
			"000000000000000000000000" + addr2
	rawData, err := hex.DecodeString(abiEncodedAddresses)
	require.NoError(t, err)
	return rawData
}

func TestAllowlist_UpdateAndCheck(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(sampleEncodedAllowlist(t), nil)
	allowlist, err := functions.NewOnchainAllowlist(client, common.Address{}, 1, logger.TestLogger(t))
	require.NoError(t, err)

	require.NoError(t, allowlist.UpdateFromContract(context.Background()))
	require.False(t, allowlist.Allow(common.Address{}))
	require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
	require.True(t, allowlist.Allow(common.HexToAddress(addr2)))
	require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
}

func TestAllowlist_UpdatePeriodically(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	client := mocks.NewClient(t)
	client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cancel()
	}).Return(sampleEncodedAllowlist(t), nil)
	allowlist, err := functions.NewOnchainAllowlist(client, common.Address{}, 1, logger.TestLogger(t))
	require.NoError(t, err)

	allowlist.UpdatePeriodically(ctx, time.Millisecond*10, time.Second*1)
	require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
	require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
}
