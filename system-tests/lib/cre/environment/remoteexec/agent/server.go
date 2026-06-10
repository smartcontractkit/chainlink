package agent

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/mount"
	dockerclient "github.com/moby/moby/client"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/chipsink"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/internal/dockerops"
)

const (
	SchemaVersionV1           = "v1"
	OperationStartComponent   = "StartComponent"
	OperationStopComponent    = "StopComponent"
	OperationDeployArtifacts  = "DeployArtifacts"
	OperationHealth           = "Health"
	ComponentTypeBlockchain   = "blockchain"
	ComponentTypeJD           = "jd"
	ComponentTypeNodeSet      = "nodeset"
	ComponentTypeChipTestSink = "chip-testsink"

	ErrCodeMethodNotAllowed      = "method_not_allowed"
	ErrCodeInvalidRequestBody    = "invalid_request_body"
	ErrCodeUnsupportedSchema     = "unsupported_schema_version"
	ErrCodeUnsupportedOperation  = "unsupported_operation"
	ErrCodeInvalidPayload        = "invalid_payload"
	ErrCodeUnsupportedComponent  = "unsupported_component_type"
	ErrCodeMissingComponentInput = "missing_component_input"
	ErrCodeDeployFailed          = "deployment_failed"
	ErrCodeTransportEncodeFailed = "transport_encode_failed"

	RemoteStartPolicyAlways         = "always"
	RemoteStartPolicyReuseIdentical = "reuse_if_identical"

	EnvKeepFailedContainers = "CRE_AGENT_KEEP_FAILED_CONTAINERS"

	defaultComponentLogsLimit       = 200
	maxComponentLogsLimit           = 1000
	componentLogsRingSize           = 2000
	inFlightOperationScopeLifecycle = "lifecycle"
	inFlightOperationScopeGeneral   = "general"
	protocolVersion                 = "1.0.0"
	capabilityComponentLogs         = "componentLogs"
	capabilityLocks                 = "locks"
	capabilityDeployArtifacts       = "deployArtifacts"
	capabilityStartComponent        = "startComponent"
	capabilityRelay                 = "relay"
	capabilityListCTFResources      = "listCTFResources"
	agentVersion                    = "dev"
)

var frameworkLogCaptureMu sync.Mutex

type StartComponentEnvelope struct {
	SchemaVersion string          `json:"schemaVersion"`
	Operation     string          `json:"operation"`
	Payload       json.RawMessage `json:"payload"`
}

type StartComponentPayload struct {
	ComponentType      string            `json:"componentType"`
	Blockchain         *blockchain.Input `json:"blockchain"`
	RegistryBlockchain map[string]any    `json:"registryBlockchain,omitempty"`
	JD                 *jd.Input         `json:"jd"`
	NodeSet            *ns.Input         `json:"nodeset,omitempty"`
	ReusePolicy        string            `json:"reusePolicy,omitempty"`
}

type DeployArtifactsPayload struct {
	NodeSetName string                `json:"nodeSetName"`
	TargetDir   string                `json:"targetDir"`
	Files       []DeployArtifactsFile `json:"files"`
}

type DeployArtifactsFile struct {
	Name          string `json:"name"`
	ContentBase64 string `json:"contentBase64"`
}

type StartComponentResponse struct {
	ComponentType string         `json:"componentType,omitempty"`
	Output        map[string]any `json:"output,omitempty"`
	Found         bool           `json:"found,omitempty"`
	Stopped       bool           `json:"stopped,omitempty"`
	AgentLogs     []string       `json:"agentLogs,omitempty"`
	ErrorCode     string         `json:"errorCode,omitempty"`
	Error         string         `json:"error,omitempty"`
}

type CTFResourcesResponse struct {
	Containers []string `json:"containers,omitempty"`
	Volumes    []string `json:"volumes,omitempty"`
}

