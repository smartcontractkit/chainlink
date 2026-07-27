package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

// setupTestServer builds a chart with two node-set instances so that
// chart-derived DON membership can be exercised for both a workflow DON
// ("workflow" don-name, containing bootstrap + worker nodes) and a gateway
// DON ("gateway-don" don-name, containing the gateway node) — mirroring how
// membership is always derived from cv.NodeNamesForDONName, never from a
// desired.toml "nodes" field.
func setupTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	dir := t.TempDir()

	desiredPath := filepath.Join(dir, "desired.toml")
	statePath := filepath.Join(dir, "state.toml")

	workflowDir := filepath.Join(dir, "deploy", "config", "workflow")
	require.NoError(t, os.MkdirAll(workflowDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(workflowDir, "shared.yaml"), []byte(`
chainlink-node:
  registerNodes:
    labels:
      don-name: "workflow"
  defaults:
    nodeType: standard
  instances:
    node-bt-0:
      nodeType: boot
    node-0: {}
    node-1: {}
    node-2: {}
    node-3: {}
`), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(workflowDir, "dev.yaml"), []byte(`
anvil:
  instances:
    anvil-1337:
      chainID: 1337
      gateway:
        hostnames:
          - anvil-1337.my-repo.dev.internal.griddle.sh
`), 0600))

	gatewayDir := filepath.Join(dir, "deploy", "config", "gateway-don")
	require.NoError(t, os.MkdirAll(gatewayDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(gatewayDir, "shared.yaml"), []byte(`
chainlink-node:
  registerNodes:
    labels:
      don-name: "gateway-don"
  defaults:
    nodeType: standard
  instances:
    node-gw-0:
      nodeType: gateway
`), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(gatewayDir, "dev.yaml"), []byte(`{}`), 0600))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "griddle.yaml"), []byte(`
deploy:
  dev:
    instances:
      - name: workflow
        namespace: my-repo-nodeset
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/workflow/shared.yaml
          - deploy/config/workflow/dev.yaml
      - name: gateway-don
        namespace: my-repo-nodeset
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/gateway-don/shared.yaml
          - deploy/config/gateway-don/dev.yaml
`), 0600))

	log := zerolog.Nop()
	server := NewServer(desiredPath, statePath, dir, "dev", "my-repo-nodeset", "", log)

	return server, dir
}

