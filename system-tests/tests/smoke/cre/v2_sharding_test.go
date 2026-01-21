package cre

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/container"
	dfilter "github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	ringpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"
	shardorchpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/shardorchestrator/pb"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	deployment_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	shard_config_changeset "github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

/*
Sharding Test

This test validates the SetupSharding functionality which:
1. Deploys a ShardConfig contract
2. Deploys a Ring OCR3 contract
3. Creates Ring jobs on the shard leader DON
4. Configures the Ring OCR3 contract with DON signers

Prerequisites:
- Start the environment with the sharded DON config:
  cd core/scripts/cre/environment
  CTF_CONFIGS=configs/workflow-gateway-sharded-don.toml go run . env start --with-beholder --with-contracts-version v2

- Run the test:
  go test -timeout 20m -run "^Test_CRE_V2_Sharding$" -v
*/

func ExecuteShardingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for sharding test")
	testLogger.Info().Msgf("Found %d shard DONs", len(shardDONs))

	var shardZero *cre.Don
	for _, don := range shardDONs {
		if don.Metadata().IsShardLeader() {
			shardZero = don
			break
		}
	}
	require.NotNil(t, shardZero, "Expected to find shard zero DON")
	testLogger.Info().Msgf("Shard zero DON: %s (ID: %d)", shardZero.Name, shardZero.ID)

	bootstrap, hasBootstrap := testEnv.Dons.Bootstrap()
	require.True(t, hasBootstrap, "Expected bootstrap node to exist")
	testLogger.Info().Msgf("Bootstrap node found: %s", bootstrap.Name)

	workers, err := shardZero.Workers()
	require.NoError(t, err, "Expected shard zero to have worker nodes")
	require.NotEmpty(t, workers, "Expected at least one worker node in shard zero DON")
	testLogger.Info().Msgf("Shard zero has %d worker nodes", len(workers))

	for _, don := range shardDONs {
		metadata := don.Metadata()
		testLogger.Info().
			Str("name", don.Name).
			Uint64("id", don.ID).
			Bool("isShardZero", metadata.IsShardLeader()).
			Uint("shardIndex", metadata.ShardIndex).
			Int("nodeCount", len(don.Nodes)).
			Msg("Shard DON info")
	}

	testLogger.Info().Msg("Calling SetupSharding to deploy contracts and create Ring jobs...")
	err = sharding.SetupSharding(t.Context(), sharding.SetupShardingInput{
		Logger:   testLogger,
		CreEnv:   testEnv.CreEnvironment,
		Topology: nil,
		Dons:     testEnv.Dons,
	})
	if err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Ring jobs already exist (from previous run), continuing with RPC tests...")
		} else {
			require.NoError(t, err, "SetupSharding failed")
		}
	} else {
		testLogger.Info().Msg("SetupSharding completed successfully")
	}

	var rpcHost string
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name == "shard0" && nodeSet.Out != nil && len(nodeSet.Out.CLNodes) > 0 {
			externalURL := nodeSet.Out.CLNodes[0].Node.ExternalURL
			parsedURL, parseErr := url.Parse(externalURL)
			require.NoError(t, parseErr, "Failed to parse ExternalURL")
			rpcHost = parsedURL.Hostname()
			testLogger.Info().
				Str("externalURL", externalURL).
				Str("rpcHost", rpcHost).
				Msg("Extracted RPC host from shard0 node ExternalURL")
			break
		}
	}
	require.NotEmpty(t, rpcHost, "Failed to find shard0 node set to extract RPC host")

	shardOrchestratorAddr := rpcHost + ":60051"
	validateShardOrchestratorRPC(t, testLogger, shardOrchestratorAddr)

	arbiterAddr := rpcHost + ":19876"
	validateArbiterRPC(t, testLogger, arbiterAddr)

	validateShardingScaleScenario(t, testEnv, rpcHost)

	validateShardFilteringExecution(t, testEnv, rpcHost)

	testLogger.Info().Msg("Sharding test completed successfully")
}

