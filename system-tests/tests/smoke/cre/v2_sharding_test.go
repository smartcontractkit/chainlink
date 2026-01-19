package cre

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	ringpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"
	shardorchpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/shardorchestrator/pb"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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
  CTF_CONFIGS=configs/workflow-gateway-sharded-don.toml go run . env start

- Run the test:
  go test -timeout 20m -run "^Test_CRE_V2_Sharding$" -v
*/

func Test_CRE_V2_Sharding(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-don.toml"),
		"--with-contracts-version", "v2",
	)

	ExecuteShardingTest(t, testEnv)
}

func ExecuteShardingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	// WIP: this is just an initial draft of the sharding test
	testLogger := framework.L

	// Verify sharding DONs exist
	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for sharding test")
	testLogger.Info().Msgf("Found %d shard DONs", len(shardDONs))

	// Find shard zero DON (ShardIndex == 0)
	var shardZero *cre.Don
	for _, don := range shardDONs {
		if don.Metadata().IsShardLeader() {
			shardZero = don
			break
		}
	}
	require.NotNil(t, shardZero, "Expected to find shard zero DON")
	testLogger.Info().Msgf("Shard zero DON: %s (ID: %d)", shardZero.Name, shardZero.ID)

	// Verify bootstrap node exists
	bootstrap, hasBootstrap := testEnv.Dons.Bootstrap()
	require.True(t, hasBootstrap, "Expected bootstrap node to exist")
	testLogger.Info().Msgf("Bootstrap node found: %s", bootstrap.Name)

	// Verify shard zero DON has worker nodes
	workers, err := shardZero.Workers()
	require.NoError(t, err, "Expected shard zero to have worker nodes")
	require.Greater(t, len(workers), 0, "Expected at least one worker node in shard zero DON")
	testLogger.Info().Msgf("Shard zero has %d worker nodes", len(workers))

	// Log information about all shard DONs
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

	// Test ShardOrchestrator and Arbiter gRPC services
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	// Get shard zero worker node to test RPC services
	shardZeroWorker := workers[0]
	shardZeroMetadata := shardZeroWorker.Metadata()

	// Test ShardOrchestrator RPC (port 50051)
	testShardOrchestratorRPC(t, ctx, shardZeroMetadata.ShardOrchestratorAddress(), testLogger)

	// Test Arbiter RPC (port 9876)
	arbiterAddr := fmt.Sprintf("%s:%d", shardZeroMetadata.Host, 9876)
	testArbiterRPC(t, ctx, arbiterAddr, testLogger)

	testLogger.Info().Msg("Sharding test completed successfully - SetupSharding was executed and RPCs are working")
}

// testShardOrchestratorRPC verifies the ShardOrchestrator gRPC service is available and responding
func testShardOrchestratorRPC(t *testing.T, ctx context.Context, addr string, logger zerolog.Logger) {
	logger.Info().Str("addr", addr).Msg("Testing ShardOrchestrator RPC")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to connect to ShardOrchestrator")
	defer conn.Close()

	client := shardorchpb.NewShardOrchestratorServiceClient(conn)

	// Call GetWorkflowShardMapping with empty workflow IDs to verify service is responding
	req := &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{},
	}

	resp, err := client.GetWorkflowShardMapping(ctx, req)
	require.NoError(t, err, "ShardOrchestrator.GetWorkflowShardMapping RPC failed")
	require.NotNil(t, resp, "ShardOrchestrator response should not be nil")

	logger.Info().
		Int("mappingsCount", len(resp.Mappings)).
		Uint64("mappingVersion", resp.MappingVersion).
		Msg("ShardOrchestrator RPC successful")
}

// testArbiterRPC verifies the Arbiter gRPC service is available and responding
func testArbiterRPC(t *testing.T, ctx context.Context, addr string, logger zerolog.Logger) {
	logger.Info().Str("addr", addr).Msg("Testing Arbiter RPC")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to connect to Arbiter")
	defer conn.Close()

	// Test ArbiterScaler.Status RPC
	scalerClient := ringpb.NewArbiterScalerClient(conn)

	statusResp, err := scalerClient.Status(ctx, &emptypb.Empty{})
	require.NoError(t, err, "ArbiterScaler.Status RPC failed")
	require.NotNil(t, statusResp, "ArbiterScaler.Status response should not be nil")

	logger.Info().
		Uint32("wantShards", statusResp.WantShards).
		Int("statusCount", len(statusResp.Status)).
		Msg("ArbiterScaler.Status RPC successful")

	// Test Arbiter.GetDesiredReplicas RPC
	arbiterClient := ringpb.NewArbiterClient(conn)

	replicasResp, err := arbiterClient.GetDesiredReplicas(ctx, &ringpb.ShardStatusRequest{
		Status: map[uint32]*ringpb.ShardStatus{},
	})
	require.NoError(t, err, "Arbiter.GetDesiredReplicas RPC failed")
	require.NotNil(t, replicasResp, "Arbiter.GetDesiredReplicas response should not be nil")

	logger.Info().
		Uint32("wantShards", replicasResp.WantShards).
		Msg("Arbiter.GetDesiredReplicas RPC successful")
}
