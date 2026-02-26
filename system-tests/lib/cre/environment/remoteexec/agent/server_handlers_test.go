package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpointReturnsOK(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "ok", rr.Body.String())
}

func TestListCTFResourcesMethodNotAllowed(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/resources/ctf", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Contains(t, rr.Body.String(), ErrCodeMethodNotAllowed)
}

func TestDeployArtifactsValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		payload     DeployArtifactsPayload
		wantCode    int
		wantErrCode string
		wantMsg     string
	}{
		{
			name:        "missing nodeset name",
			payload:     DeployArtifactsPayload{NodeSetName: "", TargetDir: "/tmp", Files: []DeployArtifactsFile{{Name: "a.txt", ContentBase64: base64.StdEncoding.EncodeToString([]byte("x"))}}},
			wantCode:    http.StatusBadRequest,
			wantErrCode: ErrCodeMissingComponentInput,
			wantMsg:     "nodeset name is required",
		},
		{
			name:        "missing target dir",
			payload:     DeployArtifactsPayload{NodeSetName: "workflow", TargetDir: "", Files: []DeployArtifactsFile{{Name: "a.txt", ContentBase64: base64.StdEncoding.EncodeToString([]byte("x"))}}},
			wantCode:    http.StatusBadRequest,
			wantErrCode: ErrCodeMissingComponentInput,
			wantMsg:     "target dir is required",
		},
		{
			name:        "no files",
			payload:     DeployArtifactsPayload{NodeSetName: "workflow", TargetDir: "/tmp"},
			wantCode:    http.StatusBadRequest,
			wantErrCode: ErrCodeMissingComponentInput,
			wantMsg:     "at least one artifact file is required",
		},
	}

	server := NewServer(zerolog.Nop(), nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envelope := StartComponentEnvelope{
				SchemaVersion: SchemaVersionV1,
				Operation:     OperationDeployArtifacts,
			}
			payloadRaw, err := json.Marshal(tt.payload)
			require.NoError(t, err)
			envelope.Payload = payloadRaw

			reqBody, err := json.Marshal(envelope)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			server.Handler().ServeHTTP(rr, req)
			require.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Body.String(), tt.wantErrCode)
			require.Contains(t, rr.Body.String(), tt.wantMsg)
		})
	}
}