//nolint:revive // AgentStatusResponse is the API contract; renaming would break external callers
type AgentStatusResponse struct {
	AgentVersion      string                      `json:"agentVersion,omitempty"`
	ProtocolVersion   string                      `json:"protocolVersion,omitempty"`
	SupportedSchemas  []string                    `json:"supportedSchemas,omitempty"`
	Capabilities      []string                    `json:"capabilities,omitempty"`
	UptimeSeconds     int64                       `json:"uptimeSeconds"`
	RuntimeComponents []string                    `json:"runtimeComponents,omitempty"`
	CachedComponents  []string                    `json:"cachedComponents,omitempty"`
	Relays            []RelayInfo                 `json:"relays,omitempty"`
	ComponentLogKeys  []string                    `json:"componentLogKeys,omitempty"`
	InFlight          []InFlightOperation         `json:"inFlight,omitempty"`
	ChipSink          *ChipTestSinkStatusResponse `json:"chipSink,omitempty"`
}

type RelayInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	RequestedPort int    `json:"requestedPort"`
	BoundPort     int    `json:"boundPort"`
}

//nolint:revive // AgentLocksResponse is the API contract; renaming would break external callers
type AgentLocksResponse struct {
	LifecycleBusy    bool                `json:"lifecycleBusy"`
	CacheEntries     int                 `json:"cacheEntries"`
	RuntimeEntries   int                 `json:"runtimeEntries"`
	RelayCount       int                 `json:"relayCount"`
	ComponentLogKeys int                 `json:"componentLogKeys"`
	InFlight         []InFlightOperation `json:"inFlight,omitempty"`
}

type InFlightOperation struct {
	ID         string `json:"id"`
	Scope      string `json:"scope"`
	StartedAt  string `json:"startedAt"`
	DurationMs int64  `json:"durationMs"`
}

type ComponentLogsResponse struct {
	ComponentKey string   `json:"componentKey"`
	TotalLines   int      `json:"totalLines"`
	Lines        []string `json:"lines,omitempty"`
}

type ChipTestSinkStartRequest struct {
	Name             string `json:"name,omitempty"`
	GRPCListen       string `json:"grpcListen,omitempty"`
	UpstreamEndpoint string `json:"upstreamEndpoint,omitempty"`
}

type ChipTestSinkStartResponse struct {
	Profile          string `json:"profile"`
	Mode             string `json:"mode"`
	Name             string `json:"name"`
	GRPCListen       string `json:"grpcListen"`
	UpstreamEndpoint string `json:"upstreamEndpoint,omitempty"`
	EventLogPath     string `json:"eventLogPath,omitempty"`
}

type ChipTestSinkStatusResponse struct {
	Profile          string `json:"profile"`
	Mode             string `json:"mode"`
	Running          bool   `json:"running"`
	Name             string `json:"name,omitempty"`
	GRPCListen       string `json:"grpcListen,omitempty"`
	UpstreamEndpoint string `json:"upstreamEndpoint,omitempty"`
	EventLogPath     string `json:"eventLogPath,omitempty"`
}

type ChipTestSinkStopResponse struct {
	Found   bool `json:"found"`
	Stopped bool `json:"stopped"`
}

type ChipTestSinkEventLogEntry struct {
	Timestamp string         `json:"timestamp"`
	Type      string         `json:"type,omitempty"`
	Event     map[string]any `json:"event,omitempty"`
}

type ChipTestSinkEventsResponse struct {
	Events []ChipTestSinkEventLogEntry `json:"events"`
}

type inFlightOperation struct {
	ID        string
	Scope     string
	StartedAt time.Time
}

type Server struct {
	lggr          zerolog.Logger
	deployers     map[blockchain.ChainFamily]blockchains.Deployer
	startedAt     time.Time
	lifecycleMu   sync.Mutex
	cacheMu       sync.Mutex
	cache         map[string]cachedStart
	runtime       map[string]runtimeState
	relayMu       sync.Mutex
	relays        map[string]*relayRegistration
	logsMu        sync.Mutex
	componentLogs map[string][]string
	opsMu         sync.Mutex
	inFlight      map[string]inFlightOperation
	chipSinkMu    sync.Mutex
	chipSink      *chipTestSinkRuntime
}

type cachedStart struct {
	PayloadHash string
	Output      map[string]any
}

type runtimeState struct {
	ComponentType string
	ContainerIDs  []string
	StopFn        func(context.Context) error
}

