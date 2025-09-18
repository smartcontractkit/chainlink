package cre

import (
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteBillingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbilling"

	billingState := getBillingAssertionState(t, testEnv.TestConfig.RelativePathToRepoRoot) // establish a baseline

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}

	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
	assertBillingStateChanged(t, billingState, timeout, 0)

	testLogger.Info().Msg("Billing test completed")
}
