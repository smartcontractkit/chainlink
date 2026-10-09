package cre

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// Automatic failover, primary recovery and centralized DON-scoped event
// routing on the static shard assignment stack (CRE-SHARD-M5-3/4/5).
//
// All three run on the sharded shared-vault failover topology (see
// configs/workflow-gateway-sharded-shared-vault-failover.toml): two workflow
// shard DONs (4 nodes each, no vault capability of their own), the one shared
// vault DON reachable by both shards, and the bootstrap gateway. The
// ShardingFailoverEnabled gate is open on both workflow shards from boot, so
// the secondary shard keeps a standby engine per workflow that denies and
// caches every trigger it does not own.
//
// The scenarios share one mechanism: ownership is re-resolved per event from
// the shard assignment job, so killing a shard's containers (silent death: no
// ExecutionStatusUpdate is ever sent), re-proposing the assignment, or
// restarting a shard's containers all take effect on the next event without
// restarting anything else.

const (
	// Settings fragments applied via t_helpers.ApplyCRESettings (global
	// scope, layered onto the boot CL_CRE_SETTINGS baseline). The 30s window
	// is the test-friendly override of the 5m production default: with the
	// 30s cron schedule the secondary then auto-executes within roughly a
	// minute of the primary going silent. The gate is re-checked both when a
	// secondary caches a trigger and when its failover deadline elapses, so a
	// runtime flip takes effect on the next event.
	shardAutoFailoverOnTOML = "ShardingFailoverAutoExecutionEnabled = 'true'\nShardingFailoverAutoWindow = '30s'"
	// The window is kept in the off fragment too: window-only changes without
	// the gate would still arm failover deadlines on freshly cached events.
	shardAutoFailoverOffTOML = "ShardingFailoverAutoExecutionEnabled = 'false'\nShardingFailoverAutoWindow = '30s'"

	// Log needle the secondary's ShardFailoverManager emits when it executes a
	// cached trigger event itself after the failover window elapsed.
	shardAutoFailoverExecutedLogNeedle = "secondary shard: auto failover executed cached trigger event"

	// shardAutoFailoverAwaitTimeout covers one sync + up to two cron ticks +
	// the failover window + slack for the auto-execution of a fresh workflow.
	shardAutoFailoverAwaitTimeout = 4 * time.Minute
	// shardNoFalseFailoverWindow is how long a healthy primary keeps serving
	// after the automatic failover gate is opened before the primary is
	// killed: window + two cron ticks + slack. No auto-execution in that
	// window proves a live primary's ExecutionStatusUpdates drain the
	// secondary's cache and no false failover fires.
	shardNoFalseFailoverWindow = 2 * time.Minute
	// shardRecoveryStableWindow is the observation window for the
	// no-execution / no-duplicate assertions: at least two cron ticks.
	shardRecoveryStableWindow = 90 * time.Second
)

// shardUserLogRecord is one user-log platform event observed by the collector,
// attributed to the DON whose node emitted it.
type shardUserLogRecord struct {
	at          time.Time
	workflowID  string
	executionID string
	donID       int32
	p2pID       string
}

// shardUserLogCollector consumes the chip test sink's user-log events for the
// whole test and records, per event, the emitting node's DON identity as
// tagged on the platform event itself (WorkflowMetadata.DonID) plus its P2P
// identity. It is the per-DON attribution every assertion below is built on.
type shardUserLogCollector struct {
	mu      sync.Mutex
	records []shardUserLogRecord
}

// startShardUserLogCollector starts a chip sink and records every vaultsecretcron
// user log (sharedVaultUserLogPrefix) with its platform-event DON tagging.
func startShardUserLogCollector(t *testing.T, testEnv *ttypes.TestEnvironment) *shardUserLogCollector {
	t.Helper()
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execCtx, cancelCause := context.WithCancelCause(t.Context())
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	c := &shardUserLogCollector{}
	go func() {
		for userLogs := range userLogsCh {
			if userLogs.M == nil {
				continue
			}
			hasUserLog := false
			for _, line := range userLogs.LogLines {
				if strings.Contains(line.Message, sharedVaultUserLogPrefix) {
					hasUserLog = true
					break
				}
			}
			if !hasUserLog {
				continue
			}
			c.mu.Lock()
			c.records = append(c.records, shardUserLogRecord{
				at:          time.Now(),
				workflowID:  userLogs.M.WorkflowID,
				executionID: userLogs.M.WorkflowExecutionID,
				donID:       userLogs.M.DonID,
				p2pID:       strings.TrimPrefix(userLogs.M.P2PID, "p2p_"),
			})
			c.mu.Unlock()
		}
	}()
	return c
}

