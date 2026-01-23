package cre

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"

	mockcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mockpb "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	llo_consumer_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/llo_streams_trigger/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// LLO Streams Trigger config path - uses the capabilities DON topology with LLO feature
// Uses a test-specific config without gateway DON and with streams-trigger capability
const lloStreamsConfigPath = "/configs/workflow-capabilities-llo-don.toml"

// To run with debug log level (shows nSubscribers and other debug logs), set:
// CTF_LOG_LEVEL=debug go test -timeout 10m -run "Test_CRE_V2_LLO_Streams_Trigger_E2E" ./smoke/cre/...

// Test_CRE_V2_LLO_Streams_Trigger_Mock runs the LLO streams-trigger test using mock capability injection
// This test injects mock LLO reports via the mock capability controller
//
// REQUIREMENTS:
// This test requires the mock capability gRPC ports (15002-15005) to be exposed
// on the capabilities DON nodes. This requires topology configuration changes.
//
// To run:
//
//	cd system-tests/tests
//	go test -timeout 15m -run "Test_CRE_V2_LLO_Streams_Trigger_Mock" ./smoke/cre/...
//
// Set LLO_MOCK_ENABLED=true to run this test
func Test_CRE_V2_LLO_Streams_Trigger_Mock(t *testing.T) {
	if os.Getenv("LLO_MOCK_ENABLED") != "true" {
		t.Skip("Skipping LLO Mock test - requires mock capability ports exposed. Set LLO_MOCK_ENABLED=true to run.")
	}

	topology := os.Getenv("TOPOLOGY_NAME")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, lloStreamsConfigPath), v2RegistriesFlags...)

	t.Run("[v2] LLO Streams Trigger (Mock) - "+topology, func(t *testing.T) {
		ExecuteLLOStreamsTriggerTest(t, testEnv)
	})
}

// Test_CRE_V2_LLO_Streams_Trigger_E2E runs the full E2E test with real LLO plugin
// This test automatically deploys the full LLO infrastructure:
// 1. LLO contracts (Configurator, ChannelConfigStore)
// 2. OCR configuration with proper encryption keys
// 3. Stream jobs to fetch data from fake price provider
// 4. LLO jobs with CRE transmitter
// 5. Channel definitions for Format 5 and Format 7
//
// The test verifies end-to-end data flow from price provider -> streams DON -> workflow
// by checking for LLO_E2E_VALUE messages in workflow logs.
//
// To run:
//
//	cd system-tests/tests
//	go test -timeout 15m -run "Test_CRE_V2_LLO_Streams_Trigger_E2E" ./smoke/cre/...
func Test_CRE_V2_LLO_Streams_Trigger_E2E(t *testing.T) {
	topology := os.Getenv("TOPOLOGY_NAME")

	// IMPORTANT: Start the fake price provider BEFORE environment setup
	// The LLO feature's PostEnvStartup hook sets the channel definitions URL on the contract,
	// and the Docker containers need to be able to reach this URL immediately.
	testLogger := framework.L
	testLogger.Info().Msg("Setting up LLO price provider (BEFORE environment setup)...")
	priceProviderURL, err := SetupLLOPriceProvider(testLogger, &fake.Input{Port: 8171}, DefaultLLOPriceConfig())
	require.NoError(t, err, "Failed to set up LLO price provider")
	testLogger.Info().Str("url", priceProviderURL).Msg("LLO price provider ready")

	// Set environment variables for the LLO feature to pick up
	channelDefsURL := GetLLOProviderChannelDefsURL()
	t.Setenv("LLO_MOCK_EA_URL", priceProviderURL)
	t.Setenv("LLO_CHANNEL_DEFS_URL", channelDefsURL)
	// Enable OnchainPublicKey signer extraction for LLO tests
	// This ensures the capability registry uses OnchainPublicKey addresses for signature verification
	// The override in keystone_llo.go checks this environment variable
	t.Setenv("USE_LLO_ONCHAIN_SIGNER", "true")

	// Now set up the test environment - the fake provider is already running
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, lloStreamsConfigPath), v2RegistriesFlags...)

	t.Run("[v2] LLO Streams Trigger E2E - "+topology, func(t *testing.T) {
		ExecuteLLOStreamsTriggerE2EWithFullLLO(t, testEnv)
	})
}

