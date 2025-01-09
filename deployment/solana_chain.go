package deployment

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	deployBinPath = "/Users/yashvardhan/chainlink-ccip/chains/solana/contracts/target/deploy"
	// keypairPath   = "/Users/yashvardhan/.config/solana/id.json" //wallet
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	// RPC cient
	Client *solRpc.Client
	// TODO: raw private key for now, need to replace with a more secure way
	DeployerKey *solana.PrivateKey
	Confirm     func(instructions []solana.Instruction, opts ...solCommomUtil.TxModifier) error
	URL         string
	WSURL       string
	KeypairPath string
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
		return strconv.FormatUint(c.Selector, 10)
	}
	return chainInfo.ChainName
}

func (c SolChain) DeployProgram(logger logger.Logger, programName string) (string, error) {
	programFile := fmt.Sprintf("%s/%s.so", deployBinPath, programName)
	programKeyPair := fmt.Sprintf("%s/%s-keypair.json", deployBinPath, programName)

	logger.Infow("c.KeypairPath", "path", c.KeypairPath)
	logger.Infow("private key", "key", c.DeployerKey)
	key, err := solana.PrivateKeyFromSolanaKeygenFile(c.KeypairPath)
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}
	logger.Infow("program key pair", "key", key)
	cmd := exec.Command("solana", "program", "deploy", programFile, "--keypair", c.KeypairPath, "--program-id", programKeyPair, "--url", c.URL)
	// cmd := exec.Command("solana", "program", "deploy", programFile, "--upgrade-authority", c.DeployerKey.PublicKey().String(), "--program-id", programKeyPair, "--url", c.URL)

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