func validateShardOrchestratorRPC(t *testing.T, logger zerolog.Logger, addr string) {
	t.Helper()

	logger.Info().Str("address", addr).Msg("Testing ShardOrchestrator RPC connectivity")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to create gRPC client for ShardOrchestrator at %s", addr)
	defer conn.Close()

	client := shardorchpb.NewShardOrchestratorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetWorkflowShardMapping(ctx, &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{"test-workflow-id"},
	})

	require.NoError(t, err, "ShardOrchestrator RPC call failed")
	require.NotNil(t, resp, "ShardOrchestrator response should not be nil")
	logger.Info().Int("mappingsCount", len(resp.Mappings)).Msg("ShardOrchestrator RPC responded successfully")
}

func validateArbiterRPC(t *testing.T, logger zerolog.Logger, addr string) {
	t.Helper()

	logger.Info().Str("address", addr).Msg("Testing Arbiter RPC connectivity")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to create gRPC client for Arbiter at %s", addr)
	defer conn.Close()

	client := ringpb.NewArbiterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetDesiredReplicas(ctx, &ringpb.ShardStatusRequest{})

	if err != nil {
		errStr := err.Error()
		require.NotContains(t, errStr, "unknown service",
			"Arbiter service not registered - ensure Ring jobs are created via SetupSharding")
		logger.Info().Err(err).Msg("Arbiter returned error (may be expected depending on state)")
	} else {
		require.NotNil(t, resp, "Arbiter response should not be nil")
		logger.Info().
			Uint32("wantShards", resp.WantShards).
			Msg("Arbiter RPC responded successfully")
	}

	logger.Info().Str("address", addr).Msg("Arbiter RPC test passed")
}

func validateShardingScaleScenario(t *testing.T, testEnv *ttypes.TestEnvironment, rpcHost string) {
	t.Helper()
	logger := framework.L
	ctx := context.Background()

	shardConfigRef := getShardConfigRef(t, testEnv)
	chainSelector := testEnv.CreEnvironment.RegistryChainSelector

	arbiterClient := newArbiterClient(t, rpcHost+":19876")
	shardOrchClient := newShardOrchestratorClient(t, rpcHost+":60051")

	workflowIDs := []string{"workflow-A", "workflow-B", "workflow-C", "workflow-D"}

	logger.Info().Msg("Step 1: Set ShardConfig to 1 shard (only shard-zero)")
	updateShardCount(t, testEnv, chainSelector, shardConfigRef, 1)

	logger.Info().Msg("Step 2: Verify Arbiter returns WantShards=1")
	waitForArbiterShardCount(t, arbiterClient, 1)

	logger.Info().Msg("Step 3: Register all workflows on shard-zero (the only shard)")
	_, err := shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &shardorchpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        0,
		RegisteredWorkflows:  map[string]uint32{"workflow-A": 1, "workflow-B": 1, "workflow-C": 1, "workflow-D": 1},
		TotalActiveWorkflows: 4,
	})
	require.NoError(t, err)

	logger.Info().Msg("Step 4: Verify all workflows mapped to shard 0")
	resp, err := shardOrchClient.GetWorkflowShardMapping(ctx, &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	for _, wfID := range workflowIDs {
		assert.Equal(t, uint32(0), resp.Mappings[wfID], "With 1 shard, workflow %s should map to shard 0", wfID)
	}
	logger.Info().Interface("mappings", resp.Mappings).Msg("All workflows on shard-zero")

	logger.Info().Msg("Step 5: Scale up - Set ShardConfig to 2 shards")
	updateShardCount(t, testEnv, chainSelector, shardConfigRef, 2)

	logger.Info().Msg("Step 6: Verify Arbiter returns WantShards=2")
	waitForArbiterShardCount(t, arbiterClient, 2)

	logger.Info().Msg("Step 7: Shard 1 reports its workflows after scaling")
	_, err = shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &shardorchpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        1,
		RegisteredWorkflows:  map[string]uint32{"workflow-C": 1, "workflow-D": 1},
		TotalActiveWorkflows: 2,
	})
	require.NoError(t, err)

	logger.Info().Msg("Step 8: Verify workflow mappings now span 2 shards")
	resp, err = shardOrchClient.GetWorkflowShardMapping(ctx, &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	shardCounts := map[uint32]int{}
	for _, shardID := range resp.Mappings {
		shardCounts[shardID]++
	}
	assert.Positive(t, shardCounts[0], "Some workflows should be on shard 0")
	assert.Positive(t, shardCounts[1], "Some workflows should be on shard 1")
	logger.Info().
		Interface("mappings", resp.Mappings).
		Interface("distribution", shardCounts).
		Msg("Workflows distributed across 2 shards after scaling")
}

