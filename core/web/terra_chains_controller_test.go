package web_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_TerraChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupTerraChainsControllerTest(t)

	newChainId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))

	minute, err := utils.NewDuration(time.Minute)
	require.NoError(t, err)
	body, err := json.Marshal(web.NewCreateChainRequest(
		newChainId,
		&db.ChainCfg{
			BlocksUntilTxTimeout:  null.IntFrom(1),
			ConfirmPollPeriod:     &minute,
			FallbackGasPriceULuna: null.StringFrom("9.999"),
			GasLimitMultiplier:    null.FloatFrom(1.55555),
			MaxMsgsPerBatch:       null.IntFrom(10),
		}))
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/chains/terra", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	chainSet := controller.app.GetChains().Terra
	dbChain, err := chainSet.Show(newChainId)
	require.NoError(t, err)

	resource := presenters.TerraChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, dbChain.ID)
	assert.Equal(t, resource.Config.BlocksUntilTxTimeout, dbChain.Cfg.BlocksUntilTxTimeout)
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
		want           func(t *testing.T, app *cltest.TestApplication) *db.Chain
	}{
		{
			inputId: validId,
			name:    "success",
			want: func(t *testing.T, app *cltest.TestApplication) *db.Chain {
				chain := db.Chain{
					ID:      validId,
					Enabled: true,
					Cfg: db.ChainCfg{
						FallbackGasPriceULuna: null.StringFrom("9.999"),
						GasLimitMultiplier:    null.FloatFrom(1.55555),
					},
				}

				return &chain
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId: "234",
			name:    "not found",
			want: func(t *testing.T, app *cltest.TestApplication) *db.Chain {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupTerraChainsControllerTestV2(t, &terra.TerraConfig{ChainID: ptr(validId), Enabled: ptr(true),
				Chain: tercfg.Chain{
					FallbackGasPriceULuna: ptr(decimal.RequireFromString("9.999")),
					GasLimitMultiplier:    ptr(decimal.RequireFromString("1.55555")),
				}})

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

	chainA := &terra.TerraConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestA-%d", rand.Int31n(999999))),
		Enabled: ptr(true),
		Chain: tercfg.Chain{
			FallbackGasPriceULuna: ptr(decimal.RequireFromString("9.999")),
		},
	}
	chainB := &terra.TerraConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestB-%d", rand.Int31n(999999))),
		Enabled: ptr(true),
		Chain: tercfg.Chain{
			GasLimitMultiplier: ptr(decimal.RequireFromString("1.55555")),
		},
	}
	controller := setupTerraChainsControllerTestV2(t, chainA, chainB)

	badResp, cleanup := controller.client.Get("/v2/chains/terra?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/terra?size=1")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

	var links jsonapi.Links

	var chains []presenters.TerraChainResource
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainA.ChainID, chains[0].ID)
	assert.Equal(t, chainA.Chain.FallbackGasPriceULuna.String(), chains[0].Config.FallbackGasPriceULuna.String)
	assert.Equal(t, chainA.Chain.GasLimitMultiplier.InexactFloat64(), chains[0].Config.GasLimitMultiplier.Float64)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.TerraChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainB.ChainID, chains[0].ID)
	assert.Equal(t, chainB.Chain.FallbackGasPriceULuna.String(), chains[0].Config.FallbackGasPriceULuna.String)
	assert.Equal(t, chainB.Chain.GasLimitMultiplier.InexactFloat64(), chains[0].Config.GasLimitMultiplier.Float64)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_TerraChainsController_Update(t *testing.T) {
	t.Parallel()

	chainUpdate := web.UpdateChainRequest[*db.ChainCfg]{
		Enabled: true,
		Config: &db.ChainCfg{
			FallbackGasPriceULuna: null.StringFrom("9.999"),
			GasLimitMultiplier:    null.FloatFrom(1.55555),
		},
	}

	validId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))

	testCases := []struct {
		name              string
		inputId           string
		wantStatusCode    int
		chainBeforeUpdate func(t *testing.T, app *cltest.TestApplication) *db.Chain
	}{
		{
			inputId: validId,
			name:    "success",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *db.Chain {
				newChainConfig := db.ChainCfg{
					FallbackGasPriceULuna: null.StringFrom("9.999"),
					GasLimitMultiplier:    null.FloatFrom(1.55555),
				}

				chain := db.Chain{
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
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *db.Chain {
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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_TerraChainsController_Delete(t *testing.T) {
	t.Parallel()

	controller := setupTerraChainsControllerTest(t)

	newChainConfig := db.ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("9.999"),
		GasLimitMultiplier:    null.FloatFrom(1.55555),
	}

	chainId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	chain := db.Chain{
		ID:      chainId,
		Enabled: true,
		Cfg:     newChainConfig,
	}
	terratest.MustInsertChain(t, controller.app.GetSqlxDB(), &chain)

	_, countBefore, err := controller.app.Chains.Terra.Index(0, 10)
	require.NoError(t, err)
	require.Equal(t, 1, countBefore)

	t.Run("non-existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete("/v2/chains/terra/121231")
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		_, countAfter, err := controller.app.Chains.Terra.Index(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, countAfter)
	})

	t.Run("existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete(
			fmt.Sprintf("/v2/chains/terra/%s", chain.ID),
		)
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, countAfter, err := controller.app.Chains.Terra.Index(0, 10)
		require.NoError(t, err)
		require.Equal(t, 0, countAfter)

		_, err = controller.app.Chains.Terra.Show(chain.ID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

type TestTerraChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func setupTerraChainsControllerTest(t *testing.T) *TestTerraChainsController {
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.TerraEnabled = null.BoolFrom(true)
	cfg.Overrides.EVMEnabled = null.BoolFrom(false)
	cfg.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestTerraChainsController{
		app:    app,
		client: client,
	}
}

func setupTerraChainsControllerTestV2(t *testing.T, cfgs ...*terra.TerraConfig) *TestTerraChainsController {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Terra = cfgs
		c.EVM = nil
	})
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestTerraChainsController{
		app:    app,
		client: client,
	}
}
