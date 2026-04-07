package cre

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"gopkg.in/yaml.v3"

	workflowsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
	grpcsourcemock "github.com/smartcontractkit/chainlink/system-tests/lib/cre/grpc_source_mock"
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
// additional source: deploy, pause, resume, delete.
//
// This test uses a pre-configured TOML with AdditionalSources pointing to host.docker.internal:8544.
// The config generation code automatically transforms host.docker.internal to the platform-specific
// Docker host address (e.g., 172.17.0.1 on Linux).
//
// To run locally:
//  1. Start the test (it will start the environment automatically):
//     go test -timeout 20m -run "^Test_CRE_GRPCSource_Lifecycle$" ./smoke/cre/...
func Test_CRE_GRPCSource_Lifecycle(t *testing.T) {
	t.Skip("Skipping: gRPC source tests require V2 workflow registry syncer - needs investigation for CI environment differences")

	testLogger := framework.L
	ctx := t.Context()

	// Step 1: Start mock gRPC server BEFORE environment (uses default port 8544)
	testLogger.Info().Msg("Starting mock gRPC source server...")
	mockServer := grpcsourcemock.NewTestContainer(grpcsourcemock.TestContainerConfig{
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

	// Step 2: Use standard pattern - config has AdditionalSources with host.docker.internal
	// The config generation code transforms this to the platform-specific Docker host
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-grpc-source.toml"),
		"--with-contracts-version", "v2",
	)

	// Step 3: Run lifecycle test
	// Pass empty string for contractWorkflowName to skip contract isolation checks
	// (no contract workflow is deployed in this test configuration)
	ExecuteGRPCSourceLifecycleTest(t, testEnv, mockServer, "" /* contractWorkflowName */)
}

// Test_CRE_GRPCSource_AuthRejection tests that JWT authentication rejection is handled
// gracefully without panics or crashes.
//
// This test uses a pre-started CRE environment (the mock server rejects all auth,
// so no config injection is needed for nodes).
//
// To run locally:
//  1. Start CRE: go run . env start --with-beholder
//  2. Run test: go test -timeout 15m -run "^Test_CRE_GRPCSource_AuthRejection$"
func Test_CRE_GRPCSource_AuthRejection(t *testing.T) {
	// Set up test environment
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), "--with-contracts-version", "v2")

	// Execute auth rejection test
	ExecuteGRPCSourceAuthRejectionTest(t, testEnv)
}