// recordsSince returns the records captured at or after cutoff.
func (c *shardUserLogCollector) recordsSince(cutoff time.Time) []shardUserLogRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	records := make([]shardUserLogRecord, 0, len(c.records))
	for _, r := range c.records {
		if !r.at.Before(cutoff) {
			records = append(records, r)
		}
	}
	return records
}

// requireOnlyDONExecuted waits out window and then requires that every user
// log for workflowIDs in the trailing window came from allowedDON, and that
// at least one did - the alive check that separates "the forbidden shard is
// silent" from "nothing is executing at all".
func (c *shardUserLogCollector) requireOnlyDONExecuted(t *testing.T, workflowIDs []string, allowedDON *cre.Don, window time.Duration) {
	t.Helper()

	time.Sleep(window)
	trailing := c.recordsSince(time.Now().Add(-window))

	allowed := 0
	for _, r := range trailing {
		if !slices.Contains(workflowIDs, r.workflowID) {
			continue
		}
		require.Equal(t, int32(allowedDON.ID), r.donID, //nolint:gosec // G115: DON IDs are small
			"workflow %s executed on DON %d (%s) but only DON %s may execute it",
			r.workflowID, r.donID, r.p2pID, allowedDON.Name)
		allowed++
	}
	require.NotEmpty(t, allowed,
		"no user logs for workflows %v from DON %s within %s - the workflow family stopped executing, which voids the only-DON assertion",
		workflowIDs, allowedDON.Name, window)
}

// requireNoUserLogsForWorkflows waits out window and then requires that none
// of workflowIDs executed anywhere in the trailing window.
func (c *shardUserLogCollector) requireNoUserLogsForWorkflows(t *testing.T, workflowIDs []string, window time.Duration) {
	t.Helper()

	time.Sleep(window)
	trailing := c.recordsSince(time.Now().Add(-window))

	for _, r := range trailing {
		require.NotContains(t, workflowIDs, r.workflowID,
			"workflow %s executed on DON %d (%s) during a window in which it must not execute at all", r.workflowID, r.donID, r.p2pID)
	}
}

// requireNoDuplicateExecutionIDs waits out window and then requires that no
// executionID in the trailing window was observed from two different DONs -
// the unbounded-duplicate-executions guard (a dual-primary would emit the same
// deterministic executionID from both shards).
func (c *shardUserLogCollector) requireNoDuplicateExecutionIDs(t *testing.T, window time.Duration) {
	t.Helper()

	time.Sleep(window)
	trailing := c.recordsSince(time.Now().Add(-window))

	byExecution := make(map[string]map[int32]struct{})
	for _, r := range trailing {
		dons, ok := byExecution[r.executionID]
		if !ok {
			dons = make(map[int32]struct{})
			byExecution[r.executionID] = dons
		}
		dons[r.donID] = struct{}{}
	}
	for executionID, dons := range byExecution {
		require.Len(t, dons, 1,
			"execution %s was observed from %d different DONs within %s - duplicate (dual-primary) executions",
			executionID, len(dons), window)
	}
}

// requireDONTaggedEvents requires that the collector observed user logs for
// each expected workflow on its expected DON, and that every record's DON tag
// (WorkflowMetadata.DonID) agrees with the DON derived from the emitting
// node's P2P identity - the platform events are tagged with DON identity and
// a consumer filtering on that tag routes them correctly (CRE-SHARD-M5-5).
func (c *shardUserLogCollector) requireDONTaggedEvents(t *testing.T, expected map[string]*cre.Don, nodeP2PIDToShardIndex map[string]uint32) {
	t.Helper()

	all := c.recordsSince(time.Time{})

	observed := make(map[string]map[uint32]int)
	for _, r := range all {
		if _, ok := expected[r.workflowID]; !ok {
			continue
		}
		donID := uint32(r.donID) //nolint:gosec // G115: DON IDs are small
		if observed[r.workflowID] == nil {
			observed[r.workflowID] = make(map[uint32]int)
		}
		observed[r.workflowID][donID]++

		derivedDON, known := nodeP2PIDToShardIndex[r.p2pID]
		require.True(t, known, "user log for workflow %s from unknown node %s", r.workflowID, r.p2pID)
		require.Equal(t, donID, derivedDON,
			"platform event DON tag (%d) disagrees with the DON derived from node %s (%d) for workflow %s",
			r.donID, r.p2pID, derivedDON, r.workflowID)
	}

	for workflowID, don := range expected {
		require.NotEmpty(t, observed[workflowID],
			"no user logs observed for workflow %s - cannot verify its DON tagging", workflowID)
		require.Contains(t, observed[workflowID], uint32(don.ID), //nolint:gosec // G115: DON IDs are small
			"workflow %s never executed on its expected DON %s (observed DONs: %v)",
			workflowID, don.Name, observed[workflowID])
	}
}

