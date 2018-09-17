package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawalsController_CreateSuccess(t *testing.T) {
	config, _ := cltest.NewConfigWithPrivateKey()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.OracleContractAddress = &oca
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()
	hash := cltest.NewHash()
	client := app.NewHTTPClient()

	ethMock := app.MockEthClient()
	sentAt := "0x5BA0"
	nonce := "0x100"

	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", nonce)
	})

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_call", "0xDE0B6B3A7640000")
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_blockNumber", sentAt)
	})

	assert.NoError(t, app.Start())

	wr := models.WithdrawalRequest{
		Address: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		Amount:  assets.NewLink(1000000000000000000),
	}

	body, err := json.Marshal(&wr)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/withdrawals", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 200)

	assert.True(t, ethMock.AllCalled(), "Not Called")
}
