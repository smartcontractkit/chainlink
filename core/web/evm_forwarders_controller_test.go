package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type TestEVMForwardersController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMForwardersControllerTest(t *testing.T) *TestEVMForwardersController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()

	return &TestEVMForwardersController{
		app:    app,
		client: client,
	}
}

func Test_EVMForwardersController_Create(t *testing.T) {
	t.Parallel()

	controller := setupEVMForwardersControllerTest(t)

	// Setting up chain
	chainId := testutils.NewRandomEVMChainID()
	chainSet := controller.app.GetChains().EVM
	dbChain, err := chainSet.ORM().CreateChain(utils.Big(*chainId), nil)
	require.NoError(t, err)

	// Build EVMForwarderRequest
	address := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")
	body, err := json.Marshal(web.CreateEVMForwarderRequest{
		EVMChainID: &dbChain.ID,
		Address:    address,
	},
	)
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/nodes/evm/forwarders", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resource := presenters.EVMForwarderResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.Address, address)
}

func Test_EVMForwardersController_Index(t *testing.T) {
	t.Parallel()

	controller := setupEVMForwardersControllerTest(t)

	// Setting up chain
	chainId := testutils.NewRandomEVMChainID()
	chainSet := controller.app.GetChains().EVM
	dbChain, err := chainSet.ORM().CreateChain(utils.Big(*chainId), nil)
	require.NoError(t, err)

	// Build EVMForwarderRequest
	fwdrs := []web.CreateEVMForwarderRequest{
		{
			EVMChainID: &dbChain.ID,
			Address:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
		},
		{
			EVMChainID: &dbChain.ID,
			Address:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d82"),
		},
	}
	for _, fwdr := range fwdrs {

		body, err := json.Marshal(web.CreateEVMForwarderRequest{
			EVMChainID: &dbChain.ID,
			Address:    fwdr.Address,
		},
		)
		require.NoError(t, err)

		resp, cleanup := controller.client.Post("/v2/nodes/evm/forwarders", bytes.NewReader(body))
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	resp, cleanup := controller.client.Get("/v2/nodes/evm/forwarders?size=2")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, len(fwdrs), metaCount)

	var links jsonapi.Links

	var fwdrcs []presenters.EVMForwarderResource
	err = web.ParsePaginatedResponse(body, &fwdrcs, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["prev"].Href)
}
