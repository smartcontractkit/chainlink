package cre

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// shadowDrainedLogMessage is the debug line the M1 central queue drainer writes per workflow per drained batch.
const shadowDrainedLogMessage = "Drained shadow trigger events"

// shadowWindowLayout matches the time range format of FeatureShadowTriggerModeActivePeriod in cresettings.
const shadowWindowLayout = "2006-01-02 15:04:05 -0700 MST"

// Test_CRE_V2_ShadowTriggerMode asserts shadow trigger mode end to end on a real DON. A shadow window delivered
// through CRE settings for one workflow sends that workflow's events through the central queue and drainer,
// while both workflows keep executing. A workflow outside the window is never shadowed. Moving the window into
// the future stops shadowing with no restart.
//
// Shadow mode only affects coordinated engines, so CoordinatedEngineEnabled is on for the whole test. It is a
// global override on the shared environment, hence serial.
//
// Start Local CRE without the ChIP Ingress stack, which takes the test sink's port. Then, from system-tests/tests:
//
//	go test ./smoke/cre -run '^Test_CRE_V2_ShadowTriggerMode$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_ShadowTriggerMode(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	coordinated := t_helpers.Global(`CoordinatedEngineEnabled = 'true'`)
	settings := t_helpers.ApplyCRESettings(t, testEnv, coordinated)

	userLogsCh := make(chan *workfloweventsv2.WorkflowUserLog, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	// executes waits for a user log of workflowID written after this call.
	executes := func(workflowID string) {
		t.Helper()
		discardPending(userLogsCh, baseMessageCh)
		t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
			"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	}
	drained := func(workflowID string) int {
		t.Helper()
		return len(logLinesForWorkflow(t, testEnv, coordinatedEngineDON, shadowDrainedLogMessage, workflowID))
	}

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	// 30s is the fastest schedule the cron trigger accepts.
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}
	shadowedID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger,
		t_helpers.UniqueWorkflowName(testEnv, "shadowcron"), &workflowConfig, workflowFileLocation)
	controlID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger,
		t_helpers.UniqueWorkflowName(testEnv, "controlcron"), &workflowConfig, workflowFileLocation)
	executes(shadowedID)
	executes(controlID)
	require.Zero(t, drained(shadowedID), "shadow mode must be off by default")

	// Each phase is one ApplyCRESettings call, because overlapping overrides fail fast. The coordinated engine
	// flag is read only when an engine is created, so the short reset does not change the running engines.
	now := time.Now().UTC()
	settings.Reset(t)
	settings = t_helpers.ApplyCRESettings(t, testEnv, coordinated,
		t_helpers.Workflow(shadowedID, shadowWindow(now.Add(-time.Hour), now.Add(time.Hour))))

	requireLogForWorkflow(t, testEnv, coordinatedEngineDON, shadowDrainedLogMessage, shadowedID)
	executes(shadowedID)
	executes(controlID)
	require.Zero(t, drained(controlID), "a workflow outside the window must not be shadowed")

	// Kill switch. Once the new window has reached the nodes, the drained count stops growing while the
	// workflow keeps executing.
	settings.Reset(t)
	t_helpers.ApplyCRESettings(t, testEnv, coordinated,
		t_helpers.Workflow(shadowedID, shadowWindow(now.Add(time.Hour), now.Add(2*time.Hour))))

	// Samples are further apart than the 30s cron period, so an equal pair means a tick passed unshadowed.
	settled := waitForStableCount(t, func() int { return drained(shadowedID) }, 40*time.Second, 4*time.Minute)
	executes(shadowedID)
	executes(shadowedID)
	require.Equal(t, settled, drained(shadowedID), "shadowing must stop once the window is in the future")
	require.Zero(t, drained(controlID), "a workflow outside the window must not be shadowed")
}

// shadowWindow returns the PerWorkflow TOML that sets the shadow trigger mode window to [from, to].
func shadowWindow(from, to time.Time) string {
	return fmt.Sprintf("[PerWorkflow]\nFeatureShadowTriggerModeActivePeriod = '[%s,%s]'",
		from.UTC().Format(shadowWindowLayout), to.UTC().Format(shadowWindowLayout))
}

// waitForStableCount samples count every interval and returns it once two samples in a row are equal. It fails
// the test if the count is still changing after timeout.
func waitForStableCount(t *testing.T, count func() int, interval, timeout time.Duration) int {
	t.Helper()
	deadline := time.Now().Add(timeout)
	last := count()
	for time.Now().Before(deadline) {
		time.Sleep(interval)
		current := count()
		if current == last {
			return current
		}
		last = current
	}
	require.FailNow(t, "count never settled", "still changing after %s, last value %d", timeout, last)
	return 0
}

// discardPending empties both chip sink channels, so the next watch only sees messages that arrive after it.
func discardPending(userLogsCh <-chan *workfloweventsv2.WorkflowUserLog, baseMessageCh <-chan *commonevents.BaseMessage) {
	for {
		select {
		case <-userLogsCh:
		case <-baseMessageCh:
		default:
			return
		}
	}
}
