package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteConsensusTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")

	channel, listenerCtx, cancelFn, startErr := t_helpers.StartWorkflowEventsSubscriber(t.Context(), t_helpers.GetStandardWorkflowEventsSubscriberConfig(testEnv, workflowID))
	require.NoError(t, startErr, "Failed to start workflow events subscriber")
	t.Cleanup(cancelFn)

	channels := t_helpers.FanOutWorkflowEvents(listenerCtx, channel, 2)
	go func() {
		t_helpers.LogWorkflowEvent(listenerCtx, channels[0])
	}()

	expectedBeholderLog := "Successfully passed all consensus tests"
	err := t_helpers.AssertWorkflowEventMatched(listenerCtx, 2, channels[1], t_helpers.GetUserLogMatcherFn(expectedBeholderLog), 4*time.Minute, testLogger)
	require.NoError(t, err, "Consensus capability test failed")
	testLogger.Info().Msg("Consensus capability test completed")
}
