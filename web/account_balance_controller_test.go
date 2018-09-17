package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
)

func TestAccountBalanceController_IndexError(t *testing.T) {
	t.Parallel()

	appWithoutAccount, cleanup := cltest.NewApplication()
	defer cleanup()
	client := appWithoutAccount.NewHTTPClient()

	resp, cleanup := client.Get("/v2/account_balance")
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "No Ethereum Accounts configured", errors.Errors[0].Detail)
}

func TestAccountBalanceController_Index(t *testing.T) {
	t.Parallel()

	appWithAccount, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := appWithAccount.NewHTTPClient()

	ethMock := appWithAccount.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100")
	ethMock.Register("eth_call", "0x0100")

	resp, cleanup := client.Get("/v2/account_balance")
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
