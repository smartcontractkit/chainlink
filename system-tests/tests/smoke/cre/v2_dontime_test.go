package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
)

func ExecuteDonTimeTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/time_consensus/main.go"
	workflowName := "timebeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	expectedBeholderLog := "Verified consensus on DON Time"
	err := assertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "DON Time test failed, Beholder should not return an error")
	testLogger.Info().Msg("DON Time test completed")
}
