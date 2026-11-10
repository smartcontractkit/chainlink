package cre

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	// evmLogTriggerUserLogNeedle is logged by the evm/logtrigger workflow once per handled log event.
	evmLogTriggerUserLogNeedle = "OnTrigger decoded message"

	// failedFailoverObserveWindow must outlast the failover window plus one event
	// emission interval, so a failover that could happen has happened by its end.
	failedFailoverObserveWindow = 75 * time.Second
)

// ExecuteShardFailoverEVMLogTriggerTest covers timeout-based failover of an EVM
// log trigger on a topology where every workflow shard is served by its own
// capability DON and capability shard 1 is broken (it hosts no EVM capability):
// configs/workflow-sharded-capabilities-failover-broken-cap.toml.
//
//  1. Assigned [1,0]: the primary (shard-1) can never register or receive the log
//     trigger, so it stays silent. The healthy secondary (shard-0 / cap-0) caches
//     the log events and, once the failover window passes with no outcome from the
//     primary, executes them itself.
//  2. Reassigned [0,1]: the healthy shard-0 executes as primary, then dies silently.
//     The failover target (shard-1) is blind (broken cap-1), so the failover cannot
//     succeed: no shard may execute the workflow again, and shard-1 must not report
//     having taken over. This is the "failed failover" case.
func ExecuteShardFailoverEVMLogTriggerTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv, "configs/workflow-sharded-capabilities-failover-broken-cap.toml")
	ownerKey := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.MustGetRootPrivateKey()
	shards.workflowOwner = strings.ToLower(crypto.PubkeyToAddress(ownerKey.PublicKey).Hex())

	healthyCapNodeSet := mustNodeSetByName(t, testEnv, "chain-capabilities-zone-a-shard-0")
	logTriggerChain := evmChainEnabledOnNodeSet(t, testEnv, healthyCapNodeSet)
	logTriggerChainID := logTriggerChain.CtfOutput().ChainID
	workflowConfig, msgEmitter := configureEVMLogTriggerWorkflow(t, testLogger, logTriggerChain)

	// The broken shard may fail to init its engine on purpose, so init errors must not fail the test.
	collector := startShardUserLogCollectorMatching(t, testEnv, evmLogTriggerUserLogNeedle, false)

	t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(shardTimeoutBasedFailoverOnTOML))

	// Phase 1: [1,0] - broken primary, healthy secondary.
	primaryOneTOML := shardFailoverAssignmentTOML(shards, false)
	proposeSharedVaultAssignment(t, testEnv, primaryOneTOML)
	requireAssignmentLandedOnShard(t, testEnv, shards.shardZeroDON, shardAssignmentUpdateHashNeedle(primaryOneTOML))
	requireAssignmentLandedOnShard(t, testEnv, shards.shardOneDON, shardAssignmentUpdateHashNeedle(primaryOneTOML))

	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger,
		t_helpers.UniqueWorkflowName(testEnv, "failover-evmlogtrigger"), &workflowConfig, "./evm/logtrigger/main.go")

	emitCtx, emitCancel := context.WithCancel(t.Context())
	defer emitCancel()
	startEVMLogTriggerEventEmitter(emitCtx, t, testLogger, logTriggerChainID, logTriggerChain, msgEmitter,
		"Data for failover log trigger chain "+logTriggerChainID)

	phase1Start := time.Now()
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle, shardAutoExecAwaitTimeout, 5*time.Second)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, shardAutoFailoverExecutedLogNeedle, shardAutoExecAwaitTimeout, 5*time.Second)
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardZeroDON, shardTickRecords, phase1Start, shardAutoExecAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardOneDON}, phase1Start)
	collector.requireNoDuplicateExecutionIDsSince(t, phase1Start)
	testLogger.Info().Msg("Phase 1: broken primary shard-1 stayed silent, healthy secondary shard-0 executed the cached log events by timeout-based failover")

	// Phase 2: [0,1] - healthy primary dies, the failover target is blind.
	primaryZeroTOML := shardFailoverAssignmentTOML(shards, true)
	proposeSharedVaultAssignment(t, testEnv, primaryZeroTOML)
	requireAssignmentLandedOnShard(t, testEnv, shards.shardZeroDON, shardAssignmentUpdateHashNeedle(primaryZeroTOML))
	requireAssignmentLandedOnShard(t, testEnv, shards.shardOneDON, shardAssignmentUpdateHashNeedle(primaryZeroTOML))

	landed := time.Now()
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardZeroDON, shardTickRecords, landed, shardEventAwaitTimeout)
	testLogger.Info().Msg("Phase 2a: shard-0 executes as the healthy primary")

	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	postDeath := death.Add(shardInflightDrainGrace)
	time.Sleep(shardInflightDrainGrace + failedFailoverObserveWindow)

	collector.requireNoRecordsForWorkflows(t, []string{workflowID}, postDeath)
	for _, line := range t_helpers.ContainerLogLinesForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle) {
		require.NotContainsf(t, line, workflowID, "shard-1 reported taking over workflow %s although its capability shard is broken", workflowID)
	}
	for _, line := range t_helpers.ContainerLogLinesForNodeset(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle) {
		require.NotContainsf(t, line, workflowID, "shard-1 cached a trigger event for workflow %s although its capability shard is broken", workflowID)
	}
	testLogger.Info().Msg("Phase 2b: primary died silently and the failover target is blind - the failover failed, nothing executed or took over")
}

func mustNodeSetByName(t *testing.T, testEnv *ttypes.TestEnvironment, name string) *cre.NodeSet {
	t.Helper()

	idx := slices.IndexFunc(testEnv.Config.NodeSets, func(ns *cre.NodeSet) bool { return ns.Name == name })
	require.GreaterOrEqualf(t, idx, 0, "node set %q not found; this test requires configs/workflow-sharded-capabilities-failover-broken-cap.toml", name)
	return testEnv.Config.NodeSets[idx]
}
