package web_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestETHKeysController_Index_Success(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	_, err := app.Store.KeyStore.NewAccount()
	require.NoError(t, err)
	require.NoError(t, app.Store.SyncDiskKeyStoreToDB())

	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Run(func(args mock.Arguments) {
		*args.Get(0).(*string) = "256"
	}).Return(nil).Once()
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Run(func(args mock.Arguments) {
		*args.Get(0).(*string) = "1"
	}).Return(nil).Once()
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(256), nil).Once()
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(1), nil).Once()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedAccounts := app.Store.KeyStore.Accounts()
	var actualBalances []presenters.ETHKey
	err = cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	assert.Len(t, actualBalances, 2)

	first := actualBalances[0]
	assert.Equal(t, expectedAccounts[0].Address.Hex(), first.Address)
	assert.Equal(t, "0.000000000000000256", first.EthBalance.String())
	assert.Equal(t, "0.000000000000000256", first.LinkBalance.String())

	second := actualBalances[1]
	assert.Equal(t, expectedAccounts[1].Address.Hex(), second.Address)
	assert.Equal(t, "0.000000000000000001", second.EthBalance.String())
	assert.Equal(t, "0.000000000000000001", second.LinkBalance.String())
}

func TestETHKeysController_Index_NoAccounts(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t, eth.NewClientWith(rpcClient, gethClient))
	defer cleanup()
	require.NoError(t, app.Start())

	err := app.Store.ORM.DB.Delete(&models.Key{}, "id = ?", app.Key.ID).Error
	require.NoError(t, err)

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()

	balances := []presenters.ETHKey{}
	err = cltest.ParseJSONAPIResponse(t, resp, &balances)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, balances, 0)
}

func TestETHKeysController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	ethClient := new(mocks.Client)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config, ethClient)
	defer cleanup()

	verify := cltest.MockApplicationEthCalls(t, app, ethClient)
	defer verify()

	ethBalanceInt := big.NewInt(100)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(ethBalanceInt, nil)
	linkBalance := assets.NewLink(42)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(linkBalance, nil)

	client := app.NewHTTPClient()

	require.NoError(t, app.StartAndConnect())

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 201)

	ethClient.AssertExpectations(t)
}
