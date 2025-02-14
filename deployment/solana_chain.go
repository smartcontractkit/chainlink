package deployment

import (
	"bytes"
	"context"
	"encoding/json"
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
		"--keypair", c.KeypairPath, // deployer keypair
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

func (c SolChain) UpgradeProgram(logger logger.Logger, programName string, programID solana.PublicKey, upgradeAuthority solana.PrivateKey) (string, error) {
	programFile := filepath.Join(c.ProgramsPath, programName+".so")
	if _, err := os.Stat(programFile); err != nil {
		return "", fmt.Errorf("program file not found: %w", err)
	}

	// Create a temporary file for the upgrade authority keypair
	tempFile, err := os.CreateTemp("", "upgrade-authority-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure cleanup

	keypairInts := make([]int, len(upgradeAuthority))
	for i, b := range upgradeAuthority {
		keypairInts[i] = int(b) // Convert to int slice to ensure JSON array format
	}

	keypairJSON, err := json.Marshal(keypairInts)
	if err != nil {
		return "", fmt.Errorf("failed to serialize private key: %w", err)
	}

	if _, err := tempFile.Write(keypairJSON); err != nil {
		tempFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close() // Close before passing to external command

	// Step 1: Get the current program account size
	currentSize, err := getProgramAccountSize(programID.String())
	if err != nil {
		return "", fmt.Errorf("error getting current program account size: %w", err)
	}
	logger.Debugw("Current program account size", "size", currentSize)

	// Step 2: Get the size of the new program binary
	newSize, err := getFileSize(programFile)
	if err != nil {
		return "", fmt.Errorf("error getting new program binary size: %w", err)
	}
	logger.Debugw("New program binary size", "size", newSize)

	// Step 3: Check if additional space is needed
	if newSize > currentSize {
		additionalSpace := newSize - currentSize
		logger.Debugw("Additional space required", "size", additionalSpace)

		// Step 4: Extend the program account
		err = extendProgramAccount(programID.String(), additionalSpace, tempFile.Name())
		if err != nil {
			return "", fmt.Errorf("error extending program account: %w", err)
		}
		logger.Debug("Program account extended successfully")
	} else {
		logger.Debug("No additional space needed for program account")
	}

	// Step 5: Write the new program binary to a buffer account
	bufferAddress, err := writeBuffer(programFile, tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error writing buffer: %w", err)
	}
	logger.Debugw("Buffer account created", "address", bufferAddress)

	// Step 6: Upgrade the program using the buffer
	newProgramID, err := upgradeProgram(programID.String(), bufferAddress, tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error upgrading program: %w", err)
	}
	logger.Debug("Program upgraded successfully")

	// Step 7: Close the buffer account to reclaim SOL
	err = closeBuffer(bufferAddress, tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error closing buffer account: %w", err)
	}
	logger.Debug("Buffer account closed successfully")
	return newProgramID, nil
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

// getProgramAccountSize retrieves the current size of the program account
func getProgramAccountSize(programID string) (int, error) {
	cmd := exec.Command("solana", "program", "show", programID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	// Parse the output to find the data length
	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Data Length:") {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				return 0, fmt.Errorf("invalid data length line: %s", line)
			}
			sizeStr := strings.TrimSpace(parts[1])
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				return 0, fmt.Errorf("failed to parse data length: %v", err)
			}
			return size, nil
		}
	}
	return 0, fmt.Errorf("data length not found in program account info")
}

// getFileSize returns the size of a file in bytes
func getFileSize(filePath string) (int, error) {
	cmd := exec.Command("stat", "-c%s", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	sizeStr := strings.TrimSpace(out.String())
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse file size: %v", err)
	}
	return size, nil
}

// extendProgramAccount extends the program account by the specified additional space
func extendProgramAccount(programID string, additionalSpace int, upgradeAuthorityKey string) error {
	cmd := exec.Command("solana", "program", "extend", "--keypair", upgradeAuthorityKey, programID, strconv.Itoa(additionalSpace))
	return cmd.Run()
}

// writeBuffer writes the new program binary to a buffer account and returns the buffer address
func writeBuffer(newProgramBinary, upgradeAuthorityKey string) (string, error) {
	cmd := exec.Command("solana", "program", "write-buffer", newProgramBinary, "--keypair", upgradeAuthorityKey)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Parse the output to extract the buffer address
	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Buffer:") {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				return "", fmt.Errorf("invalid buffer address line: %s", line)
			}
			bufferAddress := strings.TrimSpace(parts[1])
			return bufferAddress, nil
		}
	}
	return "", fmt.Errorf("buffer address not found in output")
}

// upgradeProgram upgrades the program using the buffer account
func upgradeProgram(programID, bufferAddress, upgradeAuthorityKey string) (string, error) {
	cmd := exec.Command("solana", "program", "upgrade", "--keypair", upgradeAuthorityKey, programID, bufferAddress)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Parse the output to find the new program ID
	output := out.String()
	fmt.Println(output)
	return parseProgramID(output)
}

// closeBuffer closes the buffer account and reclaims the SOL
func closeBuffer(bufferAddress, upgradeAuthorityKey string) error {
	cmd := exec.Command("solana", "program", "close", "--keypair", upgradeAuthorityKey, "--buffers", bufferAddress)
	return cmd.Run()
}
