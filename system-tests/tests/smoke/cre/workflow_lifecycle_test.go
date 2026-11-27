package cre

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// Lifecycle scenarios that pin the customer-facing lifecycle events (WorkflowActivated,
// WorkflowPaused, WorkflowDeleted): after every step the events each shard sent for a
// workflow must match exactly, so a missing and an extra event both fail.
//
// How events are counted:
//   - a node sends every lifecycle event twice, directly through Beholder and through the
//     durable emitter; only the durable copy is counted;
//   - the chip router of Local CRE returns no per-event results, so the durable emitter
//     keeps resending every event. A resend keeps its CloudEvent ID, so distinct IDs are
//     counted, not messages;
//   - an event carries no node identity, only the shard index. One emit is one ID per worker
//     node, so a shard's count is a multiple of its worker count;
//   - a WorkflowPaused or WorkflowDeleted with a "drain in progress" error is what a node
//     reports today while an execution is in flight: it is logged and not counted. Any other
//     error fails the test.

const (
	workflowActivated = "workflows.v2.WorkflowActivated"
	workflowPaused    = "workflows.v2.WorkflowPaused"
	workflowDeleted   = "workflows.v2.WorkflowDeleted"

	lifecycleEventTimeout = 2 * time.Minute

	lifecycleSettle = 30 * time.Second
)

type lifecycleMessage interface {
	proto.Message
	GetCreInfo() *workfloweventsv2.CreInfo
	GetWorkflow() *workfloweventsv2.Workflow
	GetErrorMessage() string
}

type lifecycleKey struct {
	event string
	shard uint
}

type lifecycleCounts map[lifecycleKey]int

type lifecycleEvent struct {
	lifecycleKey
	workflowID string
	errMsg     string
}

type lifecycleCollector struct {
	mu     sync.Mutex
	events map[string]lifecycleEvent // by CloudEvent ID
}

func startLifecycleCollector(t *testing.T) *lifecycleCollector {
	t.Helper()

	c := &lifecycleCollector{events: make(map[string]lifecycleEvent)}
	server := t_helpers.StartChipTestSink(t, c.publish)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server)
	})
	return c
}

func (c *lifecycleCollector) publish(_ context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
	var msg lifecycleMessage
	switch event.GetType() {
	case workflowActivated:
		msg = &workfloweventsv2.WorkflowActivated{}
	case workflowPaused:
		msg = &workfloweventsv2.WorkflowPaused{}
	case workflowDeleted:
		msg = &workfloweventsv2.WorkflowDeleted{}
	default:
		return &chippb.PublishResponse{}, nil
	}
	if event.GetAttributes()["emitter"].GetCeString() != "DurableEmitter" {
		return &chippb.PublishResponse{}, nil
	}
	if err := proto.Unmarshal(event.GetProtoData().GetValue(), msg); err != nil {
		framework.L.Error().Err(err).Str("ce_type", event.GetType()).Msg("Failed to unmarshal protobuf; skipping")
		return &chippb.PublishResponse{}, nil
	}

	e := lifecycleEvent{
		lifecycleKey: lifecycleKey{event: event.GetType(), shard: uint(msg.GetCreInfo().GetShardIndex())},
		workflowID:   msg.GetWorkflow().GetWorkflowKey().GetWorkflowID(),
		errMsg:       msg.GetErrorMessage(),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, resent := c.events[event.GetId()]; !resent && e.errMsg != "" {
		framework.L.Warn().Str("ce_type", e.event).Str("workflowID", e.workflowID).Uint("shard", e.shard).
			Str("error", e.errMsg).Msg("Lifecycle event with an error message")
	}
	c.events[event.GetId()] = e
	return &chippb.PublishResponse{}, nil
}

// table returns how many events without an error message every shard sent for workflowID,
// and the events whose error is not a drain rejection.
func (c *lifecycleCollector) table(workflowID string) (lifecycleCounts, []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	counts := lifecycleCounts{}
	var unexpected []string
	for _, e := range c.events {
		if e.workflowID != workflowID {
			continue
		}
		switch {
		case e.errMsg == "":
			counts[e.lifecycleKey]++
		case e.event == workflowActivated || !strings.Contains(e.errMsg, "drain in progress"):
			unexpected = append(unexpected, e.event+": "+e.errMsg)
		}
	}
	return counts, unexpected
}

// waitFor waits until don has sent n events of this type for workflowID.
func (c *lifecycleCollector) waitFor(t *testing.T, workflowID, event string, don *cre.Don, n int) {
	t.Helper()

	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		counts, unexpected := c.table(workflowID)
		assert.GreaterOrEqualf(ct, counts[lifecycleKey{event: event, shard: don.ShardIndex}], n,
			"%s events of workflow %s from DON %s (all events: %v, errors: %v)", event, workflowID, don.Name, counts, unexpected)
	}, lifecycleEventTimeout, time.Second)
}