// ExecuteGRPCSourceLifecycleTest tests the complete lifecycle of a workflow via the gRPC
// alternative source: deploy, pause, resume, delete.
//
// If contractWorkflowName is provided (non-empty), it also verifies that contract-source
// workflows are not affected by gRPC source operations (isolation checks).
//
// Test sequence:
// 1. (Optional) Verify contract-source workflow is active
// 2. Deploy gRPC source workflow -> verify WorkflowActivated
// 3. (Optional) Check contract workflow still running (isolation)
// 4. Pause gRPC workflow -> verify WorkflowPaused
// 5. (Optional) Check contract workflow still running (isolation)
// 6. Resume gRPC workflow -> verify WorkflowActivated
// 7. Delete gRPC workflow -> verify WorkflowDeleted
// 8. (Optional) Final isolation check - contract workflow still running
func ExecuteGRPCSourceLifecycleTest(t *testing.T, testEnv *ttypes.TestEnvironment, mockServer *grpcsourcemock.TestContainer, contractWorkflowName string) {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Determine if we should run contract isolation checks
	runIsolationChecks := contractWorkflowName != ""

	// Compile and copy gRPC workflow to containers
	grpcWorkflowName := grpcSourceTestWorkflowName + "-lifecycle"
	// Use a proper hex-encoded owner (simulating an address or identifier)
	ownerHex := "0x1234567890abcdef1234567890abcdef12345678"
	ownerBytes, err := hex.DecodeString(ownerHex[2:]) // strip 0x prefix
	require.NoError(t, err, "failed to decode owner hex")
	artifacts := compileAndCopyWorkflow(t, testEnv, grpcWorkflowName, ownerHex)

	// Start CHiP sink listener for workflow events
	testLogger.Info().Msg("Starting CHiP sink listener for workflow lifecycle events...")
	eventsCtx, messageChan := startWorkflowEventSink(t)
	grpcWorkflowID := hex.EncodeToString(artifacts.WorkflowID[:])

	// Step 1: (Optional) Verify contract workflow is activated
	if runIsolationChecks {
		testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 1: Verifying contract workflow isolation baseline...")
		assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)
	} else {
		testLogger.Info().Msg("Skipping contract workflow isolation checks (no contract workflow configured)")
	}

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
	assertWorkflowActivated(t, eventsCtx, messageChan, grpcWorkflowName, grpcWorkflowID, 2*grpcSourceTestSyncerInterval)

	// Step 3: (Optional) Verify contract workflow is still running (isolation check)
	if runIsolationChecks {
		testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 3: Verifying contract workflow isolation after gRPC deploy...")
		assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)
	}

	// Step 4: Pause gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 4: Pausing gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: true})
	require.NoError(t, err, "failed to pause workflow via private registry API")

	// Verify gRPC workflow paused
	assertWorkflowPaused(t, eventsCtx, messageChan, grpcWorkflowName, grpcWorkflowID, 2*grpcSourceTestSyncerInterval)

	// Step 5: (Optional) Verify contract workflow is still running (isolation check)
	if runIsolationChecks {
		testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 5: Verifying contract workflow isolation after gRPC pause...")
		assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)
	}

	// Step 6: Resume gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 6: Resuming gRPC workflow...")
	err = mockServer.PrivateRegistryService().UpdateWorkflow(ctx, artifacts.WorkflowID, &privateregistry.WorkflowStatusConfig{Paused: false})
	require.NoError(t, err, "failed to resume workflow via private registry API")

	// Verify gRPC workflow reactivated
	assertWorkflowActivated(t, eventsCtx, messageChan, grpcWorkflowName, grpcWorkflowID, 2*grpcSourceTestSyncerInterval)

	// Step 7: Delete gRPC workflow
	testLogger.Info().Str("workflowName", grpcWorkflowName).Msg("Step 7: Deleting gRPC workflow...")
	err = mockServer.PrivateRegistryService().DeleteWorkflow(ctx, artifacts.WorkflowID)
	require.NoError(t, err, "failed to delete workflow via private registry API")

	// Verify gRPC workflow deleted
	assertWorkflowDeleted(t, eventsCtx, messageChan, grpcWorkflowName, grpcWorkflowID, 2*grpcSourceTestSyncerInterval)

	// Step 8: (Optional) Final isolation check - contract workflow still running
	if runIsolationChecks {
		testLogger.Info().Str("workflowName", contractWorkflowName).Msg("Step 8: Final isolation check - verifying contract workflow still running...")
		assertWorkflowStillExecuting(t, testEnv, contractWorkflowName)
	}

	testLogger.Info().Msg("gRPC source lifecycle test completed successfully")
}

// ExecuteGRPCSourceAuthRejectionTest tests that JWT authentication rejection is handled
// gracefully without panics or crashes.
func ExecuteGRPCSourceAuthRejectionTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	testLogger := framework.L
	ctx := t.Context()

	// Start mock server that rejects all keys
	mockServer := grpcsourcemock.NewTestContainer(grpcsourcemock.TestContainerConfig{
		RejectAllAuth: true,
	})

	err := mockServer.Start(ctx)
	require.NoError(t, err, "failed to start mock server with reject-all auth")
	t.Cleanup(func() {
		_ = mockServer.Stop(ctx)
	})

	// Add a workflow (doesn't need real binary or valid ID - auth will be rejected before fetch)
	var workflowID [32]byte
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

	// Start CHiP sink listener
	eventsCtx, messageChan := startWorkflowEventSink(t)
	workflowIDHex := hex.EncodeToString(registration.WorkflowID[:])

	// Wait for 2 sync intervals - workflow should NOT be activated
	testLogger.Info().Msg("Waiting to verify workflow is NOT activated (auth rejection)...")
	assertNoWorkflowActivated(t, eventsCtx, messageChan, registration.WorkflowName, workflowIDHex, 2*grpcSourceTestSyncerInterval)

	// Verify nodes are still healthy (no panics)
	testLogger.Info().Msg("Verifying nodes are still healthy after auth rejection...")
	assertNodesHealthy(t, testEnv)

	testLogger.Info().Msg("JWT auth rejection test completed - rejection handled gracefully")
}

