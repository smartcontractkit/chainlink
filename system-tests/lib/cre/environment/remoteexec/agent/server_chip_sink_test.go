package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestChipSinkEventsEndpointReturnsEntriesFromLogFile(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)

	startReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/start", bytes.NewReader([]byte(`{"name":"sink-a","grpcListen":"127.0.0.1:0"}`)))
	startReq.Header.Set("Content-Type", "application/json")
	startRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(startRR, startReq)
	require.Equal(t, http.StatusOK, startRR.Code)

	var startResp ChipTestSinkStartResponse
	require.NoError(t, json.Unmarshal(startRR.Body.Bytes(), &startResp))
	require.NotEmpty(t, startResp.EventLogPath)
	t.Cleanup(func() {
		stopReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/stop", bytes.NewReader([]byte(`{}`)))
		stopReq.Header.Set("Content-Type", "application/json")
		stopRR := httptest.NewRecorder()
		server.Handler().ServeHTTP(stopRR, stopReq)
	})

	entry := ChipTestSinkEventLogEntry{
		Timestamp: time.Now().UTC().Add(1 * time.Second).Format(time.RFC3339Nano),
		Type:      "workflows.v1.UserLogs",
		Event:     map[string]any{"id": "abc"},
	}
	line, err := json.Marshal(entry)
	require.NoError(t, err)
	err = os.WriteFile(startResp.EventLogPath, append(line, '\n'), 0o600)
	require.NoError(t, err)

	eventsReq := httptest.NewRequest(http.MethodGet, "/v1/chip/sink/events?limit=10", nil)
	eventsRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(eventsRR, eventsReq)
	require.Equal(t, http.StatusOK, eventsRR.Code)

	var eventsResp ChipTestSinkEventsResponse
	require.NoError(t, json.Unmarshal(eventsRR.Body.Bytes(), &eventsResp))
	require.Len(t, eventsResp.Events, 1)
	require.Equal(t, "workflows.v1.UserLogs", eventsResp.Events[0].Type)
}

func TestStartChipSinkNormalizesBarePortListenAddress(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)

	startReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/start", bytes.NewReader([]byte(`{"name":"sink-a","grpcListen":"50052"}`)))
	startReq.Header.Set("Content-Type", "application/json")
	startRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(startRR, startReq)
	require.Equal(t, http.StatusOK, startRR.Code)

	var startResp ChipTestSinkStartResponse
	require.NoError(t, json.Unmarshal(startRR.Body.Bytes(), &startResp))
	require.True(t, strings.HasSuffix(startResp.GRPCListen, ":50052"), "expected normalized listen addr to bind port 50052, got %s", startResp.GRPCListen)

	stopReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/stop", bytes.NewReader([]byte(`{}`)))
	stopReq.Header.Set("Content-Type", "application/json")
	stopRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(stopRR, stopReq)
	require.Equal(t, http.StatusOK, stopRR.Code)
}