// shardFailoverAssignmentTOML builds the ordered per_owner assignment: the
// first shard is the primary, the second the secondary.
func shardFailoverAssignmentTOML(shards sharedVaultShardPair, primaryFirst bool) string {
	primary, secondary := shards.shardZeroIdx, shards.shardOneIdx
	if !primaryFirst {
		primary, secondary = secondary, primary
	}
	return fmt.Sprintf(`
static_default_assignment = [%d,%d]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%d,%d]
`, primary, secondary, shards.workflowOwner, primary, secondary)
}

// requireCachedEventForWorkflow requires that the nodeset logged a
// cached-trigger line carrying the workflow's ID: the secondary holds a
// standby engine for the workflow and denies (and caches) its triggers.
func requireCachedEventForWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName, workflowID string) {
	t.Helper()

	require.Eventually(t, func() bool {
		for _, line := range t_helpers.ContainerLogLinesForNodeset(t, testEnv, nodesetName, sharedVaultTriggerCachedLogNeedle) {
			if strings.Contains(line, workflowID) {
				return true
			}
		}
		return false
	}, 3*time.Minute, 5*time.Second, "nodeset %s never cached a trigger event for workflow %s", nodesetName, workflowID)
}

// ExecuteShardFailoverAutoTest covers automatic failover on silent primary
// death (CRE-SHARD-M5-3): with the ShardingFailoverAutoExecutionEnabled gate
// open, a secondary shard executes a trigger event itself when no execution
// outcome from the primary arrives within the (configurable) failover window.
// A primary that dies silently - containers stopped, so neither SUCCESS nor
// SYSTEM_ERROR is ever reported - would stall events forever under the
// report-driven failover path; the window is what covers it.
//
// Proof: the secondary's user logs carry the secret (read via the shared
// vault) through the chip sink, the "auto failover executed cached trigger
// event" log needle in the secondary's container logs, and - with the gate
// closed again - the absence of any execution for a fresh workflow, which
// pins the feature flag as the sole enabler.
func ExecuteShardFailoverAutoTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "auto-failover-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deploySharedVaultWorkflow(t, testEnv, "auto-failover-primary", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, 3*time.Minute)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	// Open the gate with a test-friendly window on BOTH shards while both are
	// alive: the settings delivery needs connected nodes. The auto-failover
	// settings are defined inline in core, not in the chainlink-common
	// cresettings catalog, so the delivery opts out of the catalog check.
	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	// A live primary drains the secondary's cache through its
	// ExecutionStatusUpdates, so the open gate alone must not fail anything
	// over: the primary keeps executing and the secondary never auto-executes.
	collector.requireOnlyDONExecuted(t, []string{workflowID}, shards.shardZeroDON, shardNoFalseFailoverWindow)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle)
	testLogger.Info().Msg("Phase 2: gate open, primary healthy - no false failover")

	// Silent primary death: stopped containers report nothing at all.
	t_helpers.StopNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	// A workflow deployed while the primary is dead can only ever execute
	// through the secondary's automatic failover.
	autoWorkflowID := deploySharedVaultWorkflow(t, testEnv, "auto-failover-exec", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{autoWorkflowID}, shards.shardOneDON, secretValue, shardAutoFailoverAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, 2*time.Minute, 5*time.Second)
	collector.requireOnlyDONExecuted(t, []string{autoWorkflowID}, shards.shardOneDON, shardRecoveryStableWindow)
	testLogger.Info().Str("workflowID", autoWorkflowID).Msg("Phase 3: primary died silently, secondary auto-executed within the window")

	// Close the gate: with the primary still dead, a fresh workflow must be
	// synced and cached on the secondary but never executed - the flag is the
	// sole enabler of automatic failover.
	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOffTOML))
	gatedWorkflowID := deploySharedVaultWorkflow(t, testEnv, "auto-failover-gated", secretKey)
	requireCachedEventForWorkflow(t, testEnv, shards.shardOneDON.Name, gatedWorkflowID)
	collector.requireNoUserLogsForWorkflows(t, []string{gatedWorkflowID}, shardNoFalseFailoverWindow)
	testLogger.Info().Str("workflowID", gatedWorkflowID).Msg("Phase 4: gate closed, primary still dead - secondary caches but never executes")
}

