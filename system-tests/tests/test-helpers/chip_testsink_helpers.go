package helpers

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	chiptestsink "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/chip-testsink"
)

const testSinkStartupTimeout = 10 * time.Second

// WaitForServerStart blocks until the CHiP test sink reports it is ready or fails fast.
func WaitForServerStart(t *testing.T, started <-chan struct{}, errCh <-chan error) {
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

// WaitForUserLog monitors workflow user logs until one contains needle or the context ends.
func WaitForUserLog(
	ctx context.Context,
	testLogger zerolog.Logger,
	publishCh <-chan *workflowevents.UserLogs,
	needle string,
) (*workflowevents.LogLine, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case logs := <-publishCh:
			for _, line := range logs.LogLines {
				if strings.Contains(line.Message, needle) {
					return line, nil
				}

				testLogger.Warn().
					Str("expected_log", needle).
					Str("found_message", strings.TrimSpace(line.Message)).
					Msg("[soft assertion] Received UserLogs message, but it does not match expected log")
			}
		}
	}
}

// FailOnBaseMessage cancels the supplied context as soon as a poison base message is observed.
func FailOnBaseMessage(
	ctx context.Context,
	cancelCause context.CancelCauseFunc,
	t *testing.T,
	testLogger zerolog.Logger,
	publishCh <-chan *commonevents.BaseMessage,
	needle string,
) {
	t.Helper()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-publishCh:
			if strings.Contains(msg.Msg, needle) {
				testLogger.Error().
					Str("expected_log", needle).
					Str("found_message", strings.TrimSpace(msg.Msg)).
					Msg("Found unexpected base message")
				cancelCause(errors.New("found unexpected base message: " + msg.Msg))
				t.FailNow()
			}
		}
	}
}

// GetPublishFn returns a CHiP publish handler that demuxes events into the provided channels.
func GetPublishFn(testLogger zerolog.Logger, userLogsCh chan *workflowevents.UserLogs, baseMessageCh chan *commonevents.BaseMessage) chiptestsink.PublishFn {
	var publishFn chiptestsink.PublishFn = func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
		switch event.Type {
		case "workflows.v1.UserLogs":
			typedMsg := &workflowevents.UserLogs{}
			if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
				testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")

				return &chippb.PublishResponse{}, nil
			}

			userLogsCh <- typedMsg
			return &chippb.PublishResponse{}, nil

		case "BaseMessage":
			typedMsg := &commonevents.BaseMessage{}
			if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
				testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")

				return &chippb.PublishResponse{}, nil
			}
			baseMessageCh <- typedMsg
			return &chippb.PublishResponse{}, nil
		default:
			// ignore
		}

		return &chippb.PublishResponse{}, nil
	}

	return publishFn
}

// StartChipTestSink boots the CHiP test sink and waits until it is accepting traffic.
func StartChipTestSink(t *testing.T, publishFn chiptestsink.PublishFn) *chiptestsink.Server {
	startCh := make(chan struct{}, 1)
	server, err := chiptestsink.NewServer(chiptestsink.Config{
		PublishFunc: publishFn,
		GRPCListen:  ":" + strconv.Itoa(config.DefaultChipIngressPort),
		Started:     startCh,
	})
	require.NoError(t, err, "failed to create new test sink server")

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()
	WaitForServerStart(t, startCh, errCh)

	return server
}

// WatchWorkflowLogs enforces that the expected log appears before timeout and that poison logs abort the test.
func WatchWorkflowLogs(
	t *testing.T,
	testLogger zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	failingBeholderLog string,
	expectedBeholderLog string,
	timeout time.Duration) {
	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, errors.New("failed to find expected user log message"))
	defer cancelFn()

	cancelCtx, cancelCauseFn := context.WithCancelCause(ctx)
	defer cancelCauseFn(nil)

	if failingBeholderLog != "" {
		go func() {
			FailOnBaseMessage(cancelCtx, cancelCauseFn, t, testLogger, baseMessageCh, failingBeholderLog)
		}()
	}
	_, err := WaitForUserLog(cancelCtx, testLogger, userLogsCh, expectedBeholderLog)
	require.NoError(t, err, "failed to find expected user log message")
}
