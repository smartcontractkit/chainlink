package solana_test

import (
	"os"
	"testing"
	"time"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

func initialDeployCS(t *testing.T, e cldf.Environment, buildConfig *ccipChangesetSolana.BuildSolanaConfig) []commonchangeset.ConfiguredChangeSet {
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	mcmsConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()
	return []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken),
			commonchangeset.DeploySolanaLinkTokenConfig{
				ChainSelector: solChainSelectors[0],
				TokenPrivKey:  solLinkTokenPrivKey,
				TokenDecimals: 9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ChainSelector:     solChainSelectors[0],
				ContractParamsPerChain: ccipChangesetSolana.ChainContractParams{
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
				MCMSWithTimelockConfig: &mcmsConfig,
				BuildConfig:            buildConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployReceiverForTest),
			ccipChangesetSolana.DeployForTestConfig{
				ChainSelector: solChainSelectors[0],
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.ExtendGlobalLookupTableChangeset),
			ccipChangesetSolana.ExtendGlobalLookupTableConfig{
				ChainSelector:   solChainSelectors[0],
				LookupTableKeys: []solana.PublicKey{feeAggregatorPubKey, solLinkTokenPrivKey.PublicKey()}, // just add some random keys
			},
		),
	}
}

// use this for a quick deploy test
func TestDeployChainContractsChangesetPreload(t *testing.T) {
	quarantine.Flaky(t, "DX-1729")
	t.Parallel()

	homeChainSel := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	progIDs := soltestutils.LoadCCIPPrograms(t, programsPath)
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, progIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	err = testhelpers.SavePreloadedSolAddresses(*env, solSelector)
	require.NoError(t, err)

	e := *env

	// empty build config means, if artifacts are not present, resolve the artifact from github based on go.mod version
	// for a simple local in memory test, they will always be present, because we need them to spin up the in memory chain
	e, _, err = commonchangeset.ApplyChangesets(t, e, initialDeployCS(t, e, nil))
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}

