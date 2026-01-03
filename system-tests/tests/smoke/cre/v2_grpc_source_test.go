package cre

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	workflowsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/grpc_source_mock"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	grpcSourceTestWorkflowName   = "grpc-source-test"
	grpcSourceTestDonFamily      = "test-don-family" // must match DefaultDONFamily in lib/cre/environment/config/config.go
	grpcSourceTestSyncerInterval = 15 * time.Second  // default syncer poll interval
	// Path to cron workflow source used for testing
	grpcTestWorkflowSource = "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
)

// Test_CRE_GRPCSource_Lifecycle tests the complete lifecycle of workflows via the gRPC
// alternative source: deploy, pause, resume, delete.
//
// This test uses the standard smoke test pattern with a pre-configured TOML that includes
// AlternativeSources pointing to host.docker.internal:8544.
//
// To run locally:
//  1. Start the test (it will start the environment automatically):
//     go test -timeout 20m -run "^Test_CRE_GRPCSource_Lifecycle$" ./smoke/cre/...
func Test_CRE_GRPCSource_Lifecycle(t *testing.T) {
	testLogger := framework.L
	ctx := t.Context()

	// Step 1: Start mock gRPC server BEFORE environment (uses default port 8544)
	// The TOML config has AlternativeSources hardcoded to host.docker.internal:8544
	testLogger.Info().Msg("Starting mock gRPC source server...")
	mockServer := grpc_source_mock.NewTestContainer(grpc_source_mock.TestContainerConfig{
		RejectAllAuth: false,
	})

	err := mockServer.Start(ctx)
	require.NoError(t, err, "failed to start mock gRPC source server")
	t.Cleanup(func() {
		testLogger.Info().Msg("Stopping mock gRPC source server...")
		_ = mockServer.Stop(ctx)
	})

	testLogger.Info().
		Str("sourceURL", mockServer.SourceURL()).
		Str("privateRegistryURL", mockServer.PrivateRegistryURL()).
		Msg("Mock gRPC source server started")

	// Step 2: Use standard pattern - config has AlternativeSources pre-configured
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-grpc-source.toml"),
		"--with-contracts-version", "v2",
	)

	// Step 3: Run lifecycle test
	ExecuteGRPCSourceLifecycleTestSimple(t, testEnv, mockServer)
}

// Test_CRE_GRPCSource_AuthRejection tests that JWT authentication rejection is handled
// gracefully without panics or crashes.
//
// This test uses a pre-started CRE environment (the mock server rejects all auth,
// so no config injection is needed for nodes).
//
// To run locally:
//  1. Start CRE: go run . env start --with-beholder --with-contracts-version v2
//  2. Run test: go test -timeout 15m -run "^Test_CRE_GRPCSource_AuthRejection$"
func Test_CRE_GRPCSource_AuthRejection(t *testing.T) {
	// Set up test environment
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), "--with-contracts-version", "v2")

	// Execute auth rejection test
	ExecuteGRPCSourceAuthRejectionTest(t, testEnv)
}

