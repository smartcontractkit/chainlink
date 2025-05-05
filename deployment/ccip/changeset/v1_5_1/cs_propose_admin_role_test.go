package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestProposeAdminRoleChangeset_Validations(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay: 0 * time.Second,
	}

	// We want an administrator to exist to force failure in the last test
	e, err := commonchangeset.Apply(t, e, timelockContracts,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
			changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
			changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	tests := []struct {
		Config changeset.TokenAdminRegistryChangesetConfig
		ErrStr string
		Msg    string
	}{
		{
			Msg: "Chain selector is invalid",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					0: {},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					5009297550715157269: {},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Ownership validation failure",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {},
				},
			},
			ErrStr: "token admin registry failed ownership validation",
		},
		{
			Msg: "Invalid pool type",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
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
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
							Version: deployment.Version1_0_0,
						},
					},
				},
			},
			ErrStr: "1.0.0 is not a known token pool version",
		},
		{
			Msg: "Admin already exists",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
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
			_, err = commonchangeset.Apply(t, e, timelockContracts,
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
			e, selectorA, selectorB, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               changeset.BurnMintTokenPool,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
				selectorB: {
					Type:               changeset.BurnMintTokenPool,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
			}, mcmsConfig != nil)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			registryOnA := state.Chains[selectorA].TokenAdminRegistry
			registryOnB := state.Chains[selectorB].TokenAdminRegistry

			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
					changeset.TokenAdminRegistryChangesetConfig{
						MCMS: mcmsConfig,
						Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
							selectorA: {
								testhelpers.TestTokenSymbol: {
									Type:    changeset.BurnMintTokenPool,
									Version: deployment.Version1_5_1,
								},
							},
							selectorB: {
								testhelpers.TestTokenSymbol: {
									Type:    changeset.BurnMintTokenPool,
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
				require.Equal(t, e.Chains[selectorA].DeployerKey.From, configOnA.PendingAdministrator)
			}

			configOnB, err := registryOnB.GetTokenConfig(nil, tokens[selectorB].Address)
			require.NoError(t, err)
			if mcmsConfig != nil {
				require.Equal(t, state.Chains[selectorB].Timelock.Address(), configOnB.PendingAdministrator)
			} else {
				require.Equal(t, e.Chains[selectorB].DeployerKey.From, configOnB.PendingAdministrator)
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
			e, selectorA, selectorB, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)
			externalAdminA := utils.RandomAddress()
			externalAdminB := utils.RandomAddress()

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               changeset.BurnMintTokenPool,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
				selectorB: {
					Type:               changeset.BurnMintTokenPool,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
			}, mcmsConfig != nil)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			registryOnA := state.Chains[selectorA].TokenAdminRegistry
			registryOnB := state.Chains[selectorB].TokenAdminRegistry

			_, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
					changeset.TokenAdminRegistryChangesetConfig{
						MCMS: mcmsConfig,
						Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
							selectorA: {
								testhelpers.TestTokenSymbol: {
									Type:          changeset.BurnMintTokenPool,
									Version:       deployment.Version1_5_1,
									ExternalAdmin: externalAdminA,
								},
							},
							selectorB: {
								testhelpers.TestTokenSymbol: {
									Type:          changeset.BurnMintTokenPool,
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