// Upgrade flows must do the following:
// 1. Build the original contracts. We cannot preload because the deployed buffers will be too small to handle an upgrade.
// We must do a deploy with .so and keypairs locally
// 2. Build the upgraded contracts. We need the declare ids to match the existing deployed programs,
// so we need to do a local build again. We cannot do a remote fetch because those artifacts will not have the same keys as step 1.
// Doing this in CI is expensive, so we skip it for now.
func TestUpgrade(t *testing.T) {
	t.Parallel()
	skipInCI(t)

	homeChainSel := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	progIDs := soltestutils.LoadCCIPPrograms(t, programsPath)
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, progIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	e := *env

	e, _, err = commonchangeset.ApplyChangesets(t, e, initialDeployCS(t, e,
		&ccipChangesetSolana.BuildSolanaConfig{
			SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
			DestinationDir:        e.BlockChains.SolanaChains()[solSelector].ProgramsPath,
			LocalBuild: ccipChangesetSolana.LocalBuildConfig{
				BuildLocally:        true,
				CleanDestinationDir: true,
				GenerateVanityKeys:  true,
			},
		},
	))
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)

	feeAggregatorPrivKey2, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey2 := feeAggregatorPrivKey2.PublicKey()

	contractParamsPerChain := ccipChangesetSolana.ChainContractParams{
		FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
			DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
		},
		OffRampParams: ccipChangesetSolana.OffRampParams{
			EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
		},
	}

	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solSelector, true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})
	upgradeAuthority := timelockSignerPDA
	// upgradeAuthority := e.BlockChains.SolanaChains()[solChainSelectors[0]].DeployerKey.PublicKey()
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	addresses, err := e.ExistingAddresses.AddressesForChain(solSelector)
	require.NoError(t, err)
	chainState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[solSelector], addresses)
	require.NoError(t, err)

	// deploy the contracts
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		// upgrade authority
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
			ccipChangesetSolana.SetUpgradeAuthorityConfig{
				ChainSelector:         solSelector,
				NewUpgradeAuthority:   upgradeAuthority,
				SetAfterInitialDeploy: true,
				SetOffRamp:            true,
				SetMCMSPrograms:       true,
				TransferKeys: []solana.PublicKey{
					state.SolChains[solSelector].CCTPTokenPool,
				},
			},
		),
		// build the upgraded contracts and deploy/replace them onchain
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ChainSelector:          solSelector,
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewFeeQuoterVersion: &deployment.Version1_1_0,
					NewRouterVersion:    &deployment.Version1_1_0,
					// test offramp upgrade in place
					NewOffRampVersion:              &deployment.Version1_0_0,
					NewMCMVersion:                  &deployment.Version1_1_0,
					NewBurnMintTokenPoolVersion:    &deployment.Version1_1_0,
					NewLockReleaseTokenPoolVersion: &deployment.Version1_1_0,
					NewCCTPTokenPoolVersion:        &deployment.Version1_1_0,
					NewRMNRemoteVersion:            &deployment.Version1_1_0,
					NewAccessControllerVersion:     &deployment.Version1_1_0,
					NewTimelockVersion:             &deployment.Version1_1_0,
					UpgradeAuthority:               upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: 1 * time.Second,
					},
				},
				// build the contracts for upgrades
				BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
					SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
					DestinationDir:        e.BlockChains.SolanaChains()[solSelector].ProgramsPath,
					LocalBuild: ccipChangesetSolana.LocalBuildConfig{
						BuildLocally:        true,
						CleanDestinationDir: true,
						CleanGitDir:         true,
						UpgradeKeys: map[cldf.ContractType]string{
							shared.Router:                  state.SolChains[solSelector].Router.String(),
							shared.FeeQuoter:               state.SolChains[solSelector].FeeQuoter.String(),
							shared.BurnMintTokenPool:       state.SolChains[solSelector].BurnMintTokenPools[shared.CLLMetadata].String(),
							shared.LockReleaseTokenPool:    state.SolChains[solSelector].LockReleaseTokenPools[shared.CLLMetadata].String(),
							shared.OffRamp:                 state.SolChains[solSelector].OffRamp.String(),
							types.AccessControllerProgram:  chainState.AccessControllerProgram.String(),
							types.RBACTimelockProgram:      chainState.TimelockProgram.String(),
							types.ManyChainMultisigProgram: chainState.McmProgram.String(),
							shared.RMNRemote:               state.SolChains[solSelector].RMNRemote.String(),
							shared.CCTPTokenPool:           state.SolChains[solSelector].CCTPTokenPool.String(),
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solSelector,
				FeeAggregator: feeAggregatorPubKey2.String(),
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)
	state, err = stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	oldOffRampAddress := state.SolChains[solSelector].OffRamp
	// add a second offramp address
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ChainSelector:          solSelector,
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewOffRampVersion: &deployment.Version1_1_0,
					UpgradeAuthority:  upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: 1 * time.Second,
					},
				},
				BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
					SolanaContractVersion: ccipChangesetSolana.VersionSolanaV1_6_0,
					DestinationDir:        e.BlockChains.SolanaChains()[solSelector].ProgramsPath,
					LocalBuild: ccipChangesetSolana.LocalBuildConfig{
						BuildLocally: true,
					},
				},
			},
		),
	})
	require.NoError(t, err)
	// verify the offramp address is different
	state, err = stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	newOffRampAddress := state.SolChains[solSelector].OffRamp
	require.NotEqual(t, oldOffRampAddress, newOffRampAddress)

	// Verify router and fee quoter upgraded in place
	// and offramp had 2nd address added
	addresses, err = e.ExistingAddresses.AddressesForChain(solSelector)
	require.NoError(t, err)
	numRouters := 0
	numFeeQuoters := 0
	numOffRamps := 0
	for _, address := range addresses {
		if address.Type == shared.Router {
			numRouters++
		}
		if address.Type == shared.FeeQuoter {
			numFeeQuoters++
		}
		if address.Type == shared.OffRamp {
			numOffRamps++
		}
	}
	require.Equal(t, 1, numRouters)
	require.Equal(t, 1, numFeeQuoters)
	require.Equal(t, 2, numOffRamps)
	require.NoError(t, err)
	// solana verification
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	t.Parallel()
	skipInCI(t)

	homeChainSel := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	progIDs := soltestutils.LoadCCIPPrograms(t, programsPath)
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, progIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	e := *env

	e, _, err = commonchangeset.ApplyChangesets(t, e, initialDeployCS(t, e,
		&ccipChangesetSolana.BuildSolanaConfig{
			SolanaContractVersion: ccipChangesetSolana.VersionSolanaV1_6_0,
			DestinationDir:        e.BlockChains.SolanaChains()[solSelector].ProgramsPath,
			LocalBuild: ccipChangesetSolana.LocalBuildConfig{
				BuildLocally:        true,
				CleanDestinationDir: true,
				GenerateVanityKeys:  true,
			},
		},
	))
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solSelector, true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			OffRamp: true,
		})

	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	// test closing the old buffers
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.CloseBuffersChangeset),
			ccipChangesetSolana.CloseBuffersConfig{
				ChainSelector: solSelector,
				Programs: []string{
					state.SolChains[solSelector].BurnMintTokenPools[shared.CLLMetadata].String(),
					state.SolChains[solSelector].Router.String(),
				},
			},
		),
		// upgrade authority
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
			ccipChangesetSolana.SetUpgradeAuthorityConfig{
				ChainSelector:       solSelector,
				NewUpgradeAuthority: timelockSignerPDA,
				SetOffRamp:          true,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.CloseBuffersChangeset),
			ccipChangesetSolana.CloseBuffersConfig{
				ChainSelector: solSelector,
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
				Programs: []string{
					state.SolChains[solSelector].OffRamp.String(),
				},
			},
		),
	})
	require.NoError(t, err)
}

func TestIDL(t *testing.T) {
	skipInCI(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	e, _, err := commonchangeset.ApplyChangesets(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UploadIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector:         solChain,
				SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
				Router:                true,
				FeeQuoter:             true,
				OffRamp:               true,
				RMNRemote:             true,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				LockReleaseTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				CCTPTokenPool:    true,
				AccessController: true,
				Timelock:         true,
				MCM:              true,
			},
		),
	})
	require.NoError(t, err)

	// deploy timelock
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain, true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
		})

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetAuthorityIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				Router:        true,
				FeeQuoter:     true,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpgradeIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				Router:        true,
				FeeQuoter:     true,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})
	require.NoError(t, err)

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpgradeIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				OffRamp:       true,
				RMNRemote:     true,
				LockReleaseTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				CCTPTokenPool:    true,
				AccessController: true,
				Timelock:         true,
				MCM:              true,
			},
		),
	})
	require.NoError(t, err)

	// Test transferring ownership of the IDL back to the deployer key and then close and recreate the PDA
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetAuthorityIDLByMCMs),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})
	require.NoError(t, err)

	// close idl account
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.CloseIDLs),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
			},
		),
	})
	require.NoError(t, err)

	// deploy idl
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UploadIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
			},
		),
	})
	require.NoError(t, err)

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetAuthorityIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector:         solChain,
				SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpgradeIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector:         solChain,
				SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})

	// Close IDL
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.CloseIDLs),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})
	require.NoError(t, err)

	// Update IDL
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UploadIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
				BurnMintTokenPoolMetadata: []string{
					shared.CLLMetadata,
				},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
			},
		),
	})
	require.NoError(t, err)
}