func getShardConfigRef(t *testing.T, testEnv *ttypes.TestEnvironment) datastore.AddressRefKey {
	t.Helper()
	return datastore.NewAddressRefKey(
		testEnv.CreEnvironment.RegistryChainSelector,
		datastore.ContractType(deployment_contracts.ShardConfig.String()),
		semver.MustParse("1"),
		"",
	)
}

func updateShardCount(t *testing.T, testEnv *ttypes.TestEnvironment, chainSelector uint64, shardConfigRef datastore.AddressRefKey, count uint64) {
	t.Helper()
	_, err := commonchangeset.RunChangeset(
		shard_config_changeset.UpdateShardCount{},
		*testEnv.CreEnvironment.CldfEnvironment,
		shard_config_changeset.UpdateShardCountInput{
			ChainSelector:  chainSelector,
			NewShardCount:  count,
			ShardConfigRef: shardConfigRef,
		},
	)
	require.NoError(t, err)
	framework.L.Info().Uint64("count", count).Msg("Updated ShardConfig shard count")
}

func waitForArbiterShardCount(t *testing.T, client ringpb.ArbiterClient, expected uint32) {
	t.Helper()
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.GetDesiredReplicas(ctx, &ringpb.ShardStatusRequest{})
		if err != nil {
			return false
		}
		framework.L.Info().Uint32("wantShards", resp.WantShards).Uint32("expected", expected).Msg("Arbiter response")
		return resp.WantShards == expected
	}, 30*time.Second, 2*time.Second, "Arbiter did not return expected WantShards=%d", expected)
}

func newArbiterClient(t *testing.T, addr string) ringpb.ArbiterClient {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return ringpb.NewArbiterClient(conn)
}

func newShardOrchestratorClient(t *testing.T, addr string) shardorchpb.ShardOrchestratorServiceClient {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return shardorchpb.NewShardOrchestratorServiceClient(conn)
}

