package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
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
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateOnRampsDestsChangeset),
					changeset.UpdateOnRampDestsConfig{
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
				),
			)
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

func TestUpdateOnRampDynamicConfig(t *testing.T) {
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
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateOnRampDynamicConfigChangeset),
					changeset.UpdateOnRampDynamicConfig{
						UpdatesByChain: map[uint64]changeset.OnRampDynamicConfigUpdate{
							source: {
								FeeAggregator: common.HexToAddress("0x1002"),
							},
							dest: {
								FeeAggregator: common.HexToAddress("0x2002"),
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// Assert the onramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OnRamp.GetDynamicConfig(&bind.CallOpts{Context: ctx})
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].FeeQuoter.Address(), sourceCfg.FeeQuoter)
			require.Equal(t, common.HexToAddress("0x1002"), sourceCfg.FeeAggregator)
			destCfg, err := state.Chains[dest].OnRamp.GetDynamicConfig(&bind.CallOpts{Context: ctx})
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].FeeQuoter.Address(), destCfg.FeeQuoter)
			require.Equal(t, common.HexToAddress("0x2002"), destCfg.FeeAggregator)
		})
	}
}

func TestUpdateOnRampAllowList(t *testing.T) {
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
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateOnRampAllowListChangeset),
					changeset.UpdateOnRampAllowListConfig{
						UpdatesByChain: map[uint64]map[uint64]changeset.OnRampAllowListUpdate{
							source: {
								dest: {
									AllowListEnabled:          true,
									AddedAllowlistedSenders:   []common.Address{common.HexToAddress("0x1001"), common.HexToAddress("0x1002")},
									RemovedAllowlistedSenders: []common.Address{common.HexToAddress("0x1002"), common.HexToAddress("0x1003")},
								},
							},
							dest: {
								source: {
									AllowListEnabled:          true,
									AddedAllowlistedSenders:   []common.Address{common.HexToAddress("0x2001"), common.HexToAddress("0x2002")},
									RemovedAllowlistedSenders: []common.Address{common.HexToAddress("0x2002"), common.HexToAddress("0x2003")},
								},
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// Assert the onramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OnRamp.GetAllowedSendersList(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Contains(t, sourceCfg.ConfiguredAddresses, common.HexToAddress("0x1001"))
			require.True(t, sourceCfg.IsEnabled)
			destCfg, err := state.Chains[dest].OnRamp.GetAllowedSendersList(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Contains(t, destCfg.ConfiguredAddresses, common.HexToAddress("0x2001"))
			require.True(t, destCfg.IsEnabled)
		})
	}
}

