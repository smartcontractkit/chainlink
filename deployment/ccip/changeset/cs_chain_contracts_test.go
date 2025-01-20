package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
)

func TestUpdateOnRampsDests(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			// Default env just has 2 chains with all contracts
			// deployed but no lanes.
			tenv, _ := testhelpers.NewMemoryEnvironment(t)
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *changeset.MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &changeset.MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
					Config: changeset.UpdateOnRampDestsConfig{
						UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
							source: {
								dest: {
									IsEnabled:        true,
									TestRouter:       true,
									AllowListEnabled: false,
								},
							},
							dest: {
								source: {
									IsEnabled:        true,
									TestRouter:       false,
									AllowListEnabled: true,
								},
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// Assert the onramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].TestRouter.Address(), sourceCfg.Router)
			require.False(t, sourceCfg.AllowlistEnabled)
			destCfg, err := state.Chains[dest].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
			require.True(t, destCfg.AllowlistEnabled)
		})
	}
}

func TestUpdateOffRampsSources(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t)
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *changeset.MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &changeset.MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOffRampSourcesChangeset),
					Config: changeset.UpdateOffRampSourcesConfig{
						UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
							source: {
								dest: {
									IsEnabled:  true,
									TestRouter: true,
								},
							},
							dest: {
								source: {
									IsEnabled:  true,
									TestRouter: false,
								},
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// Assert the offramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OffRamp.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].TestRouter.Address(), sourceCfg.Router)
			destCfg, err := state.Chains[dest].OffRamp.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
		})
	}
}

func TestUpdateFQDests(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t)
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *changeset.MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &changeset.MCMSConfig{
					MinDelay: 0,
				}
			}

			fqCfg1 := changeset.DefaultFeeQuoterDestChainConfig()
			fqCfg2 := changeset.DefaultFeeQuoterDestChainConfig()
			fqCfg2.DestGasOverhead = 1000
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
					Config: changeset.UpdateFeeQuoterDestsConfig{
						UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
							source: {
								dest: fqCfg1,
							},
							dest: {
								source: fqCfg2,
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// Assert the fq configuration is as we expect.
			source2destCfg, err := state.Chains[source].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			testhelpers.AssertEqualFeeConfig(t, fqCfg1, source2destCfg)
			dest2sourceCfg, err := state.Chains[dest].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			testhelpers.AssertEqualFeeConfig(t, fqCfg2, dest2sourceCfg)
		})
	}
}

func TestUpdateRouterRamps(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t)
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *changeset.MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &changeset.MCMSConfig{
					MinDelay: 0,
				}
			}

			// Updates test router.
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.UpdateRouterRampsChangeset),
					Config: changeset.UpdateRouterRampsConfig{
						TestRouter: true,
						UpdatesByChain: map[uint64]changeset.RouterUpdates{
							source: {
								OffRampUpdates: map[uint64]bool{
									dest: true,
								},
								OnRampUpdates: map[uint64]bool{
									dest: true,
								},
							},
							dest: {
								OffRampUpdates: map[uint64]bool{
									source: true,
								},
								OnRampUpdates: map[uint64]bool{
									source: true,
								},
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// Assert the router configuration is as we expect.
			source2destOnRampTest, err := state.Chains[source].TestRouter.GetOnRamp(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].OnRamp.Address(), source2destOnRampTest)
			source2destOnRampReal, err := state.Chains[source].Router.GetOnRamp(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, common.HexToAddress("0x0"), source2destOnRampReal)
		})
	}
}

func TestUpdateNonceManagersCS(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tenv, _ := testhelpers.NewMemoryEnvironment(t)
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *changeset.MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &changeset.MCMSConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.UpdateNonceManagersChangeset),
					Config: changeset.UpdateNonceManagerConfig{
						UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
							source: {
								RemovedAuthCallers: []common.Address{state.Chains[source].OnRamp.Address()},
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)
			// Assert the nonce manager configuration is as we expect.
			callers, err := state.Chains[source].NonceManager.GetAllAuthorizedCallers(nil)
			require.NoError(t, err)
			require.NotContains(t, callers, state.Chains[source].OnRamp.Address())
			require.Contains(t, callers, state.Chains[source].OffRamp.Address())
		})
	}
}
