package v1_5_1_test

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
)

// configureHybridTokenPoolChains sets up supported chains for hybrid token pools
func configureHybridTokenPoolChains(t *testing.T, e cldf.Environment, selectorA, selectorB uint64) {
	ratelimiterConfig := token_pool.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  big.NewInt(1e18),
		Rate:      big.NewInt(1),
	}

	tokenPoolConfig := map[uint64]v1_5_1.TokenPoolConfig{
		selectorA: {
			Type:    shared.HybridWithExternalMinterFastTransferTokenPool,
			Version: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			ChainUpdates: v1_5_1.RateLimiterPerChain{
				selectorB: v1_5_1.RateLimiterConfig{
					Inbound:  ratelimiterConfig,
					Outbound: ratelimiterConfig,
				},
			},
		},
		selectorB: {
			Type:    shared.HybridWithExternalMinterFastTransferTokenPool,
			Version: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			ChainUpdates: v1_5_1.RateLimiterPerChain{
				selectorA: v1_5_1.RateLimiterConfig{
					Inbound:  ratelimiterConfig,
					Outbound: ratelimiterConfig,
				},
			},
		},
	}

	_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
		v1_5_1.ConfigureTokenPoolContractsConfig{
			TokenSymbol: testhelpers.TestTokenSymbol,
			PoolUpdates: tokenPoolConfig,
		}))
	require.NoError(t, err)
}

func TestHybridTokenPoolUpdateGroupsChangeset_ValidationErrors(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
	externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterA,
		},
		selectorB: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	configureHybridTokenPoolChains(t, e, selectorA, selectorB)

	testCases := []struct {
		name             string
		config           v1_5_1.HybridTokenPoolUpdateGroupsConfig
		expectedErrorMsg string
	}{
		{
			name: "Empty_TokenSymbol",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     "",
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: selectorB,
							Group:               v1_5_1.BurnAndMint,
							RemoteChainSupply:   big.NewInt(1000),
						},
					},
				},
			},
			expectedErrorMsg: "token symbol must be defined",
		},
		{
			name: "Invalid_ContractType",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: selectorB,
							Group:               v1_5_1.BurnAndMint,
							RemoteChainSupply:   big.NewInt(1000),
						},
					},
				},
			},
			expectedErrorMsg: "unsupported contract type",
		},
		{
			name: "Invalid_Group_Value",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: selectorB,
							Group:               v1_5_1.Group(2),
							RemoteChainSupply:   big.NewInt(1000),
						},
					},
				},
			},
			expectedErrorMsg: "invalid group 2, must be 0 (LOCK_AND_RELEASE) or 1 (BURN_AND_MINT)",
		},
		{
			name: "Negative_RemoteChainSupply",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: selectorB,
							Group:               v1_5_1.BurnAndMint,
							RemoteChainSupply:   big.NewInt(-1),
						},
					},
				},
			},
			expectedErrorMsg: "remote chain supply cannot be negative",
		},
		{
			name: "Invalid_ChainSelector",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: 0,
							Group:               v1_5_1.BurnAndMint,
							RemoteChainSupply:   big.NewInt(1000),
						},
					},
				},
			},
			expectedErrorMsg: "invalid remote chain selector 0",
		},
		{
			name: "Empty_GroupUpdates",
			config: v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {},
				},
			},
			expectedErrorMsg: "no group updates specified for chain",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, tc.config))
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErrorMsg)
		})
	}
}

func TestHybridTokenPoolUpdateGroupsChangeset_BasicUpdates(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
	externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterA,
		},
		selectorB: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	configureHybridTokenPoolChains(t, e, selectorA, selectorB)

	testCases := []struct {
		name                string
		updates             map[uint64][]v1_5_1.GroupUpdateConfig
		expectedGroups      map[uint64]map[uint64]v1_5_1.Group
		expectedChainSupply map[uint64]map[uint64]*big.Int
	}{
		{
			name: "Single_Chain_BurnAndMint_Update",
			updates: map[uint64][]v1_5_1.GroupUpdateConfig{
				selectorA: {
					{
						RemoteChainSelector: selectorB,
						Group:               v1_5_1.BurnAndMint,
						RemoteChainSupply:   big.NewInt(5000),
					},
				},
			},
			expectedGroups: map[uint64]map[uint64]v1_5_1.Group{
				selectorA: {selectorB: v1_5_1.BurnAndMint},
			},
			expectedChainSupply: map[uint64]map[uint64]*big.Int{
				selectorA: {selectorB: big.NewInt(5000)},
			},
		},
		{
			name: "Bidirectional_Updates",
			updates: map[uint64][]v1_5_1.GroupUpdateConfig{
				selectorA: {
					{
						RemoteChainSelector: selectorB,
						Group:               v1_5_1.BurnAndMint,
						RemoteChainSupply:   big.NewInt(3000),
					},
				},
				selectorB: {
					{
						RemoteChainSelector: selectorA,
						Group:               v1_5_1.LockAndRelease,
						RemoteChainSupply:   big.NewInt(2000),
					},
				},
			},
			expectedGroups: map[uint64]map[uint64]v1_5_1.Group{
				selectorA: {selectorB: v1_5_1.BurnAndMint},
				selectorB: {selectorA: v1_5_1.LockAndRelease},
			},
			expectedChainSupply: map[uint64]map[uint64]*big.Int{
				selectorA: {selectorB: big.NewInt(3000)},
				selectorB: {selectorA: big.NewInt(2000)},
			},
		},
		{
			name: "Zero_RemoteChainSupply",
			updates: map[uint64][]v1_5_1.GroupUpdateConfig{
				selectorA: {
					{
						RemoteChainSelector: selectorB,
						Group:               v1_5_1.BurnAndMint,
						RemoteChainSupply:   big.NewInt(0),
					},
				},
			},
			expectedGroups: map[uint64]map[uint64]v1_5_1.Group{
				selectorA: {selectorB: v1_5_1.BurnAndMint},
			},
			expectedChainSupply: map[uint64]map[uint64]*big.Int{
				selectorA: {selectorB: big.NewInt(0)},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			config := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates:         tc.updates,
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, config))
			require.NoError(t, err)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)

			for chainSelector, remoteUpdates := range tc.expectedGroups {
				pool := state.Chains[chainSelector].HybridWithExternalMinterFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.HybridWithExternalMinterFastTransferTokenPoolVersion]
				for remoteChainSelector, expectedGroup := range remoteUpdates {
					currentGroup, err := pool.GetGroup(nil, remoteChainSelector)
					require.NoError(t, err)
					require.Equal(t, uint8(expectedGroup), uint8(currentGroup), "Group mismatch for chain %d -> %d", chainSelector, remoteChainSelector)
				}
			}
		})
	}
}

