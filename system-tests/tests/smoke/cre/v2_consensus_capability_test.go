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

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")

	expectedBeholderLog := "Successfully passed all consensus tests"
	err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, 4*time.Minute)
	require.NoError(t, err, "Consensus capability test failed, Beholder should not return an error")
	testLogger.Info().Msg("Consensus capability test completed")
}
