package solana_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solBinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	solRpc "github.com/gagliardetto/solana-go/rpc"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

// For remote fetching, we need to use the short sha
const (
	OldSha = "aa0756b72e7b70640a6a6235fbbd13aff407402a"
	NewSha = "f1ced171b7538afc6c9f488803f90d10ac0f0b52"
)

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
		programDataSize, err := ccipChangesetSolana.GetSolProgramSize(&e, e.SolChains[e.AllChainSelectorsSolana()[0]], programDataAccount)
		require.NoError(t, err)
		require.GreaterOrEqual(t, programDataSize, sizeBytes)
	}
}

func initialDeployCS(t *testing.T, e deployment.Environment, buildConfig *ccipChangesetSolana.BuildSolanaConfig) []commonchangeset.ConfiguredChangeSet {
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	mcmsConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	return []commonchangeset.ConfiguredChangeSet{
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
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployReceiverForTest),
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
	}
}

// use this for a quick deploy test
func TestDeployChainContractsChangesetPreload(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	solChainSelectors := e.AllChainSelectorsSolana()
	err := testhelpers.SavePreloadedSolAddresses(e, solChainSelectors[0])
	require.NoError(t, err)
	// empty build config means, if artifacts are not present, resolve the artifact from github based on go.mod version
	// for a simple local in memory test, they will always be present, because we need them to spin up the in memory chain
	e, err = commonchangeset.ApplyChangesetsV2(t, e, initialDeployCS(t, e, nil))
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}

