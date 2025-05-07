package v1_6_test

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

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
					v1_6.UpdateOnRampDestsConfig{
						UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampDynamicConfigChangeset),
					v1_6.UpdateOnRampDynamicConfig{
						UpdatesByChain: map[uint64]v1_6.OnRampDynamicConfigUpdate{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampAllowListChangeset),
					v1_6.UpdateOnRampAllowListConfig{
						UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampAllowListUpdate{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
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
					cldf.CreateLegacyChangeSet(v1_6.WithdrawOnRampFeeTokensChangeset),
					v1_6.WithdrawOnRampFeeTokensConfig{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
					v1_6.UpdateOffRampSourcesConfig{
						UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			fqCfg1 := v1_6.DefaultFeeQuoterDestChainConfig(true)
			fqCfg2 := v1_6.DefaultFeeQuoterDestChainConfig(true)
			fqCfg2.DestGasOverhead = 1000
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
					v1_6.UpdateFeeQuoterDestsConfig{
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
		name                   string
		transferHomeContracts  bool
		transferChainContracts bool
		mcmsEnabled            bool
		expectedErr            string
	}{
		{
			name:                   "MCMS enabled",
			mcmsEnabled:            true,
			transferHomeContracts:  true,
			transferChainContracts: true,
		},
		{
			name:                   "MCMS disabled",
			mcmsEnabled:            false,
			transferHomeContracts:  false,
			transferChainContracts: false,
		},
		{
			name:                   "MCMS disabled but enforced",
			mcmsEnabled:            false,
			transferHomeContracts:  true,
			transferChainContracts: false,
			expectedErr:            "MCMS is enforced for environment",
		},
		{
			name:                   "MCMS enabled & enforced, but chain contracts not transferred",
			mcmsEnabled:            true,
			transferHomeContracts:  true,
			transferChainContracts: false,
			expectedErr:            "failed to validate ownership of contracts",
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

			if tc.transferHomeContracts {
				var chains []uint64
				if tc.transferChainContracts {
					chains = []uint64{source, dest}
				}
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, chains)
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
					v1_6.UpdateRouterRampsConfig{
						UpdatesByChain: map[uint64]v1_6.RouterUpdates{
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
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				return
			}
			require.NoError(t, err)

			// Assert the router configuration is as we expect.
			source2destOnRampReal, err := state.Chains[source].Router.GetOnRamp(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].OnRamp.Address(), source2destOnRampReal)
			source2destOnRampTest, err := state.Chains[source].TestRouter.GetOnRamp(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, common.HexToAddress("0x0"), source2destOnRampTest)
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			msgInterceptor := utils.RandomAddress()
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateDynamicConfigOffRampChangeset),
					v1_6.UpdateDynamicConfigOffRampConfig{
						Updates: map[uint64]v1_6.OffRampParams{
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
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateNonceManagersChangeset),
					v1_6.UpdateNonceManagerConfig{
						UpdatesByChain: map[uint64]v1_6.NonceManagerUpdate{
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
			cldf.CreateLegacyChangeSet(v1_6.UpdateNonceManagersChangeset),
			v1_6.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]v1_6.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
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
			cldf.CreateLegacyChangeSet(v1_6.UpdateNonceManagersChangeset),
			v1_6.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]v1_6.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
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
			cldf.CreateLegacyChangeSet(v1_6.UpdateNonceManagersChangeset),
			v1_6.UpdateNonceManagerConfig{
				UpdatesByChain: map[uint64]v1_6.NonceManagerUpdate{
					srcChain.Selector: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
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
	evmContractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range allChains {
		evmContractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	var apps []commonchangeset.ConfiguredChangeSet
	// now deploy contracts
	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
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
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
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
			cldf.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset),
			v1_6.SetOCR3OffRampConfig{
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
	wrongChainConfigs := make(map[uint64]v1_6.ChainConfig)
	commitOCRConfigs := make(map[uint64]v1_6.CCIPOCRParams)
	for _, chain := range allChains {
		commitOCRConfigs[chain] = v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, e.FeedChainSel, nil, nil)
		// set wrong chain config with incorrect value of FChain
		wrongChainConfigs[chain] = v1_6.ChainConfig{
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
			cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
			v1_6.UpdateChainConfigConfig{
				HomeChainSelector: e.HomeChainSel,
				RemoteChainAdds:   wrongChainConfigs,
			},
		),
		commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain.
			cldf.CreateLegacyChangeSet(v1_6.AddDonAndSetCandidateChangeset),
			v1_6.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: e.HomeChainSel,
					FeedChainSelector: e.FeedChainSel,
				},
				PluginInfo: v1_6.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: commitOCRConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "OCR3 config FRoleDON is lower than chainConfig FChain")
}

func TestApplyFeeTokensUpdatesFeeQuoterChangeset(t *testing.T) {
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
			allChains := maps.Keys(tenv.Env.Chains)
			// deploy a new token
			ab := deployment.NewMemoryAddressBook()
			for _, selector := range allChains {
				_, err := cldf.DeployContract(tenv.Env.Logger, tenv.Env.Chains[selector], ab,
					func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
						tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
							tenv.Env.Chains[selector].DeployerKey,
							tenv.Env.Chains[selector].Client,
							string(testhelpers.TestTokenSymbol),
							string(testhelpers.TestTokenSymbol),
							testhelpers.LocalTokenDecimals,
							big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
						)
						return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
							Address:  tokenAddress,
							Contract: token,
							Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
							Tx:       tx,
							Err:      err,
						}
					},
				)
				require.NoError(t, err)
			}
			require.NoError(t, tenv.Env.ExistingAddresses.Merge(ab))
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.ApplyFeeTokensUpdatesFeeQuoterChangeset),
					v1_6.ApplyFeeTokensUpdatesConfig{
						UpdatesByChain: map[uint64]v1_6.ApplyFeeTokensUpdatesConfigPerChain{
							source: {
								TokensToAdd:    []changeset.TokenSymbol{testhelpers.TestTokenSymbol},
								TokensToRemove: []changeset.TokenSymbol{changeset.LinkSymbol},
							},
						},
						MCMSConfig: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)
			// Assert the fee quoter configuration is as we expect.
			feeTokens, err := state.Chains[source].FeeQuoter.GetFeeTokens(nil)
			require.NoError(t, err)
			tokenAddresses, err := state.Chains[source].TokenAddressBySymbol()
			require.NoError(t, err)
			require.Contains(t, feeTokens, tokenAddresses[testhelpers.TestTokenSymbol])
			require.NotContains(t, feeTokens, tokenAddresses[changeset.LinkSymbol])
		})
	}
}

func TestApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset(t *testing.T) {
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
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)
			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			// try to update PremiumMultiplierWeiPerEth for a token that does not exist
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.ApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset),
					v1_6.PremiumMultiplierWeiPerEthUpdatesConfig{
						Updates: map[uint64][]v1_6.PremiumMultiplierWeiPerEthUpdatesConfigPerChain{
							source: {
								{
									Token:                      testhelpers.TestTokenSymbol,
									PremiumMultiplierWeiPerEth: 1e18,
								},
							},
						},
						MCMS: mcmsConfig,
					}),
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "token TEST not found in state for chain")
			// deploy test new token
			ab := deployment.NewMemoryAddressBook()
			for _, selector := range allChains {
				_, err := cldf.DeployContract(tenv.Env.Logger, tenv.Env.Chains[selector], ab,
					func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
						tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
							tenv.Env.Chains[selector].DeployerKey,
							tenv.Env.Chains[selector].Client,
							string(testhelpers.TestTokenSymbol),
							string(testhelpers.TestTokenSymbol),
							testhelpers.LocalTokenDecimals,
							big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
						)
						return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
							Address:  tokenAddress,
							Contract: token,
							Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
							Tx:       tx,
							Err:      err,
						}
					},
				)
				require.NoError(t, err)
			}
			require.NoError(t, tenv.Env.ExistingAddresses.Merge(ab))
			state, err = changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)
			// now try to apply the changeset for TEST token
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.ApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset),
					v1_6.PremiumMultiplierWeiPerEthUpdatesConfig{
						Updates: map[uint64][]v1_6.PremiumMultiplierWeiPerEthUpdatesConfigPerChain{
							source: {
								{
									Token:                      testhelpers.TestTokenSymbol,
									PremiumMultiplierWeiPerEth: 1e18,
								},
							},
							dest: {
								{
									Token:                      testhelpers.TestTokenSymbol,
									PremiumMultiplierWeiPerEth: 1e18,
								},
							},
						},
						MCMS: mcmsConfig,
					}),
			)
			require.NoError(t, err)
			tokenAddress, err := state.Chains[source].TokenAddressBySymbol()
			require.NoError(t, err)
			config, err := state.Chains[source].FeeQuoter.GetPremiumMultiplierWeiPerEth(&bind.CallOpts{
				Context: testcontext.Get(t),
			}, tokenAddress[testhelpers.TestTokenSymbol])
			require.NoError(t, err)
			require.Equal(t, uint64(1e18), config)
		})
	}
}

