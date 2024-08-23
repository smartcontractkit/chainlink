package resolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/web/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type mockEvmConfig struct {
	config.EVM
	linkAddr         string
	gasEstimatorMock *mocks2.GasEstimator
}

func (m *mockEvmConfig) LinkContractAddress() string       { return m.linkAddr }
func (m *mockEvmConfig) GasEstimator() config.GasEstimator { return m.gasEstimatorMock }

func TestResolver_ETHKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetETHKeys {
			ethKeys {
				results {
					address
					isDisabled
					ethBalance
					linkBalance
					maxGasPriceWei
					createdAt
					updatedAt
					chain {
						id
					}
				}
			}
		}`

	address := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")
	secondAddress := common.HexToAddress("0x1438087186fdbfd4c256fa2df446921e30e54df8")
	keys := []ethkey.KeyV2{
		{
			Address:      address,
			EIP55Address: evmtypes.EIP55AddressFromAddress(address),
		},
		{
			Address:      secondAddress,
			EIP55Address: evmtypes.EIP55AddressFromAddress(secondAddress),
		},
	}
	gError := errors.New("error")
	keysError := fmt.Errorf("error getting unlocked keys: %v", gError)
	statesError := fmt.Errorf("error getting key states: %v", gError)

	evmMockConfig := mockEvmConfig{linkAddr: "0x5431F5F973781809D18643b87B44921b11355d81", gasEstimatorMock: mocks2.NewGasEstimator(t)}
	evmMockConfig.gasEstimatorMock.On("PriceMaxKey", mock.Anything).Return(assets.NewWeiI(1))

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ethKeys"),
		{
			name:          "success on prod",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.MustEIP55Address(address.Hex()),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *big.NewI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				cfg := configtest.NewGeneralConfig(t, nil)
				m := map[string]legacyevm.Chain{states[0].EVMChainID.String(): f.Mocks.chain}
				legacyEVMChains := legacyevm.NewLegacyChains(m, cfg.EVMConfigs())

				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.ethClient.On("LINKBalance", mock.Anything, address, linkAddr).Return(commonassets.NewLinkFromJuels(12), nil)
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.balM.On("GetEthBalance", address).Return(assets.NewEth(1))
				f.Mocks.chain.On("BalanceMonitor").Return(f.Mocks.balM)
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.relayerChainInterops.EVMChains = legacyEVMChains
				f.Mocks.relayerChainInterops.Relayers = []loop.Relayer{
					testutils.MockRelayer{
						ChainStatus: types.ChainStatus{
							ID:      "12",
							Enabled: true,
						},
						NodeStatuses: nil,
					},
				}
				f.Mocks.evmORM.PutChains(toml.EVMConfig{ChainID: &chainID})
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)

				f.Mocks.scfg.On("EVM").Return(&evmMockConfig)
			},
			query: query,
			result: `
						{
							"ethKeys": {
								"results": [
									{
										"address": "0x5431F5F973781809D18643b87B44921b11355d81",
										"isDisabled": false,
										"ethBalance": "0.000000000000000001",
										"linkBalance": "12",
										"maxGasPriceWei": "1",
										"createdAt": "2021-01-01T00:00:00Z",
										"updatedAt": "2021-01-01T00:00:00Z",
										"chain": {
											"id": "12"
										}
									}
								]
							}
						}`,
		},

		{
			name:          "success with no chains",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.MustEIP55Address(address.Hex()),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *big.NewI(12)
				f.Mocks.legacyEVMChains.On("Get", states[0].EVMChainID.String()).Return(nil, evmrelay.ErrNoChains)
				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.relayerChainInterops.EVMChains = f.Mocks.legacyEVMChains
				f.Mocks.evmORM.PutChains(toml.EVMConfig{ChainID: &chainID})
				f.Mocks.relayerChainInterops.Relayers = []loop.Relayer{
					testutils.MockRelayer{
						ChainStatus: types.ChainStatus{
							ID:      "12",
							Enabled: true,
						},
						NodeStatuses: nil,
					},
				}
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)
			},
			query: query,
			result: `
							{
								"ethKeys": {
									"results": [
										{
											"address": "0x5431F5F973781809D18643b87B44921b11355d81",
											"isDisabled": false,
											"ethBalance": null,
											"linkBalance": null,
											"maxGasPriceWei": null,
											"createdAt": "2021-01-01T00:00:00Z",
											"updatedAt": "2021-01-01T00:00:00Z",
											"chain": {
												"id": "12"
											}
										}
									]
								}
							}`,
		},

		{
			name:          "generic error on GetAll()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(nil, gError)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: keysError,
					Path:          []interface{}{"ethKeys"},
					Message:       keysError.Error(),
				},
			},
		},
		{
			name:          "generic error on GetStatesForKeys()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(nil, gError)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: statesError,
					Path:          []interface{}{"ethKeys"},
					Message:       statesError.Error(),
				},
			},
		},
		{
			name:          "generic error on Get()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.MustEIP55Address(address.Hex()),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}

				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(ethkey.KeyV2{}, gError)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"ethKeys"},
					Message:       gError.Error(),
				},
			},
		},

		{
			name:          "Empty set on legacy evm chains",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.MustEIP55Address(address.Hex()),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}

				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(ethkey.KeyV2{}, nil)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.Mocks.legacyEVMChains.On("Get", states[0].EVMChainID.String()).Return(f.Mocks.chain, gError)
				f.Mocks.relayerChainInterops.EVMChains = f.Mocks.legacyEVMChains
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query: query,
			result: `
					{
						"ethKeys": {
							"results": []
						}
					}`,
		},
		{
			name:          "generic error on GetLINKBalance()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.MustEIP55Address(address.Hex()),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *big.NewI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.Mocks.ethClient.On("LINKBalance", mock.Anything, address, linkAddr).Return(commonassets.NewLinkFromJuels(12), gError)
				f.Mocks.legacyEVMChains.On("Get", states[0].EVMChainID.String()).Return(f.Mocks.chain, nil)
				f.Mocks.relayerChainInterops.EVMChains = f.Mocks.legacyEVMChains
				f.Mocks.relayerChainInterops.Relayers = []loop.Relayer{
					testutils.MockRelayer{
						ChainStatus: types.ChainStatus{
							ID:      "12",
							Enabled: true,
						},
						NodeStatuses: nil,
					},
				}
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.balM.On("GetEthBalance", address).Return(assets.NewEth(1))
				f.Mocks.chain.On("BalanceMonitor").Return(f.Mocks.balM)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.evmORM.PutChains(toml.EVMConfig{ChainID: &chainID})
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)
				f.Mocks.scfg.On("EVM").Return(&evmMockConfig)
			},
			query: query,
			result: `
						{
							"ethKeys": {
								"results": [
									{
										"address": "0x5431F5F973781809D18643b87B44921b11355d81",
										"isDisabled": false,
										"ethBalance": "0.000000000000000001",
										"linkBalance": null,
										"maxGasPriceWei": "1",
										"createdAt": "2021-01-01T00:00:00Z",
										"updatedAt": "2021-01-01T00:00:00Z",
										"chain": {
											"id": "12"
										}
									}
								]
							}
						}`,
		},
		{
			name:          "success with no eth balance",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    evmtypes.EIP55AddressFromAddress(address),
						EVMChainID: *big.NewI(12),
						Disabled:   false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *big.NewI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.ethKs.On("GetStatesForKeys", mock.Anything, keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", mock.Anything, keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethKs.On("GetAll", mock.Anything).Return(keys, nil)
				f.Mocks.ethClient.On("LINKBalance", mock.Anything, address, linkAddr).Return(commonassets.NewLinkFromJuels(12), nil)
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.chain.On("BalanceMonitor").Return(nil)
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.legacyEVMChains.On("Get", states[0].EVMChainID.String()).Return(f.Mocks.chain, nil)
				f.Mocks.relayerChainInterops.EVMChains = f.Mocks.legacyEVMChains
				f.Mocks.evmORM.PutChains(toml.EVMConfig{ChainID: &chainID})
				f.Mocks.relayerChainInterops.Relayers = []loop.Relayer{
					testutils.MockRelayer{
						ChainStatus: types.ChainStatus{
							ID:      "12",
							Enabled: true,
						},
						NodeStatuses: nil,
					},
				}
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)
				f.Mocks.scfg.On("EVM").Return(&evmMockConfig)
			},
			query: query,
			result: `
						{
							"ethKeys": {
								"results": [
									{
										"address": "0x5431F5F973781809D18643b87B44921b11355d81",
										"isDisabled": false,
										"ethBalance": null,
										"linkBalance": "12",
										"maxGasPriceWei": "1",
										"createdAt": "2021-01-01T00:00:00Z",
										"updatedAt": "2021-01-01T00:00:00Z",
										"chain": {
											"id": "12"
										}
									}
								]
							}
						}`,
		},
	}

	RunGQLTests(t, testCases)
}
