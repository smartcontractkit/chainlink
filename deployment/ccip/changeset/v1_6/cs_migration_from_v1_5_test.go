package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_deploy "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	cciptypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
)

const (
	// PriceRegistry
	stalenessThreshold = 90_000

	// EVM2EVMOnRamp Static Config
	defaultTxGasLimit       = 200_000
	maxNumberOfTokensPerMsg = 5

	// EVM2EVMOnRamp Dynamic Config
	destGasOverhead                   = 350_000
	destGasPerPayloadByte             = 16
	destDataAvailabilityOverheadGas   = 33_596
	destGasPerDataAvailabilityByte    = 16
	destDataAvailabilityMultiplierBps = 6_840
	maxDataBytes                      = 100_000
	maxPerMsgGasLimit                 = 4_000_000
	defaultTokenFeeUSDCents           = 50
	defaultTokenDestGasOverhead       = 32

	// LINK Fee Token Config Args
	linkNetworkFeeUSDCents         = 1_00
	linkGasMultiplierWeiPerEth     = 1e18
	linkPremiumMultiplierWeiPerEth = 9e17

	// WETH Fee Token Config Args
	wethNetworkFeeUSDCents         = 2_00
	wethGasMultiplierWeiPerEth     = 1e17
	wethPremiumMultiplierWeiPerEth = 8e17

	// LINK Transfer Fee Config Args
	linkMinFeeUSDCents    = 50           // $0.5
	linkMaxFeeUSDCents    = 1_000_000_00 // $ 1 million
	linkDeciBps           = 5_0          // 5 bps
	linkDestGasOverhead   = 110_000
	linkDestBytesOverhead = 32
)

var (
	maxNopFeesJuels = big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18))
)

