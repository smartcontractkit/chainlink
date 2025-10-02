package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	migrate_seq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/migration"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	ccipocr3types "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	cciptypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
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
	maxNopFeesJuels    = big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18))
	newFeeQuoterParams = migrate_seq.NewFeeQuoterDestChainConfigParams{
		DestGasPerPayloadByteBase:      ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:      ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold: ccipevm.CalldataGasPerByteThreshold,
		DefaultTxGasLimit:              200_000,
		ChainFamilySelector:            [4]byte{0x28, 0x12, 0xd5, 0x2c},
		GasPriceStalenessThreshold:     0,
		GasMultiplierWeiPerEth:         11e17,
		NetworkFeeUSDCents:             10,
	}
)

func initMigrationEnvironment(t *testing.T, numChains int, mcmsCfg proposalutils.TimelockConfig) cldf_deploy.Environment {
	dEnv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(numChains),
		testhelpers.WithDONConfigurationSkipped(),
	)
	e := dEnv.Env
	chainSels := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily("evm"))

	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
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

		// Transfer TokenAdminRegistry, Router, & RMN Proxy to MCMS timelock on all chains
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2), commonchangeset.TransferToMCMSWithTimelockConfig{
				MCMSConfig: mcmsCfg,
				ContractsByChain: map[uint64][]common.Address{
					sel: {
						state.Chains[sel].TokenAdminRegistry.Address(),
						state.Chains[sel].RMNProxy.Address(),
						state.Chains[sel].Router.Address(),
					},
				},
			}),
		})
		if err != nil {
			t.Fatalf("Failed to transfer ownership of contracts to MCMS timelock: %v", err)
		}

		// Set LINK token on TokenAdminRegistry
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset), v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: shared.LinkSymbol,
				NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
					sel: {
						Type:               shared.BurnMintTokenPool,
						TokenAddress:       state.Chains[sel].LinkToken.Address(),
						LocalTokenDecimals: 18,
					},
				},
			}),
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset), v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: &mcmsCfg,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					sel: {
						shared.LinkSymbol: {
							Type:    shared.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			}),
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset), v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: &mcmsCfg,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					sel: {
						shared.LinkSymbol: {
							Type:    shared.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			}),
			commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset), v1_5_1.TokenAdminRegistryChangesetConfig{
				MCMS: &mcmsCfg,
				Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
					sel: {
						shared.LinkSymbol: {
							Type:    shared.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			}),
		})
		require.NoError(t, err, "Failed to set LINK token on TokenAdminRegistry")

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
							DestGasOverhead:           linkDestGasOverhead,
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

	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err, "Failed to load onchain state")

	homeChainSelector, err := state.HomeChainSelector()
	require.NoError(t, err, "Failed to get home chain selector")
	feedChainSelector := homeChainSelector // Just use home chain selector as feed chain selector for this test

	donCfgs := make(map[uint64]v1_6.DONConfig, len(e.BlockChains.EVMChains()))
	sourceChains := make(map[uint64]v1_6.SourceChainConfig)
	for _, chain := range e.BlockChains.EVMChains() {
		fqParams := make(map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams)
		for _, otherChain := range e.BlockChains.EVMChains() {
			if otherChain.Selector == chain.Selector {
				continue // Skip self
			}
			fqParams[otherChain.Selector] = newFeeQuoterParams
		}
		sourceChains[chain.Selector] = v1_6.SourceChainConfig{
			NewFeeQuoterParamsPerDest: fqParams,
		}
		nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
		require.NoError(t, err, "Failed to get node info for chain %d", chain.Selector)
		readers := nodeInfo.NonBootstraps().PeerIDs()
		donCfgs[chain.Selector] = v1_6.DONConfig{
			FeedChainSelector: feedChainSelector,
			CommitOCRParams:   v1_6.DefaultOCRParamsForCommitForETH,
			ExecOCRParams:     v1_6.DefaultOCRParamsForExecForETH,
			ChainConfig: v1_6.ChainConfig{
				Readers: readers,
				// #nosec G115 - Overflow is not a concern in this test scenario
				FChain: uint8(len(readers) / 3),
				EncodableChainConfig: chainconfig.ChainConfig{
					//nolint:staticcheck // SA1019: Type required by ChainConfig
					GasPriceDeviationPPB: ccipocr3types.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
					//nolint:staticcheck // SA1019: Type required by ChainConfig
					DAGasPriceDeviationPPB:    ccipocr3types.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
					OptimisticConfirmations:   globals.OptimisticConfirmations,
					ChainFeeDeviationDisabled: true,
				},
			},
		}
	}

	// Migrate each dest chain as its own batch
	// This will ensure that we can add DONs for a chain as a source and handle their reappearance as a dest gracefully (and vice versa).
	for i, destChain := range e.BlockChains.EVMChains() {
		destChainSel := destChain.Selector

		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(v1_6.InitChainUpgradesChangeset, v1_6.InitChainUpgradesConfig{
				HomeChainSelector: homeChainSelector,
				DONConfigs:        donCfgs,
				DestChains:        []uint64{destChainSel},
				SourceChains:      sourceChains,
				MCMSConfig:        &mcmsCfg,
			}),
		})
		require.NoError(t, err, "Failed to apply InitChainUpgradesChangeset")

		// Commit and exec candidates for all chains (source and dest) should be added when the first dest chain is migrated
		// Chain configs are added for all chains
		// FeeQuoter and NonceManager should be owned by the MCMS timelock on all chains
		commitCandidates := make(map[uint64][32]byte)
		execCandidates := make(map[uint64][32]byte)
		if i == 0 {
			for _, chain := range e.BlockChains.EVMChains() {
				// FeeQuoter is owned by the MCMS timelock on all chains
				fqOwner, err := state.Chains[chain.Selector].FeeQuoter.Owner(callOpts)
				require.NoError(t, err, "Failed to get FeeQuoter owner for chain %d", chain.Selector)
				require.Equal(t, state.Chains[chain.Selector].Timelock.Address(), fqOwner, "FeeQuoter owner should be MCMS timelock for chain %d", chain.Selector)

				// NonceManager is owned by the MCMS timelock on all chains
				nmOwner, err := state.Chains[chain.Selector].NonceManager.Owner(callOpts)
				require.NoError(t, err, "Failed to get NonceManager owner for chain %d", chain.Selector)
				require.Equal(t, state.Chains[chain.Selector].Timelock.Address(), nmOwner, "NonceManager owner should be MCMS timelock for chain %d", chain.Selector)

				// ChainConfig is set for the chain on CCIPHome
				chainConfig, err := state.Chains[homeChainSelector].CCIPHome.GetChainConfig(callOpts, chain.Selector)
				require.NoError(t, err, "Failed to get chain config for chain %d", chain.Selector)
				require.Equal(t, chainConfig.Readers, donCfgs[chain.Selector].ChainConfig.Readers, "ChainConfig readers should match for chain %d", chain.Selector)
				require.Equal(t, chainConfig.FChain, donCfgs[chain.Selector].ChainConfig.FChain, "ChainConfig FChain should match for chain %d", chain.Selector)

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

				commitCandidates[chain.Selector] = commitCandidate
				execCandidates[chain.Selector] = execCandidate
			}
		}

		// RMNRemote is owned by the MCMS timelock on dest
		owner, err := state.Chains[destChainSel].RMNRemote.Owner(callOpts)
		require.NoError(t, err, "Failed to get RMNRemote owner for chain %d", destChainSel)
		require.Equal(t, state.Chains[destChainSel].Timelock.Address(), owner, "RMNRemote owner should be MCMS timelock for chain %d", destChainSel)

		// RMNProxy is pointing at the RMNRemote on dest
		rmnOnProxyAddr, err := state.Chains[destChainSel].RMNProxy.GetARM(callOpts)
		require.NoError(t, err, "Failed to get RMNProxy ARM for chain %d", destChainSel)
		require.Equal(t, state.Chains[destChainSel].RMNRemote.Address(), rmnOnProxyAddr, "RMNProxy should point to RMNRemote for chain %d", destChainSel)

		for _, sourceChain := range e.BlockChains.EVMChains() {
			if sourceChain.Selector == destChainSel {
				continue
			}

			// FeeAggregator on OnRamp is set to Timelock on source OnRamp
			onRampDynamicConfig, err := state.Chains[sourceChain.Selector].OnRamp.GetDynamicConfig(callOpts)
			require.NoError(t, err, "Failed to get OnRamp dynamic config for chain %d", sourceChain.Selector)
			require.Equal(t, state.Chains[sourceChain.Selector].Timelock.Address(), onRampDynamicConfig.FeeAggregator, "OnRamp FeeAggregator should be MCMS timelock for chain %d", sourceChain.Selector)

			// PremiumMultiplierWeiPerEth is set for WETH and LINK
			linkPremium, err := state.Chains[sourceChain.Selector].FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, state.Chains[sourceChain.Selector].LinkToken.Address())
			require.NoError(t, err, "Failed to get PremiumMultiplierWeiPerEth for LINK on chain %d", sourceChain.Selector)
			require.Equal(t, uint64(linkPremiumMultiplierWeiPerEth), linkPremium, "LINK PremiumMultiplierWeiPerEth should match for chain %d", sourceChain.Selector)
			wethPremium, err := state.Chains[sourceChain.Selector].FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, state.Chains[sourceChain.Selector].Weth9.Address())
			require.NoError(t, err, "Failed to get PremiumMultiplierWeiPerEth for WETH on chain %d", sourceChain.Selector)
			require.Equal(t, uint64(wethPremiumMultiplierWeiPerEth), wethPremium, "WETH PremiumMultiplierWeiPerEth should match for chain %d", sourceChain.Selector)

			// TransferFeeConfigArgs are set on the FeeQuoter for LINK
			transferFeeConfigArgs, err := state.Chains[sourceChain.Selector].FeeQuoter.GetTokenTransferFeeConfig(callOpts, destChainSel, state.Chains[sourceChain.Selector].LinkToken.Address())
			require.NoError(t, err, "Failed to get LINK transfer fee config on chain %d for %d", sourceChain.Selector, destChainSel)

			require.Equal(t, uint32(linkMinFeeUSDCents), transferFeeConfigArgs.MinFeeUSDCents, "LINK MinFeeUSDCents should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(linkMaxFeeUSDCents), transferFeeConfigArgs.MaxFeeUSDCents, "LINK MaxFeeUSDCents should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint16(linkDeciBps), transferFeeConfigArgs.DeciBps, "LINK DeciBps should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(linkDestGasOverhead), transferFeeConfigArgs.DestGasOverhead, "LINK DestGasOverhead should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(linkDestBytesOverhead), transferFeeConfigArgs.DestBytesOverhead, "LINK DestBytesOverhead should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.True(t, transferFeeConfigArgs.IsEnabled, "LINK Transfer fee config should be enabled on chain %d for %d", sourceChain.Selector, destChainSel)

			// Fee tokens are set on the fee quoter
			feeTokens, err := state.Chains[sourceChain.Selector].FeeQuoter.GetFeeTokens(callOpts)
			require.NoError(t, err, "Failed to get fee tokens on chain %d", sourceChain.Selector)
			require.Len(t, feeTokens, 2, "Expected 2 fee tokens on chain %d", sourceChain.Selector)
			require.Contains(t, feeTokens, state.Chains[sourceChain.Selector].LinkToken.Address(), "Fee tokens should contain LINK on chain %d", sourceChain.Selector)
			require.Contains(t, feeTokens, state.Chains[sourceChain.Selector].Weth9.Address(), "Fee tokens should contain WETH on chain %d", sourceChain.Selector)

			// DestChainConfig is set for other chains
			fqDestChainConfig, err := state.Chains[sourceChain.Selector].FeeQuoter.GetDestChainConfig(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get dest chain config for chain %d for %d", sourceChain.Selector, destChainSel)
			require.True(t, fqDestChainConfig.IsEnabled, "DestChainConfig should be enabled on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint16(maxNumberOfTokensPerMsg), fqDestChainConfig.MaxNumberOfTokensPerMsg, "DestChainConfig MaxNumberOfTokensPerMsg should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(destGasOverhead), fqDestChainConfig.DestGasOverhead, "DestChainConfig DestGasOverhead should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(destDataAvailabilityOverheadGas), fqDestChainConfig.DestDataAvailabilityOverheadGas, "DestChainConfig DestDataAvailabilityOverheadGas should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint16(destGasPerDataAvailabilityByte), fqDestChainConfig.DestGasPerDataAvailabilityByte, "DestChainConfig DestGasPerDataAvailabilityByte should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint16(destDataAvailabilityMultiplierBps), fqDestChainConfig.DestDataAvailabilityMultiplierBps, "DestChainConfig DestDataAvailabilityMultiplierBps should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(maxDataBytes), fqDestChainConfig.MaxDataBytes, "DestChainConfig MaxDataBytes should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(maxPerMsgGasLimit), fqDestChainConfig.MaxPerMsgGasLimit, "DestChainConfig MaxPerMsgGasLimit should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint16(defaultTokenFeeUSDCents), fqDestChainConfig.DefaultTokenFeeUSDCents, "DestChainConfig DefaultTokenFeeUSDCents should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, uint32(defaultTokenDestGasOverhead), fqDestChainConfig.DefaultTokenDestGasOverhead, "DestChainConfig DefaultTokenDestGasOverhead should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.False(t, fqDestChainConfig.EnforceOutOfOrder, "DestChainConfig EnforceOutOfOrder should be false on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.ChainFamilySelector, fqDestChainConfig.ChainFamilySelector, "DestChainConfig ChainFamilySelector should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.DestGasPerPayloadByteBase, fqDestChainConfig.DestGasPerPayloadByteBase, "DestChainConfig DestGasPerPayloadByteBase should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.DestGasPerPayloadByteHigh, fqDestChainConfig.DestGasPerPayloadByteHigh, "DestChainConfig DestGasPerPayloadByteHigh should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.DestGasPerPayloadByteThreshold, fqDestChainConfig.DestGasPerPayloadByteThreshold, "DestChainConfig DestGasPerPayloadByteThreshold should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.DefaultTxGasLimit, fqDestChainConfig.DefaultTxGasLimit, "DestChainConfig DefaultTxGasLimit should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.GasPriceStalenessThreshold, fqDestChainConfig.GasPriceStalenessThreshold, "DestChainConfig GasPriceStalenessThreshold should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.NetworkFeeUSDCents, fqDestChainConfig.NetworkFeeUSDCents, "DestChainConfig NetworkFeeUSDCents should match on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, newFeeQuoterParams.GasMultiplierWeiPerEth, fqDestChainConfig.GasMultiplierWeiPerEth, "DestChainConfig GasMultiplierWeiPerEth should match on chain %d for %d", sourceChain.Selector, destChainSel)

			// NonceManager has onRamp and offRamp set for other chains
			previousRamps, err := state.Chains[sourceChain.Selector].NonceManager.GetPreviousRamps(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get previous ramps on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, state.Chains[sourceChain.Selector].EVM2EVMOnRamp[destChainSel].Address(), previousRamps.PrevOnRamp, "PrevOnRamp should match EVM2EVMOnRamp on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, state.Chains[sourceChain.Selector].EVM2EVMOffRamp[destChainSel].Address(), previousRamps.PrevOffRamp, "PrevOffRamp should match EVM2EVMOffRamp on chain %d for %d", sourceChain.Selector, destChainSel)
			previousRamps, err = state.Chains[destChainSel].NonceManager.GetPreviousRamps(callOpts, sourceChain.Selector)
			require.NoError(t, err, "Failed to get previous ramps on chain %d for %d", destChainSel, sourceChain.Selector)
			require.Equal(t, state.Chains[destChainSel].EVM2EVMOnRamp[sourceChain.Selector].Address(), previousRamps.PrevOnRamp, "PrevOnRamp should match EVM2EVMOnRamp on chain %d for %d", destChainSel, sourceChain.Selector)
			require.Equal(t, state.Chains[destChainSel].EVM2EVMOffRamp[sourceChain.Selector].Address(), previousRamps.PrevOffRamp, "PrevOffRamp should match EVM2EVMOffRamp on chain %d for %d", destChainSel, sourceChain.Selector)

			// OnRamp has destChainConfig set for other chains
			onRampDestChainConfig, err := state.Chains[sourceChain.Selector].OnRamp.GetDestChainConfig(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get dest chain config on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, state.Chains[sourceChain.Selector].TestRouter.Address(), onRampDestChainConfig.Router, "DestChainConfig Router should match TestRouter on chain %d for %d", sourceChain.Selector, destChainSel)

			// OffRamp has sourceChainConfig set for other chains
			sourceChainConfig, err := state.Chains[destChainSel].OffRamp.GetSourceChainConfig(callOpts, sourceChain.Selector)
			require.NoError(t, err, "Failed to get source chain config on chain %d for %d", destChainSel, sourceChain.Selector)
			require.Equal(t, state.Chains[destChainSel].TestRouter.Address(), sourceChainConfig.Router, "SourceChainConfig Router should match TestRouter on chain %d for %d", destChainSel, sourceChain.Selector)

			// OnRamp and OffRamp are connected to the TestRouter
			onRampOnRouter, err := state.Chains[sourceChain.Selector].TestRouter.GetOnRamp(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get onRamp on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, state.Chains[sourceChain.Selector].OnRamp.Address(), onRampOnRouter, "OnRamp on TestRouter should match OnRamp on chain %d for %d", sourceChain.Selector, destChainSel)
			isOffRamp, err := state.Chains[destChainSel].TestRouter.IsOffRamp(callOpts, sourceChain.Selector, state.Chains[destChainSel].OffRamp.Address())
			require.NoError(t, err, "Failed to check if OffRamp is connected to TestRouter on chain %d for %d", destChainSel, sourceChain.Selector)
			require.True(t, isOffRamp, "OffRamp should be connected to TestRouter on chain %d for %d", destChainSel, sourceChain.Selector)
		}

		// SetOCR3OffRampChangeset is not idempotent, so we only apply it on the first iteration
		if i == 0 {
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(cldf_deploy.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset), v1_6.SetOCR3OffRampConfig{
					HomeChainSel:       homeChainSelector,
					RemoteChainSels:    e.BlockChains.ListChainSelectors(cldf_chain.WithFamily("evm")),
					CCIPHomeConfigType: globals.ConfigTypeCandidate,
				}),
			})
			require.NoError(t, err, "Failed to apply SetOCR3OffRampChangeset")

			// SetOCR3OffRampChangeset checks
			for _, chain := range e.BlockChains.EVMChains() {
				commitCfg, err := state.Chains[chain.Selector].OffRamp.LatestConfigDetails(callOpts, uint8(cciptypes.PluginTypeCCIPCommit))
				require.NoError(t, err, "Failed to get latest commit config for chain %d", chain.Selector)
				require.Equal(t, commitCandidates[chain.Selector], commitCfg.ConfigInfo.ConfigDigest, "Commit candidate should match for chain %d", chain.Selector)
				execCfg, err := state.Chains[chain.Selector].OffRamp.LatestConfigDetails(callOpts, uint8(cciptypes.PluginTypeCCIPExec))
				require.NoError(t, err, "Failed to get latest exec config for chain %d", chain.Selector)
				require.Equal(t, execCandidates[chain.Selector], execCfg.ConfigInfo.ConfigDigest, "Exec candidate should match for chain %d", chain.Selector)
			}
		}

		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(v1_6.PromoteChainUpgradesChangeset, v1_6.PromoteChainUpgradesConfig{
				HomeChainSelector: homeChainSelector,
				DestChains:        []uint64{destChainSel},
				MCMSConfig:        &mcmsCfg,
			}),
		})
		require.NoError(t, err, "Failed to apply PromoteChainUpgradesChangeset")

		// OffRamp is owned by the MCMS timelock
		owner, err = state.Chains[destChainSel].OffRamp.Owner(callOpts)
		require.NoError(t, err, "Failed to get OffRamp owner for chain %d", destChainSel)
		require.Equal(t, state.Chains[destChainSel].Timelock.Address(), owner, "OffRamp owner should be MCMS timelock for chain %d", destChainSel)

		// PromoteChainUpgradesChangeset checks
		for _, sourceChain := range e.BlockChains.EVMChains() {
			if sourceChain.Selector == destChainSel {
				continue // Skip self
			}

			// OnRamp on other chains is owned by the MCMS timelock
			onRampOwner, err := state.Chains[sourceChain.Selector].OnRamp.Owner(callOpts)
			require.NoError(t, err, "Failed to get OnRamp owner for chain %d", sourceChain.Selector)
			require.Equal(t, state.Chains[sourceChain.Selector].Timelock.Address(), onRampOwner, "OnRamp owner should be MCMS timelock for chain %d", sourceChain.Selector)

			// OnRamp has destChainConfig set for other chains
			onRampDestChainConfig, err := state.Chains[sourceChain.Selector].OnRamp.GetDestChainConfig(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get dest chain config for OnRamp on chain %d for %d", sourceChain.Selector, destChainSel)
			require.Equal(t, state.Chains[sourceChain.Selector].Router.Address(), onRampDestChainConfig.Router, "DestChainConfig Router should match Router")

			// OffRamp has sourceChainConfig set for other chains
			sourceChainConfig, err := state.Chains[destChainSel].OffRamp.GetSourceChainConfig(callOpts, sourceChain.Selector)
			require.NoError(t, err, "Failed to get source chain config for OffRamp on chain %d for %d", destChainSel, sourceChain.Selector)
			require.Equal(t, state.Chains[destChainSel].Router.Address(), sourceChainConfig.Router, "SourceChainConfig Router should match Router")

			// OnRamp and OffRamp are connected to the MainRouter
			onRampOnRouter, err := state.Chains[sourceChain.Selector].Router.GetOnRamp(callOpts, destChainSel)
			require.NoError(t, err, "Failed to get onRamp for chain %d on %d", destChainSel, sourceChain.Selector)
			require.Equal(t, state.Chains[sourceChain.Selector].OnRamp.Address(), onRampOnRouter, "OnRamp on Router should match OnRamp")
			isOffRamp, err := state.Chains[destChainSel].Router.IsOffRamp(callOpts, sourceChain.Selector, state.Chains[destChainSel].OffRamp.Address())
			require.NoError(t, err, "Failed to check if OffRamp is connected to Router on chain %d for %d", destChainSel, sourceChain.Selector)
			require.True(t, isOffRamp, "OffRamp should be connected to Router on chain %d for %d", destChainSel, sourceChain.Selector)
		}
	}
}
