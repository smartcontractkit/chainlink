package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/chipsink"
)

const (
	defaultChipSinkName        = "default"
	defaultChipSinkGRPCListen  = "0.0.0.0:50051"
	defaultChipSinkEventsLimit = 200
	maxChipSinkEventsLimit     = 1000
)

func (s *Server) startChipTestSink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	var req ChipTestSinkStartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidRequestBody, fmt.Sprintf("invalid request body: %v", err), nil)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = defaultChipSinkName
	}
	grpcListen := strings.TrimSpace(req.GRPCListen)
	if grpcListen == "" {
		grpcListen = defaultChipSinkGRPCListen
	}
	normalizedListen, err := normalizeChipSinkListenAddress(grpcListen)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, err.Error(), nil)
		return
	}

	s.chipSinkMu.Lock()
	defer s.chipSinkMu.Unlock()

	if s.chipSink != nil {
		s.lggr.Info().
			Str("name", s.chipSink.name).
			Str("grpcListen", s.chipSink.grpcListen).
			Str("upstreamEndpoint", s.chipSink.upstreamEndpoint).
			Msg("chip test sink already running; returning existing status")
		s.respondJSONAny(w, http.StatusOK, ChipTestSinkStartResponse{
			Profile:          "sink",
			Mode:             "remote",
			Name:             s.chipSink.name,
			GRPCListen:       s.chipSink.grpcListen,
			UpstreamEndpoint: s.chipSink.upstreamEndpoint,
			EventLogPath:     s.chipSink.eventLogPath,
		})
		return
	}

	eventLogPath := defaultChipSinkEventLogPath()
	if mkdirErr := os.MkdirAll(filepath.Dir(eventLogPath), 0o755); mkdirErr != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to prepare chip sink log directory: %v", mkdirErr), nil)
		return
	}
	// Start with a clean event stream per launch.
	if removeErr := os.Remove(eventLogPath); removeErr != nil && !os.IsNotExist(removeErr) {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to reset chip sink event log: %v", removeErr), nil)
		return
	}
	var eventLogMu sync.Mutex

	started := make(chan string, 1)
	sinkServer, err := chipsink.NewServer(chipsink.Config{
		GRPCListen:       normalizedListen,
		UpstreamEndpoint: strings.TrimSpace(req.UpstreamEndpoint),
		Started:          started,
		PublishFn: func(_ context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
			_ = appendChipSinkEvent(eventLogPath, &eventLogMu, event)
			return &chippb.PublishResponse{}, nil
		},
	})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to create chip test sink server: %v", err), nil)
		return
	}

	runCtx, cancel := context.WithCancel(context.Background())
	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- sinkServer.Run()
	}()

	select {
	case boundAddr := <-started:
		s.lggr.Info().
			Str("name", name).
			Str("grpcListen", boundAddr).
			Str("upstreamEndpoint", strings.TrimSpace(req.UpstreamEndpoint)).
			Str("eventLogPath", eventLogPath).
			Msg("chip test sink started")
		s.chipSink = &chipTestSinkRuntime{
			name:             name,
			grpcListen:       boundAddr,
			upstreamEndpoint: strings.TrimSpace(req.UpstreamEndpoint),
			eventLogPath:     eventLogPath,
			server:           sinkServer,
			cancel:           cancel,
			runErrCh:         runErrCh,
		}
		s.storeRuntime(fmt.Sprintf("%s:%s", ComponentTypeChipTestSink, name), runtimeState{
			ComponentType: ComponentTypeChipTestSink,
			StopFn: func(ctx context.Context) error {
				sinkServer.Shutdown(ctx)
				cancel()
				return nil
			},
		})
	case err := <-runErrCh:
		cancel()
		s.lggr.Error().Err(err).Str("name", name).Str("grpcListen", normalizedListen).Msg("chip test sink failed to start")
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("chip test sink failed to start: %v", err), nil)
		return
	case <-time.After(10 * time.Second):
		cancel()
		sinkServer.Shutdown(context.Background())
		s.lggr.Error().Str("name", name).Str("grpcListen", normalizedListen).Msg("chip test sink startup timed out")
		s.respondError(w, http.StatusGatewayTimeout, ErrCodeDeployFailed, "timed out waiting for chip test sink to start", nil)
		return
	case <-runCtx.Done():
		cancel()
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, "chip test sink startup canceled", nil)
		return
	}

	s.respondJSONAny(w, http.StatusOK, ChipTestSinkStartResponse{
		Profile:          "sink",
		Mode:             "remote",
		Name:             name,
		GRPCListen:       s.chipSink.grpcListen,
		UpstreamEndpoint: s.chipSink.upstreamEndpoint,
		EventLogPath:     s.chipSink.eventLogPath,
	})
}

