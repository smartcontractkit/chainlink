package client

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestResolveRemoteRuntimeWithExplicitEnv(t *testing.T) {
	t.Setenv(EnvRemoteAgentURL, "http://198.51.100.20:19090")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "198.51.100.20")
	t.Setenv(EnvRemoteAgentPort, "19090")

	runtime, err := ResolveRuntime(zerolog.Nop())
	require.NoError(t, err, "expected runtime resolution to succeed")
	require.Equal(t, "http://198.51.100.20:19090", runtime.AgentBaseURL, "unexpected agent base url")
	require.Equal(t, "198.51.100.20", runtime.RemoteHostIP, "unexpected remote host ip")
	require.NotNil(t, runtime.Client, "expected resolved runtime to include component client")
}

func TestResolveRemoteRuntimeWithInputOverridesEnv(t *testing.T) {
	t.Setenv(EnvRemoteAgentURL, "http://198.51.100.20:19090")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "198.51.100.20")
	t.Setenv(EnvRemoteAgentPort, "19090")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/status", r.URL.Path)
		_ = json.NewEncoder(w).Encode(agent.AgentStatusResponse{
			ProtocolVersion: "1.0",
			Capabilities:    []string{"component_logs", "locks", "deploy_artifacts", "start_component", "relay", "list_ctf_resources"},
		})
	}))
	defer server.Close()

	runtime, err := ResolveRuntimeWithInput(zerolog.Nop(), RuntimeInput{
		AgentBaseURL: server.URL,
		RemoteHostIP: "203.0.113.22",
		AgentPort:    18081,
	})
	require.NoError(t, err)
	require.Equal(t, server.URL, runtime.AgentBaseURL)
	require.Equal(t, "203.0.113.22", runtime.RemoteHostIP)
}

func TestResolveRemoteRuntimeDerivesHostFromAgentURLWithoutAWSInputs(t *testing.T) {
	t.Setenv(EnvRemoteAgentURL, "http://198.51.100.20:19090")
	t.Setenv(runtimecfg.EnvRemoteHostIP, "")
	t.Setenv(runtimecfg.EnvRemoteAgentEC2InstanceID, "")

	runtime, err := ResolveRuntime(zerolog.Nop())
	require.NoError(t, err, "expected runtime resolution to derive host from explicit remote agent url")
	require.Equal(t, "198.51.100.20", runtime.RemoteHostIP, "expected host parsed from agent base URL")
}

func TestNewRemoteComponentClientRequiresResolvedRuntime(t *testing.T) {
	_, err := NewComponentClient(nil)
	require.Error(t, err, "expected nil runtime to fail")

	_, err = NewComponentClient(&Runtime{})
	require.Error(t, err, "expected missing agent base URL to fail")
}

func TestDescribeRemoteAgentHealthFailureMentionsResolutionHints(t *testing.T) {
	msg := describeRemoteAgentHealthFailure("http://203.0.113.10:8080")
	require.Contains(t, msg, "/v1/health")
	require.Contains(t, msg, EnvRemoteAgentPort)
	require.Contains(t, msg, EnvRemoteAgentURL)
}

func TestIsRetriableStatus(t *testing.T) {
	require.True(t, isRetriableStatus(502))
	require.True(t, isRetriableStatus(503))
	require.True(t, isRetriableStatus(504))
	require.False(t, isRetriableStatus(500))
}

func TestIsRetriableNetworkError(t *testing.T) {
	var netErr net.Error = timeoutError{}
	require.True(t, isRetriableNetworkError(netErr), "expected net.Error to be retriable")
	require.False(t, isRetriableNetworkError(errors.New("plain error")), "expected non-network error to be non-retriable")
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestStartComponentOnce_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{ComponentType: ComponentTypeBlockchain})
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	resp, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.NoError(t, err)
	require.Equal(t, ComponentTypeBlockchain, resp.ComponentType)
}

func TestStartComponentOnce_Non2xxWithAgentErrorCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{
			ErrorCode: "deployment_failed",
			Error:     "bad payload",
		})
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	_, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "remote agent error (deployment_failed)")
}

func TestStartComponentOnce_Non2xxWithoutAgentPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{})
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	_, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "start component request failed with status")
}

func TestStartComponentOnce_InvalidJSONResponseFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("{not-json"))
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	_, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode start component response")
}

func TestStartComponentOnce_200WithAgentErrorFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{
			ErrorCode: "deployment_failed",
			Error:     "start failed",
		})
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	_, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "remote agent error (deployment_failed)")
}
