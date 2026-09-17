package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type TestEVMForwardersController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMForwardersControllerTest(t *testing.T, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) *TestEVMForwardersController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplicationWithConfig(t, configtest.NewGeneralConfig(t, overrideFn))
	ctx := t.Context()
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)

	return &TestEVMForwardersController{
		app:    app,
		client: client,
	}
}

func Test_EVMForwardersController_Track(t *testing.T) {
	t.Parallel()

	chainID := sqlutil.New(testutils.NewRandomEVMChainID())
	controller := setupEVMForwardersControllerTest(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = toml.EVMConfigs{
			{ChainID: chainID, Enabled: new(true), Chain: toml.Defaults(chainID)},
		}
	})

	// Build EVMForwarderRequest
	address := utils.RandomAddress()
	body, err := json.Marshal(web.TrackEVMForwarderRequest{
		EVMChainID: chainID,
		Address:    address,
	},
	)
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/nodes/evm/forwarders/track", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resource := presenters.EVMForwarderResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.Address, address)

	require.Len(t, controller.app.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkEVM)).Slice(), 1)

	resp, cleanup = controller.client.Delete("/v2/nodes/evm/forwarders/" + resource.ID)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.NoError(t, err)
}

func Test_EVMForwardersController_Index(t *testing.T) {
	t.Parallel()

	chainID := sqlutil.New(testutils.NewRandomEVMChainID())
	controller := setupEVMForwardersControllerTest(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = toml.EVMConfigs{
			{ChainID: chainID, Enabled: new(true), Chain: toml.Defaults(chainID)},
		}
	})

	// Build EVMForwarderRequest
	fwdrs := []web.TrackEVMForwarderRequest{
		{
			EVMChainID: chainID,
			Address:    utils.RandomAddress(),
		},
		{
			EVMChainID: chainID,
			Address:    utils.RandomAddress(),
		},
	}
	for _, fwdr := range fwdrs {
		body, err := json.Marshal(web.TrackEVMForwarderRequest{
			EVMChainID: chainID,
			Address:    fwdr.Address,
		},
		)
		require.NoError(t, err)

		resp, cleanup := controller.client.Post("/v2/nodes/evm/forwarders/track", bytes.NewReader(body))
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
	require.NoError(t, err)
	assert.Empty(t, links["prev"].Href)
}
