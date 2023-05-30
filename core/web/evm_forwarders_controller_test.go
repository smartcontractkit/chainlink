package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type TestEVMForwardersController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMForwardersControllerTest(t *testing.T, lp logpoller.LogPoller, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) *TestEVMForwardersController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	mockLogPoller := mocklogpoller.NewLogPoller(t)
	mockLogPoller.On("Name").Return("mock logpoller")
	mockLogPoller.On("Start", mock.Anything).Return(nil)
	mockLogPoller.On("Ready").Return(nil)
	mockLogPoller.On("HealthReport").Return(nil)
	mockLogPoller.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	mockLogPoller.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil).Maybe()
	mockLogPoller.On("Close", mock.Anything).Return(nil)

	app := cltest.NewApplicationWithConfig(t, configtest.NewGeneralConfig(t, overrideFn), mockLogPoller)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestEVMForwardersController{
		app:    app,
		client: client,
	}
}

func Test_EVMForwardersController_Track(t *testing.T) {
	t.Parallel()

	chainId := utils.NewBig(testutils.NewRandomEVMChainID())
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := logpoller.NewORM(chainId.ToInt(), db, lggr, pgtest.NewQConfig(true))
	ec := evmtest.NewEthClientMock(t)
	lp := logpoller.NewLogPoller(orm, ec, lggr, 0, 0, 0, 0, 0)

	controller := setupEVMForwardersControllerTest(t, lp, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = ptr(true)
		c.EVM = evmcfg.EVMConfigs{
			{ChainID: chainId, Enabled: ptr(true), Chain: evmcfg.Defaults(chainId)},
		}
	})

	// Build EVMForwarderRequest
	address := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")
	body, err := json.Marshal(web.TrackEVMForwarderRequest{
		EVMChainID: chainId,
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

	require.Len(t, controller.app.Chains.EVM.Chains(), 1)

	resp, cleanup = controller.client.Delete("/v2/nodes/evm/forwarders/" + resource.ID)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.NoError(t, err)
}

func Test_EVMForwardersController_Index(t *testing.T) {
	t.Parallel()

	chainId := utils.NewBig(testutils.NewRandomEVMChainID())
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := logpoller.NewORM(chainId.ToInt(), db, lggr, pgtest.NewQConfig(true))
	lp := logpoller.NewLogPoller(orm, nil, lggr, 0, 0, 0, 0, 0)
	controller := setupEVMForwardersControllerTest(t, lp, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = ptr(true)
		c.EVM = evmcfg.EVMConfigs{
			{ChainID: chainId, Enabled: ptr(true), Chain: evmcfg.Defaults(chainId)},
		}
	})

	// Build EVMForwarderRequest
	fwdrs := []web.TrackEVMForwarderRequest{
		{
			EVMChainID: chainId,
			Address:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
		},
		{
			EVMChainID: chainId,
			Address:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d82"),
		},
	}
	for _, fwdr := range fwdrs {

		body, err := json.Marshal(web.TrackEVMForwarderRequest{
			EVMChainID: chainId,
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
	assert.NoError(t, err)
	assert.Empty(t, links["prev"].Href)
}
