package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/don_id_claimer"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

func checkConnectivity(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	selector uint64,
	remoteChainSelector uint64,
	expectedRouter *router.Router,
	expectedAllowListEnabled bool,
	expectedRMNVerificationDisabled bool,
) {
	destChainConfig, err := state.Chains[selector].OnRamp.GetDestChainConfig(nil, remoteChainSelector)
	require.NoError(t, err, "must get dest chain config from onRamp")
	require.Equal(t, expectedRouter.Address().Hex(), destChainConfig.Router.Hex(), "router must equal expected")
	require.Equal(t, expectedAllowListEnabled, destChainConfig.AllowlistEnabled, "allowListEnabled must equal expected")

	srcChainConfig, err := state.Chains[selector].OffRamp.GetSourceChainConfig(nil, remoteChainSelector)
	require.NoError(t, err, "must get src chain config from offRamp")
	require.True(t, srcChainConfig.IsEnabled, "src chain config must be enabled")
	require.Equal(t, expectedRMNVerificationDisabled, srcChainConfig.IsRMNVerificationDisabled, "rmnVerificationDisabled must equal expected")
	require.Equal(t, common.LeftPadBytes(state.Chains[remoteChainSelector].OnRamp.Address().Bytes(), 32), srcChainConfig.OnRamp, "remote onRamp must be set on offRamp")
	require.Equal(t, expectedRouter.Address().Hex(), srcChainConfig.Router.Hex(), "router must equal expected")

	isOffRamp, err := expectedRouter.IsOffRamp(nil, remoteChainSelector, state.Chains[selector].OffRamp.Address())
	require.NoError(t, err, "must check if router has offRamp")
	require.True(t, isOffRamp, "router must have offRamp")
	onRamp, err := expectedRouter.GetOnRamp(nil, remoteChainSelector)
	require.NoError(t, err, "must get onRamp from router")
	require.Equal(t, state.Chains[selector].OnRamp.Address().Hex(), onRamp.Hex(), "onRamp must equal expected")
}