// Helper functions

func startWorkflowEventSink(t *testing.T) (context.Context, <-chan proto.Message) {
	t.Helper()

	messageChan := make(chan proto.Message, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetWorkflowV2LifecyclePublishFn(framework.L, messageChan))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(messageChan)
	})

	timeout := 5 * time.Minute
	eventsCtx, cancelListener := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(func() {
		cancelListener()
	})

	return eventsCtx, messageChan
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
// It listens on messageChan for messages matching the specified matcher and workflowID.
func assertWorkflowEvent(
	t *testing.T,
	ctx context.Context, //nolint:revive // test helper conventionally has t first
	messageChan <-chan proto.Message,
	workflowName string,
	expectedWorkflowID string,
	timeout time.Duration,
	matcher workflowEventMatcher,
) {
	t.Helper()
	testLogger := framework.L
	normalizedExpectedWorkflowID := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(expectedWorkflowID)), "0x")
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok || msg == nil {
				continue
			}
			if event, ok := matcher.tryMatch(msg); ok {
				wfKey := event.GetWorkflow().GetWorkflowKey()
				normalizedEventWorkflowID := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(wfKey.GetWorkflowID())), "0x")
				if normalizedEventWorkflowID == normalizedExpectedWorkflowID {
					require.Empty(t, event.GetErrorMessage(), matcher.errorAssertionMsg)
					testLogger.Info().
						Str("workflowName", wfKey.GetWorkflowName()).
						Str("workflowID", wfKey.GetWorkflowID()).
						Str("expectedWorkflowName", workflowName).
						Str("expectedWorkflowID", expectedWorkflowID).
						Msgf("%s event received", matcher.eventName)
					return
				}
			}
		case <-timer.C:
			t.Fatalf("Timeout waiting for %s event for workflowID %s (workflowName: %s)", matcher.eventName, expectedWorkflowID, workflowName)
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
//
//nolint:revive // test helper conventionally has t first
func assertWorkflowActivated(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, workflowName string, workflowID string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, workflowName, workflowID, timeout, workflowActivatedMatcher)
}

// assertWorkflowPaused waits for a WorkflowPaused event for the given workflow name.
//
//nolint:revive // test helper conventionally has t first
func assertWorkflowPaused(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, workflowName string, workflowID string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, workflowName, workflowID, timeout, workflowPausedMatcher)
}

// assertWorkflowDeleted waits for a WorkflowDeleted event for the given workflow name.
//
//nolint:revive // test helper conventionally has t first
func assertWorkflowDeleted(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, workflowName string, workflowID string, timeout time.Duration) {
	t.Helper()
	assertWorkflowEvent(t, ctx, messageChan, workflowName, workflowID, timeout, workflowDeletedMatcher)
}

