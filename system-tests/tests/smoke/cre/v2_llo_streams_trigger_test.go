package cre

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"

	mockcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mockpb "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// LLOConsumerConfig defines the workflow configuration
type LLOConsumerConfig struct {
	StreamIDs      []uint32 `yaml:"stream_ids"`
	MaxFrequencyMs uint64   `yaml:"max_frequency_ms"`
}

// ExecuteLLOStreamsTriggerTest runs the full E2E test for streams-trigger@2.0.0
// This test verifies:
// 1. The streams-trigger capability is registered and exposed by the Streams DON
// 2. The Workflow DON can discover and subscribe to the capability
// 3. LLO reports are correctly routed from Streams DON to Workflow DON
// 4. The consumer workflow receives and processes the reports
func ExecuteLLOStreamsTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/llo-consumer/main.go"
	workflowName := "llo-consumer"

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        LLO STREAMS TRIGGER E2E TEST                                  ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Start Beholder listener to capture workflow logs
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Step 1: Deploy the LLO consumer workflow
	testLogger.Info().Msg("Step 1: Deploying LLO consumer workflow...")
	workflowConfig := LLOConsumerConfig{
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
	defer controller.CloseAll()

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
	ocrReport := &capabilitiespb.OCRTriggerReport{
		EventID:   fmt.Sprintf("mock_event_%d", seqNr),
		Timestamp: timestamp,
		StreamId:  1, // Stream ID 1 matches our workflow config
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
func getMockCapabilityAddresses(testEnv *ttypes.TestEnvironment) []string {
	// In the DON-to-DON setup, mock capabilities are exposed on port 5002
	// on the Streams DON nodes
	addresses := []string{}

	// Try common patterns for Docker-based setup
	for i := 0; i < 4; i++ {
		// Format: streams-node{i}:5002 (inside Docker network)
		// Or: localhost:15002+i (mapped port)
		addresses = append(addresses, fmt.Sprintf("localhost:%d", 15002+i))
	}

	return addresses
}

// ExecuteLLOStreamsTriggerE2EWithRealLLO runs the E2E test with actual LLO plugin
// This test requires:
// - LLO jobs deployed on Streams DON nodes
// - Stream jobs configured to fetch data
// - Channel definitions set with ReportFormat=5 (CapabilityTrigger)
func ExecuteLLOStreamsTriggerE2EWithRealLLO(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/llo-consumer/main.go"
	workflowName := "llo-consumer-real"

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║        LLO STREAMS TRIGGER E2E TEST (REAL LLO)                      ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")

	// Start Beholder listener
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Deploy workflow
	testLogger.Info().Msg("Deploying LLO consumer workflow...")
	workflowConfig := LLOConsumerConfig{
		StreamIDs:      []uint32{1},
		MaxFrequencyMs: 1000,
	}
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Wait for real LLO reports to be received
	// With real LLO, reports should arrive every ~1 second
	testLogger.Info().Msg("Waiting for real LLO reports...")
	expectedLog := "LLO_E2E_VALUE"
	timeout := 5 * time.Minute // Longer timeout for real LLO setup
	err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "Real LLO E2E test failed")

	testLogger.Info().Msg("╔══════════════════════════════════════════════════════════════════════╗")
	testLogger.Info().Msg("║  ✓ LLO STREAMS TRIGGER E2E (REAL LLO) PASSED                        ║")
	testLogger.Info().Msg("╚══════════════════════════════════════════════════════════════════════╝")
}
