package environment

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestResolveRemoteRuntimeWithExplicitEnv(t *testing.T) {
	t.Setenv(envEC2AgentURL, "http://198.51.100.20:19090")
	t.Setenv(runtimecfg.EnvEC2HostIP, "198.51.100.20")
	t.Setenv(envEC2AgentPort, "19090")

	runtime, err := resolveRemoteRuntime(zerolog.Nop())
	require.NoError(t, err, "expected runtime resolution to succeed")
	require.Equal(t, "http://198.51.100.20:19090", runtime.AgentBaseURL, "unexpected agent base url")
	require.Equal(t, "198.51.100.20", runtime.EC2HostIP, "unexpected ec2 host ip")
}

func TestResolveRemoteRuntimeRequiresHostResolution(t *testing.T) {
	t.Setenv(envEC2AgentURL, "http://198.51.100.20:19090")
	t.Setenv(runtimecfg.EnvEC2HostIP, "")
	t.Setenv(runtimecfg.EnvEC2InstanceID, "")

	_, err := resolveRemoteRuntime(zerolog.Nop())
	require.Error(t, err, "expected runtime resolution without EC2 host inputs to fail")
}

func TestNewRemoteComponentClientRequiresResolvedRuntime(t *testing.T) {
	_, err := newRemoteComponentClient(nil)
	require.Error(t, err, "expected nil runtime to fail")

	_, err = newRemoteComponentClient(&resolvedRemoteRuntime{})
	require.Error(t, err, "expected missing agent base URL to fail")
}

func TestDescribeEC2AgentHealthFailureMentionsResolutionHints(t *testing.T) {
	msg := describeEC2AgentHealthFailure("http://203.0.113.10:8080")
	require.Contains(t, msg, "/v1/health")
	require.Contains(t, msg, envEC2AgentPort)
	require.Contains(t, msg, envEC2AgentURL)
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
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{ComponentType: componentTypeBlockchain})
	}))
	defer server.Close()

	client := &httpComponentClient{baseURL: server.URL, client: server.Client()}
	resp, err := client.startComponentOnce(context.Background(), agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       json.RawMessage(`{"componentType":"blockchain"}`),
	})
	require.NoError(t, err)
	require.Equal(t, componentTypeBlockchain, resp.ComponentType)
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
