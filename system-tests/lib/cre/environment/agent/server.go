package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

const (
	SchemaVersionV1          = "v1"
	OperationStartComponent  = "StartComponent"
	OperationHealth          = "Health"
	ComponentTypeBlockchain  = "blockchain"

	ErrCodeMethodNotAllowed       = "method_not_allowed"
	ErrCodeInvalidRequestBody     = "invalid_request_body"
	ErrCodeUnsupportedSchema      = "unsupported_schema_version"
	ErrCodeUnsupportedOperation   = "unsupported_operation"
	ErrCodeInvalidPayload         = "invalid_payload"
	ErrCodeUnsupportedComponent   = "unsupported_component_type"
	ErrCodeMissingComponentInput  = "missing_component_input"
	ErrCodeDeployFailed           = "deployment_failed"
	ErrCodeTransportEncodeFailed  = "transport_encode_failed"
)

var frameworkLogCaptureMu sync.Mutex

type StartComponentEnvelope struct {
	SchemaVersion string          `json:"schemaVersion"`
	Operation     string          `json:"operation"`
	Payload       json.RawMessage `json:"payload"`
}

type StartBlockchainPayload struct {
	ComponentType string            `json:"componentType"`
	Blockchain    *blockchain.Input `json:"blockchain"`
}

type StartComponentResponse struct {
	BlockchainOutput map[string]any `json:"blockchainOutput,omitempty"`
	AgentLogs        []string       `json:"agentLogs,omitempty"`
	ErrorCode        string         `json:"errorCode,omitempty"`
	Error            string         `json:"error,omitempty"`
}

type Server struct {
	lggr      zerolog.Logger
	deployers map[blockchain.ChainFamily]blockchains.Deployer
}

func NewServer(lggr zerolog.Logger, deployers map[blockchain.ChainFamily]blockchains.Deployer) *Server {
	return &Server{
		lggr:      lggr,
		deployers: deployers,
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
	if envelope.Operation != OperationStartComponent {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedOperation, fmt.Sprintf("unsupported operation: %s", envelope.Operation), nil)
		return
	}

	var payload StartBlockchainPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeInvalidPayload, fmt.Sprintf("invalid payload: %v", err), nil)
		return
	}
	if payload.ComponentType != ComponentTypeBlockchain {
		s.respondError(w, http.StatusBadRequest, ErrCodeUnsupportedComponent, fmt.Sprintf("unsupported component type: %s", payload.ComponentType), nil)
		return
	}
	if payload.Blockchain == nil {
		s.respondError(w, http.StatusBadRequest, ErrCodeMissingComponentInput, "blockchain payload is required", nil)
		return
	}

	// Keep this stderr write explicit so startup behavior is visible when agent runs as a subprocess.
	requestLog := fmt.Sprintf("[cre-agent] starting component type=%s blockchain=%s chain_id=%s", payload.ComponentType, payload.Blockchain.Type, payload.Blockchain.ChainID)
	_, _ = fmt.Fprintln(os.Stderr, requestLog)

	var startedOutput *blockchain.Output
	capturedFrameworkLogs, startErr := captureFrameworkLogs(func() error {
		deployed, err := DeployBlockchainComponent(r.Context(), s.deployers, payload.Blockchain)
		if err != nil {
			return err
		}
		startedOutput = deployed
		return nil
	})

	agentLogs := make([]string, 0, 1+len(capturedFrameworkLogs))
	agentLogs = append(agentLogs, requestLog)
	agentLogs = append(agentLogs, capturedFrameworkLogs...)

	if startErr != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeDeployFailed, startErr.Error(), agentLogs)
		return
	}

	safeOutput, encErr := EncodeForTransport(startedOutput)
	if encErr != nil {
		s.respondError(w, http.StatusInternalServerError, ErrCodeTransportEncodeFailed, encErr.Error(), agentLogs)
		return
	}
	s.respondJSON(w, http.StatusOK, StartComponentResponse{
		BlockchainOutput: safeOutput,
		AgentLogs:        agentLogs,
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
