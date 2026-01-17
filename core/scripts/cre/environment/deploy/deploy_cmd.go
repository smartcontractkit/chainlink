package deploy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
)

var (
	rpcURL       string
	privateKey   string
	chainID      int64
	donID        uint32
	outputFile   string
)

// DeployCmd is the cobra command for deploying LLO contracts
var DeployCmd = &cobra.Command{
	Use:   "deploy-llo",
	Short: "Deploy LLO contracts (Configurator and ChannelConfigStore) for E2E testing",
	Long: `Deploy the LLO infrastructure contracts needed for the streams-trigger E2E test.

This command deploys:
1. Configurator - Stores OCR configuration for the LLO DON
2. ChannelConfigStore - Stores channel definitions

After deployment, you need to:
1. Set the OCR configuration using 'configure-llo' command
2. Update your LLO job spec with the deployed contract addresses`,
	RunE: runDeploy,
}

var ConfigureCmd = &cobra.Command{
	Use:   "configure-llo",
	Short: "Configure the deployed LLO contracts with OCR settings and channel definitions",
	RunE:  runConfigure,
}

var (
	configuratorAddr       string
	channelConfigStoreAddr string
	nodeInfoFile           string
	ocrConfigFile          string
	channelDefsFile        string
	channelDefsSHA         string
)

func init() {
	// Deploy command flags
	DeployCmd.Flags().StringVar(&rpcURL, "rpc-url", "http://localhost:8545", "RPC URL for the target chain")
	DeployCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for deployer (or use PRIVATE_KEY env var)")
	DeployCmd.Flags().Int64Var(&chainID, "chain-id", 1337, "Chain ID")
	DeployCmd.Flags().StringVar(&outputFile, "output", "llo-contracts.json", "Output file for contract addresses")

	// Configure command flags
	ConfigureCmd.Flags().StringVar(&rpcURL, "rpc-url", "http://localhost:8545", "RPC URL for the target chain")
	ConfigureCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for deployer (or use PRIVATE_KEY env var)")
	ConfigureCmd.Flags().Int64Var(&chainID, "chain-id", 1337, "Chain ID")
	ConfigureCmd.Flags().Uint32Var(&donID, "don-id", 2, "DON ID for the Streams DON")
	ConfigureCmd.Flags().StringVar(&configuratorAddr, "configurator", "", "Configurator contract address")
	ConfigureCmd.Flags().StringVar(&channelConfigStoreAddr, "channel-config-store", "", "ChannelConfigStore contract address")
	ConfigureCmd.Flags().StringVar(&nodeInfoFile, "node-info", "", "JSON file with node info (signers, transmitters)")
	ConfigureCmd.Flags().StringVar(&ocrConfigFile, "ocr-config", "", "JSON file with OCR configuration")
	ConfigureCmd.Flags().StringVar(&channelDefsFile, "channel-definitions", "", "JSON file with channel definitions")
	ConfigureCmd.Flags().StringVar(&channelDefsSHA, "channel-definitions-sha", "", "SHA3-256 hash of channel definitions")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pk := privateKey
	if pk == "" {
		pk = os.Getenv("PRIVATE_KEY")
	}
	if pk == "" {
		// Default Anvil account #0 private key
		pk = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	}
	pk = strings.TrimPrefix(pk, "0x")

	contracts, err := Deploy(ctx, DeployConfig{
		RPCURL:     rpcURL,
		PrivateKey: pk,
		ChainID:    big.NewInt(chainID),
	})
	if err != nil {
		return fmt.Errorf("failed to deploy contracts: %w", err)
	}

	// Save contract addresses
	output := map[string]string{
		"configurator":       contracts.ConfiguratorAddress.Hex(),
		"channelConfigStore": contracts.ChannelConfigStoreAddr.Hex(),
		"chainId":            fmt.Sprintf("%d", chainID),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("\nContracts deployed successfully!\n")
	fmt.Printf("Output saved to: %s\n", outputFile)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Run: go run . configure-llo --configurator=%s --channel-config-store=%s --node-info=<node-info.json>\n",
		contracts.ConfiguratorAddress.Hex(), contracts.ChannelConfigStoreAddr.Hex())
	fmt.Printf("2. Update your LLO job spec with contractID=%s\n", contracts.ConfiguratorAddress.Hex())

	return nil
}

// NodeInfo represents the information needed from each node
type NodeInfo struct {
	OCRKeyBundleID string `json:"ocrKeyBundleID"`
	CSAPublicKey   string `json:"csaPublicKey"`
	SignerAddress  string `json:"signerAddress"` // derived from OCR key
}

