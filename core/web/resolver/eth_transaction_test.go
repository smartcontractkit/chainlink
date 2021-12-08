package resolver

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestResolver_EthTransaction(t *testing.T) {
	t.Parallel()

	query := `
		query GetEthTransaction($hash: ID!) {
			ethTransaction(hash: $hash) {
				... on EthTransaction {
					from
					to
					state
					data
					gasLimit
					gasPrice
					value
					evmChainID
					chain {
						id
					}
					nonce
					hash
					hex
					sentAt
				}
				... on NotFoundError {
					code
					message
				}
			}
		}`
	variables := map[string]interface{}{
		"hash": "0x5431F5F973781809D18643b87B44921b11355d81",
	}
	hash := common.HexToHash("0x5431F5F973781809D18643b87B44921b11355d81")
	chainID := *utils.NewBigI(22)
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query, variables: variables}, "ethTransaction"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.bptxmORM.On("FindEthTxAttempt", hash).Return(&bulletprooftxmanager.EthTxAttempt{
					Hash:                    hash,
					GasPrice:                utils.NewBigI(12),
					SignedRawTx:             []byte("something"),
					BroadcastBeforeBlockNum: nil,
					EthTx: bulletprooftxmanager.EthTx{
						ToAddress:      common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						FromAddress:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						State:          bulletprooftxmanager.EthTxInProgress,
						EncodedPayload: []byte("encoded payload"),
						GasLimit:       100,
						Value:          assets.NewEthValue(100),
						EVMChainID:     *utils.NewBigI(22),
						Nonce:          nil,
					},
				}, nil)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
				f.Mocks.evmORM.On("GetChainsByIDs", []utils.Big{chainID}).Return([]types.Chain{
					{
						ID: chainID,
					},
				}, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     query,
			variables: variables,
			result: `
				{
					"ethTransaction": {
						"from": "0x5431F5F973781809D18643b87B44921b11355d81",
						"to": "0x5431F5F973781809D18643b87B44921b11355d81",
						"chain": {
							"id": "22"
						},
						"data": "0x656e636f646564207061796c6f6164",
						"state": "in_progress",
						"gasLimit": "100",
						"gasPrice": "12",
						"value": "0.000000000000000100",
						"nonce": null,
						"hash": "0x0000000000000000000000005431f5f973781809d18643b87b44921b11355d81",
						"hex": "0x736f6d657468696e67",
						"sentAt": null,
						"evmChainID": "22"
					}
				}`,
		},
		{
			name:          "success without nil values",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				num := int64(2)

				f.Mocks.bptxmORM.On("FindEthTxAttempt", hash).Return(&bulletprooftxmanager.EthTxAttempt{
					Hash:                    hash,
					GasPrice:                utils.NewBigI(12),
					SignedRawTx:             []byte("something"),
					BroadcastBeforeBlockNum: &num,
					EthTx: bulletprooftxmanager.EthTx{
						ToAddress:      common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						FromAddress:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						State:          bulletprooftxmanager.EthTxInProgress,
						EncodedPayload: []byte("encoded payload"),
						GasLimit:       100,
						Value:          assets.NewEthValue(100),
						EVMChainID:     *utils.NewBigI(22),
						Nonce:          &num,
					},
				}, nil)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
				f.Mocks.evmORM.On("GetChainsByIDs", []utils.Big{chainID}).Return([]types.Chain{
					{
						ID: chainID,
					},
				}, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     query,
			variables: variables,
			result: `
				{
					"ethTransaction": {
						"from": "0x5431F5F973781809D18643b87B44921b11355d81",
						"to": "0x5431F5F973781809D18643b87B44921b11355d81",
						"chain": {
							"id": "22"
						},
						"data": "0x656e636f646564207061796c6f6164",
						"state": "in_progress",
						"gasLimit": "100",
						"gasPrice": "12",
						"value": "0.000000000000000100",
						"nonce": "2",
						"hash": "0x0000000000000000000000005431f5f973781809d18643b87b44921b11355d81",
						"hex": "0x736f6d657468696e67",
						"sentAt": "2",
						"evmChainID": "22"
					}
				}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.bptxmORM.On("FindEthTxAttempt", hash).Return(nil, sql.ErrNoRows)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
			},
			query:     query,
			variables: variables,
			result: `
				{
					"ethTransaction": {
						"code": "NOT_FOUND",
						"message": "transaction not found"
					}
				}`,
		},
		{
			name:          "generic error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.bptxmORM.On("FindEthTxAttempt", hash).Return(nil, gError)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
			},
			query:     query,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"ethTransaction"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_EthTransactions(t *testing.T) {
	t.Parallel()

	query := `
		query GetEthTransactions {
			ethTransactions {
				results {
					from
					to
					state
					data
					gasLimit
					gasPrice
					value
					evmChainID
					nonce
					hash
					hex
					sentAt
				}
				metadata {
					total
				}
			}
		}`
	hash := common.HexToHash("0x5431F5F973781809D18643b87B44921b11355d81")
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ethTransactions"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				num := int64(2)

				f.Mocks.bptxmORM.On("EthTransactionsWithAttempts", PageDefaultOffset, PageDefaultLimit).Return([]bulletprooftxmanager.EthTx{
					{
						ToAddress:      common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						FromAddress:    common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81"),
						State:          bulletprooftxmanager.EthTxInProgress,
						EncodedPayload: []byte("encoded payload"),
						GasLimit:       100,
						Value:          assets.NewEthValue(100),
						EVMChainID:     *utils.NewBigI(22),
						EthTxAttempts: []bulletprooftxmanager.EthTxAttempt{
							{
								Hash:                    hash,
								GasPrice:                utils.NewBigI(12),
								SignedRawTx:             []byte("something"),
								BroadcastBeforeBlockNum: &num,
							},
						},
					},
				}, 1, nil)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
			},
			query: query,
			result: `
				{
					"ethTransactions": {
						"results": [{
							"from": "0x5431F5F973781809D18643b87B44921b11355d81",
							"to": "0x5431F5F973781809D18643b87B44921b11355d81",
							"data": "0x656e636f646564207061796c6f6164",
							"state": "in_progress",
							"gasLimit": "100",
							"gasPrice": "12",
							"value": "0.000000000000000100",
							"nonce": null,
							"hash": "0x0000000000000000000000005431f5f973781809d18643b87b44921b11355d81",
							"hex": "0x736f6d657468696e67",
							"sentAt": "2",
							"evmChainID": "22"
						}],
						"metadata": {
							"total": 1
						}
					}
				}`,
		},
		{
			name:          "generic error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.bptxmORM.On("EthTransactionsWithAttempts", PageDefaultOffset, PageDefaultLimit).Return(nil, 0, gError)
				f.App.On("BPTXMORM").Return(f.Mocks.bptxmORM)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"ethTransactions"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
