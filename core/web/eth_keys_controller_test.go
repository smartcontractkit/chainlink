package web_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestETHKeysController_Index_Success(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.DevMode = false
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
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, addr0).Return(assets.NewLinkFromJuels(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, addr1).Return(assets.NewLinkFromJuels(1), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, addr2).Return(assets.NewLinkFromJuels(1), nil).Once()

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

func TestETHKeysController_Index_NotDev(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})

	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(256), nil).Once()

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

	sub := evmMocks.NewSubscription(t)
	cltest.MockApplicationEthCalls(t, app, ethClient, sub)

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLinkFromJuels(42)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	require.NoError(t, app.Start(testutils.Context(t)))

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusCreated)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestETHKeysController_UpdateSuccess(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	app := cltest.NewApplicationWithConfigAndKey(t, config, ethClient)

	sub := evmMocks.NewSubscription(t)
	cltest.MockApplicationEthCalls(t, app, ethClient, sub)

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLinkFromJuels(42)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	require.NoError(t, app.Start(testutils.Context(t)))

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusCreated)

	keys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, keys)

	key := keys[0]
	resp, cleanup = client.Put("/v2/keys/eth/"+key.Address.Hex()+"?maxGasPriceGWei=777", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusOK)

	chain, err := app.Chains.EVM.Get(&cltest.FixtureChainID)
	require.NoError(t, err)

	require.Equal(t, assets.GWei(777), chain.Config().KeySpecificMaxGasPriceWei(key.Address))
}
