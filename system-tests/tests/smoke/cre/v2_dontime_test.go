package cre

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteDonTimeTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/time_consensus/main.go"
	workflowName := "timebeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &None{}, workflowFileLocation)

	expectedBeholderLog := "Verified consensus on DON Time"
	assertBeholderMessage(t, listenerCtx, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)
	testLogger.Info().Msg("DON Time test completed")
}
