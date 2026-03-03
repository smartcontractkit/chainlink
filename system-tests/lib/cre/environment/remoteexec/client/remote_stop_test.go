package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

type stubComponentClient struct {
	resp *agent.StartComponentResponse
	err  error
}

func (s *stubComponentClient) StartComponent(context.Context, agent.StartComponentEnvelope) (*agent.StartComponentResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.resp, nil
}

func TestCountRemoteStopTargets(t *testing.T) {
	cfg := &config.Config{
		Blockchains: []*config.Blockchain{
			{Input: blockchain.Input{}, Placement: config.PlacementRemote},
			{Input: blockchain.Input{}, Placement: config.PlacementLocal},
		},
		NodeSets: []*cre.NodeSet{
			{Input: &ns.Input{Name: "remote-don"}, Placement: "remote"},
			{Input: &ns.Input{Name: "local-don"}, Placement: "local"},
		},
		JD: &config.JobDistributor{Placement: config.PlacementRemote},
	}

	require.Equal(t, 3, countRemoteStopTargets(cfg), "expected only remote-targeted components to be counted")
}

func TestListRemoteCTFResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/resources/ctf", r.URL.Path)
		_, _ = w.Write([]byte(`{"containers":["c1","c2"],"volumes":["v1"]}`))
	}))
	defer server.Close()

	containers, volumes, err := listRemoteCTFResources(context.Background(), server.URL)
	require.NoError(t, err, "expected ctf resource listing to succeed")
	require.Equal(t, []string{"c1", "c2"}, containers)
	require.Equal(t, []string{"v1"}, volumes)
}

func TestListRemoteCTFResources_Non2xxFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("upstream down"))
	}))
	defer server.Close()

	_, _, err := listRemoteCTFResources(context.Background(), server.URL)
	require.Error(t, err, "expected non-2xx response to fail")
	require.Contains(t, err.Error(), "ctf resource query failed", "expected status failure context")
}

func TestStopRemoteComponents_SummaryAndResiduals(t *testing.T) {
	server := newRemoteStopTestServer(t)
	defer server.Close()

	t.Setenv(EnvRemoteAgentURL, server.URL)
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	cfg := &config.Config{
		Blockchains: []*config.Blockchain{
			{Input: blockchain.Input{}, Placement: config.PlacementRemote},
			{Input: blockchain.Input{}, Placement: config.PlacementLocal},
		},
		NodeSets: []*cre.NodeSet{
			{Input: &ns.Input{Name: "capabilities"}, Placement: "remote"},
			{Input: &ns.Input{Name: "local-don"}, Placement: "local"},
		},
		JD: &config.JobDistributor{Placement: config.PlacementRemote},
	}

	summary, err := StopRemoteComponents(context.Background(), zerolog.Nop(), cfg)
	require.NoError(t, err, "expected stop operation to succeed")
	require.Equal(t, 3, summary.Requested, "expected remote blockchain + remote nodeset + remote jd")
	require.Equal(t, 2, summary.Stopped, "expected blockchain and jd to be stopped")
	require.Equal(t, 1, summary.Missing, "expected nodeset to be missing")
	require.Equal(t, 0, summary.Failed, "expected no failed stop operations")
	require.Equal(t, []string{"leftover-container"}, summary.ResidualContainers)
	require.Equal(t, []string{"leftover-volume"}, summary.ResidualVolumes)
	require.Empty(t, summary.ResidualQueryError, "expected residual query to succeed")
}

func TestStopRemoteComponents_ResidualQueryFailureIsReportedInSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/health":
			w.WriteHeader(http.StatusOK)
		case "/v1/status":
			assert.NoError(t, json.NewEncoder(w).Encode(agent.AgentStatusResponse{ProtocolVersion: "1.0.0"}))
		case "/v1/components/start":
			resp := agent.StartComponentResponse{
				ComponentType: ComponentTypeBlockchain,
				Found:         true,
				Stopped:       true,
			}
			assert.NoError(t, json.NewEncoder(w).Encode(resp))
		case "/v1/resources/ctf":
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("ctf listing down"))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	t.Setenv(EnvRemoteAgentURL, server.URL)
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	cfg := &config.Config{
		Blockchains: []*config.Blockchain{
			{Input: blockchain.Input{}, Placement: config.PlacementRemote},
		},
	}

	summary, err := StopRemoteComponents(context.Background(), zerolog.Nop(), cfg)
	require.NoError(t, err, "stop should still succeed when residual listing fails")
	require.Equal(t, 1, summary.Requested)
	require.Equal(t, 1, summary.Stopped)
	require.NotEmpty(t, summary.ResidualQueryError, "expected residual query failure to be surfaced")
}

func TestStopRemoteComponent_UnexpectedComponentTypeFails(t *testing.T) {
	client := &stubComponentClient{
		resp: &agent.StartComponentResponse{
			ComponentType: ComponentTypeJD,
		},
	}

	_, err := stopRemoteComponent(
		context.Background(),
		zerolog.Nop(),
		client,
		agent.StartComponentPayload{ComponentType: ComponentTypeBlockchain},
		ComponentTypeBlockchain,
	)
	require.Error(t, err, "expected mismatched component type to fail")
	require.Contains(t, err.Error(), "unexpected component type")
}

func TestStopRemoteComponent_ClientErrorIsWrapped(t *testing.T) {
	client := &stubComponentClient{err: errors.New("network down")}

	_, err := stopRemoteComponent(
		context.Background(),
		zerolog.Nop(),
		client,
		agent.StartComponentPayload{ComponentType: ComponentTypeBlockchain},
		ComponentTypeBlockchain,
	)
	require.Error(t, err, "expected client failure to be returned")
	require.Contains(t, err.Error(), "failed to stop remote component type")
}

func newRemoteStopTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/health":
			w.WriteHeader(http.StatusOK)
			return
		case "/v1/status":
			assert.NoError(t, json.NewEncoder(w).Encode(agent.AgentStatusResponse{ProtocolVersion: "1.0.0"}))
			return
		case "/v1/resources/ctf":
			_, _ = w.Write([]byte(`{"containers":["leftover-container"],"volumes":["leftover-volume"]}`))
			return
		case "/v1/components/start":
			var envelope agent.StartComponentEnvelope
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&envelope))
			assert.Equal(t, agent.OperationStopComponent, envelope.Operation)

			var payload agent.StartComponentPayload
			assert.NoError(t, json.Unmarshal(envelope.Payload, &payload))

			resp := agent.StartComponentResponse{ComponentType: payload.ComponentType}
			switch payload.ComponentType {
			case ComponentTypeBlockchain:
				resp.Found = true
				resp.Stopped = true
			case ComponentTypeNodeSet:
				resp.Found = false
				resp.Stopped = false
			case ComponentTypeJD:
				resp.Found = true
				resp.Stopped = true
			default:
				t.Fatalf("unexpected component type %q", payload.ComponentType)
			}
			assert.NoError(t, json.NewEncoder(w).Encode(resp))
			return
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
}
