package environment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

const (
	relaySupervisorStateFilename = "relay_supervisor.toml"
	relaySupervisorLogFilename   = "relay_supervisor.log"
	relaySupervisorLockFilename  = "relay_supervisor.lock"
	defaultEC2AgentPort          = 18080
	defaultRelayWorkerPoolSize   = 16

	envRelaySupervisorLockPath = "CRE_RELAY_SUPERVISOR_LOCK_PATH"
)

var relaySupervisorLockFile *os.File

type relaySpec struct {
	Name string
	Port int
}

type relaySupervisorState struct {
	Version   int    `toml:"version"`
	PID       int    `toml:"pid"`
	Ports     []int  `toml:"ports"`
	StartedAt string `toml:"started_at,omitempty"`
	LogPath   string `toml:"log_path,omitempty"`
}

type localComponentRelayManager struct {
	lggr    zerolog.Logger
	baseURL string

	mu      sync.Mutex
	handles map[string]*relayHandle
}

type relayHandle struct {
	mu      sync.RWMutex
	relayID string
	name    string
	port    int
	cancel  context.CancelFunc
}

type relayOpenResponse struct {
	RelayID string `json:"relayId"`
}

type localBridgeStats struct {
	WSMessages     uint64
	WSToTCPBytes   uint64
	TCPToWSBytes   uint64
	LocalDialed    bool
	LocalDialFails uint64
}

func relaySupervisorCmd() *cobra.Command {
	var portsRaw string
	var relaySpecsRaw string
	cmd := &cobra.Command{
		Use:    "relay-supervisor",
		Short:  "Run detached mixed-mode relay supervisor",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			lockPath, err := resolveRelaySupervisorLockPath()
			if err != nil {
				return err
			}
			if err := acquireRelaySupervisorLock(lockPath); err != nil {
				return err
			}
			defer releaseRelaySupervisorLock()

			specs, err := parseRelaySpecsCSV(relaySpecsRaw)
			if err != nil {
				return err
			}
			if len(specs) == 0 {
				ports, perr := parsePortsCSV(portsRaw)
				if perr != nil {
					return perr
				}
				for _, p := range ports {
					specs = append(specs, relaySpec{
						Name: fmt.Sprintf("component-%d", p),
						Port: p,
					})
				}
			}
			if len(specs) == 0 {
				return fmt.Errorf("no relay specs or ports were provided")
			}

			manager, err := newLocalComponentRelayManager(framework.L)
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			for _, spec := range specs {
				if err := manager.EnsurePort(ctx, spec.Name, spec.Port); err != nil {
					_ = manager.Close(ctx)
					return err
				}
			}

			sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()
			<-sigCtx.Done()

			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return manager.Close(closeCtx)
		},
	}
	cmd.Flags().StringVar(&portsRaw, "ports", "", "Comma-separated list of local ports to bridge")
	cmd.Flags().StringVar(&relaySpecsRaw, "relay-specs", "", "Comma-separated list of relay specs in form name:port")
	return cmd
}

func maybeStartRelaySupervisor(relativePathToRepoRoot string, cfg *envconfig.Config) (bool, error) {
	specs := relaySpecsFromConfig(cfg)
	if len(specs) == 0 {
		if err := stopRelaySupervisor(relativePathToRepoRoot); err != nil {
			framework.L.Warn().Err(err).Msg("failed to stop stale relay supervisor")
		}
		return false, nil
	}
	framework.L.Info().Int("relaySpecs", len(specs)).Msgf("starting persistent relay supervisor with specs: %s", relaySpecsCSV(specs))
	return true, startRelaySupervisor(relativePathToRepoRoot, specs)
}

