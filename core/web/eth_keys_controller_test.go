package web_test

import (
	"math/big"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	webpresenters "github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestETHKeysController_Index_Success(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	k0, addr0 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	// disabled keys
	k1, addr1 := cltest.RandomKey{Disabled: true}.MustInsert(t, app.KeyStore.Eth())
	k2, addr2 := cltest.RandomKey{Disabled: true}.MustInsert(t, app.KeyStore.Eth())
	expectedKeys := []ethkey.KeyV2{k0, k1, k2}

	ethClient.On("BalanceAt", mock.Anything, addr0, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr1, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr2, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr0, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr1, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr2, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/keys/evm")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actualBalances []webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 3)

	for _, balance := range actualBalances {
		if balance.Address == expectedKeys[0].Address.Hex() {
			assert.Equal(t, "0.000000000000000256", balance.EthBalance.String())
			assert.Equal(t, "256", balance.LinkBalance.String())
		} else {
			assert.Equal(t, "0.000000000000000001", balance.EthBalance.String())
			assert.Equal(t, "1", balance.LinkBalance.String())
		}
	}
}

func TestETHKeysController_Index_Errors(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(nil, errors.New("fake error")).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(nil, errors.New("fake error")).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actualBalances []webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 1)

	balance := actualBalances[0]
	assert.Equal(t, addr.String(), balance.Address)
	assert.Nil(t, balance.EthBalance)
	assert.Nil(t, balance.LinkBalance)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457584007913129639935", balance.MaxGasPriceWei.String())
}

func TestETHKeysController_Index_Disabled(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(false)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actualBalances []webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 1)

	balance := actualBalances[0]
	assert.Equal(t, addr.String(), balance.Address)
	assert.Nil(t, balance.EthBalance)
	assert.Nil(t, balance.LinkBalance)
	assert.Nil(t, balance.MaxGasPriceWei)
}

func TestETHKeysController_Index_NotDev(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})

	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()

	app := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedKeys, err := app.KeyStore.Eth().GetAll(testutils.Context(t))
	require.NoError(t, err)
	var actualBalances []webpresenters.ETHKeyResource
	err = cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 1)

	only := actualBalances[0]
	assert.Equal(t, expectedKeys[0].Address.Hex(), only.Address)
	assert.Equal(t, "0.000000000000000256", only.EthBalance.String())
	assert.Equal(t, "256", only.LinkBalance.String())
}

func TestETHKeysController_Index_NoAccounts(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()

	balances := []webpresenters.ETHKeyResource{}
	err := cltest.ParseJSONAPIResponse(t, resp, &balances)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, balances, 0)
}

func TestETHKeysController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	app := cltest.NewApplicationWithConfigAndKey(t, config, ethClient)

	sub := commonmocks.NewSubscription(t)
	cltest.MockApplicationEthCalls(t, app, ethClient, sub)

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLinkFromJuels(42)
	ethClient.On("LINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient(nil)

	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	chainURL := url.URL{Path: "/v2/keys/evm"}
	query := chainURL.Query()
	query.Set("evmChainID", cltest.FixtureChainID.String())
	chainURL.RawQuery = query.Encode()

	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var balance webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &balance)
	assert.NoError(t, err)

	assert.Equal(t, ethBalanceInt, balance.EthBalance.ToInt())
	assert.Equal(t, linkBalance, balance.LinkBalance)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457584007913129639935", balance.MaxGasPriceWei.String())
}

func TestETHKeysController_ChainSuccess_UpdateNonce(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
	assert.NoError(t, err)

	assert.Equal(t, cltest.FormatWithPrefixedChainID(cltest.FixtureChainID.String(), key.Address.String()), updatedKey.ID)
	assert.Equal(t, key.Address.String(), updatedKey.Address)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, false, updatedKey.Disabled)
}

func TestETHKeysController_ChainSuccess_Disable(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	enabled := "false"
	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("enabled", enabled)

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
	assert.NoError(t, err)

	assert.Equal(t, cltest.FormatWithPrefixedChainID(updatedKey.EVMChainID.String(), key.Address.String()), updatedKey.ID)
	assert.Equal(t, key.Address.String(), updatedKey.Address)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, true, updatedKey.Disabled)
}

func TestETHKeysController_ChainSuccess_Enable(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// disabled key
	key, addr := cltest.RandomKey{Disabled: true}.MustInsert(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	enabled := "true"
	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("enabled", enabled)

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
	assert.NoError(t, err)

	assert.Equal(t, cltest.FormatWithPrefixedChainID(cltest.FixtureChainID.String(), key.Address.String()), updatedKey.ID)
	assert.Equal(t, key.Address.String(), updatedKey.Address)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, false, updatedKey.Disabled)
}