// ExecuteShardFailoverPrimaryRecoveryTest covers primary recovery
// (CRE-SHARD-M5-4): the primary dies silently, the survivor is promoted by
// re-proposing the assignment while the old primary is dead, and the former
// primary then returns as a secondary - it must stop executing and cache
// only. Without that failback the recovered shard would run on its stale
// assignment copy as a second primary, and the two primaries would emit
// unbounded duplicate executions.
//
// Proof: after the failback assignment reaches the recovered shard, a
// stability window shows user logs only from the promoted shard (the
// recovered one only caches), no executionID is observed from two DONs (the
// duplicate guard), and a final fail-back moves execution to the recovered
// shard again with the demoted shard returning to cache-only.
func ExecuteShardFailoverPrimaryRecoveryTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-recovery-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deploySharedVaultWorkflow(t, testEnv, "recovery-primary", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, 3*time.Minute)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	// The primary dies silently.
	t_helpers.StopNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	// Promote the survivor by re-proposing the swapped assignment on the
	// alive shard only (a proposal to a dead shard's nodes fails with
	// "node is not connected"). A fresh workflow proves the promoted shard
	// now executes, still reading the SAME secret from the shared vault.
	proposeAndApproveShardAssignmentJob(t, testEnv, shards.shardOneDON, shardFailoverAssignmentTOML(shards, false), testLogger)

	promotedWorkflowID := deploySharedVaultWorkflow(t, testEnv, "recovery-promoted", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{promotedWorkflowID}, shards.shardOneDON, secretValue, 4*time.Minute)
	testLogger.Info().Msg("Phase 2: primary dead, survivor promoted and executing")

	// Primary recovery: the old primary returns - as a secondary, because the
	// current assignment says so. Its engines restart fresh under the swapped
	// assignment (the shard-assignment job is re-proposed to it once its nodes
	// reconnected), so it resolves ownership per event, denies and caches.
	t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	awaitShardDONsConnectedToJD(t, testEnv)
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, false))
	requireCachedEventForWorkflow(t, testEnv, shards.shardZeroDON.Name, workflowID)
	requireCachedEventForWorkflow(t, testEnv, shards.shardZeroDON.Name, promotedWorkflowID)
	testLogger.Info().Msg("Phase 3: former primary returned, re-synced and caching as secondary")

	// Between the restart and the swapped assignment reaching the recovered
	// shard it may briefly execute on its stale assignment copy - a bounded,
	// at-least-once duplicate. What must NOT happen is that duplicate
	// executions persist once the failback assignment landed: in the
	// stability window below, only the promoted shard executes, no executionID
	// shows up on two DONs, and the recovered shard contributes nothing but
	// cached triggers.
	collector.requireOnlyDONExecuted(t, []string{workflowID, promotedWorkflowID}, shards.shardOneDON, shardRecoveryStableWindow)
	collector.requireNoDuplicateExecutionIDs(t, shardRecoveryStableWindow)
	testLogger.Info().Msg("Phase 4: recovered shard is cache-only, promoted shard is the sole executor, no duplicates")

	// Fail back after recovery: re-propose the original order on both shards.
	// The recovered shard executes again and the demoted shard returns to
	// deny-and-cache.
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	failbackWorkflowID := deploySharedVaultWorkflow(t, testEnv, "recovery-failback", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{failbackWorkflowID}, shards.shardZeroDON, secretValue, 4*time.Minute)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 5: fail-back complete, recovered shard executes again, demoted shard caches only")
}

