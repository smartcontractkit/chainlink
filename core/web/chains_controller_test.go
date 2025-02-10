package web_test

import (
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func Test_EVMChainsController_Show(t *testing.T) {
	t.Parallel()

	validID := ubig.New(testutils.NewRandomEVMChainID())

	testCases := []struct {
		name           string
		inputID        string
		wantStatusCode int
		want           *toml.EVMConfig
	}{
		{
			inputID: validID.String(),
			name:    "success",
			want: &toml.EVMConfig{
				ChainID: validID,
				Enabled: ptr(true),
				Chain: toml.Defaults(nil, &toml.Chain{
					GasEstimator: toml.GasEstimator{
						EIP1559DynamicFees: ptr(true),
						BlockHistory: toml.BlockHistoryEstimator{
							BlockHistorySize: ptr[uint16](50),
						},
					},
					RPCBlockQueryDelay:       ptr[uint16](23),
					MinIncomingConfirmations: ptr[uint32](12),
					LinkContractAddress:      ptr(types.EIP55AddressFromAddress(testutils.NewAddress())),
				}),
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputID:        "invalidid",
			name:           "invalid id",
			want:           nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			inputID:        "234",
			name:           "not found",
			want:           nil,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupEVMChainsControllerTest(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				if tc.want != nil {
					c.EVM = toml.EVMConfigs{tc.want}
				}
			}))

			wantedResult := tc.want
			resp, cleanup := controller.client.Get(
				"/v2/chains/evm/" + tc.inputID,
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.ChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, wantedResult.ChainID.String())
				toml, err := wantedResult.TOMLString()
				require.NoError(t, err)
				assert.Equal(t, toml, resource1.Config)
			}
		})
	}
}

func Test_EVMChainsController_Index(t *testing.T) {
	t.Parallel()

	// sort test chain ids to make expected comparison easy
	chainIDs := []*big.Int{testutils.NewRandomEVMChainID(), testutils.NewRandomEVMChainID(), testutils.NewRandomEVMChainID()}
	sort.Slice(chainIDs, func(i, j int) bool {
		return chainIDs[i].String() < chainIDs[j].String()
	})

	configuredChains := toml.EVMConfigs{
		{ChainID: ubig.New(chainIDs[0]), Chain: toml.Defaults(nil)},
		{
			ChainID: ubig.New(chainIDs[1]),
			Chain: toml.Defaults(nil, &toml.Chain{
				RPCBlockQueryDelay: ptr[uint16](13),
				GasEstimator: toml.GasEstimator{
					EIP1559DynamicFees: ptr(true),
					BlockHistory: toml.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](1),
					},
				},
				MinIncomingConfirmations: ptr[uint32](120),
			}),
		},
		{
			ChainID: ubig.New(chainIDs[2]),
			Chain: toml.Defaults(nil, &toml.Chain{
				RPCBlockQueryDelay: ptr[uint16](5),
				GasEstimator: toml.GasEstimator{
					EIP1559DynamicFees: ptr(false),
					BlockHistory: toml.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](2),
					},
				},
				MinIncomingConfirmations: ptr[uint32](30),
			}),
		},
	}

	assert.Len(t, configuredChains, 3)
	controller := setupEVMChainsControllerTest(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = append(c.EVM, configuredChains...)
	}))

	badResp, cleanup := controller.client.Get("/v2/chains/evm?size=asd")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusUnprocessableEntity, badResp.StatusCode)

	resp, cleanup := controller.client.Get("/v2/chains/evm?size=3")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 1+len(configuredChains), metaCount)

	var links jsonapi.Links

	var gotChains []presenters.ChainResource
	err = web.ParsePaginatedResponse(body, &gotChains, &links)
	require.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	// the difference in index value here seems to be due to the fact
	// that cltest always has a default EVM chain, which is the off-by-one
	// in the indices
	assert.Equal(t, gotChains[2].ID, configuredChains[1].ChainID.String())
	toml, err := configuredChains[1].TOMLString()
	require.NoError(t, err)
	assert.Equal(t, toml, gotChains[2].Config)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	gotChains = []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &gotChains, &links)
	require.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, gotChains[0].ID, configuredChains[2].ChainID.String())
	toml, err = configuredChains[2].TOMLString()
	require.NoError(t, err)
	assert.Equal(t, toml, gotChains[0].Config)
}

type TestEVMChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMChainsControllerTest(t *testing.T, cfg chainlink.GeneralConfig) *TestEVMChainsController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplicationWithConfig(t, cfg)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)

	return &TestEVMChainsController{
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
ComputeUnitLimitDefault = 200000
EstimateComputeUnitLimit = false
Nodes = []

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
NodeNoNewHeadsThreshold = '20s'
NoNewFinalizedHeadsThreshold = '20s'
FinalityDepth = 0
FinalityTagEnabled = true
FinalizedBlockOffset = 50
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

			controller := setupSolanaChainsControllerTestV2(t, &config.TOMLConfig{
				ChainID: ptr(validID),
				Chain: config.Chain{
					SkipPreflight: ptr(false),
					TxTimeout:     commoncfg.MustNewDuration(time.Hour),
				},
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

	chainA := &config.TOMLConfig{
		ChainID: ptr(fmt.Sprintf("ChainlinktestA-%d", rand.Int31n(999999))),
		Chain: config.Chain{
			TxTimeout: commoncfg.MustNewDuration(time.Hour),
		},
	}
	chainB := &config.TOMLConfig{
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

	chains := []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(body, &chains, &links)
	require.NoError(t, err)
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

	chains = []presenters.ChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	require.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, *chainB.ChainID, chains[0].ID)
	tomlB, err := chainB.TOMLString()
	require.NoError(t, err)
	assert.Equal(t, tomlB, chains[0].Config)
}

type TestSolanaChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupSolanaChainsControllerTestV2(t *testing.T, cfgs ...*config.TOMLConfig) *TestSolanaChainsController {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
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

func ptr[T any](t T) *T { return &t }
