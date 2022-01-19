package web_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func Test_TerraChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupTerraChainsControllerTest(t)

	const newChainId = "Chainlinktest-42"

	minute := models.MustMakeDuration(time.Minute)
	body, err := json.Marshal(web.CreateTerraChainRequest{
		ID: newChainId,
		Config: types.ChainCfg{
			BlocksUntilTxTimeout:  null.IntFrom(1),
			ConfirmMaxPolls:       null.IntFrom(10),
			ConfirmPollPeriod:     &minute,
			FallbackGasPriceULuna: null.StringFrom("9.999"),
			GasLimitMultiplier:    null.FloatFrom(1.55555),
			MaxMsgsPerBatch:       null.IntFrom(10),
		},
	})
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/chains/terra", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	chainSet := controller.app.GetChains().Terra
	dbChain, err := chainSet.ORM().Chain(newChainId)
	require.NoError(t, err)

	resource := presenters.TerraChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, dbChain.ID)
	assert.Equal(t, resource.Config.BlocksUntilTxTimeout, dbChain.Cfg.BlocksUntilTxTimeout)
	assert.Equal(t, resource.Config.ConfirmMaxPolls, dbChain.Cfg.ConfirmMaxPolls)
	assert.Equal(t, resource.Config.ConfirmPollPeriod, dbChain.Cfg.ConfirmPollPeriod)
	assert.Equal(t, resource.Config.FallbackGasPriceULuna, dbChain.Cfg.FallbackGasPriceULuna)
	assert.Equal(t, resource.Config.GasLimitMultiplier, dbChain.Cfg.GasLimitMultiplier)
	assert.Equal(t, resource.Config.MaxMsgsPerBatch, dbChain.Cfg.MaxMsgsPerBatch)
}

func Test_TerraChainsController_Show(t *testing.T) {
	t.Parallel()

	const validId = "Chainlink-12"

	testCases := []struct {
		name           string
		inputId        string
		wantStatusCode int
		want           func(t *testing.T, app *cltest.TestApplication) *types.Chain
	}{
		{
			inputId: validId,
			name:    "success",
			want: func(t *testing.T, app *cltest.TestApplication) *types.Chain {
				newChainConfig := types.ChainCfg{
					FallbackGasPriceULuna: null.StringFrom("9.999"),
					GasLimitMultiplier:    null.FloatFrom(1.55555),
				}

				chain := types.Chain{
					ID:      validId,
					Enabled: true,
					Cfg:     newChainConfig,
				}
				terratest.MustInsertChain(t, app.GetSqlxDB(), &chain)

				return &chain
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId: "234",
			name:    "not found",
			want: func(t *testing.T, app *cltest.TestApplication) *types.Chain {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupTerraChainsControllerTest(t)

			wantedResult := tc.want(t, controller.app)
			resp, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/terra/%s", tc.inputId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.TerraChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, wantedResult.ID)
				assert.Equal(t, resource1.Config.FallbackGasPriceULuna, wantedResult.Cfg.FallbackGasPriceULuna)
				assert.Equal(t, resource1.Config.GasLimitMultiplier, wantedResult.Cfg.GasLimitMultiplier)
			}
		})
	}
}

func Test_TerraChainsController_Index(t *testing.T) {
	t.Parallel()

	controller := setupTerraChainsControllerTest(t)

	newChains := []web.CreateTerraChainRequest{
		{
			ID: "Chainlinktest-24",
			Config: types.ChainCfg{
				FallbackGasPriceULuna: null.StringFrom("9.999"),
			},
		},
		{
			ID: "Chainlinktest-30",
			Config: types.ChainCfg{
				GasLimitMultiplier: null.FloatFrom(1.55555),
			},
		},
	}

	for _, newChain := range newChains {
		ch := newChain
		terratest.MustInsertChain(t, controller.app.GetSqlxDB(), &types.Chain{
			ID:      ch.ID,
			Enabled: true,
			Cfg:     ch.Config,
		})
	}

	badResp, cleanup := controller.client.Get("/v2/chains/terra?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/terra?size=1")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, len(newChains), metaCount)

	var links jsonapi.Links

	chains := []presenters.TerraChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[0].ID, chains[0].ID)
	assert.Equal(t, newChains[0].Config.FallbackGasPriceULuna, chains[0].Config.FallbackGasPriceULuna)
	assert.Equal(t, newChains[0].Config.GasLimitMultiplier, chains[0].Config.GasLimitMultiplier)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.TerraChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[1].ID, chains[0].ID)
	assert.Equal(t, newChains[1].Config.FallbackGasPriceULuna, chains[0].Config.FallbackGasPriceULuna)
	assert.Equal(t, newChains[1].Config.GasLimitMultiplier, chains[0].Config.GasLimitMultiplier)
}

