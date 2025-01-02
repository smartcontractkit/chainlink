package changeset_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"bytes"
	"fmt"
	"os/exec"

	// "github.com/stretchr/testify/require"
	// "go.uber.org/zap/zapcore"
	// "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	// "github.com/smartcontractkit/chainlink/deployment/environment/memory"
	// "github.com/smartcontractkit/chainlink/v2/core/logger"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/external_program_cpi_stub"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

var (
	PrintEvents       = true
	DefaultCommitment = rpc.CommitmentConfirmed
	StubProgram       = "EQPCTRibpsPcQNb464QVBkS1PkFfuK8kYdpd5Y17HaGh"

	// CcipRouterProgram          = solana.MustPublicKeyFromBase58("ZWzN9gkPMkNRh4CBt8n4vA6VfZxvHCD2FPa98U4gvRP")
	CcipReceiverProgram        = solana.MustPublicKeyFromBase58("CtEVnHsQzhTNWav8skikiV2oF6Xx7r7uGGa8eCDQtTjH")
	CcipReceiverAddress        = solana.MustPublicKeyFromBase58("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb")
	CcipInvalidReceiverProgram = solana.MustPublicKeyFromBase58("9Vjda3WU2gsJgE4VdU6QuDw8rfHLyigfFyWs3XDPNUn8")
	CcipTokenPoolProgram       = solana.MustPublicKeyFromBase58("GRvFSLwR7szpjgNEZbGe4HtxfJYXqySXuuRUAJDpu4WH")
	Token2022Program           = solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")

	ReceiverTargetAccountPDA, _, _           = solana.FindProgramAddress([][]byte{[]byte("counter")}, CcipReceiverProgram)
	ReceiverExternalExecutionConfigPDA, _, _ = solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, CcipReceiverProgram)
	// BillingSignerPDA, _, _                   = solana.FindProgramAddress([][]byte{[]byte("fee_billing_signer")}, CcipRouterProgram)

	BillingTokenConfigPrefix = []byte("fee_billing_token_config")
	DestChainConfigPrefix    = []byte("destination_billing_config")

	SolanaChainSelector uint64 = 15
	EvmChainSelector    uint64 = 21

	// SolanaChainStatePDA, _, _ = solana.FindProgramAddress([][]byte{[]byte("chain_state"), binary.LittleEndian.AppendUint64([]byte{}, SolanaChainSelector)}, CcipRouterProgram)
	EvmChainLE = common.Uint64ToLE(EvmChainSelector)
	// EvmChainStatePDA, _, _    = solana.FindProgramAddress([][]byte{[]byte("chain_state"), binary.LittleEndian.AppendUint64([]byte{}, EvmChainSelector)}, CcipRouterProgram)

	OnRampAddress        = []byte{1, 2, 3}
	EnableExecutionAfter = int64(1800) // 30min

	MaxOracles                      = 16
	OcrF                      uint8 = 5
	ConfigDigest                    = common.MakeRandom32ByteArray()
	Empty24Byte                     = [24]byte{}
	MaxSignersAndTransmitters       = 16
)

func setDevNet(keypairPath string) error {
	// Construct the CLI command: solana program deploy
	cmd := exec.Command("solana", "config", "set", "--url", "localhost", "--keypair", keypairPath)

	// Capture the command output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting config: %s: %s", err.Error(), stderr.String())
	}
	return nil
}

