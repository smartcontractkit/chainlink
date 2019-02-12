package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserController_UpdatePassword(t *testing.T) {
	t.Parallel()

	appWithUser, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	client := appWithUser.NewHTTPClient()

	// Invalid request
	resp, cleanup := client.Patch("/v2/user/password", bytes.NewBufferString(""))
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(resp.Body)
	require.Equal(t, 422, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)

	// Old password is wrong
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"oldPassword": "wrong password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(resp.Body)
	require.Equal(t, 409, resp.StatusCode)
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

func TestUserController_AccountBalances_NoAccounts(t *testing.T) {
	t.Parallel()

	appWithoutAccount, cleanup := cltest.NewApplication()
	defer cleanup()
	client := appWithoutAccount.NewHTTPClient()

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()

	balances := []presenters.AccountBalance{}
	err := cltest.ParseJSONAPIResponse(resp, &balances)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, balances, 0)
}

func TestUserController_AccountBalances_Success(t *testing.T) {
	t.Parallel()

	appWithAccount, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	appWithAccount.AddUnlockedKey()
	client := appWithAccount.NewHTTPClient()

	ethMock := appWithAccount.MockEthClient()
	ethMock.Context("first wallet", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getBalance", "0x0100")
		ethMock.Register("eth_call", "0x0100")
	})
	ethMock.Context("second wallet", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getBalance", "0x01")
		ethMock.Register("eth_call", "0x01")
	})

	resp, cleanup := client.Get("/v2/user/balances")
	defer cleanup()
	require.Equal(t, 200, resp.StatusCode)

	expectedAccounts := appWithAccount.Store.KeyStore.Accounts()
	actualBalances := []presenters.AccountBalance{}
	err := cltest.ParseJSONAPIResponse(resp, &actualBalances)
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
