package changeset

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
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
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
					Config: UpdateOnRampDestsConfig{
						UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
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
			require.Equal(t, false, sourceCfg.AllowlistEnabled)
			destCfg, err := state.Chains[dest].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
			require.Equal(t, true, destCfg.AllowlistEnabled)
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
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateOffRampSources),
					Config: UpdateOffRampSourcesConfig{
						UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
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
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}

			fqCfg1 := DefaultFeeQuoterDestChainConfig()
			fqCfg2 := DefaultFeeQuoterDestChainConfig()
			fqCfg2.DestGasOverhead = 1000
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterDests),
					Config: UpdateFeeQuoterDestsConfig{
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
			AssertEqualFeeConfig(t, fqCfg1, source2destCfg)
			dest2sourceCfg, err := state.Chains[dest].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			AssertEqualFeeConfig(t, fqCfg2, dest2sourceCfg)
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
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}

			// Updates test router.
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
					Config: UpdateRouterRampsConfig{
						TestRouter: true,
						UpdatesByChain: map[uint64]RouterUpdates{
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
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateNonceManagersCS),
					Config: UpdateNonceManagerConfig{
						UpdatesByChain: map[uint64]NonceManagerUpdate{
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

// TODO: skipping all mcms timelock stuff for now
func TestUpdateOnRampsDestsSolana(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	selectors := append(evmSelectors, solChainSelectors...)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	evmChainSel := selectors[0]
	solChainSelector := solChainSelectors[0]
	lggr.Info(fmt.Sprintf("EVM: %v", evmChainSel))
	lggr.Info(fmt.Sprintf("Solana: %v", solChainSelector))

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployHomeChain),
			Config: DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    solChainSelectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    solChainSelectors,
				HomeChainSelector: homeChainSel,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDestsSolana),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
					// source: {
					// 	dest: {
					// 		IsEnabled:        true,
					// 		TestRouter:       false,
					// 		AllowListEnabled: false,
					// 	},
					// },
					solChainSelector: {
						evmChainSel: {
							IsEnabled:        true,
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
				MCMS: nil,
			},
		},
	})

	// UPDATE BINDINGS
	require.NoError(t, err)
	state, _ := LoadOnchainStateSolana(e)
	for chainSel, chain := range e.SolChains {
		ccipRouterId := state.SolChains[chainSel].SolCcipRouter
		var sourceChainStateAccount ccip_router.SourceChain
		err = solCommomUtil.GetAccountDataBorshInto(e.GetContext(), chain.Client, GetEvmSourceChainStatePDA(ccipRouterId, evmChainSel), solRpc.CommitmentConfirmed, &sourceChainStateAccount)
		lggr.Infof(fmt.Sprintf("%+v", sourceChainStateAccount.State))
		require.NoError(t, err)
		require.Equal(t, uint64(1), sourceChainStateAccount.State.MinSeqNr)
		require.Equal(t, true, sourceChainStateAccount.Config.IsEnabled)
	}

}