// TestDeployProgram is a test for deploying the Solana program.
func TestDeployProgram(t *testing.T) {
	// Path to your .so file and keypair file
	// programFile := "/Users/yashvardhan/chainlink-internal-integrations/solana/contracts/target/deploy/external_program_cpi_stub.so"
	keypairPath := "/Users/yashvardhan/.config/solana/id.json" //wallet
	// programKeyPair := "/Users/yashvardhan/chainlink-internal-integrations/solana/contracts/target/deploy/external_program_cpi_stub-keypair.json"
	// keypairPath := "/Users/yashvardhan/chainlink-internal-integrations/solana/contracts/target/deploy/external_program_cpi_stub-keypair.json"

	ExternalCpiStubProgram := solana.MustPublicKeyFromBase58("EQPCTRibpsPcQNb464QVBkS1PkFfuK8kYdpd5Y17HaGh")
	solanaGoClient := rpc.New("http://127.0.0.1:8899")

	// Fetch account info for the program ID
	account, _ := solanaGoClient.GetAccountInfo(context.Background(), ExternalCpiStubProgram)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// Deploy program if it doesn't exist
	if account != nil && account.Value.Executable {
		fmt.Println("Program exists and is executable.")
	} else {
		fmt.Println("Program does not exist or is not executable.")
		// Deploy the program
		programID, err := deployment.DeploySolProgramCLI("external_program_cpi_stub")
		if err != nil {
			t.Fatalf("Failed to deploy program: %v", err)
		}
		// Verify the program ID (simple check for non-empty string)
		if programID == "" {
			t.Fatalf("Program ID is empty")
		}

		t.Logf("programID %s", programID)
	}

	// program should exist by now (either already deployed, or deployed and waited for confirmation)
	external_program_cpi_stub.SetProgramID(ExternalCpiStubProgram)

	// wallet keys
	privateKey, _ := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
	publicKey := privateKey.PublicKey()

	// this is a PDA that gets initialised when you call init on the programID
	StubAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("u8_value")}, ExternalCpiStubProgram)
	t.Logf("StubAccountPDA %s", StubAccountPDA)

	// check if the PDA is already initialised
	accountInfo, _ := solanaGoClient.GetAccountInfo(context.Background(), ExternalCpiStubProgram)

	var ix *external_program_cpi_stub.Instruction
	var err error

	if accountInfo != nil {
		fmt.Println("Account initialized successfully.")
		fmt.Printf("Account data: %v\n", accountInfo.Value.Data)
		ix, err = external_program_cpi_stub.NewEmptyInstruction().ValidateAndBuild()
	} else {
		fmt.Println("Account does not exist or has no data.")
		ix, err = external_program_cpi_stub.NewInitializeInstruction(
			StubAccountPDA,
			publicKey,
			solana.SystemProgramID, // 1111111
		).ValidateAndBuild()
	}

	_, err = common.SendAndConfirm(context.Background(), solanaGoClient, []solana.Instruction{ix}, privateKey, config.DefaultCommitment)

	require.NoError(t, err)
}

func spinUpDevNet(t *testing.T) (string, string) {
	t.Helper()
	port := "8899"
	portInt, _ := strconv.Atoi(port)

	faucetPort := "8877"
	url := "http://127.0.0.1:" + port
	wsURL := "ws://127.0.0.1:" + strconv.Itoa(portInt+1)

	args := []string{
		"--reset",
		"--rpc-port", port,
		"--faucet-port", faucetPort,
		"--ledger", t.TempDir(),
	}

	cmd := exec.Command("solana-test-validator", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		assert.NoError(t, cmd.Process.Kill())
		if err2 := cmd.Wait(); assert.Error(t, err2) {
			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
				t.Logf("solana-test-validator\n stdout: %s\n stderr: %s", stdOut.String(), stdErr.String())
			}
		}
	})

	// Wait for api server to boot
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		client := rpc.New(url)
		out, err := client.GetHealth(tests.Context(t))
		if err != nil || out != rpc.HealthOk {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		}
		ready = true
		break
	}
	if !ready {
		t.Logf("Cmd output: %s\nCmd error: %s\n", stdOut.String(), stdErr.String())
	}
	require.True(t, ready)
	t.Logf("Solana Devnet spun up successfully")

	return url, wsURL
}

func getRpcClient(t *testing.T) *rpc.Client {
	url, _ := spinUpDevNet(t)
	// url, _ := memory.solChain(t)
	return rpc.New(url)
}

