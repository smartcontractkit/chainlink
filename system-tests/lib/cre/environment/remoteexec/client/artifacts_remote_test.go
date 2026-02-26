package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/stretchr/testify/require"
)

func TestDeployArtifactsToRemoteNodeSetValidation(t *testing.T) {
	err := DeployArtifactsToRemoteNodeSet(context.Background(), zerolog.Nop(), "", "/tmp", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nodeset name is required")

	err = DeployArtifactsToRemoteNodeSet(context.Background(), zerolog.Nop(), "workflow", "", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "container target dir is required")
}

func TestDeployArtifactsToRemoteNodeSetNoFilesFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/v1/status" {
			_ = json.NewEncoder(w).Encode(agent.AgentStatusResponse{ProtocolVersion: "1.0.0"})
			return
		}
		t.Fatalf("unexpected path %s", r.URL.Path)
	}))
	defer server.Close()

	t.Setenv(EnvRemoteAgentURL, server.URL)
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	err := DeployArtifactsToRemoteNodeSet(context.Background(), zerolog.Nop(), "workflow", "/home/chainlink/workflows", []string{"", ""})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no artifact files to deploy")
}

func TestDeployArtifactsToRemoteNodeSetSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	artifactPath := filepath.Join(tmpDir, "artifact.wasm")
	require.NoError(t, os.WriteFile(artifactPath, []byte("artifact-content"), 0o600))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/health":
			w.WriteHeader(http.StatusOK)
		case "/v1/status":
			_ = json.NewEncoder(w).Encode(agent.AgentStatusResponse{ProtocolVersion: "1.0.0"})
		case "/v1/components/start":
			var envelope agent.StartComponentEnvelope
			require.NoError(t, json.NewDecoder(r.Body).Decode(&envelope))
			require.Equal(t, agent.OperationDeployArtifacts, envelope.Operation)

			var payload agent.DeployArtifactsPayload
			require.NoError(t, json.Unmarshal(envelope.Payload, &payload))
			require.Equal(t, "workflow", payload.NodeSetName)
			require.Equal(t, "/home/chainlink/workflows", payload.TargetDir)
			require.Len(t, payload.Files, 1)
			require.Equal(t, "artifact.wasm", payload.Files[0].Name)
			raw, err := base64.StdEncoding.DecodeString(payload.Files[0].ContentBase64)
			require.NoError(t, err)
			require.Equal(t, "artifact-content", string(raw))

			_ = json.NewEncoder(w).Encode(agent.StartComponentResponse{
				ComponentType: ComponentTypeNodeSet,
				AgentLogs:     []string{"artifact deployed"},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	t.Setenv(EnvRemoteAgentURL, server.URL)
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	err := DeployArtifactsToRemoteNodeSet(context.Background(), zerolog.Nop(), "workflow", "/home/chainlink/workflows", []string{artifactPath})
	require.NoError(t, err)
}
