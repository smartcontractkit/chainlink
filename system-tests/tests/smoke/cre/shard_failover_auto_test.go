package cre

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// Shard failover scenarios on the sharded shared-vault topology
// (configs/workflow-gateway-sharded-shared-vault-failover.toml), covering
// CRE-SHARD-M5-3 (automatic failover on silent primary death),
// CRE-SHARD-M5-4 (primary recovery) and CRE-SHARD-M5-5 (DON-scoped event
// routing on the centralized stack). The ShardingFailoverEnabled gate is open
// on both workflow shards from boot, so the secondary keeps a standby engine
// per workflow that denies and caches every trigger it does not own.
//
// All waits are event-driven and sized for the 30s cron cadence (the cron
// trigger's by-design floor) to fit the 7-minute CI test budget: observation
// phases wait for fresh platform events instead of sleeping fixed windows,
// and each scenario reuses one workflow across its phases.

const (
	shardFailoverSchedule = "*/30 * * * * *"

	shardAutoFailoverOnTOML = "ShardingFailoverAutoExecutionEnabled = 'true'\nShardingFailoverAutoWindow = '30s'"

	shardAutoFailoverExecutedLogNeedle = "secondary shard: auto failover executed cached trigger event"

	shardInflightDrainGrace = 10 * time.Second

	shardStallObserveWindow = 45 * time.Second

	shardTickRecords = 1

	shardEventAwaitTimeout = 2 * time.Minute

	shardFirstExecAwaitTimeout = 2 * time.Minute

	shardAutoExecAwaitTimeout = 2 * time.Minute
)

// shardUserLogRecord is one user-log platform event, attributed to the DON
// whose node emitted it. emittedAt is the node's own clock at emission
// (LogLine node_timestamp); boundary logic uses it rather than the receive
// time, because a node with a warming chip connection (e.g. right after a
// container restart) can deliver its user logs tens of seconds late.
type shardUserLogRecord struct {
	emittedAt   time.Time
	workflowID  string
	executionID string
	donID       int32
	p2pID       string
}

