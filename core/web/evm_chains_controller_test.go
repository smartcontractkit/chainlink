package web_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	corecfg "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func Test_EVMChainsController_Create(t *testing.T) {
	t.Parallel()

	controller := setupEVMChainsControllerTestLegacy(t)

	newChainId := utils.NewBig(testutils.NewRandomEVMChainID())

	body, err := json.Marshal(web.NewCreateChainRequest(newChainId,
		&types.ChainCfg{
			BlockHistoryEstimatorBlockDelay:       null.IntFrom(1),
			BlockHistoryEstimatorBlockHistorySize: null.IntFrom(12),
			EvmEIP1559DynamicFees:                 null.BoolFrom(false),
			MinIncomingConfirmations:              null.IntFrom(10),
		}))
	require.NoError(t, err)

	resp, cleanup := controller.client.Post("/v2/chains/evm", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	chainSet := controller.app.GetChains().EVM
	dbChain, err := chainSet.ORM().Chain(*newChainId)
	require.NoError(t, err)

	resource := presenters.EVMChainResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, dbChain.ID.String())
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockDelay, dbChain.Cfg.BlockHistoryEstimatorBlockDelay)
	assert.Equal(t, resource.Config.BlockHistoryEstimatorBlockHistorySize, dbChain.Cfg.BlockHistoryEstimatorBlockHistorySize)
	assert.Equal(t, resource.Config.EvmEIP1559DynamicFees, dbChain.Cfg.EvmEIP1559DynamicFees)
	assert.Equal(t, resource.Config.MinIncomingConfirmations, dbChain.Cfg.MinIncomingConfirmations)
}

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
			wantStatusCode: http.StatusUnprocessableEntity,
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
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockDelay.Int64, int64(*wantedResult.Chain.RPCBlockQueryDelay))
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockHistorySize.Int64, int64(*wantedResult.Chain.GasEstimator.BlockHistory.BlockHistorySize))
				assert.Equal(t, resource1.Config.EvmEIP1559DynamicFees.Bool, *wantedResult.Chain.GasEstimator.EIP1559DynamicFees)
				assert.Equal(t, resource1.Config.MinIncomingConfirmations.Int64, int64(*wantedResult.Chain.MinIncomingConfirmations))
				assert.Equal(t, resource1.Config.LinkContractAddress.String, wantedResult.Chain.LinkContractAddress.String())
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
	assert.Equal(t, int64(*newChains[1].Chain.RPCBlockQueryDelay), chains[2].Config.BlockHistoryEstimatorBlockDelay.Int64)
	assert.Equal(t, int64(*newChains[1].Chain.GasEstimator.BlockHistory.BlockHistorySize), chains[2].Config.BlockHistoryEstimatorBlockHistorySize.Int64)
	assert.Equal(t, *newChains[1].Chain.GasEstimator.EIP1559DynamicFees, chains[2].Config.EvmEIP1559DynamicFees.Bool)
	assert.Equal(t, int64(*newChains[1].Chain.MinIncomingConfirmations), chains[2].Config.MinIncomingConfirmations.Int64)

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
	assert.Equal(t, int64(*newChains[2].Chain.RPCBlockQueryDelay), chains[0].Config.BlockHistoryEstimatorBlockDelay.Int64)
	assert.Equal(t, int64(*newChains[2].Chain.GasEstimator.BlockHistory.BlockHistorySize), chains[0].Config.BlockHistoryEstimatorBlockHistorySize.Int64)
	assert.Equal(t, *newChains[2].Chain.GasEstimator.EIP1559DynamicFees, chains[0].Config.EvmEIP1559DynamicFees.Bool)
	assert.Equal(t, int64(*newChains[2].Chain.MinIncomingConfirmations), chains[0].Config.MinIncomingConfirmations.Int64)
}

