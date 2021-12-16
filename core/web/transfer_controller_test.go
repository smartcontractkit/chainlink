package web_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransfersController_CreateSuccess_From(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethClient.On("PendingNonceAt", mock.Anything, key.Address.Address()).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address.Address(), (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        key.Address.Address(),
		Amount:             *assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	cltest.AssertCount(t, app.GetSqlxDB(), "eth_txes", 1)
}

func TestTransfersController_TransferZeroAddressError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Amount:             *assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestTransfersController_TransferBalanceToLowError(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()

	ethClient.On("PendingNonceAt", mock.Anything, key.Address.Address()).Return(uint64(1), nil)
	ethClient.On("BalanceAt", mock.Anything, key.Address.Address(), (*big.Int)(nil)).Return(utils.NewBigI(1).ToInt(), nil)

	app := cltest.NewApplicationWithKey(t, ethClient, key)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	request := models.SendEtherRequest{
		FromAddress:        key.Address.Address(),
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             *assets.NewEth(100),
		AllowHigherAmounts: false,
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestTransfersController_JSONBindingError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer([]byte(`{"address":""}`)))
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}
