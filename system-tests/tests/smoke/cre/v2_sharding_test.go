package cre

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	ringpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/ring/pb"
	shardorchpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/shardorchestrator/pb"

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
	require.Greater(t, len(workers), 0, "Expected at least one worker node in shard zero DON")
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
	err = sharding.SetupSharding(sharding.SetupShardingInput{
		Ctx:      t.Context(),
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

	shardOrchestratorAddr := fmt.Sprintf("%s:60051", rpcHost)
	testShardOrchestratorRPC(t, testLogger, shardOrchestratorAddr)

	arbiterAddr := fmt.Sprintf("%s:19876", rpcHost)
	testArbiterRPC(t, testLogger, arbiterAddr)

	testLogger.Info().Msg("Sharding test completed successfully")
}

func testShardOrchestratorRPC(t *testing.T, logger zerolog.Logger, addr string) {
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

func testArbiterRPC(t *testing.T, logger zerolog.Logger, addr string) {
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