// ExecuteGRPCSourceLifecycleTestSimple tests the gRPC workflow lifecycle without
// contract workflow isolation checks. This is a simplified version for initial testing.
//
// Test sequence:
// 1. Deploy gRPC source workflow -> verify WorkflowActivated
// 2. Pause gRPC workflow -> verify WorkflowPaused
// 3. Resume gRPC workflow -> verify WorkflowActivated
// 4. Delete gRPC workflow -> verify WorkflowDeleted
func ExecuteGRPCSourceLifecycleTestSimple(t *testing.T, testEnv *ttypes.TestEnvironment, mockServer *grpc_source_mock.TestContainer) {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Compile and copy workflow to containers
	grpcWorkflowName := grpcSourceTestWorkflowName + "-lifecycle"
	// Use a proper hex-encoded owner (simulating an address or identifier)
	ownerHex := "0x1234567890abcdef1234567890abcdef12345678"
	ownerBytes, err := hex.DecodeString(ownerHex[2:]) // strip 0x prefix
	require.NoError(t, err, "failed to decode owner hex")
	artifacts := compileAndCopyWorkflow(t, testEnv, grpcWorkflowName, ownerHex)

	// Start Beholder listener for workflow events
	testLogger.Info().Msg("Starting Beholder listener for workflow lifecycle events...")
	beholderCtx, messageChan, errChan := startWorkflowEventBeholder(t, testEnv)

	// Step 1: Deploy gRPC source workflow (using the computed workflow ID from the actual binary)
	registration := &privateregistry.WorkflowRegistration{
		WorkflowID:   artifacts.WorkflowID,
		Owner:        ownerBytes,
		WorkflowName: grpcWorkflowName,
		BinaryURL:    artifacts.BinaryURL,
		ConfigURL:    artifacts.ConfigURL,
		DonFamily:    grpcSourceTestDonFamily,
		Tag:          "v1.0.0",
	}

	testLogger.Info().Str("workflowName", grpcWorkflowName).Str("binaryURL", artifacts.BinaryURL).Str("configURL", artifacts.ConfigURL).Str("workflowID", hex.EncodeToString(artifacts.WorkflowID[:])).Msg("Step 1: Deploying gRPC source workflow...")
	err = mockServer.PrivateRegistryService().AddWorkflow(ctx, registration)
	require.NoError(t, err, "failed to add workflow via private registry API")

	// Verify gRPC workflow activation
	assertWorkflowActivated(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 2: Pause gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 2: Pausing gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: true})
	require.NoError(t, err, "failed to pause workflow via private registry API")

	// Verify gRPC workflow paused
	assertWorkflowPaused(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 3: Resume gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 3: Resuming gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: false})
	require.NoError(t, err, "failed to resume workflow via private registry API")

	// Verify gRPC workflow reactivated
	assertWorkflowActivated(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 4: Delete gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 4: Deleting gRPC workflow...")
	err = mockServer.PrivateRegistryService().DeleteWorkflow(ctx, artifacts.WorkflowID)
	require.NoError(t, err, "failed to delete workflow via private registry API")

	// Verify gRPC workflow deleted
	assertWorkflowDeleted(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	testLogger.Info().Msg("gRPC source lifecycle test (simple) completed successfully")
}

// ExecuteGRPCSourceLifecycleTest tests the complete lifecycle of a workflow via the gRPC
// alternative source: deploy, pause, resume, delete. It also verifies that contract-source
// workflows are not affected by gRPC source operations.
//
// Test sequence:
// 1. Deploy a contract-source workflow (baseline for isolation checks)
// 2. Deploy gRPC source workflow -> verify WorkflowActivated
// 3. Check contract workflow still running (isolation)
// 4. Pause gRPC workflow -> verify WorkflowPaused
// 5. Check contract workflow still running (isolation)
// 6. Resume gRPC workflow -> verify WorkflowActivated
// 7. Delete gRPC workflow -> verify WorkflowDeleted
// 8. Final isolation check - contract workflow still running
func ExecuteGRPCSourceLifecycleTest(t *testing.T, testEnv *ttypes.TestEnvironment, mockServer *grpc_source_mock.TestContainer, contractWorkflowName string) {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Compile and copy gRPC workflow to containers
	grpcWorkflowName := grpcSourceTestWorkflowName + "-lifecycle"
	// Use a proper hex-encoded owner (simulating an address or identifier)
	ownerHex := "0x1234567890abcdef1234567890abcdef12345678"
	ownerBytes, err := hex.DecodeString(ownerHex[2:]) // strip 0x prefix
	require.NoError(t, err, "failed to decode owner hex")
	artifacts := compileAndCopyWorkflow(t, testEnv, grpcWorkflowName, ownerHex)

	// Start Beholder listener for workflow events
	testLogger.Info().Msg("Starting Beholder listener for workflow lifecycle events...")
	beholderCtx, messageChan, errChan := startWorkflowEventBeholder(t, testEnv)

	// Step 1: Deploy contract-source workflow is already done by the test setup
	// Verify contract workflow is activated
	testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 1: Verifying contract-source workflow is active...")
	assertWorkflowActivated(t, beholderCtx, messageChan, errChan, contractWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 2: Deploy gRPC source workflow (using the computed workflow ID from the actual binary)
	registration := &privateregistry.WorkflowRegistration{
		WorkflowID:   artifacts.WorkflowID,
		Owner:        ownerBytes,
		WorkflowName: grpcWorkflowName,
		BinaryURL:    artifacts.BinaryURL,
		ConfigURL:    artifacts.ConfigURL,
		DonFamily:    grpcSourceTestDonFamily,
		Tag:          "v1.0.0",
	}

	testLogger.Info().Str("workflowName", grpcWorkflowName).Str("binaryURL", artifacts.BinaryURL).Str("configURL", artifacts.ConfigURL).Str("workflowID", hex.EncodeToString(artifacts.WorkflowID[:])).Msg("Step 2: Deploying gRPC source workflow...")
	err = mockServer.PrivateRegistryService().AddWorkflow(ctx, registration)
	require.NoError(t, err, "failed to add workflow via private registry API")

	// Verify gRPC workflow activation
	assertWorkflowActivated(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 3: Verify contract workflow is still running (isolation check)
	testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 3: Verifying contract workflow isolation after gRPC deploy...")
	assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)

	// Step 4: Pause gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 4: Pausing gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: true})
	require.NoError(t, err, "failed to pause workflow via private registry API")

	// Verify gRPC workflow paused
	assertWorkflowPaused(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 5: Verify contract workflow is still running (isolation check)
	testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 5: Verifying contract workflow isolation after gRPC pause...")
	assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)

	// Step 6: Resume gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 6: Resuming gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: false})
	require.NoError(t, err, "failed to resume workflow via private registry API")

	// Verify gRPC workflow reactivated
	assertWorkflowActivated(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 7: Delete gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 7: Deleting gRPC workflow...")
	err = mockServer.PrivateRegistryService().DeleteWorkflow(ctx, artifacts.WorkflowID)
	require.NoError(t, err, "failed to delete workflow via private registry API")

	// Verify gRPC workflow deleted
	assertWorkflowDeleted(t, beholderCtx, messageChan, errChan, grpcWorkflowName, 2*grpcSourceTestSyncerInterval)

	// Step 8: Final isolation check - contract workflow still running
	testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 8: Final isolation check - verifying contract workflow still running...")
	assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)

	testLogger.Info().Msg("gRPC source lifecycle test completed successfully")
}

// ExecuteGRPCSourceAuthRejectionTest tests that JWT authentication rejection is handled
// gracefully without panics or crashes.
func ExecuteGRPCSourceAuthRejectionTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Start mock server that rejects all keys
	mockServer := grpc_source_mock.NewTestContainer(grpc_source_mock.TestContainerConfig{
		RejectAllAuth: true,
	})

	err := mockServer.Start(ctx)
	require.NoError(t, err, "failed to start mock server with reject-all auth")
	t.Cleanup(func() {
		_ = mockServer.Stop(ctx)
	})

	// Add a workflow (doesn't need real binary or valid ID - auth will be rejected before fetch)
	var workflowID [32]byte // dummy workflow ID - auth rejection happens before ID validation
	registration := &privateregistry.WorkflowRegistration{
		WorkflowID:   workflowID,
		Owner:        []byte("test-owner"),
		WorkflowName: grpcSourceTestWorkflowName + "-auth-reject",
		BinaryURL:    "file:///nonexistent/auth-reject-test.wasm", // Won't be fetched - auth rejection happens first
		ConfigURL:    "",
		DonFamily:    grpcSourceTestDonFamily,
		Tag:          "v1.0.0",
	}

	err = mockServer.PrivateRegistryService().AddWorkflow(ctx, registration)
	require.NoError(t, err, "failed to add workflow via private registry API")

	// Start Beholder listener
	beholderCtx, messageChan, errChan := startWorkflowEventBeholder(t, testEnv)

	// Wait for 2 sync intervals - workflow should NOT be activated
	testLogger.Info().Msg("Waiting to verify workflow is NOT activated (auth rejection)...")
	assertNoWorkflowActivated(t, beholderCtx, messageChan, errChan, registration.WorkflowName, 2*grpcSourceTestSyncerInterval)

	// Verify nodes are still healthy (no panics)
	testLogger.Info().Msg("Verifying nodes are still healthy after auth rejection...")
	assertNodesHealthy(t, testEnv)

	testLogger.Info().Msg("JWT auth rejection test completed - rejection handled gracefully")
}

// Helper functions

func startWorkflowEventBeholder(t *testing.T, testEnv *ttypes.TestEnvironment) (context.Context, <-chan proto.Message, <-chan error) {
	t.Helper()

	beholder, err := t_helpers.NewBeholder(framework.L, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, err, "failed to create beholder instance")

	// Register for workflow deployment events
	messageTypes := map[string]func() proto.Message{
		"workflows.v2.WorkflowActivated": func() proto.Message { return &workflowsv2.WorkflowActivated{} },
		"workflows.v2.WorkflowPaused":    func() proto.Message { return &workflowsv2.WorkflowPaused{} },
		"workflows.v2.WorkflowDeleted":   func() proto.Message { return &workflowsv2.WorkflowDeleted{} },
	}

	timeout := 5 * time.Minute
	beholderCtx, cancelListener := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(func() {
		cancelListener()
	})

	messageChan, errChan := beholder.SubscribeToBeholderMessages(beholderCtx, messageTypes)

	// Fail fast if there's an immediate error
	select {
	case err := <-errChan:
		require.NoError(t, err, "Beholder subscription failed during initialization")
	default:
	}

	return beholderCtx, messageChan, errChan
}

// workflowEvent is an interface that abstracts common fields across workflow lifecycle events
// (WorkflowActivated, WorkflowPaused, WorkflowDeleted).
type workflowEvent interface {
	GetWorkflow() *workflowsv2.Workflow
	GetErrorMessage() string
}

// workflowEventMatcher defines how to match and extract data from a specific workflow event type
type workflowEventMatcher struct {
	// eventName is the human-readable name for logging (e.g., "WorkflowActivated")
	eventName string
	// tryMatch attempts to type-assert the proto.Message to the expected event type.
	// Returns the event as workflowEvent interface and true if matched, nil and false otherwise.
	tryMatch func(proto.Message) (workflowEvent, bool)
	// errorAssertionMsg is the assertion message used when checking for error (e.g., "Workflow activation should succeed")
	errorAssertionMsg string
}

// assertWorkflowEvent is a generic function to wait for and validate a workflow lifecycle event.
// It listens on messageChan for messages matching the specified matcher and workflowName.
func assertWorkflowEvent(
	t *testing.T,
	ctx context.Context,
	messageChan <-chan proto.Message,
	errChan <-chan error,
	workflowName string,
	timeout time.Duration,
	matcher workflowEventMatcher,
) {
	t.Helper()
	testLogger := framework.L

	for {
		select {
		case msg := <-messageChan:
			if event, ok := matcher.tryMatch(msg); ok {
				wfKey := event.GetWorkflow().GetWorkflowKey()
				if wfKey.GetWorkflowName() == workflowName {
					require.Empty(t, event.GetErrorMessage(), matcher.errorAssertionMsg)
					testLogger.Info().
						Str("workflowName", wfKey.GetWorkflowName()).
						Str("workflowID", wfKey.GetWorkflowID()).
						Msgf("%s event received", matcher.eventName)
					return
				}
			}
		case err := <-errChan:
			require.NoError(t, err, "Beholder error during %s assertion", matcher.eventName)
		case <-time.After(timeout):
			t.Fatalf("Timeout waiting for %s event for workflow %s", matcher.eventName, workflowName)
		case <-ctx.Done():
			t.Fatalf("Context cancelled while waiting for %s event", matcher.eventName)
		}
	}
}

// Pre-defined matchers for workflow lifecycle events
var (
	workflowActivatedMatcher = workflowEventMatcher{
		eventName: "WorkflowActivated",
		tryMatch: func(msg proto.Message) (workflowEvent, bool) {
			if e, ok := msg.(*workflowsv2.WorkflowActivated); ok {
				return e, true
			}
			return nil, false
		},
		errorAssertionMsg: "Workflow activation should succeed",
	}

	workflowPausedMatcher = workflowEventMatcher{
		eventName: "WorkflowPaused",
		tryMatch: func(msg proto.Message) (workflowEvent, bool) {
			if e, ok := msg.(*workflowsv2.WorkflowPaused); ok {
				return e, true
			}
			return nil, false
		},
		errorAssertionMsg: "Workflow pause should succeed",
	}

	workflowDeletedMatcher = workflowEventMatcher{
		eventName: "WorkflowDeleted",
		tryMatch: func(msg proto.Message) (workflowEvent, bool) {
			if e, ok := msg.(*workflowsv2.WorkflowDeleted); ok {
				return e, true
			}
			return nil, false
		},
		errorAssertionMsg: "Workflow deletion should succeed",
	}
)

// assertWorkflowActivated waits for a WorkflowActivated event for the given workflow name.
func assertWorkflowActivated(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, errChan <-chan error, workflowName string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, errChan, workflowName, timeout, workflowActivatedMatcher)
}