// Added TestDeployLinkTokenSol -> so this is not required anymore
// func TestTokenDeploy(t *testing.T) {
// 	solanaGoClient := getRpcClient(t)
// 	adminPrivateKey := deployment.GetSolanaDeployerKey()
// 	adminPublicKey := adminPrivateKey.PublicKey()
// 	mint, _ := solana.NewRandomPrivateKey()
// 	mintPublicKey := mint.PublicKey()
// 	instructions, err := tokens.CreateToken(context.Background(), config.Token2022Program, mintPublicKey, adminPublicKey, uint8(0), solanaGoClient, DefaultCommitment)
// 	require.NoError(t, err)
// 	_, err = common.SendAndConfirm(context.Background(), solanaGoClient, instructions, adminPrivateKey, DefaultCommitment, common.AddSigners(mint))
// 	require.NoError(t, err)
// }

func TestCcipRouterDeploy(t *testing.T) {
	// Path to your .so file and keypair file
	// programFile := "/Users/ttata/dev/chainlink-ccip/chains/solana/contracts/target/deploy/ccip_router.so"
	keypairPath := "/Users/yashvardhan/.config/solana/id.json" //wallet
	// programKeyPair := "/Users/ttata/dev/chainlink-ccip/chains/solana/contracts/target/deploy/ccip_router-keypair.json"

	adminPrivateKey, _ := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
	adminPublicKey := adminPrivateKey.PublicKey()
	solanaGoClient := getRpcClient(t)
	err := setDevNet(keypairPath)
	require.NoError(t, err)
	ctx := context.Background()
	testutils.FundAccounts(ctx, []solana.PrivateKey{adminPrivateKey}, solanaGoClient, t)

	// get program data account before deploying, hitting NotFound error
	// data, err := solanaGoClient.GetAccountInfoWithOpts(ctx, CcipRouterProgram, &rpc.GetAccountInfoOpts{
	// 	Commitment: DefaultCommitment,
	// })
	// require.ErrorAs(t, err, &rpc.ErrNotFound)
	// Deploy the program
	programID, err := deployment.DeploySolProgramCLI("ccip_router")
	CcipRouterProgram := solana.MustPublicKeyFromBase58(programID)
	if err != nil {
		t.Fatalf("Failed to deploy program: %v", err)
	}
	// get program data account
	data, err := solanaGoClient.GetAccountInfoWithOpts(ctx, CcipRouterProgram, &rpc.GetAccountInfoOpts{
		Commitment: DefaultCommitment,
	})
	require.NoError(t, err)
	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	require.NoError(t, bin.UnmarshalBorsh(&programData, data.Bytes()))
	RouterConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, CcipRouterProgram)
	RouterStatePDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("state")}, CcipRouterProgram)
	ExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, CcipRouterProgram)
	ExternalTokenPoolsSignerPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_token_pools_signer")}, CcipRouterProgram)
	ccip_router.SetProgramID(CcipRouterProgram)
	instruction, err := ccip_router.NewInitializeInstruction(
		SolanaChainSelector,  // chain selector
		bin.Uint128{},        // default gas limit
		true,                 // allow out of order execution
		EnableExecutionAfter, // period to wait before allowing manual execution
		RouterConfigPDA,
		RouterStatePDA,
		adminPublicKey,
		solana.SystemProgramID,
		CcipRouterProgram,
		programData.Address,
		ExternalExecutionConfigPDA,
		ExternalTokenPoolsSignerPDA,
	).ValidateAndBuild()
	require.NoError(t, err)

	// skip preflight, txs with init PDAs will fail preflight
	_, err = common.SendAndConfirm(ctx, solanaGoClient, []solana.Instruction{instruction}, adminPrivateKey, DefaultCommitment)
	require.NoError(t, err)
	t.Logf("Program deployed successfully with ID: %s", programID)
}
