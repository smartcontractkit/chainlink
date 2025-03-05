package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmtoml "github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	chainlinkmocks "github.com/smartcontractkit/chainlink/v2/core/services/chainlink/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/testutils"
)

func TestResolver_Chains(t *testing.T) {
	var (
		chainID = *big.NewI(1)
		query   = `
			query GetChains {
				chains {
					results {
						id
						enabled
						config
						network
					}
					metadata {
						total
					}
				}
			}`
		configTOML = `ChainID = '1'
Enabled = true
AutoCreateKey = false
BlockBackfillDepth = 100
BlockBackfillSkip = true
ChainType = 'Optimism'
FinalityDepth = 42
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 13
MinContractPayment = '9.223372036854775807 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
RPCDefaultBatchSize = 17
RPCBlockQueryDelay = 10
Nodes = []

[Transactions]
ForwardersEnabled = true
MaxInFlight = 19
MaxQueued = 99
ReaperInterval = '1m0s'
ReaperThreshold = '1m0s'
ResendAfterThreshold = '1h0m0s'
`
	)
	var chain evmtoml.EVMConfig
	err := toml.Unmarshal([]byte(configTOML), &chain)
	require.NoError(t, err)

	configTOMLEscaped, err := json.Marshal(configTOML)
	require.NoError(t, err)
	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chains"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				chainConf := evmtoml.EVMConfig{
					ChainID: &chainID,
					Enabled: chain.Enabled,
					Chain:   chain.Chain,
				}

				chainConfToml, err2 := chainConf.TOMLString()
				require.NoError(t, err2)

				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{
					commontypes.RelayID{
						Network: relay.NetworkEVM,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: *chain.Enabled,
						Config:  chainConfToml,
					}},
				}})
			},
			query: query,
			result: fmt.Sprintf(`
			{
				"chains": {
					"results": [{
						"id": "1",
						"enabled": true,
						"config": %s,
						"network": "evm"
					}],
					"metadata": {
						"total": 1
					}
				}
			}`, configTOMLEscaped),
		},
		unauthorizedTestCase(GQLTestCase{query: query}, "chains"),
		{
			name:          "no chains",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{}})
			},
			query: query,
			result: `
			{
				"chains": {
					"results": [],
					"metadata": {
						"total": 0
					}
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_Chain(t *testing.T) {
	var (
		chainID = *big.NewI(1)
		query   = `
			query GetChain {
				chain(id: "1") {
					... on Chain {
						id
						enabled
						config
						network
					}
					... on NotFoundError {
						code
						message
					}
				}
			}
		`
		queryWithNetwork = `
			query GetChain($network: String)  {
				chain(id: "1", network: $network) {
					... on Chain {
						id
						enabled
						config
						network
					}
					... on NotFoundError {
						code
						message
					}
				}
			}
		`
		configTOML = `ChainID = '1'
AutoCreateKey = false
BlockBackfillDepth = 100
BlockBackfillSkip = true
ChainType = 'Optimism'
FinalityDepth = 42
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 13
MinContractPayment = '9.223372036854775807 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
RPCDefaultBatchSize = 17
RPCBlockQueryDelay = 10
Nodes = []

[Transactions]
ForwardersEnabled = true
MaxInFlight = 19
MaxQueued = 99
ReaperInterval = '1m0s'
ReaperThreshold = '1m0s'
ResendAfterThreshold = '1h0m0s'
`
	)
	var chain evmtoml.Chain
	err := toml.Unmarshal([]byte(configTOML), &chain)
	require.NoError(t, err)

	configTOMLEscaped, err := json.Marshal(configTOML)
	require.NoError(t, err)
	multipleChainError := errors.New("multiple chains found with the same chain ID")
	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chain"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				chainConf := evmtoml.EVMConfig{
					Chain:   chain,
					ChainID: &chainID,
				}

				chainConfToml, err2 := chainConf.TOMLString()
				require.NoError(t, err2)

				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{
					commontypes.RelayID{
						Network: relay.NetworkEVM,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: chainConf.IsEnabled(),
						Config:  chainConfToml,
					}},
				}})
			},
			query: query,
			result: fmt.Sprintf(`
				{
					"chain": {
						"id": "1",
						"enabled": true,
						"config": %s,
						"network": "evm"
					}
				}`, configTOMLEscaped),
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{}})
			},
			query: query,
			result: `
				{
					"chain": {
						"code": "NOT_FOUND",
						"message": "chain not found"
					}
				}`,
		},
		{
			name:          "multiple chain with same chainID found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				chainConf := evmtoml.EVMConfig{
					Chain:   chain,
					ChainID: &chainID,
				}

				chainConfToml, err2 := chainConf.TOMLString()
				require.NoError(t, err2)

				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{
					commontypes.RelayID{
						Network: relay.NetworkEVM,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: chainConf.IsEnabled(),
						Config:  chainConfToml,
					}},
					commontypes.RelayID{
						Network: relay.NetworkAptos,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: chainConf.IsEnabled(),
						Config:  chainConfToml,
					}},
				}})
			},
			query:  query,
			result: "null",
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: multipleChainError,
					Path:          []interface{}{"chain"},
					Message:       multipleChainError.Error(),
				},
			},
		},
		{
			name:          "should return aptos chain if network is aptos",
			authenticated: true,
			variables: map[string]interface{}{
				"network": relay.NetworkAptos,
			},
			before: func(ctx context.Context, f *gqlTestFramework) {
				chainConf := evmtoml.EVMConfig{
					Chain:   chain,
					ChainID: &chainID,
				}

				chainConfToml, err2 := chainConf.TOMLString()
				require.NoError(t, err2)

				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[commontypes.RelayID]loop.Relayer{
					commontypes.RelayID{
						Network: relay.NetworkEVM,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: chainConf.IsEnabled(),
						Config:  chainConfToml,
					}},
					commontypes.RelayID{
						Network: relay.NetworkAptos,
						ChainID: chainID.String(),
					}: testutils.MockRelayer{ChainStatus: commontypes.ChainStatus{
						ID:      chainID.String(),
						Enabled: chainConf.IsEnabled(),
						Config:  chainConfToml,
					}},
				}})
			},
			query: queryWithNetwork,
			result: fmt.Sprintf(`
				{
					"chain": {
						"id": "1",
						"enabled": true,
						"config": %s,
						"network": "aptos"
					}
				}`, configTOMLEscaped),
		},
	}

	RunGQLTests(t, testCases)
}
