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

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	ring_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/sharding"
	stvault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteManualShardAssignmentTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for manual assignment test")

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"

	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "org_test_manual")

	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [1]
hashed_default_assignment = false

[per_org_assignment]
  org_test_manual = [0]
`)

	shardLeaderDON := getShardZeroDon(t, testEnv)

	proposeAndApproveShardAssignmentJob(t, testEnv, shardLeaderDON, shardAssignmentTOML, testLogger)

	const numWorkflows = 5
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowName := fmt.Sprintf("manualshard%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed workflows for manual shard assignment test")

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = 0
	}

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 3 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on correct shards")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on correct shards (manual-only mode)")
}

func ExecuteRingOCROverridesTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for override test")

	var shardZero *cre.Don
	for _, don := range shardDONs {
		if don.Metadata().IsShardLeader() {
			shardZero = don
			break
		}
	}
	require.NotNil(t, shardZero, "Expected to find shard zero DON")

	topology, tErr := cre.NewTopology(testEnv.Config.NodeSets, *testEnv.Config.Infra, testEnv.Config.CapabilityConfigs)
	require.NoError(t, tErr, "Failed to recreate topology")

	err := sharding.SetupSharding(t.Context(), sharding.SetupShardingInput{
		Logger:   testLogger,
		CreEnv:   testEnv.CreEnvironment,
		Topology: topology,
		Dons:     testEnv.Dons,
	})
	if err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Ring jobs already exist (from previous run), continuing...")
		} else {
			require.NoError(t, err, "SetupSharding failed")
		}
	} else {
		testLogger.Info().Msg("SetupSharding completed successfully")
	}

	waitForRingOracleHealthy(t, shardZero)

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "org_test_override")

	shardAssignmentTOML := `
static_default_assignment = [0]
hashed_default_assignment = true

[per_org_assignment]
  org_test_override = [1]
`

	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, testLogger)
	}

	const numWorkflows = 5
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowName := fmt.Sprintf("override-shard%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed workflows for ringocr-with-overrides test")

	var rpcHost string
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name == "shard0" && nodeSet.Out != nil && len(nodeSet.Out.CLNodes) > 0 {
			externalURL := nodeSet.Out.CLNodes[0].Node.ExternalURL
			parsedURL, parseErr := url.Parse(externalURL)
			require.NoError(t, parseErr, "Failed to parse ExternalURL")
			rpcHost = parsedURL.Hostname()
			break
		}
	}
	require.NotEmpty(t, rpcHost, "Failed to find shard0 node set to extract RPC host")

	shardOrchClient := newShardOrchestratorClient(t, rpcHost+":60051")

	testLogger.Info().Msg("Reporting shard status to ALL nodes' Arbiters...")
	initializeAllArbiterStates(t, testEnv, shardZero, len(shardDONs))

	testLogger.Info().Msg("Diagnostic: Verifying store connection (direct registration)...")
	verifyStoreConnection(t, shardOrchClient)

	testLogger.Info().Msg("Diagnostic: Verifying Ring OCR rounds are completing...")
	waitForRingOCRRounds(t, shardOrchClient)

	testLogger.Info().Msg("Waiting for workflows to be registered via Ring OCR...")
	waitForWorkflowsRegistered(t, shardOrchClient, workflowIDs)

	resp, err := shardOrchClient.GetWorkflowShardMapping(t.Context(), &ringpb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	testLogger.Info().Interface("mappings", resp.Mappings).Msg("Ring OCR workflow mappings")
	require.Len(t, resp.Mappings, len(workflowIDs), "All deployed workflows should be mapped")

	const overrideShard uint32 = 1

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = overrideShard
	}

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 5 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on correct shards")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on correct shards (ringocr-with-overrides mode)")
}

func proposeAndApproveShardAssignmentJob(t *testing.T, testEnv *ttypes.TestEnvironment, targetDON *cre.Don, shardAssignmentTOML string, testLogger zerolog.Logger) {
	t.Helper()

	jobInput := ring_ops.ProposeShardAssignmentJobInput{
		Domain:          offchain.ProductLabel,
		Environment:     testEnv.CreEnvironment.CldfEnvironment.Name,
		DONName:         targetDON.Name,
		ShardAssignment: shardAssignmentTOML,
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: targetDON.Name},
		},
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: "shard-assignment"},
	}

	report, err := operations.ExecuteOperation(
		testEnv.CreEnvironment.CldfEnvironment.OperationsBundle,
		ring_ops.ProposeShardAssignmentJob,
		ring_ops.ProposeShardAssignmentJobDeps{Env: *testEnv.CreEnvironment.CldfEnvironment},
		jobInput,
	)
	if err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Shard assignment job already exists (from previous run), continuing...")
			return
		}
		require.NoError(t, err, "Failed to propose shard assignment job")
	}

	if err := jobs.Approve(t.Context(), testEnv.CreEnvironment.CldfEnvironment.Offchain, testEnv.Dons, report.Output.Specs); err != nil {
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			testLogger.Info().Msg("Shard assignment job already approved (from previous run), continuing...")
			return
		}
		require.NoError(t, err, "Failed to approve shard assignment job")
	}

	testLogger.Info().Msg("Shard assignment job proposed and approved")
}
