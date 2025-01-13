package deployment

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	// RPC client
	Client *solRpc.Client
	URL    string
	WSURL  string
	// TODO: raw private key for now, need to replace with a more secure way
	DeployerKey *solana.PrivateKey
	// deploy uses the solana CLI which needs a keyfile
	KeypairPath  string
	ProgramsPath string
	Confirm      func(instructions []solana.Instruction, opts ...solCommomUtil.TxModifier) error
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

func (c SolChain) DeployProgram(programName string) (string, error) {
	programFile := fmt.Sprintf("%s/%s.so", c.ProgramsPath, programName)
	programKeyPair := fmt.Sprintf("%s/%s-keypair.json", c.ProgramsPath, programName)

	// Construct the CLI command: solana program deploy
	// TODO: @terry doing this on the fly
	cmd := exec.Command("solana", "program", "deploy", programFile, "--keypair", c.KeypairPath, "--program-id", programKeyPair)

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
	startIdx := strings.Index(output, prefix)
	if startIdx == -1 {
		return "", errors.New("failed to find program ID in output")
	}
	startIdx += len(prefix)
	endIdx := strings.Index(output[startIdx:], "\n")
	if endIdx == -1 {
		endIdx = len(output)
	}
	return output[startIdx : startIdx+endIdx], nil
}