func initMigrationEnvironment(t *testing.T, numChains int, mcmsCfg proposalutils.TimelockConfig) cldf_deploy.Environment {
	dEnv, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = numChains
	})
	e := dEnv.Env
	chainSels := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily("evm"))

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		t.Fatalf("Failed to load onchain state: %v", err)
	}
	homeChainSel, err := state.HomeChainSelector()
	if err != nil {
		t.Fatalf("Failed to get home chain selector: %v", err)
	}

	for _, sel := range chainSels {
		var err error

		// Transfer home chain contracts to MCMS timelock
		if sel == homeChainSel {
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2), commonchangeset.TransferToMCMSWithTimelockConfig{
					MCMSConfig: mcmsCfg,
					ContractsByChain: map[uint64][]common.Address{
						sel: {
							state.Chains[sel].CapabilityRegistry.Address(),
							state.Chains[sel].CCIPHome.Address(),
						},
					},
				}),
			})
			if err != nil {
				t.Fatalf("Failed to transfer home chain contracts to MCMS timelock: %v", err)
			}
		}

		// Transfer NonceManager, FeeQuoter, Router, & RMN Proxy to MCMS timelock on all chains
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2), commonchangeset.TransferToMCMSWithTimelockConfig{
				MCMSConfig: mcmsCfg,
				ContractsByChain: map[uint64][]common.Address{
					sel: {
						state.Chains[sel].NonceManager.Address(),
						state.Chains[sel].FeeQuoter.Address(),
						state.Chains[sel].RMNProxy.Address(),
						state.Chains[sel].Router.Address(),
					},
				},
			}),
		})
		if err != nil {
			t.Fatalf("Failed to transfer NonceManager and FeeQuoter to MCMS timelock: %v", err)
		}

		// Deploy a PriceRegistry 1.2.0
		priceRegDeploy, err := cldf_deploy.DeployContract(e.Logger, e.BlockChains.EVMChains()[sel], e.ExistingAddresses, func(chain evm.Chain) cldf_deploy.ContractDeploy[*price_registry.PriceRegistry] {
			addr, tx, registry, err := price_registry.DeployPriceRegistry(chain.DeployerKey, chain.Client, []common.Address{}, []common.Address{
				state.Chains[sel].LinkToken.Address(),
				state.Chains[sel].Weth9.Address(),
			}, stalenessThreshold)
			return cldf_deploy.ContractDeploy[*price_registry.PriceRegistry]{
				Address:  addr,
				Tx:       tx,
				Tv:       cldf_deploy.NewTypeAndVersion(shared.PriceRegistry, deployment.Version1_2_0),
				Contract: registry,
				Err:      err,
			}
		})
		if err != nil {
			t.Fatalf("Failed to deploy PriceRegistry 1.2.0 on chain %d: %v", sel, err)
		}
		// Deploy one EVM2EVMOnRamp 1.5.0 & one EVM2EVMOffRamp for each of the other chains
		for _, otherSel := range chainSels {
			if otherSel == sel {
				continue // Skip self
			}
			_, err = cldf_deploy.DeployContract(e.Logger, e.BlockChains.EVMChains()[sel], e.ExistingAddresses, func(chain evm.Chain) cldf_deploy.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp] {
				addr, tx, onRamp, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(chain.DeployerKey, chain.Client,
					evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
						LinkToken:          state.Chains[sel].LinkToken.Address(),
						ChainSelector:      sel,
						DestChainSelector:  otherSel,
						DefaultTxGasLimit:  defaultTxGasLimit,
						MaxNopFeesJuels:    maxNopFeesJuels,
						PrevOnRamp:         utils.ZeroAddress,
						RmnProxy:           state.Chains[sel].RMNProxy.Address(),
						TokenAdminRegistry: state.Chains[sel].TokenAdminRegistry.Address(),
					},
					evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
						Router:                            state.Chains[sel].Router.Address(),
						MaxNumberOfTokensPerMsg:           maxNumberOfTokensPerMsg,
						DestGasOverhead:                   destGasOverhead,
						DestGasPerPayloadByte:             destGasPerPayloadByte,
						DestDataAvailabilityOverheadGas:   destDataAvailabilityOverheadGas,
						DestGasPerDataAvailabilityByte:    destGasPerDataAvailabilityByte,
						DestDataAvailabilityMultiplierBps: destDataAvailabilityMultiplierBps,
						PriceRegistry:                     priceRegDeploy.Address,
						MaxDataBytes:                      maxDataBytes,
						MaxPerMsgGasLimit:                 maxPerMsgGasLimit,
						DefaultTokenFeeUSDCents:           defaultTokenFeeUSDCents,
						DefaultTokenDestGasOverhead:       defaultTokenDestGasOverhead,
					},
					evm_2_evm_onramp.RateLimiterConfig{
						IsEnabled: false,
						Capacity:  big.NewInt(0),
						Rate:      big.NewInt(0),
					},
					[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
						{
							Token:                      state.Chains[sel].LinkToken.Address(),
							NetworkFeeUSDCents:         linkNetworkFeeUSDCents,
							GasMultiplierWeiPerEth:     linkGasMultiplierWeiPerEth,
							PremiumMultiplierWeiPerEth: linkPremiumMultiplierWeiPerEth,
							Enabled:                    true,
						},
						{
							Token:                      state.Chains[sel].Weth9.Address(),
							NetworkFeeUSDCents:         wethNetworkFeeUSDCents,
							GasMultiplierWeiPerEth:     wethGasMultiplierWeiPerEth,
							PremiumMultiplierWeiPerEth: wethPremiumMultiplierWeiPerEth,
							Enabled:                    true,
						},
					},
					[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
						{
							Token:                     state.Chains[sel].LinkToken.Address(),
							MinFeeUSDCents:            linkMinFeeUSDCents,
							MaxFeeUSDCents:            linkMaxFeeUSDCents,
							DeciBps:                   linkDeciBps,
							DestGasOverhead:           linkDeciBps,
							DestBytesOverhead:         linkDestBytesOverhead,
							AggregateRateLimitEnabled: true,
						},
					},
					[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
				)
				return cldf_deploy.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp]{
					Address:  addr,
					Tx:       tx,
					Tv:       cldf_deploy.NewTypeAndVersion(shared.OnRamp, deployment.Version1_5_0),
					Contract: onRamp,
					Err:      err,
				}
			})
			if err != nil {
				t.Fatalf("Failed to deploy EVM2EVMOnRamp 1.5.0 on chain %d for %d: %v", sel, otherSel, err)
			}
			commitStoreDeploy, err := cldf_deploy.DeployContract(e.Logger, e.BlockChains.EVMChains()[sel], e.ExistingAddresses, func(chain evm.Chain) cldf_deploy.ContractDeploy[*commit_store.CommitStore] {
				addr, tx, commitStore, err := commit_store.DeployCommitStore(chain.DeployerKey, chain.Client,
					commit_store.CommitStoreStaticConfig{
						ChainSelector:       sel,
						SourceChainSelector: otherSel,
						OnRamp:              utils.RandomAddress(), // Placeholder, not relevant for this test
						RmnProxy:            state.Chains[sel].RMNProxy.Address(),
					},
				)
				return cldf_deploy.ContractDeploy[*commit_store.CommitStore]{
					Address:  addr,
					Tx:       tx,
					Tv:       cldf_deploy.NewTypeAndVersion(shared.CommitStore, deployment.Version1_5_0),
					Contract: commitStore,
					Err:      err,
				}
			})
			if err != nil {
				t.Fatalf("Failed to deploy CommitStore 1.5.0 on chain %d for %d: %v", sel, otherSel, err)
			}
			_, err = cldf_deploy.DeployContract(e.Logger, e.BlockChains.EVMChains()[sel], e.ExistingAddresses, func(chain evm.Chain) cldf_deploy.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp] {
				addr, tx, offRamp, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(chain.DeployerKey, chain.Client,
					evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
						CommitStore:         commitStoreDeploy.Address,
						ChainSelector:       sel,
						SourceChainSelector: otherSel,
						OnRamp:              utils.RandomAddress(), // Placeholder, not relevant for this test
						PrevOffRamp:         utils.ZeroAddress,
						RmnProxy:            state.Chains[sel].RMNProxy.Address(),
						TokenAdminRegistry:  state.Chains[sel].TokenAdminRegistry.Address(),
					},
					evm_2_evm_offramp.RateLimiterConfig{
						IsEnabled: false,
						Capacity:  big.NewInt(0),
						Rate:      big.NewInt(0),
					},
				)
				return cldf_deploy.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp]{
					Address:  addr,
					Tx:       tx,
					Tv:       cldf_deploy.NewTypeAndVersion(shared.OffRamp, deployment.Version1_5_0),
					Contract: offRamp,
					Err:      err,
				}
			})
			if err != nil {
				t.Fatalf("Failed to deploy EVM2EVMOffRamp 1.5.0 on chain %d for %d: %v", sel, otherSel, err)
			}
		}
	}

	return e
}

