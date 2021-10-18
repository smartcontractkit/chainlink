package web_test

import (
	"bytes"
	"encoding/json"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/require"
)

func Test_ChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupChainsControllerTest(t)

	body, err := json.Marshal(web.CreateChainRequest{
		ID: utils.Big{},
		Config: types.ChainCfg{
			BlockHistoryEstimatorBlockDelay:       null.IntFrom(1),
			BlockHistoryEstimatorBlockHistorySize: null.IntFrom(12),
			EvmEIP1559DynamicFees:                 null.BoolFrom(false),
			MinIncomingConfirmations:              null.IntFrom(10),
		},
	})
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/chains/evm", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	chainSet := controller.app.GetChainSet()

	chains := chainSet.Chains()
	require.Len(t, chains, 1)

	chain := chains[0]

	resource := presenters.ChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, chain.ID().String())
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockDelay, chain.Config().BlockHistoryEstimatorBlockDelay())
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockHistorySize, chain.Config().BlockHistoryEstimatorBlockHistorySize())
	assert.Equal(t, resource.Config.EvmEIP1559DynamicFees, chain.Config().EvmEIP1559DynamicFees())
	assert.Equal(t, resource.Config.MinIncomingConfirmations, chain.Config().MinIncomingConfirmations())

}

type TestChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupChainsControllerTest(t *testing.T) *TestChainsController {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	return &TestChainsController{
		app:    app,
		client: client,
	}
}
