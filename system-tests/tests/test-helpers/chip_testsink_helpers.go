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

type ChipSink interface {
	Shutdown(ctx context.Context)
}

type baseMessageWatchCfg struct {
	workflowID string
	labelEq    map[string]string
	labelIn    map[string]map[string]struct{}
	labelHas   map[string]string
}

type userLogWatchCfg struct {
	workflowID string
}

// BaseMessageWatchOpt customizes base message watchers.
type BaseMessageWatchOpt func(*baseMessageWatchCfg)

// UserLogWatchOpt customizes user log watchers.
type UserLogWatchOpt func(*userLogWatchCfg)

// WithBaseMessageWorkflowID filters base messages to a specific workflow ID.
func WithBaseMessageWorkflowID(workflowID string) BaseMessageWatchOpt {
	return func(cfg *baseMessageWatchCfg) {
		cfg.workflowID = normalizeWorkflowID(workflowID)
	}
}

// WithBaseMessageLabelEquals requires a base message label to be exactly equal to value.
func WithBaseMessageLabelEquals(key, value string) BaseMessageWatchOpt {
	return func(cfg *baseMessageWatchCfg) {
		k := strings.TrimSpace(key)
		if k == "" {
			return
		}
		if cfg.labelEq == nil {
			cfg.labelEq = make(map[string]string)
		}
		cfg.labelEq[k] = strings.TrimSpace(value)
	}
}

// WithBaseMessageLabelIn requires a base message label to match one of the allowed values.
func WithBaseMessageLabelIn(key string, allowedValues ...string) BaseMessageWatchOpt {
	return func(cfg *baseMessageWatchCfg) {
		k := strings.TrimSpace(key)
		if k == "" || len(allowedValues) == 0 {
			return
		}
		if cfg.labelIn == nil {
			cfg.labelIn = make(map[string]map[string]struct{})
		}
		values := make(map[string]struct{}, len(allowedValues))
		for _, v := range allowedValues {
			values[strings.TrimSpace(v)] = struct{}{}
		}
		cfg.labelIn[k] = values
	}
}

// WithBaseMessageLabelContains requires a base message label to contain the provided substring.
func WithBaseMessageLabelContains(key, substring string) BaseMessageWatchOpt {
	return func(cfg *baseMessageWatchCfg) {
		k := strings.TrimSpace(key)
		if k == "" {
			return
		}
		if cfg.labelHas == nil {
			cfg.labelHas = make(map[string]string)
		}
		cfg.labelHas[k] = strings.TrimSpace(substring)
	}
}

// WithUserLogWorkflowID filters user logs to a specific workflow ID.
func WithUserLogWorkflowID(workflowID string) UserLogWatchOpt {
	return func(cfg *userLogWatchCfg) {
		cfg.workflowID = normalizeWorkflowID(workflowID)
	}
}

type fanoutSubscription struct {
	id string
}

func (s *fanoutSubscription) Shutdown(_ context.Context) {
	fanoutSubMu.Lock()
	defer fanoutSubMu.Unlock()
	delete(fanoutSubs, s.id)
}

var (
	fanoutOnce   sync.Once
	fanoutServer *chiptestsink.Server
	errFanout    error

	fanoutSubMu sync.Mutex
	fanoutSubs  = make(map[string]chiptestsink.PublishFn)
)

func safeSendUserLogs(ch chan *workflowevents.UserLogs, msg *workflowevents.UserLogs) {
	// In fanout mode, tests may close their log channels immediately after
	// unsubscribing during cleanup. An in-flight publish can race with that close,
	// which would panic on send. We recover to treat delivery as best-effort.
	defer func() { _ = recover() }()
	if ch == nil {
		return
	}
	ch <- msg
}

func safeSendBaseMessage(ch chan *commonevents.BaseMessage, msg *commonevents.BaseMessage) {
	// Same race as safeSendUserLogs; avoid panic on send to closed channel.
	defer func() { _ = recover() }()
	if ch == nil {
		return
	}
	ch <- msg
}

func safeSendProtoMessage(ch chan proto.Message, msg proto.Message) {
	// Same race as safeSendUserLogs; avoid panic on send to closed channel.
	defer func() { _ = recover() }()
	if ch == nil {
		return
	}
	ch <- msg
}

func ChipSinkFanoutEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("CRE_TEST_CHIP_SINK_FANOUT_ENABLED")))
	return v == "1" || v == "true" || v == "yes"
}