func TestWithdrawOnRampFeeTokens(t *testing.T) {
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

			require.NotNil(t, state.Chains[source].ProposerMcm)
			require.NotNil(t, state.Chains[dest].ProposerMcm)

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

			linkToken := state.Chains[source].LinkToken
			require.NotNil(t, linkToken)
			weth9 := state.Chains[source].Weth9
			require.NotNil(t, weth9)

			// mint some Link and deposit Weth9 to onramp on source chain
			tokenAmount := big.NewInt(100)
			onRamp := state.Chains[source].OnRamp
			config, err := onRamp.GetDynamicConfig(&bind.CallOpts{Context: ctx})
			require.NoError(t, err)
			feeAgggregator := config.FeeAggregator
			deployer := tenv.Env.Chains[source].DeployerKey

			// LINK
			tx, err := linkToken.GrantMintRole(deployer, feeAgggregator)
			require.NoError(t, err)
			_, err = tenv.Env.Chains[source].Confirm(tx)
			require.NoError(t, err)
			tx, err = linkToken.Mint(deployer, onRamp.Address(), tokenAmount)
			require.NoError(t, err)
			_, err = tenv.Env.Chains[source].Confirm(tx)
			require.NoError(t, err)

			// WETH9
			txOpts := *tenv.Env.Chains[source].DeployerKey
			txOpts.Value = tokenAmount
			tx, err = weth9.Deposit(&txOpts)
			require.NoError(t, err)
			_, err = tenv.Env.Chains[source].Confirm(tx)
			require.NoError(t, err)
			tx, err = weth9.Transfer(deployer, onRamp.Address(), tokenAmount)
			require.NoError(t, err)
			_, err = tenv.Env.Chains[source].Confirm(tx)
			require.NoError(t, err)

			// check init balances
			aggregatorInitLinks, err := linkToken.BalanceOf(&bind.CallOpts{Context: ctx}, feeAgggregator)
			require.NoError(t, err)
			require.Equal(t, int64(0), aggregatorInitLinks.Int64())
			aggregatorInitWeth, err := weth9.BalanceOf(&bind.CallOpts{Context: ctx}, feeAgggregator)
			require.NoError(t, err)
			require.Equal(t, int64(0), aggregatorInitWeth.Int64())

			onRampInitLinks, err := linkToken.BalanceOf(&bind.CallOpts{Context: ctx}, onRamp.Address())
			require.NoError(t, err)
			require.Equal(t, tokenAmount, onRampInitLinks)
			onRampInitWeth, err := weth9.BalanceOf(&bind.CallOpts{Context: ctx}, onRamp.Address())
			require.NoError(t, err)
			require.Equal(t, tokenAmount, onRampInitWeth)

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.WithdrawOnRampFeeTokensChangeset),
					changeset.WithdrawOnRampFeeTokensConfig{
						FeeTokensByChain: map[uint64][]common.Address{
							source: {linkToken.Address(), weth9.Address()},
							dest:   {state.Chains[dest].LinkToken.Address(), state.Chains[dest].Weth9.Address()},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// Assert that feeAggregator receives all fee tokens from OnRamp
			aggregatorLinks, err := linkToken.BalanceOf(&bind.CallOpts{Context: ctx}, feeAgggregator)
			require.NoError(t, err)
			assert.Equal(t, tokenAmount, aggregatorLinks)
			aggregatorWeth, err := weth9.BalanceOf(&bind.CallOpts{Context: ctx}, feeAgggregator)
			require.NoError(t, err)
			assert.Equal(t, tokenAmount, aggregatorWeth)
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
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateOffRampSourcesChangeset),
					changeset.UpdateOffRampSourcesConfig{
						UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
							source: {
								dest: {
									IsEnabled:                 true,
									TestRouter:                true,
									IsRMNVerificationDisabled: true,
								},
							},
							dest: {
								source: {
									IsEnabled:                 true,
									TestRouter:                false,
									IsRMNVerificationDisabled: true,
								},
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// Assert the offramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OffRamp.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].TestRouter.Address(), sourceCfg.Router)
			require.True(t, sourceCfg.IsRMNVerificationDisabled)
			require.True(t, sourceCfg.IsEnabled)
			destCfg, err := state.Chains[dest].OffRamp.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
			require.True(t, destCfg.IsRMNVerificationDisabled)
			require.True(t, destCfg.IsEnabled)
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

			fqCfg1 := changeset.DefaultFeeQuoterDestChainConfig(true)
			fqCfg2 := changeset.DefaultFeeQuoterDestChainConfig(true)
			fqCfg2.DestGasOverhead = 1000
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
					changeset.UpdateFeeQuoterDestsConfig{
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
				),
			)
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
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateRouterRampsChangeset),
					changeset.UpdateRouterRampsConfig{
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
				),
			)
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

func TestUpdateDynamicConfigOffRampChangeset(t *testing.T) {
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
			msgInterceptor := utils.RandomAddress()
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateDynamicConfigOffRampChangeset),
					changeset.UpdateDynamicConfigOffRampConfig{
						Updates: map[uint64]changeset.OffRampParams{
							source: {
								PermissionLessExecutionThresholdSeconds: uint32(2 * 60 * 60),
								MessageInterceptor:                      msgInterceptor,
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)
			// Assert the nonce manager configuration is as we expect.
			actualConfig, err := state.Chains[source].OffRamp.GetDynamicConfig(nil)
			require.NoError(t, err)
			require.Equal(t, uint32(2*60*60), actualConfig.PermissionLessExecutionThresholdSeconds)
			require.Equal(t, msgInterceptor, actualConfig.MessageInterceptor)
			require.Equal(t, state.Chains[source].FeeQuoter.Address(), actualConfig.FeeQuoter)
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

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.UpdateNonceManagersChangeset),
					changeset.UpdateNonceManagerConfig{
						UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
							source: {
								RemovedAuthCallers: []common.Address{state.Chains[source].OnRamp.Address()},
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)
			// Assert the nonce manager configuration is as we expect.
			callers, err := state.Chains[source].NonceManager.GetAllAuthorizedCallers(nil)
			require.NoError(t, err)
			require.NotContains(t, callers, state.Chains[source].OnRamp.Address())
			require.Contains(t, callers, state.Chains[source].OffRamp.Address())
		})
	}
}