func TestHybridTokenPoolUpdateGroupsChangeset_WithMCMS(t *testing.T) {
	testCases := []struct {
		name            string
		mcmsEnabled     bool
		contractType    cldf.ContractType
		contractVersion semver.Version
	}{
		{
			name:            "HybridWithExternalMinterFastTransferTokenPool",
			mcmsEnabled:     false,
			contractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
			contractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
		},
		{
			name:            "HybridWithExternalMinterFastTransferTokenPool with MCMS",
			mcmsEnabled:     true,
			contractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
			contractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), tc.mcmsEnabled)

			externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
			externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterA,
				},
				selectorB: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterB,
				},
			}, tc.mcmsEnabled)

			configureHybridTokenPoolChains(t, e, selectorA, selectorB)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			config := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    tc.contractType,
				ContractVersion: tc.contractVersion,
				MCMS:            mcmsConfig,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					selectorA: {
						{
							RemoteChainSelector: selectorB,
							Group:               v1_5_1.BurnAndMint,
							RemoteChainSupply:   big.NewInt(10000),
						},
					},
					selectorB: {
						{
							RemoteChainSelector: selectorA,
							Group:               v1_5_1.LockAndRelease,
							RemoteChainSupply:   big.NewInt(5000),
						},
					},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, config))
			require.NoError(t, err)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)

			poolA := state.Chains[selectorA].HybridWithExternalMinterFastTransferTokenPools[testhelpers.TestTokenSymbol][tc.contractVersion]
			groupAB, err := poolA.GetGroup(nil, selectorB)
			require.NoError(t, err)
			require.Equal(t, uint8(v1_5_1.BurnAndMint), uint8(groupAB))

			poolB := state.Chains[selectorB].HybridWithExternalMinterFastTransferTokenPools[testhelpers.TestTokenSymbol][tc.contractVersion]
			groupBA, err := poolB.GetGroup(nil, selectorA)
			require.NoError(t, err)
			require.Equal(t, uint8(v1_5_1.LockAndRelease), uint8(groupBA))
		})
	}
}

func TestHybridTokenPoolUpdateGroupsChangeset_EdgeCases(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
	externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterA,
		},
		selectorB: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	configureHybridTokenPoolChains(t, e, selectorA, selectorB)

	testCases := []struct {
		name             string
		updates          map[uint64][]v1_5_1.GroupUpdateConfig
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Large_RemoteChainSupply",
			updates: map[uint64][]v1_5_1.GroupUpdateConfig{
				selectorA: {
					{
						RemoteChainSelector: selectorB,
						Group:               v1_5_1.BurnAndMint,
						RemoteChainSupply:   new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1e18)),
					},
				},
			},
			expectError: false,
		},
		{
			name: "Nil_RemoteChainSupply_DefaultsToZero",
			updates: map[uint64][]v1_5_1.GroupUpdateConfig{
				selectorA: {
					{
						RemoteChainSelector: selectorB,
						Group:               v1_5_1.BurnAndMint,
						RemoteChainSupply:   nil,
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			config := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates:         tc.updates,
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, config))

			if tc.expectError {
				require.Error(t, err)
				if tc.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tc.expectedErrorMsg)
				}
				return
			}

			require.NoError(t, err)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)

			pool := state.Chains[selectorA].HybridWithExternalMinterFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.HybridWithExternalMinterFastTransferTokenPoolVersion]
			currentGroup, err := pool.GetGroup(nil, selectorB)
			require.NoError(t, err)
			require.Equal(t, uint8(v1_5_1.BurnAndMint), uint8(currentGroup))
		})
	}
}

func TestHybridTokenPoolUpdateGroupsChangeset_NoOpUpdate(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
	externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterA,
		},
		selectorB: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	configureHybridTokenPoolChains(t, e, selectorA, selectorB)

	firstConfig := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
		TokenSymbol:     testhelpers.TestTokenSymbol,
		ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
		ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
		Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
			selectorA: {
				{
					RemoteChainSelector: selectorB,
					Group:               v1_5_1.BurnAndMint,
					RemoteChainSupply:   big.NewInt(1000),
				},
			},
		},
	}

	_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, firstConfig))
	require.NoError(t, err)

	duplicateConfig := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
		TokenSymbol:     testhelpers.TestTokenSymbol,
		ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
		ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
		Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
			selectorA: {
				{
					RemoteChainSelector: selectorB,
					Group:               v1_5_1.BurnAndMint,
					RemoteChainSupply:   big.NewInt(2000),
				},
			},
		},
	}

	_, err = commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.HybridTokenPoolUpdateGroupsChangeset, duplicateConfig))
	require.Error(t, err)
	require.Contains(t, err.Error(), "is already in group 1")
}
