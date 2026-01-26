package cre

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
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
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	testsinkminimal "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/testsink-minimal"
	testsinkstateful "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/testsink-stateful"
)

// smoke
func ExecuteCronBeholderTest_Stateless(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	expectedBeholderLog := "Amazing workflow user log"
	logsFound := 0

	var publishFn testsinkminimal.PublishFn = func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
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
	server, err := testsinkminimal.NewServer(testsinkminimal.Config{
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

func ExecuteCronBeholderTest_State(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	expectedBeholderLog := "Amazing workflow user log"

	startCh := make(chan struct{}, 1)
	server, err := testsinkstateful.NewServer(testsinkstateful.Config{
		GRPCListen: ":50051",
		HTTPListen: ":8082",
		Started:    startCh,
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

	var lastTimeStamp int
	logsFound := 0
	httpClient := &http.Client{Timeout: 30 * time.Second}
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

			query := url.Values{
				"timestamp": {strconv.Itoa(lastTimeStamp)},
				"limit":     {"100"},
				"entity":    {"workflows.v1.UserLogs"},
			}
			reqURL := fmt.Sprintf("http://localhost:8082/events?%s", query.Encode())

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
			if err != nil {
				testLogger.Error().Err(err).Msg("failed to create HTTP request")
				continue
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				testLogger.Error().Err(err).Msg("failed to perform HTTP request")
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				testLogger.Error().Msgf("unexpected status %d: %s", resp.StatusCode, string(body))
				continue
			}

			var results []testsinkstateful.ServedEvent
			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
				testLogger.Error().Err(err).Msg("failed to decode response")
				return
			}

			testLogger.Info().Msgf("Found %d events", len(results))

			slices.SortFunc(results, func(a, b testsinkstateful.ServedEvent) int {
				return a.Timestamp - b.Timestamp
			})

			for _, event := range results {
				typedMsg := &workflowevents.UserLogs{}
				raw, err := base64.StdEncoding.DecodeString(event.Body)
				if err != nil {
					testLogger.Error().Err(err).Str("schema", event.Schema).Msg("failed to base64 decode event body")
					continue
				}

				if err := proto.Unmarshal(raw, typedMsg); err != nil {
					testLogger.Error().Err(err).Str("schema", event.Schema).Msg("Failed to unmarshal protobuf; skipping")
					continue
				}

				for _, logLine := range typedMsg.LogLines {
					if strings.Contains(logLine.Message, expectedBeholderLog) {
						testLogger.Info().
							Str("expected_log", expectedBeholderLog).
							Str("found_message", strings.TrimSpace(logLine.Message)).
							Str("workflow_id", typedMsg.M.WorkflowExecutionID).
							Msg("🎯 Found expected user log message!")

						logsFound++

						continue
					}

					testLogger.Warn().
						Str("expected_log", expectedBeholderLog).
						Str("found_message", strings.TrimSpace(logLine.Message)).
						Msg("[soft assertion] Received UserLogs message, but it does not match expected log")
				}
			}

			if len(results) == 0 {
				continue
			}

			lastTimeStamp = results[len(results)-1].Timestamp
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
