package v1_6_test

import (
	"math/big"
	"testing"
	"time"

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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
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
			}, 90000) // TODO: Staleness threshold should be configurable
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
					// TODO: These params should be configurable
					evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
						LinkToken:          state.Chains[sel].LinkToken.Address(),
						ChainSelector:      sel,
						DestChainSelector:  otherSel,
						DefaultTxGasLimit:  200_000,
						MaxNopFeesJuels:    big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
						PrevOnRamp:         utils.ZeroAddress,
						RmnProxy:           state.Chains[sel].RMNProxy.Address(),
						TokenAdminRegistry: state.Chains[sel].TokenAdminRegistry.Address(),
					},
					evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
						Router:                            state.Chains[sel].Router.Address(),
						MaxNumberOfTokensPerMsg:           5,
						DestGasOverhead:                   350_000,
						DestGasPerPayloadByte:             16,
						DestDataAvailabilityOverheadGas:   33_596,
						DestGasPerDataAvailabilityByte:    16,
						DestDataAvailabilityMultiplierBps: 6840,
						PriceRegistry:                     priceRegDeploy.Address,
						MaxDataBytes:                      1e5,
						MaxPerMsgGasLimit:                 4_000_000,
						DefaultTokenFeeUSDCents:           50,
						DefaultTokenDestGasOverhead:       32,
					},
					evm_2_evm_onramp.RateLimiterConfig{
						IsEnabled: false,
						Capacity:  big.NewInt(0),
						Rate:      big.NewInt(0),
					},
					[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
						{
							Token:                      state.Chains[sel].LinkToken.Address(),
							NetworkFeeUSDCents:         1_00,
							GasMultiplierWeiPerEth:     1e18,
							PremiumMultiplierWeiPerEth: 9e17,
							Enabled:                    true,
						},
						{
							Token:                      state.Chains[sel].Weth9.Address(),
							NetworkFeeUSDCents:         1_00,
							GasMultiplierWeiPerEth:     1e18,
							PremiumMultiplierWeiPerEth: 1e18,
							Enabled:                    true,
						},
					},
					[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
						{
							Token:                     state.Chains[sel].LinkToken.Address(),
							MinFeeUSDCents:            50,           // $0.5
							MaxFeeUSDCents:            1_000_000_00, // $ 1 million
							DeciBps:                   5_0,          // 5 bps
							DestGasOverhead:           110_000,
							DestBytesOverhead:         32,
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
					// TODO: These params should be configurable
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
					// TODO: These params should be configurable
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

	// TODO: InitChainUpgradesChangeset checks

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