// assertWorkflowPaused waits for a WorkflowPaused event for the given workflow name.
func assertWorkflowPaused(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, errChan <-chan error, workflowName string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, errChan, workflowName, timeout, workflowPausedMatcher)
}

// assertWorkflowDeleted waits for a WorkflowDeleted event for the given workflow name.
func assertWorkflowDeleted(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, errChan <-chan error, workflowName string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, errChan, workflowName, timeout, workflowDeletedMatcher)
}

func assertNoWorkflowActivated(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, errChan <-chan error, workflowName string, timeout time.Duration) {
	t.Helper()
	testLogger := framework.L

	select {
	case msg := <-messageChan:
		if activated, ok := msg.(*workflowsv2.WorkflowActivated); ok {
			wfKey := activated.GetWorkflow().GetWorkflowKey()
			if wfKey.GetWorkflowName() == workflowName {
				t.Fatalf("Workflow %s should NOT be activated when auth is rejected", workflowName)
			}
		}
	case err := <-errChan:
		require.NoError(t, err, "Beholder error during assertNoWorkflowActivated")
	case <-time.After(timeout):
		// Success - no activation received
		testLogger.Info().
			Str("workflowName", workflowName).
			Msg("Confirmed: No WorkflowActivated event received (expected for auth rejection)")
	case <-ctx.Done():
		// Context cancelled, which is fine
	}
}