func TestUpdateNonceManagersCSApplyPreviousRampsUpdates(t *testing.T) {
	e, tenv := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
			RMNConfig: &rmn_contract.RMNConfig{
				BlessWeightThreshold: 2,
				CurseWeightThreshold: 2,
				// setting dummy voters, we will permabless this later
				Voters: []rmn_contract.RMNVoter{
					{
						BlessWeight:   2,
						CurseWeight:   2,
						BlessVoteAddr: utils.RandomAddress(),
						CurseVoteAddr: utils.RandomAddress(),
					},
				},
			},
		}),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithChainIDs([]uint64{chainselectors.GETH_TESTNET.EvmChainID}))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectorsExcluding([]uint64{chainselectors.GETH_TESTNET.Selector})
	require.Contains(t, e.Env.AllChainSelectors(), chainselectors.GETH_TESTNET.Selector)
	require.Len(t, allChains, 2)
	src, dest := allChains[1], chainselectors.GETH_TESTNET.Selector
	srcChain := e.Env.Chains[src]
	destChain := e.Env.Chains[dest]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dest},
	}
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tenv, false)
	// try to apply previous ramps updates without having any previous ramps
	// it should fail
	_, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateNonceManagersChangeset),
			changeset.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: destChain.Selector,
							},
						},
					},
				},
			},
		),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no previous onramp for source chain")
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)
	// Now apply the nonce manager update
	// it should fail again as there is no offramp for the source chain
	_, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateNonceManagersChangeset),
			changeset.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: destChain.Selector,
							},
						},
					},
				},
			},
		),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no previous offramp for source chain")
	// Now apply the update with AllowEmptyOffRamp and it should pass
	_, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateNonceManagersChangeset),
			changeset.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: destChain.Selector,
								AllowEmptyOffRamp:   true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
}

func TestSetOCR3ConfigValidations(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(nil))
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	evmContractParams := make(map[uint64]changeset.ChainContractParams)
	for _, chain := range allChains {
		evmContractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	var apps []commonchangeset.ConfiguredChangeSet
	// now deploy contracts
	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployHomeChainChangeset),
			changeset.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployChainContractsChangeset),
			changeset.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: evmContractParams,
			},
		),
	}...)
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, apps)
	require.NoError(t, err)
	// try to apply ocr3config on offRamp without setting the active config on home chain
	_, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains.
			deployment.CreateLegacyChangeSet(changeset.SetOCR3OffRampChangeset),
			changeset.SetOCR3OffRampConfig{
				HomeChainSel:       e.HomeChainSel,
				RemoteChainSels:    allChains,
				CCIPHomeConfigType: globals.ConfigTypeActive,
			},
		),
	)
	// it should fail as we need to update the chainconfig on CCIPHome first
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid OCR3 config state, expected active config")

	// Build the per chain config.
	wrongChainConfigs := make(map[uint64]changeset.ChainConfig)
	ocrConfigs := make(map[uint64]changeset.CCIPOCRParams)
	for _, chain := range allChains {
		ocrParams := changeset.DeriveCCIPOCRParams(
			changeset.WithDefaultCommitOffChainConfig(e.FeedChainSel, nil),
			changeset.WithDefaultExecuteOffChainConfig(nil),
		)
		ocrConfigs[chain] = ocrParams
		// set wrong chain config with incorrect value of FChain
		wrongChainConfigs[chain] = changeset.ChainConfig{
			Readers: envNodes.NonBootstraps().PeerIDs(),
			//nolint:gosec // disable G115
			FChain: uint8(len(envNodes.NonBootstraps().PeerIDs())),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(globals.GasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(globals.DAGasPriceDeviationPPB)},
				OptimisticConfirmations: globals.OptimisticConfirmations,
			},
		}
	}
	// now set the chain config with wrong values of FChain
	// it should fail on addDonAndSetCandidateChangeset
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// Add the chain configs for the new chains.
			deployment.CreateLegacyChangeSet(changeset.UpdateChainConfigChangeset),
			changeset.UpdateChainConfigConfig{
				HomeChainSelector: e.HomeChainSel,
				RemoteChainAdds:   wrongChainConfigs,
			},
		),
		commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain.
			deployment.CreateLegacyChangeSet(changeset.AddDonAndSetCandidateChangeset),
			changeset.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					FeedChainSelector: e.FeedChainSel,
				},
				PluginInfo: changeset.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "OCR3 config FRoleDON is lower than chainConfig FChain")
}
