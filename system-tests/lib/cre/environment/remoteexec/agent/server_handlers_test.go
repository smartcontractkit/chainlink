package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
