package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

func TestAcceptAdminRoleChangeset_Validations(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

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
					0: map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					5009297550715157269: map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Invalid pool type",
			Config: v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					selectorA: map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
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
					selectorA: map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
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
			Msg: "Not pending admin",
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
			ErrStr: "is not the pending administrator",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestAcceptAdminRoleChangeset_Execution(t *testing.T) {
	for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} {
		msg := "Accept admin role with MCMS"
		if mcmsConfig == nil {
			msg = "Accept admin role without MCMS"
		}

		t.Run(msg, func(t *testing.T) {
			e, selectorA, selectorB, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)

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

			registryOnA := state.MustGetEVMChainState(selectorA).TokenAdminRegistry
			registryOnB := state.MustGetEVMChainState(selectorB).TokenAdminRegistry

			e, err = commonchangeset.Apply(t, e, timelockContracts,
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
				commonchangeset.Configure(
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
				),
			)
			require.NoError(t, err)

			configOnA, err := registryOnA.GetTokenConfig(nil, tokens[selectorA].Address)
			require.NoError(t, err)
			if mcmsConfig != nil {
				require.Equal(t, state.MustGetEVMChainState(selectorA).Timelock.Address(), configOnA.Administrator)
			} else {
				require.Equal(t, e.BlockChains.EVMChains()[selectorA].DeployerKey.From, configOnA.Administrator)
			}

			configOnB, err := registryOnB.GetTokenConfig(nil, tokens[selectorB].Address)
			require.NoError(t, err)
			if mcmsConfig != nil {
				require.Equal(t, state.MustGetEVMChainState(selectorB).Timelock.Address(), configOnB.Administrator)
			} else {
				require.Equal(t, e.BlockChains.EVMChains()[selectorB].DeployerKey.From, configOnB.Administrator)
			}
		})
	}
}
