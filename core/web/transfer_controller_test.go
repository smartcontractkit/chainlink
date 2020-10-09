package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransfersController_CreateSuccess_From(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	client := app.NewHTTPClient()
	require.NoError(t, app.StartAndConnect())

	sendKeys, err := app.GetStore().SendKeys()
	from := common.HexToAddress(string(sendKeys[0].Address))
	require.NoError(t, err)

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        from,
		Amount:             *assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	count, err := app.GetStore().CountOf(models.EthTx{})
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestTransfersController_TransferError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Amount:             *assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestTransfersController_JSONBindingError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	client := app.NewHTTPClient()
	assert.NoError(t, app.StartAndConnect())
	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer([]byte(`{"address":""}`)))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}