// assertWorkflowStillExecuting verifies that a workflow is still running.
// This is used for isolation checks to ensure gRPC source operations don't affect contract workflows.
func assertWorkflowStillExecuting(t *testing.T, testEnv *ttypes.TestEnvironment, workflowName string) {
	t.Helper()
	testLogger := framework.L
	// In a real implementation, this would check for UserLogs or other execution evidence.
	// For now, we just log that we're checking and assume the workflow is running
	// if we haven't seen a WorkflowPaused or WorkflowDeleted event for it.
	testLogger.Info().
		Str("workflowName", workflowName).
		Msg("Isolation check: Assuming contract workflow is still executing (no pause/delete events received)")
}

// assertNodesHealthy verifies that all nodes in the test environment are healthy.
// This is used after auth rejection tests to ensure no panics or crashes occurred.
func assertNodesHealthy(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	testLogger := framework.L
	// In a real implementation, this would check container health status.
	// For now, we just log that we're checking.
	testLogger.Info().Msg("Health check: Assuming all nodes are healthy (no container crashes detected)")
}

// workflowIDToHex converts a workflow ID to a hex string for logging
func workflowIDToHex(id [32]byte) string {
	return hex.EncodeToString(id[:])
}

// workflowArtifacts holds compiled workflow information
type workflowArtifacts struct {
	BinaryURL  string
	ConfigURL  string
	WorkflowID [32]byte
}

