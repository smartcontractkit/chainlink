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
	vaultsecretcron_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecretcron/config"
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
//
// Budget: the CI smoke lane gives each test 7 minutes inside a 10-minute job,
// so every wait here is event-driven and sized for the 10s cron schedule and
// the 15s failover window below - observation phases wait for fresh platform
// events to prove cron ticks passed (each tick is a full failover-window
// cycle) instead of sleeping fixed windows, and each scenario reuses ONE
// workflow across all its phases (one WASM compile, clean per-phase
// attribution from the chip-sink records).

const (
	// shardFailoverSchedule is the cron schedule for the scenarios' workflows.
	// The 30s cadence is the cron trigger's by-design floor: a faster schedule
	// is rejected at trigger registration and the workflow engine fails to
	// initialize, so every tick-bound wait costs at least 30s.
	shardFailoverSchedule = "*/30 * * * * *"

	// shardAutoFailoverOnTOML is the settings fragment applied via
	// t_helpers.ApplyCRESettings (global scope, layered onto the boot
	// CL_CRE_SETTINGS baseline). The 30s window is the test-friendly override
	// of the 5m production default: with the 30s cron schedule the secondary
	// then auto-executes within roughly a minute of the primary going silent.
	// The gate is re-checked both when a secondary caches a trigger and when
	// its failover deadline elapses, so the flip takes effect on the next
	// event. The before-phase needs no fragment: the gate's default state is
	// closed, and the boot baseline opens only ShardingFailoverEnabled.
	shardAutoFailoverOnTOML = "ShardingFailoverAutoExecutionEnabled = 'true'\nShardingFailoverAutoWindow = '30s'"

	// Log needle the secondary's ShardFailoverManager emits when it executes a
	// cached trigger event itself after the failover window elapsed.
	shardAutoFailoverExecutedLogNeedle = "secondary shard: auto failover executed cached trigger event"

	// shardStallObserveWindow is the fixed observe beat for the gate-closed
	// stall: one cron tick past the death, so at least one trigger fired and
	// was cached with nothing executing it. The stall itself is about the
	// gate, not the window - with the gate closed a cached event is never
	// armed with a deadline at all.
	shardStallObserveWindow = 35 * time.Second

	// shardTickRecords is how many fresh user-log records an observation
	// phase waits for: each record is one cron tick, and each tick is a full
	// failover-window cycle, so one fresh record proves a whole window cycle
	// elapsed under the condition being observed.
	shardTickRecords = 1

	// shardEventAwaitTimeout bounds the event-driven record waits (ticks
	// arrive every 30s; generous for a slow first sync after a restart).
	shardEventAwaitTimeout = 2 * time.Minute

	// shardFirstExecAwaitTimeout bounds waiting for a workflow's first
	// execution after deployment (sync + first cron tick + execution).
	shardFirstExecAwaitTimeout = 2 * time.Minute

	// shardAutoExecAwaitTimeout bounds waiting for the secondary's first
	// auto-executed user log after the primary's silent death (first post-death
	// cron tick + failover window + execution).
	shardAutoExecAwaitTimeout = 2 * time.Minute
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

// awaitRecordsFromDON waits until minCount fresh records for workflowIDs
// arrived from don at or after since - each record is one cron tick, and each
// tick is a full failover-window cycle. Fails after timeout; follow-up instant
// assertions use the caller's own base times (a phase's flip, death or
// failback-landed moment), so they span the whole observed phase.
func (c *shardUserLogCollector) awaitRecordsFromDON(t *testing.T, workflowIDs []string, don *cre.Don, minCount int, since time.Time, timeout time.Duration) {
	t.Helper()

	require.Eventuallyf(t, func() bool {
		count := 0
		for _, r := range c.recordsSince(since) {
			if workflowIDMatches(workflowIDs, r) && r.donID == int32(don.ID) { //nolint:gosec // G115: DON IDs are small
				count++
			}
		}
		return count >= minCount
	}, timeout, time.Second, "no %d fresh user-log records for workflows %v from DON %s within %s (the workflow family stopped ticking)",
		minCount, workflowIDs, don.Name, timeout)
}

// requireNoRecordsForWorkflows requires that none of workflowIDs produced any
// user-log record at or after since.
func (c *shardUserLogCollector) requireNoRecordsForWorkflows(t *testing.T, workflowIDs []string, since time.Time) {
	t.Helper()

	for _, r := range c.recordsSince(since) {
		require.NotContainsf(t, workflowIDs, r.workflowID,
			"workflow %s executed on DON %d (%s) in a window where it must not execute at all", r.workflowID, r.donID, r.p2pID)
	}
}

// requireNoRecordsFromDONs requires that no user-log record for workflowIDs
// came from any of forbidden at or after since.
func (c *shardUserLogCollector) requireNoRecordsFromDONs(t *testing.T, workflowIDs []string, forbidden []*cre.Don, since time.Time) {
	t.Helper()

	for _, r := range c.recordsSince(since) {
		if !workflowIDMatches(workflowIDs, r) {
			continue
		}
		for _, don := range forbidden {
			require.NotEqualf(t, int32(don.ID), r.donID, //nolint:gosec // G115: DON IDs are small
				"workflow %s executed on DON %s (%s) in a window where only its assigned shard may execute",
				r.workflowID, don.Name, r.p2pID)
		}
	}
}

// requireNoDuplicateExecutionIDsSince requires that no executionID produced
// user-log records from two different DONs at or after since - the
// unbounded-duplicate-executions guard (a dual-primary would emit the same
// deterministic executionID from both shards).
func (c *shardUserLogCollector) requireNoDuplicateExecutionIDsSince(t *testing.T, since time.Time) {
	t.Helper()

	byExecution := make(map[string]map[int32]struct{})
	for _, r := range c.recordsSince(since) {
		dons, ok := byExecution[r.executionID]
		if !ok {
			dons = make(map[int32]struct{})
			byExecution[r.executionID] = dons
		}
		dons[r.donID] = struct{}{}
	}
	for executionID, dons := range byExecution {
		require.Lenf(t, dons, 1,
			"execution %s was observed from %d different DONs - duplicate (dual-primary) executions", executionID, len(dons))
	}
}

// requireDONTaggedEvents requires that the collector observed user logs for
// workflowID on every DON in wantDONs, and that every record's DON tag
// (WorkflowMetadata.DonID) agrees with the DON derived from the emitting
// node's P2P identity - platform events are tagged with DON identity and a
// consumer filtering on that tag routes them correctly (CRE-SHARD-M5-5).
func (c *shardUserLogCollector) requireDONTaggedEvents(t *testing.T, workflowID string, wantDONs []*cre.Don, nodeP2PIDToShardIndex map[string]uint32) {
	t.Helper()

	observed := make(map[uint32]int)
	for _, r := range c.recordsSince(time.Time{}) {
		if r.workflowID != workflowID {
			continue
		}
		donID := uint32(r.donID) //nolint:gosec // G115: DON IDs are small
		observed[donID]++

		derivedDON, known := nodeP2PIDToShardIndex[r.p2pID]
		require.Truef(t, known, "user log for workflow %s from unknown node %s", workflowID, r.p2pID)
		require.Equalf(t, donID, derivedDON,
			"platform event DON tag (%d) disagrees with the DON derived from node %s (%d) for workflow %s",
			donID, r.p2pID, derivedDON, workflowID)
	}

	for _, don := range wantDONs {
		require.NotEmptyf(t, observed[uint32(don.ID)], //nolint:gosec // G115: DON IDs are small
			"workflow %s never executed on DON %s (observed DONs: %v) - its events are not attributable per DON",
			workflowID, don.Name, observed)
	}
}

func workflowIDMatches(workflowIDs []string, r shardUserLogRecord) bool {
	return slices.Contains(workflowIDs, r.workflowID)
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

// deployShardFailoverWorkflow deploys one fast-scheduled (10s cron)
// vaultsecretcron workflow that reads the given secret from the shared vault
// on every run. The scenario reuses this ONE workflow across all its phases -
// every observation is attributed per phase from the chip-sink records, and
// fresh cron ticks give every phase fresh trigger and execution IDs, so no
// phase needs its own compiled workflow.
func deployShardFailoverWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowBaseName, secretKey string) string {
	t.Helper()

	workflowConfig := vaultsecretcron_config.Config{
		Schedule:        shardFailoverSchedule,
		SecretNamespace: sharedVaultNamespace,
		SecretKey:       secretKey,
	}
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L,
		t_helpers.UniqueWorkflowName(testEnv, workflowBaseName), &workflowConfig, sharedVaultWorkflowFileLocation)
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
	}, 2*time.Minute, 5*time.Second, "nodeset %s never cached a trigger event for workflow %s", nodesetName, workflowID)
}

