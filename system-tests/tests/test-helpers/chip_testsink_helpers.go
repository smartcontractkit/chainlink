package helpers

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"

	chiptestsink "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/chip-testsink"
)

const testSinkStartupTimeout = 10 * time.Second

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
	var publishFn = func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
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

// GetLoggingPublishFn returns a CHiP publish handler that demuxes events into the provided channels and saves all events to a file.
// Useful when debugging failures of tests that depend on workflow logs.
func GetLoggingPublishFn(
	testLogger zerolog.Logger,
	userLogsCh chan *workflowevents.UserLogs,
	baseMessageCh chan *commonevents.BaseMessage,
	dumpFilePath string, // <--- best set to `./logs/your_file.txt` since `./logs` folder inside `smoke/cre` is uploaded as artifact in GH
) chiptestsink.PublishFn {
	// 1. Thread-safe helper to write generic proto messages to a file
	var fileMu sync.Mutex
	logToFile := func(eventType string, msg proto.Message) {
		if dumpFilePath == "" {
			return
		}

		fileMu.Lock()
		defer fileMu.Unlock()

		if dir := filepath.Dir(dumpFilePath); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				testLogger.Warn().Err(err).Str("path", dir).Msg("Failed to create dump directory")
				return
			}
		}

		// Serialize the proto message to JSON
		// Multiline: false ensures one event per line (easier to parse later as ndjson)
		dataBytes, err := (protojson.MarshalOptions{Multiline: false}).Marshal(msg)
		if err != nil {
			testLogger.Warn().Err(err).Str("type", eventType).Msg("Failed to marshal event for dump")
			return
		}

		// Wrap in a simple structure to preserve the event type in the log file
		entry := map[string]interface{}{
			"type":      eventType,
			"timestamp": time.Now(),
			"data":      json.RawMessage(dataBytes),
		}

		line, err := json.Marshal(entry)
		if err != nil {
			return
		}

		// Open file in Append mode
		f, err := os.OpenFile(dumpFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			testLogger.Warn().Err(err).Str("path", dumpFilePath).Msg("Failed to open dump file")
			return
		}
		defer f.Close()

		if _, err := f.Write(append(line, '\n')); err != nil {
			testLogger.Warn().Err(err).Msg("Failed to write to dump file")
		}
	}

	// Returns the actual PublishFn
	return func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
		// --- SWITCH 1: Data Persistence (Observability) ---
		if dumpFilePath != "" {
			var msgToSave proto.Message

			switch event.Type {
			// workflows.v1 events
			case "workflows.v1.CapabilityExecutionFinished":
				msgToSave = &workflowevents.CapabilityExecutionFinished{}
			case "workflows.v1.CapabilityExecutionStarted":
				msgToSave = &workflowevents.CapabilityExecutionStarted{}
			case "workflows.v1.MeteringReport":
				msgToSave = &workflowevents.MeteringReport{}
			case "workflows.v1.TransmissionsScheduledEvent":
				msgToSave = &workflowevents.TransmissionsScheduledEvent{}
			case "workflows.v1.TransmitScheduleEvent":
				msgToSave = &workflowevents.TransmitScheduleEvent{}
			case "workflows.v1.WorkflowExecutionFinished":
				msgToSave = &workflowevents.WorkflowExecutionFinished{}
			case "workflows.v1.WorkflowExecutionStarted":
				msgToSave = &workflowevents.WorkflowExecutionStarted{}
			case "workflows.v1.WorkflowStatusChanged":
				msgToSave = &workflowevents.WorkflowStatusChanged{}
			case "workflows.v1.UserLogs":
				msgToSave = &workflowevents.UserLogs{}

			// workflows.v2 events
			case "workflows.v2.CapabilityExecutionFinished":
				msgToSave = &workfloweventsv2.CapabilityExecutionFinished{}
			case "workflows.v2.CapabilityExecutionStarted":
				msgToSave = &workfloweventsv2.CapabilityExecutionStarted{}
			case "workflows.v2.TriggerExecutionStarted":
				msgToSave = &workfloweventsv2.TriggerExecutionStarted{}
			case "workflows.v2.WorkflowActivated":
				msgToSave = &workfloweventsv2.WorkflowActivated{}
			case "workflows.v2.WorkflowDeleted":
				msgToSave = &workfloweventsv2.WorkflowDeleted{}
			case "workflows.v2.WorkflowDeployed":
				msgToSave = &workfloweventsv2.WorkflowDeployed{}
			case "workflows.v2.WorkflowExecutionFinished":
				msgToSave = &workfloweventsv2.WorkflowExecutionFinished{}
			case "workflows.v2.WorkflowExecutionStarted":
				msgToSave = &workfloweventsv2.WorkflowExecutionStarted{}
			case "workflows.v2.WorkflowPaused":
				msgToSave = &workfloweventsv2.WorkflowPaused{}
			case "workflows.v2.WorkflowUpdated":
				msgToSave = &workfloweventsv2.WorkflowUpdated{}
			case "workflows.v2.WorkflowUserLog":
				msgToSave = &workfloweventsv2.WorkflowUserLog{}

			case "BaseMessage":
				msgToSave = &commonevents.BaseMessage{}
			default:
				// Optional: Log that we saw an unknown event type not in our save list?
			}

			if msgToSave != nil {
				// Unmarshal specifically for logging (safe to do redundantly for clarity)
				if err := proto.Unmarshal(event.GetProtoData().GetValue(), msgToSave); err == nil {
					logToFile(event.Type, msgToSave)
				}
			}
		}

		// --- SWITCH 2: Test Orchestration (Logic) ---
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
}

