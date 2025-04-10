package solana_test

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
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

	"github.com/gagliardetto/solana-go/programs/system"
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
	programID := solana.MustPublicKeyFromBase58("AacpQtBFpfVDWqacCqBPj59GChahCYUNLNb4Wsvft83M")

	deployerKey, err := solana.PrivateKeyFromSolanaKeygenFile("/Users/yashvardhan/.config/solana/id_devnet.json")
	require.NoError(t, err)
	// currentUpgradeAuthority := solana.MustPublicKeyFromBase58("7oZnxiocDK1aa9XAQC3CZ1VHKFkKwLuwRK8NddhU3FT2")
	// require.Equal(t, deployerKey.PublicKey().String(), currentUpgradeAuthority.String())

	// derive idl address
	base, _, err := solana.FindProgramAddress([][]byte{}, programID)
	require.NoError(t, err)
	idlAddress, err := solana.CreateWithSeed(base, "anchor:idl", programID)
	require.NoError(t, err)
	fmt.Println("IDL Address:  ", idlAddress.String())

	// build set authority instruction data
	newAuthority := solana.MustPublicKeyFromBase58("FxghvBLeWky3gxXYnDP2sHa2MuFbJ7WWiSzFT96sMZqi")
	data := binary.LittleEndian.AppendUint64([]byte{}, IDL_IX_TAG) // 4-byte Extend instruction identifier
	data = append(data, byte(4))
	data = append(data, newAuthority.Bytes()...)
	fmt.Println("Data:         ", hex.EncodeToString(data))

	instruction := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(idlAddress, true, false),
			solana.NewAccountMeta(deployerKey.PublicKey(), false, true), // Current upgrade authority (signer)
		},
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

func getCompressedIDL(t *testing.T) ([]byte, error) {
	idlBytes, err := os.ReadFile("/Users/yashvardhan/chainlink/deployment/ccip/changeset/internal/solana_contracts/access_controller.json")
	require.NoError(t, err)
	var idl ccipChangesetSolana.IDL
	err = json.Unmarshal(idlBytes, &idl)
	require.NoError(t, err)
	compressedIDL, err := serializeIdl(idl)
	require.NoError(t, err)
	return compressedIDL, nil
}

func TestIDLSetBuffer(t *testing.T) {
	// These are placeholder IDs — replace with your actual ones
	programID := solana.MustPublicKeyFromBase58("GnpXUEvp4Uu5qzTRDHzDT5oRh7oY7g5GDRagc6zJo7M5")
	payer, err := solana.PrivateKeyFromSolanaKeygenFile("/Users/yashvardhan/.config/solana/id.json")
	require.NoError(t, err)

	compressedIDL, err := getCompressedIDL(t)
	require.NoError(t, err)
	fmt.Println("Compressed IDL length: ", len(compressedIDL))

	buffer := solana.NewWallet() // new keypair
	fmt.Println("Buffer Public Key: ", buffer.PublicKey().String())
	space := (8 + 32 + 4 + len(compressedIDL)) * 2

	// Step 1: Fetch rent-exempt lamports
	client := solRpc.New(solRpc.LocalNet_RPC)
	lamports, err := client.GetMinimumBalanceForRentExemption(
		context.Background(), uint64(space), solRpc.CommitmentFinalized,
	)
	require.NoError(t, err)

	// CREATE ACCOUNT
	createAccountIx, err := system.NewCreateAccountInstruction(
		lamports,
		uint64(space),
		programID,
		payer.PublicKey(),
		buffer.PublicKey(),
	).ValidateAndBuild()
	require.NoError(t, err)

	// INITIALIZE BUFFER ACCOUNT
	data, err := buildCreateBufferData()
	require.NoError(t, err)
	createBufferIx := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(buffer.PublicKey(), true, false),
			solana.NewAccountMeta(payer.PublicKey(), false, true),
		},
		data,
	)
	_, err = solCommonUtil.SendAndConfirm(
		context.Background(), client, []solana.Instruction{createAccountIx, createBufferIx}, payer, solRpc.CommitmentConfirmed, solCommonUtil.AddSigners(buffer.PrivateKey),
	)
	require.NoError(t, err)

	// write idl to buffer
	chunk := compressedIDL[0:100]
	require.NoError(t, err)
	data = binary.LittleEndian.AppendUint64([]byte{}, IDL_IX_TAG) // 4-byte Extend instruction identifier
	data = append(data, byte(2))
	data = append(data, chunk...)
	fmt.Println("Data:         ", hex.EncodeToString(data))

	writeIx := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(buffer.PublicKey(), true, false),
			solana.NewAccountMeta(payer.PublicKey(), false, true),
		},
		data,
	)
	_, err = solCommonUtil.SendAndConfirm(
		context.Background(), client, []solana.Instruction{writeIx}, payer, solRpc.CommitmentConfirmed, // <- show logs
	)
	require.NoError(t, err)

}

