package resolver

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestResolver_ETHKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetETHKeys {
			ethKeys {
				results {
					address
					isFunding
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

	address := ethkey.EIP55Address("0x5431F5F973781809D18643b87B44921b11355d81")
	secondAddress := ethkey.EIP55Address("0x1438087186fdbfd4c256fa2df446921e30e54df8")
	keys := []ethkey.KeyV2{
		{
			Address: address,
		},
		{
			Address: secondAddress,
		},
	}
	gError := errors.New("error")
	keysError := fmt.Errorf("error getting unlocked keys: %v", gError)
	statesError := fmt.Errorf("error getting key states: %v", gError)
	chainError := fmt.Errorf("error getting EVM Chain: %v", gError)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ethKeys"),
		{
			name:          "success on dev",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  true,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
					{
						Address:    secondAddress,
						EVMChainID: *utils.NewBigI(42),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *utils.NewBigI(12)
				chainID2 := *utils.NewBigI(42)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.cfg.On("Dev").Return(true)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("GetAll").Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethKs.On("Get", keys[1].Address.Hex()).Return(keys[1], nil)
				f.Mocks.ethClient.On("GetLINKBalance", linkAddr, address.Address()).Return(assets.NewLinkFromJuels(12), nil)
				f.Mocks.ethClient.On("GetLINKBalance", linkAddr, secondAddress.Address()).Return(assets.NewLinkFromJuels(100), nil)
				f.Mocks.balM.On("GetEthBalance", address.Address()).Return(assets.NewEth(1))
				f.Mocks.balM.On("GetEthBalance", secondAddress.Address()).Return(assets.NewEth(12))
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.chain.On("BalanceMonitor").Return(f.Mocks.balM)
				f.Mocks.scfg.On("LinkContractAddress").Return("0x5431F5F973781809D18643b87B44921b11355d81")
				f.Mocks.scfg.On("KeySpecificMaxGasPriceWei", keys[0].Address.Address()).Return(big.NewInt(1))
				f.Mocks.scfg.On("KeySpecificMaxGasPriceWei", keys[1].Address.Address()).Return(big.NewInt(1))
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, nil)
				f.Mocks.chainSet.On("Get", states[1].EVMChainID.ToInt()).Return(f.Mocks.chain, nil)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID}, types.DBChain{ID: chainID2})
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"results": [
							{
								"address": "0x1438087186fdbfd4c256fa2df446921e30e54df8",
								"isFunding": false,
								"ethBalance": "0.000000000000000012",
								"linkBalance": "100",
								"maxGasPriceWei": "1",
								"createdAt": "2021-01-01T00:00:00Z",
								"updatedAt": "2021-01-01T00:00:00Z",
								"chain": {
									"id": "42"
								}
							},
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": true,
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
			name:          "success on prod",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *utils.NewBigI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethClient.On("GetLINKBalance", linkAddr, address.Address()).Return(assets.NewLinkFromJuels(12), nil)
				f.Mocks.scfg.On("LinkContractAddress").Return("0x5431F5F973781809D18643b87B44921b11355d81")
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.balM.On("GetEthBalance", address.Address()).Return(assets.NewEth(1))
				f.Mocks.chain.On("BalanceMonitor").Return(f.Mocks.balM)
				f.Mocks.scfg.On("KeySpecificMaxGasPriceWei", keys[0].Address.Address()).Return(big.NewInt(1))
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, nil)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"results": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": false,
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
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *utils.NewBigI(12)

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, evm.ErrNoChains)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"results": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": false,
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
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("Dev").Return(true)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("GetAll").Return(nil, gError)
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
			name:          "generic error on SendingKeys()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(nil, gError)
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
			before: func(f *gqlTestFramework) {
				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(nil, gError)
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
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(ethkey.KeyV2{}, gError)
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
			name:          "generic error on #chainSet.Get()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(ethkey.KeyV2{}, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, gError)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: chainError,
					Path:          []interface{}{"ethKeys"},
					Message:       chainError.Error(),
				},
			},
		},
		{
			name:          "generic error on GetLINKBalance()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *utils.NewBigI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.Mocks.ethClient.On("GetLINKBalance", linkAddr, address.Address()).Return(assets.NewLinkFromJuels(12), gError)
				f.Mocks.scfg.On("LinkContractAddress").Return("0x5431F5F973781809D18643b87B44921b11355d81")
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, nil)
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.balM.On("GetEthBalance", address.Address()).Return(assets.NewEth(1))
				f.Mocks.chain.On("BalanceMonitor").Return(f.Mocks.balM)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
				f.Mocks.scfg.On("KeySpecificMaxGasPriceWei", keys[0].Address.Address()).Return(big.NewInt(1))
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"results": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": false,
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
			before: func(f *gqlTestFramework) {
				states := []ethkey.State{
					{
						Address:    address,
						EVMChainID: *utils.NewBigI(12),
						IsFunding:  false,
						CreatedAt:  f.Timestamp(),
						UpdatedAt:  f.Timestamp(),
					},
				}
				chainID := *utils.NewBigI(12)
				linkAddr := common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys", (*big.Int)(nil)).Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.ethClient.On("GetLINKBalance", linkAddr, address.Address()).Return(assets.NewLinkFromJuels(12), nil)
				f.Mocks.scfg.On("LinkContractAddress").Return("0x5431F5F973781809D18643b87B44921b11355d81")
				f.Mocks.chain.On("Client").Return(f.Mocks.ethClient)
				f.Mocks.chain.On("BalanceMonitor").Return(nil)
				f.Mocks.scfg.On("KeySpecificMaxGasPriceWei", keys[0].Address.Address()).Return(big.NewInt(1))
				f.Mocks.chain.On("Config").Return(f.Mocks.scfg)
				f.Mocks.chainSet.On("Get", states[0].EVMChainID.ToInt()).Return(f.Mocks.chain, nil)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"results": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": false,
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
