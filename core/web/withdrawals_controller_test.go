package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
)

// verifyLinkBalanceCheck(t) is used to check that the address checked in a
// mocked call to the balance method of the LINK contract is correct
func verifyLinkBalanceCheck(address common.Address, t *testing.T) func(interface{}, ...interface{}) error {
	return func(_ interface{}, arg ...interface{}) error {
		balanceAddress :=
			cltest.ExtractTargetAddressFromERC20EthEthCallMock(t, arg)
		assert.Equal(t, balanceAddress, address)
		return nil
	}
}

func TestWithdrawalsController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(config)
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
		ethMock.Register(
			"eth_call", "0xDE0B6B3A7640000",
			verifyLinkBalanceCheck(oca, t))
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_blockNumber", sentAt)
	})

	assert.NoError(t, app.StartAndConnect())

	wr := models.WithdrawalRequest{
		DestinationAddress: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		Amount:             assets.NewLink(1000000000000000000),
	}

	body, err := json.Marshal(&wr)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/withdrawals", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 200)

	assert.True(t, ethMock.AllCalled(), "Not Called")
}

func TestWithdrawalsController_BalanceTooLow(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(config)
	defer cleanup()
	client := app.NewHTTPClient()

	contractAddress :=
		common.HexToAddress("0x3141592653589793238462643383279502884197")

	wr := models.WithdrawalRequest{
		DestinationAddress: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		ContractAddress:    contractAddress,
		Amount:             assets.NewLink(1000000000000000000),
	}

	ethMock := app.MockEthClient()

	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})
	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_call", "0x0",
			verifyLinkBalanceCheck(contractAddress, t))
	})

	assert.NoError(t, app.StartAndConnect())

	body, err := json.Marshal(&wr)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/withdrawals", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 400)
	assert.True(t, ethMock.AllCalled(), ethMock.Remaining())
}
