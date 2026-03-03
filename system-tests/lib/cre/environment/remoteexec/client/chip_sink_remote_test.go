package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

func TestStartRemoteChipTestSinkSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chip/sink/start", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		_ = json.NewEncoder(w).Encode(agent.ChipTestSinkStartResponse{
			Profile:    "sink",
			Mode:       "remote",
			Name:       "default",
			GRPCListen: "0.0.0.0:50051",
		})
	}))
	defer server.Close()

	resp, err := StartRemoteChipTestSink(context.Background(), &Runtime{AgentBaseURL: server.URL}, agent.ChipTestSinkStartRequest{})
	require.NoError(t, err)
	require.Equal(t, "sink", resp.Profile)
	require.Equal(t, "remote", resp.Mode)
}

func TestStopRemoteChipTestSinkSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chip/sink/stop", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		_ = json.NewEncoder(w).Encode(agent.ChipTestSinkStopResponse{Found: true, Stopped: true})
	}))
	defer server.Close()

	resp, err := StopRemoteChipTestSink(context.Background(), &Runtime{AgentBaseURL: server.URL})
	require.NoError(t, err)
	require.True(t, resp.Found)
	require.True(t, resp.Stopped)
}

func TestGetRemoteChipTestSinkStatusSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chip/sink/status", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_ = json.NewEncoder(w).Encode(agent.ChipTestSinkStatusResponse{
			Profile:    "sink",
			Mode:       "remote",
			Running:    true,
			Name:       "default",
			GRPCListen: "0.0.0.0:50051",
		})
	}))
	defer server.Close()

	resp, err := GetRemoteChipTestSinkStatus(context.Background(), &Runtime{AgentBaseURL: server.URL})
	require.NoError(t, err)
	require.True(t, resp.Running)
	require.Equal(t, "sink", resp.Profile)
}

func TestGetRemoteChipTestSinkEventsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chip/sink/events", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		_ = json.NewEncoder(w).Encode(agent.ChipTestSinkEventsResponse{
			Events: []agent.ChipTestSinkEventLogEntry{{Type: "workflows.v1.UserLogs"}},
		})
	}))
	defer server.Close()

	resp, err := GetRemoteChipTestSinkEvents(context.Background(), &Runtime{AgentBaseURL: server.URL}, time.Time{}, 5)
	require.NoError(t, err)
	require.Len(t, resp.Events, 1)
	require.Equal(t, "workflows.v1.UserLogs", resp.Events[0].Type)
}
