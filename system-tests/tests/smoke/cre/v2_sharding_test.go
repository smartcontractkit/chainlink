package cre

import (
	"context"
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
  CTF_CONFIGS=configs/workflow-gateway-sharded-don.toml go run . env start --with-contracts-version v2

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
	_, err := shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        0,
		RegisteredWorkflows:  map[string]uint32{"workflow-A": 1, "workflow-B": 1, "workflow-C": 1, "workflow-D": 1},
		TotalActiveWorkflows: 4,
	})
	require.NoError(t, err)

	logger.Info().Msg("Step 4: Verify all workflows mapped to shard 0")
	resp, err := shardOrchClient.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
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
	_, err = shardOrchClient.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
		SourceShardId:        1,
		RegisteredWorkflows:  map[string]uint32{"workflow-C": 1, "workflow-D": 1},
		TotalActiveWorkflows: 2,
	})
	require.NoError(t, err)

	logger.Info().Msg("Step 8: Verify workflow mappings now span 2 shards")
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

func newShardOrchestratorClient(t *testing.T, addr string) ringpb.ShardOrchestratorServiceClient {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return ringpb.NewShardOrchestratorServiceClient(conn)
}