// ExecuteShardFailoverAutoTest covers automatic failover on silent primary
// death (CRE-SHARD-M5-3): with the ShardingFailoverAutoExecutionEnabled gate
// open, a secondary shard executes a trigger event itself when no execution
// outcome from the primary arrives within the (configurable) failover window.
// A primary that dies silently - containers killed, so neither SUCCESS nor
// SYSTEM_ERROR is ever reported - would stall events forever under the
// report-driven failover path; the window is what covers it. The gate-closed
// stall (the before) is covered by ExecuteShardFailoverPrimaryRecoveryTest.
//
// Phases: the gate is opened with both shards alive - a settings delivery
// needs connected nodes, and a healthy-primary phase then proves no false
// failover: the primary's ExecutionStatusUpdates drain the secondary's cache,
// so two fresh ticks (two full failover-window cycles) arrive from the
// primary with zero executions and zero auto-failover needles from the
// secondary. The primary then dies silently, and the secondary executes the
// workflow itself within the window, reading the same secret from the shared
// vault (the replay carries the workflow's tenant identity, so the secret
// path works exactly as on the primary).
// The workflow deploys once, while both shards are alive (the deploy's
// artifact copy targets every workflow shard's containers), and serves on the
// primary until the death; every failover assertion measures post-death
// behavior.
// Proof: the secondary's user logs through the chip sink, and the
// "auto failover executed cached trigger event" log needle in the secondary's
// container logs.
func ExecuteShardFailoverAutoTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "auto-failover-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deployShardFailoverWorkflow(t, testEnv, "auto-failover", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	// Open the gate with a test-friendly window, with both shards alive (the
	// delivery needs connected nodes). The auto-failover settings are defined
	// inline in core, not in the chainlink-common cresettings catalog, so the
	// delivery opts out of the catalog check.
	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	// No false failover while the primary is healthy: its
	// ExecutionStatusUpdates drain the secondary's cache, so a fresh tick (a
	// full failover-window cycle) arrives from the primary with zero
	// executions and zero auto-failover needles from the secondary.
	gateArmed := time.Now()
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardZeroDON, shardTickRecords, gateArmed, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardOneDON}, gateArmed)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle)
	testLogger.Info().Msg("Phase 2: gate open, primary healthy - no false failover")

	// The silent death: stopped-by-SIGKILL containers report nothing at all.
	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	// The secondary executes the cached triggers itself within the failover
	// window, reading the same secret from the shared vault.
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardAutoExecAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, time.Minute, 5*time.Second)
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, death, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, death)
	collector.requireNoDuplicateExecutionIDsSince(t, death)
	testLogger.Info().Msg("Phase 3: primary died silently, secondary auto-executed within the window")
}

