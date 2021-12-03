package resolver

import (
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestResolver_ETHKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetETHKeys {
			ethKeys {
				keys {
					address
					isFunding
					createdAt
					updatedAt
					chain {
						id
					}
				}
			}
		}`

	address := ethkey.EIP55Address("0x5431F5F973781809D18643b87B44921b11355d81")
	keys := []ethkey.KeyV2{
		{
			Address: address,
		},
	}
	gError := errors.New("error")
	keysError := fmt.Errorf("error getting unlocked keys: %v", gError)
	statesError := fmt.Errorf("error getting key states: %v", gError)

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
				}
				chainID := *utils.NewBigI(12)

				f.Mocks.cfg.On("Dev").Return(true)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("GetAll").Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.evmORM.On("GetChainsByIDs", []utils.Big{chainID}).Return([]types.Chain{
					{
						ID: chainID,
					},
				}, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"keys": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": true,
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

				f.Mocks.cfg.On("Dev").Return(false)
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.Mocks.ethKs.On("SendingKeys").Return(keys, nil)
				f.Mocks.ethKs.On("GetStatesForKeys", keys).Return(states, nil)
				f.Mocks.ethKs.On("Get", keys[0].Address.Hex()).Return(keys[0], nil)
				f.Mocks.evmORM.On("GetChainsByIDs", []utils.Big{chainID}).Return([]types.Chain{
					{
						ID: chainID,
					},
				}, nil)
				f.Mocks.keystore.On("Eth").Return(f.Mocks.ethKs)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query: query,
			result: `
				{
					"ethKeys": {
						"keys": [
							{
								"address": "0x5431F5F973781809D18643b87B44921b11355d81",
								"isFunding": false,
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
				f.Mocks.ethKs.On("SendingKeys").Return(nil, gError)
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
				f.Mocks.ethKs.On("SendingKeys").Return(keys, nil)
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
				f.Mocks.ethKs.On("SendingKeys").Return(keys, nil)
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
	}

	RunGQLTests(t, testCases)
}