// ExecuteLLOStreamsTriggerTest runs the full E2E test for streams-trigger@2.0.0
// This test verifies:
// 1. The streams-trigger capability is registered and exposed by the Streams DON
// 2. The Workflow DON can discover and subscribe to the capability
// 3. LLO reports are correctly routed from Streams DON to Workflow DON
// 4. The consumer workflow receives and processes the reports
func ExecuteLLOStreamsTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "llo_consumer/main.go"
	workflowName := "llo-consumer"

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        LLO STREAMS TRIGGER E2E TEST                                  ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Start Beholder listener to capture workflow logs
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Step 1: Deploy the LLO consumer workflow
	testLogger.Info().Msg("Step 1: Deploying LLO consumer workflow...")
	workflowConfig := llo_consumer_config.LLOConsumerConfig{
		StreamIDs:      []uint32{1}, // Subscribe to stream ID 1
		MaxFrequencyMs: 1000,        // 1 second max frequency
	}
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
	testLogger.Info().Msg("✓ Workflow deployed successfully")

	// Step 2: Wait for capability registration and DON-to-DON discovery
	testLogger.Info().Msg("Step 2: Waiting for capability discovery...")
	time.Sleep(10 * time.Second) // Allow time for capability registration
	testLogger.Info().Msg("✓ Capability discovery period completed")

	// Step 3: Inject mock LLO reports via the mock capability controller
	testLogger.Info().Msg("Step 3: Injecting mock LLO reports...")
	err := injectMockLLOReports(t, testEnv, 5) // Inject 5 reports
	require.NoError(t, err, "Failed to inject mock LLO reports")
	testLogger.Info().Msg("✓ Mock reports injected")

	// Step 4: Verify workflow received the reports
	testLogger.Info().Msg("Step 4: Verifying workflow received reports...")
	expectedLog := "LLO_E2E_VALUE"
	timeout := 2 * time.Minute
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - workflow did not receive reports")

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║  ✓ LLO STREAMS TRIGGER E2E TEST PASSED                              ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")
}

// injectMockLLOReports sends mock LLO reports via the mock capability controller
func injectMockLLOReports(t *testing.T, testEnv *ttypes.TestEnvironment, count int) error {
	// Get mock capability controller addresses from environment
	// The mock capability is exposed on ports 5002 (gRPC) on each Streams DON node
	mockAddresses := getMockCapabilityAddresses(testEnv)
	if len(mockAddresses) == 0 {
		return fmt.Errorf("no mock capability addresses found")
	}

	// Create mock capability controller
	controller := mockcap.NewMockCapabilityController(framework.L)
	err := controller.ConnectAll(mockAddresses, true, false)
	if err != nil {
		return fmt.Errorf("failed to connect to mock capability controllers: %w", err)
	}
	// Note: gRPC clients don't need explicit cleanup in Go

	// Send mock reports
	for i := 0; i < count; i++ {
		report, err := generateMockLLOReport(uint64(i + 1))
		if err != nil {
			return fmt.Errorf("failed to generate mock report %d: %w", i, err)
		}

		// Marshal the streams.Report into an anypb.Any
		anyReport, err := anypb.New(report)
		if err != nil {
			return fmt.Errorf("failed to marshal report to Any: %w", err)
		}

		req := &mockpb.SendTriggerEventRequest{
			TriggerID:   "streams-trigger@2.0.0",
			TriggerType: "llo",
			ID:          fmt.Sprintf("mock-llo-report-%d", i+1),
			Payload:     anyReport,
		}

		err = controller.SendTrigger(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to send trigger event %d: %w", i, err)
		}

		framework.L.Info().Msgf("Sent mock LLO report %d/%d", i+1, count)
		time.Sleep(500 * time.Millisecond) // Space out reports
	}

	return nil
}

// generateMockLLOReport creates a mock LLO report that matches the streams.Report structure
func generateMockLLOReport(seqNr uint64) (*streams.Report, error) {
	timestamp := uint64(time.Now().UnixNano())
	configDigest, _ := hex.DecodeString("00091599c39d29821b4949b9ba237d2d1d9b7369087a71283c921034898320b0")

	// Create OCR trigger report (the inner report payload)
	// Note: OCRTriggerReport only has EventID, Timestamp, and Outputs fields
	ocrReport := &capabilitiespb.OCRTriggerReport{
		EventID:   fmt.Sprintf("streams_1_%d_f5", timestamp), // Format: streams_DONID_TIMESTAMP_f5
		Timestamp: timestamp,
	}

	ocrReportBytes, err := proto.Marshal(ocrReport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OCR report: %w", err)
	}

	// Create the streams.Report wrapper
	return &streams.Report{
		ConfigDigest: configDigest,
		SeqNr:        seqNr,
		Report:       ocrReportBytes,
		Sigs: []*streams.OCRSignature{
			{Signer: 0, Signature: []byte("mock_signature_0")},
			{Signer: 1, Signature: []byte("mock_signature_1")},
		},
	}, nil
}

// getMockCapabilityAddresses returns the addresses of mock capability controllers
// TODO: Extract addresses from testEnv.NodeSets once mock capability port mapping is available
func getMockCapabilityAddresses(_ *ttypes.TestEnvironment) []string {
	// In the DON-to-DON setup, mock capabilities are exposed on port 5002 (gRPC)
	// Default: 4 nodes with ports 15002-15005 mapped to host
	const (
		numNodes = 4
		basePort = 15002
	)

	addresses := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		addresses[i] = fmt.Sprintf("localhost:%d", basePort+i)
	}
	return addresses
}

