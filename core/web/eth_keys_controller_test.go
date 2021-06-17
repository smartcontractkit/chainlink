package web_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
	cfg := cltest.NewTestEVMConfig(t)
	cfg.GeneralConfig.Overrides.Dev = null.BoolFrom(true)
	cfg.Overrides.EvmNonceAutoSync = null.BoolFrom(false)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient)
	t.Cleanup(cleanup)

	cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), true)

	expectedKeys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)

	ethClient.On("BalanceAt", mock.Anything, expectedKeys[0].Address.Address(), mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("BalanceAt", mock.Anything, expectedKeys[1].Address.Address(), mock.Anything).Return(big.NewInt(1), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, expectedKeys[0].Address.Address()).Return(assets.NewLink(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, expectedKeys[1].Address.Address()).Return(assets.NewLink(1), nil).Once()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actualBalances []webpresenters.ETHKeyResource
	err = cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	require.Len(t, actualBalances, 2)

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

func TestETHKeysController_Index_NotDev(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)
	cfg := cltest.NewTestEVMConfig(t)
	cfg.GeneralConfig.Overrides.Dev = null.BoolFrom(false)
	cfg.Overrides.EvmNonceAutoSync = null.BoolFrom(false)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient)
	t.Cleanup(cleanup)

	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(256), nil).Once()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
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

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

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

	config := cltest.NewTestEVMConfig(t)
	ethClient := cltest.NewEthClientMock(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config, ethClient)
	t.Cleanup(cleanup)

	verify := cltest.MockApplicationEthCalls(t, app, ethClient)
	defer verify()

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLink(42)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient()

	require.NoError(t, app.Start())

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 201)

	ethClient.AssertExpectations(t)
}
