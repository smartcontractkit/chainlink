package cre

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	deployment_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	shard_config_changeset "github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
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

	testLogger.Info().Msg("Verifying Ring OCR Oracle health on shard0 nodes...")
	waitForRingOracleHealthy(t, shardZero)

	const numWorkflows = 10
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	var workflowIDs []string
	for i := 0; i < numWorkflows; i++ {
		workflowName := fmt.Sprintf("shardtest%d", i)
		workflowConfig := crontypes.WorkflowConfig{
			Schedule: "*/30 * * * * *",
		}
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed real workflows for sharding test")

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

	testLogger.Info().Msg("Reporting shard status to ALL nodes' Arbiters...")

	// TODO: we should modify arbiter not to report or report something else when health data is not available from Scaler
	initializeAllArbiterStates(t, testEnv, shardZero, len(shardDONs))

	// TODO: check edge case what happens when the GetWorkflowMappingsBatch returns mappings with missing workflows
	// TODO: check what happens if workflows is being removed, if there are no artifacts in the routing state (ringStore)
	validateShardingScaleScenario(t, testEnv, rpcHost, workflowIDs)

	testLogger.Info().Msg("Sharding test completed successfully")
}

func validateShardOrchestratorRPC(t *testing.T, logger zerolog.Logger, addr string) {
	t.Helper()

	logger.Info().Str("address", addr).Msg("Testing ShardOrchestrator RPC connectivity")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to create gRPC client for ShardOrchestrator at %s", addr)
	defer conn.Close()

	client := ringpb.NewShardOrchestratorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{"test-workflow-id"},
	})

	require.NoError(t, err, "ShardOrchestrator RPC call failed")
	require.NotNil(t, resp, "ShardOrchestrator response should not be nil")
	logger.Info().Int("mappingsCount", len(resp.Mappings)).Msg("ShardOrchestrator RPC responded successfully")
}

func initializeAllArbiterStates(t *testing.T, testEnv *ttypes.TestEnvironment, shardZero *cre.Don, numShards int) {
	t.Helper()
	logger := framework.L

	shardStatus := make(map[uint32]*ringpb.ShardStatus)
	for i := 0; i < numShards; i++ {
		shardStatus[uint32(i)] = &ringpb.ShardStatus{IsHealthy: true}
	}

	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name != shardZero.Name || nodeSet.Out == nil {
			continue
		}
		for _, clNode := range nodeSet.Out.CLNodes {
			parsedURL, parseErr := url.Parse(clNode.Node.ExternalURL)
			if parseErr != nil {
				logger.Warn().Err(parseErr).Str("url", clNode.Node.ExternalURL).Msg("Failed to parse node URL")
				continue
			}
			arbiterAddr := parsedURL.Hostname() + ":19876"

			conn, err := grpc.NewClient(arbiterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Warn().Err(err).Str("addr", arbiterAddr).Msg("Failed to connect to Arbiter")
				continue
			}

			client := ringpb.NewArbiterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			resp, err := client.GetDesiredReplicas(ctx, &ringpb.ShardStatusRequest{
				Status: shardStatus,
			})
			cancel()
			conn.Close()

			if err != nil {
				logger.Warn().Err(err).Str("addr", arbiterAddr).Msg("Failed to call GetDesiredReplicas")
			} else {
				logger.Info().
					Str("addr", arbiterAddr).
					Uint32("wantShards", resp.WantShards).
					Msg("Initialized Arbiter state")
			}
		}
	}
	logger.Info().Int("numShards", numShards).Msg("Arbiter states initialized on all shard0 nodes")
}

func validateShardingScaleScenario(t *testing.T, testEnv *ttypes.TestEnvironment, rpcHost string, workflowIDs []string) {
	t.Helper()
	logger := framework.L
	ctx := context.Background()

	shardConfigRef := getShardConfigRef(t, testEnv)
	chainSelector := testEnv.CreEnvironment.RegistryChainSelector

	arbiterClient := newArbiterClient(t, rpcHost+":19876")
	shardOrchClient := newShardOrchestratorClient(t, rpcHost+":60051")

	logger.Info().Msg("Diagnostic: Verifying store connection (direct registration)...")
	verifyStoreConnection(t, shardOrchClient)

	logger.Info().Msg("Diagnostic: Verifying Ring OCR rounds are completing...")
	waitForRingOCRRounds(t, shardOrchClient)

	logger.Info().Msg("Waiting for real workflows to be registered via Ring OCR consensus...")
	waitForWorkflowsRegistered(t, shardOrchClient, workflowIDs)

	logger.Info().Msg("Step 1: Verify all real workflows are mapped")
	resp, err := shardOrchClient.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Mappings, len(workflowIDs), "All deployed workflows should be mapped")
	logger.Info().Interface("mappings", resp.Mappings).Int("count", len(resp.Mappings)).Msg("Real workflows mapped")

	logger.Info().Msg("Step 2: Set ShardConfig to 1 shard (only shard-zero)")
	updateShardCount(t, testEnv, chainSelector, shardConfigRef, 1)

	logger.Info().Msg("Step 3: Verify Arbiter returns WantShards=1")
	waitForArbiterShardCount(t, arbiterClient, 1)

	logger.Info().Msg("Step 4: Verify all workflows mapped to shard 0 with single shard")
	resp, err = shardOrchClient.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
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

	logger.Info().Msg("Step 7: Wait for workflow redistribution after scaling")
	waitForWorkflowsDistributed(t, shardOrchClient, workflowIDs, 2)

	logger.Info().Msg("Step 8: Verify workflow mappings span 2 shards")
	resp, err = shardOrchClient.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
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
		Msg("Real workflows distributed across 2 shards after scaling")
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