type chipTestSinkRuntime struct {
	name             string
	grpcListen       string
	upstreamEndpoint string
	eventLogPath     string
	server           *chipsink.Server
	cancel           context.CancelFunc
	runErrCh         chan error
}

func NewServer(lggr zerolog.Logger, deployers map[blockchain.ChainFamily]blockchains.Deployer) *Server {
	return &Server{
		lggr:          lggr,
		deployers:     deployers,
		startedAt:     time.Now(),
		cache:         make(map[string]cachedStart),
		runtime:       make(map[string]runtimeState),
		relays:        make(map[string]*relayRegistration),
		componentLogs: make(map[string][]string),
		inFlight:      make(map[string]inFlightOperation),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health", s.health)
	mux.HandleFunc("/v1/components/start", s.startComponent)
	mux.HandleFunc("/v1/relay/open", s.openRelay)
	mux.HandleFunc("/v1/relay/close", s.closeRelay)
	mux.HandleFunc("/v1/relay/connect", s.connectRelay)
	mux.HandleFunc("/v1/resources/ctf", s.listCTFResources)
	mux.HandleFunc("/v1/status", s.status)
	mux.HandleFunc("/v1/locks", s.locks)
	mux.HandleFunc("/v1/components/logs", s.componentLogsHandler)
	mux.HandleFunc("/v1/chip/sink/start", s.startChipTestSink)
	mux.HandleFunc("/v1/chip/sink/stop", s.stopChipTestSink)
	mux.HandleFunc("/v1/chip/sink/status", s.chipTestSinkStatus)
	mux.HandleFunc("/v1/chip/sink/events", s.chipTestSinkEvents)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) listCTFResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	client, err := dockerclient.New(dockerclient.FromEnv)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to create docker client: %v", err), nil)
		return
	}
	defer client.Close()

	filterArgs := make(dockerclient.Filters).Add("label", "framework=ctf")
	listRes, err := client.ContainerList(r.Context(), dockerclient.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to list ctf containers: %v", err), nil)
		return
	}
	containerNames := make([]string, 0, len(listRes.Items))
	for _, c := range listRes.Items {
		if len(c.Names) > 0 {
			containerNames = append(containerNames, strings.TrimPrefix(c.Names[0], "/"))
			continue
		}
		containerNames = append(containerNames, c.ID)
	}
	slices.Sort(containerNames)

	volResp, err := client.VolumeList(r.Context(), dockerclient.VolumeListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to list ctf volumes: %v", err), nil)
		return
	}
	volumeNames := make([]string, 0, len(volResp.Items))
	for _, v := range volResp.Items {
		if strings.TrimSpace(v.Name) == "" {
			continue
		}
		volumeNames = append(volumeNames, v.Name)
	}
	slices.Sort(volumeNames)

	s.respondJSONAny(w, http.StatusOK, CTFResourcesResponse{
		Containers: containerNames,
		Volumes:    volumeNames,
	})
}

