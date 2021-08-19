package web_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestETHKeysController_Index_Success(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.Dev = null.BoolFrom(true)
	cfg.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
	cfg.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
	app, cleanup := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	t.Cleanup(cleanup)

	_, addr1 := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth(), true)
	_, addr2 := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth(), false)

	ethClient.On("BalanceAt", mock.Anything, addr1, mock.Anything).Return(big.NewInt(256), nil).Twice()
	ethClient.On("BalanceAt", mock.Anything, addr2, mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(1), nil).Once()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedKeys, err := app.KeyStore.Eth().AllKeys()
	require.NoError(t, err)
	var actualBalances []webpresenters.ETHKeyResource
	err = cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 2)

	first := actualBalances[0]
	assert.Equal(t, expectedKeys[0].Address.Hex(), first.Address)
	assert.Equal(t, "0.000000000000000256", first.EthBalance.String())
	assert.Equal(t, "256", first.LinkBalance.String())

	second := actualBalances[1]
	assert.Equal(t, expectedKeys[1].Address.Hex(), second.Address)
	assert.Equal(t, "0.000000000000000001", second.EthBalance.String())
	assert.Equal(t, "1", second.LinkBalance.String())
}

func TestETHKeysController_Index_NotDev(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.Dev = null.BoolFrom(false)
	cfg.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
	cfg.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient)
	t.Cleanup(cleanup)

	cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth(), true)

	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedKeys, err := app.KeyStore.Eth().AllKeys()
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

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	err := app.Store.ORM.DB.Delete(&ethkey.Key{}, "id = ?", app.Key.ID).Error
	require.NoError(t, err)

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()

	balances := []webpresenters.ETHKeyResource{}
	err = cltest.ParseJSONAPIResponse(t, resp, &balances)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, balances, 0)
}

func TestETHKeysController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config, ethClient)
	t.Cleanup(cleanup)

	verify := cltest.MockApplicationEthCalls(t, app, ethClient)
	defer verify()

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLinkFromJuels(42)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient()

	require.NoError(t, app.Start())

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 201)

	ethClient.AssertExpectations(t)
}
