package web_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func Test_CosmosChainsController_Show(t *testing.T) {
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
						FallbackGasPriceUAtom: null.StringFrom("9.999"),
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

			controller := setupCosmosChainsControllerTestV2(t, &cosmos.CosmosConfig{ChainID: ptr(validId), Enabled: ptr(true),
				Chain: coscfg.Chain{
					FallbackGasPriceUAtom: ptr(decimal.RequireFromString("9.999")),
					GasLimitMultiplier:    ptr(decimal.RequireFromString("1.55555")),
				}})

			wantedResult := tc.want(t, controller.app)
			resp, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/cosmos/%s", tc.inputId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.CosmosChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, wantedResult.ID)
				assert.Equal(t, resource1.Config.FallbackGasPriceUAtom, wantedResult.Cfg.FallbackGasPriceUAtom)
				assert.Equal(t, resource1.Config.GasLimitMultiplier, wantedResult.Cfg.GasLimitMultiplier)
			}
		})
	}
}

func Test_CosmosChainsController_Index(t *testing.T) {
	t.Parallel()

	chainA := &cosmos.CosmosConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestA-%d", rand.Int31n(999999))),
		Enabled: ptr(true),
		Chain: coscfg.Chain{
			FallbackGasPriceUAtom: ptr(decimal.RequireFromString("9.999")),
		},
	}
	chainB := &cosmos.CosmosConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestB-%d", rand.Int31n(999999))),
		Enabled: ptr(true),
		Chain: coscfg.Chain{
			GasLimitMultiplier: ptr(decimal.RequireFromString("1.55555")),
		},
	}
	controller := setupCosmosChainsControllerTestV2(t, chainA, chainB)

	badResp, cleanup := controller.client.Get("/v2/chains/cosmos?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/cosmos?size=1")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

	var links jsonapi.Links

	var chains []presenters.CosmosChainResource
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainA.ChainID, chains[0].ID)
	assert.Equal(t, chainA.Chain.FallbackGasPriceUAtom.String(), chains[0].Config.FallbackGasPriceUAtom.String)
	assert.Equal(t, chainA.Chain.GasLimitMultiplier.InexactFloat64(), chains[0].Config.GasLimitMultiplier.Float64)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.CosmosChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainB.ChainID, chains[0].ID)
	assert.Equal(t, chainB.Chain.FallbackGasPriceUAtom.String(), chains[0].Config.FallbackGasPriceUAtom.String)
	assert.Equal(t, chainB.Chain.GasLimitMultiplier.InexactFloat64(), chains[0].Config.GasLimitMultiplier.Float64)
}

type TestCosmosChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupCosmosChainsControllerTestV2(t *testing.T, cfgs ...*cosmos.CosmosConfig) *TestCosmosChainsController {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cfgs
		c.EVM = nil
	})
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestCosmosChainsController{
		app:    app,
		client: client,
	}
}
