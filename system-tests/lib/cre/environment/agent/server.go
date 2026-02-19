package agent

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	dockerevents "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

const (
	SchemaVersionV1          = "v1"
	OperationStartComponent  = "StartComponent"
	OperationStopComponent   = "StopComponent"
	OperationHealth          = "Health"
	ComponentTypeBlockchain  = "blockchain"
	ComponentTypeJD          = "jd"
	ComponentTypeNodeSet     = "nodeset"

	ErrCodeMethodNotAllowed       = "method_not_allowed"
	ErrCodeInvalidRequestBody     = "invalid_request_body"
	ErrCodeUnsupportedSchema      = "unsupported_schema_version"
	ErrCodeUnsupportedOperation   = "unsupported_operation"
	ErrCodeInvalidPayload         = "invalid_payload"
	ErrCodeUnsupportedComponent   = "unsupported_component_type"
	ErrCodeMissingComponentInput  = "missing_component_input"
	ErrCodeDeployFailed           = "deployment_failed"
	ErrCodeTransportEncodeFailed  = "transport_encode_failed"

	RemoteStartPolicyAlways       = "always"
	RemoteStartPolicyReuseIdentical = "reuse_if_identical"

	EnvKeepFailedContainers = "CRE_AGENT_KEEP_FAILED_CONTAINERS"
)

var frameworkLogCaptureMu sync.Mutex

type StartComponentEnvelope struct {
	SchemaVersion string          `json:"schemaVersion"`
	Operation     string          `json:"operation"`
	Payload       json.RawMessage `json:"payload"`
}

type StartComponentPayload struct {
	ComponentType      string             `json:"componentType"`
	Blockchain         *blockchain.Input  `json:"blockchain"`
	RegistryBlockchain map[string]any     `json:"registryBlockchain,omitempty"`
	JD                 *jd.Input          `json:"jd"`
	NodeSet            *ns.Input          `json:"nodeset,omitempty"`
	ReusePolicy        string             `json:"reusePolicy,omitempty"`
}

type StartComponentResponse struct {
	ComponentType string         `json:"componentType,omitempty"`
	Output        map[string]any `json:"output,omitempty"`
	Found         bool           `json:"found,omitempty"`
	Stopped       bool           `json:"stopped,omitempty"`
	AgentLogs        []string       `json:"agentLogs,omitempty"`
	ErrorCode        string         `json:"errorCode,omitempty"`
	Error            string         `json:"error,omitempty"`
}

type Server struct {
	lggr      zerolog.Logger
	deployers map[blockchain.ChainFamily]blockchains.Deployer
	lifecycleMu sync.Mutex
	cacheMu   sync.Mutex
	cache     map[string]cachedStart
	runtime   map[string]runtimeState
}

type cachedStart struct {
	PayloadHash string
	Output      map[string]any
}

type runtimeState struct {
	ComponentType string
	ContainerIDs  []string
}

func NewServer(lggr zerolog.Logger, deployers map[blockchain.ChainFamily]blockchains.Deployer) *Server {
	return &Server{
		lggr:      lggr,
		deployers: deployers,
		cache:     make(map[string]cachedStart),
		runtime:   make(map[string]runtimeState),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health", s.health)
	mux.HandleFunc("/v1/components/start", s.startComponent)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
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
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedSchema, fmt.Sprintf("unsupported schema version: %s", envelope.SchemaVersion), nil)
		return
	}
	var payload StartComponentPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("invalid payload: %v", err), nil)
		return
	}
	if payload.ComponentType != ComponentTypeBlockchain && payload.ComponentType != ComponentTypeJD && payload.ComponentType != ComponentTypeNodeSet {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedComponent, fmt.Sprintf("unsupported component type: %s", payload.ComponentType), nil)
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
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedOperation, fmt.Sprintf("unsupported operation: %s", envelope.Operation), nil)
		return
	}
	payloadHash := hashPayload(envelope.Payload)

	// Keep this stderr write explicit so startup behavior is visible when agent runs as a subprocess.
	requestLog := fmt.Sprintf("[cre-agent] starting component type=%s key=%s", payload.ComponentType, componentKey)
	_, _ = fmt.Fprintln(os.Stderr, requestLog)
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
				AgentLogs: []string{requestLog, reuseLog},
			})
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
		return
	}

	var output map[string]any
	var encErr error
	if blockchainOutput != nil {
		output, encErr = EncodeForTransport(blockchainOutput)
	} else if jdOutput != nil {
		output, encErr = EncodeForTransport(jdOutput)
	} else if nodeSetOutput != nil {
		output, encErr = EncodeForTransport(nodeSetOutput)
	}
	if encErr != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeTransportEncodeFailed, encErr.Error(), agentLogs)
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
}

