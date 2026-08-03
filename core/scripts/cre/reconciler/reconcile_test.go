package reconciler

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

func TestStateFile_PhaseDiscoveryRoundtrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := dir + "/state.toml"

	original := &domain.StateFile{
		Phase: domain.PhaseDiscovery,
		NodeRuntime: map[string]domain.NodeRuntimeInfo{
			"node-0": {APIURL: "https://node-0.example", CSAKey: "abc"},
		},
	}
	require.NoError(t, original.Store(path))

	loaded, err := domain.LoadState(path)
	require.NoError(t, err)
	require.Equal(t, domain.PhaseDiscovery, loaded.Phase)
	require.Equal(t, "https://node-0.example", loaded.NodeRuntime["node-0"].APIURL)
}

func TestOnChainComplete(t *testing.T) {
	t.Parallel()

	r := &Reconciler{state: &domain.StateFile{Phase: domain.PhaseDone}}
	require.False(t, r.onChainComplete())

	r.state.SetAddress(domain.AddressRef{Type: "CapabilitiesRegistry", Address: "0xabc"})
	require.False(t, r.onChainComplete())

	r.state.DONIDs = map[string]uint64{"workflow": 1}
	require.False(t, r.onChainComplete())

	r.state.WorkflowReg = &domain.WorkflowRegState{ChainSelector: 1}
	require.True(t, r.onChainComplete())
}

func TestAdvancePhaseAfterDiscovery(t *testing.T) {
	t.Parallel()

	fresh := &domain.StateFile{}
	advancePhaseAfterDiscovery(fresh)
	require.Equal(t, domain.PhaseDiscovery, fresh.Phase)

	again := &domain.StateFile{Phase: domain.PhaseDiscovery}
	advancePhaseAfterDiscovery(again)
	require.Equal(t, domain.PhaseDiscovery, again.Phase)

	later := &domain.StateFile{Phase: domain.PhaseTOML}
	advancePhaseAfterDiscovery(later)
	require.Equal(t, domain.PhaseTOML, later.Phase)
}

func TestRequireDiscoveryComplete(t *testing.T) {
	t.Parallel()

	r := &Reconciler{state: &domain.StateFile{}}
	require.Error(t, r.requireDiscoveryComplete())
	require.Contains(t, r.requireDiscoveryComplete().Error(), "discovery data missing")

	r.state.NodeRuntime = map[string]domain.NodeRuntimeInfo{}
	require.Error(t, r.requireDiscoveryComplete())

	r.state.NodeRuntime["node-0"] = domain.NodeRuntimeInfo{APIURL: "https://node-0.example"}
	require.NoError(t, r.requireDiscoveryComplete())
}

func TestRequireJDAccessToken(t *testing.T) {
	// Not parallel: mutates the process-wide GRIDDLE_JD_ACCESS_TOKEN env var via t.Setenv.

	t.Run("JD not configured — no token required", func(t *testing.T) {
		t.Setenv(infra.JDAccessTokenEnv, "")
		r := &Reconciler{desired: &domain.DesiredState{JD: domain.JDConfig{GRPC: ""}}}
		require.NoError(t, r.requireJDAccessToken())
	})

	t.Run("JD configured, token set — passes", func(t *testing.T) {
		t.Setenv(infra.JDAccessTokenEnv, "some-token")
		r := &Reconciler{desired: &domain.DesiredState{JD: domain.JDConfig{GRPC: "grpc-jd.example.com:443"}}}
		require.NoError(t, r.requireJDAccessToken())
	})

	t.Run("JD configured, token unset — fails with actionable message", func(t *testing.T) {
		t.Setenv(infra.JDAccessTokenEnv, "")
		r := &Reconciler{desired: &domain.DesiredState{JD: domain.JDConfig{GRPC: "grpc-jd.example.com:443"}}}
		err := r.requireJDAccessToken()
		require.Error(t, err)
		require.Contains(t, err.Error(), infra.JDAccessTokenEnv)
	})

	t.Run("JD configured, token set to whitespace only — treated as unset", func(t *testing.T) {
		t.Setenv(infra.JDAccessTokenEnv, "   ")
		r := &Reconciler{desired: &domain.DesiredState{JD: domain.JDConfig{GRPC: "grpc-jd.example.com:443"}}}
		require.Error(t, r.requireJDAccessToken())
	})
}

// TestRun_FailsFastWhenJDConfiguredWithoutAccessToken proves the token check runs
// BEFORE any other work in Run — including K8s discovery. state, cv, and k8s are
// deliberately left nil: if requireJDAccessToken were not the first thing Run does,
// runDiscoveryPhase would nil-pointer panic (on r.state.NodeRuntime) instead of this
// test observing a clean, actionable error.
func TestRun_FailsFastWhenJDConfiguredWithoutAccessToken(t *testing.T) {
	t.Setenv(infra.JDAccessTokenEnv, "")

	r := &Reconciler{
		desired: &domain.DesiredState{JD: domain.JDConfig{GRPC: "grpc-jd.example.com:443"}},
		log:     zerolog.Nop(),
	}

	err := r.Run(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), infra.JDAccessTokenEnv)
}

