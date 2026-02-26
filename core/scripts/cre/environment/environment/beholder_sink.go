package environment

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/chipsink"
	remoteclient "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/client"
)

const (
	chipSinkStateFilename  = "chip_testsink.toml"
	chipSinkLogFilename    = "chip_testsink.log"
	chipSinkEventsFilename = "chip_testsink_events.ndjson"
	defaultLocalSinkListen = "127.0.0.1:50051"
)

type chipSinkLocalState struct {
	Version          int    `toml:"version"`
	PID              int    `toml:"pid"`
	GRPCListen       string `toml:"grpc_listen"`
	UpstreamEndpoint string `toml:"upstream_endpoint,omitempty"`
	EventLogPath     string `toml:"event_log_path,omitempty"`
	StartedAt        string `toml:"started_at,omitempty"`
}

func beholderSinkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sink",
		Short: "Manage chip test sink lifecycle",
	}
	cmd.AddCommand(startBeholderSinkCmd())
	cmd.AddCommand(stopBeholderSinkCmd())
	cmd.AddCommand(statusBeholderSinkCmd())
	cmd.AddCommand(eventsBeholderSinkCmd())
	cmd.AddCommand(runLocalBeholderSinkCmd())
	return cmd
}

func startBeholderSinkCmd() *cobra.Command {
	var placement, grpcListen, upstream string
	cmd := &cobra.Command{
		Use:              "start",
		Short:            "Start chip test sink (local or remote)",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch normalizePlacement(placement) {
			case "local":
				return startLocalChipSink(grpcListen, upstream)
			case "remote":
				return startRemoteChipSink(cmd.Context(), grpcListen, upstream)
			default:
				return fmt.Errorf("invalid placement %q (expected local or remote)", placement)
			}
		},
	}
	cmd.Flags().StringVar(&placement, "placement", "local", "Sink placement: local or remote")
	cmd.Flags().StringVar(&grpcListen, "grpc-listen", defaultLocalSinkListen, "Sink gRPC listen address")
	cmd.Flags().StringVar(&upstream, "upstream-endpoint", "", "Optional upstream Chip Ingress endpoint")
	return cmd
}

func stopBeholderSinkCmd() *cobra.Command {
	var placement string
	cmd := &cobra.Command{
		Use:              "stop",
		Short:            "Stop chip test sink (local or remote)",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch normalizePlacement(placement) {
			case "local":
				return stopLocalChipSink()
			case "remote":
				return stopRemoteChipSink(cmd.Context())
			default:
				return fmt.Errorf("invalid placement %q (expected local or remote)", placement)
			}
		},
	}
	cmd.Flags().StringVar(&placement, "placement", "local", "Sink placement: local or remote")
	return cmd
}

func statusBeholderSinkCmd() *cobra.Command {
	var placement string
	cmd := &cobra.Command{
		Use:              "status",
		Short:            "Show chip test sink status (local or remote)",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch normalizePlacement(placement) {
			case "local":
				return statusLocalChipSink()
			case "remote":
				return statusRemoteChipSink(cmd.Context())
			default:
				return fmt.Errorf("invalid placement %q (expected local or remote)", placement)
			}
		},
	}
	cmd.Flags().StringVar(&placement, "placement", "local", "Sink placement: local or remote")
	return cmd
}

func eventsBeholderSinkCmd() *cobra.Command {
	var (
		placement string
		limit     int
		sinceRaw  string
	)
	cmd := &cobra.Command{
		Use:              "events",
		Short:            "Read chip test sink events (local or remote)",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var since time.Time
			if strings.TrimSpace(sinceRaw) != "" {
				parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(sinceRaw))
				if err != nil {
					return fmt.Errorf("invalid --since value %q (expected RFC3339Nano)", sinceRaw)
				}
				since = parsed
			}
			switch normalizePlacement(placement) {
			case "local":
				return readLocalChipSinkEvents(since, limit)
			case "remote":
				return readRemoteChipSinkEvents(cmd.Context(), since, limit)
			default:
				return fmt.Errorf("invalid placement %q (expected local or remote)", placement)
			}
		},
	}
	cmd.Flags().StringVar(&placement, "placement", "local", "Sink placement: local or remote")
	cmd.Flags().IntVar(&limit, "limit", 200, "Max number of events to return")
	cmd.Flags().StringVar(&sinceRaw, "since", "", "Filter events after RFC3339Nano timestamp")
	return cmd
}