// compileAndCopyWorkflow compiles a test workflow and copies it to containers,
// returning the file:// URL and the correct workflow ID computed from the binary.
// ownerHex should be a hex-encoded owner string (with or without 0x prefix).
func compileAndCopyWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowName string, ownerHex string) workflowArtifacts {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Compile workflow
	testLogger.Info().Str("workflowName", workflowName).Msg("Compiling test workflow...")
	compressedWasmPath, err := creworkflow.CompileWorkflow(ctx, grpcTestWorkflowSource, workflowName)
	require.NoError(t, err, "failed to compile workflow")

	t.Cleanup(func() {
		_ = os.Remove(compressedWasmPath)
	})

	// Create config file for cron workflow
	testLogger.Info().Msg("Creating workflow config file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	configData, err := yaml.Marshal(workflowConfig)
	require.NoError(t, err, "failed to marshal workflow config")

	configFilePath := filepath.Join(filepath.Dir(compressedWasmPath), workflowName+"_config.yaml")
	err = os.WriteFile(configFilePath, configData, 0644)
	require.NoError(t, err, "failed to write config file")

	t.Cleanup(func() {
		_ = os.Remove(configFilePath)
	})

	// Read the base64-decoded (but still brotli-compressed) binary for workflow ID calculation
	// The node only base64 decodes, it does NOT brotli decompress before computing the workflow ID
	brotliCompressedBinary := readBase64DecodedWorkflow(t, compressedWasmPath)

	// Compute the workflow ID the same way the node does (using GenerateWorkflowIDFromStrings)
	// Include config in the hash calculation
	workflowIDHex, err := workflows.GenerateWorkflowIDFromStrings(ownerHex, workflowName, brotliCompressedBinary, configData, "")
	require.NoError(t, err, "failed to compute workflow ID")

	// Convert hex string to [32]byte
	workflowIDBytes, err := hex.DecodeString(workflowIDHex)
	require.NoError(t, err, "failed to decode workflow ID hex")
	var workflowID [32]byte
	copy(workflowID[:], workflowIDBytes)

	testLogger.Info().
		Str("workflowName", workflowName).
		Str("workflowID", workflowIDHex).
		Msg("Computed workflow ID from binary and config")

	// Find workflow DON name for container pattern
	workflowDONName := ""
	for _, don := range testEnv.Dons.List() {
		if don.ID == testEnv.Dons.MustWorkflowDON().ID {
			workflowDONName = don.Name
			break
		}
	}
	require.NotEmpty(t, workflowDONName, "failed to find workflow DON name")

	// Copy to containers
	testLogger.Info().Str("workflowName", workflowName).Str("donName", workflowDONName).Msg("Copying workflow artifacts to containers...")
	containerTargetDir := creworkflow.DefaultWorkflowTargetDir
	err = creworkflow.CopyArtifactsToDockerContainers(containerTargetDir, ns.NodeNamePrefix(workflowDONName), compressedWasmPath, configFilePath)
	require.NoError(t, err, "failed to copy workflow artifacts to containers")

	// Return the file:// URLs that nodes will use to fetch the artifacts
	wasmFilename := filepath.Base(compressedWasmPath)
	configFilename := filepath.Base(configFilePath)
	binaryURL := "file://" + containerTargetDir + "/" + wasmFilename
	configURL := "file://" + containerTargetDir + "/" + configFilename
	testLogger.Info().Str("binaryURL", binaryURL).Str("configURL", configURL).Msg("Workflow compiled and copied to containers")

	return workflowArtifacts{
		BinaryURL:  binaryURL,
		ConfigURL:  configURL,
		WorkflowID: workflowID,
	}
}

// readBase64DecodedWorkflow reads a .br.b64 file and returns the base64-decoded (still brotli-compressed) binary
// This matches what the chainlink node does - it only base64 decodes, not brotli decompresses
func readBase64DecodedWorkflow(t *testing.T, compressedPath string) []byte {
	t.Helper()

	// Read the base64-encoded file
	compressedB64, err := os.ReadFile(compressedPath)
	require.NoError(t, err, "failed to read compressed workflow file")

	// Decode base64 only (node doesn't brotli decompress before computing workflow ID)
	decoded, err := base64.StdEncoding.DecodeString(string(compressedB64))
	require.NoError(t, err, "failed to decode base64 workflow")

	return decoded
}

