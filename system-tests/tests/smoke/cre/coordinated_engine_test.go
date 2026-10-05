package cre

import (
	"context"
	"encoding/hex"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	workflow_registry_wrapper_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// coordinatedEngineDON is the workflow DON of the default topology
// (configs/workflow-gateway-capabilities-don.toml).
const coordinatedEngineDON = "workflow"

// Test_CRE_V2_CoordinatedEngine runs a cron workflow with CoordinatedEngineEnabled
// on, so its triggers are registered, delivered and ACKed by the node's trigger
// coordinator instead of by the engine itself.
//
// The flag is read once at engine creation, so it is applied before deploying.
// It is a global override on the shared environment, hence serial.
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_V2_CoordinatedEngine$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_CoordinatedEngine(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`CoordinatedEngineEnabled = 'true'`))

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "coordinatedcron")
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))

	// The execution alone would also pass on the legacy engine; these prove it ran on the coordinated path.
	// Both log lines carry workflowID=<hex>, so the needles attribute them to this workflow.
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, coordinatedEngineDON,
		"Routing workflow to the coordinated engine", time.Minute, 5*time.Second)
	requireCoordinatorLogForWorkflow(t, testEnv, workflowID)
}

// Test_CRE_V2_CoordinatedEngine_MixedModeDeploy starts a workflow on the legacy
// engine (flag off), then enables CoordinatedEngineEnabled and deploys a second
// workflow. Both must keep executing: the first on its legacy engine, the second
// through the trigger coordinator.
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_V2_CoordinatedEngine_MixedModeDeploy$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_CoordinatedEngine_MixedModeDeploy(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}

	// Flag off: the first workflow registers and runs on the legacy engine.
	legacyName := t_helpers.UniqueWorkflowName(testEnv, "legacycron")
	legacyID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, legacyName, &workflowConfig, workflowFileLocation)
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(legacyID))

	// Flag on: only engines created from now on are coordinated.
	t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`CoordinatedEngineEnabled = 'true'`))

	coordinatedName := t_helpers.UniqueWorkflowName(testEnv, "coordinatedcron")
	coordinatedID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, coordinatedName, &workflowConfig, workflowFileLocation)
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(coordinatedID))

	// The legacy workflow must still be executing after the flag flip and the second deploy.
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(legacyID))

	// The coordinator registration log must be attributable to the second workflow only.
	requireCoordinatorLogForWorkflow(t, testEnv, coordinatedID)
	require.Empty(t, coordinatorLogLinesForWorkflow(t, testEnv, legacyID),
		"legacy workflow must not have registered triggers via the coordinator")
}

// Test_CRE_V2_CoordinatedEngine_MixedModeReRegister starts a workflow on the
// legacy engine (flag off), enables CoordinatedEngineEnabled, then pauses and
// re-activates the legacy workflow on the registry. The syncer tears it down and
// re-registers it, and the new engine is created with the flag on — so both
// workflows end up coordinated.
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_V2_CoordinatedEngine_MixedModeReRegister$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_CoordinatedEngine_MixedModeReRegister(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}

	// Flag off: the first workflow registers and runs on the legacy engine.
	legacyName := t_helpers.UniqueWorkflowName(testEnv, "legacycron")
	legacyID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, legacyName, &workflowConfig, workflowFileLocation)
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(legacyID))

	// Flag on before the pause/activate so the re-registered engine is coordinated.
	t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`CoordinatedEngineEnabled = 'true'`))

	coordinatedName := t_helpers.UniqueWorkflowName(testEnv, "coordinatedcron")
	coordinatedID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, coordinatedName, &workflowConfig, workflowFileLocation)
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(coordinatedID))

	// The engine reads the flag once at creation, so the legacy workflow must be
	// recreated to migrate: pause removes its engine, activate builds a new one
	// with the flag on. The wait matters because a pause with an in-flight
	// execution is deferred, and an activate processed while the old engine is
	// still registered takes the happy path (engine exists, ready, active) and
	// returns without recreating it — the legacy engine would keep running until
	// a later reconcile tick happens to rebuild it.
	pauseWorkflow(t, testEnv, legacyID)
	waitForEngineTeardown(t, testEnv, legacyID)
	activateWorkflow(t, testEnv, legacyID)

	// Both workflows must execute after the re-registration.
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(legacyID))
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(coordinatedID))

	// The re-registered legacy workflow must have gone through the coordinator too:
	// its registration log is distinct from the second workflow's earlier one.
	requireCoordinatorLogForWorkflow(t, testEnv, legacyID)
	requireCoordinatorLogForWorkflow(t, testEnv, coordinatedID)
}