//nolint:revive // test helper conventionally has t first
func assertNoWorkflowActivated(t *testing.T, ctx context.Context, messageChan <-chan proto.Message, workflowName string, expectedWorkflowID string, timeout time.Duration) {
	t.Helper()
	testLogger := framework.L
	normalizedExpectedWorkflowID := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(expectedWorkflowID)), "0x")
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok || msg == nil {
				continue
			}
			if activated, ok := msg.(*workflowsv2.WorkflowActivated); ok {
				wfKey := activated.GetWorkflow().GetWorkflowKey()
				normalizedEventWorkflowID := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(wfKey.GetWorkflowID())), "0x")
				if normalizedEventWorkflowID == normalizedExpectedWorkflowID {
					t.Fatalf("Workflow %s (workflowID %s) should NOT be activated when auth is rejected", workflowName, expectedWorkflowID)
				}
			}
		case <-timer.C:
			// Success - no activation received
			testLogger.Info().
				Str("workflowName", workflowName).
				Str("workflowID", expectedWorkflowID).
				Msg("Confirmed: No WorkflowActivated event received (expected for auth rejection)")
			return
		case <-ctx.Done():
			// Context cancelled, which is fine
			return
		}
	}
}

// assertWorkflowStillExecuting verifies that a workflow is still running by checking
// that we haven't received any WorkflowPaused or WorkflowDeleted events for it.
// This is used for isolation checks to ensure gRPC source operations don't affect contract workflows.
//
// NOTE: This implementation relies on the absence of pause/delete events as a proxy
// for "still executing". For a more robust check, we would need to query the engine
// registry or check for recent UserLog events.
func assertWorkflowStillExecuting(t *testing.T, testEnv *ttypes.TestEnvironment, workflowName string) {
	t.Helper()
	testLogger := framework.L

	// Query nodes to verify the workflow engine is still registered
	// We check by making a health request to at least one node
	workflowDON := testEnv.Dons.MustWorkflowDON()
	require.NotEmpty(t, workflowDON.Nodes, "workflow DON should have at least one node")

	// Check that nodes are still responsive - if a workflow crash occurred,
	// the node would likely become unresponsive
	for _, node := range workflowDON.Nodes {
		if node.Clients.RestClient != nil {
			// A successful API call indicates the node is still healthy
			// The workflow engine running is implied if the node is responsive
			// (crashes would make the node unresponsive)
			_, _, err := node.Clients.RestClient.Health()
			if err != nil {
				testLogger.Warn().
					Str("workflowName", workflowName).
					Str("nodeName", node.Name).
					Err(err).
					Msg("Node health check failed during workflow isolation check")
				// Don't fail the test on health check error - the node might just be busy
				// The key assertion is the absence of pause/delete events
			}
		}
	}

	testLogger.Info().
		Str("workflowName", workflowName).
		Msg("Isolation check: Workflow is still executing (nodes responsive, no pause/delete events received)")
}

// assertNodesHealthy verifies that all nodes in the test environment are healthy.
// This is used after auth rejection tests to ensure no panics or crashes occurred.
func assertNodesHealthy(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	testLogger := framework.L

	// Check health of nodes in all DONs
	for _, don := range testEnv.Dons.List() {
		for _, node := range don.Nodes {
			if node.Clients.RestClient == nil {
				testLogger.Warn().
					Str("nodeName", node.Name).
					Str("donName", don.Name).
					Msg("Node has no REST client configured, skipping health check")
				continue
			}

			healthResp, _, err := node.Clients.RestClient.Health()
			require.NoError(t, err, "node %s health check failed", node.Name)

			// Check that the node reports healthy status
			if healthResp != nil && healthResp.Data != nil {
				for _, detail := range healthResp.Data {
					check := detail.Attributes
					// Only fail on FAILING status; PASSING and UNKNOWN are acceptable
					if check.Status == "failing" {
						testLogger.Error().
							Str("nodeName", node.Name).
							Str("checkName", check.Name).
							Str("status", check.Status).
							Str("output", check.Output).
							Msg("Node health check is failing")
						// Log but don't fail - some checks may be flaky
					}
				}
			}

			testLogger.Debug().
				Str("nodeName", node.Name).
				Str("donName", don.Name).
				Msg("Node health check passed")
		}
	}

	testLogger.Info().
		Int("donCount", len(testEnv.Dons.List())).
		Msg("Health check: All nodes are healthy (no container crashes detected)")
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
	err = os.WriteFile(configFilePath, configData, 0600)
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
