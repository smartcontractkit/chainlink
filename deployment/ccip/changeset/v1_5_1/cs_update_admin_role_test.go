package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUpdateAdminRoleChangesetV2_Validations(t *testing.T) {
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
		Config v1_5_1.UpdateAdminRoleConfig
		MsgStr string
		ErrStr string
	}{
		{
			ErrStr: "admin address cannot be the same as token address",
			MsgStr: "Admin address same as token address",
			Config: v1_5_1.UpdateAdminRoleConfig{
				MCMS: mcmsConfig,
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: tokenAddress, // Same as token
						},
					},
				},
			},
		},
		{
			ErrStr: "does not exist in state",
			MsgStr: "Chain selector is invalid",
			Config: v1_5_1.UpdateAdminRoleConfig{
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
					0: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
		},
		{
			ErrStr: "does not exist in state",
			MsgStr: "Chain selector doesn't exist in environment",
			Config: v1_5_1.UpdateAdminRoleConfig{
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
					5009297550715157269: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
		},
		{
			ErrStr: "admin address cannot be zero",
			MsgStr: "Zero admin address",
			Config: v1_5_1.UpdateAdminRoleConfig{
				MCMS: mcmsConfig,
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorA: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.ZeroAddress,
						},
					},
				},
			},
		},
		{
			ErrStr: "token admin registry failed ownership validation",
			MsgStr: "Ownership validation failure without MCMS",
			Config: v1_5_1.UpdateAdminRoleConfig{
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
					selectorB: {
						{
							TokenAddress: tokenAddress,
							AdminAddress: utils.RandomAddress(),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.MsgStr, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					v1_5_1.UpdateAdminRoleChangesetV2,
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestUpdateAdminRoleChangesetV2_EmptyConfigIsGracefullyHandled(t *testing.T) {
	t.Parallel()

	e, _, _, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.UpdateAdminRoleChangesetV2,
			v1_5_1.UpdateAdminRoleConfig{},
		),
	)
	require.NoError(t, err)
}

func TestUpdateAdminRoleChangesetV2_ExecutionWithoutMCMS(t *testing.T) {
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

	// Manually accept the admin roles by calling acceptAdminRole directly on the registry
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Get chain states
	chainStateA := state.Chains[selectorA]
	chainStateB := state.Chains[selectorB]

	// Only accept admin role for tokenAddressA
	tx, err := chainStateA.TokenAdminRegistry.AcceptAdminRole(
		e.BlockChains.EVMChains()[selectorA].DeployerKey,
		tokenAddressA,
	)
	require.NoError(t, err)
	_, err = e.BlockChains.EVMChains()[selectorA].Confirm(tx)
	require.NoError(t, err)

	// Verify that the deployer is now the active administrator for tokenAddressA
	configA, err := chainStateA.TokenAdminRegistry.GetTokenConfig(&bind.CallOpts{Context: e.GetContext()}, tokenAddressA)
	require.NoError(t, err)
	require.Equal(t, e.BlockChains.EVMChains()[selectorA].DeployerKey.From, configA.Administrator, "deployer should be the active administrator for tokenA")

	// Verify that the deployer is NOT the active administrator for tokenAddressB (should still be pending)
	configB, err := chainStateB.TokenAdminRegistry.GetTokenConfig(&bind.CallOpts{Context: e.GetContext()}, tokenAddressB)
	require.NoError(t, err)
	require.Equal(t, e.BlockChains.EVMChains()[selectorB].DeployerKey.From, configB.PendingAdministrator, "deployer should be the pending administrator for tokenB")

	// This should run TransferAdminRoleChangesetV2 for tokenAddressA and override the pending admin for tokenAddressB
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_5_1.UpdateAdminRoleChangesetV2,
			v1_5_1.UpdateAdminRoleConfig{
				// Need to override the pending admin (DeployerKey) for tokenAddressB
				OverridePendingAdmin: true,
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
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
	state, err = stateview.LoadOnchainState(e)
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

func TestUpdateAdminRoleChangesetV2_ExecutionWithMCMS(t *testing.T) {
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
			v1_5_1.UpdateAdminRoleChangesetV2,
			v1_5_1.UpdateAdminRoleConfig{
				MCMS: mcmsConfig,
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
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

func TestUpdateAdminRoleChangesetV2_MultipleTokensPerChain(t *testing.T) {
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
			v1_5_1.UpdateAdminRoleChangesetV2,
			v1_5_1.UpdateAdminRoleConfig{
				MCMS: mcmsConfig,
				ChainUpdates: map[uint64][]v1_5_1.TokenAdminInfo{
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
