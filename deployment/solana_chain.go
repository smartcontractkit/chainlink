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

	"github.com/gagliardetto/solana-go/rpc"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	ProgramIDPrefix                 = "Program Id: "
	BufferIDPrefix                  = "Buffer: "
	SolDefaultCommitment            = rpc.CommitmentConfirmed
	RouterProgramName               = "ccip_router"
	OffRampProgramName              = "ccip_offramp"
	FeeQuoterProgramName            = "fee_quoter"
	BurnMintTokenPoolProgramName    = "burnmint_token_pool"
	LockReleaseTokenPoolProgramName = "lockrelease_token_pool"
	AccessControllerProgramName     = "access_controller"
	TimelockProgramName             = "timelock"
	McmProgramName                  = "mcm"
	RMNRemoteProgramName            = "rmn_remote"
	ReceiverProgramName             = "test_ccip_receiver"
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

// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
// https://docs.google.com/document/d/1nCNuam0ljOHiOW0DUeiZf4ntHf_1Bw94Zi7ThPGoKR4/edit?tab=t.0#heading=h.hju45z55bnqd
func GetSolanaProgramBytes() map[string]int {
	return map[string]int{
		RouterProgramName:               5 * 1024 * 1024,
		OffRampProgramName:              0 * 1024 * 1024, // router should be redeployed but it does support upgrades if required (big fixes etc.)
		FeeQuoterProgramName:            5 * 1024 * 1024,
		BurnMintTokenPoolProgramName:    3 * 1024 * 1024,
		LockReleaseTokenPoolProgramName: 3 * 1024 * 1024,
		AccessControllerProgramName:     1 * 1024 * 1024,
		TimelockProgramName:             1 * 1024 * 1024,
		McmProgramName:                  1 * 1024 * 1024,
		RMNRemoteProgramName:            3 * 1024 * 1024,
	}
}

func (c SolChain) DeployProgram(logger logger.Logger, programName string, isUpgrade bool) (string, error) {
	programFile := filepath.Join(c.ProgramsPath, programName+".so")
	if _, err := os.Stat(programFile); err != nil {
		return "", fmt.Errorf("program file not found: %w", err)
	}
	programKeyPair := filepath.Join(c.ProgramsPath, programName+"-keypair.json")

	cliCommand := "deploy"
	prefix := ProgramIDPrefix
	if isUpgrade {
		cliCommand = "write-buffer"
		prefix = BufferIDPrefix
	}

	// Base command with required args
	baseArgs := []string{
		"program", cliCommand,
		programFile,                // .so file
		"--keypair", c.KeypairPath, // deployer keypair
		"--url", c.URL, // rpc url
		"--use-rpc", // use rpc for deployment
	}

	var cmd *exec.Cmd
	// We need to specify the program ID on the initial deploy but not on upgrades
	// Upgrades happen in place so we don't need to supply the keypair
	// It will write the .so file to a buffer and then deploy it to the existing keypair
	if !isUpgrade {
		logger.Infow("Deploying program with existing keypair",
			"programFile", programFile,
			"programKeyPair", programKeyPair)
		baseArgs = append(baseArgs, "--program-id", programKeyPair)
		totalBytes := GetSolanaProgramBytes()[programName]
		if totalBytes > 0 {
			baseArgs = append(baseArgs, "--max-len", strconv.Itoa(totalBytes))
		}
		cmd = exec.Command("solana", baseArgs...) // #nosec G204
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
	return parseProgramID(output, prefix)
}

func (c SolChain) GetAccountDataBorshInto(ctx context.Context, pubkey solana.PublicKey, accountState interface{}) error {
	err := solCommonUtil.GetAccountDataBorshInto(ctx, c.Client, pubkey, SolDefaultCommitment, accountState)
	if err != nil {
		return err
	}
	return nil
}

// parseProgramID parses the program ID from the deploy output.
func parseProgramID(output string, prefix string) (string, error) {
	// Look for the program ID in the CLI output
	// Example output: "Program Id: <PROGRAM_ID>"
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
