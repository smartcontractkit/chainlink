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

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func Test_SolanaChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupSolanaChainsControllerTest(t)

	newChainId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))

	second := models.MustMakeDuration(time.Second)
	minute := models.MustMakeDuration(time.Minute)
	hour := models.MustMakeDuration(time.Hour)
	body, err := json.Marshal(web.NewCreateChainRequest(
		newChainId,
		&db.ChainCfg{
			BalancePollPeriod:   &second,
			ConfirmPollPeriod:   &minute,
			OCR2CachePollPeriod: &minute,
			OCR2CacheTTL:        &second,
			TxTimeout:           &hour,
			SkipPreflight:       null.BoolFrom(false),
			Commitment:          null.StringFrom(string(rpc.CommitmentRecent)),
		}))

	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/chains/solana", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	chainSet := controller.app.GetChains().Solana
	dbChain, err := chainSet.Show(newChainId)
	require.NoError(t, err)

	resource := presenters.SolanaChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, dbChain.ID)
	assert.Equal(t, resource.Config.BalancePollPeriod, dbChain.Cfg.BalancePollPeriod)
	assert.Equal(t, resource.Config.ConfirmPollPeriod, dbChain.Cfg.ConfirmPollPeriod)
	assert.Equal(t, resource.Config.OCR2CachePollPeriod, dbChain.Cfg.OCR2CachePollPeriod)
	assert.Equal(t, resource.Config.OCR2CacheTTL, dbChain.Cfg.OCR2CacheTTL)
	assert.Equal(t, resource.Config.TxTimeout, dbChain.Cfg.TxTimeout)
	assert.Equal(t, resource.Config.SkipPreflight, dbChain.Cfg.SkipPreflight)
	assert.Equal(t, resource.Config.Commitment, dbChain.Cfg.Commitment)
}

func Test_SolanaChainsController_Show(t *testing.T) {
	t.Parallel()

	const validId = "Chainlink-12"

	hour := models.MustMakeDuration(time.Hour)
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
				newChainConfig := db.ChainCfg{
					SkipPreflight: null.BoolFrom(false),
					TxTimeout:     &hour,
				}

				chain := db.Chain{
					ID:      validId,
					Enabled: true,
					Cfg:     newChainConfig,
				}
				solanatest.MustInsertChain(t, app.GetSqlxDB(), &chain)

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

			controller := setupSolanaChainsControllerTest(t)

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

	controller := setupSolanaChainsControllerTest(t)

	hour := models.MustMakeDuration(time.Hour)
	newChains := []web.CreateChainRequest[string, *db.ChainCfg]{
		{
			ID: fmt.Sprintf("ChainlinktestA-%d", rand.Int31n(999999)),
			Config: &db.ChainCfg{
				TxTimeout: &hour,
			},
		},
		{
			ID: fmt.Sprintf("ChainlinktestB-%d", rand.Int31n(999999)),
			Config: &db.ChainCfg{
				SkipPreflight: null.BoolFrom(false),
			},
		},
	}

	for _, newChain := range newChains {
		ch := newChain
		solanatest.MustInsertChain(t, controller.app.GetSqlxDB(), &db.Chain{
			ID:      ch.ID,
			Enabled: true,
			Cfg:     *ch.Config,
		})
	}

	badResp, cleanup := controller.client.Get("/v2/chains/solana?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/solana?size=1")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, len(newChains), metaCount)

	var links jsonapi.Links

	chains := []presenters.SolanaChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[0].ID, chains[0].ID)
	assert.Equal(t, newChains[0].Config.SkipPreflight, chains[0].Config.SkipPreflight)
	assert.Equal(t, newChains[0].Config.TxTimeout, chains[0].Config.TxTimeout)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.SolanaChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[1].ID, chains[0].ID)
	assert.Equal(t, newChains[1].Config.SkipPreflight, chains[0].Config.SkipPreflight)
	assert.Equal(t, newChains[1].Config.TxTimeout, chains[0].Config.TxTimeout)
}

func Test_SolanaChainsController_Update(t *testing.T) {
	t.Parallel()

	hour := models.MustMakeDuration(time.Hour)
	chainUpdate := web.UpdateChainRequest[*db.ChainCfg]{
		Enabled: true,
		Config: &db.ChainCfg{
			SkipPreflight: null.BoolFrom(false),
			TxTimeout:     &hour,
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
					SkipPreflight: null.BoolFrom(false),
					TxTimeout:     &hour,
				}

				chain := db.Chain{
					ID:      validId,
					Enabled: true,
					Cfg:     newChainConfig,
				}
				solanatest.MustInsertChain(t, app.GetSqlxDB(), &chain)

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

			controller := setupSolanaChainsControllerTest(t)

			beforeUpdate := tc.chainBeforeUpdate(t, controller.app)

			body, err := json.Marshal(chainUpdate)
			require.NoError(t, err)

			resp, cleanup := controller.client.Patch(
				fmt.Sprintf("/v2/chains/solana/%s", tc.inputId),
				bytes.NewReader(body),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if beforeUpdate != nil {
				resource1 := presenters.SolanaChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, beforeUpdate.ID)
				assert.Equal(t, resource1.Enabled, chainUpdate.Enabled)
				assert.Equal(t, resource1.Config.SkipPreflight, chainUpdate.Config.SkipPreflight)
				assert.Equal(t, resource1.Config.TxTimeout, chainUpdate.Config.TxTimeout)
			}
		})
	}
}

func Test_SolanaChainsController_Delete(t *testing.T) {
	t.Parallel()

	controller := setupSolanaChainsControllerTest(t)

	hour := models.MustMakeDuration(time.Hour)
	newChainConfig := db.ChainCfg{
		SkipPreflight: null.BoolFrom(false),
		TxTimeout:     &hour,
	}

	chainId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	chain := db.Chain{
		ID:      chainId,
		Enabled: true,
		Cfg:     newChainConfig,
	}
	solanatest.MustInsertChain(t, controller.app.GetSqlxDB(), &chain)

	_, countBefore, err := controller.app.Chains.Solana.Index(0, 10)
	require.NoError(t, err)
	require.Equal(t, 1, countBefore)

	t.Run("non-existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete("/v2/chains/solana/121231")
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		_, countAfter, err := controller.app.Chains.Solana.Index(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, countAfter)
	})

	t.Run("existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete(
			fmt.Sprintf("/v2/chains/solana/%s", chain.ID),
		)
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, countAfter, err := controller.app.Chains.Solana.Index(0, 10)
		require.NoError(t, err)
		require.Equal(t, 0, countAfter)

		_, err = controller.app.Chains.Solana.Show(chain.ID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

type TestSolanaChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupSolanaChainsControllerTest(t *testing.T) *TestSolanaChainsController {
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.SolanaEnabled = null.BoolFrom(true)
	cfg.Overrides.EVMEnabled = null.BoolFrom(false)
	cfg.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()

	return &TestSolanaChainsController{
		app:    app,
		client: client,
	}
}
