package deployment

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/rpc"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	SolDefaultCommitment        = rpc.CommitmentConfirmed
	SolDefaultGasLimit          = solBinary.Uint128{Lo: 3000, Hi: 0, Endianness: nil}
	SolDefaultMaxFeeJuelsPerMsg = solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil}
	SPL2022Tokens               = "SPL2022Tokens"
	SPLTokens                   = "SPLTokens"
	EnableExecutionAfter        = int64(1800) // 30min
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
	Confirm     func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error

	// deploy uses the solana CLI which needs a keyfile
	KeypairPath  string
	ProgramsPath string
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
	programFile := filepath.Join(c.ProgramsPath, programName+".so")
	if _, err := os.Stat(programFile); err != nil {
		return "", fmt.Errorf("program file not found: %w", err)
	}
	programKeyPair := filepath.Join(c.ProgramsPath, programName+"-keypair.json")

	// Base command with required args
	baseArgs := []string{
		"program", "deploy",
		programFile,                // .so file
		"--keypair", c.KeypairPath, // program keypair
		"--url", c.URL, // rpc url
	}

	var cmd *exec.Cmd
	if _, err := os.Stat(programKeyPair); err == nil {
		// Keypair exists, include program-id
		logger.Infow("Deploying program with existing keypair",
			"programFile", programFile,
			"programKeyPair", programKeyPair)
		cmd = exec.Command("solana", append(baseArgs, "--program-id", programKeyPair)...) // #nosec G204
	} else {
		// Keypairs wont be created for devenvs
		logger.Infow("Deploying new program",
			"programFile", programFile)
		cmd = exec.Command("solana", baseArgs...) // #nosec G204
	}

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

	// TODO: obviously need to do this better
	time.Sleep(5 * time.Second)
	return parseProgramID(output)
}

func (c SolChain) GetAccountDataBorshInto(ctx context.Context, pubkey solana.PublicKey, accountState interface{}) error {
	err := solCommonUtil.GetAccountDataBorshInto(ctx, c.Client, pubkey, SolDefaultCommitment, accountState)
	if err != nil {
		return err
	}
	return nil
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
