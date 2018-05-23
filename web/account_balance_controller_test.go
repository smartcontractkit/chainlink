package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func TestAccountBalanceController_Index(t *testing.T) {
	t.Parallel()

	appWithoutAccount, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(appWithoutAccount.Server.URL + "/v2/account_balance")
	body := cltest.ParseErrorsJSON(resp.Body)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "No Ethereum Accounts configured", body.Errors[0])

	appWithAccount, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := appWithAccount.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100")
	ethMock.Register("eth_call", "0x0100")

	resp = cltest.BasicAuthGet(appWithAccount.Server.URL + "/v2/account_balance")
	assert.Equal(t, 200, resp.StatusCode)

	account, err := appWithAccount.Store.KeyStore.GetAccount()
	assert.NoError(t, err)

	ab := presenters.AccountBalance{}
	err = web.ParseResponse(cltest.ParseResponseBody(resp), &ab)
	assert.NoError(t, err)

	assert.Equal(t, account.Address.Hex(), ab.Address)
	assert.Equal(t, "0.000000000000000256", ab.EthBalance.String())
	assert.Equal(t, "0.000000000000000256", ab.LinkBalance.String())
}