func (s *Server) startComponent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed", nil)
		return
	}

	var envelope StartComponentEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidRequestBody, fmt.Sprintf("invalid request body: %v", err), nil)
		return
	}

	if envelope.SchemaVersion != SchemaVersionV1 {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedSchema, "unsupported schema version: "+envelope.SchemaVersion, nil)
		return
	}
	if envelope.Operation == OperationDeployArtifacts {
		s.deployArtifacts(w, r, envelope.Payload)
		return
	}
	var payload StartComponentPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("invalid payload: %v", err), nil)
		return
	}
	if payload.ComponentType != ComponentTypeBlockchain && payload.ComponentType != ComponentTypeJD && payload.ComponentType != ComponentTypeNodeSet {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedComponent, "unsupported component type: "+payload.ComponentType, nil)
		return
	}

	componentKey, inputErr := componentCacheKey(payload)
	if inputErr != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, inputErr.Error(), nil)
		return
	}
	if envelope.Operation == OperationStopComponent {
		s.stopComponentByKey(w, r, payload.ComponentType, componentKey)
		return
	}
	if envelope.Operation != OperationStartComponent {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedOperation, "unsupported operation: "+envelope.Operation, nil)
		return
	}
	payloadHash := hashPayload(envelope.Payload)

	// Keep this stderr write explicit so startup behavior is visible when agent runs as a subprocess.
	requestLog := fmt.Sprintf("[cre-agent] starting component type=%s key=%s", payload.ComponentType, componentKey)
	_, _ = fmt.Fprintln(os.Stderr, requestLog)
	s.beginInFlight("start:"+componentKey, inFlightOperationScopeLifecycle)
	defer s.endInFlight("start:" + componentKey)
	preStartLogs := make([]string, 0, 2)
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	if shouldRestartBeforeStart(payload.ComponentType, payload.ReusePolicy) {
		stopped, stopErr := s.stopTrackedComponentLocked(r.Context(), componentKey)
		if stopErr != nil {
			s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to stop existing component before start: %v", stopErr), []string{requestLog})
			return
		}
		if stopped {
			preStartLogs = append(preStartLogs, "[cre-agent] stopped existing component before start")
		} else {
			preStartLogs = append(preStartLogs, "[cre-agent] no existing component to stop before start")
		}
	}
	if shouldReuseRemoteStart(payload.ComponentType, payload.ReusePolicy) {
		if cached, ok := s.lookupCachedStart(componentKey, payloadHash); ok {
			reuseLog := fmt.Sprintf("[cre-agent] reusing existing component for key=%s (payload hash matched)", componentKey)
			_, _ = fmt.Fprintln(os.Stderr, reuseLog)
			s.respondJSON(w, http.StatusOK, StartComponentResponse{
				ComponentType: payload.ComponentType,
				Output:        cached.Output,
				AgentLogs:     []string{requestLog, reuseLog},
			})
			s.appendComponentLogs(componentKey, []string{requestLog, reuseLog})
			return
		}
	}

	agentLogs := make([]string, 0, 8)
	agentLogs = append(agentLogs, requestLog)
	agentLogs = append(agentLogs, preStartLogs...)
	var blockchainOutput *blockchain.Output
	var jdOutput *jd.Output
	var nodeSetOutput *ns.Output
	trackedContainers, startErr := s.discoverOwnedContainers(r.Context(), func() error {
		capturedFrameworkLogs, runErr := captureFrameworkLogs(func() error {
			switch payload.ComponentType {
			case ComponentTypeBlockchain:
				deployed, err := DeployBlockchainComponent(r.Context(), s.deployers, payload.Blockchain)
				if err != nil {
					return err
				}
				blockchainOutput = deployed
			case ComponentTypeJD:
				deployed, err := DeployJDComponent(r.Context(), payload.JD)
				if err != nil {
					return err
				}
				jdOutput = deployed
			case ComponentTypeNodeSet:
				registryOutput, err := DecodeFromTransport[blockchain.Output](payload.RegistryBlockchain)
				if err != nil {
					return fmt.Errorf("failed to decode registry blockchain payload for nodeset: %w", err)
				}
				deployed, err := DeployNodeSetComponent(r.Context(), payload.NodeSet, registryOutput)
				if err != nil {
					return err
				}
				nodeSetOutput = deployed
			}
			return nil
		})
		agentLogs = append(agentLogs, capturedFrameworkLogs...)
		return runErr
	})

	if startErr != nil {
		if len(trackedContainers) > 0 && shouldCleanupFailedContainers() {
			cleanupErr := stopContainers(r.Context(), trackedContainers)
			if cleanupErr != nil {
				agentLogs = append(agentLogs, fmt.Sprintf("[cre-agent] failed startup cleanup for %d tracked container(s): %v", len(trackedContainers), cleanupErr))
			} else {
				agentLogs = append(agentLogs, fmt.Sprintf("[cre-agent] cleaned up %d tracked container(s) after failed startup", len(trackedContainers)))
			}
		} else if len(trackedContainers) > 0 {
			agentLogs = append(agentLogs, fmt.Sprintf("[cre-agent] preserving %d tracked container(s) after failed startup because %s is enabled", len(trackedContainers), EnvKeepFailedContainers))
		}
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, startErr.Error(), agentLogs)
		s.appendComponentLogs(componentKey, agentLogs)
		return
	}

	var output map[string]any
	var encErr error
	switch {
	case blockchainOutput != nil:
		output, encErr = EncodeForTransport(blockchainOutput)
	case jdOutput != nil:
		output, encErr = EncodeForTransport(jdOutput)
	case nodeSetOutput != nil:
		output, encErr = EncodeForTransport(nodeSetOutput)
	}
	if encErr != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeTransportEncodeFailed, encErr.Error(), agentLogs)
		s.appendComponentLogs(componentKey, agentLogs)
		return
	}
	if shouldReuseRemoteStart(payload.ComponentType, payload.ReusePolicy) {
		s.cacheSuccessfulStart(componentKey, payloadHash, output)
	}
	s.storeRuntime(componentKey, runtimeState{
		ComponentType: payload.ComponentType,
		ContainerIDs:  trackedContainers,
	})
	s.respondJSON(w, http.StatusOK, StartComponentResponse{
		ComponentType: payload.ComponentType,
		Output:        output,
		AgentLogs:     agentLogs,
	})
	s.appendComponentLogs(componentKey, agentLogs)
}

