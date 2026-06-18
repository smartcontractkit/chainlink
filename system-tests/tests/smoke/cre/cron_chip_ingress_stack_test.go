package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// smoke
func ExecuteCronChipIngressStackTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "cronchipingressstack")

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartChipIngressStack(t, testLogger, testEnv)

	t.Cleanup(func() {
		// stop ChIP Ingress after the test to free the port, on which other tests will start the ChiP Test Sink
		err := t_helpers.StopChipIngressStack(testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
		require.NoError(t, err, "Failed to stop Chip Ingress stack")
	})

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	expectedUserLog := "Amazing workflow user log"

	err := t_helpers.AssertChipIngressStackMessage(listenerCtx, t, expectedUserLog, testLogger, messageChan, kafkaErrChan, 4*time.Minute)
	require.NoError(t, err, "Cron (Chip Ingress stack) test failed")
	testLogger.Info().Msg("Cron (Chip Ingress stack) test completed")
}