func relaySpecsFromConfig(cfg *envconfig.Config) []relaySpec {
	if cfg == nil {
		return nil
	}
	hasRemoteNodeSets := false
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) == string(envconfig.PlacementRemote) {
			hasRemoteNodeSets = true
			break
		}
	}
	if !hasRemoteNodeSets {
		return nil
	}

	specByPort := map[int]relaySpec{}
	addSpec := func(name string, port int) {
		if port <= 0 || port > 65535 {
			return
		}
		if _, exists := specByPort[port]; exists {
			return
		}
		specByPort[port] = relaySpec{Name: name, Port: port}
	}
	for _, blockchainCfg := range cfg.Blockchains {
		if blockchainCfg == nil || blockchainCfg.Placement != envconfig.PlacementLocal {
			continue
		}
		if blockchainCfg.Out != nil {
			for nodeIdx, node := range blockchainCfg.Out.Nodes {
				if node == nil {
					continue
				}
				if p, ok := endpointPort(node.ExternalHTTPUrl); ok {
					addSpec(fmt.Sprintf("blockchain-%s-http-%d", blockchainCfg.ChainID, nodeIdx), p)
				}
				if p, ok := endpointPort(node.ExternalWSUrl); ok {
					addSpec(fmt.Sprintf("blockchain-%s-ws-%d", blockchainCfg.ChainID, nodeIdx), p)
				}
			}
			continue
		}
		for _, p := range inferLocalBlockchainPortsFromInput(blockchainCfg.Input) {
			addSpec(fmt.Sprintf("blockchain-%s-port-%d", blockchainCfg.ChainID, p), p)
		}
	}

	if cfg.JD != nil && cfg.JD.Placement == envconfig.PlacementLocal {
		if cfg.JD.Out != nil {
			if p, ok := endpointPort(cfg.JD.Out.ExternalGRPCUrl); ok {
				addSpec("jd-grpc", p)
			}
			if p, ok := endpointPort(cfg.JD.Out.ExternalWSRPCUrl); ok {
				addSpec("jd-wsrpc", p)
			}
		} else {
			ports := inferLocalJDPortsFromInput(cfg.JD.Input)
			for idx, p := range ports {
				if idx == 0 {
					addSpec("jd-grpc", p)
					continue
				}
				if idx == 1 {
					addSpec("jd-wsrpc", p)
					continue
				}
				addSpec(fmt.Sprintf("jd-port-%d", p), p)
			}
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet == nil || strings.TrimSpace(nodeSet.Placement) != string(envconfig.PlacementLocal) {
			continue
		}
		for idx, p := range inferLocalNodeSetOCR2Ports(nodeSet) {
			addSpec(fmt.Sprintf("%s-ocr-%d", strings.TrimSpace(nodeSet.Name), idx), p)
		}
	}

	specs := make([]relaySpec, 0, len(specByPort))
	for _, spec := range specByPort {
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i, j int) bool {
		if specs[i].Port == specs[j].Port {
			return specs[i].Name < specs[j].Name
		}
		return specs[i].Port < specs[j].Port
	})
	return specs
}

func inferLocalBlockchainPortsFromInput(in blockchain.Input) []int {
	portSet := map[int]struct{}{}
	add := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}
		p, err := strconv.Atoi(raw)
		if err == nil && p > 0 && p <= 65535 {
			portSet[p] = struct{}{}
		}
	}
	chainType := strings.ToLower(strings.TrimSpace(in.Type))
	switch chainType {
	case "anvil", "":
		add(firstNonEmpty(in.Port, "8545"))
		// Anvil WS is served on the same port.
		add(firstNonEmpty(in.WSPort, in.Port, "8545"))
	default:
		// Best effort for other families: infer from explicit configured ports only.
		add(in.Port)
		add(in.WSPort)
	}
	out := make([]int, 0, len(portSet))
	for p := range portSet {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

func inferLocalJDPortsFromInput(in jd.Input) []int {
	const (
		defaultJDGRPC  = "14231"
		defaultJDWSRPC = "8080"
	)
	portSet := map[int]struct{}{}
	add := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}
		p, err := strconv.Atoi(raw)
		if err == nil && p > 0 && p <= 65535 {
			portSet[p] = struct{}{}
		}
	}
	add(firstNonEmpty(in.GRPCPort, defaultJDGRPC))
	add(firstNonEmpty(in.WSRPCPort, defaultJDWSRPC))
	out := make([]int, 0, len(portSet))
	for p := range portSet {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

func hasBootstrapRole(roles []string) bool {
	for _, role := range roles {
		if strings.EqualFold(strings.TrimSpace(role), "bootstrap") {
			return true
		}
	}
	return false
}

func inferLocalNodeSetOCR2Ports(nodeSet *cre.NodeSet) []int {
	if nodeSet == nil {
		return nil
	}
	nodeCount := nodeSet.Nodes
	if nodeCount <= 0 {
		nodeCount = len(nodeSet.NodeSpecs)
	}
	if nodeCount <= 0 {
		return nil
	}
	base := nodeSet.OCR2P2PRangeStart
	if base == 0 {
		httpStart := nodeSet.HTTPPortRangeStart
		if httpStart == 0 {
			httpStart = ns.DefaultHTTPPortStaticRangeStart
		}
		base = httpStart + (ns.DefaultOCR2P2PStaticRangeStart - ns.DefaultHTTPPortStaticRangeStart)
	}
	out := make([]int, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		p := base + i
		if p <= 0 || p > 65535 {
			continue
		}
		out = append(out, p)
	}
	return out
}

func endpointPort(raw string) (int, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil || parsed.Port() == "" {
			return 0, false
		}
		port, convErr := strconv.Atoi(parsed.Port())
		if convErr != nil || port <= 0 || port > 65535 {
			return 0, false
		}
		return port, true
	}
	_, portRaw, err := net.SplitHostPort(trimmed)
	if err != nil {
		return 0, false
	}
	port, convErr := strconv.Atoi(portRaw)
	if convErr != nil || port <= 0 || port > 65535 {
		return 0, false
	}
	return port, true
}