func (s *Server) deployArtifacts(w http.ResponseWriter, r *http.Request, rawPayload json.RawMessage) {
	s.beginInFlight("deploy-artifacts", inFlightOperationScopeGeneral)
	defer s.endInFlight("deploy-artifacts")

	var payload DeployArtifactsPayload
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("invalid payload: %v", err), nil)
		return
	}
	if strings.TrimSpace(payload.NodeSetName) == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "nodeset name is required", nil)
		return
	}
	if strings.TrimSpace(payload.TargetDir) == "" {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "target dir is required", nil)
		return
	}
	if len(payload.Files) == 0 {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "at least one artifact file is required", nil)
		return
	}

	tmpDir, err := os.MkdirTemp("", "cre-agent-artifacts")
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to create temp dir: %v", err), nil)
		return
	}
	defer os.RemoveAll(tmpDir)

	filePaths := make([]string, 0, len(payload.Files))
	for idx, f := range payload.Files {
		if strings.TrimSpace(f.Name) == "" {
			s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("artifact %d has empty name", idx), nil)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(f.ContentBase64)
		if err != nil {
			s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("artifact %s has invalid base64 content: %v", f.Name, err), nil)
			return
		}
		target := filepath.Join(tmpDir, filepath.Base(f.Name))
		if err := os.WriteFile(target, decoded, 0o600); err != nil {
			s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to write artifact %s: %v", f.Name, err), nil)
			return
		}
		filePaths = append(filePaths, target)
	}

	containerPrefix := ns.NodeNamePrefix(payload.NodeSetName)
	if err := dockerops.CopyFilesToContainers(r.Context(), containerPrefix, payload.TargetDir, filePaths); err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to copy artifacts to containers: %v", err), nil)
		return
	}

	s.respondJSON(w, http.StatusOK, StartComponentResponse{
		AgentLogs: []string{
			fmt.Sprintf("[cre-agent] copied %d artifact(s) to all containers for nodeset %s", len(filePaths), payload.NodeSetName),
		},
	})
	s.appendComponentLogs(fmt.Sprintf("%s:%s", ComponentTypeNodeSet, payload.NodeSetName), []string{
		fmt.Sprintf("[cre-agent] copied %d artifact(s) to all containers for nodeset %s", len(filePaths), payload.NodeSetName),
	})
}

func (s *Server) stopComponentByKey(w http.ResponseWriter, r *http.Request, componentType, componentKey string) {
	s.beginInFlight("stop:"+componentKey, inFlightOperationScopeLifecycle)
	defer s.endInFlight("stop:" + componentKey)

	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	requestLog := fmt.Sprintf("[cre-agent] stopping component type=%s key=%s", componentType, componentKey)
	_, _ = fmt.Fprintln(os.Stderr, requestLog)

	stopped, err := s.stopTrackedComponentLocked(r.Context(), componentKey)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, fmt.Sprintf("failed to stop component containers: %v", err), []string{requestLog})
		return
	}
	if !stopped {
		s.deleteCachedStart(componentKey)
		s.respondJSON(w, http.StatusOK, StartComponentResponse{
			ComponentType: componentType,
			Found:         false,
			Stopped:       false,
			AgentLogs:     []string{requestLog, "[cre-agent] nothing to stop (component not found)"},
		})
		s.appendComponentLogs(componentKey, []string{requestLog, "[cre-agent] nothing to stop (component not found)"})
		return
	}
	s.deleteCachedStart(componentKey)
	s.respondJSON(w, http.StatusOK, StartComponentResponse{
		ComponentType: componentType,
		Found:         true,
		Stopped:       true,
		AgentLogs:     []string{requestLog, "[cre-agent] stopped existing component"},
	})
	s.appendComponentLogs(componentKey, []string{requestLog, "[cre-agent] stopped existing component"})
}