func runLocalBeholderSinkCmd() *cobra.Command {
	var grpcListen, upstream, eventsFile string
	cmd := &cobra.Command{
		Use:    "run-local",
		Short:  "Run local chip test sink server",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if strings.TrimSpace(eventsFile) == "" {
				return errors.New("events-file is required")
			}
			normalizedListen, err := normalizeLocalSinkListenAddress(grpcListen)
			if err != nil {
				return err
			}
			started := make(chan string, 1)
			var eventsMu sync.Mutex
			sinkServer, err := chipsink.NewServer(chipsink.Config{
				GRPCListen:       normalizedListen,
				UpstreamEndpoint: strings.TrimSpace(upstream),
				Started:          started,
				PublishFn: func(_ context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
					if err := appendLocalChipSinkEvent(eventsFile, &eventsMu, event); err != nil {
						framework.L.Warn().Err(err).Str("eventsFile", eventsFile).Msg("failed to append local chip sink event")
					}
					return &chippb.PublishResponse{}, nil
				},
			})
			if err != nil {
				return err
			}
			errCh := make(chan error, 1)
			go func() {
				errCh <- sinkServer.Run()
			}()

			select {
			case addr := <-started:
				framework.L.Info().Str("grpcListen", addr).Msg("local chip test sink started")
				fmt.Printf("local chip sink started: grpcListen=%s eventsFile=%s\n", addr, eventsFile)
			case err := <-errCh:
				return err
			case <-time.After(10 * time.Second):
				sinkServer.Shutdown(context.Background())
				return errors.New("timed out waiting for local chip test sink to start")
			}

			sigCtx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			select {
			case <-sigCtx.Done():
				sinkServer.Shutdown(context.Background())
				fmt.Printf("local chip sink stopped: eventsFile=%s\n", eventsFile)
				return nil
			case err := <-errCh:
				return err
			}
		},
	}
	cmd.Flags().StringVar(&grpcListen, "grpc-listen", defaultLocalSinkListen, "Sink gRPC listen address")
	cmd.Flags().StringVar(&upstream, "upstream-endpoint", "", "Optional upstream Chip Ingress endpoint")
	cmd.Flags().StringVar(&eventsFile, "events-file", "", "Path to NDJSON file with captured sink events")
	return cmd
}

func startLocalChipSink(grpcListen, upstream string) error {
	normalizedListen, err := normalizeLocalSinkListenAddress(grpcListen)
	if err != nil {
		return err
	}
	existing, err := loadChipSinkLocalState()
	if err == nil && existing.PID > 0 && processExists(existing.PID) {
		framework.L.Info().Int("pid", existing.PID).Str("grpcListen", existing.GRPCListen).Str("eventsFile", existing.EventLogPath).Msg("local chip test sink already running")
		fmt.Printf("local chip sink already running: pid=%d grpcListen=%s eventsFile=%s\n", existing.PID, existing.GRPCListen, existing.EventLogPath)
		return nil
	}

	executablePath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "resolve executable path for local chip sink")
	}
	statePath := chipSinkStatePath()
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return errors.Wrap(err, "create chip sink state directory")
	}
	logPath := filepath.Join(filepath.Dir(statePath), chipSinkLogFilename)
	eventsPath := chipSinkEventsPath()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return errors.Wrap(err, "open chip sink log file")
	}
	defer logFile.Close()
	if err := os.Remove(eventsPath); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "reset local chip sink events file")
	}

	args := []string{"env", "beholder", "sink", "run-local", "--grpc-listen", normalizedListen, "--events-file", eventsPath}
	if strings.TrimSpace(upstream) != "" {
		args = append(args, "--upstream-endpoint", strings.TrimSpace(upstream))
	}
	cmd := exec.Command(executablePath, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "start local chip sink process")
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()
	if !waitForPIDAlive(pid, 1500*time.Millisecond) {
		return fmt.Errorf("local chip sink process exited too quickly (pid=%d)", pid)
	}
	if err := waitForLocalSinkReady(pid, normalizedListen, 5*time.Second, logPath); err != nil {
		_ = stopPID(pid)
		return err
	}
	if err := storeChipSinkLocalState(&chipSinkLocalState{
		Version:          1,
		PID:              pid,
		GRPCListen:       normalizedListen,
		UpstreamEndpoint: strings.TrimSpace(upstream),
		EventLogPath:     eventsPath,
		StartedAt:        time.Now().UTC().Format(time.RFC3339Nano),
	}); err != nil {
		return err
	}
	fmt.Printf("local chip sink started in background: pid=%d grpcListen=%s eventsFile=%s\n", pid, normalizedListen, eventsPath)
	return nil
}

