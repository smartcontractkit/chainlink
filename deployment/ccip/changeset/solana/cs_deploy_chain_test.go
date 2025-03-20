package solana_test

import (
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solBinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	OldSha = "e5f38e1c557eda4bc4a0436b69646a534ae16d39"
	NewSha = "ed22f75fe8de58634a46a137b386d4c508fde14f"
)

func TestDeployChainContractsChangesetSolana(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]ccipChangeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, ccipChangeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	feeAggregatorPrivKey2, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey2 := feeAggregatorPrivKey2.PublicKey()

	contractParamsPerChain := map[uint64]ccipChangesetSolana.ChainContractParams{
		solChainSelectors[0]: {
			FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
				DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
			},
			OffRampParams: ccipChangesetSolana.OffRampParams{
				EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
			},
		},
	}

	initialDeployConfig := ccipChangesetSolana.DeployChainContractsConfig{
		HomeChainSelector:      homeChainSel,
		ContractParamsPerChain: contractParamsPerChain,
	}

	ci := os.Getenv("CI") == "true"
	// we can't upgrade in place locally if we preload addresses so we have to change where we build
	// we also don't want to incur two builds in CI, so only do it locally
	if ci {
		testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	} else {
		initialDeployConfig.BuildConfig = ccipChangesetSolana.BuildSolanaConfig{
			GitCommitSha:        OldSha,
			DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
			CleanDestinationDir: true,
		}
		require.NoError(t, err)
	}

	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectorsSolana(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			initialDeployConfig,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployForTest),
			ccipChangesetSolana.DeployForTestConfig{
				ChainSelector: solChainSelectors[0],
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
	})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)

	// test upgrade flow
	if ci {
		// Expensive to run in CI
		return
	}

	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolana(t, &e, solChainSelectors[0], true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})
	upgradeAuthority := timelockSignerPDA
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	verifyProgramSizes(t, e)
	addresses, err := e.ExistingAddresses.AddressesForChain(e.AllChainSelectorsSolana()[0])
	require.NoError(t, err)
	chainState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[e.AllChainSelectorsSolana()[0]], addresses)
	require.NoError(t, err)

	// deploy the contracts
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		// upgrade authority
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
			ccipChangesetSolana.SetUpgradeAuthorityConfig{
				ChainSelector:         solChainSelectors[0],
				NewUpgradeAuthority:   upgradeAuthority,
				SetAfterInitialDeploy: true,
				SetMCMSPrograms:       true,
			},
		),
		// build the upgraded contracts and deploy/replace them onchain
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewFeeQuoterVersion: &deployment.Version1_1_0,
					NewRouterVersion:    &deployment.Version1_1_0,
					NewMCMVersion:       &deployment.Version1_1_0,
					UpgradeAuthority:    upgradeAuthority,
					SpillAddress:        upgradeAuthority,
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
				},
				// build the contracts for upgrades
				BuildConfig: ccipChangesetSolana.BuildSolanaConfig{
					GitCommitSha:        NewSha,
					DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
					CleanDestinationDir: true,
					CleanGitDir:         true,
					UpgradeKeys: map[deployment.ContractType]string{
						ccipChangeset.Router:               state.SolChains[solChainSelectors[0]].Router.String(),
						ccipChangeset.FeeQuoter:            state.SolChains[solChainSelectors[0]].FeeQuoter.String(),
						ccipChangeset.BurnMintTokenPool:    state.SolChains[solChainSelectors[0]].BurnMintTokenPool.String(),
						ccipChangeset.LockReleaseTokenPool: state.SolChains[solChainSelectors[0]].LockReleaseTokenPool.String(),
						types.AccessControllerProgram:      chainState.AccessControllerProgram.String(),
						types.RBACTimelockProgram:          chainState.TimelockProgram.String(),
						types.ManyChainMultisigProgram:     chainState.McmProgram.String(),
						ccipChangeset.RMNRemote:            state.SolChains[solChainSelectors[0]].RMNRemote.String(),
					},
				},
			},
		),
		// Split the upgrade to avoid txn size limits. No need to build again.
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewBurnMintTokenPoolVersion:    &deployment.Version1_1_0,
					NewLockReleaseTokenPoolVersion: &deployment.Version1_1_0,
					NewRMNRemoteVersion:            &deployment.Version1_1_0,
					UpgradeAuthority:               upgradeAuthority,
					SpillAddress:                   upgradeAuthority,
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
				},
			},
		),
		// Split the upgrade to avoid txn size limits. No need to build again.
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewAccessControllerVersion: &deployment.Version1_1_0,
					NewTimelockVersion:         &deployment.Version1_1_0,
					UpgradeAuthority:           upgradeAuthority,
					SpillAddress:               upgradeAuthority,
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey2.String(),
				MCMSSolana: &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				},
			},
		),
	})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
	state, err = ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	oldOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
	// add a second offramp address
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
							DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
						},
						OffRampParams: ccipChangesetSolana.OffRampParams{
							EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
						},
					},
				},
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewOffRampVersion: &deployment.Version1_1_0,
					UpgradeAuthority:  upgradeAuthority,
					SpillAddress:      upgradeAuthority,
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
				},
			},
		),
	})
	require.NoError(t, err)
	// verify the offramp address is different
	state, err = ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	newOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
	require.NotEqual(t, oldOffRampAddress, newOffRampAddress)

	// Verify router and fee quoter upgraded in place
	// and offramp had 2nd address added
	addresses, err = e.ExistingAddresses.AddressesForChain(solChainSelectors[0])
	require.NoError(t, err)
	numRouters := 0
	numFeeQuoters := 0
	numOffRamps := 0
	for _, address := range addresses {
		if address.Type == ccipChangeset.Router {
			numRouters++
		}
		if address.Type == ccipChangeset.FeeQuoter {
			numFeeQuoters++
		}
		if address.Type == ccipChangeset.OffRamp {
			numOffRamps++
		}
	}
	require.Equal(t, 1, numRouters)
	require.Equal(t, 1, numFeeQuoters)
	require.Equal(t, 2, numOffRamps)
	require.NoError(t, err)
	// solana verification
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}

