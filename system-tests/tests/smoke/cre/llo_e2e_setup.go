// Package cre provides helpers for CRE (Chainlink Runtime Environment) testing.
//
// This file provides LLO E2E setup including contract deployment and OCR configuration.
// It includes embedded deployment logic to avoid cross-module dependencies.
package cre

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	// Import data-streams-deploy changesets for LLO contract deployment
	configurator_v0_5_0 "github.com/smartcontractkit/data-streams-deploy/changeset/configurator/v0_5_0"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// LLOContracts holds the deployed contract addresses and instances
type LLOContracts struct {
	ConfiguratorAddress common.Address
	Configurator        *configurator.Configurator
}

// LLOInfrastructure holds all the deployed LLO infrastructure components
type LLOInfrastructure struct {
	Contracts *LLOContracts
	DonID     uint32
	// StopChannelDefsServer is a no-op since we use inline channel definitions
	StopChannelDefsServer func()
}

// SetupLLOInfrastructure deploys LLO contracts only.
// OCR configuration should be set AFTER LLO jobs are deployed using SetOCRConfiguration().
// This ensures LogPoller is running and can catch the ProductionConfigSet event.
func SetupLLOInfrastructure(
	t *testing.T,
	ctx context.Context,
	testLogger zerolog.Logger,
	testEnv *ttypes.TestEnvironment,
	donID uint32,
) (*LLOInfrastructure, error) {
	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        SETTING UP LLO INFRASTRUCTURE                                ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Get the first EVM blockchain
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain type")
	evmBlockchain := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)

	// Wait for Anvil to be ready before accessing it
	testLogger.Info().Msg("Waiting for Anvil to be ready...")
	err := waitForAnvil(ctx, evmBlockchain.SethClient.Client, testLogger, 30*time.Second)
	require.NoError(t, err, "Anvil did not become ready in time")

	// Get chain ID from the Seth client
	chainID, err := evmBlockchain.SethClient.Client.ChainID(ctx)
	require.NoError(t, err, "failed to get chain ID")
	testLogger.Info().Str("chainID", chainID.String()).Msg("Using blockchain")

	// Get chain selector from chain ID
	chainSelector, err := chainselectors.SelectorFromChainId(chainID.Uint64())
	require.NoError(t, err, "failed to get chain selector from chain ID")
	testLogger.Info().Uint64("chainSelector", chainSelector).Msg("Using chain selector")

	// Get CLD environment
	if testEnv.CreEnvironment.CldfEnvironment == nil {
		return nil, fmt.Errorf("CLDF environment is nil - cannot use changesets")
	}
	cldfEnv := testEnv.CreEnvironment.CldfEnvironment

	// Deploy LLO contracts using CLD changesets
	testLogger.Info().Msg("Deploying LLO contracts (Configurator) using CLD changesets...")
	contracts, err := deployLLOContractsWithChangesets(ctx, testLogger, cldfEnv, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy LLO contracts: %w", err)
	}
	testLogger.Info().
		Str("configurator", contracts.ConfiguratorAddress.Hex()).
		Msg("✓ LLO contracts deployed")

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║  ✓ LLO INFRASTRUCTURE SETUP COMPLETE                                ║")
	testLogger.Info().Msg("║  NOTE: Call SetOCRConfiguration() AFTER deploying LLO jobs!       ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")
	testLogger.Info().
		Str("configurator", contracts.ConfiguratorAddress.Hex()).
		Uint32("donID", donID).
		Msg("LLO Infrastructure Summary")

	return &LLOInfrastructure{
		Contracts:             contracts,
		DonID:                 donID,
		StopChannelDefsServer: func() {}, // No server with inline channel definitions
	}, nil
}

