package web_test

import (
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/pkg/errors"

	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
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

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// disabled key
	k0, addr0 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)
	// enabled keys
	k1, addr1 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), false)
	k2, addr2 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), false)
	expectedKeys := []ethkey.KeyV2{k0, k1, k2}

	ethClient.On("BalanceAt", mock.Anything, addr0, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr1, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr2, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr0, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr1, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr2, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	resp, cleanup := client.Get("/v2/keys/eth")
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

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(nil, errors.New("fake error")).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(nil, errors.New("fake error")).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(false)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})

	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()

	app := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedKeys, err := app.KeyStore.Eth().GetAll()
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
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

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

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	require.NoError(t, app.Start(testutils.Context(t)))

	resp, cleanup := client.Post("/v2/keys/evm", nil)
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

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	nextNonce := 52
	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("nextNonce", strconv.Itoa(nextNonce))

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &updatedKey)
	assert.NoError(t, err)

	assert.Equal(t, key.ID(), updatedKey.ID)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, false, updatedKey.Disabled)
	assert.Equal(t, int64(nextNonce), updatedKey.NextNonce)
}

func TestETHKeysController_ChainSuccess_Disable(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	assert.Equal(t, key.ID(), updatedKey.ID)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, int64(0), updatedKey.NextNonce)
	assert.Equal(t, true, updatedKey.Disabled)
}

func TestETHKeysController_ChainSuccess_Enable(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// disabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), false)

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	assert.Equal(t, key.ID(), updatedKey.ID)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, int64(0), updatedKey.NextNonce)
	assert.Equal(t, false, updatedKey.Disabled)
}

func TestETHKeysController_ChainSuccess_ResetWithAbandon(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled key
	key, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, addr, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	chain := app.GetChains().EVM.Chains()[0]
	subject := uuid.New()
	strategy := commontxmmocks.NewTxStrategy(t)
	strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
	strategy.On("PruneQueue", mock.AnythingOfType("*txmgr.evmTxStore"), mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)
	_, err := chain.TxManager().CreateTransaction(txmgr.TxRequest{
		FromAddress:    addr,
		ToAddress:      testutils.NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		FeeLimit:       uint32(1000),
		Meta:           nil,
		Strategy:       strategy,
	})
	assert.NoError(t, err)

	var count int
	err = app.GetSqlxDB().Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = 'fatal_error'`, addr)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	assert.Equal(t, key.ID(), updatedKey.ID)
	assert.Equal(t, cltest.FixtureChainID.String(), updatedKey.EVMChainID.String())
	assert.Equal(t, int64(0), updatedKey.NextNonce)
	assert.Equal(t, false, updatedKey.Disabled)

	var s string
	err = app.GetSqlxDB().Get(&s, `SELECT error FROM eth_txes WHERE from_address = $1 AND state = 'fatal_error'`, addr)
	require.NoError(t, err)
	assert.Equal(t, "abandoned", s)
}

func TestETHKeysController_ChainFailure_InvalidAddress(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	nextNonce := 52
	query.Set("address", "invalid address")
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("nextNonce", strconv.Itoa(nextNonce))

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	// TODO once cleared, update to http.StatusBadRequest https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_MissingAddress(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	nextNonce := 52
	query.Set("address", testutils.NewAddress().Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("nextNonce", strconv.Itoa(nextNonce))

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	// TODO once cleared, update to http.StatusNotFound https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_InvalidChainID(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	nextNonce := 52
	query.Set("address", testutils.NewAddress().Hex())
	query.Set("evmChainID", "bad chain ID")
	query.Set("nextNonce", strconv.Itoa(nextNonce))

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	// TODO once cleared, update to http.StatusBadRequest https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_MissingChainID(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	nextNonce := 52
	query.Set("address", addr.Hex())
	query.Set("evmChainID", "123456789")
	query.Set("nextNonce", strconv.Itoa(nextNonce))

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	// TODO once cleared, update to http.StatusNotFound https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_ChainFailure_InvalidNonce(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)

	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled key
	_, addr := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	query.Set("address", addr.Hex())
	query.Set("evmChainID", cltest.FixtureChainID.String())
	query.Set("nextNonce", "bad nonce")

	chainURL.RawQuery = query.Encode()
	resp, cleanup := client.Post(chainURL.String(), nil)
	defer cleanup()
	// TODO once cleared, update to http.StatusBadRequest https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_DeleteSuccess(t *testing.T) {
	t.Parallel()
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	// enabled keys
	key0, addr0 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)
	_, addr1 := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, addr0, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, addr1, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr0, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()
	ethClient.On("LINKBalance", mock.Anything, addr1, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/" + addr0.Hex()}
	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()
	t.Log(resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var deletedKey webpresenters.ETHKeyResource
	err := cltest.ParseJSONAPIResponse(t, resp, &deletedKey)
	assert.NoError(t, err)

	assert.Equal(t, key0.ID(), deletedKey.ID)
	assert.Equal(t, cltest.FixtureChainID.String(), deletedKey.EVMChainID.String())
	assert.Equal(t, false, deletedKey.Disabled)
	assert.Equal(t, int64(0), deletedKey.NextNonce)

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
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm" + "/bad_address"}

	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()
	// TODO once cleared, update to http.StatusBadRequest https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestETHKeysController_DeleteFailure_KeyMissing(t *testing.T) {
	t.Parallel()
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.KeyStore.Unlock(cltest.Password))

	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	chainURL := url.URL{Path: "/v2/keys/evm/" + testutils.NewAddress().Hex()}

	resp, cleanup := client.Delete(chainURL.String())
	defer cleanup()
	t.Log(resp)
	// TODO once cleared, update to http.StatusNotFound https://smartcontract-it.atlassian.net/browse/BCF-2346
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
