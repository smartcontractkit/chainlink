package cre

import (
	"fmt"
	"testing"
	"time"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

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
	{"oversized payload", "Consensus - oversized payload", expectedConsensusError},
}

func ConsensusFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, consensusNegativeTest consensusNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./consensus/main.go"

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		chainID := bcOutput.CtfOutput().ChainID

		testLogger.Info().Msg("Creating Consensus Fail workflow configuration...")
		workflowName := t_helpers.UniqueWorkflowName(
			testEnv,
			fmt.Sprintf("consensus-fail-%s-%s", chainID, consensusNegativeTest.name),
		)
		feedID := "018e16c38e000320000000000000000000000000000000000000000000000000" // 32 hex characters (16 bytes)
		workflowConfig := consensus_negative_config.Config{
			CaseToTrigger: consensusNegativeTest.caseToTrigger,
			FeedID:        feedID,
			PayloadSizeKB: 101, // only used for oversized payload test
		}
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

		expectedError := consensusNegativeTest.expectedError

		t_helpers.WatchWorkflowLogs(
			t,
			testLogger,
			userLogsCh,
			baseMessageCh,
			t_helpers.WorkflowEngineInitErrorLog,
			expectedError,
			2*time.Minute,
			t_helpers.WithUserLogWorkflowID(workflowID),
		)
		testLogger.Info().Str("test case", consensusNegativeTest.name).Msg("Consensus Fail test successfully completed")
	}
}