// ExecuteLLOStreamsTriggerE2EWithFullLLO runs the E2E test with actual LLO plugin
// This test deploys LLO infrastructure:
// - LLO contracts (Configurator, ChannelConfigStore)
// - OCR configuration with proper encryption keys
// - Channel definitions for Format 5 and Format 7
// - Stream jobs fetching from fake price provider
// - LLO jobs with CRE transmitter
//
// The test verifies end-to-end data flow by checking for magic numbers in workflow logs.
func ExecuteLLOStreamsTriggerE2EWithFullLLO(t *testing.T, testEnv *ttypes.TestEnvironment) {
	ctx := context.Background()
	testLogger := framework.L
	workflowFileLocation := "llo_consumer/main.go"
	workflowName := "llo-consumer-e2e"

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        LLO STREAMS TRIGGER E2E TEST (FULL LLO)                       ║")
	testLogger.Info().Msg("║  Magic Numbers: Format5=424242, Format7=555555                       ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Step 1: Setup LLO Infrastructure (deploy contracts only)
	testLogger.Info().Msg("Step 1: Setting up LLO Infrastructure (contracts only)...")
	const donID uint32 = 2 // Use DON ID 2 for LLO
	infra, err := SetupLLOInfrastructure(t, ctx, testLogger, testEnv, donID)
	require.NoError(t, err, "Failed to setup LLO infrastructure")
	defer infra.StopChannelDefsServer()
	testLogger.Info().Msg("✓ LLO Infrastructure ready")

	// Step 2: Deploy stream jobs
	testLogger.Info().Msg("Step 2: Deploying stream jobs...")
	err = DeployStreamJobs(ctx, testLogger, testEnv)
	require.NoError(t, err, "Failed to deploy stream jobs")
	testLogger.Info().Msg("✓ Stream jobs deployed")

	// Step 3: Deploy LLO jobs with CRE transmitter (this starts LogPoller)
	testLogger.Info().Msg("Step 3: Deploying LLO jobs with CRE transmitter...")
	err = DeployLLOJobs(ctx, testLogger, testEnv, infra)
	require.NoError(t, err, "Failed to deploy LLO jobs")
	testLogger.Info().Msg("✓ LLO jobs deployed")

	// Step 4: Set OCR config AFTER LLO jobs are deployed
	// This is critical because LogPoller only starts indexing when the LLO job
	// is deployed. Setting config now ensures the ProductionConfigSet event is
	// at a block that the LogPoller can see.
	testLogger.Info().Msg("Step 4: Setting OCR configuration...")
	err = SetOCRConfiguration(t, ctx, testLogger, testEnv, infra)
	require.NoError(t, err, "Failed to set OCR configuration")
	testLogger.Info().Msg("✓ OCR configuration set")

	// Step 5: Wait for LLO jobs to detect config and start OCR rounds
	testLogger.Info().Msg("Step 5: Waiting for LLO jobs to detect config and start OCR rounds...")
	time.Sleep(30 * time.Second) // Give LLO jobs time to detect config and start OCR rounds

	// Step 6: Wait for CRE Transmitter to register
	testLogger.Info().Msg("Step 6: Waiting for CRE Transmitter to register...")
	time.Sleep(10 * time.Second)

	// Step 7: Wait for capability discovery BEFORE deploying workflow
	// The workflow engine initialization happens immediately when the workflow is deployed,
	// so we need to ensure the streams-trigger capability is registered and discoverable first.
	testLogger.Info().Msg("Step 7: Waiting for capability discovery (streams-trigger must be registered)...")
	time.Sleep(40 * time.Second) // Allow time for cross-DON capability discovery

	// Start Beholder listener to capture workflow logs
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Step 8: Deploy the LLO consumer workflow
	testLogger.Info().Msg("Step 8: Deploying LLO consumer workflow...")
	workflowConfig := llo_consumer_config.LLOConsumerConfig{
		StreamIDs:      []uint32{1, 4}, // Subscribe to both Format 5 and Format 7 streams
		MaxFrequencyMs: 1000,
	}
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
	testLogger.Info().Msg("✓ Workflow deployed")

	// Step 9: Wait for LLO reports with magic numbers
	testLogger.Info().Msg("Step 9: Waiting for LLO reports with magic numbers...")

	expectedLog := "LLO_E2E_VALUE"
	timeout := 3 * time.Minute
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - workflow did not receive reports")
	testLogger.Info().Msg("✓ Workflow received LLO reports with magic numbers")

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║  ✓ LLO STREAMS TRIGGER E2E TEST PASSED                               ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")
}

