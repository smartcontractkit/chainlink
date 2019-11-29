package web_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"

	"chainlink/core/assets"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()
	client := app.NewHTTPClient()

	wr := models.WithdrawalRequest{
		DestinationAddress: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		Amount:             assets.NewLink(1000000000000000000),
	}

	subscription := cltest.EmptyMockSubscription()
	txManager := new(mocks.TxManager)
	txManager.On("SubscribeToNewHeads", mock.Anything).Maybe().Return(subscription, nil)
	txManager.On("GetChainID").Maybe().Return(big.NewInt(3), nil)
	txManager.On("Register", mock.Anything).Return(big.NewInt(3), nil)
	txManager.On("ContractLINKBalance", wr).Return(*wr.Amount, nil)
	txManager.On("WithdrawLINK", wr).Return(cltest.NewHash(), nil)
	app.Store.TxManager = txManager

	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)

	assert.NoError(t, app.StartAndConnect())

	body, err := json.Marshal(&wr)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/withdrawals", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusOK)
	txManager.AssertExpectations(t)
}

func TestWithdrawalsController_BalanceTooLow(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()
	client := app.NewHTTPClient()

	contractAddress :=
		common.HexToAddress("0x3141592653589793238462643383279502884197")

	wr := models.WithdrawalRequest{
		DestinationAddress: common.HexToAddress("0xDEADEAFDEADEAFDEADEAFDEADEAFDEAD00000000"),
		ContractAddress:    contractAddress,
		Amount:             assets.NewLink(1000000000000000000),
	}

	ethMock := app.MockCallerSubscriberClient()

	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
		ethMock.Register("eth_chainId", config.ChainID())
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

	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
	assert.True(t, ethMock.AllCalled(), ethMock.Remaining())
}