func stopLocalChipSink() error {
	state, err := loadChipSinkLocalState()
	if err != nil {
		if os.IsNotExist(err) {
			framework.L.Info().Msg("local chip test sink is not running")
			return nil
		}
		return err
	}
	if state.PID <= 0 || !processExists(state.PID) {
		return removeChipSinkLocalState()
	}
	proc, err := os.FindProcess(state.PID)
	if err != nil {
		return err
	}
	_ = proc.Signal(syscall.SIGTERM)
	deadline := time.Now().Add(2 * time.Second)
	for processExists(state.PID) && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}
	if processExists(state.PID) {
		_ = proc.Signal(syscall.SIGKILL)
	}
	if processExists(state.PID) {
		return fmt.Errorf("local chip sink pid=%d did not stop", state.PID)
	}
	fmt.Printf("local chip sink stopped: pid=%d eventsFile=%s\n", state.PID, state.EventLogPath)
	return removeChipSinkLocalState()
}

func statusLocalChipSink() error {
	state, err := loadChipSinkLocalState()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("chip sink status: local running=false")
			return nil
		}
		return err
	}
	running := state.PID > 0 && processExists(state.PID)
	if !running {
		fmt.Printf("chip sink status: local running=false pid=%d grpcListen=%s eventsFile=%s (stale state)\n", state.PID, state.GRPCListen, state.EventLogPath)
		return nil
	}
	fmt.Printf("chip sink status: local running=true pid=%d grpcListen=%s eventsFile=%s\n", state.PID, state.GRPCListen, state.EventLogPath)
	return nil
}

func startRemoteChipSink(ctx context.Context, grpcListen, upstream string) error {
	runtime, err := remoteclient.ResolveRuntime(framework.L)
	if err != nil {
		return err
	}
	resp, err := remoteclient.StartRemoteChipTestSink(ctx, runtime, agent.ChipTestSinkStartRequest{
		Name:             "default",
		GRPCListen:       grpcListen,
		UpstreamEndpoint: strings.TrimSpace(upstream),
	})
	if err != nil {
		return err
	}
	if err := storeRemoteAgentStateSnapshot(relativePathToRepoRoot); err != nil {
		framework.L.Warn().Err(err).Msg("failed to persist remote agent state snapshot")
	}
	fmt.Printf("chip sink status: remote running=true grpcListen=%s\n", resp.GRPCListen)
	return nil
}

func stopRemoteChipSink(ctx context.Context) error {
	return withResolvedRemoteRuntime(ctx, func(ctx context.Context, runtime *remoteclient.Runtime) error {
		resp, err := remoteclient.StopRemoteChipTestSink(ctx, runtime)
		if err != nil {
			return err
		}
		fmt.Printf("chip sink stop: remote found=%t stopped=%t\n", resp.Found, resp.Stopped)
		return nil
	})
}

func statusRemoteChipSink(ctx context.Context) error {
	return withResolvedRemoteRuntime(ctx, func(ctx context.Context, runtime *remoteclient.Runtime) error {
		resp, err := remoteclient.GetRemoteChipTestSinkStatus(ctx, runtime)
		if err != nil {
			return err
		}
		fmt.Printf("chip sink status: remote running=%t grpcListen=%s\n", resp.Running, resp.GRPCListen)
		return nil
	})
}

func normalizePlacement(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "local":
		return "local"
	case "remote":
		return "remote"
	default:
		return strings.ToLower(strings.TrimSpace(v))
	}
}

func chipSinkStatePath() string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, envconfig.StateDirname, chipSinkStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for chip sink state file: %w", err))
	}
	return absPath
}

func chipSinkEventsPath() string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, envconfig.StateDirname, chipSinkEventsFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for chip sink events file: %w", err))
	}
	return absPath
}