func TestAPI_Nodes(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/nodes", nil)
	w := httptest.NewRecorder()
	server.handleNodes(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp NodesResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.NotEmpty(t, resp.Nodes)
	require.Equal(t, "my-repo-nodeset", resp.Namespace)
	require.Equal(t, "node-bt-0", resp.Bootstrap)
	require.Contains(t, resp.Gateways, "node-gw-0")
	require.NotEmpty(t, resp.ChartDir)

	for _, n := range resp.Nodes {
		if n.Name == "node-0" {
			require.Equal(t, "workflow", n.ChartDonName)
		}
		if n.Name == "node-gw-0" {
			require.Equal(t, "gateway-don", n.ChartDonName)
		}
	}
}

func TestAPI_ChainsDiscovered_Empty(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/chains/discovered", nil)
	w := httptest.NewRecorder()
	server.handleChainsDiscovered(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Chains []ChainResponse `json:"chains"`
	}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	// The test chart's nodes have no "configuration" TOML layers with [[EVM]]
	// sections, so nothing should be discovered.
	require.Empty(t, resp.Chains)
}

func TestAPI_Desired_GetEmpty(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/desired", nil)
	w := httptest.NewRecorder()
	server.handleDesired(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp DesiredResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.Empty(t, resp.DONs)
}

func TestAPI_Desired_SaveAndLoad(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	// Save a desired state. "nodes" is no longer part of the schema — DON
	// membership is always chart-derived from the don-name label, never from
	// desired.toml.
	payload := `{
		"infra": {"type": "griddle", "chartValues": "deploy/config/workflow", "namespace": "my-repo-nodeset"},
		"jd": {"grpc": "grpc-jd:443", "domain": "cre", "environment": "dev"},
		"chains": [
			{"chainId": 1337, "wsUrl": "wss://anvil-1337.example.com", "httpUrl": "https://anvil-1337.example.com", "registry": true}
		],
		"dons": [
			{
				"name": "workflow",
				"donTypes": ["workflow"],
				"capabilities": ["cron", "evm-1337"],
				"exposesRemoteCapabilities": false
			}
		],
		"capabilityConfigs": {
			"cron": {"binary_name": "cron"},
			"evm": {"binary_name": "evm"}
		}
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/desired", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleDesired(w, req)

	require.Equal(t, http.StatusOK, w.Code, "POST failed: %s", w.Body.String())

	// Load it back — membership is derived server-side from the chart, not
	// from anything persisted in desired.toml.
	req2 := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/desired", nil)
	w2 := httptest.NewRecorder()
	server.handleDesired(w2, req2)

	require.Equal(t, http.StatusOK, w2.Code)

	var resp DesiredResponse
	require.NoError(t, json.NewDecoder(w2.Body).Decode(&resp))
	require.Len(t, resp.DONs, 1)
	require.Equal(t, "workflow", resp.DONs[0].Name)
	require.Contains(t, resp.DONs[0].Capabilities, "cron")
	require.ElementsMatch(t, []string{"node-bt-0", "node-0", "node-1", "node-2", "node-3"}, resp.DONs[0].Nodes)
	require.Len(t, resp.Chains, 1)
	require.Equal(t, uint64(1337), resp.Chains[0].ChainID)
	require.True(t, resp.Chains[0].Registry)
}

func TestAPI_Desired_IgnoresPostedNodes(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	// Even if a legacy client posts "nodes", the server must not persist or
	// trust it — membership always comes from the chart.
	payload := `{
		"infra": {"type": "griddle", "chartValues": "deploy/config/workflow", "namespace": "my-repo-nodeset"},
		"jd": {"grpc": "grpc-jd:443", "domain": "cre", "environment": "dev"},
		"chains": [
			{"chainId": 1337, "wsUrl": "wss://anvil-1337.example.com", "httpUrl": "https://anvil-1337.example.com", "registry": true}
		],
		"dons": [
			{
				"name": "workflow",
				"donTypes": ["workflow"],
				"capabilities": ["cron"],
				"nodes": ["totally-made-up-node"],
				"exposesRemoteCapabilities": false
			}
		],
		"capabilityConfigs": {
			"cron": {"binary_name": "cron"}
		}
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/desired", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleDesired(w, req)
	require.Equal(t, http.StatusOK, w.Code, "POST failed: %s", w.Body.String())

	raw, err := os.ReadFile(server.desiredPath)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "totally-made-up-node", "posted nodes must never be persisted to desired.toml")

	req2 := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/desired", nil)
	w2 := httptest.NewRecorder()
	server.handleDesired(w2, req2)
	require.Equal(t, http.StatusOK, w2.Code)

	var resp DesiredResponse
	require.NoError(t, json.NewDecoder(w2.Body).Decode(&resp))
	require.ElementsMatch(t, []string{"node-bt-0", "node-0", "node-1", "node-2", "node-3"}, resp.DONs[0].Nodes)
}

func TestAPI_Desired_GatewayNodesRoundTrip(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	payload := `{
		"infra": {"type": "griddle", "chartValues": "deploy/config/workflow", "namespace": "my-repo-nodeset"},
		"jd": {"grpc": "grpc-jd:443", "domain": "cre", "environment": "dev"},
		"chains": [
			{"chainId": 1337, "wsUrl": "wss://anvil-1337.example.com", "httpUrl": "https://anvil-1337.example.com", "registry": true}
		],
		"dons": [
			{
				"name": "workflow",
				"donTypes": ["workflow"],
				"capabilities": ["cron", "http-action"],
				"exposesRemoteCapabilities": false
			},
			{
				"name": "gateway-don",
				"donTypes": ["gateway"],
				"capabilities": [],
				"exposesRemoteCapabilities": false
			}
		],
		"gatewayNodes": [
			{"node": "node-gw-0", "don": "workflow"}
		],
		"capabilityConfigs": {
			"cron": {"binary_name": "cron"},
			"http-action": {"binary_name": "http_action"}
		}
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/desired", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleDesired(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	req2 := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/desired", nil)
	w2 := httptest.NewRecorder()
	server.handleDesired(w2, req2)
	require.Equal(t, http.StatusOK, w2.Code)

	var resp DesiredResponse
	require.NoError(t, json.NewDecoder(w2.Body).Decode(&resp))
	require.Len(t, resp.GatewayNodes, 1)
	require.Equal(t, "node-gw-0", resp.GatewayNodes[0].Node)
	require.Equal(t, "workflow", resp.GatewayNodes[0].DON)

	for _, don := range resp.DONs {
		if don.Name == "gateway-don" {
			require.ElementsMatch(t, []string{"node-gw-0"}, don.Nodes)
		}
	}
}

func TestAPI_Desired_ValidationFails(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	payload := `{
		"infra": {"type": "docker", "chartValues": "x", "namespace": "ns"},
		"jd": {"grpc": "x", "domain": "cre", "environment": "dev"},
		"dons": []
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/desired", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleDesired(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_State_Empty(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	server.handleState(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp StateResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.False(t, resp.HasState)
}

func TestAPI_State_WithContent(t *testing.T) {
	t.Parallel()

	server, dir := setupTestServer(t)

	// Write a state file
	statePath := filepath.Join(dir, "state.toml")
	s := &domain.StateFile{
		Phase: domain.PhaseProvisioned,
		Addresses: []domain.AddressRef{
			{ChainSelector: 1337, Address: "0xabc", Type: "CapabilitiesRegistry", Version: "v2"},
		},
		DONIDs: map[string]uint64{"workflow": 1},
	}
	require.NoError(t, s.Store(statePath))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	server.handleState(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp StateResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.True(t, resp.HasState)
	require.Equal(t, "provisioned", resp.Phase)
	require.NotEmpty(t, resp.Addresses)
	require.Equal(t, uint64(1), resp.DONIDs["workflow"])
}

func TestAPI_Capabilities(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/capabilities", nil)
	w := httptest.NewRecorder()
	server.handleCapabilities(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))

	caps, ok := resp["capabilities"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, caps)

	for _, c := range caps {
		cm := c.(map[string]any)
		if cm["name"] == "evm" {
			require.True(t, cm["chainScoped"].(bool))
		}
		if cm["name"] == "cron" {
			require.False(t, cm["chainScoped"].(bool))
		}
	}
}

func TestAPI_PreviewTOML(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	payload := `{
		"infra": {"type": "griddle", "chartValues": "deploy/config/workflow", "namespace": "my-repo-nodeset"},
		"jd": {"grpc": "grpc-jd:443", "domain": "cre", "environment": "dev"},
		"dons": [
			{
				"name": "workflow",
				"donTypes": ["workflow"],
				"capabilities": ["cron"],
				"exposesRemoteCapabilities": false
			}
		],
		"capabilityConfigs": {
			"cron": {"binary_name": "cron"}
		}
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/preview-toml", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handlePreviewTOML(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	require.Contains(t, body, "workflow")
	require.Contains(t, body, "cron")
}

// TestAPI_JDCheck_FailsFastWithoutAccessToken proves that clicking "Check JD
// Connectivity" with no GRIDDLE_JD_ACCESS_TOKEN on the server returns an
// immediate, actionable error instead of attempting a dial that can only ever
// fail with a confusing gRPC/auth error — JD always requires this token.
func TestAPI_JDCheck_FailsFastWithoutAccessToken(t *testing.T) {
	t.Setenv(infra.JDAccessTokenEnv, "")
	server, _ := setupTestServer(t)

	payload := `{"grpc": "grpc-job-distributor.example.com:443", "useTLS": true}`
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/jd/check", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleJDCheck(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp JDCheckResponseFull
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	require.False(t, resp.Connected)
	require.Contains(t, resp.Error, infra.JDAccessTokenEnv)
}

func TestAPI_Health(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	server.handleHealth(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_StaticFiles(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	// Index HTML
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	server.handleStatic(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "CRE Reconciler")

	// CSS
	req2 := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/styles.css", nil)
	w2 := httptest.NewRecorder()
	server.handleStatic(w2, req2)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w2.Body.String(), "tab-active")

	// JS
	req3 := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/app.js", nil)
	w3 := httptest.NewRecorder()
	server.handleStatic(w3, req3)

	require.Equal(t, http.StatusOK, w3.Code)
	require.Contains(t, w3.Body.String(), "loadNodes")
}

func TestResponseToDesiredState(t *testing.T) {
	t.Parallel()

	resp := DesiredResponse{
		Infra: struct {
			Type        string `json:"type"`
			ChartValues string `json:"chartValues"`
			Namespace   string `json:"namespace"`
		}{Type: "griddle", ChartValues: "x", Namespace: "ns"},
		JD: struct {
			GRPC        string `json:"grpc"`
			Domain      string `json:"domain"`
			Environment string `json:"environment"`
			UseTLS      bool   `json:"useTLS"`
		}{GRPC: "g:443", Domain: "cre", Environment: "dev"},
		DONs: []DONResponse{
			{Name: "w", DONTypes: []string{"workflow"}, Capabilities: []string{"cron"}},
		},
		CapabilityConfigs: map[string]domain.CapabilityConfig{"cron": {BinaryName: "cron"}},
		GatewayNodes: []GatewayNodeResponse{
			{Node: "node-gw-0", DON: "workflow"},
		},
	}

	ds := responseToDesiredState(resp)
	require.Equal(t, "griddle", ds.Infra.Type)
	require.Equal(t, "w", ds.DONs[0].Name)
	require.Contains(t, ds.DONs[0].Capabilities, "cron")
	require.Equal(t, "cron", ds.CapabilityConfigs["cron"].BinaryName)
	require.Len(t, ds.GatewayNodes, 1)
	require.Equal(t, "node-gw-0", ds.GatewayNodes[0].Node)
	require.Equal(t, "workflow", ds.GatewayNodes[0].DON)
}

func TestNormalizeJSONNumberValue(t *testing.T) {
	t.Parallel()

	require.Equal(t, int64(1500000000), normalizeJSONNumberValue(float64(1500000000)))
	require.InDelta(t, 100.5, normalizeJSONNumberValue(100.5), 0.0001, "non-whole floats must be preserved")
	require.Equal(t, "evm", normalizeJSONNumberValue("evm"), "non-numeric values must be untouched")

	nested := map[string]any{"inner": float64(42)}
	require.Equal(t, map[string]any{"inner": int64(42)}, normalizeJSONNumberValue(nested))

	list := []any{float64(1), float64(2.5), "x"}
	require.Equal(t, []any{int64(1), 2.5, "x"}, normalizeJSONNumberValue(list))
}

func TestAPI_Desired_SaveAndLoad_NormalizesCapabilityConfigNumbers(t *testing.T) {
	t.Parallel()

	server, _ := setupTestServer(t)

	// LogTriggerPollInterval arrives as a JSON number with no decimal point,
	// which encoding/json still decodes to float64 for a map[string]any target
	// (the exact scenario that broke evm.go's `%d`-formatted job spec).
	payload := `{
		"infra": {"type": "griddle", "chartValues": "deploy/config/workflow", "namespace": "my-repo-nodeset"},
		"jd": {"grpc": "grpc-jd:443", "domain": "cre", "environment": "dev"},
		"chains": [
			{"chainId": 1337, "wsUrl": "wss://anvil-1337.example.com", "httpUrl": "https://anvil-1337.example.com", "registry": true}
		],
		"dons": [
			{
				"name": "workflow",
				"donTypes": ["workflow"],
				"capabilities": ["evm-1337"],
				"exposesRemoteCapabilities": false
			}
		],
		"capabilityConfigs": {
			"evm": {"binary_name": "evm", "values": {"LogTriggerPollInterval": 1500000000, "ReceiverGasMinimum": 100.5}}
		}
	}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/desired", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleDesired(w, req)
	require.Equal(t, http.StatusOK, w.Code, "POST failed: %s", w.Body.String())

	raw, err := os.ReadFile(server.desiredPath)
	require.NoError(t, err)
	require.Contains(t, string(raw), "LogTriggerPollInterval = 1500000000",
		"whole-number value must be persisted as a TOML integer, not a float (e.g. 1.5e+09)")
	require.Contains(t, string(raw), "ReceiverGasMinimum = 100.5",
		"non-whole value must still be persisted as a float")

	ds, err := domain.LoadDesiredState(server.desiredPath)
	require.NoError(t, err)
	values := ds.CapabilityConfigs["evm"].Values
	require.IsType(t, int64(0), values["LogTriggerPollInterval"])
	require.Equal(t, int64(1500000000), values["LogTriggerPollInterval"])
}
