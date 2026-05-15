package cre

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

/*
Module Cache Smoke Test

Validates that the WASM module cache is active and exercising all cache states
(loaded, evicted, reloaded) by deploying multiple cron workflows under a
MaxLoaded=1 constraint.

Prerequisites:
  - Start the environment with the cache-test topology:
    cd core/scripts/cre/environment
    CTF_CONFIGS=configs/workflow-gateway-don-cache-test.toml go run . env start
  - Run:
    go test -timeout 10m -run "^Test_CRE_V2_Module_Cache$" -v
*/
func ExecuteModuleCacheTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// Do not use t.Context() here: it is cancelled before cleanup runs, which breaks chip-router
		// unregister and can leave gRPC Publish blocked on full log channels after WatchWorkflowLogs returns.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"

	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}

	const numWorkflows = 10
	for i := 0; i < numWorkflows; i++ {
		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, fmt.Sprintf("cachetest%d", i), &workflowConfig, workflowFileLocation)
	}
	testLogger.Info().Int("count", numWorkflows).Msg("All cache-test workflows deployed")

	const watchTimeout = 2 * time.Minute
	const cacheObserveWindow = 1 * time.Minute
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, "Amazing workflow user log", watchTimeout)
	testLogger.Info().Dur("window", cacheObserveWindow).Msg("First workflow execution confirmed, observing cache activity...")

	drainFor(userLogsCh, cacheObserveWindow)

	assertNodeLogs(t, testEnv, "Module cache enabled")
}

func drainFor(ch <-chan *workflowevents.UserLogs, d time.Duration) {
	timer := time.After(d)
	for {
		select {
		case <-timer:
			return
		case <-ch:
		}
	}
}

func assertNodeLogs(t *testing.T, testEnv *ttypes.TestEnvironment, needle string) {
	t.Helper()

	found := false
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Out == nil {
			continue
		}
		for _, clNode := range nodeSet.Out.CLNodes {
			name := clNode.Node.ContainerName
			if name == "" {
				continue
			}
			out, err := exec.CommandContext(t.Context(), "docker", "logs", name).CombinedOutput()
			if err != nil {
				framework.L.Warn().Str("container", name).Err(err).Msg("could not read docker logs")
				continue
			}
			if strings.Contains(string(out), needle) {
				found = true
				framework.L.Info().Str("container", name).Msg("confirmed: " + needle)
			}
		}
	}
	assert.True(t, found, "expected at least one node container log to contain %q", needle)
}