func loadChipSinkLocalState() (*chipSinkLocalState, error) {
	data, err := os.ReadFile(chipSinkStatePath())
	if err != nil {
		return nil, err
	}
	state := &chipSinkLocalState{}
	if err := toml.Unmarshal(data, state); err != nil {
		return nil, err
	}
	return state, nil
}

func storeChipSinkLocalState(state *chipSinkLocalState) error {
	data, err := toml.Marshal(state)
	if err != nil {
		return err
	}
	path := chipSinkStatePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func removeChipSinkLocalState() error {
	if err := os.Remove(chipSinkStatePath()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func appendLocalChipSinkEvent(path string, mu *sync.Mutex, event *pb.CloudEvent) error {
	if event == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	eventData := localChipSinkEventData(event)
	entry := map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"type":      strings.TrimSpace(event.Type),
		"data":      eventData,
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

func readLocalChipSinkEvents(since time.Time, limit int) error {
	eventsPath := chipSinkEventsPath()
	file, err := os.Open(eventsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return printDebugJSON(map[string]any{"events": []any{}})
		}
		return err
	}
	defer file.Close()

	events := make([]map[string]any, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if !since.IsZero() {
			tsRaw, _ := entry["timestamp"].(string)
			ts, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(tsRaw))
			if err != nil || !ts.After(since) {
				continue
			}
		}
		events = append(events, entry)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if limit > 0 && len(events) > limit {
		events = events[len(events)-limit:]
	}
	return printDebugJSON(map[string]any{"events": events})
}

func readRemoteChipSinkEvents(ctx context.Context, since time.Time, limit int) error {
	return withResolvedRemoteRuntime(ctx, func(ctx context.Context, runtime *remoteclient.Runtime) error {
		resp, err := remoteclient.GetRemoteChipTestSinkEvents(ctx, runtime, since, limit)
		if err != nil {
			return err
		}
		return printDebugJSON(resp)
	})
}

func localChipSinkEventData(event *pb.CloudEvent) any {
	return chipsink.EventData(event)
}

func normalizeLocalSinkListenAddress(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return defaultLocalSinkListen, nil
	}
	// Accept bare port for convenience, e.g. "50052".
	if _, err := strconv.Atoi(trimmed); err == nil {
		return net.JoinHostPort("127.0.0.1", trimmed), nil
	}
	// Accept ":50052" and normalize to explicit host.
	if strings.HasPrefix(trimmed, ":") {
		return net.JoinHostPort("127.0.0.1", strings.TrimPrefix(trimmed, ":")), nil
	}
	_, port, err := net.SplitHostPort(trimmed)
	if err != nil || strings.TrimSpace(port) == "" {
		return "", fmt.Errorf("invalid --grpc-listen %q: expected host:port or port", raw)
	}
	return trimmed, nil
}

func waitForLocalSinkReady(pid int, listenAddr string, timeout time.Duration, logPath string) error {
	probeAddr, err := probeAddressForListen(listenAddr)
	if err != nil {
		return err
	}
	deadline := time.Now().Add(timeout)
	var lastDialErr error
	for time.Now().Before(deadline) {
		if !processExists(pid) {
			return fmt.Errorf("local chip sink process exited before becoming ready (pid=%d); check log: %s", pid, logPath)
		}
		conn, dialErr := net.DialTimeout("tcp", probeAddr, 200*time.Millisecond)
		if dialErr == nil {
			_ = conn.Close()
			return nil
		}
		lastDialErr = dialErr
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("local chip sink failed readiness probe on %s within %s (pid=%d, last error: %v); check log: %s", probeAddr, timeout, pid, lastDialErr, logPath)
}

func probeAddressForListen(listenAddr string) (string, error) {
	host, port, err := net.SplitHostPort(strings.TrimSpace(listenAddr))
	if err != nil {
		return "", fmt.Errorf("invalid normalized listen address %q: %w", listenAddr, err)
	}
	host = strings.TrimSpace(host)
	switch host {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port), nil
}

func stopPID(pid int) error {
	if pid <= 0 {
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	_ = proc.Signal(syscall.SIGTERM)
	deadline := time.Now().Add(2 * time.Second)
	for processExists(pid) && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}
	if processExists(pid) {
		_ = proc.Signal(syscall.SIGKILL)
	}
	return nil
}
