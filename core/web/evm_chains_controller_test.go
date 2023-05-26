package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func Test_EVMChainsController_Show(t *testing.T) {
	t.Parallel()

	validId := utils.NewBig(testutils.NewRandomEVMChainID())

	testCases := []struct {
		name           string
		inputId        string
		wantStatusCode int
		want           *evmcfg.EVMConfig
	}{
		{
			inputId: validId.String(),
			name:    "success",
			want: &evmcfg.EVMConfig{
				ChainID: validId,
				Enabled: ptr(true),
				Chain: evmcfg.Defaults(nil, &evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						EIP1559DynamicFees: ptr(true),
						BlockHistory: evmcfg.BlockHistoryEstimator{
							BlockHistorySize: ptr[uint16](50),
						},
					},
					RPCBlockQueryDelay:       ptr[uint16](23),
					MinIncomingConfirmations: ptr[uint32](12),
					LinkContractAddress:      ptr(ethkey.EIP55AddressFromAddress(testutils.NewAddress())),
				}),
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId:        "invalidid",
			name:           "invalid id",
			want:           nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			inputId:        "234",
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
					c.EVM = evmcfg.EVMConfigs{tc.want}
				}
			}))

			wantedResult := tc.want
			resp, cleanup := controller.client.Get(
				fmt.Sprintf("/v2/chains/evm/%s", tc.inputId),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if wantedResult != nil {
				resource1 := presenters.EVMChainResource{}
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

	newChains := evmcfg.EVMConfigs{
		{ChainID: utils.NewBig(testutils.NewRandomEVMChainID()), Chain: evmcfg.Defaults(nil)},
		{
			ChainID: utils.NewBig(testutils.NewRandomEVMChainID()),
			Chain: evmcfg.Defaults(nil, &evmcfg.Chain{
				RPCBlockQueryDelay: ptr[uint16](13),
				GasEstimator: evmcfg.GasEstimator{
					EIP1559DynamicFees: ptr(true),
					BlockHistory: evmcfg.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](1),
					},
				},
				MinIncomingConfirmations: ptr[uint32](120),
			}),
		},
		{
			ChainID: utils.NewBig(testutils.NewRandomEVMChainID()),
			Chain: evmcfg.Defaults(nil, &evmcfg.Chain{
				RPCBlockQueryDelay: ptr[uint16](5),
				GasEstimator: evmcfg.GasEstimator{
					EIP1559DynamicFees: ptr(false),
					BlockHistory: evmcfg.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](2),
					},
				},
				MinIncomingConfirmations: ptr[uint32](30),
			}),
		},
	}

	controller := setupEVMChainsControllerTest(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = append(c.EVM, newChains...)
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
	require.Equal(t, 1+len(newChains), metaCount)

	var links jsonapi.Links

	var chains []presenters.EVMChainResource
	err = web.ParsePaginatedResponse(body, &chains, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[1].ChainID.String(), chains[2].ID)
	toml, err := newChains[1].TOMLString()
	require.NoError(t, err)
	assert.Equal(t, toml, chains[2].Config)

	resp, cleanup = controller.client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	chains = []presenters.EVMChainResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &chains, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"].Href)
	assert.NotEmpty(t, links["prev"].Href)

	assert.Len(t, links, 1)
	assert.Equal(t, newChains[2].ChainID.String(), chains[0].ID)
	toml, err = newChains[2].TOMLString()
	require.NoError(t, err)
	assert.Equal(t, toml, chains[0].Config)
}

type TestEVMChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMChainsControllerTest(t *testing.T, cfg chainlink.GeneralConfig) *TestEVMChainsController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplicationWithConfig(t, cfg)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestEVMChainsController{
		app:    app,
		client: client,
	}
}

func ptr[T any](t T) *T { return &t }