// SetOCRConfiguration sets the OCR configuration on the Configurator contract.
// This should be called AFTER LLO jobs are deployed so that LogPoller is running
// and can catch the ProductionConfigSet event.
func SetOCRConfiguration(
	t *testing.T,
	ctx context.Context,
	testLogger zerolog.Logger,
	testEnv *ttypes.TestEnvironment,
	infra *LLOInfrastructure,
) error {
	testLogger.Info().Msg("Setting OCR configuration on Configurator using CLD changesets...")

	// Get the first EVM blockchain
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain type")
	evmBlockchain := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)

	// Wait for Anvil to be ready before accessing it
	testLogger.Info().Msg("Waiting for Anvil to be ready...")
	err := waitForAnvil(ctx, evmBlockchain.SethClient.Client, testLogger, 30*time.Second)
	if err != nil {
		return fmt.Errorf("Anvil did not become ready: %w", err)
	}

	// Get chain ID and selector
	chainID, err := evmBlockchain.SethClient.Client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}
	chainSelector, err := chainselectors.SelectorFromChainId(chainID.Uint64())
	if err != nil {
		return fmt.Errorf("failed to get chain selector: %w", err)
	}

	// Get CLD environment
	if testEnv.CreEnvironment.CldfEnvironment == nil {
		return fmt.Errorf("CLDF environment is nil - cannot use changesets")
	}
	cldfEnv := testEnv.CreEnvironment.CldfEnvironment

	// Wait briefly for LogPoller filter to be registered by the LLO job
	testLogger.Info().Msg("Waiting 3s for LogPoller filter to be registered...")
	time.Sleep(3 * time.Second)

	// Set OCR configuration using changesets
	err = setOCRConfigurationWithChangesets(ctx, testLogger, testEnv, cldfEnv, infra.Contracts, chainSelector, infra.DonID)
	if err != nil {
		return fmt.Errorf("failed to set OCR configuration: %w", err)
	}
	testLogger.Info().Msg("✓ OCR configuration set")

	// Wait for LogPoller to index the event
	testLogger.Info().Msg("Waiting 10s for LogPoller to pick up the config event...")
	time.Sleep(10 * time.Second)

	// Mine some blocks to ensure confirmation
	err = mineBlocksAndWait(ctx, evmBlockchain.SethClient.URL, 5, testLogger)
	if err != nil {
		testLogger.Warn().Err(err).Msg("Failed to mine blocks - continuing anyway")
	}

	testLogger.Info().Msg("✓ OCR configuration complete")
	return nil
}