func ensureFanoutServer(t *testing.T) {
	t.Helper()

	fanoutOnce.Do(func() {
		grpcListenAddr := ":" + chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT
		startCh := make(chan struct{}, 1)
		fanoutServer, errFanout = chiptestsink.NewServer(chiptestsink.Config{
			GRPCListen: grpcListenAddr,
			Started:    startCh,
			PublishFunc: func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
				fanoutSubMu.Lock()
				snapshot := make([]chiptestsink.PublishFn, 0, len(fanoutSubs))
				for _, fn := range fanoutSubs {
					snapshot = append(snapshot, fn)
				}
				fanoutSubMu.Unlock()

				for _, fn := range snapshot {
					if _, err := fn(ctx, event); err != nil {
						// Best-effort delivery: one subscriber must not fail all.
						continue
					}
				}
				return &chippb.PublishResponse{}, nil
			},
		})
		if errFanout != nil {
			return
		}

		errCh := make(chan error, 1)
		go func() {
			errCh <- fanoutServer.Run()
		}()

		select {
		case <-startCh:
		case err := <-errCh:
			errFanout = err
		case <-time.After(testSinkStartupTimeout):
			errFanout = errors.New("timeout waiting for fanout sink server to start")
		}
	})

	require.NoError(t, errFanout, "failed to start fanout sink server")
}

// WaitForUserLog monitors workflow user logs until one contains needle or the context ends.
func WaitForUserLog(
	ctx context.Context,
	testLogger zerolog.Logger,
	publishCh <-chan *workflowevents.UserLogs,
	needle string,
	opts ...UserLogWatchOpt,
) (*workflowevents.LogLine, error) {
	cfg := userLogWatchCfg{}
	for _, opt := range opts {
		opt(&cfg)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case logs := <-publishCh:
			if logs == nil {
				continue
			}
			if cfg.workflowID != "" && !userLogsHasWorkflowID(logs, cfg.workflowID) {
				continue
			}
			for _, line := range logs.LogLines {
				if strings.Contains(line.Message, needle) {
					testLogger.Info().Str("expected_log", needle).Str("actual_log", strings.TrimSpace(line.Message)).Msg("Found expected user log")
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
	opts ...BaseMessageWatchOpt,
) {
	t.Helper()
	cfg := baseMessageWatchCfg{}
	for _, opt := range opts {
		opt(&cfg)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-publishCh:
			// Channel can be closed during cleanup; closed or nil messages should exit.
			if !ok || msg == nil {
				return
			}
			if cfg.workflowID != "" && !baseMessageHasWorkflowID(msg, cfg.workflowID) {
				continue
			}
			if strings.Contains(msg.Msg, needle) {
				ok, reason := baseMessageMatchesLabelFilters(msg, cfg)
				if !ok {
					testLogger.Warn().
						Str("expected_log", needle).
						Str("found_message", strings.TrimSpace(msg.Msg)).
						Str("filter_reason", reason).
						Msg("[soft assertion] Ignoring poison BaseMessage because source labels do not match")
					continue
				}
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

			safeSendUserLogs(userLogsCh, typedMsg)
			return &chippb.PublishResponse{}, nil

		case "BaseMessage":
			typedMsg := &commonevents.BaseMessage{}
			if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
				testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")

				return &chippb.PublishResponse{}, nil
			}
			safeSendBaseMessage(baseMessageCh, typedMsg)
			return &chippb.PublishResponse{}, nil
		default:
			// ignore
		}

		return &chippb.PublishResponse{}, nil
	}

	return publishFn
}

// GetWorkflowV2LifecyclePublishFn returns a CHiP publish handler that demuxes workflow lifecycle v2 events.
func GetWorkflowV2LifecyclePublishFn(testLogger zerolog.Logger, workflowEventCh chan proto.Message) chiptestsink.PublishFn {
	return func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
		var typedMsg proto.Message
		switch event.Type {
		case "workflows.v2.WorkflowActivated":
			typedMsg = &workfloweventsv2.WorkflowActivated{}
		case "workflows.v2.WorkflowPaused":
			typedMsg = &workfloweventsv2.WorkflowPaused{}
		case "workflows.v2.WorkflowDeleted":
			typedMsg = &workfloweventsv2.WorkflowDeleted{}
		default:
			return &chippb.PublishResponse{}, nil
		}

		if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
			testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")
			return &chippb.PublishResponse{}, nil
		}

		safeSendProtoMessage(workflowEventCh, typedMsg)
		return &chippb.PublishResponse{}, nil
	}
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
			safeSendUserLogs(userLogsCh, typedMsg)
			return &chippb.PublishResponse{}, nil

		case "BaseMessage":
			typedMsg := &commonevents.BaseMessage{}
			if err := proto.Unmarshal(event.GetProtoData().GetValue(), typedMsg); err != nil {
				testLogger.Error().Err(err).Str("ce_type", event.Type).Msg("Failed to unmarshal protobuf; skipping")
				return &chippb.PublishResponse{}, nil
			}
			safeSendBaseMessage(baseMessageCh, typedMsg)
			return &chippb.PublishResponse{}, nil

		default:
			// ignore
		}

		return &chippb.PublishResponse{}, nil
	}
}

