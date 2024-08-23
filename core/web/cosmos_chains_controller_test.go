package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func Test_CosmosChainsController_Show(t *testing.T) {
	t.Parallel()

	const validId = "Chainlink-12"

	testCases := []struct {
		name           string
		inputId        string
		wantStatusCode int
		want           func(t *testing.T, app *cltest.TestApplication) *types.ChainStatus
	}{
		{
			inputId: validId,
			name:    "success",
			want: func(t *testing.T, app *cltest.TestApplication) *types.ChainStatus {
				return &types.ChainStatus{
					ID:      validId,
					Enabled: true,
					Config: `ChainID = 'Chainlink-12'
Enabled = true
Bech32Prefix = 'wasm'
BlockRate = '6s'
BlocksUntilTxTimeout = 30
ConfirmPollPeriod = '1s'
FallbackGasPrice = '9.999'
GasToken = 'ucosm'
GasLimitMultiplier = '1.55555'
MaxMsgsPerBatch = 100
OCR2CachePollPeriod = '4s'
OCR2CacheTTL = '1m0s'
TxMsgTimeout = '10m0s'
Nodes = []
`,
				}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId: "234",
			name:    "not found",
			want: func(t *testing.T, app *cltest.TestApplication) *types.ChainStatus {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupCosmosChainsControllerTestV2(t, &coscfg.TOMLConfig{
				ChainID: ptr(validId),
				Enabled: ptr(true),
				Chain: coscfg.Chain{
					FallbackGasPrice:   ptr(decimal.RequireFromString("9.999")),
					GasLimitMultiplier: ptr(decimal.RequireFromString("1.55555")),
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

				assert.Equal(t, wantedResult.ID, resource1.ID)
				assert.Equal(t, wantedResult.Config, resource1.Config)
			}
		})
	}
}

func Test_CosmosChainsController_Index(t *testing.T) {
	t.Parallel()

	chainA := &coscfg.TOMLConfig{
		ChainID: ptr("a" + cosmostest.RandomChainID()),
		Enabled: ptr(true),
		Chain: coscfg.Chain{
			FallbackGasPrice: ptr(decimal.RequireFromString("9.999")),
		},
	}

	chainB := &coscfg.TOMLConfig{
		ChainID: ptr("b" + cosmostest.RandomChainID()),
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
	tomlA, err := chainA.TOMLString()
	require.NoError(t, err)
	assert.Equal(t, tomlA, chains[0].Config)

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
	tomlB, err := chainB.TOMLString()
	require.NoError(t, err)
	assert.Equal(t, tomlB, chains[0].Config)
}

type TestCosmosChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupCosmosChainsControllerTestV2(t *testing.T, cfgs ...*coscfg.TOMLConfig) *TestCosmosChainsController {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cfgs
		c.EVM = nil
	})
	app := cltest.NewApplicationWithConfig(t, cfg)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)

	return &TestCosmosChainsController{
		app:    app,
		client: client,
	}
}