func TestInitAndPromoteChainUpgrades(t *testing.T) {
	mcmsCfg := proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: mcmstypes.TimelockActionSchedule,
	}
	callOpts := &bind.CallOpts{Context: t.Context()}

	e := initMigrationEnvironment(t, 3, mcmsCfg)
	require.Len(t, e.BlockChains.EVMChains(), 3, "Expected 3 EVM chains in the environment")

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err, "Failed to load onchain state")

	homeChainSelector, err := state.HomeChainSelector()
	require.NoError(t, err, "Failed to get home chain selector")
	feedChainSelector := homeChainSelector // Just use home chain selector as feed chain selector for this test

	chainUpgradeCfgs := make(map[uint64]v1_6.ChainUpgradeConfig, len(e.BlockChains.EVMChains()))
	for _, chain := range e.BlockChains.EVMChains() {
		chainUpgradeCfgs[chain.Selector] = v1_6.ChainUpgradeConfig{
			FeedChainSelector: feedChainSelector,
			CommitOCRParams:   v1_6.DefaultOCRParamsForCommitForETH,
			ExecOCRParams:     v1_6.DefaultOCRParamsForExecForETH,
		}
	}

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(v1_6.InitChainUpgradesChangeset, v1_6.InitChainUpgradesConfig{
			HomeChainSelector: homeChainSelector,
			ChainsToUpgrade:   chainUpgradeCfgs,
			MCMSConfig:        &mcmsCfg,
		}),
	})
	require.NoError(t, err, "Failed to apply InitChainUpgradesChangeset")

	// InitChainUpgradesChangeset checks
	state, err = stateview.LoadOnchainState(e)
	require.NoError(t, err, "Failed to load onchain state")
	for _, chain := range e.BlockChains.EVMChains() {
		// Commit and exec candidates are set on CCIPHome
		donID, err := internal.DonIDForChain(
			state.Chains[homeChainSelector].CapabilityRegistry,
			state.Chains[homeChainSelector].CCIPHome,
			chain.Selector,
		)
		require.NoError(t, err, "Failed to get DON ID for chain %d", chain.Selector)
		commitCandidate, err := state.Chains[homeChainSelector].CCIPHome.GetCandidateDigest(callOpts, donID, uint8(cciptypes.PluginTypeCCIPCommit))
		require.NoError(t, err, "Failed to get commit candidate for chain %d", chain.Selector)
		require.NotEqual(t, [32]byte{}, commitCandidate, "Commit candidate should not be empty for chain %d", chain.Selector)
		execCandidate, err := state.Chains[homeChainSelector].CCIPHome.GetCandidateDigest(callOpts, donID, uint8(cciptypes.PluginTypeCCIPExec))
		require.NoError(t, err, "Failed to get exec candidate for chain %d", chain.Selector)
		require.NotEqual(t, [32]byte{}, execCandidate, "Exec candidate should not be empty for chain %d", chain.Selector)

		// RMNRemote is owned by the MCMS timelock
		owner, err := state.Chains[chain.Selector].RMNRemote.Owner(callOpts)
		require.NoError(t, err, "Failed to get RMNRemote owner for chain %d", chain.Selector)
		require.Equal(t, state.Chains[chain.Selector].Timelock.Address(), owner, "RMNRemote owner should be MCMS timelock for chain %d", chain.Selector)

		// RMNProxy is pointing at the RMNRemote
		rmnOnProxyAddr, err := state.Chains[chain.Selector].RMNProxy.GetARM(callOpts)
		require.NoError(t, err, "Failed to get RMNProxy ARM for chain %d", chain.Selector)
		require.Equal(t, state.Chains[chain.Selector].RMNRemote.Address(), rmnOnProxyAddr, "RMNProxy should point to RMNRemote for chain %d", chain.Selector)

		// PremiumMultiplierWeiPerEth is set for WETH and LINK
		/*
			TODO: Needs fixing

			premiumMultiplierWeiPerEth, err := state.Chains[chain.Selector].FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, state.Chains[chain.Selector].LinkToken.Address())
			require.NoError(t, err, "Failed to get PremiumMultiplierWeiPerEth for LINK on chain %d", chain.Selector)
			require.Equal(t, uint64(linkPremiumMultiplierWeiPerEth), premiumMultiplierWeiPerEth, "LINK PremiumMultiplierWeiPerEth should match for chain %d", chain.Selector)
			premiumMultiplierWeiPerEth, err = state.Chains[chain.Selector].FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, state.Chains[chain.Selector].Weth9.Address())
			require.NoError(t, err, "Failed to get PremiumMultiplierWeiPerEth for WETH on chain %d", chain.Selector)
			require.Equal(t, uint64(wethPremiumMultiplierWeiPerEth), premiumMultiplierWeiPerEth, "WETH PremiumMultiplierWeiPerEth should match for chain %d", chain.Selector)
		*/

		for _, otherChain := range e.BlockChains.EVMChains() {
			if otherChain.Selector == chain.Selector {
				continue // Skip self
			}

			// TransferFeeConfigArgs are set on the FeeQuoter for LINK
			/*
				TODO: Needs fixing

				transferFeeConfigArgs, err := state.Chains[chain.Selector].FeeQuoter.GetTokenTransferFeeConfig(callOpts, otherChain.Selector, state.Chains[chain.Selector].LinkToken.Address())
				require.NoError(t, err, "Failed to get LINK transfer fee config for chain %d for %d", chain.Selector, otherChain.Selector)

				require.Equal(t, uint32(linkMinFeeUSDCents), transferFeeConfigArgs.MinFeeUSDCents, "LINK MinFeeUSDCents should match for chain %d for %d", chain.Selector, otherChain.Selector)
				require.Equal(t, uint32(linkMaxFeeUSDCents), transferFeeConfigArgs.MaxFeeUSDCents, "LINK MaxFeeUSDCents should match for chain %d for %d", chain.Selector, otherChain.Selector)
				require.Equal(t, uint16(linkDeciBps), transferFeeConfigArgs.DeciBps, "LINK DeciBps should match for chain %d for %d", chain.Selector, otherChain.Selector)
				require.Equal(t, uint32(linkDestGasOverhead), transferFeeConfigArgs.DestGasOverhead, "LINK DestGasOverhead should match for chain %d for %d", chain.Selector, otherChain.Selector)
				require.Equal(t, uint32(linkDestBytesOverhead), transferFeeConfigArgs.DestBytesOverhead, "LINK DestBytesOverhead should match for chain %d for %d", chain.Selector, otherChain.Selector)
				require.True(t, transferFeeConfigArgs.IsEnabled, "LINK Transfer fee config should be enabled for chain %d for %d", chain.Selector, otherChain.Selector)
			*/

			// Fee tokens are set on the fee quoter
			feeTokens, err := state.Chains[chain.Selector].FeeQuoter.GetFeeTokens(callOpts)
			require.NoError(t, err, "Failed to get fee tokens for chain %d", chain.Selector)
			require.Len(t, feeTokens, 2, "Expected 2 fee tokens for chain %d", chain.Selector)
			require.Contains(t, feeTokens, state.Chains[chain.Selector].LinkToken.Address(), "Fee tokens should contain LINK for chain %d", chain.Selector)
			require.Contains(t, feeTokens, state.Chains[chain.Selector].Weth9.Address(), "Fee tokens should contain WETH for chain %d", chain.Selector)

			// DestChainConfig is set for other chains
			fqDestChainConfig, err := state.Chains[chain.Selector].FeeQuoter.GetDestChainConfig(callOpts, otherChain.Selector)
			require.NoError(t, err, "Failed to get dest chain config for chain %d for %d", chain.Selector, otherChain.Selector)
			require.True(t, fqDestChainConfig.IsEnabled, "DestChainConfig should be enabled for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint16(maxNumberOfTokensPerMsg), fqDestChainConfig.MaxNumberOfTokensPerMsg, "DestChainConfig MaxNumberOfTokensPerMsg should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint32(destGasOverhead), fqDestChainConfig.DestGasOverhead, "DestChainConfig DestGasOverhead should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint32(destDataAvailabilityOverheadGas), fqDestChainConfig.DestDataAvailabilityOverheadGas, "DestChainConfig DestDataAvailabilityOverheadGas should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint16(destGasPerDataAvailabilityByte), fqDestChainConfig.DestGasPerDataAvailabilityByte, "DestChainConfig DestGasPerDataAvailabilityByte should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint16(destDataAvailabilityMultiplierBps), fqDestChainConfig.DestDataAvailabilityMultiplierBps, "DestChainConfig DestDataAvailabilityMultiplierBps should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint32(maxDataBytes), fqDestChainConfig.MaxDataBytes, "DestChainConfig MaxDataBytes should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint32(maxPerMsgGasLimit), fqDestChainConfig.MaxPerMsgGasLimit, "DestChainConfig MaxPerMsgGasLimit should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint16(defaultTokenFeeUSDCents), fqDestChainConfig.DefaultTokenFeeUSDCents, "DestChainConfig DefaultTokenFeeUSDCents should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, uint32(defaultTokenDestGasOverhead), fqDestChainConfig.DefaultTokenDestGasOverhead, "DestChainConfig DefaultTokenDestGasOverhead should match for chain %d for %d", chain.Selector, otherChain.Selector)
			require.False(t, fqDestChainConfig.EnforceOutOfOrder, "DestChainConfig EnforceOutOfOrder should be false for chain %d for %d", chain.Selector, otherChain.Selector)

			// NonceManager has onRamp and offRamp set for other chains
			previousRamps, err := state.Chains[chain.Selector].NonceManager.GetPreviousRamps(callOpts, otherChain.Selector)
			require.NoError(t, err, "Failed to get previous ramps for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, state.Chains[chain.Selector].EVM2EVMOnRamp[otherChain.Selector].Address(), previousRamps.PrevOnRamp, "PrevOnRamp should match EVM2EVMOnRamp for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, state.Chains[chain.Selector].EVM2EVMOffRamp[otherChain.Selector].Address(), previousRamps.PrevOffRamp, "PrevOffRamp should match EVM2EVMOffRamp for chain %d for %d", chain.Selector, otherChain.Selector)

			// OnRamp has destChainConfig set for other chains
			onRampDestChainConfig, err := state.Chains[chain.Selector].OnRamp.GetDestChainConfig(callOpts, otherChain.Selector)
			require.NoError(t, err, "Failed to get dest chain config for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, state.Chains[chain.Selector].TestRouter.Address(), onRampDestChainConfig.Router, "DestChainConfig Router should match TestRouter for chain %d for %d", chain.Selector, otherChain.Selector)

			// OffRamp has sourceChainConfig set for other chains
			sourceChainConfig, err := state.Chains[chain.Selector].OffRamp.GetSourceChainConfig(callOpts, otherChain.Selector)
			require.NoError(t, err, "Failed to get source chain config for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, state.Chains[chain.Selector].TestRouter.Address(), sourceChainConfig.Router, "SourceChainConfig Router should match TestRouter for chain %d for %d", chain.Selector, otherChain.Selector)

			// OnRamp and OffRamp are connected to the TestRouter
			onRampOnRouter, err := state.Chains[chain.Selector].TestRouter.GetOnRamp(callOpts, otherChain.Selector)
			require.NoError(t, err, "Failed to get onRamp for chain %d for %d", chain.Selector, otherChain.Selector)
			require.Equal(t, state.Chains[chain.Selector].OnRamp.Address(), onRampOnRouter, "OnRamp on TestRouter should match OnRamp for chain %d for %d", chain.Selector, otherChain.Selector)
			isOffRamp, err := state.Chains[chain.Selector].TestRouter.IsOffRamp(callOpts, otherChain.Selector, state.Chains[chain.Selector].OffRamp.Address())
			require.NoError(t, err, "Failed to check if OffRamp is connected to TestRouter for chain %d for %d", chain.Selector, otherChain.Selector)
			require.True(t, isOffRamp, "OffRamp should be connected to TestRouter for chain %d for %d", chain.Selector, otherChain.Selector)
		}
	}

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset), v1_6.SetOCR3OffRampConfig{
			HomeChainSel:       homeChainSelector,
			RemoteChainSels:    e.BlockChains.ListChainSelectors(cldf_chain.WithFamily("evm")),
			CCIPHomeConfigType: globals.ConfigTypeCandidate,
		}),
	})
	require.NoError(t, err, "Failed to apply SetOCR3OffRampChangeset")

	// TODO: SetOCR3OffRampChangeset checks

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(v1_6.PromoteChainUpgradesChangeset, v1_6.PromoteChainUpgradesConfig{
			HomeChainSelector: homeChainSelector,
			ChainsToPromote:   e.BlockChains.ListChainSelectors(cldf_chain.WithFamily("evm")),
			MCMSConfig:        &mcmsCfg,
		}),
	})
	require.NoError(t, err, "Failed to apply PromoteChainUpgradesChangeset")

	// TODO: PromoteChainUpgradesChangeset checks
}
