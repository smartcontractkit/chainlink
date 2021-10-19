package web_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

			resp2, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/%s", newHexChainId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp2.StatusCode)

			if tc.want != nil {
				resource1 := presenters.ChainResource{}
				resource2 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)
				err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp2), &resource2)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, dbChain.ID.String())
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockDelay, dbChain.Cfg.BlockHistoryEstimatorBlockDelay)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockHistorySize, dbChain.Cfg.BlockHistoryEstimatorBlockHistorySize)
				assert.Equal(t, resource1.Config.EvmEIP1559DynamicFees, dbChain.Cfg.EvmEIP1559DynamicFees)
				assert.Equal(t, resource1.Config.MinIncomingConfirmations, dbChain.Cfg.MinIncomingConfirmations)
				assert.Equal(t, resource2.ID, dbChain.ID.String())
				assert.Equal(t, resource2.Config.BlockHistoryEstimatorBlockDelay, dbChain.Cfg.BlockHistoryEstimatorBlockDelay)
				assert.Equal(t, resource2.Config.BlockHistoryEstimatorBlockHistorySize, dbChain.Cfg.BlockHistoryEstimatorBlockHistorySize)
				assert.Equal(t, resource2.Config.EvmEIP1559DynamicFees, dbChain.Cfg.EvmEIP1559DynamicFees)
				assert.Equal(t, resource2.Config.MinIncomingConfirmations, dbChain.Cfg.MinIncomingConfirmations)
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
