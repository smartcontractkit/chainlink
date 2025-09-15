package cre

import (
	"testing"
	"time"

	cronconfig "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteBeholderTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := cronconfig.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	expectedBeholderLog := "Amazing workflow user log"
	timeout := 2 * time.Minute
	assertBeholderMessage(t, listenerCtx, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)

	testLogger.Info().Msg("Beholder test completed")
}
