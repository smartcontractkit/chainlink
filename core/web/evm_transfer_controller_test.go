package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransfersController_CreateSuccess_From(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        key.Address,
		Amount:             amount,
		SkipWaitTxAttempt:  true,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	validateTxCount(t, app.GetDB(), 1)
}

func TestTransfersController_CreateSuccess_From_WEI(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	balance, err := assets.NewEthValueS("2")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount := assets.NewEthValue(1000000000000000000)

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        key.Address,
		Amount:             amount,
		SkipWaitTxAttempt:  true,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	validateTxCount(t, app.GetDB(), 1)
}

func TestTransfersController_CreateSuccess_From_BalanceMonitorDisabled(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})

	app := cltest.NewApplicationWithConfigAndKey(t, config, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        key.Address,
		Amount:             amount,
		SkipWaitTxAttempt:  true,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	validateTxCount(t, app.GetDB(), 1)
}

func TestTransfersController_TransferZeroAddressError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	client := app.NewHTTPClient(nil)
	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Amount:             amount,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestTransfersController_TransferBalanceToLowError(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(assets.NewEth(10).ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	request := models.SendEtherRequest{
		FromAddress:        key.Address,
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             amount,
		AllowHigherAmounts: false,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestTransfersController_TransferBalanceToLowError_ZeroBalance(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	balance, err := assets.NewEthValueS("0")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	request := models.SendEtherRequest{
		FromAddress:        key.Address,
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             amount,
		AllowHigherAmounts: false,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestTransfersController_JSONBindingError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBufferString(`{"address":""}`))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestTransfersController_CreateSuccess_eip1559(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient := cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)
	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil).Maybe()

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		c.EVM[0].ChainID = (*ubig.Big)(testutils.FixtureChainID)
		// NOTE: FallbackPollInterval is used in this test to quickly create TxAttempts
		// Testing triggers requires committing transactions and does not work with transactional tests
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(time.Second)
	})

	app := cltest.NewApplicationWithConfigAndKey(t, config, ethClient, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	amount, err := assets.NewEthValueS("100")
	require.NoError(t, err)

	timeout := 5 * time.Second
	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        key.Address,
		Amount:             amount,
		WaitAttemptTimeout: &timeout,
		EVMChainID:         ubig.New(evmtest.MustGetDefaultChainID(t, config.EVMConfigs())),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resource := presenters.EthTxResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	assert.NoError(t, err)

	validateTxCount(t, app.GetDB(), 1)

	// check returned data
	assert.NotEmpty(t, resource.Hash)
	assert.NotEmpty(t, resource.To)
	assert.NotEmpty(t, resource.From)
	assert.NotEmpty(t, resource.Nonce)
	assert.NotEqual(t, "unstarted", resource.State)
}

func TestTransfersController_FindTxAttempt(t *testing.T) {
	tx := txmgr.Tx{ID: 1}
	attempt := txmgr.TxAttempt{ID: 2}
	txWithAttempt := txmgr.Tx{ID: 1, TxAttempts: []txmgr.TxAttempt{attempt}}

	// happy path
	t.Run("happy_path", func(t *testing.T) {
		ctx := testutils.Context(t)
		timeout := 5 * time.Second
		var done bool
		find := func(_ context.Context, _ int64) (txmgr.Tx, error) {
			if !done {
				done = true
				return tx, nil
			}
			return txWithAttempt, nil
		}
		a, err := web.FindTxAttempt(ctx, timeout, tx, find)
		require.NoError(t, err)
		assert.Equal(t, tx.ID, a.Tx.ID)
		assert.Equal(t, attempt.ID, a.ID)
	})

	// failed to find tx
	t.Run("failed to find tx", func(t *testing.T) {
		ctx := testutils.Context(t)
		find := func(_ context.Context, _ int64) (txmgr.Tx, error) {
			return txmgr.Tx{}, fmt.Errorf("ERRORED")
		}
		_, err := web.FindTxAttempt(ctx, time.Second, tx, find)
		assert.ErrorContains(t, err, "failed to find transaction")
	})

	// timeout
	t.Run("timeout", func(t *testing.T) {
		ctx := testutils.Context(t)
		find := func(_ context.Context, _ int64) (txmgr.Tx, error) {
			return tx, nil
		}
		_, err := web.FindTxAttempt(ctx, time.Second, tx, find)
		assert.ErrorContains(t, err, "context deadline exceeded")
	})

	// context canceled
	t.Run("context canceled", func(t *testing.T) {
		ctx := testutils.Context(t)
		find := func(_ context.Context, _ int64) (txmgr.Tx, error) {
			return tx, nil
		}

		ctx, cancel := context.WithCancel(ctx)
		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		_, err := web.FindTxAttempt(ctx, 5*time.Second, tx, find)
		assert.ErrorContains(t, err, "context canceled")
	})
}

func validateTxCount(t *testing.T, ds sqlutil.DataSource, count int) {
	txStore := txmgr.NewTxStore(ds, logger.TestLogger(t))

	txes, err := txStore.GetAllTxes(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, txes, count)
}
