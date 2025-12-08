package cre

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	consensus_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/consensus/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// regression
type consensusNegativeTest struct {
	name          string
	caseToTrigger string
	expectedError string
}

const (
	expectedConsensusError = "consensus calculation failed"
)

var consensusNegativeTestsGenerateReport = []consensusNegativeTest{
	// Consensus - generate report with random timestamps
	{"random timestamps", "Consensus - random timestamps", expectedConsensusError},
	{"inconsistent feedIDs", "Consensus - inconsistent feedIDs", expectedConsensusError},
	{"inconsistent prices", "Consensus - inconsistent prices", expectedConsensusError},
}

func ConsensusFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, consensusNegativeTest consensusNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./consensus/main.go"

	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID

		testLogger.Info().Msg("Creating Consensus Fail workflow configuration...")
		workflowName := fmt.Sprintf("consensus-fail-workflow-%s-%04d", chainID, rand.Intn(10000))
		feedID := "018e16c38e000320000000000000000000000000000000000000000000000000" // 32 hex characters (16 bytes)
		workflowConfig := consensus_negative_config.Config{
			CaseToTrigger: consensusNegativeTest.caseToTrigger,
			FeedID:        feedID,
		}
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		channel, listenerCtx, cancelFn, startErr := t_helpers.StartWorkflowEventsSubscriber(t.Context(), t_helpers.GetStandardWorkflowEventsSubscriberConfig(testEnv, workflowID))
		require.NoError(t, startErr, "Failed to start workflow events subscriber")

		timeout := 90 * time.Second
		expectedError := consensusNegativeTest.expectedError
		err := t_helpers.AssertWorkflowEventMatched(listenerCtx, cancelFn, 2, channel, t_helpers.GetUserLogMatcherFn(expectedError), timeout, testLogger)
		require.NoError(t, err, "Consensus Fail test failed")
		testLogger.Info().Msg("Consensus Fail test successfully completed")
	}
}