func (s *Server) respondJSON(w http.ResponseWriter, code int, body StartComponentResponse) {
	s.respondJSONAny(w, code, body)
}

func (s *Server) respondJSONAny(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func (s *Server) respondError(w http.ResponseWriter, code int, errorCode string, message string, logs []string) {
	s.respondJSON(w, code, StartComponentResponse{
		AgentLogs: logs,
		ErrorCode: errorCode,
		Error:     message,
	})
}

func captureFrameworkLogs(fn func() error) ([]string, error) {
	frameworkLogCaptureMu.Lock()
	defer frameworkLogCaptureMu.Unlock()

	var buf bytes.Buffer
	originalLogger := framework.L
	framework.L = originalLogger.Output(io.MultiWriter(os.Stderr, &buf))
	defer func() {
		framework.L = originalLogger
	}()

	err := fn()

	logs := make([]string, 0)
	for _, line := range strings.Split(buf.String(), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		logs = append(logs, trimmed)
	}

	return logs, err
}

func shouldReuseRemoteStart(componentType, policy string) bool {
	if componentType == ComponentTypeJD {
		return false
	}
	normalized := strings.TrimSpace(strings.ToLower(policy))
	if normalized == "" {
		normalized = RemoteStartPolicyReuseIdentical
	}
	return normalized == RemoteStartPolicyReuseIdentical
}

func shouldRestartBeforeStart(componentType, policy string) bool {
	if componentType == ComponentTypeJD {
		return true
	}
	normalized := strings.TrimSpace(strings.ToLower(policy))
	return normalized == RemoteStartPolicyAlways
}

func (s *Server) stopTrackedComponentLocked(ctx context.Context, componentKey string) (bool, error) {
	state, ok := s.takeRuntime(componentKey)
	if !ok {
		return false, nil
	}
	if state.StopFn != nil {
		if err := state.StopFn(ctx); err != nil {
			return false, err
		}
		return true, nil
	}
	if err := stopContainers(ctx, state.ContainerIDs); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Server) discoverOwnedContainers(ctx context.Context, fn func() error) ([]string, error) {
	client, err := dockerclient.New(dockerclient.FromEnv)
	if err != nil {
		s.lggr.Warn().Err(err).Msg("Docker unavailable for component ownership tracking; continuing without tracked dependencies")
		if runErr := fn(); runErr != nil {
			return nil, runErr
		}
		return []string{}, nil
	}
	defer client.Close()

	before, err := listContainerIDSet(ctx, client)
	if err != nil {
		return nil, err
	}

	eventsCtx, cancelEvents := context.WithCancel(ctx)
	defer cancelEvents()
	eventsResult := client.Events(eventsCtx, dockerclient.EventsListOptions{
		Filters: make(dockerclient.Filters).Add("type", "container"),
	})
	events := eventsResult.Messages
	errs := eventsResult.Err

	var wg sync.WaitGroup
	eventIDs := make([]string, 0)
	var eventMu sync.Mutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case msg, ok := <-events:
				if !ok {
					return
				}
				if msg.Action == "create" || msg.Action == "start" {
					eventMu.Lock()
					eventIDs = append(eventIDs, msg.Actor.ID)
					eventMu.Unlock()
				}
			case evtErr, ok := <-errs:
				if !ok || evtErr == nil {
					return
				}
				return
			case <-eventsCtx.Done():
				return
			}
		}
	}()

	runErr := fn()
	time.Sleep(150 * time.Millisecond)
	cancelEvents()
	wg.Wait()

	after, err := listContainerIDSet(ctx, client)
	if err != nil {
		if runErr != nil {
			return nil, runErr
		}
		return nil, err
	}

	owned := make([]string, 0)
	seen := make(map[string]struct{})
	for id := range after {
		if _, existed := before[id]; existed {
			continue
		}
		owned = append(owned, id)
		seen[id] = struct{}{}
	}
	eventMu.Lock()
	for _, id := range eventIDs {
		if _, ok := after[id]; !ok {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		owned = append(owned, id)
		seen[id] = struct{}{}
	}
	eventMu.Unlock()
	slices.Sort(owned)
	if runErr != nil {
		return owned, runErr
	}
	return owned, nil
}

