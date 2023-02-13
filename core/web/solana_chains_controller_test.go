package web_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func Test_SolanaChainsController_Show(t *testing.T) {
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
						SkipPreflight: null.BoolFrom(false),
						TxTimeout:     utils.MustNewDuration(time.Hour),
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

			controller := setupSolanaChainsControllerTestV2(t, &solana.SolanaConfig{ChainID: ptr(validId),
				Chain: config.Chain{
					SkipPreflight: ptr(false),
					TxTimeout:     utils.MustNewDuration(time.Hour),
				},
			})

			wantedResult := tc.want(t, controller.app)
			resp, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/solana/%s", tc.inputId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.SolanaChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, wantedResult.ID)
				assert.Equal(t, resource1.Config.SkipPreflight, wantedResult.Cfg.SkipPreflight)
				assert.Equal(t, resource1.Config.TxTimeout, wantedResult.Cfg.TxTimeout)
			}
		})
	}
}

func Test_SolanaChainsController_Index(t *testing.T) {
	t.Parallel()

	chainA := &solana.SolanaConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestA-%d", rand.Int31n(999999))),
		Chain: config.Chain{
			TxTimeout: utils.MustNewDuration(time.Hour),
		},
	}
	chainB := &solana.SolanaConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestB-%d", rand.Int31n(999999))),
		Chain: config.Chain{
			SkipPreflight: ptr(false),
		},
	}
	controller := setupSolanaChainsControllerTestV2(t, chainA, chainB)

	badResp, cleanup := controller.client.Get("/v2/chains/solana?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/solana?size=1")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

	var links jsonapi.Links

	chains := []presenters.SolanaChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainA.ChainID, chains[0].ID)
	assert.Equal(t, *chainA.Chain.SkipPreflight, chains[0].Config.SkipPreflight.Bool)
	assert.Equal(t, chainA.Chain.TxTimeout.Duration(), chains[0].Config.TxTimeout.Duration())

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.SolanaChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainB.ChainID, chains[0].ID)
	assert.Equal(t, *chainB.Chain.SkipPreflight, chains[0].Config.SkipPreflight.Bool)
	assert.Equal(t, chainB.Chain.TxTimeout.Duration(), chains[0].Config.TxTimeout.Duration())
}

type TestSolanaChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupSolanaChainsControllerTestV2(t *testing.T, cfgs ...*solana.SolanaConfig) *TestSolanaChainsController {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = cfgs
		c.EVM = nil
	})
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestSolanaChainsController{
		app:    app,
		client: client,
	}
}