// shardUserLogCollector records every vaultsecretcron user log with its
// platform-event DON tagging; it is the per-DON attribution every assertion
// below is built on.
type shardUserLogCollector struct {
	mu      sync.Mutex
	records []shardUserLogRecord
}

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
			emittedAt := time.Time{}
			hasUserLog := false
			for _, line := range userLogs.LogLines {
				if !strings.Contains(line.Message, sharedVaultUserLogPrefix) {
					continue
				}
				hasUserLog = true
				if ts, err := time.Parse(time.RFC3339Nano, line.NodeTimestamp); err == nil {
					emittedAt = ts
					break
				}
			}
			if !hasUserLog {
				continue
			}
			if emittedAt.IsZero() {
				emittedAt = time.Now()
			}
			c.mu.Lock()
			c.records = append(c.records, shardUserLogRecord{
				emittedAt:   emittedAt,
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

func (c *shardUserLogCollector) recordsSince(cutoff time.Time) []shardUserLogRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	records := make([]shardUserLogRecord, 0, len(c.records))
	for _, r := range c.records {
		if !r.emittedAt.Before(cutoff) {
			records = append(records, r)
		}
	}
	return records
}

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

func (c *shardUserLogCollector) requireNoRecordsForWorkflows(t *testing.T, workflowIDs []string, since time.Time) {
	t.Helper()

	for _, r := range c.recordsSince(since) {
		require.NotContainsf(t, workflowIDs, r.workflowID,
			"workflow %s executed on DON %d (%s) in a window where it must not execute at all", r.workflowID, r.donID, r.p2pID)
	}
}

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

// requireDONTaggedEvents requires user logs for workflowID on every DON in
// wantDONs, with every record's DON tag (WorkflowMetadata.DonID) agreeing with
// the DON derived from the emitting node's P2P identity.
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
// first shard is the primary, the second the secondary. Every spec carries a
// unique comment: every shard-assignment proposal uses the same fixed job
// UUID, and a proposal byte-identical to an earlier proposal of that job is a
// silent no-op, so repeated orders (a failback, or re-runs on a shared local
// environment) must be distinct spec strings.
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

# proposed %s
`, primary, secondary, shards.workflowOwner, primary, secondary, time.Now().Format(time.RFC3339Nano))
}

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

// shardAssignmentUpdateHashNeedle returns the log needle proving a node applied
// the shard-assignment spec built from toml: the job template stamps the spec
// with sha256 of the config_type marker plus the TOML, and every node logs the
// hash on apply. Waiting for it on every worker container of a shard is the
// node-complete landing proof.
func shardAssignmentUpdateHashNeedle(toml string) string {
	settings := "config_type = \"shard_assignment\"\n" + toml
	sum := sha256.Sum256([]byte(settings))
	return hex.EncodeToString(sum[:])
}

func requireAssignmentLandedOnShard(t *testing.T, testEnv *ttypes.TestEnvironment, don *cre.Don, hashNeedle string) {
	t.Helper()

	workers, err := don.Workers()
	require.NoErrorf(t, err, "failed to get worker nodes of DON %q", don.Name)
	require.NotEmpty(t, workers)

	require.Eventuallyf(t, func() bool {
		return t_helpers.CountNodesetContainersWithLog(t, testEnv, don.Name, hashNeedle) >= len(workers)
	}, 2*time.Minute, 3*time.Second,
		"the shard-assignment update (hash %s) did not land on all %d worker containers of DON %s",
		hashNeedle, len(workers), don.Name)
}

// ExecuteShardFailoverAutoTest covers automatic failover on silent primary
// death (CRE-SHARD-M5-3): with the ShardingFailoverAutoExecutionEnabled gate
// open, a secondary shard executes a trigger event itself when no execution
// outcome from the primary arrives within the (configurable) failover window.
func ExecuteShardFailoverAutoTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv, "configs/workflow-gateway-sharded-shared-vault-failover.toml")
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "auto-failover-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deployShardFailoverWorkflow(t, testEnv, "auto-failover", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	gateArmed := time.Now()
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardZeroDON, shardTickRecords, gateArmed, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardOneDON}, gateArmed)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle)
	testLogger.Info().Msg("Phase 2: gate open, primary healthy - no false failover")

	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	postDeath := death.Add(shardInflightDrainGrace)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardAutoExecAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, time.Minute, 5*time.Second)
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, postDeath, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, postDeath)
	collector.requireNoDuplicateExecutionIDsSince(t, postDeath)
	testLogger.Info().Msg("Phase 3: primary died silently, secondary auto-executed within the window")
}

// ExecuteShardFailoverPrimaryRecoveryTest covers primary recovery
// (CRE-SHARD-M5-4) and the gate-closed stall (the before of CRE-SHARD-M5-3):
// the primary dies silently and the workflow stalls with the automatic
// failover gate at its default closed state; the survivor is promoted by
// re-proposing the assignment while the old primary is dead; the former
// primary returns as a secondary - it must stop executing and cache only -
// and a final fail-back moves execution back to it.
func ExecuteShardFailoverPrimaryRecoveryTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv, "configs/workflow-gateway-sharded-shared-vault-failover.toml")
	fixture := setupSharedVaultShardFixture(t, testEnv)
	shards.workflowOwner = strings.ToLower(fixture.owner)

	secretValue := "failover-recovery-shared-vault-secret"
	secretKey := fixture.createSharedVaultSecret(t, secretValue)

	proposeSharedVaultAssignment(t, testEnv, shardFailoverAssignmentTOML(shards, true))

	workflowID := deployShardFailoverWorkflow(t, testEnv, "shard-recovery", secretKey)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 1: primary executes, secondary denies and caches")

	collector := startShardUserLogCollector(t, testEnv)

	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	time.Sleep(shardStallObserveWindow)
	collector.requireNoRecordsForWorkflows(t, []string{workflowID}, death.Add(shardInflightDrainGrace))
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle)
	testLogger.Info().Msg("Phase 2: gate closed, primary dead - workflow stalls, nothing auto-executes")

	proposeAndApproveShardAssignmentJob(t, testEnv, shards.shardOneDON, shardFailoverAssignmentTOML(shards, false), testLogger)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 3: primary dead, survivor promoted and executing")

	// The recovered shard restarts under its stale assignment copy and briefly
	// executes as a second primary; the swap lands on each node independently,
	// so wait for it on every worker container before measuring the recovered
	// shard's behavior.
	swappedAssignmentTOML := shardFailoverAssignmentTOML(shards, false)
	t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	awaitShardDONsConnectedToJD(t, testEnv)
	proposeSharedVaultAssignment(t, testEnv, swappedAssignmentTOML)
	requireAssignmentLandedOnShard(t, testEnv, shards.shardZeroDON, shardAssignmentUpdateHashNeedle(swappedAssignmentTOML))
	requireCachedEventForWorkflow(t, testEnv, shards.shardZeroDON.Name, workflowID)
	failbackLanded := time.Now()
	testLogger.Info().Msg("Phase 4: former primary returned, re-synced and caching as secondary")

	stableSince := failbackLanded.Add(shardInflightDrainGrace)
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, stableSince, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, stableSince)
	collector.requireNoDuplicateExecutionIDsSince(t, stableSince)
	testLogger.Info().Msg("Phase 5: recovered shard is cache-only, promoted shard is the sole executor, no duplicates")

	failbackTOML := shardFailoverAssignmentTOML(shards, true)
	proposeSharedVaultAssignment(t, testEnv, failbackTOML)
	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardZeroDON, secretValue, shardFirstExecAwaitTimeout)
	testLogger.Info().Msg("Phase 6: fail-back complete, recovered shard executes again")
}

// ExecuteShardFailoverCentralizedEventRoutingTest covers DON-scoped event
// routing on the centralized stack (CRE-SHARD-M5-5): platform events are
// tagged with the emitting DON so a consumer filtering on that tag routes each
// shard's events correctly, and duplicate execution events stay bounded.
// Requires the observability stack (`env start --with-dashboards` locally; in
// CI the workflow's observability step starts it).
func ExecuteShardFailoverCentralizedEventRoutingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shards := mustSharedVaultShardPair(t, testEnv, "configs/workflow-gateway-sharded-shared-vault-failover.toml")
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

	t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.AllowUncatalogedSettings(),
		t_helpers.Global(shardAutoFailoverOnTOML))

	death := time.Now()
	t_helpers.KillNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	t.Cleanup(func() {
		t_helpers.StartNodesetContainers(t, testEnv, shards.shardZeroDON.Name)
	})

	awaitSharedVaultWorkflowExecution(t, testEnv, []string{workflowID}, shards.shardOneDON, secretValue, shardAutoExecAwaitTimeout)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, shards.shardOneDON.Name, shardAutoFailoverExecutedLogNeedle, time.Minute, 5*time.Second)
	testLogger.Info().Msg("Phase 2: primary died silently, secondary auto-executing")

	collector.requireDONTaggedEvents(t, workflowID, []*cre.Don{shards.shardZeroDON, shards.shardOneDON}, nodeP2PIDToShardIndex)

	postDeath := death.Add(shardInflightDrainGrace)
	collector.awaitRecordsFromDON(t, []string{workflowID}, shards.shardOneDON, shardTickRecords, postDeath, shardEventAwaitTimeout)
	collector.requireNoRecordsFromDONs(t, []string{workflowID}, []*cre.Don{shards.shardZeroDON}, postDeath)
	collector.requireNoDuplicateExecutionIDsSince(t, postDeath)
	testLogger.Info().Msg("Phase 3: platform events DON-tagged, duplicates bounded")

	t_helpers.RequirePrometheusQueryEventually(t, "failover_auto_execution_total", 2*time.Minute, 5*time.Second,
		func(sum float64) bool { return sum > 0 },
		"failover_auto_execution_total did not increase after the secondary auto-executed; is the observability stack up?")
	testLogger.Info().Msg("Phase 4: failover_auto_execution_total observable via Prometheus")

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
