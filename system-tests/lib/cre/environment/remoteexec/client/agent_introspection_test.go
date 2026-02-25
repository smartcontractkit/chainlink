package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/stretchr/testify/require"
)

func TestGetAgentStatusSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(agent.AgentStatusResponse{UptimeSeconds: 7})
	}))
	defer server.Close()

	resp, err := GetAgentStatus(context.Background(), &Runtime{AgentBaseURL: server.URL})
	require.NoError(t, err)
	require.Equal(t, int64(7), resp.UptimeSeconds)
}

func TestGetAgentLocksSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(agent.AgentLocksResponse{LifecycleBusy: true, RelayCount: 2})
	}))
	defer server.Close()

	resp, err := GetAgentLocks(context.Background(), &Runtime{AgentBaseURL: server.URL})
	require.NoError(t, err)
	require.True(t, resp.LifecycleBusy)
	require.Equal(t, 2, resp.RelayCount)
}

func TestGetComponentLogsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "nodeset:workflow", r.URL.Query().Get("componentKey"))
		require.Equal(t, "5", r.URL.Query().Get("limit"))
		_ = json.NewEncoder(w).Encode(agent.ComponentLogsResponse{
			ComponentKey: "nodeset:workflow",
			TotalLines:   8,
			Lines:        []string{"a", "b"},
		})
	}))
	defer server.Close()

	resp, err := GetComponentLogs(context.Background(), &Runtime{AgentBaseURL: server.URL}, "nodeset:workflow", 5)
	require.NoError(t, err)
	require.Equal(t, "nodeset:workflow", resp.ComponentKey)
	require.Equal(t, 8, resp.TotalLines)
	require.Equal(t, []string{"a", "b"}, resp.Lines)
}

func TestGetComponentLogsRequiresComponentKey(t *testing.T) {
	_, err := GetComponentLogs(context.Background(), &Runtime{AgentBaseURL: "http://127.0.0.1:1"}, "", 10)
	require.Error(t, err)
	require.Contains(t, err.Error(), "componentKey is required")
}

func TestGetAgentStatusPropagatesAgentError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{
			ErrorCode: "invalid_payload",
			Error:     "bad request",
		})
	}))
	defer server.Close()

	_, err := GetAgentStatus(context.Background(), &Runtime{AgentBaseURL: server.URL})
	require.Error(t, err)
	require.Contains(t, err.Error(), "remote agent error (invalid_payload): bad request")
}

func TestGetAgentStatusRequiresRuntime(t *testing.T) {
	_, err := GetAgentStatus(context.Background(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime is nil")
}
