package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestProposeAdminRoleChangeset_Validations(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay: 0 * time.Second,
	}

	// We want an administrator to exist to force failure in the last test
	e, err := commonchangeset.Apply(t, e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			MCMS: mcmsConfig,
			Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
				selectorA: {
					testhelpers.TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
			},
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			MCMS: mcmsConfig,
			Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
				selectorA: {
					testhelpers.TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
			},
		},
	))
	require.NoError(t, err)

	tests := []struct {
		Config v1_5_1.TokenAdminRegistryChangesetConfig
		ErrStr string
		Msg    string
	}{
		{
			Msg: "Chain selector is invalid",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					0: {},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					5009297550715157269: {},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Ownership validation failure",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: {},
				},
			},
			ErrStr: "token admin registry failed ownership validation",
		},
		{
			Msg: "Invalid pool type",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    "InvalidType",
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
			ErrStr: "InvalidType is not a known token pool type",
		},
		{
			Msg: "Invalid pool version",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    shared.BurnMintTokenPool,
							Version: deployment.Version1_0_0,
						},
					},
				},
			},
			ErrStr: "1.0.0 is not a known token pool version",
		},
		{
			Msg: "Admin already exists",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    shared.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
			ErrStr: "token already has an administrator",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			_, err = commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestProposeAdminRoleChangeset_ExecutionWithoutExternalAdmin(t *testing.T) {
	for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} {
		msg := "Propose admin role without external admin with MCMS"
		if mcmsConfig == nil {
			msg = "Propose admin role without external admin without MCMS"
		}

		t.Run(msg, func(t *testing.T) {
			e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               shared.BurnMintTokenPool,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
				selectorB: {
					Type:               shared.BurnMintTokenPool,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
			}, mcmsConfig != nil)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)

			registryOnA := state.Chains[selectorA].TokenAdminRegistry
			registryOnB := state.Chains[selectorB].TokenAdminRegistry

			e, err = commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
					v1_5_1.TokenAdminRegistryChangesetConfig{
						MCMS: mcmsConfig,
						Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
							selectorA: {
								testhelpers.TestTokenSymbol: {
									Type:    shared.BurnMintTokenPool,
									Version: deployment.Version1_5_1,
								},
							},
							selectorB: {
								testhelpers.TestTokenSymbol: {
									Type:    shared.BurnMintTokenPool,
									Version: deployment.Version1_5_1,
								},
							},
						},
					},
				),
			)
			require.NoError(t, err)

			configOnA, err := registryOnA.GetTokenConfig(nil, tokens[selectorA].Address)
			require.NoError(t, err)
			if mcmsConfig != nil {
				require.Equal(t, state.Chains[selectorA].Timelock.Address(), configOnA.PendingAdministrator)
			} else {
				require.Equal(t, e.BlockChains.EVMChains()[selectorA].DeployerKey.From, configOnA.PendingAdministrator)
			}

			configOnB, err := registryOnB.GetTokenConfig(nil, tokens[selectorB].Address)
			require.NoError(t, err)
			if mcmsConfig != nil {
				require.Equal(t, state.Chains[selectorB].Timelock.Address(), configOnB.PendingAdministrator)
			} else {
				require.Equal(t, e.BlockChains.EVMChains()[selectorB].DeployerKey.From, configOnB.PendingAdministrator)
			}
		})
	}
}

func TestProposeAdminRoleChangeset_ExecutionWithExternalAdmin(t *testing.T) {
	for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} {
		msg := "Propose admin role with external admin with MCMS"
		if mcmsConfig == nil {
			msg = "Propose admin role with external admin without MCMS"
		}

		t.Run(msg, func(t *testing.T) {
			e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)
			externalAdminA := utils.RandomAddress()
			externalAdminB := utils.RandomAddress()

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               shared.BurnMintTokenPool,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
				selectorB: {
					Type:               shared.BurnMintTokenPool,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
			}, mcmsConfig != nil)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)

			registryOnA := state.Chains[selectorA].TokenAdminRegistry
			registryOnB := state.Chains[selectorB].TokenAdminRegistry

			_, err = commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
					v1_5_1.TokenAdminRegistryChangesetConfig{
						MCMS: mcmsConfig,
						Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
							selectorA: {
								testhelpers.TestTokenSymbol: {
									Type:          shared.BurnMintTokenPool,
									Version:       deployment.Version1_5_1,
									ExternalAdmin: externalAdminA,
								},
							},
							selectorB: {
								testhelpers.TestTokenSymbol: {
									Type:          shared.BurnMintTokenPool,
									Version:       deployment.Version1_5_1,
									ExternalAdmin: externalAdminB,
								},
							},
						},
					},
				),
			)
			require.NoError(t, err)

			configOnA, err := registryOnA.GetTokenConfig(nil, tokens[selectorA].Address)
			require.NoError(t, err)
			require.Equal(t, externalAdminA, configOnA.PendingAdministrator)

			configOnB, err := registryOnB.GetTokenConfig(nil, tokens[selectorB].Address)
			require.NoError(t, err)
			require.Equal(t, externalAdminB, configOnB.PendingAdministrator)
		})
	}
}