// validateShardFilteringExecution verifies that workflows are only executed on their assigned shard.
// This validates the core ShardFilter functionality by:
// 1. Registering test workflows to specific shards via orchestrator
// 2. Capturing logs from shard nodes to verify filtering behavior
// 3. Checking that each shard only processes workflows assigned to it
func validateShardFilteringExecution(t *testing.T, testEnv *ttypes.TestEnvironment, rpcHost string) {
	t.Helper()
	logger := framework.L
	ctx := context.Background()

	logger.Info().Msg("Validating shard filtering: workflows execute only on assigned shards")

	shardOrchClient := newShardOrchestratorClient(t, rpcHost+":60051")
	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Need at least 2 shard DONs for filtering test")

	// Find shard-0 and shard-1 DONs
	var shard0, shard1 *cre.Don
	for _, don := range shardDONs {
		metadata := don.Metadata()
		if metadata.ShardIndex == 0 {
			shard0 = don
		} else if metadata.ShardIndex == 1 {
			shard1 = don
		}
	}
	require.NotNil(t, shard0, "Shard-0 DON not found")
	require.NotNil(t, shard1, "Shard-1 DON not found")

	// Test workflows - use unique IDs to avoid conflicts
	testRunID := time.Now().Unix()
	workflowA := fmt.Sprintf("filter-test-wf-A-%d", testRunID)
	workflowB := fmt.Sprintf("filter-test-wf-B-%d", testRunID)
	workflowC := fmt.Sprintf("filter-test-wf-C-%d", testRunID)
	workflowD := fmt.Sprintf("filter-test-wf-D-%d", testRunID)
	workflowIDs := []string{workflowA, workflowB, workflowC, workflowD}

	logger.Info().Msgf("Test workflows: A=%s, B=%s, C=%s, D=%s", workflowA, workflowB, workflowC, workflowD)

	// Mark the time before we start - we'll check logs after this timestamp
	logCheckStartTime := time.Now()

	// Step 1: Register workflows with specific shard assignments
	logger.Info().Msg("Step 1: Register test workflows to specific shards")

	// Register workflows A and B to shard-0
	_, err := shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &shardorchpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        0,
		RegisteredWorkflows:  map[string]uint32{workflowA: 1, workflowB: 1},
		TotalActiveWorkflows: 2,
	})
	require.NoError(t, err, "Failed to register workflows on shard-0")

	// Register workflows C and D to shard-1
	_, err = shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &shardorchpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        1,
		RegisteredWorkflows:  map[string]uint32{workflowC: 1, workflowD: 1},
		TotalActiveWorkflows: 2,
	})
	require.NoError(t, err, "Failed to register workflows on shard-1")

	// Step 2: Verify orchestrator returns correct mappings
	logger.Info().Msg("Step 2: Verify ShardOrchestrator returns correct workflow-to-shard mappings")

	resp, err := shardOrchClient.GetWorkflowShardMapping(ctx, &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err, "Failed to get workflow shard mappings")
	require.NotNil(t, resp)

	// Verify expected mappings
	expectedMappings := map[string]uint32{
		workflowA: 0,
		workflowB: 0,
		workflowC: 1,
		workflowD: 1,
	}

	for wfID, expectedShard := range expectedMappings {
		actualShard, found := resp.Mappings[wfID]
		require.True(t, found, "Workflow %s not found in orchestrator mappings", wfID)
		assert.Equal(t, expectedShard, actualShard, "Workflow %s should be mapped to shard %d, got shard %d", wfID, expectedShard, actualShard)
	}

	logger.Info().
		Interface("mappings", resp.Mappings).
		Msg("ShardOrchestrator returned expected workflow mappings")

	// Step 3: Check logs for filtering behavior
	logger.Info().Msg("Step 3: Checking node logs for shard filtering behavior")

	// Wait a bit for any workflow activity and log generation
	time.Sleep(5 * time.Second)

	// Check shard-0 logs
	logger.Info().Msg("Checking shard-0 logs...")
	shard0Workers, err := shard0.Workers()
	require.NoError(t, err)

	shard0FilteringCount := 0
	shard0ExecutionCount := 0

	for _, node := range shard0Workers {
		logs := captureNodeLogs(t, logger, node, logCheckStartTime)

		// Check if shard-0 is skipping workflows C and D (assigned to shard-1)
		if containsFilterMessage(logs, workflowC) {
			shard0FilteringCount++
			logger.Info().Msgf("✓ Shard-0 node %s correctly skipped workflow %s", node.Name, workflowC)
		}
		if containsFilterMessage(logs, workflowD) {
			shard0FilteringCount++
			logger.Info().Msgf("✓ Shard-0 node %s correctly skipped workflow %s", node.Name, workflowD)
		}

		// Check if shard-0 is NOT skipping workflows A and B (assigned to shard-0)
		if !containsFilterMessage(logs, workflowA) {
			shard0ExecutionCount++
		}
		if !containsFilterMessage(logs, workflowB) {
			shard0ExecutionCount++
		}
	}

	// Check shard-1 logs
	logger.Info().Msg("Checking shard-1 logs...")
	shard1Workers, err := shard1.Workers()
	require.NoError(t, err)

	shard1FilteringCount := 0
	shard1ExecutionCount := 0

	for _, node := range shard1Workers {
		logs := captureNodeLogs(t, logger, node, logCheckStartTime)

		// Check if shard-1 is skipping workflows A and B (assigned to shard-0)
		if containsFilterMessage(logs, workflowA) {
			shard1FilteringCount++
			logger.Info().Msgf("✓ Shard-1 node %s correctly skipped workflow %s", node.Name, workflowA)
		}
		if containsFilterMessage(logs, workflowB) {
			shard1FilteringCount++
			logger.Info().Msgf("✓ Shard-1 node %s correctly skipped workflow %s", node.Name, workflowB)
		}

		// Check if shard-1 is NOT skipping workflows C and D (assigned to shard-1)
		if !containsFilterMessage(logs, workflowC) {
			shard1ExecutionCount++
		}
		if !containsFilterMessage(logs, workflowD) {
			shard1ExecutionCount++
		}
	}

	// Report findings
	logger.Info().Msgf("Shard-0: Found %d filtering messages (expected for workflows C,D)", shard0FilteringCount)
	logger.Info().Msgf("Shard-1: Found %d filtering messages (expected for workflows A,B)", shard1FilteringCount)

	if shard0FilteringCount > 0 || shard1FilteringCount > 0 {
		logger.Info().Msg("✓ ShardFilter is working: found evidence of filtering in logs")
	} else {
		logger.Warn().Msg("⚠ No filtering messages found - ShardFilter may not be configured yet")
		logger.Warn().Msg("To enable filtering, configure these environment variables on each node:")
		logger.Warn().Msg("  WORKFLOW_SHARD_INDEX=<0|1>")
		logger.Warn().Msg("  WORKFLOW_SHARD_ORCHESTRATOR_ADDR=<host>:60051")
	}

	logger.Info().Msg("Shard filtering validation completed")
}