func TestComponentCacheKeyVariants(t *testing.T) {
	key, err := componentCacheKey(StartComponentPayload{
		ComponentType: ComponentTypeJD,
		JD:            &jd.Input{Image: "job-distributor:0.22.1"},
	})
	require.NoError(t, err)
	require.Contains(t, key, ComponentTypeJD)

	key, err = componentCacheKey(StartComponentPayload{
		ComponentType: ComponentTypeNodeSet,
		NodeSet:       &ns.Input{Name: "workflow"},
	})
	require.NoError(t, err)
	require.Equal(t, "nodeset:workflow", key)

	_, err = componentCacheKey(StartComponentPayload{ComponentType: ComponentTypeNodeSet})
	require.Error(t, err)
	require.Contains(t, err.Error(), "nodeset payload is required")

	_, err = componentCacheKey(StartComponentPayload{ComponentType: "unknown"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported component type")
}

func TestStatusEndpointReturnsAgentState(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	server.cacheSuccessfulStart("blockchain:anvil:1337", "hash-a", map[string]any{"ok": true})
	server.storeRuntime("nodeset:workflow", runtimeState{ComponentType: ComponentTypeNodeSet})
	server.appendComponentLogs("nodeset:workflow", []string{"line-a"})
	server.beginInFlight("start:nodeset:workflow", inFlightOperationScopeLifecycle)
	defer server.endInFlight("start:nodeset:workflow")

	openReq := httptest.NewRequest(http.MethodPost, "/v1/relay/open", bytes.NewReader([]byte(`{"name":"workflow-ocr-0","requestedPort":0}`)))
	openReq.Header.Set("Content-Type", "application/json")
	openRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(openRR, openReq)
	require.Equal(t, http.StatusOK, openRR.Code)

	req := httptest.NewRequest(http.MethodGet, "/v1/status", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp AgentStatusResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.GreaterOrEqual(t, resp.UptimeSeconds, int64(0))
	require.Contains(t, resp.CachedComponents, "blockchain:anvil:1337")
	require.Contains(t, resp.RuntimeComponents, "nodeset:workflow")
	require.Contains(t, resp.ComponentLogKeys, "nodeset:workflow")
	require.Len(t, resp.Relays, 1)
	require.Equal(t, "workflow-ocr-0", resp.Relays[0].Name)
	require.Greater(t, resp.Relays[0].BoundPort, 0)
	require.Len(t, resp.InFlight, 1)
}

func TestLocksEndpointShowsLifecycleBusy(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	server.cacheSuccessfulStart("blockchain:anvil:1337", "hash-a", map[string]any{"ok": true})
	server.storeRuntime("nodeset:workflow", runtimeState{ComponentType: ComponentTypeNodeSet})
	server.appendComponentLogs("nodeset:workflow", []string{"line-a"})
	server.beginInFlight("start:nodeset:workflow", inFlightOperationScopeLifecycle)
	defer server.endInFlight("start:nodeset:workflow")

	req := httptest.NewRequest(http.MethodGet, "/v1/locks", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp AgentLocksResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.True(t, resp.LifecycleBusy)
	require.Equal(t, 1, resp.CacheEntries)
	require.Equal(t, 1, resp.RuntimeEntries)
	require.Equal(t, 1, resp.ComponentLogKeys)
	require.Len(t, resp.InFlight, 1)
}

func TestComponentLogsEndpointValidationAndLimit(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	server.appendComponentLogs("nodeset:workflow", []string{"line-a", "line-b", "line-c"})
	time.Sleep(1 * time.Millisecond)

	reqMissingKey := httptest.NewRequest(http.MethodGet, "/v1/components/logs", nil)
	rrMissingKey := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrMissingKey, reqMissingKey)
	require.Equal(t, http.StatusBadRequest, rrMissingKey.Code)
	require.Contains(t, rrMissingKey.Body.String(), "componentKey query parameter is required")

	reqInvalidLimit := httptest.NewRequest(http.MethodGet, "/v1/components/logs?componentKey=nodeset:workflow&limit=abc", nil)
	rrInvalidLimit := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrInvalidLimit, reqInvalidLimit)
	require.Equal(t, http.StatusBadRequest, rrInvalidLimit.Code)
	require.Contains(t, rrInvalidLimit.Body.String(), "limit query parameter must be a positive integer")

	req := httptest.NewRequest(http.MethodGet, "/v1/components/logs?componentKey=nodeset:workflow&limit=2", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var resp ComponentLogsResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "nodeset:workflow", resp.ComponentKey)
	require.Equal(t, 3, resp.TotalLines)
	require.Equal(t, []string{"line-b", "line-c"}, resp.Lines)
}

func TestChipTestSinkLifecycleEndpoints(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)

	startReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/start", bytes.NewReader([]byte(`{"name":"sink-a","grpcListen":"127.0.0.1:0"}`)))
	startReq.Header.Set("Content-Type", "application/json")
	startRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(startRR, startReq)
	require.Equal(t, http.StatusOK, startRR.Code)

	var startResp ChipTestSinkStartResponse
	require.NoError(t, json.Unmarshal(startRR.Body.Bytes(), &startResp))
	require.Equal(t, "sink", startResp.Profile)
	require.Equal(t, "remote", startResp.Mode)
	require.Equal(t, "sink-a", startResp.Name)
	require.NotEmpty(t, startResp.GRPCListen)

	statusReq := httptest.NewRequest(http.MethodGet, "/v1/chip/sink/status", nil)
	statusRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(statusRR, statusReq)
	require.Equal(t, http.StatusOK, statusRR.Code)

	var statusResp ChipTestSinkStatusResponse
	require.NoError(t, json.Unmarshal(statusRR.Body.Bytes(), &statusResp))
	require.True(t, statusResp.Running)
	require.Equal(t, "sink-a", statusResp.Name)

	stopReq := httptest.NewRequest(http.MethodPost, "/v1/chip/sink/stop", bytes.NewReader([]byte(`{}`)))
	stopReq.Header.Set("Content-Type", "application/json")
	stopRR := httptest.NewRecorder()
	server.Handler().ServeHTTP(stopRR, stopReq)
	require.Equal(t, http.StatusOK, stopRR.Code)

	var stopResp ChipTestSinkStopResponse
	require.NoError(t, json.Unmarshal(stopRR.Body.Bytes(), &stopResp))
	require.True(t, stopResp.Found)
	require.True(t, stopResp.Stopped)
}
