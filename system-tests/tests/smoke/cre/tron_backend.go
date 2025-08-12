package cre

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TronBackend manages the Tron simulated backend for testing
type TronBackend struct {
	ChainID            string
	FullNodeURL        string
	SolidityURL        string
	JSONRPCFullNodeURL string
	JSONRPCSolidityURL string
	started            bool
	logger             zerolog.Logger
}

const (
	DefaultTronChainID      = "728126428"                                                        // Tron Mainnet chain ID
	DefaultTronPrivateKey   = "da146374a75310b9666e834ee4ad0866d6f4035967bfc76217c5a495fff9f0d0" // Default Tron private key
	TronFullNodePort        = "16667"
	TronSolidityPort        = "16668"
	TronJSONRPCFullNodePort = "16671"
	TronJSONRPCSolidityPort = "16672"
	// Local RPCs
	DefaultInternalFullNodeUrl     = "http://127.0.0.1:16667/wallet"
	DefaultInternalSolidityNodeUrl = "http://127.0.0.1:16668/walletsolidity"

	// Testnet RPCs
	// Urls can be found at https://developers.tron.network/reference/background
	ShastaFullNodeUrl     = "https://api.shasta.trongrid.io/wallet"
	ShastaSolidityNodeUrl = "https://api.shasta.trongrid.io/walletsolidity"

	NileFullNodeUrl     = "https://nile.trongrid.io/wallet"
	NileSolidityNodeUrl = "https://nile.trongrid.io/walletsolidity"

	// Configs for TXM
	DevnetFeeLimit                  = 1_000_000_000
	DevnetMaxWaitTime               = 30 //seconds
	DevnetPollFrequency             = 1  //seconds
	DevnetOcrTransmissionFrequency  = 5 * time.Second
	TestnetFeeLimit                 = 10_000_000_000
	TestnetMaxWaitTime              = 90 //seconds
	TestnetPollFrequency            = 5  //seconds
	TestnetOcrTransmissionFrequency = 10 * time.Second

	// Testing network names
	Shasta = "shasta"
	Nile   = "nile"
	Devnet = "devnet"
)

// NewTronBackend creates a new Tron backend manager
func NewTronBackend(logger zerolog.Logger, chainID string) *TronBackend {
	if chainID == "" {
		chainID = DefaultTronChainID
	}

	internalFullNodeURL := fmt.Sprintf("http://127.0.0.1:%s/wallet", TronFullNodePort)
	internalSolidityURL := fmt.Sprintf("http://127.0.0.1:%s/walletsolidity", TronSolidityPort)

	// Convert to JSON-RPC endpoints (Tron uses JSON-RPC for HTTP, no WebSocket support)
	jsonRpcFullNodeUrl := fmt.Sprintf("http://127.0.0.1:%s/jsonrpc", TronJSONRPCFullNodePort)
	jsonRpcSolidityUrl := fmt.Sprintf("http://127.0.0.1:%s/jsonrpc", TronJSONRPCSolidityPort)

	return &TronBackend{
		ChainID:            chainID,
		FullNodeURL:        internalFullNodeURL,
		SolidityURL:        internalSolidityURL,
		JSONRPCFullNodeURL: jsonRpcFullNodeUrl,
		JSONRPCSolidityURL: jsonRpcSolidityUrl,
		started:            false,
		logger:             logger,
	}
}

func StartTronNode(genesisAddress string) error {
	gitRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find Git root: %+w", err)
	}

	scriptPath := filepath.Join(gitRoot, "system-tests/tests/smoke/cre/scripts/java-tron.sh")
	cmd := exec.Command(scriptPath, genesisAddress)

	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Failed to start java-tron, dumping output:\n%s\n", string(output))
			return fmt.Errorf("Failed to start java-tron, bad exit code: %v", exitError.ExitCode())
		}
		return fmt.Errorf("Failed to start java-tron: %+w", err)
	}

	return nil
}

func StopTronNode() error {
	gitRoot, err := FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find Git root: %+w", err)
	}

	scriptPath := filepath.Join(gitRoot, "system-tests/tests/smoke/cre/scripts/java-tron.down.sh")
	cmd := exec.Command(scriptPath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Failed to stop java-tron, dumping output:\n%s\n", string(output))
			return fmt.Errorf("Failed to start java-tron, bad exit code: %v", exitError.ExitCode())
		}
		return fmt.Errorf("Failed to stop java-tron: %+w", err)
	}

	return nil
}

func ShutdownTronBackend(t *testing.T) {
	err := StopTronNode()
	require.NoError(t, err, "failed to stop Tron node")
}

// Finds the closest git repo root, assuming that a directory with a .git directory is a git repo.
func FindGitRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(currentDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", fmt.Errorf("no Git repository found")
		}

		currentDir = parentDir
	}
}

func GetTronNodeIpAddress() string {
	if runtime.GOOS == "darwin" {
		return "127.0.0.1"
	} else {
		return "172.255.0.101"
	}
}

// GetChainlinkConfig returns the Chainlink configuration for Tron
func (tb *TronBackend) GetChainlinkConfig() string {
	return fmt.Sprintf(`
[[EVM]]
Enabled = true
ChainID = '%s'
ChainType = 'tron'
FinalityTagEnabled = false
LogBroadcasterEnabled = false
FinalityDepth = 1
LogPollInterval = '10s'
NoNewHeadsThreshold = '3m0s'
FinalizedBlockOffset = 0

[[EVM.Nodes]]
Name = 'primary'
HTTPURL = '%s'

[EVM.BalanceMonitor]
Enabled = true

[EVM.GasEstimator]
PriceMin = '210 wei'
LimitDefault = 6000000
LimitMax = 6000000

[EVM.NodePool]
NewHeadsPollInterval = '10s'

[EVM.HeadTracker]
FinalityTagBypass = false

[Feature]
LogPoller = true

[Log]
Level = 'debug'

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
DeltaDial = '5s'
DeltaReconcile = '5s'
ListenAddresses = ['0.0.0.0:6691']

[WebServer]
HTTPPort = 6688
[WebServer.TLS]
HTTPSPort = 0
`, tb.ChainID, tb.JSONRPCFullNodeURL)
}

// SetTronPrivateKeyEnv sets the PRIVATE_KEY environment variable for Tron
func SetTronPrivateKeyEnv() error {
	if os.Getenv("PRIVATE_KEY") == "" {
		err := os.Setenv("PRIVATE_KEY", DefaultTronPrivateKey)
		if err != nil {
			return fmt.Errorf("failed to set PRIVATE_KEY env var: %w", err)
		}
	}
	return nil
}