func TestRunDiscoveryPhase_FailureDoesNotAdvancePhase(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		state: &domain.StateFile{Phase: domain.PhaseNone},
		log:   zerolog.Nop(),
	}

	err := r.runDiscoveryPhase(context.Background())
	require.Error(t, err)
	require.Equal(t, domain.PhaseNone, r.state.Phase)
}

func TestRunDiscoveryPhase_SkipsWhenAlreadyComplete(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		state: &domain.StateFile{
			Phase: domain.PhaseDiscovery,
			NodeRuntime: map[string]domain.NodeRuntimeInfo{
				"node-0": {APIURL: "https://node-0.example"},
				"node-1": {APIURL: "https://node-1.example"},
			},
		},
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{
				{Name: "node-0"},
				{Name: "node-1"},
			},
		},
		log: zerolog.Nop(),
		// k8s nil — would fail if discover() ran
	}

	require.NoError(t, r.runDiscoveryPhase(context.Background()))
	require.Equal(t, domain.PhaseDiscovery, r.state.Phase)
}

func TestDiscoveryComplete(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		state: &domain.StateFile{NodeRuntime: map[string]domain.NodeRuntimeInfo{"node-0": {APIURL: "https://x"}}},
		cv:    &domain.ChartValues{Nodes: []domain.ChartNodeInfo{{Name: "node-0"}, {Name: "node-1"}}},
	}
	require.False(t, r.discoveryComplete())

	r.state.NodeRuntime["node-1"] = domain.NodeRuntimeInfo{APIURL: "https://y"}
	require.True(t, r.discoveryComplete())
}

func TestValidateDONNamesAgainstChart_HappyPath(t *testing.T) {
	t.Parallel()

	ds := &domain.DesiredState{
		DONs: []domain.DON{
			{Name: "zone-a-workflow"},
			{Name: "zone-a-gateway"},
		},
	}
	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", DONName: "zone-a-workflow"},
			{Name: "node-1", DONName: "zone-a-workflow"},
			{Name: "gateway-0", DONName: "zone-a-gateway"},
		},
	}

	require.NoError(t, validateDONNamesAgainstChart(ds, cv))
}

func TestValidateDONNamesAgainstChart_NameMismatch(t *testing.T) {
	t.Parallel()

	// DON membership is always chart-derived from don.Name now, so a "name
	// mismatch" means the DON's name simply doesn't match any chart node-set's
	// don-name label.
	ds := &domain.DesiredState{
		DONs: []domain.DON{
			{Name: "wf"},
		},
	}
	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", DONName: "zone-a-workflow"},
			{Name: "node-1", DONName: "zone-a-workflow"},
		},
	}

	err := validateDONNamesAgainstChart(ds, cv)
	require.Error(t, err)
	require.Contains(t, err.Error(), `no Griddle node-set has don-name "wf"`)
}

func TestValidateDONNamesAgainstChart_DuplicateNodeSetClaim(t *testing.T) {
	t.Parallel()

	ds := &domain.DesiredState{
		DONs: []domain.DON{
			{Name: "zone-a-workflow"},
			{Name: "zone-a-workflow"},
		},
	}
	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", DONName: "zone-a-workflow"},
		},
	}

	err := validateDONNamesAgainstChart(ds, cv)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is already claimed by dons entry")
}

func TestResolveGlobalBootstrapNodeName_BootstrapOnlyDON(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		desired: &domain.DesiredState{
			DONs: []domain.DON{
				{Name: "workflow"},
				{Name: "bootstrap", DONTypes: []string{"bootstrap"}, BootstrapNode: "node-bt-0"},
			},
		},
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{
				{Name: "node-0", NodeType: domain.RoleStandard, DONName: "workflow"},
				{Name: "node-1", NodeType: domain.RoleStandard, DONName: "workflow"},
				{Name: "node-bt-0", NodeType: domain.RoleBootstrap, DONName: "bootstrap"},
			},
		},
	}

	name, err := r.resolveGlobalBootstrapNodeName()
	require.NoError(t, err)
	require.Equal(t, "node-bt-0", name)
}

