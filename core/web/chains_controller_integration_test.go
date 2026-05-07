//go:build integration

package web_test

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"

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

	const validID = "Chainlink-12"

	testCases := []struct {
		name           string
		inputID        string
		wantStatusCode int
		want           func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus
	}{
		{
			inputID: validID,
			name:    "success",
			want: func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus {
				return &commonTypes.ChainStatus{
					ID:      validID,
					Enabled: true,
					Config: `ChainID = 'Chainlink-12'
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

[[Nodes]]
Name = 'primary'
TendermintURL = 'http://tender.mint'
`,
				}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputID: "234",
			name:    "not found",
			want: func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupCosmosChainsControllerTestV2(t, chainlink.RawConfig{
				"ChainID":            validID,
				"FallbackGasPrice":   ptr(decimal.RequireFromString("9.999")),
				"GasLimitMultiplier": ptr(decimal.RequireFromString("1.55555")),
				"Nodes": []map[string]any{{
					"Name":          "primary",
					"TendermintURL": "http://tender.mint",
				}},
			})

			wantedResult := tc.want(t, controller.app)
			resp, cleanup := controller.client.Get(
				"/v2/chains/cosmos/" + tc.inputID,
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, wantedResult.ID, resource1.ID)
				assert.Equal(t, wantedResult.Config, resource1.Config)
			}
		})
	}
}

func configForChain(rc chainlink.RawConfig) string {
	gasPrice := cmp.Or(rc["FallbackGasPrice"], "0.015")
	gasLimitMult := cmp.Or(rc["GasLimitMultiplier"], "1.5")

	return fmt.Sprintf(`ChainID = '%s'
Bech32Prefix = 'wasm'
BlockRate = '6s'
BlocksUntilTxTimeout = 30
ConfirmPollPeriod = '1s'
FallbackGasPrice = '%s'
GasToken = 'ucosm'
GasLimitMultiplier = '%s'
MaxMsgsPerBatch = 100
OCR2CachePollPeriod = '4s'
OCR2CacheTTL = '1m0s'
TxMsgTimeout = '10m0s'

[[Nodes]]
Name = 'primary'
TendermintURL = 'http://tender.mint'
`, rc.ChainID(), gasPrice, gasLimitMult)
}

func Test_CosmosChainsController_Index(t *testing.T) {
	t.Parallel()

	chainA := chainlink.RawConfig{
		"ChainID":          "a" + cosmostest.RandomChainID(),
		"FallbackGasPrice": ptr(decimal.RequireFromString("9.999")),
		"Nodes": []map[string]any{{
			"Name":          "primary",
			"TendermintURL": "http://tender.mint",
		}},
	}

	chainB := chainlink.RawConfig{
		"ChainID":            "b" + cosmostest.RandomChainID(),
		"GasLimitMultiplier": ptr(decimal.RequireFromString("1.55555")),
		"Nodes": []map[string]any{{
			"Name":          "primary",
			"TendermintURL": "http://tender.mint",
		}},
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

	var chains []presenters.ChainResource
	err = web.ParsePaginatedResponse(body, &chains, &links)
	require.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, chainA.ChainID(), chains[0].ID)
	assert.Equal(t, configForChain(chainA), chains[0].Config)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	require.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, chainB.ChainID(), chains[0].ID)
	assert.Equal(t, configForChain(chainB), chains[0].Config)
}

type TestCosmosChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupCosmosChainsControllerTestV2(t *testing.T, cfgs ...chainlink.RawConfig) *TestCosmosChainsController {
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

func Test_SolanaChainsController_Show(t *testing.T) {
	t.Parallel()

	const validID = "Chainlink-12"

	testCases := []struct {
		name           string
		inputID        string
		wantStatusCode int
		want           func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus
	}{
		{
			inputID: validID,
			name:    "success",
			want: func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus {
				return &commonTypes.ChainStatus{
					ID:      validID,
					Enabled: true,
					Config: `ChainID = 'Chainlink-12'
Enabled = true
BlockTime = '500ms'
BalancePollPeriod = '5s'
ConfirmPollPeriod = '500ms'
OCR2CachePollPeriod = '1s'
OCR2CacheTTL = '1m0s'
TxTimeout = '1h0m0s'
TxRetryTimeout = '10s'
TxConfirmTimeout = '30s'
TxExpirationRebroadcast = false
TxRetentionTimeout = '0s'
SkipPreflight = false
Commitment = 'confirmed'
MaxRetries = 0
FeeEstimatorMode = 'fixed'
ComputeUnitPriceMax = 1000
ComputeUnitPriceMin = 0
ComputeUnitPriceDefault = 0
FeeBumpPeriod = '3s'
BlockHistoryPollPeriod = '5s'
BlockHistorySize = 1
BlockHistoryBatchLoadSize = 20
ComputeUnitLimitDefault = 200000
EstimateComputeUnitLimit = false
LogPollerStartingLookback = '24h0m0s'
LogPollerCPIEventsEnabled = true
LogPollerSlotsBatchSize = 1000

[Workflow]
AcceptanceTimeout = '45s'
ForwarderAddress = '11111111111111111111111111111111'
ForwarderState = '11111111111111111111111111111111'
FromAddress = '11111111111111111111111111111111'
GasLimitDefault = 300000
Local = false
PollPeriod = '3s'
TxAcceptanceState = 3

[MultiNode]
Enabled = false
PollFailureThreshold = 5
PollInterval = '15s'
SelectionMode = 'PriorityLevel'
SyncThreshold = 10
NodeIsSyncingEnabled = false
LeaseDuration = '1m0s'
NewHeadsPollInterval = '5s'
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '20s'
VerifyChainID = true
NodeNoNewHeadsThreshold = '20s'
NoNewFinalizedHeadsThreshold = '20s'
FinalityDepth = 0
FinalityTagEnabled = true
FinalizedBlockOffset = 50

[[Nodes]]
Name = 'primary'
URL = 'http://solana.example'
SendOnly = false
`,
				}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputID: "234",
			name:    "not found",
			want: func(t *testing.T, app *cltest.TestApplication) *commonTypes.ChainStatus {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupSolanaChainsControllerTestV2(t, chainlink.RawConfig{
				"ChainID":       validID,
				"SkipPreflight": false,
				"TxTimeout":     "1h0m0s",
				"Nodes": []map[string]any{{
					"Name": "primary",
					"URL":  "http://solana.example",
				}},
			})

			wantedResult := tc.want(t, controller.app)
			resp, cleanup := controller.client.Get(
				"/v2/chains/solana/" + tc.inputID,
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, wantedResult.ID, resource1.ID)
				assert.Equal(t, wantedResult.Enabled, resource1.Enabled)
				assert.Equal(t, wantedResult.Config, resource1.Config)
			}
		})
	}
}

func Test_SolanaChainsController_Index(t *testing.T) {
	t.Parallel()

	chainA := chainlink.RawConfig{
		"ChainID":   fmt.Sprintf("ChainlinktestA-%d", rand.Int32N(999999)),
		"TxTimeout": "1h0m0s",
		"Nodes": []map[string]any{{
			"Name": "primary",
			"URL":  "http://solana.example",
		}},
	}
	chainB := chainlink.RawConfig{
		"ChainID":       fmt.Sprintf("ChainlinktestB-%d", rand.Int32N(999999)),
		"SkipPreflight": false,
		"Nodes": []map[string]any{{
			"Name": "primary",
			"URL":  "http://solana.example",
		}},
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

	chains := []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	require.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, chainA.ChainID(), chains[0].ID)
	assert.NotEmpty(t, chains[0].Config)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	require.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, chainB.ChainID(), chains[0].ID)
	assert.NotEmpty(t, chains[0].Config)
}

type TestSolanaChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupSolanaChainsControllerTestV2(t *testing.T, cfgs ...chainlink.RawConfig) *TestSolanaChainsController {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = cfgs
		c.EVM = nil
	})
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	return &TestSolanaChainsController{
		app:    app,
		client: client,
	}
}