// NodesConfig holds all node information for configuration
type NodesConfig struct {
	Nodes []NodeInfo `json:"nodes"`
}

// OCRConfigInput represents the full OCR configuration from the configure script
type OCRConfigInput struct {
	DonID                 uint32   `json:"donId"`
	Signers               []string `json:"signers"`               // On-chain signing addresses
	OffchainTransmitters  []string `json:"offchainTransmitters"`  // CSA public keys
	OffchainPublicKeys    []string `json:"offchainPublicKeys"`    // OCR offchain public keys
	ConfigPublicKeys      []string `json:"configPublicKeys"`      // Config public keys
	PeerIDs               []string `json:"peerIds"`               // P2P peer IDs
	F                     int      `json:"f"`                     // Fault tolerance
	OffchainConfigVersion int      `json:"offchainConfigVersion"` // Config version
}

func runConfigure(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pk := privateKey
	if pk == "" {
		pk = os.Getenv("PRIVATE_KEY")
	}
	if pk == "" {
		pk = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	}
	pk = strings.TrimPrefix(pk, "0x")

	if configuratorAddr == "" || channelConfigStoreAddr == "" {
		return fmt.Errorf("both --configurator and --channel-config-store are required")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}

	auth, _, err := GetTransactorFromPrivateKey(pk, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// For testing, we'll create a simple configuration
	// In production, this would come from node info

	// Create default channel definitions for testing
	channelDefs := llotypes.ChannelDefinitions{
		1: {
			ReportFormat: llotypes.ReportFormatJSON,
			Streams: []llotypes.Stream{
				{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
			},
		},
	}

	// Serialize channel definitions
	channelDefsJSON, err := json.Marshal(channelDefs)
	if err != nil {
		return fmt.Errorf("failed to marshal channel definitions: %w", err)
	}

	// Start a simple HTTP server to host channel definitions
	channelDefsURL, channelDefsSHA, stopServer := startChannelDefsServer(channelDefsJSON)
	defer stopServer()

	fmt.Printf("Channel definitions server started at: %s\n", channelDefsURL)
	fmt.Printf("Channel definitions SHA: %x\n", channelDefsSHA)

	// Create LLOContracts instance with existing contracts
	contracts, err := loadContracts(client, configuratorAddr, channelConfigStoreAddr)
	if err != nil {
		return err
	}

	// Set channel definitions
	err = contracts.SetChannelDefinitions(ctx, client, auth, donID, channelDefsURL, channelDefsSHA)
	if err != nil {
		return fmt.Errorf("failed to set channel definitions: %w", err)
	}

	// For OCR config, we need the OCR config file with all key information
	if ocrConfigFile != "" {
		err = configureOCR(ctx, client, auth, contracts, ocrConfigFile)
		if err != nil {
			return fmt.Errorf("failed to configure OCR: %w", err)
		}
	} else {
		fmt.Printf("\nWarning: No --ocr-config provided, skipping OCR configuration.\n")
		fmt.Printf("To complete setup, run the configure-llo-ocr.sh script first to gather key info.\n")
	}

	fmt.Printf("\nConfiguration complete!\n")
	
	// Don't block forever - the channel definitions server will be stopped when the process exits
	// For E2E testing, the nodes will have already fetched the channel definitions
	fmt.Printf("Channel definitions were served at: %s\n", channelDefsURL)
	fmt.Printf("Note: Server is shutting down. For long-running tests, use inline channelDefinitions in job spec.\n")
	
	return nil
}

func loadContracts(client *ethclient.Client, configuratorHex, channelConfigStoreHex string) (*LLOContracts, error) {
	confAddr := common.HexToAddress(configuratorHex)
	channelAddr := common.HexToAddress(channelConfigStoreHex)

	return LoadExistingContracts(client, confAddr, channelAddr)
}

func configureOCR(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	contracts *LLOContracts,
	ocrConfigPath string,
) error {
	data, err := os.ReadFile(ocrConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read OCR config file: %w", err)
	}

	var ocrInput OCRConfigInput
	if err := json.Unmarshal(data, &ocrInput); err != nil {
		return fmt.Errorf("failed to parse OCR config: %w", err)
	}

	n := len(ocrInput.Signers)
	if n < 4 {
		return fmt.Errorf("need at least 4 nodes for OCR, got %d", n)
	}

	// Build Oracle identities for the config helper
	oracles := make([]confighelper.OracleIdentityExtra, n)
	for i := 0; i < n; i++ {
		// Decode offchain public key
		offchainPKBytes, err := hex.DecodeString(ocrInput.OffchainPublicKeys[i])
		if err != nil {
			return fmt.Errorf("failed to decode offchain public key for node %d: %w", i, err)
		}
		var offchainPK types.OffchainPublicKey
		copy(offchainPK[:], offchainPKBytes)

		// Decode on-chain public key (signer)
		onchainPKBytes, err := hex.DecodeString(strings.TrimPrefix(ocrInput.Signers[i], "0x"))
		if err != nil {
			return fmt.Errorf("failed to decode signer for node %d: %w", i, err)
		}

		// Decode config public key
		configPKBytes, err := hex.DecodeString(ocrInput.ConfigPublicKeys[i])
		if err != nil {
			return fmt.Errorf("failed to decode config public key for node %d: %w", i, err)
		}
		var configPK types.ConfigEncryptionPublicKey
		copy(configPK[:], configPKBytes)

		// CSA key as transmitter
		csaBytes, err := hex.DecodeString(ocrInput.OffchainTransmitters[i])
		if err != nil {
			return fmt.Errorf("failed to decode transmitter for node %d: %w", i, err)
		}

		oracles[i] = confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPKBytes,
				OffchainPublicKey: offchainPK,
				PeerID:            ocrInput.PeerIDs[i],
				TransmitAccount:   types.Account(hex.EncodeToString(csaBytes)),
			},
			ConfigEncryptionPublicKey: configPK,
		}
	}

	// Create offchain config for LLO
	lloOffchainConfig := datastreamsllo.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: uint64(1 * time.Second),
	}
	lloOffchainConfigBytes, err := lloOffchainConfig.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode LLO offchain config: %w", err)
	}

	// Create the LLO onchain config (required by the Configurator contract)
	lloOnchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil, // nil = production config (no predecessor)
	})
	if err != nil {
		return fmt.Errorf("failed to encode LLO onchain config: %w", err)
	}

	f := ocrInput.F
	if f == 0 {
		f = (n - 1) / 3 // BFT tolerance
	}

	// Generate the full OCR3 config
	maxDurationInit := 0 * time.Second
	signers, transmitters, fVal, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,   // DeltaProgress
		20*time.Second,  // DeltaResend
		800*time.Millisecond, // DeltaInitial
		500*time.Millisecond, // DeltaRound
		200*time.Millisecond, // DeltaGrace
		1*time.Second,   // DeltaCertifiedCommitRequest
		10*time.Second,  // DeltaStage
		3,               // RMax
		[]int{n},        // S - number of oracles that sign
		oracles,
		lloOffchainConfigBytes, // ReportingPluginConfig
		&maxDurationInit,       // MaxDurationInitialization
		250*time.Millisecond,   // MaxDurationQuery
		250*time.Millisecond,   // MaxDurationObservation
		250*time.Millisecond,   // MaxDurationShouldAcceptAttestedReport
		250*time.Millisecond,   // MaxDurationShouldTransmitAcceptedReport
		f,
		lloOnchainConfig, // OnchainConfig - encoded LLO onchain config
	)
	if err != nil {
		return fmt.Errorf("failed to generate OCR3 config: %w", err)
	}

	fmt.Printf("Generated OCR config: f=%d, %d signers, %d transmitters\n", fVal, len(signers), len(transmitters))

	// Convert to the format expected by SetProductionConfig
	onchainPubKeys := make([][]byte, len(signers))
	for i, signer := range signers {
		onchainPubKeys[i] = signer
	}

	offchainTransmittersBytes := make([][32]byte, len(transmitters))
	for i, t := range transmitters {
		txBytes, err := hex.DecodeString(string(t))
		if err != nil {
			return fmt.Errorf("failed to decode transmitter %d: %w", i, err)
		}
		copy(offchainTransmittersBytes[i][:], txBytes)
	}

	ocrCfg := OCRConfig{
		DonID:                 ocrInput.DonID,
		Signers:               onchainPubKeys,
		Transmitters:          offchainTransmittersBytes,
		F:                     fVal,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return contracts.SetProductionConfig(ctx, client, auth, ocrCfg)
}

func startChannelDefsServer(channelDefs []byte) (url string, sha [32]byte, stop func()) {
	sha = sha256.Sum256(channelDefs)

	mux := http.NewServeMux()
	mux.HandleFunc("/channel-definitions.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(channelDefs)
	})

	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Channel definitions server error: %v\n", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	return "http://localhost:18080/channel-definitions.json", sha, func() {
		server.Shutdown(context.Background())
	}
}
