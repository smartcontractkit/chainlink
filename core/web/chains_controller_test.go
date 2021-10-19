package web_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/manyminds/api2go/jsonapi"
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

	newChainId := *utils.NewBigI(42)

	body, err := json.Marshal(web.CreateChainRequest{
		ID: newChainId,
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
	dbChain, err := chainSet.ORM().Chain(newChainId)
	require.NoError(t, err)

	resource := presenters.ChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, dbChain.ID.String())
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockDelay, dbChain.Cfg.BlockHistoryEstimatorBlockDelay)
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockHistorySize, dbChain.Cfg.BlockHistoryEstimatorBlockHistorySize)
	assert.Equal(t, resource.Config.EvmEIP1559DynamicFees, dbChain.Cfg.EvmEIP1559DynamicFees)
	assert.Equal(t, resource.Config.MinIncomingConfirmations, dbChain.Cfg.MinIncomingConfirmations)
}

func Test_ChainsController_Show(t *testing.T) {
	t.Parallel()

	var dbChain types.Chain

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *types.Chain
		wantStatusCode int
	}{
		{
			name: "success",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				newChainConfig := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay:       null.IntFrom(23),
					BlockHistoryEstimatorBlockHistorySize: null.IntFrom(50),
					EvmEIP1559DynamicFees:                 null.BoolFrom(true),
					MinIncomingConfirmations:              null.IntFrom(12),
				}
				chainId := utils.HexToBig(*id)
				chain, err := app.GetChainSet().Add(chainId, newChainConfig)
				require.NoError(t, err)

				dbChain = chain
				*id = chainId.String()
			},
			want:           &dbChain,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid id",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "invalidid"
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "not found",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "234"
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupChainsControllerTest(t)

			newHexChainId := hex.EncodeToString([]byte("1"))
			if tc.before != nil {
				tc.before(t, controller.app, &newHexChainId)
			}

			resp, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/evm/%s", newHexChainId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resource1 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, tc.want.ID.String())
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockDelay, tc.want.Cfg.BlockHistoryEstimatorBlockDelay)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockHistorySize, tc.want.Cfg.BlockHistoryEstimatorBlockHistorySize)
				assert.Equal(t, resource1.Config.EvmEIP1559DynamicFees, tc.want.Cfg.EvmEIP1559DynamicFees)
				assert.Equal(t, resource1.Config.MinIncomingConfirmations, tc.want.Cfg.MinIncomingConfirmations)
			}
		})
	}
}

func Test_ChainsController_Index(t *testing.T) {
	t.Parallel()

	controller := setupChainsControllerTest(t)

	newChains := []web.CreateChainRequest{
		{
			ID: *utils.NewBigI(30),
			Config: types.ChainCfg{
				BlockHistoryEstimatorBlockDelay:       null.IntFrom(5),
				BlockHistoryEstimatorBlockHistorySize: null.IntFrom(2),
				EvmEIP1559DynamicFees:                 null.BoolFrom(false),
				MinIncomingConfirmations:              null.IntFrom(30),
			},
		},
		{
			ID: *utils.NewBigI(24),
			Config: types.ChainCfg{
				BlockHistoryEstimatorBlockDelay:       null.IntFrom(13),
				BlockHistoryEstimatorBlockHistorySize: null.IntFrom(1),
				EvmEIP1559DynamicFees:                 null.BoolFrom(true),
				MinIncomingConfirmations:              null.IntFrom(120),
			},
		},
	}

	for _, newChain := range newChains {
		_, err := controller.app.GetChainSet().Add(newChain.ID.ToInt(), newChain.Config)
		require.NoError(t, err)
	}

	badResp, cleanup := controller.client.Get("/v2/chains/evm?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/evm?size=3")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	// Apparently there are 2 chains by default...
	require.Equal(t, 4, metaCount)

	var links jsonapi.Links

	chains := []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[1].ID.String(), chains[2].ID)
	assert.Equal(t, newChains[1].Config.BlockHistoryEstimatorBlockDelay, chains[2].Config.BlockHistoryEstimatorBlockDelay)
	assert.Equal(t, newChains[1].Config.BlockHistoryEstimatorBlockHistorySize, chains[2].Config.BlockHistoryEstimatorBlockHistorySize)
	assert.Equal(t, newChains[1].Config.EvmEIP1559DynamicFees, chains[2].Config.EvmEIP1559DynamicFees)
	assert.Equal(t, newChains[1].Config.MinIncomingConfirmations, chains[2].Config.MinIncomingConfirmations)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[0].ID.String(), chains[0].ID)
	assert.Equal(t, newChains[0].Config.BlockHistoryEstimatorBlockDelay, chains[0].Config.BlockHistoryEstimatorBlockDelay)
	assert.Equal(t, newChains[0].Config.BlockHistoryEstimatorBlockHistorySize, chains[0].Config.BlockHistoryEstimatorBlockHistorySize)
	assert.Equal(t, newChains[0].Config.EvmEIP1559DynamicFees, chains[0].Config.EvmEIP1559DynamicFees)
	assert.Equal(t, newChains[0].Config.MinIncomingConfirmations, chains[0].Config.MinIncomingConfirmations)
}

func Test_ChainsController_Update(t *testing.T) {
	t.Parallel()

	var dbChain types.Chain

	chainUpdate := web.UpdateChainRequest{
		Enabled: true,
		Config: types.ChainCfg{
			BlockHistoryEstimatorBlockDelay:       null.IntFrom(55),
			BlockHistoryEstimatorBlockHistorySize: null.IntFrom(33),
			EvmEIP1559DynamicFees:                 null.BoolFrom(true),
			MinIncomingConfirmations:              null.IntFrom(100),
		},
	}

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *types.Chain
		wantStatusCode int
	}{
		{
			name: "success",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				newChainConfig := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay:       null.IntFrom(5),
					BlockHistoryEstimatorBlockHistorySize: null.IntFrom(2),
					EvmEIP1559DynamicFees:                 null.BoolFrom(false),
					MinIncomingConfirmations:              null.IntFrom(30),
				}

				chainId := utils.HexToBig(*id)
				chain, err := app.GetChainSet().Add(chainId, newChainConfig)
				require.NoError(t, err)

				dbChain = chain
				*id = chainId.String()
			},
			want:           &dbChain,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid id",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "invalidid"
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "not found",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "341212"
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupChainsControllerTest(t)

			newHexChainId := hex.EncodeToString([]byte("40"))
			if tc.before != nil {
				tc.before(t, controller.app, &newHexChainId)
			}

			body, err := json.Marshal(chainUpdate)
			require.NoError(t, err)

			resp, cleanup := controller.client.Patch(
				fmt.Sprintf("/v2/chains/evm/%s", newHexChainId),
				bytes.NewReader(body),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resource1 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, tc.want.ID.String())
				assert.Equal(t, resource1.Enabled, chainUpdate.Enabled)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockDelay, chainUpdate.Config.BlockHistoryEstimatorBlockDelay)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockHistorySize, chainUpdate.Config.BlockHistoryEstimatorBlockHistorySize)
				assert.Equal(t, resource1.Config.EvmEIP1559DynamicFees, chainUpdate.Config.EvmEIP1559DynamicFees)
				assert.Equal(t, resource1.Config.MinIncomingConfirmations, chainUpdate.Config.MinIncomingConfirmations)
			}
		})
	}
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