func TestProposeAdminRoleChangesetV2_Validations(t *testing.T) {
	t.Parallel()

	e, _, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	tokenAddress := utils.RandomAddress()

	tests := []struct {
		Config v1_5_1.ProposeAdminRoleConfig
		ErrStr string
		Msg    string
	}{
		{
			Msg: "Empty ProposeAdminByChain map",
			Config: v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{},
			},
			ErrStr: "at least one chain with token admin info must be specified",
		},

		{
			Msg: "Duplicate token addresses on same chain",
			Config: v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
						{
							TokenAddress: tokenAddress, // Same token address
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "duplicate token address",
		},
		{
			Msg: "Admin address same as token address",
			Config: v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: tokenAddress, // Same as token address
						},
					},
				},
			},
			ErrStr: "admin address cannot be the same as token address",
		},
		{
			Msg: "Chain selector is invalid",
			Config: v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					0: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "does not exist in state",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Config: v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					5009297550715157269: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "does not exist in state",
		},
		{
			Msg: "Empty token admin info array",
			Config: v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {},
				},
			},
			ErrStr: "no token admin info provided for chain selector",
		},
		{
			Msg: "Zero token address",
			Config: v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: utils.ZeroAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "token address cannot be zero for propose admin role",
		},
		{
			Msg: "Zero admin address",
			Config: v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.ZeroAddress,
						},
					},
				},
			},
			ErrStr: "admin address cannot be zero for propose admin role",
		},
		{
			Msg: "Ownership validation failure without MCMS",
			Config: v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "token admin registry failed ownership validation",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					v1_5_1.ProposeAdminRoleChangesetV2,
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestProposeAdminRoleChangesetV2_ExecutionWithMCMS(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	// Use random addresses - no need to deploy actual tokens
	tokenAddressA := utils.RandomAddress()
	tokenAddressB := utils.RandomAddress()
	adminAddressA := utils.RandomAddress()
	adminAddressB := utils.RandomAddress()

	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddressA,
							AdminAddress: adminAddressA,
						},
					},
					selectorB: {
						{
							TokenAddress: tokenAddressB,
							AdminAddress: adminAddressB,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify that the admin proposals were created correctly
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry
	registryOnB := state.Chains[selectorB].TokenAdminRegistry

	configOnA, err := registryOnA.GetTokenConfig(nil, tokenAddressA)
	require.NoError(t, err)
	require.Equal(t, adminAddressA, configOnA.PendingAdministrator)

	configOnB, err := registryOnB.GetTokenConfig(nil, tokenAddressB)
	require.NoError(t, err)
	require.Equal(t, adminAddressB, configOnB.PendingAdministrator)
}

func TestProposeAdminRoleChangesetV2_ExecutionWithoutMCMS(t *testing.T) {
	t.Parallel()

	e, selectorA, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	// Use random addresses - no need to deploy actual tokens
	tokenAddressA := utils.RandomAddress()
	tokenAddressB := utils.RandomAddress()
	adminAddressA := utils.RandomAddress()
	adminAddressB := utils.RandomAddress()

	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddressA,
							AdminAddress: adminAddressA,
						},
					},
					selectorB: {
						{
							TokenAddress: tokenAddressB,
							AdminAddress: adminAddressB,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify that the admin proposals were created correctly
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry
	registryOnB := state.Chains[selectorB].TokenAdminRegistry

	configOnA, err := registryOnA.GetTokenConfig(nil, tokenAddressA)
	require.NoError(t, err)
	require.Equal(t, adminAddressA, configOnA.PendingAdministrator)

	configOnB, err := registryOnB.GetTokenConfig(nil, tokenAddressB)
	require.NoError(t, err)
	require.Equal(t, adminAddressB, configOnB.PendingAdministrator)
}

func TestProposeAdminRoleChangesetV2_MultipleTokensPerChain(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	// Use random addresses for multiple tokens - no need to deploy actual tokens
	tokenAddress1 := utils.RandomAddress()
	tokenAddress2 := utils.RandomAddress()
	tokenAddress3 := utils.RandomAddress()
	adminAddress1 := utils.RandomAddress()
	adminAddress2 := utils.RandomAddress()
	adminAddress3 := utils.RandomAddress()

	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress1,
							AdminAddress: adminAddress1,
						},
						{
							TokenAddress: tokenAddress2,
							AdminAddress: adminAddress2,
						},
						{
							TokenAddress: tokenAddress3,
							AdminAddress: adminAddress3,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify that all admin proposals were created correctly
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry

	config1, err := registryOnA.GetTokenConfig(nil, tokenAddress1)
	require.NoError(t, err)
	require.Equal(t, adminAddress1, config1.PendingAdministrator)

	config2, err := registryOnA.GetTokenConfig(nil, tokenAddress2)
	require.NoError(t, err)
	require.Equal(t, adminAddress2, config2.PendingAdministrator)

	config3, err := registryOnA.GetTokenConfig(nil, tokenAddress3)
	require.NoError(t, err)
	require.Equal(t, adminAddress3, config3.PendingAdministrator)
}

func TestProposeAdminRoleChangesetV2_EmptyConfigReturnsError(t *testing.T) {
	t.Parallel()

	e, _, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	// Test that empty config returns error as expected by the validation logic
	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{},
			},
		),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "at least one chain with token admin info must be specified")
}