func startRelaySupervisor(relativePathToRepoRoot string, specs []relaySpec) error {
	if len(specs) == 0 {
		return nil
	}
	ports := make([]int, 0, len(specs))
	for _, spec := range specs {
		ports = append(ports, spec.Port)
	}
	ports = uniqueSortedPorts(ports)
	if err := stopRelaySupervisor(relativePathToRepoRoot); err != nil {
		framework.L.Warn().Err(err).Msg("failed to stop existing relay supervisor before restart")
	}

	executablePath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "resolve executable path for relay supervisor")
	}

	statePath := relaySupervisorStatePath(relativePathToRepoRoot)
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return errors.Wrap(err, "create relay supervisor state directory")
	}
	logPath := filepath.Join(filepath.Dir(statePath), relaySupervisorLogFilename)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return errors.Wrap(err, "open relay supervisor log file")
	}
	defer logFile.Close()

	cmd := exec.Command(executablePath, "env", "relay-supervisor", "--relay-specs", relaySpecsCSV(specs))
	lockPath := filepath.Join(filepath.Dir(statePath), relaySupervisorLockFilename)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", envRelaySupervisorLockPath, lockPath))
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "start relay supervisor process")
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	if !waitForPIDAlive(pid, 1500*time.Millisecond) {
		return fmt.Errorf("relay supervisor process exited too quickly (pid=%d)", pid)
	}

	state := relaySupervisorState{
		Version:   1,
		PID:       pid,
		Ports:     ports,
		StartedAt: time.Now().UTC().Format(time.RFC3339Nano),
		LogPath:   logPath,
	}
	return storeRelaySupervisorState(relativePathToRepoRoot, &state)
}