// Test_CRE_V2_CoordinatedEngine_RollbackToLegacy starts a workflow with
// CoordinatedEngineEnabled on, then disables the flag and pauses/re-activates the
// workflow. The syncer tears down the coordinated engine and re-registers it, and
// the new engine is created with the flag off — the workflow ends up on the legacy
// engine.
//
// This is the rollback path. The flag is read once at engine creation, so a flip
// alone does not migrate running workflows: ops must also pause/activate each
// workflow (as here) or restart the nodes, which rebuild every engine from the
// registry state with the flag read fresh.
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_V2_CoordinatedEngine_RollbackToLegacy$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_CoordinatedEngine_RollbackToLegacy(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	// Flag on: the workflow registers and runs through the trigger coordinator.
	settings := t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`CoordinatedEngineEnabled = 'true'`))

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "coordinatedcron")
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	requireCoordinatorLogForWorkflow(t, testEnv, workflowID)
	registrationsBefore := countCoordinatorLogLinesForWorkflow(t, testEnv, workflowID)

	// Flag off: engines created from now on are legacy. The baseline has the flag
	// off, so resetting the override disables it — a second ApplyCRESettings call
	// would fail: only one override may be active per test.
	settings.Reset(t)

	// Recreate the engine so it picks up the disabled flag: pause removes it,
	// activate builds a new one. The wait matters because a pause with an
	// in-flight execution is deferred, and an activate processed while the old
	// engine is still registered takes the happy path and returns without
	// recreating it.
	pauseWorkflow(t, testEnv, workflowID)
	waitForEngineTeardown(t, testEnv, workflowID)
	activateWorkflow(t, testEnv, workflowID)

	// The workflow must keep executing on the legacy engine.
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))

	// No new coordinator registration may appear for the re-created engine:
	// the line count must be unchanged since before the pause.
	require.Equal(t, registrationsBefore, countCoordinatorLogLinesForWorkflow(t, testEnv, workflowID),
		"the re-created engine must not register triggers via the coordinator")
}

// coordinatorLogMessage is the syncer log line emitted once a workflow's triggers
// are registered through the coordinator.
const coordinatorLogMessage = "Registered triggers via coordinator"

// coordinatorLogLinesForWorkflow filters the coordinator registration lines down to one
// workflow. The message and the workflow ID are matched in two passes, so fields injected
// between them (e.g. the node logger's version field) do not break the match. workflowID
// is lowercase hex with no 0x — the format CompileAndDeployWorkflow returns.
func coordinatorLogLinesForWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) []string {
	t.Helper()
	lines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, coordinatedEngineDON, coordinatorLogMessage)
	return slices.DeleteFunc(lines, func(line string) bool { return !strings.Contains(line, workflowID) })
}

// requireCoordinatorLogForWorkflow waits for the coordinator registration log of
// workflowID to appear on the workflow DON.
func requireCoordinatorLogForWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) {
	t.Helper()
	require.Eventually(t, func() bool {
		return len(coordinatorLogLinesForWorkflow(t, testEnv, workflowID)) > 0
	}, time.Minute, 5*time.Second, "expected a %q line for workflow %s", coordinatorLogMessage, workflowID)
}

// countCoordinatorLogLinesForWorkflow counts the coordinator registration lines of
// workflowID across the workflow DON: one per node per registration through the
// coordinator. Note the count is per-node — 4 nodes each log one line per
// registration, so a single registration yields 4 lines.
func countCoordinatorLogLinesForWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) int {
	t.Helper()
	return len(coordinatorLogLinesForWorkflow(t, testEnv, workflowID))
}

// workflowRegistry binds the WorkflowRegistry contract for the test chain.
func workflowRegistry(t *testing.T, testEnv *ttypes.TestEnvironment) *workflow_registry_wrapper_v2.WorkflowRegistry {
	t.Helper()

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	registryAddr := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	registry, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(common.HexToAddress(registryAddr), sc.Client)
	require.NoError(t, err, "failed to bind WorkflowRegistry contract")
	return registry
}

// workflowIDBytes decodes a hex workflow ID (the format CompileAndDeployWorkflow returns)
// into the [32]byte the registry contract expects.
func workflowIDBytes(t *testing.T, workflowID string) [32]byte {
	t.Helper()

	idBytes, err := hex.DecodeString(strings.TrimPrefix(workflowID, "0x"))
	require.NoError(t, err, "workflow ID is not valid hex")
	require.Len(t, idBytes, 32, "workflow ID must be 32 bytes")
	return [32]byte(idBytes)
}

// pauseWorkflow pauses a workflow on the WorkflowRegistry.
func pauseWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) {
	t.Helper()

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	_, err := sc.Decode(workflowRegistry(t, testEnv).PauseWorkflow(sc.NewTXOpts(), workflowIDBytes(t, workflowID)))
	require.NoError(t, err, "failed to pause workflow")
}

// activateWorkflow re-activates a paused workflow on the WorkflowRegistry.
func activateWorkflow(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) {
	t.Helper()

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	_, err := sc.Decode(workflowRegistry(t, testEnv).ActivateWorkflow(sc.NewTXOpts(), workflowIDBytes(t, workflowID), testEnv.Dons.MustWorkflowDON().DonFamily()))
	require.NoError(t, err, "failed to activate workflow")
}

// waitForEngineTeardown waits until the pause has been fully processed on every
// workflow node: the syncer logs "handled event (WorkflowPaused)" only after the
// engine is drained, closed and popped. One line per node must appear.
func waitForEngineTeardown(t *testing.T, testEnv *ttypes.TestEnvironment, workflowID string) {
	t.Helper()

	require.Eventually(t, func() bool {
		lines := t_helpers.ContainerLogLinesForNodeset(t, testEnv, coordinatedEngineDON, "handled event (WorkflowPaused)")
		lines = slices.DeleteFunc(lines, func(line string) bool { return !strings.Contains(line, workflowID) })
		return len(lines) >= 4 // one per workflow node
	}, 2*time.Minute, 5*time.Second, "workflow %s was not torn down on all workflow nodes", workflowID)
}
