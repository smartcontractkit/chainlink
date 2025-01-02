package deployment

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
)

// TODO: hard coding these for test, need to figure out the dynamic way like evm
var (
	SolanaChainSelector uint64 = 16423721717087811551                        //devnet
	keypairPath                = "/Users/yashvardhan/.config/solana/id.json" //wallet
	deployBinPath              = "/Users/yashvardhan/chainlink-ccip/chains/solana/contracts/target/deploy"
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	// RPC cient
	Client *rpc.Client
	// TODO: raw private key for now, need to replace with a more secure way
	DeployerKey *solana.PrivateKey
	Confirm     func(instructions []solana.Instruction, opts ...solCommomUtil.TxModifier) error
}

func (c SolChain) String() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	return fmt.Sprintf("%s (%d)", chainInfo.ChainName, chainInfo.ChainSelector)
}

func (c SolChain) Name() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	if chainInfo.ChainName == "" {
		return fmt.Sprintf("%d", c.Selector)
	}
	return chainInfo.ChainName
}

type SolClient interface {
}

type ContractDeploySolana struct {
	ProgramID *solana.PublicKey // We leave this incase a Go binding doesn't have Address()
	Tv        TypeAndVersion
	Err       error
}

func DeploySolProgramCLI(programName string) (string, error) {
	programFile := fmt.Sprintf("%s/%s.so", deployBinPath, programName)
	programKeyPair := fmt.Sprintf("%s/%s-keypair.json", deployBinPath, programName)

	// Construct the CLI command: solana program deploy
	// TODO: @terry doing this on the fly
	cmd := exec.Command("solana", "program", "deploy", programFile, "--keypair", keypairPath, "--program-id", programKeyPair)

	// Capture the command output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error deploying program: %s: %s", err.Error(), stderr.String())
	}

	// Parse and return the program ID
	output := stdout.String()

	time.Sleep(5 * time.Second) // obviously need to do this better
	return parseProgramID(output)
}

// parseProgramID parses the program ID from the deploy output.
func parseProgramID(output string) (string, error) {
	// Look for the program ID in the CLI output
	// Example output: "Program Id: <PROGRAM_ID>"
	const prefix = "Program Id: "
	startIdx := bytes.Index([]byte(output), []byte(prefix))
	if startIdx == -1 {
		return "", fmt.Errorf("failed to find program ID in output")
	}
	startIdx += len(prefix)
	endIdx := bytes.Index([]byte(output[startIdx:]), []byte("\n"))
	if endIdx == -1 {
		endIdx = len(output)
	}
	return output[startIdx : startIdx+endIdx], nil
}

func GetSolanaDeployerKey() solana.PrivateKey {
	adminPrivateKey, _ := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
	return adminPrivateKey
}
