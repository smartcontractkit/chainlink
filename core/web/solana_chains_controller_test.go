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

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_SolanaChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupSolanaChainsControllerTest(t)

	newChainId := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))

	body, err := json.Marshal(web.NewCreateChainRequest(
		newChainId,
		&db.ChainCfg{
			BalancePollPeriod:   utils.MustNewDuration(time.Second),
			ConfirmPollPeriod:   utils.MustNewDuration(time.Minute),
			OCR2CachePollPeriod: utils.MustNewDuration(time.Minute),
			OCR2CacheTTL:        utils.MustNewDuration(time.Second),
			TxTimeout:           utils.MustNewDuration(time.Hour),
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
				Enabled: ptr(true), Chain: config.Chain{
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
		Enabled: ptr(true),
		Chain: config.Chain{
			TxTimeout: utils.MustNewDuration(time.Hour),
		},
	}
	chainB := &solana.SolanaConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestB-%d", rand.Int31n(999999))),
		Enabled: ptr(true),
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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_SolanaChainsController_Update(t *testing.T) {
	t.Parallel()

	chainUpdate := web.UpdateChainRequest[*db.ChainCfg]{
		Enabled: true,
		Config: &db.ChainCfg{
			SkipPreflight: null.BoolFrom(false),
			TxTimeout:     utils.MustNewDuration(time.Hour),
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
					TxTimeout:     utils.MustNewDuration(time.Hour),
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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_SolanaChainsController_Delete(t *testing.T) {
	t.Parallel()

	controller := setupSolanaChainsControllerTest(t)

	newChainConfig := db.ChainCfg{
		SkipPreflight: null.BoolFrom(false),
		TxTimeout:     utils.MustNewDuration(time.Hour),
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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func setupSolanaChainsControllerTest(t *testing.T) *TestSolanaChainsController {
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.SolanaEnabled = null.BoolFrom(true)
	cfg.Overrides.EVMEnabled = null.BoolFrom(false)
	cfg.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestSolanaChainsController{
		app:    app,
		client: client,
	}
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