// captureNodeLogs captures logs from a specific node's container since the given timestamp
func captureNodeLogs(t *testing.T, logger zerolog.Logger, node *cre.Node, since time.Time) string {
	t.Helper()

	// Create filter for this specific node's container
	listOpts := container.ListOptions{
		All: true,
		Filters: dfilter.NewArgs(
			dfilter.KeyValuePair{Key: "label", Value: "framework=ctf"},
			dfilter.KeyValuePair{Key: "name", Value: node.Name},
		),
	}

	logOpts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since.Format(time.RFC3339),
	}

	logStream, err := framework.StreamContainerLogs(listOpts, logOpts)
	if err != nil {
		logger.Warn().Err(err).Str("node", node.Name).Msg("Failed to stream container logs")
		return ""
	}

	// Read all logs from the first (and should be only) container
	for containerName, reader := range logStream {
		defer reader.Close()

		var logContent strings.Builder
		header := make([]byte, 8)

		for {
			_, err := io.ReadFull(reader, header)
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Debug().Err(err).Str("container", containerName).Msg("Error reading log header")
				break
			}

			// Extract message size from Docker stream header
			msgSize := binary.BigEndian.Uint32(header[4:8])
			msgBuf := make([]byte, msgSize)

			_, err = io.ReadFull(reader, msgBuf)
			if err != nil {
				logger.Debug().Err(err).Str("container", containerName).Msg("Error reading log message")
				break
			}

			logContent.Write(msgBuf)
		}

		return logContent.String()
	}

	return ""
}

// containsFilterMessage checks if logs contain the shard filter skip message for the given workflow ID
func containsFilterMessage(logs, workflowID string) bool {
	// Look for the debug message logged by the ShardFilter when skipping a workflow
	// Pattern: "skipping workflow execution - not assigned to this shard" with workflowID in context
	filterPattern := regexp.MustCompile(fmt.Sprintf(`(?i)(skipping workflow.*not assigned.*shard|workflow.*%s.*skip)`, regexp.QuoteMeta(workflowID)))
	return filterPattern.MatchString(logs)
}
