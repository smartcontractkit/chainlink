package resolver

import (
	"database/sql"
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/stretchr/testify/mock"
)

func Test_CreateFeedsManagerChainConfig(t *testing.T) {
	var (
		mgrID       = int64(100)
		cfgID       = int64(1)
		chainID     = "42"
		accountAddr = "0x0000001"
		adminAddr   = "0x0000002"
		peerID      = null.StringFrom("p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw")
		keyBundleID = null.StringFrom("6fdb8235e16e099de91df7ef8a8088e9deea0ed6ae106b133e5d985a8a9e1562")

		mutation = `
			mutation CreateFeedsManagerChainConfig($input: CreateFeedsManagerChainConfigInput!) {
				createFeedsManagerChainConfig(input: $input) {
					... on CreateFeedsManagerChainConfigSuccess {
						chainConfig {
							id
						  }
					}
					... on NotFoundError {
						message
						code
					}
					... on InputErrors {
						errors {
							path
							message
							code
						}
					}
				}
			}`
		variables = map[string]interface{}{
			"input": map[string]interface{}{
				"feedsManagerID":     stringutils.FromInt64(mgrID),
				"chainID":            chainID,
				"chainType":          "EVM",
				"accountAddr":        accountAddr,
				"adminAddr":          adminAddr,
				"fluxMonitorEnabled": false,
				"ocr1Enabled":        true,
				"ocr1IsBootstrap":    false,
				"ocr1P2PPeerID":      peerID.String,
				"ocr1KeyBundleID":    keyBundleID.String,
				"ocr2Enabled":        false,
			},
		}
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createFeedsManagerChainConfig"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CreateChainConfig", feeds.ChainConfig{
					FeedsManagerID: mgrID,
					ChainType:      feeds.ChainTypeEVM,
					ChainID:        chainID,
					AccountAddress: accountAddr,
					AdminAddress:   adminAddr,
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: false,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled:     true,
						P2PPeerID:   peerID,
						KeyBundleID: keyBundleID,
					},
				}).Return(cfgID, nil)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(&feeds.ChainConfig{
					ID:             cfgID,
					ChainType:      feeds.ChainTypeEVM,
					ChainID:        chainID,
					AccountAddress: accountAddr,
					AdminAddress:   adminAddr,
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: false,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled:     true,
						P2PPeerID:   peerID,
						KeyBundleID: keyBundleID,
					},
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManagerChainConfig": {
					"chainConfig": {
						"id": "1"
					}
				}
			}`,
		},
		{
			name:          "create call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CreateChainConfig", mock.IsType(feeds.ChainConfig{})).Return(int64(0), sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "get call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CreateChainConfig", mock.IsType(feeds.ChainConfig{})).Return(cfgID, nil)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_DeleteFeedsManagerChainConfig(t *testing.T) {
	var (
		cfgID = int64(1)

		mutation = `
			mutation DeleteFeedsManagerChainConfig($id: ID!) {
				deleteFeedsManagerChainConfig(id: $id) {
					... on DeleteFeedsManagerChainConfigSuccess {
						chainConfig {
							id
						}
					}
					... on NotFoundError {
						message
						code
					}
				}
			}`
		variables = map[string]interface{}{
			"id": "1",
		}
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteFeedsManagerChainConfig"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(&feeds.ChainConfig{
					ID: cfgID,
				}, nil)
				f.Mocks.feedsSvc.On("DeleteChainConfig", cfgID).Return(cfgID, nil)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"deleteFeedsManagerChainConfig": {
					"chainConfig": {
						"id": "1"
					}
				}
			}`,
		},
		{
			name:          "get call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"deleteFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "delete call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(&feeds.ChainConfig{
					ID: cfgID,
				}, nil)
				f.Mocks.feedsSvc.On("DeleteChainConfig", cfgID).Return(int64(0), sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"deleteFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_UpdateFeedsManagerChainConfig(t *testing.T) {
	var (
		cfgID = int64(1)

		mutation = `
			mutation UpdateFeedsManagerChainConfig($id: ID!, $input: UpdateFeedsManagerChainConfigInput!) {
				updateFeedsManagerChainConfig(id: $id, input: $input) {
					... on UpdateFeedsManagerChainConfigSuccess {
						chainConfig {
							id
						}
					}
					... on NotFoundError {
						message
						code
					}
					... on InputErrors {
						errors {
							path
							message
							code
						}
					}
				}
			}`
		variables = map[string]interface{}{
			"id": "1",
			"input": map[string]interface{}{
				"accountAddr":        "0x0000001",
				"adminAddr":          "0x0000001",
				"fluxMonitorEnabled": false,
				"ocr1Enabled":        true,
				"ocr1IsBootstrap":    false,
				"ocr1P2PPeerID":      "p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw",
				"ocr1KeyBundleID":    "6fdb8235e16e099de91df7ef8a8088e9deea0ed6ae106b133e5d985a8a9e1562",
				"ocr2Enabled":        false,
			},
		}
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateFeedsManagerChainConfig"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateChainConfig", feeds.ChainConfig{
					ID:             cfgID,
					AccountAddress: "0x0000001",
					AdminAddress:   "0x0000001",
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: false,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled:     true,
						P2PPeerID:   null.StringFrom("p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
						KeyBundleID: null.StringFrom("6fdb8235e16e099de91df7ef8a8088e9deea0ed6ae106b133e5d985a8a9e1562"),
					},
				}).Return(cfgID, nil)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(&feeds.ChainConfig{
					ID:             cfgID,
					AccountAddress: "0x0000001",
					AdminAddress:   "0x0000001",
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: false,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled:     true,
						P2PPeerID:   null.StringFrom("p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
						KeyBundleID: null.StringFrom("6fdb8235e16e099de91df7ef8a8088e9deea0ed6ae106b133e5d985a8a9e1562"),
					},
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateFeedsManagerChainConfig": {
					"chainConfig": {
						"id": "1"
					}
				}
			}`,
		},
		{
			name:          "update call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateChainConfig", mock.IsType(feeds.ChainConfig{})).Return(int64(0), sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "get call not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateChainConfig", mock.IsType(feeds.ChainConfig{})).Return(cfgID, nil)
				f.Mocks.feedsSvc.On("GetChainConfig", cfgID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateFeedsManagerChainConfig": {
					"message": "chain config not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