func verifyProgramSizes(t *testing.T, e deployment.Environment) {
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	addresses, err := e.ExistingAddresses.AddressesForChain(e.AllChainSelectorsSolana()[0])
	require.NoError(t, err)
	chainState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[e.AllChainSelectorsSolana()[0]], addresses)
	require.NoError(t, err)
	programsToState := map[string]solana.PublicKey{
		deployment.RouterProgramName:               state.SolChains[e.AllChainSelectorsSolana()[0]].Router,
		deployment.OffRampProgramName:              state.SolChains[e.AllChainSelectorsSolana()[0]].OffRamp,
		deployment.FeeQuoterProgramName:            state.SolChains[e.AllChainSelectorsSolana()[0]].FeeQuoter,
		deployment.BurnMintTokenPoolProgramName:    state.SolChains[e.AllChainSelectorsSolana()[0]].BurnMintTokenPool,
		deployment.LockReleaseTokenPoolProgramName: state.SolChains[e.AllChainSelectorsSolana()[0]].LockReleaseTokenPool,
		deployment.AccessControllerProgramName:     chainState.AccessControllerProgram,
		deployment.TimelockProgramName:             chainState.TimelockProgram,
		deployment.McmProgramName:                  chainState.McmProgram,
		deployment.RMNRemoteProgramName:            state.SolChains[e.AllChainSelectorsSolana()[0]].RMNRemote,
	}
	for program, sizeBytes := range deployment.GetSolanaProgramBytes() {
		t.Logf("Verifying program %s size is at least %d bytes", program, sizeBytes)
		programDataAccount, _, _ := solana.FindProgramAddress([][]byte{programsToState[program].Bytes()}, solana.BPFLoaderUpgradeableProgramID)
		programDataSize, err := ccipChangesetSolana.SolProgramSize(&e, e.SolChains[e.AllChainSelectorsSolana()[0]], programDataAccount)
		require.NoError(t, err)
		require.GreaterOrEqual(t, programDataSize, sizeBytes)
	}
}

func TestDeployChainContractsChangesetLocal(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		return
	}
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]ccipChangeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, ccipChangeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectors(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectorsSolana(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			cfg,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangeset.DeployPrerequisitesChangeset),
			ccipChangeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
							DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
						},
						OffRampParams: ccipChangesetSolana.OffRampParams{
							EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
	})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}
