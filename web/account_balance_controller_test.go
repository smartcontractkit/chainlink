package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
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
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(resp), &ab)
	assert.NoError(t, err)

	assert.Equal(t, account.Address.Hex(), ab.Address)
	assert.Equal(t, "0.000000000000000256", ab.EthBalance.String())
	assert.Equal(t, "0.000000000000000256", ab.LinkBalance.String())
}

func TestAccountBalanceController_Withdraw(t *testing.T) {
	config, _ := cltest.NewConfigWithPrivateKey()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.OracleContractAddress = &oca
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()
	client := app.NewHTTPClient()

	ethMock := app.MockEthClient()
	sentAt := "0x5BA0"
	nonce := "0x100"

	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", nonce)
	})

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_call", "0xDE0B6B3A7640000")
		ethMock.Register("eth_blockNumber", sentAt)
	})

	assert.NoError(t, app.Start())

	wr := models.WithdrawalRequest{
		Address: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		Amount:  assets.NewLink(1000000000000000000),
	}

	body, err := json.Marshal(&wr)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/withdraw", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 200)

	assert.True(t, ethMock.AllCalled(), "Not Called")
}