func listContainerIDSet(ctx context.Context, client *dockerclient.Client) (map[string]struct{}, error) {
	listRes, err := client.ContainerList(ctx, dockerclient.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list docker containers: %w", err)
	}
	ids := make(map[string]struct{}, len(listRes.Items))
	for _, c := range listRes.Items {
		ids[c.ID] = struct{}{}
	}
	return ids, nil
}

func stopContainers(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	client, err := dockerclient.New(dockerclient.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client for stop: %w", err)
	}
	defer client.Close()

	namedVolumes, err := discoverNamedVolumesForContainers(ctx, client, ids)
	if err != nil {
		return err
	}

	for i := len(ids) - 1; i >= 0; i-- {
		_, err := client.ContainerRemove(ctx, ids[i], dockerclient.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil && !cerrdefs.IsNotFound(err) {
			return fmt.Errorf("failed to remove container %s: %w", ids[i], err)
		}
	}

	var removeVolumeErrors []error
	for _, volumeName := range namedVolumes {
		_, err := client.VolumeRemove(ctx, volumeName, dockerclient.VolumeRemoveOptions{Force: true})
		if err != nil && !cerrdefs.IsNotFound(err) {
			removeVolumeErrors = append(removeVolumeErrors, fmt.Errorf("remove volume %s: %w", volumeName, err))
		}
	}
	if len(removeVolumeErrors) > 0 {
		return fmt.Errorf("failed to remove one or more named volumes: %w", errors.Join(removeVolumeErrors...))
	}
	return nil
}

func discoverNamedVolumesForContainers(ctx context.Context, client *dockerclient.Client, ids []string) ([]string, error) {
	volumes := make(map[string]struct{})
	for _, id := range ids {
		inspectRes, err := client.ContainerInspect(ctx, id, dockerclient.ContainerInspectOptions{})
		if err != nil {
			if cerrdefs.IsNotFound(err) {
				continue
			}
			return nil, fmt.Errorf("inspect container %s before removal: %w", id, err)
		}
		for _, mountPoint := range inspectRes.Container.Mounts {
			if mountPoint.Type != mount.TypeVolume {
				continue
			}
			name := strings.TrimSpace(mountPoint.Name)
			if name == "" {
				continue
			}
			volumes[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(volumes))
	for name := range volumes {
		out = append(out, name)
	}
	slices.Sort(out)
	return out, nil
}

func hashPayload(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func shouldCleanupFailedContainers() bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(EnvKeepFailedContainers)))
	return raw == "" || (raw != "1" && raw != "true" && raw != "yes" && raw != "on")
}

func componentCacheKey(payload StartComponentPayload) (string, error) {
	switch payload.ComponentType {
	case ComponentTypeBlockchain:
		if payload.Blockchain == nil {
			return "", errors.New("blockchain payload is required")
		}
		return fmt.Sprintf("%s:%s:%s", payload.ComponentType, payload.Blockchain.Type, payload.Blockchain.ChainID), nil
	case ComponentTypeJD:
		if payload.JD == nil {
			return "", errors.New("jd payload is required")
		}
		return fmt.Sprintf("%s:%s", payload.ComponentType, payload.JD.Image), nil
	case ComponentTypeNodeSet:
		if payload.NodeSet == nil {
			return "", errors.New("nodeset payload is required")
		}
		return fmt.Sprintf("%s:%s", payload.ComponentType, payload.NodeSet.Name), nil
	default:
		return "", fmt.Errorf("unsupported component type: %s", payload.ComponentType)
	}
}

func Run(ctx context.Context, addr string, srv *Server) error {
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return httpSrv.Shutdown(context.Background())
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
