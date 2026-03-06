package cre

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/stretchr/testify/require"

	maxlimitsconfig "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/maxlimits/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteMaxLimitsTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	ballastPath := filepath.Join("maxlimits", "ballast.dat")
	generateBallast(t, ballastPath, 14*1024*1024)
	t.Cleanup(func() { os.Remove(ballastPath) })

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	workflowConfig := maxlimitsconfig.Config{
		CronSchedule:         "*/30 * * * * *",
		ConsensusRounds:      1,
		ConsensusPayloadSize: 100,
		LogEventCount:        5,
	}

	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "maxlimitsworkflow",
		&workflowConfig, "./maxlimits/main.go")

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		"[maxlimits] cron handler completed all phases",
		4*time.Minute)
}

func generateBallast(t *testing.T, path string, sizeBytes int) {
	t.Helper()
	data := make([]byte, sizeBytes)
	_, err := rand.Read(data)
	require.NoError(t, err, "failed to generate ballast data")
	require.NoError(t, os.WriteFile(path, data, 0644), "failed to write ballast.dat")
}