func (s *Server) stopChipTestSink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	s.chipSinkMu.Lock()
	defer s.chipSinkMu.Unlock()

	if s.chipSink == nil {
		s.lggr.Info().Msg("chip test sink stop requested; nothing running")
		s.respondJSONAny(w, http.StatusOK, ChipTestSinkStopResponse{Found: false, Stopped: false})
		return
	}

	runtime := s.chipSink
	s.lggr.Info().
		Str("name", runtime.name).
		Str("grpcListen", runtime.grpcListen).
		Msg("stopping chip test sink")
	runtime.server.Shutdown(r.Context())
	runtime.cancel()
	_, _ = s.takeRuntime(fmt.Sprintf("%s:%s", ComponentTypeChipTestSink, runtime.name))
	s.chipSink = nil
	s.lggr.Info().Str("name", runtime.name).Msg("chip test sink stopped")

	s.respondJSONAny(w, http.StatusOK, ChipTestSinkStopResponse{Found: true, Stopped: true})
}

func (s *Server) chipTestSinkStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	s.chipSinkMu.Lock()
	defer s.chipSinkMu.Unlock()
	if s.chipSink == nil {
		s.respondJSONAny(w, http.StatusOK, ChipTestSinkStatusResponse{
			Profile: "sink",
			Mode:    "remote",
			Running: false,
		})
		return
	}

	s.respondJSONAny(w, http.StatusOK, ChipTestSinkStatusResponse{
		Profile:          "sink",
		Mode:             "remote",
		Running:          true,
		Name:             s.chipSink.name,
		GRPCListen:       s.chipSink.grpcListen,
		UpstreamEndpoint: s.chipSink.upstreamEndpoint,
		EventLogPath:     s.chipSink.eventLogPath,
	})
}

func (s *Server) currentChipSinkStatus() *ChipTestSinkStatusResponse {
	s.chipSinkMu.Lock()
	defer s.chipSinkMu.Unlock()
	if s.chipSink == nil {
		return nil
	}
	return &ChipTestSinkStatusResponse{
		Profile:          "sink",
		Mode:             "remote",
		Running:          true,
		Name:             s.chipSink.name,
		GRPCListen:       s.chipSink.grpcListen,
		UpstreamEndpoint: s.chipSink.upstreamEndpoint,
		EventLogPath:     s.chipSink.eventLogPath,
	}
}

func (s *Server) chipTestSinkEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	s.chipSinkMu.Lock()
	runtime := s.chipSink
	s.chipSinkMu.Unlock()
	if runtime == nil {
		s.respondJSONAny(w, http.StatusOK, ChipTestSinkEventsResponse{Events: []ChipTestSinkEventLogEntry{}})
		return
	}

	limit := defaultChipSinkEventsLimit
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, "limit query parameter must be a positive integer", nil)
			return
		}
		if parsed > maxChipSinkEventsLimit {
			parsed = maxChipSinkEventsLimit
		}
		limit = parsed
	}

	var since time.Time
	if rawSince := strings.TrimSpace(r.URL.Query().Get("since")); rawSince != "" {
		parsed, err := time.Parse(time.RFC3339Nano, rawSince)
		if err != nil {
			s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, "since query parameter must be RFC3339Nano timestamp", nil)
			return
		}
		since = parsed
	}

	events, err := readChipSinkEvents(runtime.eventLogPath, since, limit)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to read chip sink events: %v", err), nil)
		return
	}
	s.respondJSONAny(w, http.StatusOK, ChipTestSinkEventsResponse{Events: events})
}

func defaultChipSinkEventLogPath() string {
	return filepath.Join(os.TempDir(), "cre-agent-chip-sink-events.ndjson")
}

func normalizeChipSinkListenAddress(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return defaultChipSinkGRPCListen, nil
	}
	// Accept bare port for convenience, e.g. "50052".
	if _, err := strconv.Atoi(trimmed); err == nil {
		return net.JoinHostPort("0.0.0.0", trimmed), nil
	}
	// Accept ":50052" and normalize to explicit host.
	if strings.HasPrefix(trimmed, ":") {
		return net.JoinHostPort("0.0.0.0", strings.TrimPrefix(trimmed, ":")), nil
	}
	_, port, err := net.SplitHostPort(trimmed)
	if err != nil || strings.TrimSpace(port) == "" {
		return "", fmt.Errorf("invalid grpcListen %q: expected host:port or port", raw)
	}
	return trimmed, nil
}

func appendChipSinkEvent(path string, mu *sync.Mutex, event *pb.CloudEvent) error {
	if event == nil {
		return nil
	}
	entry := ChipTestSinkEventLogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Type:      strings.TrimSpace(event.Type),
		Event:     chipsink.EventData(event),
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(line, '\n')); err != nil {
		return err
	}
	return nil
}

func readChipSinkEvents(path string, since time.Time, limit int) ([]ChipTestSinkEventLogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ChipTestSinkEventLogEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	events := make([]ChipTestSinkEventLogEntry, 0, limit)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry ChipTestSinkEventLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if !since.IsZero() {
			ts, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(entry.Timestamp))
			if err != nil || !ts.After(since) {
				continue
			}
		}
		events = append(events, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(events) <= limit {
		return events, nil
	}
	return events[len(events)-limit:], nil
}
