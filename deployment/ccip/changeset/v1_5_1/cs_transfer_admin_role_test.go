package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

func TestTransferAdminRoleChangeset_Validations(t *testing.T) {
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
			Msg: "External admin undefined",
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
			ErrStr: "external admin must be defined",
		},
		{
			Msg: "Not admin",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:          shared.BurnMintTokenPool,
							Version:       deployment.Version1_5_1,
							ExternalAdmin: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "is not the administrator",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.TransferAdminRoleChangeset),
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestTransferAdminRoleChangeset_Execution(t *testing.T) {
	for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} {
		msg := "Transfer admin role with MCMS"
		if mcmsConfig == nil {
			msg = "Transfer admin role without MCMS"
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

			_, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
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
						selectorB: {
							testhelpers.TestTokenSymbol: {
								Type:    shared.BurnMintTokenPool,
								Version: deployment.Version1_5_1,
							},
						},
					},
				},
			), commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.TransferAdminRoleChangeset),
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
			))
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

func TestTransferAdminRoleChangesetV2_EmptyConfigReturnsError(t *testing.T) {
	t.Parallel()

	e, _, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.TransferAdminRoleChangesetV2,
			v1_5_1.TransferAdminRoleConfig{},
		),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "at least one chain with token admin info must be specified")
}

func TestTransferAdminRoleChangesetV2_ExecutionWithMCMS(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	tokenAddressA := utils.RandomAddress()
	tokenAddressB := utils.RandomAddress()
	newAdminA := utils.RandomAddress()
	newAdminB := utils.RandomAddress()

	// First propose admin roles - the timelock becomes the pending admin
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddressA,
							AdminAddress: newAdminA,
						},
					},
					selectorB: {
						{
							TokenAddress: tokenAddressB,
							AdminAddress: newAdminB,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify the proposed admin roles exist
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry
	configOnA, err := registryOnA.GetTokenConfig(nil, tokenAddressA)
	require.NoError(t, err)
	require.Equal(t, newAdminA, configOnA.PendingAdministrator)

	registryOnB := state.Chains[selectorB].TokenAdminRegistry
	configOnB, err := registryOnB.GetTokenConfig(nil, tokenAddressB)
	require.NoError(t, err)
	require.Equal(t, newAdminB, configOnB.PendingAdministrator)
}

func TestTransferAdminRoleChangesetV2_ExecutionWithoutMCMS(t *testing.T) {
	t.Parallel()

	e, selectorA, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	tokenAddressA := utils.RandomAddress()
	tokenAddressB := utils.RandomAddress()
	newAdminA := utils.RandomAddress()
	newAdminB := utils.RandomAddress()

	// First propose admin roles with deployer key (no MCMS) - deployer becomes pending admin
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddressA,
							AdminAddress: e.BlockChains.EVMChains()[selectorA].DeployerKey.From,
						},
					},
					selectorB: {
						{
							TokenAddress: tokenAddressB,
							AdminAddress: e.BlockChains.EVMChains()[selectorB].DeployerKey.From,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Now transfer admin roles to new addresses directly from the deployer (who owns the registry)
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.TransferAdminRoleChangesetV2,
			v1_5_1.TransferAdminRoleConfig{
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddressA,
							AdminAddress: newAdminA,
						},
					},
					selectorB: {
						{
							TokenAddress: tokenAddressB,
							AdminAddress: newAdminB,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify the transfers - new admins should be in pending state
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registryOnA := state.Chains[selectorA].TokenAdminRegistry
	configOnA, err := registryOnA.GetTokenConfig(nil, tokenAddressA)
	require.NoError(t, err)
	require.Equal(t, newAdminA, configOnA.PendingAdministrator)

	registryOnB := state.Chains[selectorB].TokenAdminRegistry
	configOnB, err := registryOnB.GetTokenConfig(nil, tokenAddressB)
	require.NoError(t, err)
	require.Equal(t, newAdminB, configOnB.PendingAdministrator)
}

func TestTransferAdminRoleChangesetV2_MultipleTokensPerChain(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	token1 := utils.RandomAddress()
	token2 := utils.RandomAddress()
	token3 := utils.RandomAddress()
	newAdmin := utils.RandomAddress()

	// First propose admin roles for multiple tokens - timelock becomes pending admin
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{TokenAddress: token1, AdminAddress: newAdmin},
						{TokenAddress: token2, AdminAddress: newAdmin},
						{TokenAddress: token3, AdminAddress: newAdmin},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify all transfers - all tokens should have the new admin as pending
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	registry := state.Chains[selectorA].TokenAdminRegistry
	for _, token := range []common.Address{token1, token2, token3} {
		config, err := registry.GetTokenConfig(nil, token)
		require.NoError(t, err)
		require.Equal(t, newAdmin, config.PendingAdministrator)
	}
}

func TestTransferAdminRoleChangesetV2_Validations(t *testing.T) {
	t.Parallel()

	mcmsConfig := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	e, selectorA, selectorB, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	tokenAddress := utils.RandomAddress()

	// First, set up a token with a pending admin by proposing (timelock becomes pending admin)
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.ProposeAdminRoleChangesetV2,
			v1_5_1.ProposeAdminRoleConfig{
				MCMS: mcmsConfig,
				ProposeAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	tests := []struct {
		Msg    string
		Config v1_5_1.TransferAdminRoleConfig
		ErrStr string
	}{
		{
			Msg: "Empty TransferAdminByChain map",
			Config: v1_5_1.TransferAdminRoleConfig{
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{},
			},
			ErrStr: "at least one chain with token admin info must be specified",
		},
		{
			Msg: "Duplicate token addresses on same chain",
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
						{
							TokenAddress: tokenAddress, // Duplicate
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "duplicate token address",
		},
		{
			Msg: "Admin address same as token address",
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: tokenAddress, // Same as token
						},
					},
				},
			},
			ErrStr: "admin address cannot be the same as token address",
		},
		{
			Msg: "Chain selector is invalid",
			Config: v1_5_1.TransferAdminRoleConfig{
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
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
			Config: v1_5_1.TransferAdminRoleConfig{
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
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
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {},
				},
			},
			ErrStr: "no token admin info provided for chain selector",
		},
		{
			Msg: "Zero token address",
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: utils.ZeroAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "token address cannot be zero",
		},
		{
			Msg: "Zero admin address",
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.ZeroAddress,
						},
					},
				},
			},
			ErrStr: "admin address cannot be zero",
		},
		{
			Msg: "Token with no administrator",
			Config: v1_5_1.TransferAdminRoleConfig{
				MCMS: mcmsConfig,
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: utils.RandomAddress(), // Random token with no admin
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "current administrator is 0x0000000000000000000000000000000000000000",
		},
		{
			Msg: "Ownership validation failure without MCMS",
			Config: v1_5_1.TransferAdminRoleConfig{
				TransferAdminByChain: map[uint64][]v1_5_1.TokenAdminInfo{
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
					v1_5_1.TransferAdminRoleChangesetV2,
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}