// ExecuteShardFailoverPrimaryRecoveryTest covers primary recovery
// (CRE-SHARD-M5-4) and the gate-closed stall (the before of CRE-SHARD-M5-3):
// the primary dies silently and, with the automatic failover gate at its
// default CLOSED state (the boot CL_CRE_SETTINGS baseline opens only
// ShardingFailoverEnabled), the workflow stalls - the secondary keeps denying
// and caching, nothing ever executes. The survivor is then promoted by
// re-proposing the assignment while the old primary is dead, and the former
// primary later returns as a secondary - it must stop executing and cache
// only. Without that failback the recovered shard would run on its stale
// assignment copy as a second primary, and the two primaries would emit
// unbounded duplicate executions.
//
// Proof: nothing executes between the death and the promotion (the stall);
// after the failback assignment reaches the recovered shard, fresh ticks
// arrive only from the promoted shard (the recovered one contributes nothing
// but cached-trigger needles), no executionID is observed from two DONs, and
// a final fail-back moves execution to the recovered shard again.
// The workflow deploys once, while both shards are alive (the deploy's
// artifact copy targets every workflow shard's containers).
func ExecuteShardFailoverPrimaryRecoveryTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-recovery-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deployShardFailoverWorkflow(t, testEnv, "shard-recovery", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	// The primary dies silently, and with the automatic failover gate at its
	// default CLOSED state the workflow stalls - nothing executes it anywhere.
	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	time.Sleep(shardStallObserveWindow)
	collector.requireNoRecordsForWorkflows(t, []string{workflowID}, death)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle)
	testLogger.Info().Msg("Phase 2: gate closed, primary dead - workflow stalls, nothing auto-executes")

	// The survivor is promoted by re-proposing the swapped assignment on the
	// alive shard only (a proposal to a dead shard's nodes fails with "node
	// is not connected"). The promoted shard then executes the workflow, still
	// reading the SAME secret from the shared vault.
	proposeAndApproveShardAssignmentJob(t, testEnv, shards.shardOneDON, shardFailoverAssignmentTOML(shards, false), testLogger)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 3: primary dead, survivor promoted and executing")

	// Primary recovery: the old primary returns - as a secondary, because the
	// current assignment says so. Its engines restart fresh under the swapped
	// assignment (the shard-assignment job is re-proposed to it once its nodes
	// reconnected), so it resolves ownership per event, denies and caches.
	// Between the restart and the swapped assignment landing on it, the
	// recovered shard may briefly execute on its stale assignment copy - a
	// bounded, at-least-once duplicate; the stability assertions below
	// measure after the failback landed.
	t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	awaitShardDONsConnectedToJD(t, testEnv)
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, false))
	requireCachedEventForWorkflow(t, testEnv, shards.shardZeroDON.Name, workflowID)
	failbackLanded := time.Now()
	testLogger.Info().Msg("Phase 4: former primary returned, re-synced and caching as secondary")

	// Recovery assertions: from the failback landing on, fresh ticks arrive
	// only from the promoted shard, and no executionID is seen on two DONs.
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, failbackLanded, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, failbackLanded)
	collector.requireNoDuplicateExecutionIDsSince(t, failbackLanded)
	testLogger.Info().Msg("Phase 5: recovered shard is cache-only, promoted shard is the sole executor, no duplicates")

	// Fail back after recovery: re-propose the original order on both shards.
	// The recovered shard executes again.
	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 6: fail-back complete, recovered shard executes again")
}