// deployLLOContractsWithChangesets deploys the Configurator contract using CLD changesets
func deployLLOContractsWithChangesets(
	ctx context.Context,
	testLogger zerolog.Logger,
	cldfEnv *cldf.Environment,
	chainSelector uint64,
) (*LLOContracts, error) {
	testLogger.Info().Msg("Deploying Configurator contract using data-streams-deploy changeset...")

	// Deploy Configurator contract using CLD changeset
	configuratorOutput, err := configurator_v0_5_0.DeployConfiguratorChangeset.Apply(*cldfEnv, configurator_v0_5_0.DeployConfiguratorConfig{
		ChainsToDeploy: []uint64{chainSelector},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Configurator: %w", err)
	}

	// Extract Configurator address from datastore
	configuratorAddr, err := extractContractAddress(configuratorOutput.DataStore, chainSelector, "Configurator")
	if err != nil {
		return nil, fmt.Errorf("failed to get Configurator address: %w", err)
	}
	testLogger.Info().Str("address", configuratorAddr.Hex()).Msg("Configurator deployed")

	// Merge the changeset output datastore into the CLD environment's datastore
	// This is required so that subsequent changesets can find the deployed contract
	if configuratorOutput.DataStore != nil {
		mergedDS := ds.NewMemoryDataStore()
		// Merge the changeset output datastore
		if err := mergedDS.Merge(configuratorOutput.DataStore.Seal()); err != nil {
			return nil, errors.Wrap(err, "failed to merge changeset datastore")
		}
		// Merge the existing CLD environment datastore
		if cldfEnv.DataStore != nil {
			if err := mergedDS.Merge(cldfEnv.DataStore); err != nil {
				return nil, errors.Wrap(err, "failed to merge existing datastore")
			}
		}
		// Update the CLD environment with the merged datastore
		cldfEnv.DataStore = mergedDS.Seal()
		testLogger.Debug().Msg("Merged Configurator deployment datastore into CLD environment")
	}

	// Get the contract instance for later use (if needed)
	// Note: We don't need the contract instance for changesets, but keeping it for compatibility
	// We can get the client from the chain provider if needed, but for now we'll just store the address
	// The contract instance can be created later if needed using the address and a client
	return &LLOContracts{
		ConfiguratorAddress: configuratorAddr,
		Configurator:        nil, // Will be created on-demand if needed
	}, nil
}

// getCapabilitiesDON returns the DON that has type "capabilities" (streams-trigger / LLO).
// It uses the DON type flag so topology is resolved reliably regardless of DON name.
// Returns a clear error listing available DONs if the capabilities DON is not in the topology.
func getCapabilitiesDON(testEnv *ttypes.TestEnvironment) (*cre.Don, error) {
	caps := testEnv.Dons.DonsWithFlag(cre.CapabilitiesDON)
	if len(caps) > 0 {
		return caps[0], nil
	}
	names := make([]string, 0, len(testEnv.Dons.List()))
	for _, d := range testEnv.Dons.List() {
		names = append(names, d.Metadata().Name)
	}
	return nil, fmt.Errorf("capabilities DON not found (topology has DONs: %v). Ensure the environment was started with workflow-capabilities-llo-don.toml", names)
}

// setOCRConfigurationWithChangesets gathers node keys and sets the OCR configuration on the Configurator using CLD changesets
func setOCRConfigurationWithChangesets(
	ctx context.Context,
	testLogger zerolog.Logger,
	testEnv *ttypes.TestEnvironment,
	cldfEnv *cldf.Environment,
	contracts *LLOContracts,
	chainSelector uint64,
	donID uint32,
) error {
	capabilitiesDON, err := getCapabilitiesDON(testEnv)
	if err != nil {
		return err
	}

	// Get worker nodes (exclude bootstrap)
	workers, err := capabilitiesDON.Workers()
	if err != nil {
		return fmt.Errorf("failed to get worker nodes: %w", err)
	}

	testLogger.Info().Int("workers", len(workers)).Msg("Found worker nodes in capabilities DON")

	// Build oracle identities from node keys
	oracles := make([]confighelper.OracleIdentityExtra, len(workers))
	for i, node := range workers {
		// Get CSA key
		if node.Keys.CSAKey == nil {
			return fmt.Errorf("node %s has no CSA key", node.Name)
		}
		csaKeyHex := strings.TrimPrefix(node.Keys.CSAKey.Key, "csa_")
		csaBytes, err := hex.DecodeString(csaKeyHex)
		if err != nil {
			return fmt.Errorf("failed to decode CSA key for node %s: %w", node.Name, err)
		}

		// Get peer ID - must NOT have the p2p_ prefix for OracleIdentity
		peerID := strings.TrimPrefix(node.Keys.P2PKey.PeerID.String(), "p2p_")

		// Fetch real OCR2 keys from the node via REST API
		ocr2Keys, err := node.Clients.RestClient.MustReadOCR2Keys()
		if err != nil {
			return fmt.Errorf("failed to read OCR2 keys from node %s: %w", node.Name, err)
		}

		// Find the EVM OCR2 key bundle
		var offchainPK types.OffchainPublicKey
		var configPK types.ConfigEncryptionPublicKey
		var onchainPK []byte
		foundEVMKey := false
		for _, keyData := range ocr2Keys.Data {
			if keyData.Attributes.ChainType == "evm" {
				// Extract keys, stripping prefixes
				offchainHex := strings.TrimPrefix(keyData.Attributes.OffChainPublicKey, "ocr2off_evm_")
				configHex := strings.TrimPrefix(keyData.Attributes.ConfigPublicKey, "ocr2cfg_evm_")
				onchainHex := strings.TrimPrefix(keyData.Attributes.OnChainPublicKey, "ocr2on_evm_")

				offchainBytes, err := hex.DecodeString(offchainHex)
				if err != nil {
					return fmt.Errorf("failed to decode offchain key for node %s: %w", node.Name, err)
				}
				copy(offchainPK[:], offchainBytes)

				configBytes, err := hex.DecodeString(configHex)
				if err != nil {
					return fmt.Errorf("failed to decode config key for node %s: %w", node.Name, err)
				}
				copy(configPK[:], configBytes)

				onchainPK, err = hex.DecodeString(onchainHex)
				if err != nil {
					return fmt.Errorf("failed to decode onchain key for node %s: %w", node.Name, err)
				}

				foundEVMKey = true
				testLogger.Debug().
					Str("node", node.Name).
					Str("offchainKey", offchainHex[:16]+"...").
					Str("configKey", configHex[:16]+"...").
					Str("onchainKey", onchainHex).
					Msg("Found EVM OCR2 key bundle")
				break
			}
		}
		if !foundEVMKey {
			return fmt.Errorf("no EVM OCR2 key bundle found for node %s", node.Name)
		}

		// Format CSA key as 0x-prefixed hex string (66 characters total: 0x + 64 hex chars)
		// The changeset expects TransmitAccount in format "0x<64-hex-chars>"
		csaKeyFormatted := "0x" + hex.EncodeToString(csaBytes)
		oracles[i] = confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPK,
				OffchainPublicKey: offchainPK,
				PeerID:            peerID,
				TransmitAccount:   types.Account(csaKeyFormatted),
			},
			ConfigEncryptionPublicKey: configPK,
		}

		testLogger.Debug().
			Str("node", node.Name).
			Str("peerID", peerID).
			Str("csaKey", csaKeyHex[:16]+"...").
			Msg("Gathered node identity with real OCR2 keys")
	}

	// Calculate f (fault tolerance)
	n := len(oracles)
	f := (n - 1) / 3
	if f < 1 {
		f = 1
	}

	// Build configuration using the builder pattern from data-streams-deploy
	// The changeset builder handles OCR3 config generation internally
	// DON config ID must be a bytes32 hex string (0x + 64 hex chars)
	var donConfigIDBytes [32]byte
	big.NewInt(int64(donID)).FillBytes(donConfigIDBytes[:])
	donConfigID := "0x" + hex.EncodeToString(donConfigIDBytes[:])
	fInt := int(f) // Convert uint8 to int for changeset
	configParams := configurator_v0_5_0.NewConfiguratorConfig(configurator_v0_5_0.ConfiguratorSetParamsOptions{
		DONConfigID:         &donConfigID,
		ConfiguratorAddress: &contracts.ConfiguratorAddress,
		OCROptions: &configurator_v0_5_0.OCR3DataStreamsOptions{
			Oracles: oracles,
			F:       &fInt,
			S:       []int{len(oracles)},
			OnchainConfigOptions: &datastreamsllo.OnchainConfig{
				Version:                 1,
				PredecessorConfigDigest: nil,
			},
			OffchainConfigOptions: datastreamsllo.OffchainConfig{
				ProtocolVersion:                     1,
				DefaultMinReportIntervalNanoseconds: uint64(1 * time.Second),
				EnableObservationCompression:        true,
			},
		},
	})

	testLogger.Info().
		Int("oracles", len(oracles)).
		Int("f", fInt).
		Msg("Built OCR3 configuration for changeset")

	// Set production config using CLD changeset
	_, err = configurator_v0_5_0.SetProductionConfigChangeset.Apply(*cldfEnv, configurator_v0_5_0.SetProductionConfig{
		ConfigurationsByChain: map[uint64][]configurator_v0_5_0.ConfiguratorSetParams{
			chainSelector: {*configParams},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set production config: %w", err)
	}

	testLogger.Info().Msg("✓ Production config set using CLD changeset")
	return nil
}

// waitForAnvil waits for Anvil to be ready by attempting to get the chain ID
func waitForAnvil(ctx context.Context, client interface {
	ChainID(context.Context) (*big.Int, error)
}, logger zerolog.Logger, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		_, err := client.ChainID(ctx)
		if err == nil {
			logger.Info().Msg("Anvil is ready")
			return nil
		}
		logger.Debug().Err(err).Msg("Anvil not ready yet, retrying...")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Continue loop
		}
	}
	return fmt.Errorf("Anvil did not become ready within %v", timeout)
}