// requireCounts requires the events sent for workflowID to be exactly want: an event type and
// shard missing from want must have no events.
func (c *lifecycleCollector) requireCounts(t *testing.T, workflowID string, want lifecycleCounts) {
	t.Helper()

	got, unexpected := c.table(workflowID)
	require.Emptyf(t, unexpected, "lifecycle events of workflow %s carry an unexpected error", workflowID)
	require.Equalf(t, want, got, "lifecycle events of workflow %s", workflowID)
}

// ExecuteWorkflowLifecycleTest walks cron workflows through register, pause, activate and
// delete on shard 0 and then moves their owner to shard 1, with failover off.
func ExecuteWorkflowLifecycleTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	shard0, shard1 := lifecycleShardDONs(t, testEnv)
	w0, w1 := shard0.WorkersCount(), shard1.WorkersCount()
	owner := lifecycleWorkflowOwner(testEnv)
	events := startLifecycleCollector(t)

	assignOwnerToShards(t, testEnv, owner, []*cre.Don{shard0}, []*cre.Don{shard1})
	wfA := deployLifecycleWorkflow(t, testEnv, "lifecycle-a")
	wfB := deployLifecycleWorkflow(t, testEnv, "lifecycle-b")
	wfC := deployLifecycleWorkflow(t, testEnv, "lifecycle-c")

	want := map[string]lifecycleCounts{wfA: {}, wfB: {}, wfC: {}}
	// expect waits until don has sent n events of this type for the workflow in total, then
	// requires the whole table of the workflow.
	expect := func(workflowID, event string, don *cre.Don, n int) {
		t.Helper()
		want[workflowID][lifecycleKey{event: event, shard: don.ShardIndex}] = n
		events.waitFor(t, workflowID, event, don, n)
		events.requireCounts(t, workflowID, want[workflowID])
	}

	expect(wfA, workflowActivated, shard0, w0)
	expect(wfB, workflowActivated, shard0, w0)
	expect(wfC, workflowActivated, shard0, w0)

	pauseWorkflow(t, testEnv, wfA)
	pauseWorkflow(t, testEnv, wfC)
	expect(wfA, workflowPaused, shard0, w0)
	expect(wfC, workflowPaused, shard0, w0)

	activateWorkflow(t, testEnv, wfA)
	expect(wfA, workflowActivated, shard0, 2*w0)

	// Today a node reports a user delete of a running workflow twice: the diff reports it and
	// the orphan sweep reports it again in the same tick.
	deleteWorkflow(t, testEnv, wfA)
	expect(wfA, workflowDeleted, shard0, 2*w0)

	// Moving the owner is not a user action on a workflow: shard 0 releases wfB and wfC and
	// shard 1 starts wfB, and neither shard may report a lifecycle event for it.
	assignOwnerToShards(t, testEnv, owner, []*cre.Don{shard1}, []*cre.Don{shard1})
	requireLogOnEveryWorker(t, testEnv, shard1, "handled event (WorkflowActivated)", wfB)
	requireLogOnEveryWorker(t, testEnv, shard0, "handled event (WorkflowDeleted)", wfB)
	requireLogOnEveryWorker(t, testEnv, shard0, "handled event (WorkflowDeleted)", wfC)
	time.Sleep(lifecycleSettle)
	events.requireCounts(t, wfA, want[wfA])
	events.requireCounts(t, wfB, want[wfB])
	events.requireCounts(t, wfC, want[wfC])

	deleteWorkflow(t, testEnv, wfB)
	expect(wfB, workflowDeleted, shard1, 2*w1)
}