// StartChipTestSink boots the CHiP test sink and waits until it is accepting traffic.
// In fanout mode (CRE_TEST_CHIP_SINK_FANOUT_ENABLED=1), a singleton sink is started and each test
// registers its own publish function as a fanout subscriber.
func StartChipTestSink(t *testing.T, publishFn chiptestsink.PublishFn) ChipSink {
	if ChipSinkFanoutEnabled() {
		ensureFanoutServer(t)
		subID := t.Name() + "-" + time.Now().Format("150405.000000000")
		fanoutSubMu.Lock()
		fanoutSubs[subID] = publishFn
		fanoutSubMu.Unlock()
		return &fanoutSubscription{id: subID}
	}

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
	timeout time.Duration,
	opts ...UserLogWatchOpt,
) {
	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, errors.New("failed to find expected user log message"))
	defer cancelFn()

	cancelCtx, cancelCauseFn := context.WithCancelCause(ctx)
	defer cancelCauseFn(nil)
	userCfg := userLogWatchCfg{}
	for _, opt := range opts {
		opt(&userCfg)
	}

	if failingBeholderLog != "" {
		go func() {
			if userCfg.workflowID != "" {
				FailOnBaseMessage(cancelCtx, cancelCauseFn, t, testLogger, baseMessageCh, failingBeholderLog, WithBaseMessageWorkflowID(userCfg.workflowID))
				return
			}
			FailOnBaseMessage(cancelCtx, cancelCauseFn, t, testLogger, baseMessageCh, failingBeholderLog)
		}()
	}
	_, err := WaitForUserLog(cancelCtx, testLogger, userLogsCh, expectedBeholderLog, opts...)
	require.NoError(t, err, "failed to find expected user log message")
}

// WaitForBaseMessage blocks until the base message channel emits a message containing needle.
func WaitForBaseMessage(
	ctx context.Context,
	testLogger zerolog.Logger,
	publishCh <-chan *commonevents.BaseMessage,
	needle string,
	opts ...BaseMessageWatchOpt,
) (*commonevents.BaseMessage, error) {
	cfg := baseMessageWatchCfg{}
	for _, opt := range opts {
		opt(&cfg)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case msg := <-publishCh:
			if msg == nil {
				continue
			}
			if cfg.workflowID != "" && !baseMessageHasWorkflowID(msg, cfg.workflowID) {
				continue
			}
			if strings.Contains(msg.Msg, needle) {
				ok, reason := baseMessageMatchesLabelFilters(msg, cfg)
				if !ok {
					testLogger.Warn().
						Str("expected_log", needle).
						Str("found_message", strings.TrimSpace(msg.Msg)).
						Str("filter_reason", reason).
						Msg("[soft assertion] Received BaseMessage with expected message, but source labels do not match")
					continue
				}
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
	opts ...BaseMessageWatchOpt,
) *commonevents.BaseMessage {
	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, errors.New("failed to find expected base message"))
	defer cancelFn()

	msg, err := WaitForBaseMessage(ctx, testLogger, baseMessageCh, expectedMessage, opts...)
	require.NoError(t, err, "failed to find expected base message")
	testLogger.Info().Msgf("Found expected base message: %s", expectedMessage)

	return msg
}

func normalizeWorkflowID(workflowID string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(workflowID)), "0x")
}

func baseMessageHasWorkflowID(msg *commonevents.BaseMessage, workflowID string) bool {
	if msg == nil || workflowID == "" {
		return false
	}

	for _, key := range []string{"workflowID", "workflow_id", "workflowId"} {
		if value, ok := msg.Labels[key]; ok && normalizeWorkflowID(value) == workflowID {
			return true
		}
	}

	// Some messages carry workflow id only inside the "err" label payload.
	if errLabel, ok := msg.Labels["err"]; ok && strings.Contains(strings.ToLower(errLabel), workflowID) {
		return true
	}

	return false
}

func baseMessageMatchesLabelFilters(msg *commonevents.BaseMessage, cfg baseMessageWatchCfg) (bool, string) {
	if msg == nil {
		return false, "base message is nil"
	}

	for key, expected := range cfg.labelEq {
		actual, ok := msg.Labels[key]
		if !ok {
			return false, "missing label " + key
		}
		if actual != expected {
			return false, "label " + key + " value mismatch"
		}
	}

	for key, allowedSet := range cfg.labelIn {
		actual, ok := msg.Labels[key]
		if !ok {
			return false, "missing label " + key
		}
		if _, ok := allowedSet[actual]; !ok {
			return false, "label " + key + " not in allowed values"
		}
	}

	for key, needle := range cfg.labelHas {
		actual, ok := msg.Labels[key]
		if !ok {
			return false, "missing label " + key
		}
		if !strings.Contains(actual, needle) {
			return false, "label " + key + " does not contain expected substring"
		}
	}

	return true, ""
}

func userLogsHasWorkflowID(logs *workflowevents.UserLogs, workflowID string) bool {
	if logs == nil || logs.M == nil || workflowID == "" {
		return false
	}
	return normalizeWorkflowID(logs.M.WorkflowID) == workflowID
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