// mineBlocksAndWait mines a specified number of blocks on Anvil and waits for them
func mineBlocksAndWait(ctx context.Context, rpcURL string, numBlocks int, logger zerolog.Logger) error {
	// Connect directly to Anvil RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	for i := 0; i < numBlocks; i++ {
		// Send a dummy self-transaction to trigger block mining on Anvil
		// Using anvil_mine RPC call for faster mining
		var result interface{}
		err := client.Client().CallContext(ctx, &result, "anvil_mine", 1)
		if err != nil {
			// Fallback: try evm_mine
			err = client.Client().CallContext(ctx, &result, "evm_mine")
			if err != nil {
				return fmt.Errorf("failed to mine block %d: %w", i+1, err)
			}
		}
	}
	logger.Debug().Int("blocks", numBlocks).Msg("Mined additional blocks")
	return nil
}

// DeployStreamJobs deploys stream jobs to the capabilities DON nodes via Job Distributor
// Stream jobs use hardcoded values (via memo task) to avoid bridge connectivity issues
func DeployStreamJobs(
	ctx context.Context,
	testLogger zerolog.Logger,
	testEnv *ttypes.TestEnvironment,
) error {
	testLogger.Info().Msg("Deploying stream jobs via Job Distributor...")

	capabilitiesDON, err := getCapabilitiesDON(testEnv)
	if err != nil {
		return err
	}

	// Stream job specs for each stream - using hardcoded values
	// Stream 1: For ReportFormat 5 (CapabilityTrigger) - TEST/USD
	// Streams 2,3,4: For ReportFormat 7 (requires at least 3 streams) - NATIVE/USD, LINK/USD, DATA/USD
	streamJobs := []struct {
		streamID uint32
		name     string
		pair     string
	}{
		{1, "stream-test-usd", "TEST/USD"},
		{2, "stream-native-usd", "NATIVE/USD"},
		{3, "stream-link-usd", "LINK/USD"},
		{4, "stream-data-usd", "DATA/USD"},
	}

	workers, err := capabilitiesDON.Workers()
	if err != nil {
		return fmt.Errorf("failed to get worker nodes: %w", err)
	}

	var jobSpecs cre.DonJobs

	for _, node := range workers {
		for _, sj := range streamJobs {
			// Generate unique external job ID using a proper UUID
			externalJobID := uuid.New().String()

			// Make job name unique per node to avoid conflicts
			// Format: stream-{pair}-{node-name} (e.g., "stream-test-usd-capabilities-node0")
			uniqueJobName := fmt.Sprintf("%s-%s", sj.name, node.Name)

			// Build job spec with hardcoded values (bridgeName parameter is ignored by BuildStreamJobSpec)
			// Note: streamID must be a top-level field, not under [streamSpec]
			jobSpec := BuildStreamJobSpec(uniqueJobName, sj.streamID, "", externalJobID)

			jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
				NodeId: node.JobDistributorDetails.NodeID,
				Spec:   jobSpec,
			})

			testLogger.Debug().
				Str("node", node.Name).
				Str("nodeId", node.JobDistributorDetails.NodeID).
				Uint32("streamID", sj.streamID).
				Str("jobName", uniqueJobName).
				Msg("Prepared stream job spec with hardcoded values")
		}
	}

	// Create jobs via Job Distributor
	if testEnv.CreEnvironment.CldfEnvironment == nil {
		return fmt.Errorf("CLDF environment is nil - cannot deploy jobs via Job Distributor")
	}

	err = createJobs(ctx, testLogger, testEnv.CreEnvironment.CldfEnvironment.Offchain, testEnv.Dons, jobSpecs)
	if err != nil {
		return fmt.Errorf("failed to create stream jobs: %w", err)
	}

	testLogger.Info().Int("count", len(jobSpecs)).Msg("✓ Stream jobs deployed")
	return nil
}