// ExecuteShardFailoverCentralizedEventRoutingTest covers DON-scoped event
// routing on the centralized stack (CRE-SHARD-M5-5): platform events are
// tagged with the emitting DON, so a consumer filtering on that tag routes
// each shard's events correctly - wrong-shard events are distinguishable
// instead of looking like data loss or ghost executions - and duplicate
// execution events, which failover's at-least-once semantics can produce
// transiently at a swap, stay labeled by DON and bounded rather than silent
// and unbounded.
//
// Proof: the chip-sink user-log events carry WorkflowMetadata.DonID, which
// must agree with the DON derived from the emitting node's P2P identity and
// with the shard that owned each phase's execution (shard 0, then shard 1
// after a swap, then shard 1 again through an automatic failover); no
// executionID is observed from two DONs in a stabilized window; and, with
// the observability stack up (`env start --with-dashboards`), the metrics
// ride the centralized pipeline (failover_auto_execution_total and
// platform_engine_workflow_execution_started_count) and the platform logs
// reach Loki (beholder_data_type=zap_log_message).
func ExecuteShardFailoverCentralizedEventRoutingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "centralized-routing-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	collector := startShardUserLogCollector(t, testEnv)
	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	// Phase 1: shard 0 primary - its events must be DON-tagged as shard 0.
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	primaryWorkflowID := deploySharedVaultWorkflow(t, testEnv, "centralized-routing-primary", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{primaryWorkflowID}, shards.shardZeroDON, secretValue, 3*time.Minute)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: shard 0 executing")

	// Phase 2: swap both ways with both shards alive. The same workflow family
	// now executes on shard 1, and the transient in-flight window at the swap
	// is the duplicate-execution source a consumer must be able to see.
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, false))

	swapWorkflowID := deploySharedVaultWorkflow(t, testEnv, "centralized-routing-swap", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{swapWorkflowID}, shards.shardOneDON, secretValue, 4*time.Minute)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardZeroDON.Name, sharedVaultTriggerCachedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 2: swap, shard 1 executing")

	// DON-scoped routing: every platform event is tagged with the emitting
	// DON, the tag agrees with the node's P2P identity, and each workflow
	// executed on the shard the assignment routed it to.
	collector.requireDONTaggedEvents(t, map[string]*cre.Don{
		primaryWorkflowID: shards.shardZeroDON,
		swapWorkflowID:    shards.shardOneDON,
	}, nodeP2PIDToShardIndex)

	// Duplicates stay bounded: once both shards settled on the swapped
	// assignment, no executionID is seen from two DONs.
	collector.requireNoDuplicateExecutionIDs(t, shardRecoveryStableWindow)
	testLogger.Info().Msg("Phase 3: platform events DON-tagged, duplicates bounded")

	// Phase 4: an automatic failover leg (shard 0 dies silently, shard 1
	// auto-executes) so the auto-failover counter is observable on the
	// centralized metric pipeline. The auto-failover settings are core-only
	// (not in the chainlink-common cresettings catalog), so the delivery opts
	// out of the catalog check.
	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	t_helpers.StopNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	autoWorkflowID := deploySharedVaultWorkflow(t, testEnv, "centralized-routing-auto", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{autoWorkflowID}, shards.shardOneDON, secretValue, shardAutoFailoverAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, 2*time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 4: automatic failover executed on shard 1")

	// Metric proof on the centralized pipeline: the auto-failover counter
	// increased and the engine's execution counter is queryable.
	t_helpers.RequirePrometheusQueryEventually(t, "failover_auto_execution_total", 2*time.Minute, 5*time.Second,
		func(sum float64) bool { return sum > 0 },
		"failover_auto_execution_total did not increase after the secondary auto-executed; is the observability stack up (--with-dashboards)?")
	t_helpers.RequirePrometheusQueryEventually(t, "platform_engine_workflow_execution_started_count", 2*time.Minute, 5*time.Second,
		func(sum float64) bool { return sum > 0 },
		"platform_engine_workflow_execution_started_count not observable on the centralized metric pipeline; is the observability stack up (--with-dashboards)?")
	testLogger.Info().Msg("Phase 5: failover_auto_execution_total and execution counters observable via Prometheus")

	// Log proof on the centralized stack: the platform logs reached Loki.
	require.Eventually(t, func() bool {
		count, err := queryLokiForBeholderLogs(t.Context(), framework.LocalLokiBaseURL, 600)
		if err != nil {
			testLogger.Debug().Err(err).Msg("Error querying Loki")
			return false
		}
		return count > 0
	}, 2*time.Minute, 5*time.Second, "Expected beholder zap_log_message logs in Loki while the observability stack is up")
	testLogger.Info().Msg("Phase 6: platform logs flowing to Loki")
}