func Test_TerraChainsController_Update(t *testing.T) {
	t.Parallel()

	chainUpdate := web.UpdateTerraChainRequest{
		Enabled: true,
		Config: types.ChainCfg{
			FallbackGasPriceULuna: null.StringFrom("9.999"),
			GasLimitMultiplier:    null.FloatFrom(1.55555),
		},
	}

	const validId = "Chainlinktest-12"

	testCases := []struct {
		name              string
		inputId           string
		wantStatusCode    int
		chainBeforeUpdate func(t *testing.T, app *cltest.TestApplication) *types.Chain
	}{
		{
			inputId: validId,
			name:    "success",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *types.Chain {
				newChainConfig := types.ChainCfg{
					FallbackGasPriceULuna: null.StringFrom("9.999"),
					GasLimitMultiplier:    null.FloatFrom(1.55555),
				}

				chain := types.Chain{
					ID:      validId,
					Enabled: true,
					Cfg:     newChainConfig,
				}
				terratest.MustInsertChain(t, app.GetSqlxDB(), &chain)

				return &chain
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId: "341212",
			name:    "not found",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *types.Chain {
				return nil
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupTerraChainsControllerTest(t)

			beforeUpdate := tc.chainBeforeUpdate(t, controller.app)

			body, err := json.Marshal(chainUpdate)
			require.NoError(t, err)

			resp, cleanup := controller.client.Patch(
				fmt.Sprintf("/v2/chains/terra/%s", tc.inputId),
				bytes.NewReader(body),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if beforeUpdate != nil {
				resource1 := presenters.TerraChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, beforeUpdate.ID)
				assert.Equal(t, resource1.Enabled, chainUpdate.Enabled)
				assert.Equal(t, resource1.Config.FallbackGasPriceULuna, chainUpdate.Config.FallbackGasPriceULuna)
				assert.Equal(t, resource1.Config.GasLimitMultiplier, chainUpdate.Config.GasLimitMultiplier)
			}
		})
	}
}

func Test_TerraChainsController_Delete(t *testing.T) {
	t.Parallel()

	controller := setupTerraChainsControllerTest(t)

	newChainConfig := types.ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("9.999"),
		GasLimitMultiplier:    null.FloatFrom(1.55555),
	}

	const chainId = "Chainlinktest-50"
	chain := types.Chain{
		ID:      chainId,
		Enabled: true,
		Cfg:     newChainConfig,
	}
	terratest.MustInsertChain(t, controller.app.GetSqlxDB(), &chain)

	_, countBefore, err := controller.app.TerraORM().Chains(0, 10)
	require.NoError(t, err)
	require.Equal(t, 1, countBefore)

	t.Run("non-existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete("/v2/chains/terra/121231")
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		_, countAfter, err := controller.app.TerraORM().Chains(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, countAfter)
	})

	t.Run("existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete(
			fmt.Sprintf("/v2/chains/terra/%s", chain.ID),
		)
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, countAfter, err := controller.app.TerraORM().Chains(0, 10)
		require.NoError(t, err)
		require.Equal(t, 0, countAfter)

		_, err = controller.app.TerraORM().Chain(chain.ID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

type TestTerraChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupTerraChainsControllerTest(t *testing.T) *TestTerraChainsController {
	// Using this instead of `NewApplicationTerraDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplication(t)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	return &TestTerraChainsController{
		app:    app,
		client: client,
	}
}