func TestResolveGlobalBootstrapForTOML_FromBootstrapOnlyDON(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		desired: &domain.DesiredState{
			DONs: []domain.DON{
				{Name: "workflow"},
				{Name: "bootstrap", DONTypes: []string{"bootstrap"}, BootstrapNode: "node-bt-0"},
			},
		},
		cv: &domain.ChartValues{
			Namespace: "zone-a",
			Nodes: []domain.ChartNodeInfo{
				{Name: "node-0", NodeType: domain.RoleStandard, DONName: "workflow"},
				{Name: "node-bt-0", NodeType: domain.RoleBootstrap, DONName: "bootstrap"},
			},
		},
		state: &domain.StateFile{
			NodeRuntime: map[string]domain.NodeRuntimeInfo{
				"node-bt-0": {PeerID: "16Uiu2HAbtv1rGFUqYhFvuk6z6b86P"},
			},
		},
	}

	info, err := r.resolveGlobalBootstrapForTOML()
	require.NoError(t, err)
	require.Equal(t, "node-bt-0", info.nodeName)
	require.Equal(t, "16Uiu2HAbtv1rGFUqYhFvuk6z6b86P", info.peerID)
	require.Equal(t, "node-bt-0.zone-a.svc.cluster.local", info.host)
}

func TestResolveGlobalBootstrapForTOML_BootstrapInDifferentNamespace(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		desired: &domain.DesiredState{
			DONs: []domain.DON{
				{Name: "zone-c-workflow"},
				{Name: "zone-c-bootstrap", DONTypes: []string{"bootstrap"}, BootstrapNode: "node-bt-1"},
			},
		},
		cv: &domain.ChartValues{
			Namespace: "zone-c", // primary namespace comes from the first-loaded node-set instance
			Nodes: []domain.ChartNodeInfo{
				{Name: "node-0", NodeType: domain.RoleStandard, Namespace: "zone-c", DONName: "zone-c-workflow"},
				{Name: "node-bt-1", NodeType: domain.RoleBootstrap, Namespace: "zone-c-bt", DONName: "zone-c-bootstrap"},
			},
		},
		state: &domain.StateFile{
			NodeRuntime: map[string]domain.NodeRuntimeInfo{
				"node-bt-1": {PeerID: "16Uiu2HAbtv1rGFUqYhFvuk6z6b86P"},
			},
		},
	}

	info, err := r.resolveGlobalBootstrapForTOML()
	require.NoError(t, err)
	require.Equal(t, "node-bt-1", info.nodeName)
	require.Equal(t, "node-bt-1.zone-c-bt.svc.cluster.local", info.host,
		"bootstrap host must use the bootstrap node's own node-set namespace, not the chart's primary namespace")
}

func TestResolveGlobalBootstrapForTOML_MissingPeerID(t *testing.T) {
	t.Parallel()

	r := &Reconciler{
		desired: &domain.DesiredState{
			DONs: []domain.DON{
				{Name: "bootstrap", DONTypes: []string{"bootstrap"}, BootstrapNode: "node-bt-0"},
			},
		},
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{{Name: "node-bt-0", NodeType: domain.RoleBootstrap, DONName: "bootstrap"}},
		},
		state: &domain.StateFile{
			NodeRuntime: map[string]domain.NodeRuntimeInfo{
				"node-bt-0": {},
			},
		},
	}

	_, err := r.resolveGlobalBootstrapForTOML()
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing peer_id")
}

// fakeK8s implements K8sAPI for tests.
type fakeK8s struct {
	restartedNodes []string
	restartErr     error
}

func (f *fakeK8s) GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*infra.NodeAPIInfo, error) {
	return nil, nil
}

func (f *fakeK8s) GetNodeSecretsToml(ctx context.Context, nodeName, namespace string) (string, error) {
	return "", nil
}

func (f *fakeK8s) RestartNodePods(ctx context.Context, nodeName, namespace string) error {
	if f.restartErr != nil {
		return f.restartErr
	}
	f.restartedNodes = append(f.restartedNodes, nodeName)
	return nil
}

func (f *fakeK8s) CopyFilesToPod(ctx context.Context, namespace, podName, container, destDir string, localPaths []string) error {
	return nil
}

func TestRestartWorkerNodePods_SkipsBootstrapAndGateway(t *testing.T) {
	t.Parallel()

	fake := &fakeK8s{}
	r := &Reconciler{
		log: zerolog.Nop(),
		k8s: fake,
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{
				{Name: "node-bt-0", NodeType: domain.RoleBootstrap},
				{Name: "node-gw-0", NodeType: domain.RoleGateway},
				{Name: "node-0", NodeType: domain.RoleStandard},
				{Name: "node-1", NodeType: domain.RoleStandard},
			},
		},
	}

	err := r.restartWorkerNodePods(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"node-0", "node-1"}, fake.restartedNodes)
}

func TestRestartWorkerNodePods_PropagatesError(t *testing.T) {
	t.Parallel()

	fake := &fakeK8s{restartErr: errors.New("boom")}
	r := &Reconciler{
		log: zerolog.Nop(),
		k8s: fake,
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{{Name: "node-0", NodeType: domain.RoleStandard}},
		},
	}

	err := r.restartWorkerNodePods(context.Background())
	require.Error(t, err)
}