// Upgrade flows must do the following:
// 1. Build the original contracts. We cannot preload because the deployed buffers will be too small to handle an upgrade.
// We must do a deploy with .so and keypairs locally
// 2. Build the upgraded contracts. We need the declare ids to match the existing deployed programs,
// so we need to do a local build again. We cannot do a remote fetch because those artifacts will not have the same keys as step 1.
// Doing this in CI is expensive, so we skip it for now.
func TestUpgrade(t *testing.T) {
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
	e, err := commonchangeset.ApplyChangesetsV2(t, e, initialDeployCS(t, e,
		&ccipChangesetSolana.BuildSolanaConfig{
			GitCommitSha:   OldSha,
			DestinationDir: e.SolChains[solChainSelectors[0]].ProgramsPath,
			LocalBuild: ccipChangesetSolana.LocalBuildConfig{
				BuildLocally:        true,
				CleanDestinationDir: true,
				GenerateVanityKeys:  true,
			},
		},
	))
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)

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

	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolana(t, &e, solChainSelectors[0], true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})
	upgradeAuthority := timelockSignerPDA
	// upgradeAuthority := e.SolChains[solChainSelectors[0]].DeployerKey.PublicKey()
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
				ChainSelector:          solChainSelectors[0],
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewFeeQuoterVersion: &deployment.Version1_1_0,
					NewRouterVersion:    &deployment.Version1_1_0,
					NewMCMVersion:       &deployment.Version1_1_0,
					UpgradeAuthority:    upgradeAuthority,
					SpillAddress:        upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: 1 * time.Second,
					},
				},
				// build the contracts for upgrades
				BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
					GitCommitSha:   NewSha,
					DestinationDir: e.SolChains[solChainSelectors[0]].ProgramsPath,
					LocalBuild: ccipChangesetSolana.LocalBuildConfig{
						BuildLocally:        true,
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
			},
		),
		// Split the upgrade to avoid txn size limits. No need to build again.
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ChainSelector:          solChainSelectors[0],
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewBurnMintTokenPoolVersion:    &deployment.Version1_1_0,
					NewLockReleaseTokenPoolVersion: &deployment.Version1_1_0,
					NewRMNRemoteVersion:            &deployment.Version1_1_0,
					UpgradeAuthority:               upgradeAuthority,
					SpillAddress:                   upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
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
				ChainSelector:          solChainSelectors[0],
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewAccessControllerVersion: &deployment.Version1_1_0,
					NewTimelockVersion:         &deployment.Version1_1_0,
					UpgradeAuthority:           upgradeAuthority,
					SpillAddress:               upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
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
					MCMS: &proposalutils.TimelockConfig{
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
				HomeChainSelector:      homeChainSel,
				ChainSelector:          solChainSelectors[0],
				ContractParamsPerChain: contractParamsPerChain,
				UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
					NewOffRampVersion: &deployment.Version1_1_0,
					UpgradeAuthority:  upgradeAuthority,
					SpillAddress:      upgradeAuthority,
					MCMS: &proposalutils.TimelockConfig{
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

func TestIDLUpgrade(t *testing.T) {
	// programID := solana.MustPublicKeyFromBase58("AacpQtBFpfVDWqacCqBPj59GChahCYUNLNb4Wsvft83M")
	// idlAddress, _, _ := solana.FindProgramAddress([][]byte{[]byte("anchor:idl")}, programID)
	// solana.CreateWithSeed(base, )
	// fmt.Println(idlAddress.String())

	programID := solana.MustPublicKeyFromBase58("AacpQtBFpfVDWqacCqBPj59GChahCYUNLNb4Wsvft83M")

	deployerKey, err := solana.PrivateKeyFromSolanaKeygenFile("/Users/yashvardhan/.config/solana/id_devnet.json")
	require.NoError(t, err)
	currentUpgradeAuthority := solana.MustPublicKeyFromBase58("7oZnxiocDK1aa9XAQC3CZ1VHKFkKwLuwRK8NddhU3FT2")
	require.Equal(t, deployerKey.PublicKey().String(), currentUpgradeAuthority.String())

	// Step 1: Get the base PDA using FindProgramAddress with no seeds
	base, _, err := solana.FindProgramAddress([][]byte{}, programID)
	if err != nil {
		fmt.Printf("failed to find base PDA: %v", err)
	}

	// Step 2: Create the derived address with a seed
	seed := "anchor:idl"
	idlAddress, err := solana.CreateWithSeed(base, seed, programID)
	if err != nil {
		fmt.Printf("failed to create with seed: %v", err)
	}

	fmt.Println("Base PDA:     ", base.String())
	fmt.Println("IDL Address:  ", idlAddress.String())

	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(idlAddress, true, false),
		solana.NewAccountMeta(currentUpgradeAuthority, false, true), // New upgrade authority
	}

	newAuthority := solana.MustPublicKeyFromBase58("FxghvBLeWky3gxXYnDP2sHa2MuFbJ7WWiSzFT96sMZqi")

	// Step 1: Build IDL instruction data
	data, err := serializeIdlSetAuthority(newAuthority)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Data:         ", hex.EncodeToString(data))

	instruction := solana.NewInstruction(
		programID,
		keys,
		data,
	)
	client := solRpc.New(solRpc.DevNet.RPC)
	_, err = solCommonUtil.SendAndConfirm(
		context.Background(), client, []solana.Instruction{instruction}, deployerKey, solRpc.CommitmentConfirmed,
	)
	require.NoError(t, err)
}

// Sha256("anchor:idl")[..8] = 0x0a69e9a778bcf440
const IDL_IX_TAG uint64 = 0x0a69e9a778bcf440

func serializeIdlSetAuthority(newAuthority solana.PublicKey) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Step 1: Write IDL_IX_TAG as little endian
	if err := binary.Write(buf, binary.LittleEndian, IDL_IX_TAG); err != nil {
		return nil, fmt.Errorf("failed to write IDL_IX_TAG: %w", err)
	}

	// Step 2: Write 8-byte discriminator for SetAuthority
	discriminator := idlDiscriminator("SetAuthority")
	buf.Write(discriminator[:])

	// Step 3: Write new authority pubkey (32 bytes)
	buf.Write(newAuthority.Bytes())

	return buf.Bytes(), nil
}

// Generate the 8-byte Anchor discriminator for an IDL instruction
func idlDiscriminator(name string) [8]byte {
	hash := sha256.Sum256([]byte("global:" + name))
	var out [8]byte
	copy(out[:], hash[:8])
	return out
}