func Test_EVMChainsController_Update(t *testing.T) {
	t.Parallel()

	chainUpdate := web.UpdateChainRequest[*types.ChainCfg]{
		Enabled: true,
		Config: &types.ChainCfg{
			BlockHistoryEstimatorBlockDelay:       null.IntFrom(55),
			BlockHistoryEstimatorBlockHistorySize: null.IntFrom(33),
			EvmEIP1559DynamicFees:                 null.BoolFrom(true),
			MinIncomingConfirmations:              null.IntFrom(100),
			LinkContractAddress:                   null.StringFrom(utils.ZeroAddress.String()),
		},
	}

	validId := utils.NewBig(testutils.NewRandomEVMChainID())

	testCases := []struct {
		name              string
		inputId           string
		wantStatusCode    int
		chainBeforeUpdate func(t *testing.T, app *cltest.TestApplication) *types.DBChain
	}{
		{
			inputId: validId.String(),
			name:    "success",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *types.DBChain {
				newChainConfig := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay:       null.IntFrom(5),
					BlockHistoryEstimatorBlockHistorySize: null.IntFrom(2),
					EvmEIP1559DynamicFees:                 null.BoolFrom(false),
					MinIncomingConfirmations:              null.IntFrom(30),
				}

				chain := types.DBChain{
					ID:      *validId,
					Enabled: true,
					Cfg:     &newChainConfig,
				}
				evmtest.MustInsertChain(t, app.GetSqlxDB(), &chain)

				return &chain
			},
			wantStatusCode: http.StatusOK,
		},
		{
			inputId: "invalidid",
			name:    "invalid id",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *types.DBChain {
				return nil
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			inputId: "341212",
			name:    "not found",
			chainBeforeUpdate: func(t *testing.T, app *cltest.TestApplication) *types.DBChain {
				return nil
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			controller := setupEVMChainsControllerTestLegacy(t)

			beforeUpdate := tc.chainBeforeUpdate(t, controller.app)

			body, err := json.Marshal(chainUpdate)
			require.NoError(t, err)

			resp, cleanup := controller.client.Patch(
				fmt.Sprintf("/v2/chains/evm/%s", tc.inputId),
				bytes.NewReader(body),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if beforeUpdate != nil {
				resource1 := presenters.EVMChainResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource1)
				require.NoError(t, err)

				assert.Equal(t, resource1.ID, beforeUpdate.ID.String())
				assert.Equal(t, resource1.Enabled, chainUpdate.Enabled)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockDelay, chainUpdate.Config.BlockHistoryEstimatorBlockDelay)
				assert.Equal(t, resource1.Config.BlockHistoryEstimatorBlockHistorySize, chainUpdate.Config.BlockHistoryEstimatorBlockHistorySize)
				assert.Equal(t, resource1.Config.EvmEIP1559DynamicFees, chainUpdate.Config.EvmEIP1559DynamicFees)
				assert.Equal(t, resource1.Config.MinIncomingConfirmations, chainUpdate.Config.MinIncomingConfirmations)
				assert.Equal(t, resource1.Config.LinkContractAddress, chainUpdate.Config.LinkContractAddress)
			}
		})
	}
}

func Test_EVMChainsController_Delete(t *testing.T) {
	t.Parallel()

	controller := setupEVMChainsControllerTestLegacy(t)

	newChainConfig := types.ChainCfg{
		BlockHistoryEstimatorBlockDelay:       null.IntFrom(5),
		BlockHistoryEstimatorBlockHistorySize: null.IntFrom(2),
		EvmEIP1559DynamicFees:                 null.BoolFrom(false),
		MinIncomingConfirmations:              null.IntFrom(30),
	}

	chainId := *utils.NewBig(testutils.NewRandomEVMChainID())
	chain := types.DBChain{
		ID:      chainId,
		Enabled: true,
		Cfg:     &newChainConfig,
	}
	evmtest.MustInsertChain(t, controller.app.GetSqlxDB(), &chain)

	_, countBefore, err := controller.app.EVMORM().Chains(0, 10)
	require.NoError(t, err)
	// 3 with the default chains
	require.Equal(t, 3, countBefore)

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		resp, cleanup := controller.client.Delete("/v2/chains/evm/invalid_id")
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("non-existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete("/v2/chains/evm/121231")
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		_, countAfter, err := controller.app.EVMORM().Chains(0, 10)
		require.NoError(t, err)
		// 3 with the default chains
		require.Equal(t, 3, countAfter)
	})

	t.Run("existing chain", func(t *testing.T) {
		resp, cleanup := controller.client.Delete(
			fmt.Sprintf("/v2/chains/evm/%d", chain.ID.ToInt()),
		)
		t.Cleanup(cleanup)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, countAfter, err := controller.app.EVMORM().Chains(0, 10)
		require.NoError(t, err)
		// 3 with the default chains
		require.Equal(t, 2, countAfter)

		_, err = controller.app.EVMORM().Chain(chain.ID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

type TestEVMChainsController struct {
	app    *cltest.TestApplication
	client cltest.HTTPClientCleaner
}

func setupEVMChainsControllerTest(t *testing.T, cfg corecfg.GeneralConfig) *TestEVMChainsController {
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

// setupEVMChainsControllerTestLegacy exists to support legacy-only tests until it is removed.
func setupEVMChainsControllerTestLegacy(t *testing.T) *TestEVMChainsController {
	// Using this instead of `NewApplicationEVMDisabled` since we need the chain set to be loaded in the app
	// for the sake of the API endpoints to work properly
	app := cltest.NewApplicationWithConfig(t, cltest.NewTestGeneralConfig(t))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return &TestEVMChainsController{
		app:    app,
		client: client,
	}
}

func ptr[T any](t T) *T { return &t }