func stopRelaySupervisor(relativePathToRepoRoot string) error {
	state, err := loadRelaySupervisorState(relativePathToRepoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if state.PID <= 0 {
		return removeRelaySupervisorState(relativePathToRepoRoot)
	}
	if !processExists(state.PID) {
		return removeRelaySupervisorState(relativePathToRepoRoot)
	}
	isRelayProc, verifyErr := isRelaySupervisorProcess(state.PID)
	if verifyErr != nil {
		return verifyErr
	}
	if !isRelayProc {
		return fmt.Errorf("refusing to kill non-relay process pid=%d from relay supervisor state", state.PID)
	}
	proc, findErr := os.FindProcess(state.PID)
	if findErr != nil {
		return findErr
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
		return fmt.Errorf("relay supervisor pid=%d did not stop", state.PID)
	}
	return removeRelaySupervisorState(relativePathToRepoRoot)
}

func loadRelaySupervisorState(relativePathToRepoRoot string) (*relaySupervisorState, error) {
	data, err := os.ReadFile(relaySupervisorStatePath(relativePathToRepoRoot))
	if err != nil {
		return nil, err
	}
	state := &relaySupervisorState{}
	if err := toml.Unmarshal(data, state); err != nil {
		return nil, err
	}
	if state.Version == 0 {
		state.Version = 1
	}
	return state, nil
}

func storeRelaySupervisorState(relativePathToRepoRoot string, state *relaySupervisorState) error {
	data, err := toml.Marshal(state)
	if err != nil {
		return err
	}
	path := relaySupervisorStatePath(relativePathToRepoRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func removeRelaySupervisorState(relativePathToRepoRoot string) error {
	path := relaySupervisorStatePath(relativePathToRepoRoot)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func relaySupervisorStatePath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, remoteStateDirname, relaySupervisorStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for relay supervisor state file: %w", err))
	}
	return absPath
}

func resolveRelaySupervisorLockPath() (string, error) {
	if configured := strings.TrimSpace(os.Getenv(envRelaySupervisorLockPath)); configured != "" {
		return configured, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "resolve working directory for relay supervisor lock")
	}
	return filepath.Join(wd, remoteStateDirname, relaySupervisorLockFilename), nil
}

func acquireRelaySupervisorLock(lockPath string) error {
	if relaySupervisorLockFile != nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return errors.Wrap(err, "create relay supervisor lock directory")
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return errors.Wrap(err, "open relay supervisor lock file")
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return fmt.Errorf("relay supervisor already running (lock file in use: %s)", lockPath)
		}
		return errors.Wrap(err, "acquire relay supervisor file lock")
	}
	if err := f.Truncate(0); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return errors.Wrap(err, "truncate relay supervisor lock file")
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return errors.Wrap(err, "seek relay supervisor lock file")
	}
	_, _ = f.WriteString(fmt.Sprintf("pid=%d\nstarted_at=%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339Nano)))
	_ = f.Sync()
	relaySupervisorLockFile = f
	return nil
}

func releaseRelaySupervisorLock() {
	if relaySupervisorLockFile == nil {
		return
	}
	_ = syscall.Flock(int(relaySupervisorLockFile.Fd()), syscall.LOCK_UN)
	_ = relaySupervisorLockFile.Close()
	relaySupervisorLockFile = nil
}

func isRelaySupervisorProcess(pid int) (bool, error) {
	out, err := exec.Command("ps", "-o", "command=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return false, err
	}
	cmd := strings.TrimSpace(string(out))
	if cmd == "" {
		return false, nil
	}
	return strings.Contains(cmd, "relay-supervisor"), nil
}

func waitForPIDAlive(pid int, maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	for time.Now().Before(deadline) {
		if processExists(pid) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return processExists(pid)
}

func portsCSV(ports []int) string {
	if len(ports) == 0 {
		return ""
	}
	parts := make([]string, 0, len(ports))
	for _, port := range ports {
		parts = append(parts, strconv.Itoa(port))
	}
	return strings.Join(parts, ",")
}

func parsePortsCSV(raw string) ([]int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		port, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", part, err)
		}
		if port <= 0 || port > 65535 {
			return nil, fmt.Errorf("invalid port %d", port)
		}
		out = append(out, port)
	}
	return uniqueSortedPorts(out), nil
}

func parseRelaySpecsCSV(raw string) ([]relaySpec, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	specByPort := make(map[int]relaySpec, len(parts))
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}
		idx := strings.LastIndex(token, ":")
		if idx <= 0 || idx >= len(token)-1 {
			return nil, fmt.Errorf("invalid relay spec %q; expected name:port", token)
		}
		name := strings.TrimSpace(token[:idx])
		portRaw := strings.TrimSpace(token[idx+1:])
		if name == "" {
			return nil, fmt.Errorf("invalid relay spec %q; name is empty", token)
		}
		port, err := strconv.Atoi(portRaw)
		if err != nil || port <= 0 || port > 65535 {
			return nil, fmt.Errorf("invalid relay port %q in spec %q", portRaw, token)
		}
		if _, exists := specByPort[port]; exists {
			continue
		}
		specByPort[port] = relaySpec{Name: name, Port: port}
	}
	specs := make([]relaySpec, 0, len(specByPort))
	for _, spec := range specByPort {
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i, j int) bool {
		if specs[i].Port == specs[j].Port {
			return specs[i].Name < specs[j].Name
		}
		return specs[i].Port < specs[j].Port
	})
	return specs, nil
}