// ExecuteWorkflowLifecycleFailoverTest registers, swaps and deletes a cron workflow with
// ShardingFailoverEnabled on. Both shards of the assignment hold the workflow, so both report
// its registration and its delete, and swapping primary and secondary reports nothing.
func ExecuteWorkflowLifecycleFailoverTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	primary, secondary := lifecycleShardDONs(t, testEnv)
	shards := []*cre.Don{primary, secondary}
	owner := lifecycleWorkflowOwner(testEnv)
	events := startLifecycleCollector(t)

	assignOwnerToShards(t, testEnv, owner, shards, shards)
	workflowID := deployLifecycleWorkflow(t, testEnv, "lifecycle-failover")

	want := lifecycleCounts{}
	for _, don := range shards {
		want[lifecycleKey{event: workflowActivated, shard: don.ShardIndex}] = don.WorkersCount()
		events.waitFor(t, workflowID, workflowActivated, don, don.WorkersCount())
	}
	events.requireCounts(t, workflowID, want)

	swapped := []*cre.Don{secondary, primary}
	assignOwnerToShards(t, testEnv, owner, swapped, swapped)
	time.Sleep(lifecycleSettle)
	events.requireCounts(t, workflowID, want)

	// Two per node, as in ExecuteWorkflowLifecycleTest.
	deleteWorkflow(t, testEnv, workflowID)
	for _, don := range shards {
		want[lifecycleKey{event: workflowDeleted, shard: don.ShardIndex}] = 2 * don.WorkersCount()
		events.waitFor(t, workflowID, workflowDeleted, don, 2*don.WorkersCount())
	}
	events.requireCounts(t, workflowID, want)
}

// lifecycleShardDONs returns shard zero and the other shard DON.
func lifecycleShardDONs(t *testing.T, testEnv *ttypes.TestEnvironment) (*cre.Don, *cre.Don) {
	t.Helper()

	shardZero := getShardZeroDon(t, testEnv)
	for _, don := range testEnv.Dons.DonsWithFlag(cre.ShardDON) {
		if don.ID != shardZero.ID {
			return shardZero, don
		}
	}
	require.FailNow(t, "Expected to find a second shard DON")
	return nil, nil
}

func lifecycleWorkflowOwner(testEnv *ttypes.TestEnvironment) string {
	return strings.ToLower(testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.MustGetRootKeyAddress().Hex())
}

// assignOwnerToShards pins owner to ownerShards, with defaultShards as the static default, on
// every shard DON and waits until every worker node applied it.
func assignOwnerToShards(t *testing.T, testEnv *ttypes.TestEnvironment, owner string, ownerShards, defaultShards []*cre.Don) {
	t.Helper()

	// The comment keeps every spec distinct, see shardFailoverAssignmentTOML.
	shardAssignmentTOML := fmt.Sprintf(`
static_default_assignment = [%s]
hashed_default_assignment = false

[per_owner_assignment]
  %q = [%s]

# proposed %s
`, shardIndexList(defaultShards), owner, shardIndexList(ownerShards), time.Now().Format(time.RFC3339Nano))

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, shardAssignmentTOML, framework.L)
	}
	for _, don := range shardDONs {
		requireAssignmentLandedOnShard(t, testEnv, don, shardAssignmentUpdateHashNeedle(shardAssignmentTOML))
	}
}

func shardIndexList(dons []*cre.Don) string {
	indexes := make([]string, 0, len(dons))
	for _, don := range dons {
		indexes = append(indexes, strconv.FormatUint(uint64(don.ShardIndex), 10))
	}
	return strings.Join(indexes, ",")
}

// deployLifecycleWorkflow deploys the cron example under a name unique to this run. The
// workflow ID is derived from the name, and nodes keep resending the events of earlier runs.
func deployLifecycleWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, name string) string {
	t.Helper()

	workflowName := fmt.Sprintf("%s-%d", name, time.Now().Unix()%100000)
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L, workflowName, &workflowConfig,
		"../../../../core/scripts/cre/environment/examples/workflows/cron/main.go")
}

// deleteWorkflow deletes a workflow from the WorkflowRegistry.
func deleteWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) {
	t.Helper()

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	_, err := sc.Decode(workflowRegistry(t, testEnv).DeleteWorkflow(sc.NewTXOpts(), workflowIDBytes(t, workflowID)))
	require.NoError(t, err, "failed to delete workflow")
}

// requireLogOnEveryWorker waits until every worker node of don has logged message for workflowID.
func requireLogOnEveryWorker(t *testing.T, testEnv *ttypes.TestEnvironment, don *cre.Don, message, workflowID string) {
	t.Helper()

	require.Eventuallyf(t, func() bool {
		return t_helpers.CountNodesetContainersWithLog(t, testEnv, don.Name, message, workflowID) >= don.WorkersCount()
	}, lifecycleEventTimeout, 5*time.Second, "not every worker node of DON %s logged %q for workflow %s", don.Name, message, workflowID)
}