// ExecuteShardFailoverCentralizedEventRoutingTest covers DON-scoped event
// routing on the centralized stack (CRE-SHARD-M5-5): platform events are
// tagged with the emitting DON, so a consumer filtering on that tag routes
// each shard's events correctly - wrong-shard events are distinguishable
// instead of looking like data loss or ghost executions - and duplicate
// execution events stay bounded rather than silent and unbounded.
//
// Phases: the workflow serves on the primary (shard 0), the gate is opened
// with both shards alive, the primary dies silently and the secondary
// auto-executes - the same workflow family executing on BOTH shards over the
// test's lifetime, with the failover boundary in between. Proof: the
// chip-sink user-log events carry WorkflowMetadata.DonID, which must agree
// with the DON derived from the emitting node's P2P identity and cover both
// shards; no executionID is observed from two DONs after the death; and, with
// the observability stack up (`env start --with-dashboards`), the scenario's
// own metric rides the centralized pipeline (failover_auto_execution_total)
// and the platform logs reach Loki (beholder_data_type=zap_log_message).
// The workflow deploys once, while both shards are alive (the deploy's
// artifact copy targets every workflow shard's containers).
func ExecuteShardFailoverCentralizedEventRoutingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv)
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "centralized-routing-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	collector := startShardUserLogCollector(t, testEnv)
	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deployShardFailoverWorkflow(t, testEnv, "centralized-routing", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, sharedVaultTriggerCachedLogNeedle, time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 1: shard 0 executing")

	// Phase 2: open the gate with a test-friendly window, with both shards
	// alive (the delivery needs connected nodes; the auto-failover settings
	// are core-only, so the delivery opts out of the catalog check).
	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	// Phase 3: the primary dies silently and the secondary auto-executes - the
	// workflow family now executed on BOTH shards, across the failover
	// boundary.
	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardAutoExecAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 2: primary died silently, secondary auto-executing")

	// DON-scoped routing: every platform event is tagged with the emitting
	// DON, the tag agrees with the node's P2P identity, and both shards'
	// executions are attributable - wrong-shard events distinguishable.
	collector.requireDONTaggedEvents(t, workflowID, []*cre.Don{shards.shardZeroDON, shards.shardOneDON}, nodeP2PIDToShardIndex)

	// Duplicates stay bounded: fresh ticks only from the secondary, and no
	// executionID from two DONs after the death.
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, death, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, death)
	collector.requireNoDuplicateExecutionIDsSince(t, death)
	testLogger.Info().Msg("Phase 3: platform events DON-tagged, duplicates bounded")

	// Metric proof on the centralized pipeline: the scenario's own
	// auto-failover counter is observable and non-zero after the secondary
	// auto-executed.
	t_helpers.RequirePrometheusQueryEventually(t, "failover_auto_execution_total", 2*time.Minute, 5*time.Second,
		func(sum float64) bool { return sum > 0 },
		"failover_auto_execution_total did not increase after the secondary auto-executed; is the observability stack up (--with-dashboards)?")
	testLogger.Info().Msg("Phase 4: failover_auto_execution_total observable via Prometheus")

	// Log proof on the centralized stack: the platform logs reached Loki.
	require.Eventually(t, func() bool {
		count, err := queryLokiForBeholderLogs(t.Context(), framework.LocalLokiBaseURL, 600)
		if err != nil {
			testLogger.Debug().Err(err).Msg("Error querying Loki")
			return false
		}
		return count > 0
	}, time.Minute, 5*time.Second, "Expected beholder zap_log_message logs in Loki while the observability stack is up")
	testLogger.Info().Msg("Phase 5: platform logs flowing to Loki")
}