func TestConnectNewChain(t *testing.T) {
	t.Parallel()

	mustHaveOwner := func(t *testing.T, ownable commonchangeset.Ownable, expectedOwner string) {
		owner, err := ownable.Owner(nil)
		require.NoError(t, err, "must get owner")
		require.Equal(t, expectedOwner, owner.Hex(), "owner must be "+expectedOwner)
	}

	type test struct {
		Msg                        string
		TransferRemoteChainsToMCMS bool
		TestRouter                 bool
		MCMS                       *proposalutils.TimelockConfig
		ErrStr                     string
	}

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: mcmstypes.TimelockActionSchedule,
	}

	tests := []test{
		{
			Msg:                        "Use production router (with MCMS)",
			TransferRemoteChainsToMCMS: true,
			TestRouter:                 false,
			MCMS:                       mcmsConfig,
		},
		{
			Msg:                        "Use production router (without MCMS)",
			TransferRemoteChainsToMCMS: false,
			TestRouter:                 false,
			MCMS:                       nil,
		},
		{
			Msg:                        "Use test router (with MCMS)",
			TransferRemoteChainsToMCMS: true,
			TestRouter:                 true,
			MCMS:                       mcmsConfig,
		},
		{
			Msg:                        "Use test router (without MCMS)",
			TransferRemoteChainsToMCMS: false,
			TestRouter:                 true,
			MCMS:                       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.Chains = 3
			})
			e := deployedEnvironment.Env

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
			var newSelector uint64
			remoteChainSelectors := make([]uint64, 0, len(selectors)-1)
			for _, selector := range selectors {
				if selector != deployedEnvironment.HomeChainSel && newSelector == 0 {
					newSelector = selector // Just take any non-home chain selector
					continue
				}
				remoteChainSelectors = append(remoteChainSelectors, selector)
			}

			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(selectors))
			for _, selector := range selectors {
				// Assemble map of addresses required for Timelock scheduling & execution
				timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
					Timelock:  state.Chains[selector].Timelock,
					CallProxy: state.Chains[selector].CallProxy,
				}
			}

			if test.TransferRemoteChainsToMCMS {
				// onRamp, offRamp, and router on non-new chains are assumed to be owned by the timelock
				contractsToTransfer := make(map[uint64][]common.Address, len(remoteChainSelectors))
				for _, selector := range remoteChainSelectors {
					contractsToTransfer[selector] = []common.Address{
						state.Chains[selector].OnRamp.Address(),
						state.Chains[selector].OffRamp.Address(),
						state.Chains[selector].Router.Address(),
					}
				}
				e, err = commonchangeset.Apply(t, e, timelockContracts,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
						commoncs.TransferToMCMSWithTimelockConfig{
							ContractsByChain: contractsToTransfer,
							MCMSConfig: proposalutils.TimelockConfig{
								MinDelay: 0 * time.Second,
							},
						},
					),
				)
				require.NoError(t, err, "must apply TransferToMCMSWithTimelock")
			}

			remoteChains := make(map[uint64]v1_6.ConnectionConfig, len(remoteChainSelectors))
			for _, selector := range remoteChainSelectors {
				remoteChains[selector] = v1_6.ConnectionConfig{
					RMNVerificationDisabled: false,
					AllowListEnabled:        false,
				}
			}

			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.ConnectNewChainChangeset,
					v1_6.ConnectNewChainConfig{
						TestRouter:       &test.TestRouter,
						RemoteChains:     remoteChains,
						NewChainSelector: newSelector,
						NewChainConnectionConfig: v1_6.ConnectionConfig{
							RMNVerificationDisabled: true,
							AllowListEnabled:        true,
						},
						MCMSConfig: test.MCMS,
					},
				),
			)
			if test.ErrStr != "" {
				require.ErrorContains(t, err, test.ErrStr, "expected ConnectNewChainChangeset error")
				return
			}
			require.NoError(t, err, "must apply ConnectNewChainChangeset")

			for _, selector := range selectors {
				expectedAllowListEnabled := true
				expectedRMNVerificationDisabled := true
				remoteSelectors := []uint64{newSelector}
				if selector == newSelector {
					expectedAllowListEnabled = false
					expectedRMNVerificationDisabled = false
					remoteSelectors = remoteChainSelectors
					if !test.TestRouter && test.MCMS != nil {
						// New chain must have all contracts owned by timelock
						mustHaveOwner(t, state.Chains[selector].OnRamp, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].OffRamp, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].FeeQuoter, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].RMNProxy, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].NonceManager, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].TokenAdminRegistry, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].Router, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].RMNRemote, state.Chains[selector].Timelock.Address().Hex())

						// Admin role for deployer key should be revoked
						adminRole, err := state.Chains[selector].Timelock.ADMINROLE(nil)
						require.NoError(t, err, "must get admin role")
						hasRole, err := state.Chains[selector].Timelock.HasRole(nil, adminRole, e.BlockChains.EVMChains()[selector].DeployerKey.From)
						require.NoError(t, err, "must get admin role")
						require.False(t, hasRole, "deployer key must not have admin role")
					} else {
						// onRamp, offRamp, and router should still be owned by deployer key
						mustHaveOwner(t, state.Chains[selector].OnRamp, e.BlockChains.EVMChains()[selector].DeployerKey.From.Hex())
						mustHaveOwner(t, state.Chains[selector].OffRamp, e.BlockChains.EVMChains()[selector].DeployerKey.From.Hex())
						mustHaveOwner(t, state.Chains[selector].Router, e.BlockChains.EVMChains()[selector].DeployerKey.From.Hex())
					}
				}

				for _, remoteChainSelector := range remoteSelectors {
					expectedRouter := state.Chains[selector].Router
					if test.TestRouter {
						expectedRouter = state.Chains[selector].TestRouter
					}

					checkConnectivity(t, e, state, selector, remoteChainSelector, expectedRouter, expectedAllowListEnabled, expectedRMNVerificationDisabled)
				}
			}
		})
	}
}