// DeployLLOJobs deploys LLO OCR2 jobs with CRE transmitter to the capabilities DON via Job Distributor
// This uses inline channelDefinitions (most reliable) instead of fetching from a URL.
// IMPORTANT: This also deploys a bootstrap job for LLO to the bootstrap node - this is essential
// for the OCR protocol to start. Without the bootstrap job, the Blue instance won't run.
func DeployLLOJobs(
	ctx context.Context,
	testLogger zerolog.Logger,
	testEnv *ttypes.TestEnvironment,
	infra *LLOInfrastructure,
) error {
	testLogger.Info().Msg("Deploying LLO jobs with CRE transmitter via Job Distributor...")

	capabilitiesDON, err := getCapabilitiesDON(testEnv)
	if err != nil {
		return err
	}

	workers, err := capabilitiesDON.Workers()
	if err != nil {
		return fmt.Errorf("failed to get worker nodes: %w", err)
	}

	// Get bootstrap peer info from the global bootstrap node (shared across all DONs)
	bootstrap, found := testEnv.Dons.Bootstrap()
	if !found {
		return fmt.Errorf("bootstrap node not found in topology")
	}
	_, ocrPeeringCfg, peerErr := cre.PeeringCfgs(bootstrap)
	if peerErr != nil {
		return fmt.Errorf("failed to get peering configs: %w", peerErr)
	}
	bootstrapAddr := fmt.Sprintf("%s@%s:%d", ocrPeeringCfg.OCRBootstraperPeerID, ocrPeeringCfg.OCRBootstraperHost, ocrPeeringCfg.Port)

	var jobSpecs cre.DonJobs

	// Step 1: Deploy bootstrap job for LLO to the bootstrap node
	// This is CRITICAL - without the bootstrap job, the LLO OCR protocol won't start.
	// The bootstrap job provides the P2P networking foundation for the LLO network.
	bootstrapJobID := uuid.New().String()
	bootstrapJobSpec := fmt.Sprintf(`type = "bootstrap"
schemaVersion = 1
name = "llo-bootstrap"
externalJobID = "%s"
contractID = "%s"
contractConfigTrackerPollInterval = "1s"
relay = "evm"

[relayConfig]
chainID = "1337"
fromBlock = 0
lloDonID = %d
lloConfigMode = "bluegreen"
providerType = "llo"
`,
		bootstrapJobID,
		infra.Contracts.ConfiguratorAddress.Hex(),
		infra.DonID,
	)

	jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
		NodeId: bootstrap.JobDistributorDetails.NodeID,
		Spec:   bootstrapJobSpec,
	})
	testLogger.Info().
		Str("nodeId", bootstrap.JobDistributorDetails.NodeID).
		Str("contractID", infra.Contracts.ConfiguratorAddress.Hex()).
		Msg("Prepared LLO bootstrap job spec")

	// Step 2: Deploy worker LLO jobs
	for _, node := range workers {
		ocrKeyBundleID := node.Keys.OCR2BundleIDs["evm"]
		// Generate unique external job ID using a proper UUID
		externalJobID := uuid.New().String()

		// Get CSA public key for transmitterID (strip the "csa_" prefix)
		csaPubKeyHex := strings.TrimPrefix(node.Keys.CSAKey.Key, "csa_")

		// Use inline channelDefinitions - much more reliable than fetching from URL
		// Channel 1: ReportFormat 5 (CapabilityTrigger) - Stream 1 (TEST/USD) → MAGIC 424242
		// Channel 2: ReportFormat 7 (EVMABIEncodeUnpackedExpr) - Stream 4 (DATA/USD) base value 111111
		//   Calculated stream 5 multiplies stream 4 by 5: Mul(s4, 5) → 555555
		// Job spec matches the working version from commit 73215a5
		jobSpec := fmt.Sprintf(`type = "offchainreporting2"
schemaVersion = 1
name = "llo-streams-don"
externalJobID = "%s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%s"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%s"
p2pv2Bootstrappers = ["%s"]
relay = "evm"
pluginType = "llo"
transmitterID = "%s"

[pluginConfig]
# Channel Definitions for E2E test with both report formats
# Note: ReportFormat 7 requires at least 3 streams
# Stream 4 has base value 111111, calculated stream 5 multiplies by 5 to get 555555
channelDefinitions = """
{
  "1": {
    "reportFormat": 5,
    "streams": [{"streamId": 1, "aggregator": "median"}],
    "opts": {}
  },
  "2": {
    "reportFormat": 7,
    "streams": [
      {"streamId": 2, "aggregator": "median"},
      {"streamId": 3, "aggregator": "median"},
      {"streamId": 4, "aggregator": "median"},
      {"streamId": 5, "aggregator": "calculated"}
    ],
    "opts": {
      "feedId": "0x0001000000000000000000000000000000000000000000000000000000000001",
      "baseUSDFee": "0.1",
      "expirationWindow": 3600,
      "abi": [
        {"type": "int192", "expression": "Mul(s4, 5)", "expressionStreamId": 5}
      ]
    }
  }
}
"""
donID = %d

[[pluginConfig.transmitters]]
type = "cre"
[pluginConfig.transmitters.opts]
triggerCapabilityName = "streams-trigger"
triggerCapabilityVersion = "2.0.0"
triggerTickerMinResolutionMs = 1000
triggerSendChannelBufferSize = 1000
# Top-of-window delay: sends wait until next wall-clock boundary (non-zero exercises delayed-send path)
transmissionWindowMs = 100

[relayConfig]
chainID = "1337"
lloDonID = %d
lloConfigMode = "bluegreen"
fromBlock = 1
`,
			externalJobID,
			infra.Contracts.ConfiguratorAddress.Hex(),
			ocrKeyBundleID,
			bootstrapAddr,
			csaPubKeyHex,
			infra.DonID,
			infra.DonID,
		)

		jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
			NodeId: node.JobDistributorDetails.NodeID,
			Spec:   jobSpec,
		})

		testLogger.Debug().
			Str("node", node.Name).
			Str("nodeId", node.JobDistributorDetails.NodeID).
			Str("ocrKeyBundleID", ocrKeyBundleID).
			Msg("Prepared LLO job spec")
	}

	// Create jobs via Job Distributor
	if testEnv.CreEnvironment.CldfEnvironment == nil {
		return fmt.Errorf("CLDF environment is nil - cannot deploy jobs via Job Distributor")
	}

	err = createJobs(ctx, testLogger, testEnv.CreEnvironment.CldfEnvironment.Offchain, testEnv.Dons, jobSpecs)
	if err != nil {
		return fmt.Errorf("failed to create LLO jobs: %w", err)
	}

	testLogger.Info().Int("count", len(jobSpecs)).Msg("✓ LLO jobs deployed")
	return nil
}