// StartChipTestSink boots the CHiP test sink and waits until it is accepting traffic.
func StartChipTestSink(t *testing.T, publishFn chiptestsink.PublishFn) *chiptestsink.Server {
	grpcListenAddr := ":" + chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT
	if !isPortAvailable(grpcListenAddr) {
		t.Fatalf(`failed to start ChIP Ingress Test Sink. Port %s is already taken. Most probably an instance of ChIP Ingress is already running.
If you want to use both together start ChIP Ingress on a different port with '--grpc-port' flag
and make sure that the sink is pointing to correct upstream endpoint ('localhost:<grpc-port>' in most cases)`, chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT)
	}

	startCh := make(chan struct{}, 1)
	server, err := chiptestsink.NewServer(chiptestsink.Config{
		PublishFunc: publishFn,
		GRPCListen:  grpcListenAddr,
		Started:     startCh, // signals that server is indeed listening on the GRPC port
		// UpstreamEndpoint: "localhost:50052", // uncomment to forward events to ChIP, remember to start ChIP on a different port config.DefaultChipIngressPort (=50051)
	})
	require.NoError(t, err, "failed to create new test sink server")

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	select {
	case <-startCh:
	case err := <-errCh:
		require.NoError(t, err, "test sink server failed while starting")
	case <-time.After(testSinkStartupTimeout):
		require.FailNow(t, "timeout waiting for test sink server to start")
	}

	return server
}

func isPortAvailable(addr string) bool {
	lc := net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return false // already in use or permission denied
	}
	_ = l.Close()
	return true
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

// WaitForBaseMessage blocks until the base message channel emits a message containing needle.
func WaitForBaseMessage(
	ctx context.Context,
	testLogger zerolog.Logger,
	publishCh <-chan *commonevents.BaseMessage,
	needle string,
) (*commonevents.BaseMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case msg := <-publishCh:
			if strings.Contains(msg.Msg, needle) {
				return msg, nil
			}
			if strings.Contains(msg.Msg, "heartbeat") {
				continue
			}
			testLogger.Warn().
				Str("expected_log", needle).
				Str("found_message", strings.TrimSpace(msg.Msg)).
				Msg("[soft assertion] Received BaseMessage message, but it does not match expected log")
		}
	}
}

// WatchBaseMessages requires that the expected base message arrives before the timeout.
func WatchBaseMessages(
	t *testing.T,
	testLogger zerolog.Logger,
	baseMessageCh <-chan *commonevents.BaseMessage,
	expectedMessage string,
	timeout time.Duration,
) *commonevents.BaseMessage {
	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, errors.New("failed to find expected base message"))
	defer cancelFn()

	msg, err := WaitForBaseMessage(ctx, testLogger, baseMessageCh, expectedMessage)
	require.NoError(t, err, "failed to find expected base message")

	return msg
}

// IgnoreUserLogs drains user log traffic so publishers never block when tests do not care about logs.
func IgnoreUserLogs(ctx context.Context, userLogsCh <-chan *workflowevents.UserLogs) {
	go func() {
		defer func() { _ = recover() }() // in case channel closes
		for {
			select {
			case <-ctx.Done():
				return
			case <-userLogsCh:
				// noop
			}
		}
	}()
}