func TestAddAndPromoteCandidatesForNewChain(t *testing.T) {
	t.Parallel()

	type test struct {
		Msg         string
		MCMS        *proposalutils.TimelockConfig
		DonIDOffSet *uint32
	}

	offset := uint32(0)

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: mcmstypes.TimelockActionSchedule,
	}

	testRouter := true

	tests := []test{
		{
			Msg:  "Remote chains owned by MCMS",
			MCMS: mcmsConfig,
		},
		{
			Msg:  "Remote chains not owned by MCMS",
			MCMS: nil,
		},
		{
			Msg:         "Remote chains with donID offset (Sync with capReg reg after wrong donIDClaim)",
			DonIDOffSet: &offset,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			chainIDs := []uint64{
				chain_selectors.ETHEREUM_MAINNET.EvmChainID,
				chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID,
				chain_selectors.ETHEREUM_MAINNET_OPTIMISM_1.EvmChainID,
			}

			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.ChainIDs = chainIDs
			})
			e := deployedEnvironment.Env
			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			// Identify and delete addresses from the new chain
			var newChainSelector uint64
			var linkAddress common.Address
			remoteChainSelectors := make([]uint64, 0, len(chainIDs)-1)
			addressesByChain := make(map[uint64]map[string]cldf.TypeAndVersion, len(chainIDs)-1)
			for _, selector := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
				if selector != deployedEnvironment.HomeChainSel && newChainSelector == 0 {
					newChainSelector = selector
					linkAddress = state.Chains[selector].LinkToken.Address()
				} else {
					remoteChainSelectors = append(remoteChainSelectors, selector)
					addrs, err := e.ExistingAddresses.AddressesForChain(selector)
					require.NoError(t, err, "must get addresses for chain")
					addressesByChain[selector] = addrs
				}
			}
			e.ExistingAddresses = cldf.NewMemoryAddressBookFromMap(addressesByChain)
			state, err = stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			// Identify and delete the DON ID for the new chain
			donID, err := internal.DonIDForChain(
				state.Chains[deployedEnvironment.HomeChainSel].CapabilityRegistry,
				state.Chains[deployedEnvironment.HomeChainSel].CCIPHome,
				newChainSelector,
			)
			require.NoError(t, err, "must get DON ID for chain")
			tx, err := state.Chains[deployedEnvironment.HomeChainSel].CapabilityRegistry.RemoveDONs(
				e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel].DeployerKey,
				[]uint32{donID},
			)
			require.NoError(t, err, "must remove DON ID")
			_, err = e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel].Confirm(tx)
			require.NoError(t, err, "must confirm DON ID removal")

			// Remove chain config on CCIPHome
			tx, err = state.Chains[deployedEnvironment.HomeChainSel].CCIPHome.ApplyChainConfigUpdates(
				e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel].DeployerKey,
				[]uint64{newChainSelector},
				[]ccip_home.CCIPHomeChainConfigArgs{},
			)
			require.NoError(t, err, "must remove chain config from CCIPHome")
			_, err = e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel].Confirm(tx)
			require.NoError(t, err, "must confirm chain config removal")

			// Assemble map of addresses required for Timelock scheduling & execution
			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(remoteChainSelectors))
			for _, selector := range remoteChainSelectors {
				timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
					Timelock:  state.Chains[selector].Timelock,
					CallProxy: state.Chains[selector].CallProxy,
				}
			}

			// Transfer remote contracts to MCMS if an MCMS config is supplied
			if test.MCMS != nil {
				contractsToTransfer := make(map[uint64][]common.Address, len(remoteChainSelectors))
				for _, selector := range remoteChainSelectors {
					contractsToTransfer[selector] = []common.Address{
						state.Chains[selector].OnRamp.Address(),
						state.Chains[selector].OffRamp.Address(),
						state.Chains[selector].Router.Address(),
						state.Chains[selector].FeeQuoter.Address(),
						state.Chains[selector].RMNProxy.Address(),
						state.Chains[selector].NonceManager.Address(),
						state.Chains[selector].TokenAdminRegistry.Address(),
						state.Chains[selector].RMNRemote.Address(),
					}
				}
				contractsToTransfer[deployedEnvironment.HomeChainSel] = append(
					contractsToTransfer[deployedEnvironment.HomeChainSel],
					state.Chains[deployedEnvironment.HomeChainSel].CCIPHome.Address(),
				)
				contractsToTransfer[deployedEnvironment.HomeChainSel] = append(
					contractsToTransfer[deployedEnvironment.HomeChainSel],
					state.Chains[deployedEnvironment.HomeChainSel].CapabilityRegistry.Address(),
				)
				e, err = commonchangeset.Apply(t, e, timelockContracts,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
						commoncs.TransferToMCMSWithTimelockConfig{
							ContractsByChain: contractsToTransfer,
							MCMSConfig: proposalutils.TimelockConfig{
								MinDelay: 0 * time.Second,
							},
						},
					),
				)
				require.NoError(t, err, "must apply TransferToMCMSWithTimelock")
			}

			// Build remote chain configurations
			remoteChains := make([]v1_6.ChainDefinition, len(remoteChainSelectors))
			for i, selector := range remoteChainSelectors {
				remoteChains[i] = v1_6.ChainDefinition{
					ConnectionConfig: v1_6.ConnectionConfig{
						RMNVerificationDisabled: true,
						AllowListEnabled:        false,
					},
					Selector:                 selector,
					GasPrice:                 big.NewInt(1e17),
					TokenPrices:              map[common.Address]*big.Int{},
					FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
				}
			}

			// Build new chain configuration
			nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
			require.NoError(t, err, "must get node info")
			mcmsDeploymentCfg := proposalutils.SingleGroupTimelockConfigV2(t)
			newChain := v1_6.NewChainDefinition{
				ChainDefinition: v1_6.ChainDefinition{
					ConnectionConfig: v1_6.ConnectionConfig{
						RMNVerificationDisabled: true,
						AllowListEnabled:        false,
					},
					Selector:                 newChainSelector,
					GasPrice:                 big.NewInt(1e17),
					TokenPrices:              map[common.Address]*big.Int{},
					FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
				},
				ChainContractParams: ccipseq.ChainContractParams{
					FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
					OffRampParams:   ccipops.DefaultOffRampParams(),
				},
				ExistingContracts: commoncs.ExistingContractsConfig{
					ExistingContracts: []commoncs.Contract{
						{
							Address:        linkAddress.Hex(),
							TypeAndVersion: cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
							ChainSelector:  newChainSelector,
						},
					},
				},
				ConfigOnHome: v1_6.ChainConfig{
					Readers: nodeInfo.NonBootstraps().PeerIDs(),
					FChain:  uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3), // #nosec G115 - Overflow is not a concern in this test scenario
					EncodableChainConfig: chainconfig.ChainConfig{
						GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
						DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
						OptimisticConfirmations: globals.OptimisticConfirmations,
					},
				},
				CommitOCRParams: v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, deployedEnvironment.FeedChainSel, nil, nil),
				ExecOCRParams:   v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
				// RMNRemoteConfig:   &v1_6.RMNRemoteConfig{...}, // TODO: Enable?
			}

			// deploy donIDClaimer
			e, err = commonchangeset.Apply(t, e, nil,
				commonchangeset.Configure(
					v1_6.DeployDonIDClaimerChangeset,
					v1_6.DeployDonIDClaimerConfig{},
				))
			require.NoError(t, err, "must deploy donIDClaimer contract")

			state, err = stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			if test.DonIDOffSet != nil {
				tx, err := state.Chains[deployedEnvironment.HomeChainSel].DonIDClaimer.ClaimNextDONId(e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel].DeployerKey)
				require.NoError(t, err)

				_, err = cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[deployedEnvironment.HomeChainSel], tx, don_id_claimer.DonIDClaimerABI, err)
				require.NoError(t, err)
			}

			// Apply AddCandidatesForNewChainChangeset
			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.AddCandidatesForNewChainChangeset,
					v1_6.AddCandidatesForNewChainConfig{
						HomeChainSelector:    deployedEnvironment.HomeChainSel,
						FeedChainSelector:    deployedEnvironment.FeedChainSel,
						NewChain:             newChain,
						RemoteChains:         remoteChains,
						MCMSDeploymentConfig: &mcmsDeploymentCfg,
						MCMSConfig:           test.MCMS,
						DonIDOffSet:          test.DonIDOffSet,
					},
				),
			)
			require.NoError(t, err, "must apply AddCandidatesForNewChainChangeset")
			state, err = stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			capReg := state.Chains[deployedEnvironment.HomeChainSel].CapabilityRegistry
			ccipHome := state.Chains[deployedEnvironment.HomeChainSel].CCIPHome
			rmnProxy := state.Chains[newChainSelector].RMNProxy
			rmnRemote := state.Chains[newChainSelector].RMNRemote
			feeQuoter := state.Chains[newChainSelector].FeeQuoter

			donID, err = internal.DonIDForChain(capReg, ccipHome, newChainSelector)
			require.NoError(t, err, "must get DON ID for chain")

			digests, err := ccipHome.GetConfigDigests(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
			candidateDigest := digests.CandidateConfigDigest
			require.NoError(t, err, "must get config digests")
			require.Empty(t, digests.ActiveConfigDigest, "active config digest must be empty")

			rmn, err := rmnProxy.GetARM(nil)
			require.NoError(t, err, "must get ARM")
			require.Equal(t, rmnRemote.Address(), rmn, "RMN must be set on RMNProxy")

			for _, remoteChain := range remoteChains {
				destChainConfig, err := feeQuoter.GetDestChainConfig(nil, remoteChain.Selector)
				require.NoError(t, err, "must get dest chain config from feeQuoter")
				require.Equal(t, remoteChain.FeeQuoterDestChainConfig, destChainConfig, "dest chain config must equal expected")

				gasPrice, err := feeQuoter.GetDestinationChainGasPrice(nil, remoteChain.Selector)
				require.NoError(t, err, "must get dest chain gas price from feeQuoter")
				require.Equal(t, remoteChain.GasPrice.String(), gasPrice.Value.String(), "gas price must equal expected")
			}

			// Apply PromoteNewChainForConfigChangeset
			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.PromoteNewChainForConfigChangeset,
					v1_6.PromoteNewChainForConfig{
						HomeChainSelector: deployedEnvironment.HomeChainSel,
						NewChain:          newChain,
						RemoteChains:      remoteChains,
						TestRouter:        &testRouter,
						MCMSConfig:        test.MCMS,
					},
				),
			)
			require.NoError(t, err, "must apply PromoteNewChainForConfigChangeset")

			digests, err = ccipHome.GetConfigDigests(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
			require.NoError(t, err, "must get config digests")
			require.Empty(t, digests.CandidateConfigDigest, "candidate config digest must be empty")
			require.Equal(t, candidateDigest, digests.ActiveConfigDigest, "active config digest must equal old candidate digest")

			testRouter := state.Chains[newChain.Selector].TestRouter
			for _, remoteChain := range remoteChains {
				feeQuoterOnRemote := state.Chains[remoteChain.Selector].FeeQuoter
				testRouterOnRemote := state.Chains[remoteChain.Selector].TestRouter

				destChainConfig, err := feeQuoterOnRemote.GetDestChainConfig(nil, newChain.Selector)
				require.NoError(t, err, "must get dest chain config from feeQuoter")
				require.Equal(t, newChain.FeeQuoterDestChainConfig, destChainConfig, "dest chain config must equal expected")

				gasPrice, err := feeQuoterOnRemote.GetDestinationChainGasPrice(nil, newChain.Selector)
				require.NoError(t, err, "must get dest chain gas price from feeQuoter")
				require.Equal(t, newChain.GasPrice.String(), gasPrice.Value.String(), "gas price must equal expected")

				checkConnectivity(t, e, state, remoteChain.Selector, newChain.Selector, testRouterOnRemote, false, true)
				checkConnectivity(t, e, state, newChain.Selector, remoteChain.Selector, testRouter, false, true)
			}

			// Apply ConnectNewChainChangeset to activate on prod routers
			noTestRouter := false
			remoteConnectionConfigs := make(map[uint64]v1_6.ConnectionConfig, len(remoteChainSelectors))
			for _, remoteChain := range remoteChains {
				remoteConnectionConfigs[remoteChain.Selector] = remoteChain.ConnectionConfig
			}
			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.ConnectNewChainChangeset,
					v1_6.ConnectNewChainConfig{
						NewChainSelector:         newChain.Selector,
						NewChainConnectionConfig: newChain.ConnectionConfig,
						RemoteChains:             remoteConnectionConfigs,
						TestRouter:               &noTestRouter,
						MCMSConfig:               test.MCMS,
					},
				),
			)
			require.NoError(t, err, "must apply ConnectNewChainChangeset")

			router := state.Chains[newChain.Selector].Router
			for _, remoteChain := range remoteChains {
				routerOnRemote := state.Chains[remoteChain.Selector].Router
				checkConnectivity(t, e, state, remoteChain.Selector, newChain.Selector, routerOnRemote, false, true)
				checkConnectivity(t, e, state, newChain.Selector, remoteChain.Selector, router, false, true)
			}
		})
	}
}
