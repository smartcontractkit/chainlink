package deployment

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	SolDefaultCommitment = rpc.CommitmentConfirmed
	SolDefaultGasLimit   = solBinary.Uint128{Lo: 3000, Hi: 0, Endianness: nil}
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

	if _, err := os.Stat(programKeyPair); err == nil {
		baseArgs = append(baseArgs, "--program-id", programKeyPair)
		logger.Infow("Deploying program with existing keypair",
			"programFile", programFile,
			"programKeyPair", programKeyPair)
	} else {
		logger.Infow("Deploying new program", "programFile", programFile)
	}

	// Create context with timeout to run the deploy command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create command without context
	cmd := exec.Command("solana", baseArgs...)
	logger.Infow("Running deploy program command", "cmd", strings.Join(cmd.Args, " "))

	// Connect standard streams with both buffer and real-time logging
	var stdout, stderr bytes.Buffer
	cmd.Stdin = os.Stdin

	// Set up pipe for stdout
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdout pipe: %w", err)
	}
	// Set up pipe for stderr
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting program deployment: %w", err)
	}

	// Create a channel to signal command completion
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Copy output to both buffer and logger in real-time
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			logger.Infow("Program deployment stdout", "line", line)
			stdout.WriteString(line + "\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			logger.Infow("Program deployment stderr", "line", line)
			stderr.WriteString(line + "\n")
		}
	}()

	// Wait for either completion or timeout
	select {
	case <-ctx.Done():
		logger.Errorw("Program deployment timed out",
			"stdout", stdout.String(),
			"stderr", stderr.String())
		// Try to kill the process and its children
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGTERM)
		}
		cmd.Process.Kill()
		return "", fmt.Errorf("deployment timed out after 5 minutes")
	case err := <-done:
		if err != nil {
			logger.Errorw("Program deployment failed",
				"error", err,
				"stdout", stdout.String(),
				"stderr", stderr.String())
			return "", fmt.Errorf("error deploying program: %s: %s", err.Error(), stderr.String())
		}
	}

	outputStr := stdout.String()
	logger.Infow("Program deployment successful",
		"stdout", outputStr,
		"stderr", stderr.String())
	return parseProgramID(outputStr)
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

// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName string) (solana.PublicKey, error) {
	tokenPrograms := map[string]solana.PublicKey{
		"spl-token":      solana.TokenProgramID,
		"spl-token-2022": solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: spl-token, spl-token-2022", programName)
	}
	return programID, nil
}

// GetPoolType returns the token pool type constant for the given string
func GetPoolType(poolType string) (token_pool.PoolType, error) {
	poolTypes := map[string]token_pool.PoolType{
		"LockAndRelease": token_pool.LockAndRelease_PoolType,
		"BurnAndMint":    token_pool.BurnAndMint_PoolType,
	}

	poolTypeConstant, ok := poolTypes[poolType]
	if !ok {
		return 0, fmt.Errorf("invalid pool type: %s. Must be one of: LockAndRelease, BurnAndMint", poolType)
	}
	return poolTypeConstant, nil
}

func FindTokenAddress(e Environment, chainSelector uint64, tokenName string) (solana.PublicKey, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return solana.PublicKey{}, err
	}

	tv := NewTypeAndVersion(ContractType(tokenName), Version1_0_0)
	for address, tvStr := range addresses {
		if tvStr == tv {
			return solana.MustPublicKeyFromBase58(address), nil
		}
	}
	return solana.PublicKey{}, fmt.Errorf("token address not found in address book: %s", tokenName)
}

func (c SolChain) GetAccountDataBorshInto(ctx context.Context, pubkey solana.PublicKey, accountState interface{}) error {
	err := solCommonUtil.GetAccountDataBorshInto(ctx, c.Client, pubkey, SolDefaultCommitment, accountState)
	if err != nil {
		return err
	}
	return nil
}
