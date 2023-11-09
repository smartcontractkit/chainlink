package functions_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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
	config := functions.OnchainAllowlistConfig{
		ContractVersion:    1,
		ContractAddress:    common.Address{},
		BlockConfirmations: 1,
	}
	allowlist, err := functions.NewOnchainAllowlist(client, config, logger.TestLogger(t))
	require.NoError(t, err)

	err = allowlist.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, allowlist.Close())
	})

	require.NoError(t, allowlist.UpdateFromContract(testutils.Context(t)))
	require.False(t, allowlist.Allow(common.Address{}))
	require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
	require.True(t, allowlist.Allow(common.HexToAddress(addr2)))
	require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
}

func TestAllowlist_UnsupportedVersion(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	config := functions.OnchainAllowlistConfig{
		ContractVersion:    0,
		ContractAddress:    common.Address{},
		BlockConfirmations: 1,
	}
	_, err := functions.NewOnchainAllowlist(client, config, logger.TestLogger(t))
	require.Error(t, err)
}

func TestAllowlist_UpdatePeriodically(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(testutils.Context(t))
	client := mocks.NewClient(t)
	client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)
	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cancel()
	}).Return(sampleEncodedAllowlist(t), nil)
	config := functions.OnchainAllowlistConfig{
		ContractAddress:    common.Address{},
		ContractVersion:    1,
		BlockConfirmations: 1,
		UpdateFrequencySec: 1,
		UpdateTimeoutSec:   1,
	}
	allowlist, err := functions.NewOnchainAllowlist(client, config, logger.TestLogger(t))
	require.NoError(t, err)

	err = allowlist.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, allowlist.Close())
	})

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return allowlist.Allow(common.HexToAddress(addr1)) && !allowlist.Allow(common.HexToAddress(addr3))
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}