func (s *Server) stopComponentByKey(w http.ResponseWriter, r *http.Request, componentType, componentKey string) {
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
		return
	}
	s.deleteCachedStart(componentKey)
	s.respondJSON(w, http.StatusOK, StartComponentResponse{
		ComponentType: componentType,
		Found:         true,
		Stopped:       true,
		AgentLogs:     []string{requestLog, "[cre-agent] stopped existing component"},
	})
}

func (s *Server) respondJSON(w http.ResponseWriter, code int, body StartComponentResponse) {
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

func (s *Server) lookupCachedStart(componentKey, payloadHash string) (*cachedStart, bool) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	start, ok := s.cache[componentKey]
	if !ok || start.PayloadHash != payloadHash {
		return nil, false
	}
	return &start, true
}

func (s *Server) cacheSuccessfulStart(componentKey, payloadHash string, output map[string]any) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache[componentKey] = cachedStart{
		PayloadHash: payloadHash,
		Output:      output,
	}
}

func (s *Server) deleteCachedStart(componentKey string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	delete(s.cache, componentKey)
}

func (s *Server) storeRuntime(componentKey string, state runtimeState) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.runtime[componentKey] = state
}

func (s *Server) takeRuntime(componentKey string) (runtimeState, bool) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	state, ok := s.runtime[componentKey]
	if ok {
		delete(s.runtime, componentKey)
	}
	return state, ok
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
	if err := stopContainers(ctx, state.ContainerIDs); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Server) discoverOwnedContainers(ctx context.Context, fn func() error) ([]string, error) {
	client, err := dockerclient.NewClientWithOpts(dockerclient.WithAPIVersionNegotiation())
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
	events, errs := client.Events(eventsCtx, dockerevents.ListOptions{
		Filters: filters.NewArgs(filters.Arg("type", "container")),
	})

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
					eventIDs = append(eventIDs, msg.ID)
					eventMu.Unlock()
				}
			case err, ok := <-errs:
				if !ok || err == nil {
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
	containers, err := client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list docker containers: %w", err)
	}
	ids := make(map[string]struct{}, len(containers))
	for _, c := range containers {
		ids[c.ID] = struct{}{}
	}
	return ids, nil
}

func stopContainers(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	client, err := dockerclient.NewClientWithOpts(dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client for stop: %w", err)
	}
	defer client.Close()

	for i := len(ids) - 1; i >= 0; i-- {
		err := client.ContainerRemove(ctx, ids[i], container.RemoveOptions{Force: true})
		if err != nil && !cerrdefs.IsNotFound(err) {
			return fmt.Errorf("failed to remove container %s: %w", ids[i], err)
		}
	}
	return nil
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
			return "", fmt.Errorf("blockchain payload is required")
		}
		return fmt.Sprintf("%s:%s:%s", payload.ComponentType, payload.Blockchain.Type, payload.Blockchain.ChainID), nil
	case ComponentTypeJD:
		if payload.JD == nil {
			return "", fmt.Errorf("jd payload is required")
		}
		return fmt.Sprintf("%s:%s", payload.ComponentType, payload.JD.Image), nil
	case ComponentTypeNodeSet:
		if payload.NodeSet == nil {
			return "", fmt.Errorf("nodeset payload is required")
		}
		return fmt.Sprintf("%s:%s", payload.ComponentType, payload.NodeSet.Name), nil
	default:
		return "", fmt.Errorf("unsupported component type: %s", payload.ComponentType)
	}
}

func Run(ctx context.Context, addr string, srv *Server) error {
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: srv.Handler(),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return httpSrv.Shutdown(context.Background())
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
