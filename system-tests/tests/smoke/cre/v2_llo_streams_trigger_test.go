package cre

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	crevalues "github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	mockcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mockpb "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	llo_consumer_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/llo_streams_trigger/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// LLO Streams Trigger config path - uses the capabilities DON topology with LLO feature
// Uses a test-specific config without gateway DON and with streams-trigger capability
const lloStreamsConfigPath = "/configs/workflow-capabilities-llo-don.toml"

// Mock test uses a config where capabilities DON has only "mock" (no "streams-trigger")
// so the workflow subscribes to the mock; mock must register as streams-trigger@2.0.0
const lloStreamsMockConfigPath = "/configs/workflow-capabilities-llo-don-mock.toml"

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
	// Use mock-only config so workflow subscribes to mock (not real streams-trigger)
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, lloStreamsMockConfigPath), v2RegistriesFlags...)

	t.Run("[v2] LLO Streams Trigger (Mock) - "+topology, func(t *testing.T) {
		ExecuteLLOStreamsTriggerTest(t, testEnv)
	})
}

// Test_CRE_V2_LLO_Streams_Trigger_E2E runs the full E2E test with real LLO plugin.
// This test automatically deploys the full LLO infrastructure:
// 1. LLO contracts (Configurator, ChannelConfigStore)
// 2. OCR configuration with proper encryption keys
// 3. Stream jobs and LLO jobs with CRE transmitter
// 4. Channel definitions (from LLO infrastructure setup)
//
// The test verifies end-to-end data flow by checking for magic numbers in workflow logs.
// Triggering is done by the LLO plugin once OCR rounds produce reports (stream jobs -> LLO jobs -> CRE transmitter).
//
// To run:
//
//	cd system-tests/tests
//	go test -timeout 15m -run "Test_CRE_V2_LLO_Streams_Trigger_E2E" ./smoke/cre/...
func Test_CRE_V2_LLO_Streams_Trigger_E2E(t *testing.T) {
	topology := os.Getenv("TOPOLOGY_NAME")

	// Enable OnchainPublicKey signer extraction for LLO tests
	// This ensures the capability registry uses OnchainPublicKey addresses for signature verification
	// The override in keystone_llo.go checks this environment variable
	t.Setenv("USE_LLO_ONCHAIN_SIGNER", "true")

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

	// Step 1: Deploy the LLO consumer workflow (mock test: subscribe to mock@1.0.0)
	testLogger.Info().Msg("Step 1: Deploying LLO consumer workflow...")
	workflowConfig := llo_consumer_config.LLOConsumerConfig{
		StreamIDs:           []uint32{1},  // Subscribe to stream ID 1
		MaxFrequencyMs:      1000,         // 1 second max frequency
		TriggerCapabilityID: "mock@1.0.0", // Mock test uses mock-only DON; workflow must subscribe to mock trigger
	}
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
	testLogger.Info().Msg("✓ Workflow deployed successfully")

	// Step 2: Wait for capability registration and DON-to-DON discovery
	testLogger.Info().Msg("Step 2: Waiting for capability discovery...")
	time.Sleep(15 * time.Second) // Allow time for capability registration and workflow subscription
	testLogger.Info().Msg("✓ Capability discovery period completed")

	// Step 2b: Wait for the workflow to subscribe to streams-trigger on the mock (mock-only config).
	// The mock capability only has a trigger after RegisterTrigger is called; SendTriggerEvent
	// returns "cannot find trigger" until then.
	testLogger.Info().Msg("Step 2b: Waiting for trigger subscribers (workflow must subscribe first)...")
	controller := mockcap.NewMockCapabilityController(framework.L)
	if err := controller.ConnectAll(getMockCapabilityAddresses(testEnv), true, false); err != nil {
		require.NoError(t, err, "Failed to connect to mock capability controllers")
	}
	err := controller.WaitForTriggerSubscribers(context.Background(), "mock@1.0.0", 30*time.Second)
	require.NoError(t, err, "Timeout waiting for workflow to subscribe to mock@1.0.0 on mock - use workflow-capabilities-llo-don-mock.toml (capabilities DON has only mock)")

	// Step 3: Inject mock LLO reports via the mock capability controller
	testLogger.Info().Msg("Step 3: Injecting mock LLO reports...")
	err = injectMockLLOReports(t, testEnv, 5) // Inject 5 reports
	require.NoError(t, err, "Failed to inject mock LLO reports")
	testLogger.Info().Msg("✓ Mock reports injected")

	// Step 4: Verify workflow received the reports
	testLogger.Info().Msg("Step 4: Verifying workflow received reports...")
	expectedLog := "LLO_E2E_VALUE"
	timeout := 45 * time.Second
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
		eventID := fmt.Sprintf("mock-llo-report-%d", i+1)
		report, err := generateMockLLOReport(uint64(i+1), eventID)
		if err != nil {
			return fmt.Errorf("failed to generate mock report %d: %w", i, err)
		}

		// Build OCRTriggerEvent so the mock can set Event.Outputs (required by SignedReportRemoteAggregator).
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: report.ConfigDigest,
			SeqNr:        report.SeqNr,
			Report:       report.Report,
			Sigs:         make([]capabilities.OCRAttributedOnchainSignature, len(report.Sigs)),
		}
		for j, s := range report.Sigs {
			ocrEvent.Sigs[j] = capabilities.OCRAttributedOnchainSignature{Signer: s.Signer, Signature: s.Signature}
		}
		outputsMap, err := ocrEvent.ToMap()
		if err != nil {
			return fmt.Errorf("failed to build OCR event map: %w", err)
		}
		// Serialize map to bytes in the format the mock's BytesToMap expects (proto Value with MapValue).
		outputsBytes, err := proto.Marshal(crevalues.Proto(outputsMap))
		if err != nil {
			return fmt.Errorf("failed to marshal outputs map: %w", err)
		}

		// Marshal the streams.Report into an anypb.Any (Payload is passed through by the mock).
		anyReport, err := anypb.New(report)
		if err != nil {
			return fmt.Errorf("failed to marshal report to Any: %w", err)
		}

		req := &mockpb.SendTriggerEventRequest{
			TriggerID:   "mock@1.0.0", // Must match trigger ID the workflow subscribed to (mock test uses mock@1.0.0)
			TriggerType: "llo",
			ID:          eventID, // Must match OCRTriggerReport.EventID so aggregator accepts the response
			Payload:     anyReport,
			Outputs:     outputsBytes, // Aggregator expects Event.Outputs to parse as OCRTriggerEvent
			OCREvent: &mockpb.OCRTriggerEvent{
				ConfigDigest: report.ConfigDigest,
				SeqNr:        report.SeqNr,
				Report:       report.Report,
				Sigs:         convertToMockSigs(report.Sigs),
			},
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
// and Format 5 (CapabilityTrigger) layout so the workflow's decodeFormat5Report succeeds.
// eventID must match the trigger event ID (req.ID) so SignedReportRemoteAggregator accepts the response.
// configDigest is the OCR config digest (32-byte SHA256). For mock tests any fixed digest is fine.
func generateMockLLOReport(seqNr uint64, eventID string) (*streams.Report, error) {
	timestamp := uint64(time.Now().UnixNano())
	configDigest, _ := hex.DecodeString("00091599c39d29821b4949b9ba237d2d1d9b7369087a71283c921034898320b0")

	// Format 5 expects Outputs.Payload = list of stream maps, each with StreamID and Decimal (shopspring binary).
	// Workflow uses stream ID 1 and expects MAGIC_NUMBER_FORMAT5 (424242).
	dec := decimal.NewFromInt(424242)
	decBytes, err := dec.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal decimal: %w", err)
	}
	payloadList := valuespb.NewListValue([]*valuespb.Value{
		valuespb.NewMapValue(map[string]*valuespb.Value{
			"StreamID": valuespb.NewInt64Value(1),
			"Decimal":  valuespb.NewBytesValue(decBytes),
		}),
	})
	// payloadList is a *Value wrapping List; we need Outputs = Map with "Payload" -> that Value
	outputsMap := &valuespb.Map{
		Fields: map[string]*valuespb.Value{
			"Payload": payloadList,
		},
	}

	ocrReport := &capabilitiespb.OCRTriggerReport{
		EventID:   eventID, // Must match req.ID so aggregator's aggregateProtobuf accepts (rep.EventID == triggerEventID)
		Timestamp: timestamp,
		Outputs:   outputsMap,
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

// convertToMockSigs converts streams.Report signatures to mock pb type for SendTriggerEventRequest.OCREvent.
func convertToMockSigs(sigs []*streams.OCRSignature) []*mockpb.OCRAttributedOnchainSignature {
	if len(sigs) == 0 {
		return nil
	}
	out := make([]*mockpb.OCRAttributedOnchainSignature, len(sigs))
	for i, s := range sigs {
		out[i] = &mockpb.OCRAttributedOnchainSignature{Signer: s.Signer, Signature: s.Signature}
	}
	return out
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

// ExecuteLLOStreamsTriggerE2EWithFullLLO runs the E2E test with actual LLO plugin.
// LLO infrastructure (Configurator, ChannelConfigStore, channel defs, stream jobs,
// LLO jobs, OCR config) is deployed by StreamsTrigger.PostEnvStartup during env startup.
// This test deploys the consumer workflow and verifies end-to-end data flow via magic numbers.
func ExecuteLLOStreamsTriggerE2EWithFullLLO(t *testing.T, testEnv *ttypes.TestEnvironment) {
	var err error
	testLogger := framework.L
	workflowFileLocation := "llo_consumer/main.go"
	workflowName := "llo-consumer-e2e"

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        LLO STREAMS TRIGGER E2E TEST (FULL LLO)                       ║")
	testLogger.Info().Msg("║  LLO deployed via PostEnvStartup; Magic: Format5=424242, Format7=555555 ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Step 1: Wait for LLO jobs to detect config and start OCR rounds (PostEnvStartup already ran)
	testLogger.Info().Msg("Step 1: Waiting for LLO jobs to detect config and start OCR rounds...")
	time.Sleep(30 * time.Second)

	// Step 2: Wait for CRE Transmitter to register
	testLogger.Info().Msg("Step 2: Waiting for CRE Transmitter to register...")
	time.Sleep(10 * time.Second)

	// Step 3: Wait for capability discovery before deploying workflow
	testLogger.Info().Msg("Step 3: Waiting for capability discovery (streams-trigger must be registered)...")
	time.Sleep(40 * time.Second)

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Step 4: Deploy the LLO consumer workflow
	testLogger.Info().Msg("Step 4: Deploying LLO consumer workflow...")
	workflowConfig := llo_consumer_config.LLOConsumerConfig{
		StreamIDs:      []uint32{1, 4}, // Subscribe to both Format 5 and Format 7 streams
		MaxFrequencyMs: 1000,
	}
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
	testLogger.Info().Msg("✓ Workflow deployed")

	// Step 5: Wait for LLO reports with magic numbers (Format 5=424242, Format 7=555555)
	testLogger.Info().Msg("Step 5: Waiting for LLO reports with magic numbers...")
	timeout := 3 * time.Minute

	// First, explicitly wait for Format 5 (424242)
	// Check for "Format=5" to ensure we get Format 5 reports specifically
	expectedLogFormat5 := "Format=5"
	testLogger.Info().Msg("Waiting for Format 5 reports (424242)...")
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLogFormat5, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - workflow did not receive Format 5 reports")

	// Verify Format 5 value is correct (424242)
	expectedLogFormat5Value := "Value=424242"
	testLogger.Info().Msg("Verifying Format 5 value (424242)...")
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLogFormat5Value, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - Format 5 value is not 424242")
	testLogger.Info().Msg("✓ Workflow received Format 5 reports (424242)")

	// Then, explicitly wait for Format 7 (555555 from calculated stream: 111111 * 5)
	// Check for Format=7 specifically to verify calculated streams are working
	expectedLogFormat7 := "Format=7"
	testLogger.Info().Msg("Waiting for Format 7 reports with calculated stream (555555 = 111111 * 5)...")
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLogFormat7, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - workflow did not receive Format 7 reports with calculated stream")

	// Also verify the calculated value is correct (555555)
	expectedLogFormat7Value := "Value=555555"
	testLogger.Info().Msg("Verifying calculated stream value (555555)...")
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLogFormat7Value, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "LLO Streams Trigger E2E test failed - calculated stream value is not 555555")
	testLogger.Info().Msg("✓ Workflow received Format 7 reports (555555)")

	testLogger.Info().Msg("✓ Workflow received LLO reports with magic numbers (both Format 5 and Format 7 with calculated stream)")

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║  ✓ LLO STREAMS TRIGGER E2E TEST PASSED                               ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")
}
