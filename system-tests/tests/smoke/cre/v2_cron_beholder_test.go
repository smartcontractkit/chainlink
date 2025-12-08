package cre

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// smoke
func ExecuteCronBeholderTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	channel, listenerCtx, cancelFn, startErr := t_helpers.StartWorkflowEventsSubscriber(t.Context(), t_helpers.WorkflowEventsSubscriberConfig{
		WorkflowID:   workflowID,
		Don:          testEnv.Dons.MustWorkflowDON(),
		PollInterval: 5 * time.Second,
		F:            2,
		Logger:       testLogger,
	})
	require.NoError(t, startErr, "Failed to start workflow events subscriber")

	matcherFn := func(msg proto.Message) bool {
		if log, ok := msg.(*workflowevents.UserLogs); ok {
			for _, line := range log.LogLines {
				// we use contains, because the log line has additional context (e.g., timestamps) and some escape characters
				if strings.Contains(line.Message, "Amazing workflow user log") {
					return true
				}
			}
		}

		return false
	}

	err := t_helpers.AssertWorkflowEventMatched(listenerCtx, cancelFn, 2, channel, matcherFn, 2*time.Minute, testLogger)
	require.NoError(t, err, "Cron (Beholder) test failed")
	testLogger.Info().Msg("Cron (Beholder) test completed")
}