func relaySpecsCSV(specs []relaySpec) string {
	if len(specs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(specs))
	for _, spec := range specs {
		if spec.Port <= 0 || spec.Port > 65535 {
			continue
		}
		name := strings.TrimSpace(spec.Name)
		if name == "" {
			name = fmt.Sprintf("component-%d", spec.Port)
		}
		parts = append(parts, fmt.Sprintf("%s:%d", name, spec.Port))
	}
	return strings.Join(parts, ",")
}

func uniqueSortedPorts(in []int) []int {
	if len(in) == 0 {
		return nil
	}
	set := make(map[int]struct{}, len(in))
	for _, p := range in {
		if p > 0 && p <= 65535 {
			set[p] = struct{}{}
		}
	}
	out := make([]int, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

func newLocalComponentRelayManager(lggr zerolog.Logger) (*localComponentRelayManager, error) {
	baseURL, err := resolveAgentBaseURLForRelay()
	if err != nil {
		return nil, err
	}
	return &localComponentRelayManager{
		lggr:    lggr,
		baseURL: baseURL,
		handles: make(map[string]*relayHandle),
	}, nil
}

func (m *localComponentRelayManager) EnsurePort(ctx context.Context, relayName string, localPort int) error {
	if m == nil || localPort <= 0 {
		return nil
	}
	// Deduplicate by port. HTTP and WS for the same endpoint can share one listener.
	key := strconv.Itoa(localPort)

	m.mu.Lock()
	if _, ok := m.handles[key]; ok {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	relayID, err := openRelay(ctx, m.baseURL, relayName, localPort)
	if err != nil {
		return err
	}

	workerCtx, cancel := context.WithCancel(context.Background())
	localAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(localPort))
	handle := &relayHandle{
		relayID: relayID,
		name:    relayName,
		port:    localPort,
		cancel:  cancel,
	}
	for i := 0; i < defaultRelayWorkerPoolSize; i++ {
		go relayWorker(workerCtx, m.lggr, m.baseURL, handle, localAddr, i)
	}

	m.mu.Lock()
	m.handles[key] = handle
	m.mu.Unlock()
	m.lggr.Info().Str("relayName", relayName).Int("port", localPort).Msg("ensured persistent mixed component relay")
	return nil
}

func (m *localComponentRelayManager) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	handles := make([]*relayHandle, 0, len(m.handles))
	for _, h := range m.handles {
		handles = append(handles, h)
	}
	m.handles = map[string]*relayHandle{}
	m.mu.Unlock()

	var firstErr error
	for _, h := range handles {
		h.cancel()
		if err := closeRelay(ctx, m.baseURL, h.getRelayID()); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (h *relayHandle) getRelayID() string {
	if h == nil {
		return ""
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.relayID
}

func (h *relayHandle) setRelayID(relayID string) {
	if h == nil {
		return
	}
	h.mu.Lock()
	h.relayID = relayID
	h.mu.Unlock()
}

func resolveAgentBaseURLForRelay() (string, error) {
	if v := strings.TrimSpace(os.Getenv("CRE_EC2_AGENT_URL")); v != "" {
		return v, nil
	}
	hostIP, err := runtimecfg.DirectHostIP()
	if err == nil {
		port, err := resolveEC2AgentPortForRelay()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("http://%s:%d", hostIP, port), nil
	}
	return "", fmt.Errorf("cannot resolve agent base URL for relay; set CRE_EC2_AGENT_URL or provide EC2 discovery envs: %w", err)
}

func resolveEC2AgentPortForRelay() (int, error) {
	raw := strings.TrimSpace(os.Getenv("CRE_EC2_AGENT_PORT"))
	if raw == "" {
		return defaultEC2AgentPort, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port <= 0 || port > 65535 {
		return 0, fmt.Errorf("invalid CRE_EC2_AGENT_PORT: %q", raw)
	}
	return port, nil
}

func openRelay(ctx context.Context, baseURL, name string, requestedPort int) (string, error) {
	body, _ := json.Marshal(map[string]any{"name": name, "requestedPort": requestedPort})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/v1/relay/open", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("open relay failed: status %s body %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var out relayOpenResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", err
	}
	if strings.TrimSpace(out.RelayID) == "" {
		return "", fmt.Errorf("open relay returned empty relayId")
	}
	return out.RelayID, nil
}

func closeRelay(ctx context.Context, baseURL, relayID string) error {
	body, _ := json.Marshal(map[string]any{"relayId": relayID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/v1/relay/close", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("close relay failed: status %s body %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	return nil
}

func relayWorker(ctx context.Context, lggr zerolog.Logger, baseURL string, handle *relayHandle, localAddr string, workerIndex int) {
	backoff := 250 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		relayID := handle.getRelayID()
		wsURL, err := relayConnectWSURL(baseURL, relayID)
		if err != nil {
			lggr.Warn().
				Err(err).
				Str("relayId", relayID).
				Str("relayName", handle.name).
				Int("workerIndex", workerIndex).
				Msg("relay worker failed to construct websocket URL")
			time.Sleep(backoff)
			continue
		}
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			if isBadHandshakeError(err) {
				reopenCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				newRelayID, reopenErr := openRelay(reopenCtx, baseURL, handle.name, handle.port)
				cancel()
				if reopenErr != nil {
					lggr.Warn().
						Err(reopenErr).
						Str("relayId", relayID).
						Str("relayName", handle.name).
						Int("requestedPort", handle.port).
						Int("workerIndex", workerIndex).
						Msg("relay worker failed to reopen relay after websocket bad handshake")
				} else {
					handle.setRelayID(newRelayID)
					lggr.Info().
						Str("oldRelayId", relayID).
						Str("newRelayId", newRelayID).
						Str("relayName", handle.name).
						Int("requestedPort", handle.port).
						Int("workerIndex", workerIndex).
						Msg("relay worker refreshed relay id after websocket bad handshake")
					backoff = 250 * time.Millisecond
					continue
				}
			}
			lggr.Warn().
				Err(err).
				Str("relayId", relayID).
				Str("relayName", handle.name).
				Int("workerIndex", workerIndex).
				Msg("relay worker failed to connect websocket bridge")
			time.Sleep(backoff)
			continue
		}
		lggr.Info().
			Str("relayId", relayID).
			Str("relayName", handle.name).
			Str("localAddr", localAddr).
			Int("workerIndex", workerIndex).
			Msg("relay worker established websocket bridge; waiting for payload to dial local endpoint")
		bridgeStarted := time.Now()
		stats, bridgeErr := bridgeRelayStream(ctx, lggr, handle.name, relayID, workerIndex, ws, localAddr)
		_ = ws.Close()
		if bridgeErr != nil && !errors.Is(bridgeErr, context.Canceled) {
			lggr.Warn().
				Err(bridgeErr).
				Str("relayId", relayID).
				Str("relayName", handle.name).
				Int("workerIndex", workerIndex).
				Uint64("wsMessages", stats.WSMessages).
				Uint64("wsToTCPBytes", stats.WSToTCPBytes).
				Uint64("tcpToWSBytes", stats.TCPToWSBytes).
				Bool("localDialed", stats.LocalDialed).
				Uint64("localDialFails", stats.LocalDialFails).
				Dur("duration", time.Since(bridgeStarted)).
				Msg("relay worker bridge ended with error")
		} else {
			lggr.Info().
				Str("relayId", relayID).
				Str("relayName", handle.name).
				Int("workerIndex", workerIndex).
				Uint64("wsMessages", stats.WSMessages).
				Uint64("wsToTCPBytes", stats.WSToTCPBytes).
				Uint64("tcpToWSBytes", stats.TCPToWSBytes).
				Bool("localDialed", stats.LocalDialed).
				Uint64("localDialFails", stats.LocalDialFails).
				Dur("duration", time.Since(bridgeStarted)).
				Msg("relay worker bridge ended")
		}
		if backoff < 2*time.Second {
			backoff *= 2
		}
	}
}

func isBadHandshakeError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "bad handshake")
}

func relayConnectWSURL(baseURL, relayID string) (string, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported agent url scheme: %s", u.Scheme)
	}
	u.Path = "/v1/relay/connect"
	q := u.Query()
	q.Set("relayId", relayID)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func bridgeRelayStream(
	ctx context.Context,
	lggr zerolog.Logger,
	relayName, relayID string,
	workerIndex int,
	ws *websocket.Conn,
	localAddr string,
) (*localBridgeStats, error) {
	errCh := make(chan error, 2)
	stats := &localBridgeStats{}
	writeMu := &sync.Mutex{}
	localReady := make(chan net.Conn, 1)
	var localConn net.Conn
	var localConnMu sync.Mutex
	keepAliveCtx, keepAliveCancel := context.WithCancel(ctx)
	defer keepAliveCancel()
	go relayKeepAlive(keepAliveCtx, ws, writeMu, errCh)
	getLocalConn := func() net.Conn {
		localConnMu.Lock()
		defer localConnMu.Unlock()
		return localConn
	}
	setLocalConn := func(conn net.Conn) {
		localConnMu.Lock()
		localConn = conn
		localConnMu.Unlock()
	}
	ensureLocalConn := func() (net.Conn, error) {
		if existing := getLocalConn(); existing != nil {
			return existing, nil
		}
		conn, err := net.DialTimeout("tcp", localAddr, 2*time.Second)
		if err != nil {
			atomic.AddUint64(&stats.LocalDialFails, 1)
			lggr.Warn().
				Err(err).
				Str("relayId", relayID).
				Str("relayName", relayName).
				Int("workerIndex", workerIndex).
				Str("localAddr", localAddr).
				Msg("relay worker lazy local dial failed")
			return nil, err
		}
		stats.LocalDialed = true
		lggr.Info().
			Str("relayId", relayID).
			Str("relayName", relayName).
			Int("workerIndex", workerIndex).
			Str("localAddr", localAddr).
			Msg("relay worker lazy local dial succeeded")
		setLocalConn(conn)
		select {
		case localReady <- conn:
		default:
		}
		return conn, nil
	}
	defer func() {
		if conn := getLocalConn(); conn != nil {
			_ = conn.Close()
		}
	}()

	go func() {
		var conn net.Conn
		select {
		case conn = <-localReady:
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		}
		if conn == nil {
			errCh <- fmt.Errorf("local relay connection was nil")
			return
		}
		buf := make([]byte, 32*1024)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				atomic.AddUint64(&stats.TCPToWSBytes, uint64(n))
				writeMu.Lock()
				wErr := ws.WriteMessage(websocket.BinaryMessage, buf[:n])
				writeMu.Unlock()
				if wErr != nil {
					errCh <- wErr
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	go func() {
		for {
			msgType, payload, err := ws.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
				continue
			}
			if len(payload) == 0 {
				continue
			}
			atomic.AddUint64(&stats.WSMessages, 1)
			atomic.AddUint64(&stats.WSToTCPBytes, uint64(len(payload)))
			if stats.WSMessages == 1 {
				lggr.Info().
					Str("relayId", relayID).
					Str("relayName", relayName).
					Int("workerIndex", workerIndex).
					Int("payloadBytes", len(payload)).
					Msg("relay worker received first websocket payload")
			}
			conn, dialErr := ensureLocalConn()
			if dialErr != nil {
				errCh <- dialErr
				return
			}
			if _, wErr := conn.Write(payload); wErr != nil {
				errCh <- wErr
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		return stats, ctx.Err()
	case err := <-errCh:
		return stats, err
	}
}

func relayKeepAlive(ctx context.Context, ws *websocket.Conn, writeMu *sync.Mutex, errCh chan<- error) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			writeMu.Lock()
			err := ws.WriteControl(websocket.PingMessage, []byte("keepalive"), time.Now().Add(5*time.Second))
			writeMu.Unlock()
			if err != nil {
				select {
				case errCh <- fmt.Errorf("keepalive ping failed: %w", err):
				default:
				}
				return
			}
		}
	}
}