func TestUpdateTokenPriceFeedsFeeQuoterChangeset(t *testing.T) {
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
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]
			// deploy a new token
			ab := deployment.NewMemoryAddressBook()
			_, err := cldf.DeployContract(tenv.Env.Logger, tenv.Env.Chains[source], ab,
				func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
					tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
						tenv.Env.Chains[source].DeployerKey,
						tenv.Env.Chains[source].Client,
						string(testhelpers.TestTokenSymbol),
						string(testhelpers.TestTokenSymbol),
						testhelpers.LocalTokenDecimals,
						big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
					)
					return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			require.NoError(t, err)
			require.NoError(t, tenv.Env.ExistingAddresses.Merge(ab))
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			// try to update price feed for this it will fail as there is no price feed deployed for this token
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateTokenPriceFeedsFeeQuoterChangeset),
					v1_6.UpdateTokenPriceFeedsConfig{
						Updates: map[uint64][]v1_6.UpdateTokenPriceFeedsConfigPerChain{
							source: {
								{
									SourceToken: testhelpers.TestTokenSymbol,
									IsEnabled:   true,
								},
							},
						},
						FeedChainSelector: tenv.FeedChainSel,
						MCMS:              mcmsConfig,
					}),
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "price feed for token TEST not found in state for chain")
			// now try to apply the changeset for link token, there is already a price feed deployed for link token
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateTokenPriceFeedsFeeQuoterChangeset),
					v1_6.UpdateTokenPriceFeedsConfig{
						Updates: map[uint64][]v1_6.UpdateTokenPriceFeedsConfigPerChain{
							source: {
								{
									SourceToken: changeset.LinkSymbol,
									IsEnabled:   true,
								},
							},
							dest: {
								{
									SourceToken: changeset.LinkSymbol,
									IsEnabled:   true,
								},
							},
						},
						FeedChainSelector: tenv.FeedChainSel,
						MCMS:              mcmsConfig,
					}),
			)
			require.NoError(t, err)
			tokenAddress, err := state.Chains[source].TokenAddressBySymbol()
			require.NoError(t, err)
			tokenDetails, err := state.Chains[source].TokenDetailsBySymbol()
			require.NoError(t, err)
			decimals, err := tokenDetails[changeset.LinkSymbol].Decimals(&bind.CallOpts{Context: testcontext.Get(t)})
			require.NoError(t, err)
			config, err := state.Chains[source].FeeQuoter.GetTokenPriceFeedConfig(&bind.CallOpts{
				Context: testcontext.Get(t),
			}, tokenAddress[changeset.LinkSymbol])
			require.NoError(t, err)
			require.True(t, config.IsEnabled)
			require.Equal(t, state.Chains[tenv.FeedChainSel].USDFeeds[changeset.LinkSymbol].Address(), config.DataFeedAddress)
			require.Equal(t, decimals, config.TokenDecimals)
		})
	}
}

func TestApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset(t *testing.T) {
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
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]
			state, err := changeset.LoadOnchainState(tenv.Env)
			require.NoError(t, err)
			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, []uint64{source, dest})
			}

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset),
					v1_6.ApplyTokenTransferFeeConfigUpdatesConfig{
						UpdatesByChain: map[uint64]v1_6.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
							source: {
								TokenTransferFeeConfigRemoveArgs: []v1_6.TokenTransferFeeConfigRemoveArg{
									{
										DestChain: dest,
										Token:     changeset.LinkSymbol,
									},
								},
							},
							dest: {
								TokenTransferFeeConfigArgs: []v1_6.TokenTransferFeeConfigArg{
									{
										DestChain: source,
										TokenTransferFeeConfigPerToken: map[changeset.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig{
											changeset.LinkSymbol: {
												MinFeeUSDCents:    1,
												MaxFeeUSDCents:    1,
												DeciBps:           1,
												DestGasOverhead:   1,
												DestBytesOverhead: 1,
												IsEnabled:         true,
											},
										},
									},
								},
							},
						},
						MCMS: mcmsConfig,
					}),
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "min fee must be less than max fee for token")
			_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset),
					v1_6.ApplyTokenTransferFeeConfigUpdatesConfig{
						UpdatesByChain: map[uint64]v1_6.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
							source: {
								TokenTransferFeeConfigRemoveArgs: []v1_6.TokenTransferFeeConfigRemoveArg{
									{
										DestChain: dest,
										Token:     changeset.LinkSymbol,
									},
								},
							},
							dest: {
								TokenTransferFeeConfigArgs: []v1_6.TokenTransferFeeConfigArg{
									{
										DestChain: source,
										TokenTransferFeeConfigPerToken: map[changeset.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig{
											changeset.LinkSymbol: {
												MinFeeUSDCents:    1,
												MaxFeeUSDCents:    2,
												DeciBps:           1,
												DestGasOverhead:   1,
												DestBytesOverhead: 64,
												IsEnabled:         true,
											},
										},
									},
								},
							},
						},
						MCMS: mcmsConfig,
					}),
			)
			require.NoError(t, err)
		})
	}
}