// Build instruction data for IdlInstruction::CreateBuffer
func buildCreateBufferData() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Step 1: Write the IDL_IX_TAG (8 bytes LE)
	if err := binary.Write(buf, binary.LittleEndian, IDL_IX_TAG); err != nil {
		return nil, err
	}

	// Step 2: Write the discriminator for "CreateBuffer"
	discriminator := byte(1)
	buf.Write([]byte{discriminator})

	return buf.Bytes(), nil
}

func TestConvertLegacyIdl(t *testing.T) {
	idlBytes, err := os.ReadFile("/Users/yashvardhan/chainlink/deployment/ccip/changeset/internal/solana_contracts/access_controller.json")
	require.NoError(t, err)
	var idl ccipChangesetSolana.IDL
	if err := json.Unmarshal(idlBytes, &idl); err != nil {
		fmt.Errorf("failed to parse legacy IDL: %w", err)
	}
	fmt.Println(idl)
	serialized, err := serializeIdl(idl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serialized & compressed IDL (%d bytes): %x\n", len(serialized), serialized)
}

func serializeIdl(idl any) ([]byte, error) {
	// Step 1: JSON encode
	jsonBytes, err := json.Marshal(idl)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}

	// Step 2: Compress using zlib (default compression level)
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(jsonBytes); err != nil {
		return nil, fmt.Errorf("zlib write failed: %w", err)
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("zlib close failed: %w", err)
	}

	return buf.Bytes(), nil
}

func TestIDLSetAuthorityLocal(t *testing.T) {
	programID := solana.MustPublicKeyFromBase58("GnpXUEvp4Uu5qzTRDHzDT5oRh7oY7g5GDRagc6zJo7M5")

	// derive idl address
	base, _, err := solana.FindProgramAddress([][]byte{}, programID)
	require.NoError(t, err)
	idlAddress, err := solana.CreateWithSeed(base, "anchor:idl", programID)
	require.NoError(t, err)
	fmt.Println("IDL Address:  ", idlAddress.String())

	deployerPrivKey, err := solana.PrivateKeyFromSolanaKeygenFile("/Users/yashvardhan/.config/solana/id.json")
	require.NoError(t, err)
	deployerKey := deployerPrivKey.PublicKey()
	newAuthorityPrivKey, err := solana.PrivateKeyFromSolanaKeygenFile("/Users/yashvardhan/.config/solana/id_devnet.json")
	require.NoError(t, err)
	newAuthority := newAuthorityPrivKey.PublicKey()

	data := binary.LittleEndian.AppendUint64([]byte{}, IDL_IX_TAG) // 4-byte Extend instruction identifier
	data = append(data, byte(4))
	data = append(data, newAuthority.Bytes()...)
	fmt.Println("Data:         ", hex.EncodeToString(data))

	instruction := solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(idlAddress, true, false),
			solana.NewAccountMeta(deployerKey, false, true), // Current upgrade authority (signer)
		},
		data,
	)
	client := solRpc.New(solRpc.LocalNet_RPC)
	_, err = solCommonUtil.SendAndConfirm(
		context.Background(), client, []solana.Instruction{instruction}, deployerPrivKey, solRpc.CommitmentConfirmed,
	)
	require.NoError(t, err)
}

func TestUploadIDL(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	// evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]
	_, err := commonchangeset.ApplyChangesetsV2(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.UploadIDL),
			ccipChangesetSolana.IDLConfig{
				ChainSelector: solChain,
			},
		),
	})
	require.NoError(t, err)
}
