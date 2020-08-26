package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserController_UpdatePassword(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	// Invalid request
	resp, cleanup := client.Patch("/v2/user/password", bytes.NewBufferString(""))
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)

	// Old password is wrong
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"oldPassword": "wrong password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(t, resp.Body)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)
	assert.Equal(t, "old password does not match", errors.Errors[0].Detail)

	// Success
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"newPassword": "password", "oldPassword": "password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Len(t, errors.Errors, 0)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserController_AccountBalances_NoAccounts(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	kst := new(mocks.KeyStoreInterface)
	kst.On("Accounts").Return([]accounts.Account{})
	app.Store.KeyStore = kst
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()

	balances := []presenters.AccountBalance{}
	err := cltest.ParseJSONAPIResponse(t, resp, &balances)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, balances, 0)
	kst.AssertExpectations(t)
}

func TestUserController_AccountBalances_Success(t *testing.T) {
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

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	expectedAccounts := app.Store.KeyStore.Accounts()
	actualBalances := []presenters.AccountBalance{}
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

func TestUserController_NewAPIToken(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	req, err := json.Marshal(models.ChangeAuthTokenRequest{
		Password: cltest.Password,
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token", bytes.NewBuffer(req))
	defer cleanup()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var authToken auth.Token
	err = cltest.ParseJSONAPIResponse(t, resp, &authToken)
	require.NoError(t, err)
	assert.NotEmpty(t, authToken.AccessKey)
	assert.NotEmpty(t, authToken.Secret)
}

func TestUserController_NewAPIToken_unauthorized(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	req, err := json.Marshal(models.ChangeAuthTokenRequest{
		Password: "wrong-password",
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token", bytes.NewBuffer(req))
	defer cleanup()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUserController_DeleteAPIKey(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	req, err := json.Marshal(models.ChangeAuthTokenRequest{
		Password: cltest.Password,
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token/delete", bytes.NewBuffer(req))
	defer cleanup()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserController_DeleteAPIKey_unauthorized(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	req, err := json.Marshal(models.ChangeAuthTokenRequest{
		Password: "wrong-password",
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token/delete", bytes.NewBuffer(req))
	defer cleanup()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
