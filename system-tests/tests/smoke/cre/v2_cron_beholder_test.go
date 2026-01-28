package cre

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"

	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	chiptestsink "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/chip-testsink"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// smoke
func ExecuteCronBeholderTest_Stateless(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	expectedBeholderLog := "Amazing workflow user log"
	logsFound := 0

	var publishFn chiptestsink.PublishFn = func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
		switch event.Type {
		case "workflows.v1.UserLogs":
			typedMsg := &workflowevents.UserLogs{}
			if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
				testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")

				return &chippb.PublishResponse{}, nil
			}

			for _, logLine := range typedMsg.LogLines {
				if strings.Contains(logLine.Message, expectedBeholderLog) {
					testLogger.Info().
						Str("expected_log", expectedBeholderLog).
						Str("found_message", strings.TrimSpace(logLine.Message)).
						Str("workflow_id", typedMsg.M.WorkflowExecutionID).
						Msg("🎯 Found expected user log message!")

					logsFound++

					return &chippb.PublishResponse{}, nil
				}

				testLogger.Warn().
					Str("expected_log", expectedBeholderLog).
					Str("found_message", strings.TrimSpace(logLine.Message)).
					Msg("[soft assertion] Received UserLogs message, but it does not match expected log")
			}
		default:
			// ignore
		}

		return &chippb.PublishResponse{}, nil
	}

	startCh := make(chan struct{}, 1)
	server, err := chiptestsink.NewServer(chiptestsink.Config{
		PublishFunc: publishFn,
		GRPCListen:  ":50051",
		Started:     startCh,
	})
	require.NoError(t, err, "failed to create new test sink server")

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()
	waitForServerStart(t, startCh, errCh)

	t.Cleanup(func() {
		server.Shutdown(t.Context())
	})

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	var timeout time.Duration = 2 * time.Minute

	ticker := time.NewTicker(2 * time.Second)
	ctx, cancelFn := context.WithTimeout(t.Context(), timeout)
	defer cancelFn()
	for {
		select {
		case <-ctx.Done():
			require.GreaterOrEqual(t, logsFound, 1, "Cron (Beholder) test failed")
		case <-ticker.C:
			if logsFound > 0 {
				ticker.Stop()
				testLogger.Info().Msg("Cron (Beholder) test completed")
				return
			}
			testLogger.Info().Msg("Waiting for expected user log message...")
		}
	}
}

const testSinkStartupTimeout = 10 * time.Second

func waitForServerStart(t *testing.T, started <-chan struct{}, errCh <-chan error) {
	t.Helper()

	select {
	case <-started:
		return
	case err := <-errCh:
		require.NoError(t, err, "test sink server failed while starting")
	case <-time.After(testSinkStartupTimeout):
		require.FailNow(t, "timeout waiting for test sink server to start")
	}
}
