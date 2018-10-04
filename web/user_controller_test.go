package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
)

func TestUserController_UpdatePassword(t *testing.T) {
	appWithUser, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := appWithUser.NewHTTPClient()

	// Invalid request
	resp, cleanup := client.Patch("/v2/user/password", bytes.NewBufferString(""))
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 422, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)

	// Old password is wrong
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"oldPassword": "wrong password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 409, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)
	assert.Equal(t, "Old password does not match", errors.Errors[0].Detail)

	// Success
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"newPassword": "password", "oldPassword": "password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUserController_AccountBalances_Error(t *testing.T) {
	t.Parallel()

	appWithoutAccount, cleanup := cltest.NewApplication()
	defer cleanup()
	client := appWithoutAccount.NewHTTPClient()

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "No Ethereum Accounts configured", errors.Errors[0].Detail)
}

func TestUserController_AccountBalances_Success(t *testing.T) {
	t.Parallel()

	appWithAccount, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := appWithAccount.NewHTTPClient()

	ethMock := appWithAccount.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100")
	ethMock.Register("eth_call", "0x0100")

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode)

	account, err := appWithAccount.Store.KeyStore.GetAccount()
	assert.NoError(t, err)

	ab := presenters.AccountBalance{}
	err = cltest.ParseJSONAPIResponse(resp, &ab)
	assert.NoError(t, err)

	assert.Equal(t, account.Address.Hex(), ab.Address)
	assert.Equal(t, "0.000000000000000256", ab.EthBalance.String())
	assert.Equal(t, "0.000000000000000256", ab.LinkBalance.String())
}