func TestProposeAdminRoleChangesetV2_PendingAdminValidation(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	tokenAddress := utils.RandomAddress()
	adminAddress1 := utils.RandomAddress()
	adminAddress2 := utils.RandomAddress()

	// First, propose an admin for the token
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: adminAddress1,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Now try to propose another admin for the same token without override - this should fail
	_, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: adminAddress2,
						},
					},
				},
			},
		),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "already has a pending administrator")
	require.ErrorContains(t, err, "Set OverridePendingAdmin=true to override")

	// Now try with override enabled - this should succeed
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS:                 mcmsConfig,
				OverridePendingAdmin: true, // Enable override
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: adminAddress2,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify that the pending admin was actually changed
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry
	config, err := registryOnA.GetTokenConfig(nil, tokenAddress)
	require.NoError(t, err)
	require.Equal(t, adminAddress2, config.PendingAdministrator, "Pending administrator should be updated to the new admin address")
}

func TestProposeAdminRoleChangesetV2_OverrideFunctionality(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	// Test case 1: Override with same admin address should fail
	tokenAddress1 := utils.RandomAddress()
	adminAddress := utils.RandomAddress()

	// First, propose an admin
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress1,
							AdminAddress: adminAddress,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Try to override with the same admin address - should fail
	_, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS:                 mcmsConfig,
				OverridePendingAdmin: true,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress1,
							AdminAddress: adminAddress, // Same admin address
						},
					},
				},
			},
		),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "is already the pending administrator")

	// Test case 2: Multiple tokens, some with pending admins, some without
	tokenAddress2 := utils.RandomAddress() // New token, no pending admin
	tokenAddress3 := utils.RandomAddress() // Will have pending admin
	newAdminAddress := utils.RandomAddress()

	// Set up a pending admin for tokenAddress3
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress3,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Now propose admins for multiple tokens: one new, one with override
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS:                 mcmsConfig,
				OverridePendingAdmin: true,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress2, // New token
							AdminAddress: newAdminAddress,
						},
						{
							TokenAddress: tokenAddress3, // Override existing pending admin
							AdminAddress: newAdminAddress,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify both tokens have the correct pending admin
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registry := state.Chains[selectorA].TokenAdminRegistry

	config2, err := registry.GetTokenConfig(nil, tokenAddress2)
	require.NoError(t, err)
	require.Equal(t, newAdminAddress, config2.PendingAdministrator)

	config3, err := registry.GetTokenConfig(nil, tokenAddress3)
	require.NoError(t, err)
	require.Equal(t, newAdminAddress, config3.PendingAdministrator)
}