// createJobs creates jobs via Job Distributor with rate limiting and auto-approval
func createJobs(ctx context.Context, testLogger zerolog.Logger, offChainClient cldf_offchain.Client, dons *cre.Dons, jobSpecs cre.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	for _, jobReq := range jobSpecs {
		timeout := time.Second * 60
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

		_, err := offChainClient.ProposeJob(ctxWithTimeout, jobReq)
		cancel()
		if err != nil {
			return fmt.Errorf("failed to propose job for node %s: %w", jobReq.NodeId, err)
		}

		// Find the node and accept the job
		for _, don := range dons.List() {
			for _, node := range don.Nodes {
				if node.JobDistributorDetails.NodeID != jobReq.NodeId {
					continue
				}

				if err := node.AcceptJob(ctx, jobReq.Spec); err != nil {
					// Workflow specs get auto approved, so this error is expected
					if strings.Contains(err.Error(), "cannot approve an approved spec") {
						continue
					}
					// Handle foreign key constraint error - this happens when trying to delete an old job
					// that still has job proposals referencing it. This is usually from a previous test run
					// that didn't clean up properly. The job might already be in the desired state.
					// Using FRESH_ENV=true ensures a clean database and avoids this issue.
					if strings.Contains(err.Error(), "violates foreign key constraint") ||
						strings.Contains(err.Error(), "job_proposals_job_id_fkey") {
						// Log a warning but continue - the job proposal might already be approved
						// and the job might already be running, which is fine for our test
						testLogger.Warn().
							Str("node", node.Name).
							Str("error", err.Error()).
							Msg("Foreign key constraint error when accepting job - likely leftover state from previous run. Continuing...")
						continue
					}
					// Handle duplicate job name error - this happens when a job with the same name
					// already exists. This is usually from a previous test run that didn't clean up properly.
					// The job might already be in the desired state, so we can continue.
					if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
						strings.Contains(err.Error(), "idx_jobs_name") {
						testLogger.Warn().
							Str("node", node.Name).
							Str("error", err.Error()).
							Msg("Duplicate job name error when accepting job - likely leftover state from previous run. Continuing...")
						continue
					}
					return fmt.Errorf("failed to accept job for node %s: %w", node.Name, err)
				}
			}
		}

		// Small delay between jobs to avoid overwhelming the system
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
