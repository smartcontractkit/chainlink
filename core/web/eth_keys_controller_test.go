package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestETHKeysController_Index_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	app.AddUnlockedKey()
	client := app.NewHTTPClient()

	ethMock := app.EthMock
	ethMock.Context("first wallet", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getBalance", "0x100")
		ethMock.Register("eth_call", "0x100")
	})
	ethMock.Context("second wallet", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getBalance", "0x1")
		ethMock.Register("eth_call", "0x1")
	})

	app.Store.SyncDiskKeyStoreToDB()

	resp, cleanup := client.Get("/v2/keys/eth")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedAccounts := app.Store.KeyStore.Accounts()
	actualBalances := []presenters.ETHKey{}
	err := cltest.ParseJSONAPIResponse(t, resp, &actualBalances)
	assert.NoError(t, err)

	first := actualBalances[0]
	assert.Equal(t, expectedAccounts[0].Address.Hex(), first.Address)
	assert.Equal(t, "0.000000000000000256", first.EthBalance.String())
	assert.Equal(t, "0.000000000000000256", first.LinkBalance.String())

	second := actualBalances[1]
	assert.Equal(t, expectedAccounts[1].Address.Hex(), second.Address)
	assert.Equal(t, "0.000000000000000001", second.EthBalance.String())
	assert.Equal(t, "0.000000000000000001", second.LinkBalance.String())
}

// func TestETHKeysController_CreateSuccess(t *testing.T) {
// 	t.Parallel()

// 	config, _ := cltest.NewConfig(t)
// 	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
// 		cltest.EthMockRegisterChainID,
// 		cltest.EthMockRegisterGetBalance,
// 	)
// 	defer cleanup()

// 	ethMock := app.EthMock
// 	client := app.NewHTTPClient()

// 	assert.NoError(t, app.StartAndConnect())

// 	body, err := json.Marshal(&request)
// 	assert.NoError(t, err)

// 	resp, cleanup := client.Post("/v2/keys/eth", nil)
// 	defer cleanup()

// 	cltest.AssertServerResponse(t, resp, 201)

// 	ethMock.AllCalled()
// }

func TestETHKeysController_Create_JSONBindingError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	resp, cleanup := client.Post("/v2/keys/eth", nil)
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 422)
}
