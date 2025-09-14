package cre

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteBeholderTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &None{}, workflowFileLocation)

	expectedBeholderLog := "Amazing workflow user log"
	timeout := 2 * time.Minute
	assertBeholderMessage(t, listenerCtx, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)

	testLogger.Info().Msg("Beholder test completed")
}