func newShardOrchestratorClient(t *testing.T, addr string) ringpb.ShardOrchestratorServiceClient {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return ringpb.NewShardOrchestratorServiceClient(conn)
}

func waitForRingOracleHealthy(t *testing.T, shardZero *cre.Don) {
	t.Helper()
	logger := framework.L

	workers, err := shardZero.Workers()
	require.NoError(t, err, "Failed to get shard0 workers")
	require.NotEmpty(t, workers, "No worker nodes in shard0")

	node := workers[0]
	logger.Info().Str("node", node.Name).Msg("Checking available health components...")

	health, _, healthErr := node.Clients.RestClient.Health()
	if healthErr != nil {
		logger.Error().Err(healthErr).Msg("Failed to get health status")
	} else if health != nil && health.Data != nil {
		var ocrComponents []string
		for _, check := range health.Data {
			if strings.Contains(strings.ToLower(check.Attributes.Name), "ocr") ||
				strings.Contains(strings.ToLower(check.Attributes.Name), "oracle") ||
				strings.Contains(strings.ToLower(check.Attributes.Name), "ring") {
				ocrComponents = append(ocrComponents, check.Attributes.Name+" ("+check.Attributes.Status+")")
			}
		}
		if len(ocrComponents) > 0 {
			logger.Info().Strs("ocrComponents", ocrComponents).Msg("Found OCR-related health components")
		} else {
			logger.Warn().Int("totalComponents", len(health.Data)).Msg("No OCR/Oracle/Ring health components found - Ring Oracle likely not running")
			var allNames []string
			for _, check := range health.Data {
				allNames = append(allNames, check.Attributes.Name)
			}
			logger.Debug().Strs("allComponents", allNames).Msg("All available health components")
		}
	}

	logger.Info().Msg("Ring Oracle health check complete - if no OCR components found, Ring job may not be running")
}

func verifyStoreConnection(t *testing.T, client ringpb.ShardOrchestratorServiceClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testWorkflowID := "test-store-connection-workflow"
	_, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        0,
		RegisteredWorkflows:  map[string]uint32{testWorkflowID: 0},
		TotalActiveWorkflows: 1,
	})
	require.NoError(t, err, "Failed to register test workflow via ReportWorkflowTriggerRegistration")

	resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{testWorkflowID},
	})
	require.NoError(t, err, "Failed to get workflow mapping")
	require.Contains(t, resp.Mappings, testWorkflowID, "Store connection test: directly registered workflow should appear in mappings")
	framework.L.Info().
		Str("testWorkflowID", testWorkflowID).
		Uint64("mappingVersion", resp.MappingVersion).
		Msg("Store connection verified - direct registration works")
}

func waitForRingOCRRounds(t *testing.T, client ringpb.ShardOrchestratorServiceClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initialResp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{"probe-workflow"},
	})
	require.NoError(t, err)
	initialVersion := initialResp.MappingVersion

	framework.L.Info().
		Uint64("initialVersion", initialVersion).
		Msg("Waiting for Ring OCR rounds to complete (mapping_version should increase)")

	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"probe-workflow"},
		})
		if err != nil {
			framework.L.Warn().Err(err).Msg("Failed to get workflow mapping during Ring OCR check")
			return false
		}
		framework.L.Info().
			Uint64("currentVersion", resp.MappingVersion).
			Uint64("initialVersion", initialVersion).
			Int("mappingsCount", len(resp.Mappings)).
			Msg("Ring OCR round check")
		return resp.MappingVersion > initialVersion
	}, 90*time.Second, 5*time.Second, "Ring OCR rounds not completing - mapping_version not increasing. Initial: %d", initialVersion)
}

func waitForWorkflowsRegistered(t *testing.T, client ringpb.ShardOrchestratorServiceClient, workflowIDs []string) {
	t.Helper()
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: workflowIDs,
		})
		if err != nil {
			return false
		}
		framework.L.Info().
			Int("registered", len(resp.Mappings)).
			Int("expected", len(workflowIDs)).
			Uint64("mappingVersion", resp.MappingVersion).
			Msg("Waiting for workflows")
		return len(resp.Mappings) == len(workflowIDs)
	}, 2*time.Minute, 5*time.Second, "Workflows not registered within timeout")
}

func waitForWorkflowsDistributed(t *testing.T, client ringpb.ShardOrchestratorServiceClient, workflowIDs []string, minShards int) {
	t.Helper()
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: workflowIDs,
		})
		if err != nil {
			return false
		}
		shardsSeen := make(map[uint32]bool)
		for _, shardID := range resp.Mappings {
			shardsSeen[shardID] = true
		}
		framework.L.Info().Int("shardsUsed", len(shardsSeen)).Int("minShards", minShards).Msg("Waiting for distribution")
		return len(shardsSeen) >= minShards
	}, 2*time.Minute, 5*time.Second, "Workflows not distributed across %d shards within timeout", minShards)
}