func TestETHKeysController_ChainSuccess_ResetWithAbandon(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	chain := app.GetRelayers().LegacyEVMChains().Slice()[0]
	subject := uuid.New()
	strategy := commontxmmocks.NewTxStrategy(t)
	strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
	strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("*txmgr.evmTxStore")).Return(nil, nil)
	_, err := chain.TxManager().CreateTransaction(testutils.Context(t), txmgr.TxRequest{
		FromAddress:    addr,
		ToAddress:      testutils.NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		FeeLimit:       uint64(1000),
		Meta:           nil,
		Strategy:       strategy,
	})
	assert.NoError(t, err)

	txStore := txmgr.NewTxStore(app.GetDB(), logger.TestLogger(t))

	txes, err := txStore.FindTxesByFromAddressAndState(testutils.Context(t), addr, "fatal_error")
	require.NoError(t, err)
	require.Len(t, txes, 0)

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("abandon", "true")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedKey webpresenters.ETHKeyResource
	err = cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
	assert.NoError(t, err)

	assert.Equal(t, cltest.FormatWithPrefixedChainID(cltest.FixtureChainID.String(), key.Address.String()), updatedKey.ID)
	assert.Equal(t, key.Address.String(), updatedKey.Address)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, false, updatedKey.Disabled)

	txes, err = txStore.FindTxesByFromAddressAndState(testutils.Context(t), addr, "fatal_error")
	require.NoError(t, err)
	require.Len(t, txes, 1)

	tx := txes[0]
	assert.Equal(t, "abandoned", tx.Error.String)
}

