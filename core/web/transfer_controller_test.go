package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransfersController_CreateSuccess(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
		ethMock.Register("eth_getBlockByNumber", models.BlockHeader{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(0))
		ethMock.Register("eth_chainId", *cltest.Int(config.ChainID()))
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	ethMock.AllCalled()
}

func TestTransfersController_CreateSuccess_From(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
		ethMock.Register("eth_getBlockByNumber", models.BlockHeader{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(0))
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
		ethMock.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
	})

	client := app.NewHTTPClient()

	require.NoError(t, app.StartAndConnect())

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		FromAddress:        app.Store.TxManager.NextActiveAccount().Address,
		Amount:             assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, errors.Errors, 0)

	ethMock.AllCalled()
}

func TestTransfersController_TransferError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
		ethMock.Register("eth_getBlockByNumber", models.BlockHeader{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(0))
		ethMock.Register("eth_chainId", *cltest.Int(config.ChainID()))
		ethMock.RegisterError("eth_sendRawTransaction", "No dice")
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	request := models.SendEtherRequest{
		DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
		Amount:             assets.NewEth(100),
	}

	body, err := json.Marshal(&request)
	assert.NoError(t, err)

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer(body))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 400)

	ethMock.AllCalled()
}

func TestTransfersController_JSONBindingError(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
		ethMock.Register("eth_getBlockByNumber", models.BlockHeader{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(0))
		ethMock.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
	})

	client := app.NewHTTPClient()

	assert.NoError(t, app.StartAndConnect())

	resp, cleanup := client.Post("/v2/transfers", bytes.NewBuffer([]byte(`{"address":""}`)))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, 400)

	ethMock.AllCalled()
}