func TestETHKeysController_ChainFailure_InvalidAbandon(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("abandon", "invalid")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_AbandonUnstarted(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	t.Run("invalid abandonStarted: returns StatusBadRequest", func(t *testing.T) {
		client := app.NewHTTPClient(nil)
		chainURL := url.URL{Path: "/v2/keys/evm/chain"}
		query := chainURL.Query()

		query.Set("address", addr.Hex())
		query.Set("evmChainID", cltest.FixtureChainID.String())
		query.Set("abandon", "true")
		query.Set("abandonUnstarted", "invalid")
		query.Set("subject", "")

		chainURL.RawQuery = query.Encode()
		resp, cleanup := client.Post(chainURL.String(), nil)
		defer cleanup()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("valid abandonStarted with invalid subject: returns StatusBadRequest", func(t *testing.T) {
		client := app.NewHTTPClient(nil)
		chainURL := url.URL{Path: "/v2/keys/evm/chain"}
		query := chainURL.Query()

		query.Set("address", addr.Hex())
		query.Set("evmChainID", cltest.FixtureChainID.String())
		query.Set("abandon", "true")
		query.Set("abandonUnstarted", "true")
		query.Set("subject", "invalid")

		chainURL.RawQuery = query.Encode()
		resp, cleanup := client.Post(chainURL.String(), nil)
		defer cleanup()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestETHKeysController_ChainSuccess_AbandonUnstarted(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	txStore := txmgr.NewTxStore(app.GetDB(), logger.TestLogger(t))

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Times(3)
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Times(3)

	require.NoError(t, app.Start(ctx))

	chain := app.GetRelayers().LegacyEVMChains().Slice()[0]

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()
	query.Set("address", addr.Hex())
	query.Set("evmChainID", chain.ID().String())
	query.Set("abandon", "false")
	query.Set("abandonUnstarted", "true")

	t.Run("abandonUnstarted with empty subject provided and no txes: handled gracefully", func(t *testing.T) {
		// Makes the request
		query.Set("subject", "")
		chainURL.RawQuery = query.Encode()
		resp, cleanup := client.Post(chainURL.String(), nil)
		defer cleanup()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedKey webpresenters.ETHKeyResource
		err := cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
		assert.NoError(t, err)

		assert.Equal(t, cltest.FormatWithPrefixedChainID(chain.ID().String(), key.Address.String()), updatedKey.ID)
		assert.Equal(t, key.Address.String(), updatedKey.Address)
		assert.Equal(t, chain.ID().String(), updatedKey.EVMChainID.String())
		assert.Equal(t, false, updatedKey.Disabled)

		txesCount, err := txStore.CountTransactionsByState(testutils.Context(t), "unstarted", chain.ID())
		require.NoError(t, err)
		require.Equal(t, txesCount, uint32(0))
	})

	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	txReq := txmgr.TxRequest{
		FromAddress:    addr,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		FeeLimit:       gasLimit,
		Meta:           nil,
		Strategy:       nil,
		SignalCallback: true,
	}

	strategyWithoutSubject := txmgrcommon.NewDropOldestStrategy(uuid.UUID{}, uint32(5))
	randomSubject := uuid.New()
	strategyWithSubject := txmgrcommon.NewDropOldestStrategy(randomSubject, uint32(5))

	t.Run("abandonUnstarted with empty subject provided and txes with diff subjects: deletes all", func(t *testing.T) {
		txReq.Strategy = strategyWithoutSubject
		// Create and assert 2 unstarted transactions without subject
		_, err := txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		assertUnstartedTxesCountPerSubject(t, txStore, 2, uuid.UUID{})

		// Create and assert 2 unstarted transactions with subject
		txReq.Strategy = strategyWithSubject
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		assertUnstartedTxesCountPerSubject(t, txStore, 2, randomSubject)

		// Makes the request
		query.Set("subject", "")

		chainURL.RawQuery = query.Encode()
		resp, cleanup := client.Post(chainURL.String(), nil)
		defer cleanup()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedKey webpresenters.ETHKeyResource
		err = cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
		assert.NoError(t, err)

		assert.Equal(t, cltest.FormatWithPrefixedChainID(chain.ID().String(), key.Address.String()), updatedKey.ID)
		assert.Equal(t, key.Address.String(), updatedKey.Address)
		assert.Equal(t, chain.ID().String(), updatedKey.EVMChainID.String())
		assert.Equal(t, false, updatedKey.Disabled)

		// Asserts that all have been deleted.
		txesCount, err := txStore.CountTransactionsByState(testutils.Context(t), "unstarted", chain.ID())
		require.NoError(t, err)
		require.Equal(t, txesCount, uint32(0))
	})

	t.Run("abandonUnstarted with subject provided and txes with diff subjects: deletes only those subject txes", func(t *testing.T) {
		txReq.Strategy = strategyWithoutSubject
		// Create and assert 2 unstarted transactions without subject
		_, err := txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		assertUnstartedTxesCountPerSubject(t, txStore, 2, uuid.UUID{})

		// Create and assert 2 unstarted transactions with subject
		txReq.Strategy = strategyWithSubject
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		_, err = txStore.CreateTransaction(tests.Context(t), txReq, chain.ID())
		assert.NoError(t, err)
		assertUnstartedTxesCountPerSubject(t, txStore, 2, randomSubject)

		// Makes the request
		query.Set("subject", randomSubject.String())

		chainURL.RawQuery = query.Encode()
		resp, cleanup := client.Post(chainURL.String(), nil)
		defer cleanup()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedKey webpresenters.ETHKeyResource
		err = cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
		assert.NoError(t, err)

		assert.Equal(t, cltest.FormatWithPrefixedChainID(chain.ID().String(), key.Address.String()), updatedKey.ID)
		assert.Equal(t, key.Address.String(), updatedKey.Address)
		assert.Equal(t, chain.ID().String(), updatedKey.EVMChainID.String())
		assert.Equal(t, false, updatedKey.Disabled)

		// Asserts that only the txes of the provided subject have been deleted.
		assertUnstartedTxesCountPerSubject(t, txStore, 2, uuid.UUID{})
		assertUnstartedTxesCountPerSubject(t, txStore, 0, randomSubject)
	})
}

func TestETHKeysController_ChainFailure_InvalidEnabled(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("enabled", "invalid")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_InvalidAddress(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", "invalid address")
	query.Set("evmChainID", cltest.FixtureChainID.String())

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_MissingAddress(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", testutils.NewAddress().Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_InvalidChainID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", testutils.NewAddress().Hex())
	query.Set("evmChainID", "bad chain ID")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_MissingChainID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", "123456789")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestETHKeysController_DeleteSuccess(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	// enabled keys
	key0, addr0 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	_, addr1 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	ethClient.On("BalanceAt", mock.Anything, addr0, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr1, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr0, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr1, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/" + addr0.Hex()}
	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()
	t.Log(resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var deletedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &deletedKey)
	assert.NoError(t, err)

	assert.Equal(t, cltest.FormatWithPrefixedChainID(cltest.FixtureChainID.String(), key0.Address.String()), deletedKey.ID)
	assert.Equal(t, key0.Address.String(), deletedKey.Address)
	assert.Equal(t, cltest.FixtureChainID.String(), deletedKey.EVMChainID.String())
	assert.Equal(t, false, deletedKey.Disabled)

	resp, cleanup2 := client.Get("/v2/keys/evm")
	defer cleanup2()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actualBalances []webpresenters.ETHKeyResource
	err = cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 1)

	balance := actualBalances[0]
	assert.Equal(t, addr1.String(), balance.Address)
}

func TestETHKeysController_DeleteFailure_InvalidAddress(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm" + "/bad_address"}

	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestETHKeysController_DeleteFailure_KeyMissing(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(ctx, cltest.Password))

	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	chainURL := url.URL{Path: "/v2/keys/evm/" + testutils.NewAddress().Hex()}

	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()
	t.Log(resp)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func assertUnstartedTxesCountPerSubject(t *testing.T, txStore txmgr.TestEvmTxStore, expected int64, subject uuid.UUID) {
	t.Helper()
	count, err := txStore.CountTxesByStateAndSubject(tests.Context(t), "unstarted", subject)
	require.NoError(t, err)
	require.Equal(t, int(expected), count)
}
